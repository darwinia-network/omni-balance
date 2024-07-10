package helix_liquidity

import (
	"context"
	"omni-balance/utils/bot"
	"omni-balance/utils/provider"
	"omni-balance/utils/provider/bridge/helix_liquidity_claim"

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
	var result []bot.Task
	for index, v := range records {
		result = append(result, bot.Task{
			Remark:            v.Channel,
			Wallet:            args.Info.Wallet.GetAddress().Hex(),
			CurrentChainName:  v.FromChain,
			TokenInName:       v.TokenName,
			TokenOutName:      v.TokenName,
			TokenInChainName:  v.FromChain,
			TokenOutChainName: v.ToChain,
			Amount:            v.TotalAmount,
			Status:            provider.TxStatusPending,
			ProviderType:      claim.Type(),
			ProviderName:      claim.Name(),
			Order:             records[index],
		})
	}
	return result, bot.Queue, nil
}

func (h HelixLiquidity) Name() string {
	return "helix_liquidity"
}
