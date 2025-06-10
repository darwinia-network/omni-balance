package daemons

import (
	"context"
	"fmt"
	"omni-balance/internal/db"
	"omni-balance/internal/models"
	"omni-balance/utils/chains"
	"omni-balance/utils/configs"
	"omni-balance/utils/constant"
	"omni-balance/utils/provider"
	"omni-balance/utils/wallets"
	"strings"
	"sync"

	log "omni-balance/utils/logging"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/schollz/progressbar/v3"
	"github.com/shopspring/decimal"
)

func CreateSwapParams(order models.Order, orderProcess models.OrderProcess, wallet wallets.Wallets) provider.SwapParams {
	sourceChain := order.CurrentChainName
	if order.SourceChainName != "" {
		sourceChain = order.SourceChainName
	}
	return provider.SwapParams{
		OrderId:          order.ID,
		SourceChain:      sourceChain,
		SourceChainNames: order.TokenInChainNames,
		Sender:           wallet,
		Receiver:         order.Wallet,
		TargetChain:      order.TargetChainName,
		SourceToken:      order.TokenInName,
		TargetToken:      order.TokenOutName,
		Amount:           order.Amount,
		LastHistory:      createLastHistory(orderProcess),
		RecordFn:         createRecordFunction(order),
		SaveOrderFn:      createSaveOrderFunction(order),
		Order:            order.Order,
		Remark:           order.Remark,
		CurrentBalance:   order.CurrentBalance,
	}
}

func createLastHistory(orderProcess models.OrderProcess) provider.SwapHistory {
	return provider.SwapHistory{
		Actions:      orderProcess.Action,
		Status:       orderProcess.Status,
		CurrentChain: orderProcess.CurrentChainName,
		Amount:       orderProcess.Amount,
		Tx:           orderProcess.Tx,
	}
}

func createSaveOrderFunction(order models.Order) func(update map[string]interface{}) {
	return func(update map[string]interface{}) {
		db.DB().Model(&models.Order{}).Where("id = ?", order.ID).Updates(update)
	}
}

func createRecordFunction(order models.Order) func(s provider.SwapHistory, errs ...error) {
	return func(s provider.SwapHistory, errs ...error) {
		op := createOrderProcess(order, s)
		if len(errs) != 0 && errs[0] != nil {
			op.Error = errs[0].Error()
		}
		if err := db.DB().Create(op).Error; err != nil {
			log.Errorf("save %s bridge provider error: %s", order.TokenOutName, err.Error())
			return
		}
		if op.Error != "" {
			log.Errorf("#%d action %s status is %s error: %s",
				order.ID, s.Actions, s.Status, op.Error)
		} else {
			log.Infof("#%d action %s status is %s",
				order.ID, s.Actions, s.Status)
		}
	}
}

func createOrderProcess(order models.Order, s provider.SwapHistory) *models.OrderProcess {
	var (
		providerType = order.ProviderType
		providerName = order.ProviderName
	)
	if s.ProviderType != "" {
		providerType = configs.ProviderType(s.ProviderType)
	}
	if s.ProviderName != "" {
		providerName = s.ProviderName
	}
	return &models.OrderProcess{
		OrderId:          order.ID,
		ProviderType:     providerType,
		ProviderName:     providerName,
		Status:           s.Status,
		Action:           s.Actions,
		Amount:           s.Amount,
		CurrentChainName: s.CurrentChain,
		Tx:               s.Tx,
	}
}

func FindAllChainBalance(_ context.Context, confPath string, needPrintProgress bool) (map[string]map[string]map[string]decimal.Decimal, error) {
	conf := new(configs.Config)
	if err := cleanenv.ReadConfig(confPath, conf); err != nil {
		return nil, err
	}
	conf = conf.Init()
	tokens := make(map[string]struct{})
	type task struct {
		tokenName     string
		chainName     string
		walletAddress string
	}

	type value struct {
		balance decimal.Decimal
		remarks string
	}

	var (
		chainNames         []string
		chainQueueMap      = make(map[string]chan task)
		tokenBalanceResult sync.Map
		w                  sync.WaitGroup
		bar                *progressbar.ProgressBar
	)

	getKey := func(chainName, tokenname, walletAddress string) string {
		return fmt.Sprintf("%s:%s:%s", chainName, tokenname, walletAddress)
	}
	parseKey := func(key string) (chainName string, tokenname string, walletAddress string) {
		split := strings.Split(key, ":")
		return split[0], split[1], split[2]
	}

	for _, v := range conf.Chains {
		w.Add(1)
		for _, token := range v.Tokens {
			tokens[token.Name] = struct{}{}
		}
		chainNames = append(chainNames, v.Name)
		chainQueueMap[v.Name] = make(chan task, 5)

		go func(chainName string, queue chan task) {
			defer w.Done()
			chain := conf.GetChainConfig(chainName)
			client, err := chains.NewTryClient(context.Background(), chain.RpcEndpoints)
			if err != nil {
				panic(err)
			}
			for v := range queue {
				token := conf.GetTokenInfoOnChainNil(v.tokenName, chainName)
				key := getKey(chainName, v.tokenName, v.walletAddress)
				if token.Name == "" {
					tokenBalanceResult.Store(key, value{balance: decimal.RequireFromString("-1"), remarks: "not found"})
					continue
				}
				var errorCount int64
				for {
					balance, err := chains.GetTokenBalance(context.Background(), client, token.ContractAddress, v.walletAddress, token.Decimals)
					if err != nil {
						errorCount++
						if errorCount > 3 {
							log.Fatalf("get %s balance on %s error: %s", v.tokenName, chainName, err)
							return
						}
						continue
					}
					errorCount = 0
					tokenBalanceResult.Store(key, value{balance: balance})
					if bar != nil {
						_ = bar.Add(1)
					}
					break
				}
			}
		}(v.Name, chainQueueMap[v.Name])
	}

	var tasks []task
	for _, chainName := range chainNames {
		for tokenName := range tokens {
			for _, v := range conf.Wallets {
				tasks = append(tasks, task{
					tokenName:     tokenName,
					chainName:     chainName,
					walletAddress: v.Address,
				})
				if v.Operator.Address.Cmp(constant.ZeroAddress) != 0 {
					tasks = append(tasks, task{
						tokenName:     tokenName,
						chainName:     chainName,
						walletAddress: v.Operator.Address.Hex(),
					})
				}

				if v.Operator.Operator.Cmp(constant.ZeroAddress) != 0 {
					tasks = append(tasks, task{
						tokenName:     tokenName,
						chainName:     chainName,
						walletAddress: v.Operator.Operator.Hex(),
					})
				}
			}
		}
	}
	if needPrintProgress {
		bar = progressbar.New(len(tasks))
	}
	for _, v := range tasks {
		chainQueueMap[v.chainName] <- v
	}
	for k := range chainQueueMap {
		close(chainQueueMap[k])
	}
	w.Wait()

	var result = make(map[string]map[string]map[string]decimal.Decimal)

	tokenBalanceResult.Range(func(key, v any) bool {
		chainName, tokenName, walletAddress := parseKey(key.(string))
		if result[chainName] == nil {
			result[chainName] = make(map[string]map[string]decimal.Decimal)
		}
		if result[chainName][tokenName] == nil {
			result[chainName][tokenName] = make(map[string]decimal.Decimal)
		}
		result[chainName][tokenName][walletAddress] = v.(value).balance
		return true
	})
	return result, nil
}
