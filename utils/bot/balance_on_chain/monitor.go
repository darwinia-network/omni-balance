package balance_on_chain

import (
	"context"
	"omni-balance/utils/bot"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

func init() {
	boc := BalanceOnChain{}
	bot.Register(boc.Name(), boc)
}

type BalanceOnChain struct {
}

func (b BalanceOnChain) Name() string {
	return "balance_on_chain"
}

func (b BalanceOnChain) Check(ctx context.Context, args bot.Params) ([]bot.Task, bot.ProcessType, error) {
	var (
		config = args.Conf
		info   = args.Info
		tasks  []bot.Task
		log    = logrus.WithFields(logrus.Fields{
			"wallet":    args.Info.Wallet.GetAddress(),
			"chain":     args.Info.Chain,
			"tokenName": args.Info.TokenName,
			"job":       b.Name(),
		})
	)
	token := config.GetTokenInfoOnChain(info.TokenName, args.Info.Chain)
	balance, err := info.Wallet.GetExternalBalance(ctx, common.HexToAddress(token.ContractAddress), token.Decimals, args.Client)
	if err != nil {
		return nil, "", errors.Wrap(err, "get balance error")
	}
	threshold := config.GetTokenThreshold(info.Wallet.GetAddress().Hex(), info.TokenName, info.Chain)
	if balance.GreaterThan(threshold) {
		log.Debugf("balance on chain is enough, balance: %s;threshold: %s", balance.String(), threshold.String())
		return nil, "", nil
	}
	amount := config.GetTokenPurchaseAmount(info.Wallet.GetAddress().Hex(), info.TokenName, info.Chain)
	if balance.Add(amount).LessThanOrEqual(threshold) {
		newAmount := threshold.Add(threshold.Mul(decimal.RequireFromString("0.3")))
		log.Infof("The %s current balance is %s, amount in config is %s, balance(%s) + amount(%s) <= threshold(%s), so set amount to %s",
			info.Wallet.GetAddress(), balance, amount, balance, amount, threshold, newAmount)
		amount = newAmount
	}
	log.Debugf("balance on chain is not enough, balance: %s;threshold: %s", balance.String(), threshold.String())
	tasks = append(tasks, bot.Task{
		TokenOutChainName: args.Info.Chain,
		Wallet:            info.Wallet.GetAddress().Hex(),
		TokenOutName:      info.TokenName,
		Amount:            amount,
	})
	return tasks, bot.Queue, nil
}
