package helix_liquidity

import (
	"context"
	"omni-balance/utils/chains"
	"omni-balance/utils/constant"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

var (
	aaveAddressBook = map[string]AvaeConfig{
		constant.Arbitrum: AvaeConfig{
			Chain:     constant.Arbitrum,
			L2Pool:    common.HexToAddress("0x794a61358D6845594F94dc1DB02A252b5b4814aD"),
			Oracle:    common.HexToAddress("0xb56c2F0B653B2e0b10C9b928C8580Ac5Df02C7C7"),
			Multicall: common.HexToAddress("0xcA11bde05977b3631167028862bE2a173976CA11"),
			DebtTokens: map[string]debtTokens{
				"USDT": debtTokens{
					Name:            "USDT",
					AToken:          common.HexToAddress("0x6ab707Aca953eDAeFBc4fD23bA73294241490620"),
					VToken:          common.HexToAddress("0xfb00AC187a8Eb5AFAE4eACE434F493Eb62672df7"),
					UnderlyingToken: common.HexToAddress("0xFd086bC7CD5C481DCC9C85ebE478A1C0b69FCbb9"),
					Decimals:        6,
				},
				"WETH": debtTokens{
					Name:            "WETH",
					AToken:          common.HexToAddress("0xe50fA9b3c56FfB159cB0FCA61F5c9D750e8128c8"),
					VToken:          common.HexToAddress("0x0c84331e39d6658Cd6e6b9ba04736cC4c4734351"),
					UnderlyingToken: common.HexToAddress("0x82aF49447D8a07e3bd95BD0d56f35241523fBab1"),
					Decimals:        18,
				},
			},
		},
		constant.Optimism: AvaeConfig{
			Chain:     constant.Optimism,
			L2Pool:    common.HexToAddress("0x794a61358D6845594F94dc1DB02A252b5b4814aD"),
			Oracle:    common.HexToAddress("0xD81eb3728a631871a7eBBaD631b5f424909f0c77"),
			Multicall: common.HexToAddress("0xcA11bde05977b3631167028862bE2a173976CA11"),
			DebtTokens: map[string]debtTokens{
				"USDT": debtTokens{
					Name:            "USDT",
					AToken:          common.HexToAddress("0x6ab707Aca953eDAeFBc4fD23bA73294241490620"),
					VToken:          common.HexToAddress("0xfb00AC187a8Eb5AFAE4eACE434F493Eb62672df7"),
					UnderlyingToken: common.HexToAddress("0x94b008aA00579c1307B0EF2c499aD98a8ce58e58"),
					Decimals:        6,
				},
				"WETH": debtTokens{
					Name:            "WETH",
					AToken:          common.HexToAddress("0xe50fA9b3c56FfB159cB0FCA61F5c9D750e8128c8"),
					VToken:          common.HexToAddress("0x0c84331e39d6658Cd6e6b9ba04736cC4c4734351"),
					UnderlyingToken: common.HexToAddress("0x4200000000000000000000000000000000000006"),
					Decimals:        18,
				},
			},
		},
		constant.ArbitrumSepolia: AvaeConfig{
			Chain: constant.ArbitrumSepolia,
			DebtTokens: map[string]debtTokens{
				"USDC": debtTokens{
					Name:            "USDC",
					AToken:          common.HexToAddress("0x460b97BD498E1157530AEb3086301d5225b91216"),
					VToken:          common.HexToAddress("0x4fBE3A94C60A5085dA6a2D309965DcF34c36711d"),
					UnderlyingToken: common.HexToAddress("0x75faf114eafb1BDbe2F0316DF893fd58CE46AA4d"),
					Decimals:        6,
				},
			},
		},
	}
)

type AvaeConfig struct {
	Chain      string
	L2Pool     common.Address
	Oracle     common.Address
	Multicall  common.Address
	DebtTokens map[string]debtTokens
}

type debtTokens struct {
	Name            string
	AToken          common.Address
	VToken          common.Address
	UnderlyingToken common.Address
	Decimals        int32
}

type Aave struct {
}

func (a Aave) BalanceOf(ctx context.Context, args DebtParams) (decimal.Decimal, error) {
	conf, ok := aaveAddressBook[args.Chain]
	if !ok || conf.Chain == "" || conf.DebtTokens[args.Token].Name == "" {
		return decimal.Zero, errors.Errorf("chain %s not support", args.Chain)
	}
	atokenBalance, err := chains.GetTokenBalance(ctx, args.Client, conf.DebtTokens[args.Token].AToken.Hex(),
		args.Address.Hex(), conf.DebtTokens[args.Token].Decimals)
	if err != nil {
		return decimal.Zero, errors.Wrap(err, "get atoken balance error")
	}

	vtokenBalance, err := chains.GetTokenBalance(ctx, args.Client, conf.DebtTokens[args.Token].VToken.Hex(),
		args.Address.Hex(), conf.DebtTokens[args.Token].Decimals)
	if err != nil {
		return decimal.Zero, errors.Wrap(err, "get vtoken balance error")
	}
	return atokenBalance.Sub(vtokenBalance), nil
}

func (a Aave) Name() string {
	return "aave"
}
