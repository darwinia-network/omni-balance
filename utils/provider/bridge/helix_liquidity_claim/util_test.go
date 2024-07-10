package helix_liquidity_claim

import (
	"context"
	"math/big"
	"omni-balance/utils/chains"
	"omni-balance/utils/configs"
	"omni-balance/utils/constant"
	_ "omni-balance/utils/wallets/safe"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func Test_encodeWithdrawLiquidity(t *testing.T) {
	conf := &configs.Config{
		Debug: true,
		Chains: []configs.Chain{
			{
				Id:           constant.GetChainId(constant.Arbitrum),
				Name:         constant.Polygon,
				NativeToken:  "MATIC",
				RpcEndpoints: []string{"https://polygon-rpc.com"},
				Tokens: []configs.Token{
					{
						ContractAddress: "0x0000000000000000000000000000000000000000",
						Decimals:        18,
						Name:            "MATIC",
					},
					{
						ContractAddress: "0xc2132D05D31c914a87C6611C10748AEb04B58e8F",
						Decimals:        6,
						Name:            "USDT",
					},
				},
			},
			{
				Id:           constant.GetChainId(constant.Arbitrum),
				Name:         constant.Arbitrum,
				NativeToken:  "ETH",
				RpcEndpoints: []string{"https://arb1.arbitrum.io/rpc"},
				Tokens: []configs.Token{
					{
						ContractAddress: "0x0000000000000000000000000000000000000000",
						Decimals:        18,
						Name:            "ETH",
					},
					{
						ContractAddress: "0xFd086bC7CD5C481DCC9C85ebE478A1C0b69FCbb9",
						Decimals:        6,
						Name:            "USDT",
					},
				},
			},
		},
		Wallets: []configs.Wallet{
			{

				Address: "0x000000000Bb6a011dB294ce3F3423f00EAc4959e",
				Operator: configs.Operator{
					Address:       common.HexToAddress("0x6dB77823ECa92dAcbbcfcC976dccF03370F22874"),
					Operator:      common.HexToAddress("0xC5a809900B5BFb46B1B3892e419e69331B8FBC6c"),
					MultiSignType: "safe",
				},
				MultiSignType: "safe",
				Tokens: []configs.WalletToken{
					{
						Name:      "USDT",
						Amount:    decimal.RequireFromString("1"),
						Threshold: decimal.RequireFromString("1"),
						Chains: []string{
							"arbitrum",
							"polygon",
						},
					},
				},
			},
		},
	}
	conf.Init()

	c, err := New(*conf)
	assert.NoError(t, err)
	c1 := c.(Claim)
	appPayload, err := c1.encodeWithdrawLiquidity(
		[]string{"0xd84be70f3b8d54cc248da2aae1d73f2536b1e4c02fa53697991e284fd88bb737", "0x8bf1db785ff2091cf02cd03c399d441d46289729c074a99c06f7e90cfd3f77ad"},
		big.NewInt(137),
		common.HexToAddress("0x000000000Bb6a011dB294ce3F3423f00EAc4959e"))
	assert.NoError(t, err)
	assert.Equal(t, "7425b8b500000000000000000000000000000000000000000000000000000000000000600000000000000000000000000000000000000000000000000000000000000089000000000000000000000000000000000bb6a011db294ce3f3423f00eac4959e0000000000000000000000000000000000000000000000000000000000000002d84be70f3b8d54cc248da2aae1d73f2536b1e4c02fa53697991e284fd88bb7378bf1db785ff2091cf02cd03c399d441d46289729c074a99c06f7e90cfd3f77ad", common.Bytes2Hex(appPayload))
}

func Test_ClaimBuildTx(t *testing.T) {

	client, err := chains.NewTryClient(context.TODO(), []string{
		"https://polygon-rpc.com",
	})
	assert.NoError(t, err)
	conf := &configs.Config{
		Debug: true,
		Chains: []configs.Chain{
			{
				Id:           constant.GetChainId(constant.Arbitrum),
				Name:         constant.Polygon,
				NativeToken:  "MATIC",
				RpcEndpoints: []string{"https://polygon-rpc.com"},
				Tokens: []configs.Token{
					{
						ContractAddress: "0x0000000000000000000000000000000000000000",
						Decimals:        18,
						Name:            "MATIC",
					},
					{
						ContractAddress: "0xc2132D05D31c914a87C6611C10748AEb04B58e8F",
						Decimals:        6,
						Name:            "USDT",
					},
				},
			},
			{
				Id:           constant.GetChainId(constant.Arbitrum),
				Name:         constant.Arbitrum,
				NativeToken:  "ETH",
				RpcEndpoints: []string{"https://arb1.arbitrum.io/rpc"},
				Tokens: []configs.Token{
					{
						ContractAddress: "0x0000000000000000000000000000000000000000",
						Decimals:        18,
						Name:            "ETH",
					},
					{
						ContractAddress: "0xFd086bC7CD5C481DCC9C85ebE478A1C0b69FCbb9",
						Decimals:        6,
						Name:            "USDT",
					},
				},
			},
		},
		Wallets: []configs.Wallet{
			{

				Address: "0x000000000Bb6a011dB294ce3F3423f00EAc4959e",
				Operator: configs.Operator{
					Address:       common.HexToAddress("0x6dB77823ECa92dAcbbcfcC976dccF03370F22874"),
					Operator:      common.HexToAddress("0xC5a809900B5BFb46B1B3892e419e69331B8FBC6c"),
					MultiSignType: "safe",
				},
				MultiSignType: "safe",
				Tokens: []configs.WalletToken{
					{
						Name:      "USDT",
						Amount:    decimal.RequireFromString("1"),
						Threshold: decimal.RequireFromString("1"),
						Chains: []string{
							"arbitrum",
							"polygon",
						},
					},
				},
			},
		},
	}
	conf.Init()

	c, err := New(*conf)
	if err != nil {
		panic(err)
	}
	c1 := c.(Claim)
	tx, err := c1.BuildTx(context.TODO(), client, conf.GetWallet("0x000000000Bb6a011dB294ce3F3423f00EAc4959e"), NeedWithdrawRecords{
		ClaimInfo: ClaimInfo{
			TransferIds: []string{
				"0xd84be70f3b8d54cc248da2aae1d73f2536b1e4c02fa53697991e284fd88bb737",
				"0x8bf1db785ff2091cf02cd03c399d441d46289729c074a99c06f7e90cfd3f77ad",
			},
			TotalAmount: decimal.RequireFromString("51"),
			Channel:     "layerzero",
		},
		TokenName: "USDT",
		FromChain: "arbitrum",
		ToChain:   "polygon",
	})
	assert.NoError(t, err)
	expected := "29f085f4000000000000000000000000000000000000000000000000000000000000a4b10000000000000000000000000000000000000000000000000000000000000080000000000000000000000000000000000bb6a011db294ce3f3423f00eac4959e00000000000000000000000000000000000000000000000000000000000000e00000000000000000000000000000000000000000000000000000000000000002d84be70f3b8d54cc248da2aae1d73f2536b1e4c02fa53697991e284fd88bb7378bf1db785ff2091cf02cd03c399d441d46289729c074a99c06f7e90cfd3f77ad0000000000000000000000000000000000000000000000000000000000000014000000000bb6a011db294ce3f3423f00eac4959e000000000000000000000000"
	assert.Equal(t, expected, common.Bytes2Hex(tx.Data))
	assert.Equal(t, "0xbA5D580B18b6436411562981e02c8A9aA1776D10", tx.To.Hex())
	assert.GreaterOrEqual(t, tx.Value.Int64(), int64(0))
}
