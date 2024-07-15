package balance_on_chain

import (
	"context"
	"omni-balance/utils/bot"
	"omni-balance/utils/configs"
	"omni-balance/utils/constant"
	"omni-balance/utils/wallets/wallet_mocks"
	"strconv"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestBalanceOnChain_Check(t *testing.T) {
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
			configs.Wallet{
				Address: constant.ZeroAddress.Hex(),
				Tokens: []configs.WalletToken{
					configs.WalletToken{
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
	b := new(BalanceOnChain)
	w := wallet_mocks.NewWallets(t)
	w.On("GetExternalBalance", context.Background(), constant.ZeroAddress, int32(18), nil).Return(decimal.RequireFromString("1000"), nil)
	w.On("GetAddress").Return(constant.ZeroAddress)
	tasks, Type, err := b.Check(context.Background(), bot.Params{
		Conf: *testConf,
		Info: bot.Config{
			Wallet:    w,
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
	assert.Equal(t, constant.Ethereum, tasks[0].TokenOutChainName)

	w = wallet_mocks.NewWallets(t)
	w.On("GetExternalBalance", context.Background(), constant.ZeroAddress, int32(18), nil).Return(decimal.RequireFromString("1"), nil)
	w.On("GetAddress").Return(constant.ZeroAddress)
	tasks, Type, err = b.Check(context.Background(), bot.Params{
		Conf: *testConf,
		Info: bot.Config{
			Wallet:    w,
			TokenName: "ETH",
			Chain:     constant.Ethereum,
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, bot.Queue, Type)
	assert.Len(t, tasks, 1)
	// 1000 + (1000 * 0.3)
	assert.Equal(t, strconv.Itoa(1000+(1000*0.3)), tasks[0].Amount.String())
	assert.Equal(t, constant.ZeroAddress.Hex(), tasks[0].Wallet)
	assert.Equal(t, "ETH", tasks[0].TokenOutName)
	assert.Equal(t, constant.Ethereum, tasks[0].TokenOutChainName)
}
