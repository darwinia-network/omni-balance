package gate_liquidity

import (
	"context"
	"omni-balance/utils/bot"
	"omni-balance/utils/bot/balance_on_chain"
	"omni-balance/utils/configs"

	"github.com/shopspring/decimal"
	"github.com/spf13/cast"
)

func init() {
	bot.Register(GateLiquidity{}.Name(), GateLiquidity{})
}

type GateLiquidity struct {
	Bot bot.Bot
}

func (g GateLiquidity) Name() string {
	return "gate_liquidity"
}

func (b GateLiquidity) Balance(ctx context.Context, args bot.Params) (decimal.Decimal, error) {
	return balance_on_chain.BalanceOnChain{}.Balance(ctx, args)
}

func (b GateLiquidity) Check(ctx context.Context, args bot.Params) ([]bot.Task, bot.ProcessType, error) {
	if b.Bot == nil {
		b.Bot = balance_on_chain.BalanceOnChain{}
	}
	tasks, processType, err := b.Bot.Check(ctx, args)
	if err != nil {
		return nil, processType, err
	}
	var result []bot.Task
	for _, task := range tasks {
		botConf := args.Conf.GetBotConfig(task.Wallet, task.TokenOutChainName, task.TokenOutName, b.Name())
		if len(botConf) == 0 {
			panic("gate liquidity bot config not found")
		}
		result = append(result, bot.Task{
			Wallet:            task.Wallet,
			TokenInName:       args.Info.TokenName,
			TokenOutName:      args.Info.TokenName,
			TokenInChainName:  args.Info.Chain,
			TokenOutChainName: cast.ToString(botConf["toChain"]),
			Amount:            task.Amount,
			ProviderType:      configs.CEX,
			ProviderName:      "gate.io",
			Remark:            b.Name(),
		})
	}
	return result, processType, nil
}
