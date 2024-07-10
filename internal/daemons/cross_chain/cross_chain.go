package cross_chain

import (
	"context"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"omni-balance/internal/daemons"
	"omni-balance/internal/db"
	"omni-balance/internal/models"
	"omni-balance/utils"
	"omni-balance/utils/configs"
	"omni-balance/utils/provider"
	"sync"
	"time"
)

func init() {
	daemons.RegisterIntervalTask(daemons.Task{
		Name:            "crossChain",
		TaskFunc:        Run,
		DefaultInterval: time.Second * 3,
		Description:     "Responsible for cross-chaining unfinished tokens from the Rebalance task to the target chain.",
	})
}

func Run(ctx context.Context, conf configs.Config) error {
	orders, err := retrieveOrders()
	if err != nil {
		return err
	}
	if len(orders) == 0 {
		return nil
	}
	return processOrders(ctx, conf, orders)
}

func retrieveOrders() ([]*models.Order, error) {
	var orders []*models.Order
	return orders, db.DB().Where("status = ? AND current_chain_name != target_chain_name",
		models.OrderStatusWaitCrossChain).Find(&orders).Error
}

func processOrders(ctx context.Context, conf configs.Config, orders []*models.Order) error {
	var w sync.WaitGroup
	for _, order := range orders {
		w.Add(1)

		go func(order *models.Order) {
			defer utils.Recover()
			defer w.Done()
			log := order.GetLogs()
			if err := start(utils.SetLogToCtx(ctx, log), order, conf, log.WithField("order_id", order.ID)); err != nil {
				log.Errorf("start cross chain error: %s", err)
			}
		}(order)
	}
	w.Wait()
	return nil
}

func start(ctx context.Context, order *models.Order, conf configs.Config, log *logrus.Entry) error {
	if order.Lock(db.DB()) {
		log.Infof("order #%d locked, unlock time is %s", order.ID, time.Unix(order.LockTime+60*60*1, 0))
		return nil
	}
	defer order.UnLock(db.DB())
	if order == nil {
		return errors.New("order is nil")
	}
	wallet := conf.GetWallet(order.Wallet)
	bridge, err := getBridge(ctx, order, conf)
	if err != nil {
		return errors.Wrap(err, "get bridge error")
	}
	orderProcess := models.GetLastOrderProcess(ctx, db.DB(), order.ID)
	swapParams := daemons.CreateSwapParams(*order, orderProcess, log, wallet)
	if order.CurrentChainName != "" && order.CurrentChainName != swapParams.SourceChain {
		swapParams.SourceToken = order.CurrentChainName
	}
	result, err := bridge.Swap(utils.SetLogToCtx(ctx, order.GetLogs()), swapParams)
	if err != nil {
		return errors.Wrap(err, "swap error")
	}
	update := createUpdateLog(order, result, log)
	return db.DB().Model(&models.Order{}).Where("id = ?", order.ID).Updates(update).Error
}

func createUpdateLog(order *models.Order, result provider.SwapResult, log *logrus.Entry) map[string]interface{} {
	update := map[string]interface{}{
		"error":              result.Error,
		"current_chain_name": result.CurrentChain,
		"status":             result.Status,
	}
	if result.Status == "" {
		update["status"] = provider.TxStatusFailed
	}
	log.Infof("order #%d token %s cross from %s to %s status is %s", order.ID, order.TokenOutName,
		order.CurrentChainName, order.TargetChainName, result.Status)
	return update
}

func getBridge(ctx context.Context, order *models.Order, conf configs.Config) (provider.Provider, error) {
	var (
		bridges []provider.Provider
		log     = order.GetLogs()
	)
	for _, providerInitFunc := range provider.LiquidityProviderTypeAndConf(configs.Bridge, conf) {
		bridge, err := provider.InitializeBridge(providerInitFunc, conf)
		if err != nil {
			return nil, errors.Wrap(err, "init bridge error")
		}
		tokenInCosts, err := bridge.GetCost(context.Background(), provider.SwapParams{
			SourceToken: order.TokenInName,
			Sender:      conf.GetWallet(order.Wallet),
			TargetToken: order.TokenOutName,
			Receiver:    order.Wallet,
			TargetChain: order.TargetChainName,
			Amount:      order.Amount,
		})
		if err != nil {
			log.Warnf("check bridge error: %v, not use %s bridge", err, bridge.Name())
			continue
		}
		if len(tokenInCosts) == 0 {
			continue
		}
		log.Debugf("check bridge %s success", bridge.Name())
		bridges = append(bridges, bridge)
	}
	if len(bridges) == 0 {
		return nil, errors.New("no bridge found")
	}
	return bridges[0], nil
}
