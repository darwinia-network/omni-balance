package gate_liquidity

import (
	"context"
	"omni-balance/utils/bot"
	"omni-balance/utils/bot_mocks"
	"omni-balance/utils/configs"
	"omni-balance/utils/constant"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGateLiquidity_Check(t *testing.T) {
	testConf := &configs.Config{
		Chains: []configs.Chain{
			{
				Id:          1,
				Name:        constant.Ethereum,
				NativeToken: "ETH",
				Tokens: []configs.Token{
					{
						ContractAddress: "0x0000000000000000000000000000000000000000",
						Decimals:        18,
						Name:            "ETH",
					},
				},
			},
		},
		Wallets: []configs.Wallet{
			{
				Address: constant.ZeroAddress.Hex(),
				BotTypes: []configs.BotConfig{
					{
						Name: "gate_liquidity",
						TokenChains: map[string][]string{
							"ETH": {constant.Ethereum},
						},
						Config: map[string]interface{}{
							"toChain": constant.Arbitrum,
						},
					},
				},
				Tokens: []configs.WalletToken{
					{

						Name:      "ETH",
						Amount:    decimal.RequireFromString("1"),
						Threshold: decimal.RequireFromString("1000"),
						Chains:    []string{constant.Ethereum},
					},
				},
			},
		},
	}
	testConf.Init()
	g := new(GateLiquidity)
	b := bot_mocks.NewBot(t)
	b.On("Check", mock.Anything, mock.Anything).Return(
		[]bot.Task{
			{
				Wallet:            constant.ZeroAddress.Hex(),
				TokenInName:       "ETH",
				TokenOutName:      "ETH",
				TokenOutChainName: constant.Ethereum,
				Amount:            decimal.RequireFromString("1"),
			},
		},
		bot.Queue,
		nil,
	)
	g.Bot = b
	tasks, Type, err := g.Check(context.Background(), bot.Params{
		Conf: *testConf,
		Info: bot.Config{
			TokenName: "ETH",
			Chain:     constant.Ethereum,
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, bot.Queue, Type)
	assert.Len(t, tasks, 1)
	assert.Equal(t, "1", tasks[0].Amount.String())
	assert.Equal(t, constant.ZeroAddress.Hex(), tasks[0].Wallet)
	assert.Equal(t, "ETH", tasks[0].TokenOutName)
	assert.Equal(t, constant.Arbitrum, tasks[0].TokenOutChainName)
	assert.Equal(t, constant.Ethereum, tasks[0].TokenInChainName)
}
