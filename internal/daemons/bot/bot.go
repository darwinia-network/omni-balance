package bot

import (
	"context"
	"fmt"
	"omni-balance/internal/daemons/market"
	"omni-balance/utils/bot"
	"omni-balance/utils/chains"
	"omni-balance/utils/configs"
	"sync"

	"github.com/ethereum/go-ethereum/ethclient/simulated"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func Run(ctx context.Context, conf configs.Config) error {
	existBuyTokens, err := getExistingBuyTokens()
	if err != nil {
		return errors.Wrap(err, "find buy tokens error")
	}
	var (
		clients = make(map[string]simulated.Client)
	)
	for _, chain := range conf.Chains {
		client, err := chains.NewTryClient(ctx, chain.RpcEndpoints)
		if err != nil {
			return errors.Wrap(err, "create client error")
		}
		clients[chain.Name] = client
	}
	defer func() {
		for _, client := range clients {
			client.(*chains.Client).Close()
		}
	}()

	var (
		ignoreTokens = createIgnoreTokens(existBuyTokens)
		w            sync.WaitGroup
	)

	for _, wallet := range conf.Wallets {
		for _, token := range wallet.Tokens {
			for _, chainName := range token.Chains {
				if ignoreTokens.Contains(token.Name, chainName, wallet.Address) {
					logrus.Debugf("ignore token %s on chain %s", token.Name, chainName)
					continue
				}
				var botTypes = conf.ListBotNames(wallet.Address, chainName, token.Name)

				if len(botTypes) == 0 {
					botTypes = append(botTypes, "balance_on_chain")
				}

				for _, botType := range botTypes {
					logrus.Debugf("start check %s on %s use %s bot", token.Name, chainName, botType)

					w.Add(1)
					go func(f func() ([]bot.Task, bot.ProcessType, error), botType string) {
						defer w.Done()
						tasks, processType, err := f()
						if err != nil {
							logrus.Errorf("bot error: %s", err)
							return
						}
						if len(tasks) == 0 {
							return
						}
						orders, taskId, err := createOrder(ctx, tasks, processType)
						if err != nil {
							logrus.Errorf("create order error: %s", err)
							return
						}
						if len(orders) == 0 {
							return
						}
						logrus.Infof("create %d tasks, based %s on %s using %s bot", len(tasks), tasks[0].TokenOutName,
							tasks[0].TokenOutChainName, botType)
						market.PushTask(market.Task{
							Id:          taskId,
							ProcessType: processType,
						})
					}(process(ctx, conf, wallet.Address, token.Name, chainName, botType, clients[chainName]), botType)
				}
			}
		}
	}
	w.Wait()
	return nil
}

func process(ctx context.Context, conf configs.Config, walletAddress, tokenName, chainName, botType string,
	client simulated.Client) func() ([]bot.Task, bot.ProcessType, error) {
	m := bot.GetMonitor(botType)
	if m == nil {
		panic(fmt.Sprintf("%s botType not found", botType))
	}
	return func() ([]bot.Task, bot.ProcessType, error) {
		return m.Check(ctx, bot.Params{
			Conf: conf,
			Info: bot.Config{
				Wallet:    conf.GetWallet(walletAddress),
				TokenName: tokenName,
				Chain:     chainName,
			},
			Client: client,
		})
	}
}
