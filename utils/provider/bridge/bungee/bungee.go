package bungee

import (
	"context"
	"math/big"
	"omni-balance/utils/chains"
	"omni-balance/utils/configs"
	"omni-balance/utils/constant"
	"omni-balance/utils/error_types"
	"omni-balance/utils/provider"

	log "omni-balance/utils/logging"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/tidwall/gjson"
)

func init() {
	provider.Register(configs.Bridge, New)
}

type Bungee struct {
	conf configs.Config
}

func New(conf configs.Config, _ ...bool) (provider.Provider, error) {
	return Bungee{
		conf: conf,
	}, nil
}

func (r Bungee) GetCost(ctx context.Context, args provider.SwapParams) (provider.TokenInCosts, error) {
	var (
		err        error
		costAmount decimal.Decimal
	)
	if _, ok := SupportedChain[args.TargetChain]; !ok {
		return nil, error_types.ErrUnsupportedTokenAndChain
	}

	if args.SourceChain == "" || args.SourceToken == "" {
		var quote gjson.Result
		args.SourceToken, args.SourceChain, costAmount, quote, err = r.GetBestQuote(ctx, args)
		if err != nil {
			return nil, err
		}
		if len(quote.Get("result.routes").Array()) == 0 {
			return nil, error_types.ErrUnsupportedTokenAndChain
		}
	}

	if args.SourceChain == "" || args.SourceToken == "" || costAmount.IsZero() {
		return nil, error_types.ErrUnsupportedTokenAndChain
	}
	return provider.TokenInCosts{
		provider.TokenInCost{
			TokenName:  args.SourceToken,
			CostAmount: costAmount,
		},
	}, nil
}

func (r Bungee) CheckToken(ctx context.Context, tokenName, tokenInChainName, tokenOutChainName string, amount decimal.Decimal) (bool, error) {
	var (
		tokenIn      = r.conf.GetTokenInfoOnChain(tokenName, tokenInChainName)
		tokenOut     = r.conf.GetTokenInfoOnChain(tokenName, tokenOutChainName)
		tokenInChain = r.conf.GetChainConfig(tokenInChainName)
		tokenOuChain = r.conf.GetChainConfig(tokenOutChainName)
	)

	if _, ok := SupportedChain[tokenOuChain.Name]; !ok {
		return false, error_types.ErrUnsupportedTokenAndChain
	}

	if _, ok := SupportedChain[tokenInChain.Name]; !ok {
		return false, error_types.ErrUnsupportedTokenAndChain
	}

	quote, err := r.Quote(ctx, QuoteParams{
		FromTokenAddress: common.HexToAddress(tokenIn.ContractAddress),
		ToTokenAddress:   common.HexToAddress(tokenOut.ContractAddress),
		AmountWei:        decimal.NewFromBigInt(chains.EthToWei(amount, tokenIn.Decimals), 0),
		FromTokenChainId: tokenInChain.Id,
		ToTokenChainId:   tokenOuChain.Id,
	})
	if len(quote.Get("result").Get("routes").Array()) == 0 {
		return false, error_types.ErrUnsupportedTokenAndChain
	}
	return err == nil, err
}

func (r Bungee) Swap(ctx context.Context, args provider.SwapParams) (provider.SwapResult, error) {
	var (
		err               error
		targetChain       = r.conf.GetChainConfig(args.TargetChain)
		tokenOut          = r.conf.GetTokenInfoOnChain(args.TargetToken, targetChain.Name)
		tokenOutAmount    = args.Amount
		tokenOutAmountWei = decimal.NewFromBigInt(chains.EthToWei(tokenOutAmount, tokenOut.Decimals), 0)
		tokenIn           configs.Token
		tokenInAmount     decimal.Decimal
		tokenInAmountWei  decimal.Decimal
		costAmount        decimal.Decimal
		history           = args.LastHistory
		tx                = history.Tx
		actionNumber      = Action2Int(history.Actions)
		isActionSuccess   = history.Status == provider.TxStatusSuccess.String()
		sourceChain       configs.Chain
		quote             gjson.Result
	)

	args.SourceToken, args.SourceChain, costAmount, _, err = r.GetBestQuote(ctx, args)
	if err != nil {
		return provider.SwapResult{}, err
	}

	if args.SourceChain == "" || args.SourceToken == "" {
		log.Fatalf("#%d %s source chain and token is required", args.OrderId, args.TargetToken)
	}
	tokenIn = r.conf.GetTokenInfoOnChain(args.SourceToken, args.SourceChain)
	sourceChain = r.conf.GetChainConfig(args.SourceChain)
	tokenInAmount = costAmount
	tokenInAmountWei = decimal.NewFromBigInt(chains.EthToWei(tokenInAmount, tokenIn.Decimals), 0)
	isTokenInNative := r.conf.IsNativeToken(sourceChain.Name, tokenIn.Name)

	quotes, err := r.Quote(ctx, QuoteParams{
		FromTokenAddress: common.HexToAddress(tokenIn.ContractAddress),
		ToTokenAddress:   common.HexToAddress(tokenOut.ContractAddress),
		AmountWei:        decimal.NewFromBigInt(chains.EthToWei(costAmount, tokenIn.Decimals), 0),
		FromTokenChainId: sourceChain.Id,
		ToTokenChainId:   constant.GetChainId(args.TargetChain),
	})
	if err != nil {
		return provider.SwapResult{}, err
	}

	if len(quotes.Get("result").Get("routes").Array()) == 0 {
		return provider.SwapResult{}, error_types.ErrUnsupportedTokenAndChain
	}
	quote = quotes.Get("result").Get("routes").Array()[0]
	tokenOutAmountWei = decimal.RequireFromString(quote.Get("toAmount").String())
	tokenOutAmount = chains.WeiToEth(tokenOutAmountWei.BigInt(), tokenOut.Decimals)
	args.Amount = tokenOutAmount
	if tokenOutAmount.LessThanOrEqual(decimal.Zero) {
		return provider.SwapResult{}, errors.New("token out amount is zero")
	}

	ctx = context.WithValue(ctx, constant.ChainNameKeyInCtx, sourceChain.Name)
	client, err := chains.NewTryClient(ctx, sourceChain.RpcEndpoints)
	if err != nil {
		return provider.SwapResult{}, err
	}
	var (
		sr = new(provider.SwapResult).
			SetTokenInName(tokenIn.Name).
			SetTokenInChainName(args.SourceChain).
			SetProviderName(r.Name()).
			SetProviderType(r.Type()).
			SetCurrentChain(args.SourceChain).
			SetTx(args.LastHistory.Tx).
			SetReciever(args.Receiver)
		sh = &provider.SwapHistory{
			ProviderName: r.Name(),
			ProviderType: string(r.Type()),
			Amount:       args.Amount,
			CurrentChain: args.SourceChain,
			Tx:           history.Tx,
		}
	)

	if !isTokenInNative && actionNumber <= 1 && !isActionSuccess {
		log.Debugf("#%d %s is not native token, need approve", args.OrderId, tokenIn.Name)
		args.RecordFn(sh.SetActions(ApproveTransactionAction).SetStatus(provider.TxStatusPending).Out())
		spender := quote.Get("userTxs.0.approvalData.allowanceTarget").String()
		if spender == "" {
			return provider.SwapResult{}, errors.New("spender is empty")
		}
		err = chains.TokenApprove(ctx,
			chains.TokenApproveParams{
				ChainId:         int64(sourceChain.Id),
				TokenAddress:    common.HexToAddress(tokenIn.ContractAddress),
				Owner:           args.Sender.GetAddress(true),
				SendTransaction: args.Sender.SendTransaction,
				WaitTransaction: args.Sender.WaitTransaction,
				Spender:         common.HexToAddress(spender),
				AmountWei:       tokenInAmountWei.Mul(decimal.RequireFromString("1.02")),
				IsNotWaitTx:     r.conf.GetWalletConfig(string(args.Sender.GetAddress().Hex())).MultiSignType != "",
				Client:          client,
			})
		if err != nil {
			args.RecordFn(sh.SetActions(ApproveTransactionAction).SetStatus(provider.TxStatusFailed).Out(), err)
			return sr.SetError(err).SetStatus(provider.TxStatusFailed).Out(), err
		}
		log.Debugf("#%d %s approve transaction success", args.OrderId, tokenIn.Name)
		args.RecordFn(sh.SetActions(ApproveTransactionAction).SetStatus(provider.TxStatusSuccess).Out())
	}

	if actionNumber <= 2 && (!isActionSuccess || actionNumber == 1) {
		amount := args.Amount.Copy()
		args.Amount = tokenInAmount
		ctx = provider.WithNotify(ctx, provider.WithNotifyParams{
			OrderId:         args.OrderId,
			Receiver:        common.HexToAddress(args.Receiver),
			TokenIn:         tokenIn.Name,
			TokenOut:        tokenOut.Name,
			TokenInChain:    args.SourceChain,
			TokenOutChain:   args.TargetChain,
			ProviderName:    r.Name(),
			TokenInAmount:   tokenInAmount,
			TokenOutAmount:  tokenOutAmount,
			TransactionType: provider.SwapTransactionAction,
		})

		buildTx, err := r.BuildTx(ctx, quote, args.Sender.GetAddress(true), common.HexToAddress(args.Receiver))
		if err != nil {
			return sr.SetError(err).SetStatus(provider.TxStatusFailed).Out(), errors.Wrap(err, "build tx error")
		}
		args.Amount = amount
		if !isTokenInNative && buildTx.Value.Cmp(big.NewInt(0)) != 0 {
			err = errors.Errorf("tokenin is not native token, but value is not zero")
			return sr.SetError(err).SetStatus(provider.TxStatusFailed).Out(), err
		}
		if isTokenInNative && buildTx.Value.Cmp(big.NewInt(0)) == 0 {
			err = errors.Errorf("tokenin is native token, but value is zero")
			return sr.SetError(err).SetStatus(provider.TxStatusFailed).Out(), err
		}
		if isTokenInNative && decimal.NewFromBigInt(buildTx.Value, 0).GreaterThan(tokenInAmountWei.Mul(decimal.RequireFromString("1.5"))) {
			err = errors.Errorf("tx value is too high, tx value: %s, amount: %s", buildTx.Value.String(), tokenInAmountWei.String())
			return sr.SetError(err).SetStatus(provider.TxStatusFailed).Out(), err
		}
		sr.SetOrder(buildTx)
		args.RecordFn(sh.SetActions(SourceChainSendingAction).SetStatus(provider.TxStatusPending).Out())
		txHash, err := args.Sender.SendTransaction(ctx, buildTx, client)
		if err != nil {
			args.RecordFn(sh.SetActions(SourceChainSendingAction).SetStatus(provider.TxStatusFailed).Out(), err)
			return sr.SetError(err).SetStatus(provider.TxStatusFailed).Out(), errors.Wrap(err, "send tx error")
		}
		log.Debugf("#%d %s sending tx on chain success", args.OrderId, tokenIn.Name)
		sh = sh.SetTx(txHash.Hex())
		sr = sr.SetTx(txHash.Hex()).SetOrderId(txHash.Hex())
		args.RecordFn(sh.SetActions(SourceChainSendingAction).SetStatus(provider.TxStatusSuccess).Out())
		tx = txHash.Hex()
	}

	args.RecordFn(sh.SetActions(WaitForTxAction).SetStatus(provider.TxStatusPending).Out())
	log.Debugf("#%d %s waiting for tx on chain", args.OrderId, tokenIn.Name)
	if err := args.Sender.WaitTransaction(ctx, common.HexToHash(tx), client); err != nil {
		args.RecordFn(sh.SetActions(WaitForTxAction).SetStatus(provider.TxStatusFailed).Out(), err)
		return sr.SetError(err).SetStatus(provider.TxStatusFailed).Out(), errors.Wrap(err, "wait tx error")
	}
	realHash, err := args.Sender.GetRealHash(ctx, common.HexToHash(tx), client)
	if err != nil {
		args.RecordFn(sh.SetActions(WaitForTxAction).SetStatus(provider.TxStatusFailed).Out(), err)
		return sr.SetError(err).SetStatus(provider.TxStatusFailed).Out(), errors.Wrap(err, "get real hash error")
	}
	log.Debugf("#%d %s waiting for tx in router nitro", args.OrderId, tokenIn.Name)
	if err := r.WaitForTx(ctx, sourceChain.Id, targetChain.Id, realHash); err != nil {
		args.RecordFn(sh.SetActions(WaitForTxAction).SetStatus(provider.TxStatusFailed).Out(), err)
		return sr.SetError(err).SetStatus(provider.TxStatusFailed).Out(), errors.Wrap(err, "wait router nitro error")
	}
	log.Debugf("#%d %s waiting for tx success in routernitro", args.OrderId, tokenIn.Name)
	args.RecordFn(sh.SetActions(WaitForTxAction).SetStatus(provider.TxStatusSuccess).SetCurrentChain(targetChain.Name).Out())
	return sr.SetStatus(provider.TxStatusSuccess).SetCurrentChain(targetChain.Name).Out(), nil
}

func (r Bungee) Help() []string {
	return []string{
		// BUNGEE_API_KEY env
		"You can Set Custom Bungee API Key by setting the BUNGEE_API_KEY environment variable. The default value is " + apiKey,
		"see https://www.bungee.exchange",
	}
}

func (r Bungee) Name() string {
	return "bungee"
}

func (r Bungee) Type() configs.ProviderType {
	return configs.Bridge
}
