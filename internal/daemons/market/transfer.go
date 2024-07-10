package market

import (
	"context"
	"github.com/ethereum/go-ethereum/ethclient/simulated"
	"github.com/pkg/errors"
	"omni-balance/internal/models"
	"omni-balance/utils/configs"
	"omni-balance/utils/constant"
	"omni-balance/utils/error_types"
	"omni-balance/utils/provider"
)

func transfer(ctx context.Context, order models.Order, args provider.SwapParams,
	conf configs.Config, client simulated.Client) (bool, error) {
	ctx = context.WithValue(ctx, constant.ChainNameKeyInCtx, order.TargetChainName)
	if order.Status != models.OrderStatusWaitTransferFromOperator {
		return false, errors.Errorf("order #%d status is %s, not wait transfer from operator", order.ID, order.Status)
	}
	result, err := provider.Transfer(ctx, conf, args, client)
	if errors.Is(err, error_types.ErrNativeTokenInsufficient) ||
		errors.Is(err, error_types.ErrWalletLocked) ||
		errors.Is(err, context.Canceled) {
		return true, errors.Wrap(err, "transfer error")
	}
	if err == nil {
		return true, createUpdateLog(ctx, order, result, conf, client)
	}
	if !errors.Is(errors.Unwrap(err), error_types.ErrInsufficientBalance) &&
		!errors.Is(errors.Unwrap(err), error_types.ErrInsufficientLiquidity) {
		return false, errors.Wrap(err, "transfer not is insufficient balance")
	}
	return false, nil
}
