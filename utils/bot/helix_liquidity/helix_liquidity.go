package helix_liquidity

import (
	"context"
	"omni-balance/utils/bot"
	"omni-balance/utils/chains"
	"omni-balance/utils/provider"
	"omni-balance/utils/provider/bridge/helix_liquidity_claim"

	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

func init() {
	bot.Register(HelixLiquidity{}.Name(), HelixLiquidity{})
}

type HelixLiquidity struct {
}

func (h HelixLiquidity) Check(ctx context.Context, args bot.Params) ([]bot.Task, bot.ProcessType, error) {
	log := logrus.WithFields(logrus.Fields{
		"bot":    h.Name(),
		"wallet": args.Info.Wallet.GetAddress().Hex(),
		"chain":  args.Info.Chain,
		"token":  args.Info.TokenName,
	})
	hlc, err := helix_liquidity_claim.New(args.Conf)
	if err != nil {
		return nil, bot.Parallel, err
	}
	claim := hlc.(helix_liquidity_claim.Claim)
	log.Debugf("start check need withdraw helix records")
	records, err := claim.ListNeedWithdrawRecords(ctx, args.Info.Wallet.GetAddress(), args.Info.Chain, args.Info.TokenName)
	if err != nil {
		return nil, bot.Parallel, err
	}
	token := args.Conf.GetTokenInfoOnChain(args.Info.TokenName, args.Info.Chain)
	client, err := chains.NewTryClient(ctx, args.Conf.GetChainConfig(args.Info.Chain).RpcEndpoints)
	if err != nil {
		return nil, bot.Parallel, err
	}
	defer client.Close()
	balance, err := args.Info.Wallet.GetExternalBalance(ctx, common.HexToAddress(token.ContractAddress), token.Decimals, client)
	if err != nil {
		return nil, bot.Parallel, err
	}
	log.Debugf("balance of %s is %s on chain", args.Info.Wallet.GetAddress(), balance)
	threshold := args.Conf.GetTokenThreshold(args.Info.Wallet.GetAddress().Hex(), args.Info.TokenName, args.Info.Chain)
	purchaseAmount := args.Conf.GetTokenPurchaseAmount(args.Info.Wallet.GetAddress().Hex(), args.Info.TokenName, args.Info.Chain)
	var (
		result []bot.Task
		total  = balance
	)
	for _, v := range records {
		total = total.Add(v.TotalAmount)
	}
	log.Debugf("total balance of %s is %s on chain plus need withdraw helix", args.Info.Wallet.GetAddress(), total)
	for _, v := range debtImpl {
		balance, err := v.BalanceOf(ctx, DebtParams{
			Address: args.Info.Wallet.GetAddress(),
			Chain:   args.Info.Chain,
			Client:  client,
			Token:   token.Name,
		})
		if err != nil {
			log.Warnf("get balance of %s error: %s", v.Name(), err)
			continue
		}
		total = total.Add(balance)
	}

	log.Debugf("total balance of %s is %s on chain plus need withdraw helix plus debt", args.Info.Wallet.GetAddress(), total)
	if !total.GreaterThan(threshold) {
		log.Debugf("total balance of %s is %s on chain plus need withdraw helix plus debt less than threshold, skip", args.Info.Wallet.GetAddress(), total)
		return nil, bot.Parallel, nil
	}

	if total.Add(purchaseAmount).LessThanOrEqual(threshold) {
		newAmount := threshold.Add(threshold.Mul(decimal.RequireFromString("0.3")))
		log.Infof("The %s current balance is %s, amount in config is %s, balance(%s) + amount(%s) <= threshold(%s), so set amount to %s",
			args.Info.Wallet.GetAddress(), total, purchaseAmount, balance, purchaseAmount, threshold, newAmount)
		purchaseAmount = newAmount
	}

	result = append(result, bot.Task{
		Wallet:            args.Info.Wallet.GetAddress().Hex(),
		TokenOutName:      args.Info.TokenName,
		TokenOutChainName: args.Info.Chain,
		Amount:            purchaseAmount,
		Status:            provider.TxStatusPending,
	})
	return result, bot.Queue, nil
}

func (h HelixLiquidity) Name() string {
	return "helix_liquidity"
}
