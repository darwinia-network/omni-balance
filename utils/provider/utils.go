package provider

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient/simulated"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"omni-balance/utils"
	"omni-balance/utils/chains"
	"omni-balance/utils/configs"
	"omni-balance/utils/notice"
)

type TransactionType string

var (
	ApproveTransactionAction  TransactionType = "Approve"
	TransferTransactionAction TransactionType = "Transfer"
	SwapTransactionAction     TransactionType = "Swap"
)

var (
	transferSent     = "transferSent"
	transferReceived = "transferReceived"
)

func action2Int(action string) int {
	switch action {
	case transferSent:
		return 0
	case transferReceived:
		return 1
	default:
		return -1
	}
}

func Transfer(ctx context.Context, conf configs.Config, args SwapParams, client simulated.Client) (SwapResult, error) {

	var (
		log          = utils.GetLogFromCtx(ctx)
		last         = args.LastHistory
		actionNumber = action2Int(args.LastHistory.Actions)
		token        = conf.GetTokenInfoOnChain(args.TargetToken, args.TargetChain)
		result       = &SwapResult{
			TokenInChainName: args.TargetChain,
			ProviderType:     "direct",
			TokenInName:      args.TargetToken,
			ProviderName:     "transfer",
			CurrentChain:     args.TargetChain,
		}
		recordFn = func(s SwapHistory, errs ...error) {
			s.ProviderType = string(result.ProviderType)
			s.ProviderName = result.ProviderName
			s.Amount = args.Amount
			s.CurrentChain = args.TargetChain
			args.RecordFn(s, errs...)
		}
	)

	log = log.WithFields(logrus.Fields{
		"last": last,
	})
	if !utils.InArray(args.LastHistory.Actions, []string{transferSent, transferReceived}) {
		actionNumber = 0
		args.LastHistory.Status = ""
	}

	var txHash = last.Tx
	if actionNumber < 1 && args.LastHistory.Status != TxStatusSuccess.String() {
		transaction, err := chains.BuildSendToken(ctx, chains.SendTokenParams{
			Client:       client,
			Sender:       args.Sender.GetAddress(true),
			GetBalance:   args.Sender.GetBalance,
			TokenAddress: common.HexToAddress(token.ContractAddress),
			ToAddress:    args.Sender.GetAddress(),
			AmountWei:    decimal.NewFromBigInt(chains.EthToWei(args.Amount, token.Decimals), 0),
		})

		if err != nil {
			return *result, errors.Wrap(err, "send token error")
		}
		log.Debugf("%s transfer %s %s to %s", args.Sender.GetAddress(true), args.TargetToken,
			args.Amount, args.Sender.GetAddress())
		recordFn(SwapHistory{
			Actions: transferSent,
			Status:  TxStatusPending.String(),
		})
		tx, err := args.Sender.SendTransaction(ctx, transaction, client)
		if err != nil {
			recordFn(SwapHistory{
				Actions: transferSent,
				Status:  TxStatusFailed.String(),
			})
			return *result, errors.Wrap(err, "send transaction error")
		}
		recordFn(SwapHistory{
			Actions: transferSent,
			Status:  TxStatusSuccess.String(),
			Tx:      tx.Hex(),
		})
		txHash = tx.Hex()
	}

	if txHash == "" {
		return *result, errors.New("tx is empty")
	}
	log.Debugf("waiting for tx %s", txHash)
	err := args.Sender.WaitTransaction(ctx, common.HexToHash(txHash), client)
	if err != nil {
		recordFn(SwapHistory{
			Actions: transferReceived,
			Status:  TxStatusFailed.String(),
			Tx:      txHash,
		})
		return *result, errors.Wrap(err, "wait for tx error")
	}
	result.Tx = txHash
	result.OrderId = txHash
	result.Status = TxStatusSuccess
	recordFn(SwapHistory{
		Actions: transferReceived,
		Status:  TxStatusSuccess.String(),
		Tx:      txHash,
	})
	return *result, nil
}

type GetTokenCrossChainProvidersParams struct {
	SourceChains []string
	Conf         configs.Config
	TokenName    string
	SourceChain  string
	TargetChain  string
	Amount       decimal.Decimal
}

// GetTokenCrossChainProviders retrieves the cross-chain providers for the given token and parameters.
func GetTokenCrossChainProviders(ctx context.Context, args GetTokenCrossChainProvidersParams) (providers []Provider, err error) {
	for _, fn := range LiquidityProviderTypeAndConf(configs.Bridge, args.Conf) {
		bridge, err := fn(args.Conf)
		if err != nil {
			continue
		}
		var hasConfig bool
		for _, v := range args.Conf.LiquidityProviders {
			if v.Type != configs.Bridge || v.LiquidityName == bridge.Name() {
				continue
			}
			hasConfig = true
			break
		}
		if !hasConfig {
			continue
		}
		ok, err := bridge.CheckToken(ctx, args.TokenName, args.SourceChain, args.TargetChain, args.Amount)
		if err != nil {
			logrus.Warnf("check token %s in %s to %s failed: %s", args.TokenName, args.SourceChain, args.TargetChain, err)
			continue
		}
		if !ok {
			continue
		}
		providers = append(providers, bridge)
	}
	return providers, nil
}

type WithNotifyParams struct {
	TokenIn, TokenOut, TokenInChain, TokenOutChain, ProviderName string
	TokenInAmount, TokenOutAmount                                decimal.Decimal
	TransactionType                                              TransactionType
}

func WithNotify(ctx context.Context, args WithNotifyParams) context.Context {
	var fields = map[string]string{
		"transactionType": string(args.TransactionType),
		"providerName":    args.ProviderName,
	}
	if args.TokenIn != "" {
		fields["tokenIn"] = args.TokenIn
		if !args.TokenInAmount.IsZero() {
			fields["tokenIn"] = fmt.Sprintf("%s %s", args.TokenInAmount, fields["tokenIn"])
		}
		if args.TokenInChain != "" {
			fields["tokenIn"] = fmt.Sprintf("%s on %s", fields["tokenIn"], args.TokenInChain)
		}
	}

	if args.TokenOut != "" {
		fields["tokenOut"] = args.TokenOut
		if !args.TokenOutAmount.IsZero() {
			fields["tokenOut"] = fmt.Sprintf("%s %s", args.TokenOutAmount, fields["tokenOut"])
		}
		if args.TokenOutChain != "" {
			fields["tokenOut"] = fmt.Sprintf("%s on %s", fields["tokenOut"], args.TokenOutChain)
		}
	}
	return notice.WithFields(ctx, fields)
}
