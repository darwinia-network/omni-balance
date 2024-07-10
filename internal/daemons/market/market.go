package market

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"omni-balance/internal/daemons"
	"omni-balance/internal/db"
	"omni-balance/internal/models"
	"omni-balance/utils"
	"omni-balance/utils/bot"
	"omni-balance/utils/chains"
	"omni-balance/utils/configs"
	"omni-balance/utils/notice"
	"omni-balance/utils/provider"
	"sync"
	"sync/atomic"
	"time"
)

var (
	hasRunQueue  = atomic.Bool{}
	processTasks sync.Map
)

func Run(ctx context.Context, conf configs.Config) error {
	if !hasRunQueue.Load() {
		go runFromQueue(ctx, conf)
	}
	hasRunQueue.Store(true)
	tasks, err := models.ListNotSuccessTasks(ctx, db.DB(), func(order models.Order) bool {
		if order.HasLocked() {
			return false
		}
		return !utils.InArrayFold(order.Status.String(), []string{
			models.OrderStatusWaitTransferFromOperator.String(),
			models.OrderStatusWaitCrossChain.String(),
		})
	})
	if err != nil {
		return errors.Wrap(err, "list tasks error")
	}
	for index := range tasks {
		if _, ok := processTasks.Load(tasks[index].Id); ok {
			continue
		}
		PushTask(Task{
			Id:          tasks[index].Id,
			ProcessType: bot.ProcessType(tasks[index].ProviderType),
		})
	}
	return nil
}

func runFromQueue(ctx context.Context, conf configs.Config) {
	defer func() {
		if err := recover(); err != nil {
			logrus.Errorf("market job error: %v, restart after 5 second", err)
			time.Sleep(time.Second * 5)
			runFromQueue(ctx, conf)
		}
	}()
	for {
		select {
		case <-ctx.Done():
			return

		case task := <-taskQueue:
			if _, ok := processTasks.Load(task.Id); ok {
				return
			}
			orders, err := models.ListOrdersByTaskId(ctx, db.DB(), task.Id)
			if err != nil {
				logrus.Errorf("get orders by task id error: %s", err.Error())
				continue
			}
			if len(orders) == 0 {
				continue
			}
			logrus.Debugf("Start task %s. There are %d orders", task.Id, len(orders))
			fn := func(orders []models.Order, task Task) func() {
				processTasks.Store(task.Id, struct{}{})
				return func() {
					defer func() { processTasks.Delete(task.Id) }()

					var taskWait sync.WaitGroup
					for index := range orders {
						if orders[index].HasLocked() {
							continue
						}
						if task.ProcessType == bot.Parallel {
							taskWait.Add(1)
							go func(order models.Order) {
								taskWait.Done()
								do(ctx, order, conf)
							}(orders[index])
							continue
						}
						do(ctx, orders[index], conf)
					}
					taskWait.Wait()
					logrus.Infof("Finish all orders in task %s", task.Id)
				}
			}
			utils.Go(fn(orders, task))
		}
	}
}

func do(ctx context.Context, order models.Order, conf configs.Config) {
	defer utils.Recover()
	log := order.GetLogs()
	subCtx, cancel := context.WithCancel(utils.SetLogToCtx(ctx, log))
	defer cancel()
	utils.Go(func() {
		defer cancel()
		var t = time.NewTicker(time.Second * 5)
		defer t.Stop()
		for {
			select {
			case <-subCtx.Done():
				return
			case <-t.C:
				var count int64
				_ = db.DB().Model(&models.Order{}).Where("id = ?", order.ID).Count(&count)
				if count == 0 {
					log.Infof("order #%d not found, exit this order rebalance", order.ID)
					return
				}
			}
		}
	})

	err := processOrder(subCtx, order, conf)
	if errors.Is(err, context.Canceled) || errors.Is(ctx.Err(), context.Canceled) || errors.Is(subCtx.Err(), context.Canceled) {
		return
	}
	if err != nil {
		if addOrderError(order.ID) >= 10 {
			_ = notice.Send(
				provider.WithNotify(ctx, provider.WithNotifyParams{
					TaskId:        order.TaskId,
					OrderId:       order.ID,
					Receiver:      common.HexToAddress(order.Wallet),
					TokenOut:      order.TokenOutName,
					TokenOutChain: order.TargetChainName,
				}),
				"process order error",
				fmt.Sprintf("order #%d error: %s, please check the provider configuration then restart the application.", order.ID, err),
				logrus.ErrorLevel,
			)
			return
		}
		log.Errorf(" order #%d error: %s", order.ID, err)
		return
	}
	removeOrderError(order.ID)
	order = models.GetOrder(ctx, db.DB(), order.ID)
	err = notice.Send(
		provider.WithNotify(ctx, provider.WithNotifyParams{
			TaskId:         order.TaskId,
			OrderId:        order.ID,
			Receiver:       common.HexToAddress(order.Wallet),
			CurrentBalance: order.CurrentBalance,
		}),
		fmt.Sprintf("rebalance %s on %s success", order.TokenOutName, order.TargetChainName),
		fmt.Sprintf("rebalance %s %s from %s to %s use %s %s",
			order.TokenOutName, order.Amount, order.SourceChainName, order.TargetChainName,
			order.ProviderName, order.ProviderType),
		logrus.InfoLevel,
	)
	if err != nil {
		log.Debugf("notice error: %s", err)
	}
	log.Infof(" order #%d success", order.ID)
}

func processOrder(ctx context.Context, order models.Order, conf configs.Config) error {
	log := utils.GetLogFromCtx(ctx)
	if order.Lock(db.DB()) {
		return errors.Errorf("order #%d locked, unlock time is %s", order.ID, time.Unix(order.LockTime+60*60*1, 0))
	}
	defer order.UnLock(db.DB())
	var (
		orderProcess = models.GetLastOrderProcess(ctx, db.DB(), order.ID)
		args         = daemons.CreateSwapParams(order, orderProcess, log, conf.GetWallet(order.Wallet))
		wallet       = conf.GetWallet(order.Wallet)
		token        = conf.GetTokenInfoOnChain(order.TokenOutName, order.TargetChainName)
		chain        = conf.GetChainConfig(order.TargetChainName)
		client, err  = chains.NewTryClient(ctx, chain.RpcEndpoints)
	)

	if err != nil {
		return errors.Wrap(err, "new evm client error")
	}
	defer client.Close()
	if wallet.IsDifferentAddress() && order.Status == models.OrderStatusWaitTransferFromOperator {
		ok, err := transfer(ctx, order, args, conf, client)
		if err != nil && ok {
			return errors.Wrap(err, "transfer error")
		}
		if ok {
			return nil
		}
		log.Infof("cannot use transfer, try other providers.")
	}

	balance, err := wallet.GetExternalBalance(ctx, common.HexToAddress(token.ContractAddress), token.Decimals, client)
	if err != nil {
		return errors.Wrap(err, "check balance error")
	}

	for _, v := range conf.GetWalletConfig(order.Wallet).Tokens {
		if !utils.InArray(order.TargetChainName, v.Chains) {
			continue
		}
		if order.TokenOutName != v.Name {
			continue
		}
		if !balance.GreaterThan(balance) {
			break
		}
		log.Infof("%s balance on %s is enough, skip", v.Name, order.TargetChainName)
		if err := order.Success(db.DB(), "", nil, balance); err != nil {
			return errors.Wrap(err, "update order success error")
		}
		return nil
	}

	providerObj, err := getBestProvider(ctx, order, conf)
	if err != nil {
		return errors.Wrap(err, "get  provider error")
	}

	if err := order.SaveProvider(db.DB(), providerObj.Type(), providerObj.Name()); err != nil {
		return errors.Wrap(err, "save provider error")
	}

	log.Infof("start  #%d %s on %s use %s provider", order.ID, order.TokenOutName,
		order.TargetChainName, providerObj.Name())
	result, providerErr := providerObj.Swap(ctx, args)
	if errors.Is(providerErr, context.Canceled) {
		return nil
	}
	if result.Status == "" {
		return errors.New("the result status is empty")
	}
	if result.CurrentChain != args.TargetChain && providerErr == nil {
		result.Status = models.OrderStatusWaitCrossChain
	}
	if err := createUpdateLog(ctx, order, result, conf, client); err != nil {
		return errors.Wrap(err, "create update log error")
	}
	if providerErr != nil {
		return errors.Wrap(providerErr, "provider error")
	}
	if args.Receiver != result.Receiver && result.Receiver != "" {
		order = models.GetOrder(ctx, db.DB(), order.ID)
		if order.ID == 0 {
			return errors.New("order not found")
		}
		_, err = transfer(ctx, order, daemons.CreateSwapParams(order, orderProcess, log, conf.GetWallet(order.Wallet)), conf, client)
		if err != nil {
			return errors.Wrap(err, "transfer error")
		}
	}
	return nil
}
