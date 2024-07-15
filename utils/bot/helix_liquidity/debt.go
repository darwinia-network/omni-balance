package helix_liquidity

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient/simulated"
	"github.com/shopspring/decimal"
)

type Debt interface {
	BalanceOf(ctx context.Context, args DebtParams) (decimal.Decimal, error)
	Name() string
}

type DebtParams struct {
	Address common.Address
	Token   string
	Client  simulated.Client
	Chain   string
}

var (
	debtImpl = []Debt{
		Aave{},
	}
)
