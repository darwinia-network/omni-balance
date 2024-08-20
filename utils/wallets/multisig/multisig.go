package multisig

import (
	"context"
	"encoding/json"
	"log"
	"omni-balance/utils/constant"
	"omni-balance/utils/error_types"
	"omni-balance/utils/wallets"
	"omni-balance/utils/wallets/dsafe"
	"omni-balance/utils/wallets/mantle_safe"
	"omni-balance/utils/wallets/safe"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient/simulated"
	"github.com/shopspring/decimal"
	"github.com/spf13/cast"
)

type Multisig struct {
	conf wallets.WalletConfig
}

func NewMultisig(conf wallets.WalletConfig) wallets.Wallets {
	return Multisig{conf: conf}
}

func (s Multisig) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"address": s.conf.Address.Hex(),
		"operator": map[string]interface{}{
			"address":         s.conf.Operator.Address.Hex(),
			"operator":        s.conf.Operator.Operator.Hex(),
			"multi_sign_type": s.conf.MultiSignType,
		},
		"multi_sign_type": s.conf.MultiSignType,
	})
}

func (s Multisig) GetChainIdByCtx(ctx context.Context) int {
	chainId := constant.GetChainId(cast.ToString(ctx.Value(constant.ChainNameKeyInCtx)))
	if chainId == 0 {
		log.Fatalf("chain name not found in context")
	}
	return chainId
}

func (s Multisig) GetChainNameByCtx(ctx context.Context) string {
	return cast.ToString(ctx.Value(constant.ChainNameKeyInCtx))
}

func (m Multisig) getRealInstance(ctx context.Context) wallets.Wallets {
	switch m.GetChainNameByCtx(ctx) {
	case constant.DarwiniaDvm:
		return dsafe.NewDsafe(m.conf)
	case constant.Mantle:
		return mantle_safe.NewMantleSafe(m.conf)
	default:
		return safe.NewSafe(m.conf)
	}
}

func (m Multisig) CheckFullAccess(ctx context.Context) error {
	return m.getRealInstance(ctx).CheckFullAccess(ctx)
}

func (m Multisig) GetAddress(isReal ...bool) common.Address {
	if len(isReal) > 0 && isReal[0] {
		if m.conf.Operator.MultiSignType == "" {
			return m.conf.Address
		}
		return m.conf.Operator.Address
	}
	return m.conf.Address
}

func (Multisig) IsDifferentAddress() bool {
	return true
}

func (m Multisig) IsSupportEip712() bool {
	return false
}

func (m Multisig) SignRawMessage(msg []byte) (sig []byte, err error) {
	return nil, error_types.ErrUnsupportedWalletType
}

func (m Multisig) GetExternalBalance(ctx context.Context, tokenAddress common.Address, decimals int32, client simulated.Client) (decimal.Decimal, error) {
	return m.getRealInstance(ctx).GetExternalBalance(ctx, tokenAddress, decimals, client)
}

func (m Multisig) GetBalance(ctx context.Context, tokenAddress common.Address, decimals int32, client simulated.Client) (decimal.Decimal, error) {
	return m.getRealInstance(ctx).GetBalance(ctx, tokenAddress, decimals, client)
}

func (m Multisig) GetNonce(ctx context.Context, client simulated.Client) (uint64, error) {
	return m.getRealInstance(ctx).GetNonce(ctx, client)
}

func (m Multisig) SendTransaction(ctx context.Context, tx *types.LegacyTx, client simulated.Client) (common.Hash, error) {
	return m.getRealInstance(ctx).SendTransaction(ctx, tx, client)
}

func (m Multisig) WaitTransaction(ctx context.Context, txHash common.Hash, client simulated.Client) error {
	return m.getRealInstance(ctx).WaitTransaction(ctx, txHash, client)
}

func (m Multisig) GetRealHash(ctx context.Context, txHash common.Hash, client simulated.Client) (common.Hash, error) {
	return m.getRealInstance(ctx).GetRealHash(ctx, txHash, client)
}

func (m Multisig) Name(ctx context.Context) string {
	return m.getRealInstance(ctx).Name(ctx)
}
