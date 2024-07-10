package li

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"net/url"
	"omni-balance/utils"
	"omni-balance/utils/chains"
	"omni-balance/utils/constant"
	"omni-balance/utils/error_types"
	"omni-balance/utils/provider"
	"strconv"
	"strings"
	"time"
)

// Quote see https://apidocs.li.fi/reference/get_quote
func (l Li) Quote(ctx context.Context, args QuoteParams) (Quote, error) {
	var params = &url.Values{}
	params.Set("fromChain", strconv.Itoa(args.FromChainId))
	params.Set("toChain", strconv.Itoa(args.ToChainId))
	params.Set("fromToken", args.FromToken.Hex())
	params.Set("toToken", args.ToToken.Hex())
	params.Set("fromAmount", args.FromAmountWei.String())
	params.Set("fromAddress", args.FromAddress.Hex())
	if args.ToAddress.Cmp(constant.ZeroAddress) != 0 {
		params.Set("toAddress", args.ToAddress.Hex())
	}
	params.Set("order", "RECOMMENDED")
	params.Set("slippage", "0.01")
	var result Quote
	err := utils.Request(ctx, "GET", "https://li.quest/v1/quote?"+params.Encode(), nil, &result)
	if err != nil {
		return result, err
	}
	return result, result.Error()
}

func (l Li) GetBestQuote(ctx context.Context, args provider.SwapParams) (tokenInName, tokenInChainName string,
	tokenInAmount decimal.Decimal, quote Quote, err error) {
	log := utils.GetLogFromCtx(ctx).WithField("name", l.Name())
	if args.TargetToken == "" || args.TargetChain == "" {
		return tokenInName, tokenInChainName, tokenInAmount, Quote{}, errors.New("target token or target chain is empty")
	}

	var (
		tokenOut = l.conf.GetTokenInfoOnChain(args.TargetToken, args.TargetChain)
	)

	getQuote := func(chainName, token string) error {
		sourceToken := l.conf.GetTokenInfoOnChain(token, chainName)
		currentLog := log.WithField("TokenIn", sourceToken.Name).WithField("sourceChain", chainName).
			WithField("TargetToken", args.TargetToken).WithField("TargetChain", args.TargetChain)
		currentLog.Debug("start check tokenIn")
		chain := l.conf.GetChainConfig(chainName)
		tokenIn := l.conf.GetTokenInfoOnChain(sourceToken.Name, chainName)
		if tokenIn.ContractAddress == "" {
			return errors.New("tokenIn contract address is empty")
		}
		client, err := chains.NewTryClient(ctx, chain.RpcEndpoints)
		if err != nil {
			currentLog.Warnf("get chain %s client error: %s", chain.Name, err)
			return err
		}

		tokenInTestBalance := decimal.RequireFromString("1")
		tokenInTestBalanceWei := decimal.NewFromBigInt(chains.EthToWei(tokenInTestBalance, tokenIn.Decimals), 0)

		quoteData, err := l.Quote(ctx, QuoteParams{
			FromChainId:   chain.Id,
			ToChainId:     constant.GetChainId(args.TargetChain),
			FromToken:     common.HexToAddress(tokenIn.ContractAddress),
			ToToken:       common.HexToAddress(tokenOut.ContractAddress),
			FromAmountWei: tokenInTestBalanceWei,
			FromAddress:   args.Sender.GetAddress(true),
			ToAddress:     common.HexToAddress(args.Receiver),
		})
		if err != nil {
			currentLog.Debugf("get quote error: %s", err)
			return errors.Wrap(err, "get quote")
		}
		currentLog = currentLog.WithField("quote", utils.ToMap(quoteData))

		minimumReceived := chains.WeiToEth(quoteData.Estimate.ToAmountMin.BigInt(), tokenOut.Decimals)
		needBalance := tokenInTestBalance.Div(minimumReceived).Mul(args.Amount)

		balance, err := chains.GetTokenBalance(ctx, client, tokenIn.ContractAddress,
			args.Sender.GetAddress(true).Hex(), tokenIn.Decimals)
		if err != nil {
			currentLog.Debugf("get balance error: %s", err)
			return errors.Wrap(err, "get balance")
		}

		log.Debugf("need %s balance: %s, wallet %s balance: %s on %s",
			tokenIn.Name, needBalance, tokenIn.Name, balance, chainName)
		if needBalance.GreaterThan(balance) {
			currentLog.Debugf("%s need balance: %s, balance: %s", tokenIn.Name, needBalance, balance)
			return errors.New("not enough balance")
		}
		if tokenInAmount.Equal(decimal.Zero) {
			tokenInAmount = needBalance
		}
		if tokenInAmount.GreaterThan(needBalance) {
			currentLog.Debugf("need balance: %s, balance: %s", needBalance, balance)
			return errors.New("not enough balance")
		}
		tokenInAmount = needBalance
		tokenInName = sourceToken.Name
		tokenInChainName = chainName
		quote = quoteData
		return nil
	}

	if args.SourceChain != "" && args.SourceToken != "" {
		err = getQuote(args.SourceChain, args.SourceToken)
		return
	}

	for _, sourceToken := range l.conf.SourceTokens {
		for _, v := range sourceToken.Chains {
			if strings.EqualFold(v, args.TargetChain) && sourceToken.Name == args.TargetToken {
				continue
			}
			if err := getQuote(v, sourceToken.Name); err != nil {
				continue
			}
		}
	}

	if tokenInChainName == "" || tokenInName == "" || tokenInAmount.IsZero() {
		return "", "", tokenInAmount, Quote{}, error_types.ErrUnsupportedTokenAndChain
	}
	tokenIn := l.conf.GetTokenInfoOnChain(tokenInName, tokenInChainName)

	quoteData, err := l.Quote(ctx, QuoteParams{
		FromChainId:   constant.GetChainId(tokenInChainName),
		ToChainId:     constant.GetChainId(args.TargetChain),
		FromToken:     common.HexToAddress(tokenIn.ContractAddress),
		ToToken:       common.HexToAddress(tokenOut.ContractAddress),
		FromAmountWei: decimal.NewFromBigInt(chains.EthToWei(tokenInAmount, tokenIn.Decimals), 0),
		FromAddress:   args.Sender.GetAddress(true),
		ToAddress:     common.HexToAddress(args.Receiver),
	})
	if err != nil {
		return tokenInName, tokenInChainName, tokenInAmount, Quote{}, err
	}
	quote = quoteData
	log.Debugf("best tokenInName: %s, tokenInChainName: %s, tokenInAmount: %s", tokenInName, tokenInChainName, tokenInAmount)
	return
}

func (l Li) WaitForTx(ctx context.Context, hash common.Hash) error {
	log := utils.GetLogFromCtx(ctx)
	var (
		t     = time.NewTicker(time.Second * 10)
		count int64
	)
	defer t.Stop()
	for count < 600 {
		select {
		case <-ctx.Done():
			log.Debugf("wait for tx %s canceled", hash)
			return context.Canceled
		case <-t.C:
			log.Debugf("start check tx %s", hash)
			status, err := l.Status(ctx, hash)
			if errors.Is(err, error_types.ErrNotFound) {
				count++
				log.Debugf("tx %s not found, count: %d", hash, count)
				continue
			}

			log.WithFields(utils.ToMap(status)).Infof("tx %s status: %s, substatus: %s", hash, status.Status, status.Substatus)
			if status.Status == "" {
				log.Debugf("tx %s status is empty, try it 1 minute later", hash)
				time.Sleep(time.Minute)
				continue
			}
			switch status.Status {
			case "PENDING":
				count = 0
				continue
			case "DONE":
				if status.Substatus != "COMPLETED" {
					return errors.New(subStatus2Text[status.Substatus])
				}
				return nil
			case "FAILED":
				return errors.New(subStatus2Text[status.Substatus])
			default:
				count++
				continue
			}
		}
	}
	return errors.New("wait for tx timeout")
}

func (l Li) Status(ctx context.Context, hash common.Hash) (Status, error) {
	var params = &url.Values{}
	params.Set("txHash", hash.Hex())
	var status Status
	err := utils.Request(ctx, "GET", "https://li.quest/v1/status?"+params.Encode(), nil, &status)
	if err != nil {
		return Status{}, err
	}
	if err := status.Error(); err != nil {
		return Status{}, err
	}
	return status, nil
}
