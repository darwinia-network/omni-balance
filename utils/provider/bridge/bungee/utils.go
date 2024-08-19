package bungee

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/url"
	"omni-balance/utils"
	"omni-balance/utils/chains"
	"omni-balance/utils/constant"
	"omni-balance/utils/error_types"
	"omni-balance/utils/provider"
	"strconv"
	"strings"
	"time"

	log "omni-balance/utils/logging"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/tidwall/gjson"
)

var (
	pendingStatus = "pending"
)

type QuoteParams struct {
	Sender           common.Address  `json:"sender"`
	Receiver         common.Address  `json:"receiver"`
	FromTokenAddress common.Address  `json:"fromTokenAddress"`
	ToTokenAddress   common.Address  `json:"toTokenAddress"`
	AmountWei        decimal.Decimal `json:"amount"`
	FromTokenChainId int             `json:"fromTokenChainId"`
	ToTokenChainId   int             `json:"toTokenChainId"`
}

func (r Bungee) Quote(ctx context.Context, args QuoteParams) (gjson.Result, error) {
	params := &url.Values{}
	params.Set("fromChainId", strconv.Itoa(args.FromTokenChainId))
	params.Set("fromTokenAddress", StandardizeZeroAddress(args.FromTokenAddress).Hex())
	params.Set("toChainId", strconv.Itoa(args.ToTokenChainId))
	params.Set("toTokenAddress", StandardizeZeroAddress(args.ToTokenAddress).Hex())
	params.Set("fromAmount", args.AmountWei.String())
	params.Set("userAddress", args.Sender.Hex())
	params.Set("sort", "output")
	params.Set("recipient", args.Receiver.Hex())
	params.Set("singleTxOnly", "true")
	u, _ := url.Parse("https://api.socket.tech/v2/quote")
	u.RawQuery = params.Encode()
	data, err := utils.RequestBinary(ctx, "GET", u.String(), nil, "API-KEY", apiKey)
	if err != nil {
		return gjson.Result{}, err
	}
	result := gjson.Parse(string(data))
	if !result.Get("success").Bool() {
		return gjson.Result{}, errors.Errorf("quote error: %s", result.Get("message").String())
	}
	return result, nil
}

func (r Bungee) BuildTx(ctx context.Context, quote gjson.Result, sender, receiver common.Address) (*types.LegacyTx, error) {
	var router = make(map[string]interface{})
	if err := json.Unmarshal([]byte(quote.Raw), &router); err != nil {
		return nil, err
	}
	body := map[string]interface{}{
		"route": router,
	}
	reqData, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	data, err := utils.RequestBinary(ctx, "POST",
		"https://api.socket.tech/v2/build-tx", bytes.NewReader(reqData), "API-KEY", apiKey, "Content-Type", "application/json")
	if err != nil {
		return nil, err
	}
	result := gjson.Parse(string(data))
	if !result.Get("success").Bool() {
		return nil, errors.Errorf("build tx error: %s", result.Get("message").String())
	}
	txData := result.Get("result.txData").String()
	if txData == "" {
		return nil, errors.New("build tx error: txData is empty")
	}
	txTarget := result.Get("result.txTarget").String()
	if txTarget == "" {
		return nil, errors.New("build tx error: txTarget is empty")
	}
	value := big.NewInt(0)
	value.SetString(txTarget, 16)
	to := common.HexToAddress(txTarget)
	return &types.LegacyTx{
		To:    &to,
		Value: value,
		Data:  common.Hex2Bytes(strings.TrimPrefix(txData, "0x")),
	}, nil
}

func (r Bungee) GetBestQuote(ctx context.Context, args provider.SwapParams) (tokenInName, tokenInChainName string,
	tokenInAmount decimal.Decimal, quote gjson.Result, err error) {
	if args.TargetToken == "" || args.TargetChain == "" {
		return tokenInName, tokenInChainName, tokenInAmount, gjson.Result{}, errors.New("target token or target chain is empty")
	}

	var (
		tokenOut = r.conf.GetTokenInfoOnChain(args.TargetToken, args.TargetChain)
		msg      = fmt.Sprintf("wallet %s rebalance %s on %s", args.Receiver, args.TargetToken, args.TargetChain)
	)

	getQuote := func(chainName, tokenName string) {
		sourceToken := r.conf.GetTokenInfoOnChain(tokenName, chainName)
		chain := r.conf.GetChainConfig(chainName)
		tokenIn := r.conf.GetTokenInfoOnChain(sourceToken.Name, chainName)
		if tokenIn.ContractAddress == "" {
			log.Warnf("#%d %s %s tokenIn contract address is empty", args.OrderId, msg, tokenIn.Name)
			return
		}
		client, err := chains.NewTryClient(ctx, chain.RpcEndpoints)
		if err != nil {
			log.Warnf("#%d %s %s get chain %s client error: %s", args.OrderId, msg, tokenIn.Name, chain.Name, err)
			return
		}
		defer client.Close()
		balance, err := chains.GetTokenBalance(ctx, client, tokenIn.ContractAddress,
			args.Sender.GetAddress(true).Hex(), tokenIn.Decimals)
		if err != nil {
			log.Warnf("get %s on %s balance error: %s", tokenName, chainName, err)
			return
		}
		if balance.LessThanOrEqual(decimal.Zero) {
			return
		}

		quoteData, err := r.Quote(ctx, QuoteParams{
			FromTokenAddress: common.HexToAddress(tokenOut.ContractAddress),
			ToTokenAddress:   common.HexToAddress(tokenIn.ContractAddress),
			AmountWei:        decimal.NewFromBigInt(chains.EthToWei(args.Amount, tokenOut.Decimals), 0),
			FromTokenChainId: constant.GetChainId(args.TargetChain),
			ToTokenChainId:   chain.Id,
		})
		if err != nil {
			log.Debugf("#%d %s %s get quote error: %s", args.OrderId, msg, tokenIn.Name, err)
			return
		}

		if !quoteData.Get("success").Bool() {
			log.Warnf("#%d %s %s get quote error: %s", args.OrderId, msg, tokenIn.Name, quoteData.Raw)
			return
		}
		if len(quoteData.Get("result.routes").Array()) == 0 {
			return
		}
		tokenAmount := quoteData.Get("result.routes.0.toAmount").String()
		if tokenAmount == "" {
			log.Warnf("#%d %s %s get quote error: destination token amount is empty", args.OrderId, msg, tokenIn.Name)
			return
		}

		tokenOutAmount := chains.WeiToEth(decimal.RequireFromString(tokenAmount).BigInt(), tokenIn.Decimals)
		// 0.2% slippage + 0.2% fee
		needTokenInAmount := tokenOutAmount.Add(tokenOutAmount.Mul(decimal.RequireFromString("0.004")))

		if needTokenInAmount.GreaterThan(balance) {
			log.Debugf("%s need %s on %s balance is greater than balance, need: %s, balance: %s",
				args.Sender.GetAddress(true).Hex(), tokenName, chainName, needTokenInAmount.String(), balance.String())
			return
		}
		if tokenInAmount.Equal(decimal.Zero) {
			tokenInAmount = needTokenInAmount
		}
		if tokenInAmount.GreaterThan(needTokenInAmount) {
			return
		}

		tokenInAmount = needTokenInAmount
		tokenInName = sourceToken.Name
		tokenInChainName = chainName
		quote = quoteData
	}

	if args.SourceChain != "" && args.SourceToken != "" {
		getQuote(args.SourceChain, args.SourceToken)
		if tokenInChainName == "" || tokenInName == "" || tokenInAmount.IsZero() {
			return "", "", tokenInAmount, gjson.Result{}, error_types.ErrUnsupportedTokenAndChain
		}
		log.Debugf("#%d get best route for %s %s is use %s %s token from %s chain by specify source chain and token", args.OrderId, args.TargetChain, args.TargetToken, tokenInAmount, tokenInName, tokenInChainName)
		return
	}

	if args.SourceToken != "" && len(args.SourceChainNames) > 0 && args.SourceChain == "" {
		for _, v := range args.SourceChainNames {
			getQuote(v, args.SourceToken)
		}
		if tokenInChainName == "" || tokenInName == "" || tokenInAmount.IsZero() {
			return "", "", tokenInAmount, gjson.Result{}, error_types.ErrUnsupportedTokenAndChain
		}
		log.Debugf("#%d best route for %s %s is use %s %s token from %s chain by specify source chains and token", args.OrderId, args.TargetChain, args.TargetToken, tokenInAmount, tokenInName, tokenInChainName)
		return
	}

	for _, sourceToken := range r.conf.SourceTokens {
		if args.SourceToken != "" && sourceToken.Name != args.SourceToken {
			continue
		}
		for _, v := range sourceToken.Chains {
			if args.SourceChain != "" && v != args.SourceChain {
				continue
			}
			getQuote(v, sourceToken.Name)
		}
	}
	if tokenInChainName == "" || tokenInName == "" || tokenInAmount.IsZero() {
		return "", "", tokenInAmount, gjson.Result{}, error_types.ErrUnsupportedTokenAndChain
	}
	log.Debugf("#%d best route for %s %s is use %s %s token from %s chain", args.OrderId, args.TargetChain, args.TargetToken, tokenInAmount, tokenInName, tokenInChainName)
	return
}

func (r Bungee) WaitForTx(ctx context.Context, fromChainId, toChainId int, hash common.Hash) error {
	var (
		t     = time.NewTicker(time.Second * 2)
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
			result, err := r.Status(ctx, fromChainId, toChainId, hash)
			if errors.Is(err, error_types.ErrNotFound) {
				count++
				log.Debugf("tx %s not found, count: %d", hash, count)
				continue
			}
			sourceTxStatus := result.Get("result.sourceTxStatus").String()
			destinationTxStatus := result.Get("result.destinationTxStatus").String()
			log.Infof("tx %s sourceTxStatus: %s; destinationTxStatus: %s",
				hash, sourceTxStatus, destinationTxStatus)
			if destinationTxStatus == "success" && sourceTxStatus == "success" {
				return nil
			}
			if utils.InArrayFold(sourceTxStatus, []string{pendingStatus, "PENDING"}) {
				continue
			}
		}
	}
	return errors.New("wait for tx timeout")
}

func (r Bungee) Status(ctx context.Context, fromChainId, toChainId int, hash common.Hash) (gjson.Result, error) {
	u, _ := url.Parse("https://api.socket.tech/v2/bridge-status")
	query := u.Query()
	query.Set("fromChainId", strconv.Itoa(fromChainId))
	query.Set("toChainId", strconv.Itoa(toChainId))
	query.Set("transactionHash", hash.Hex())
	u.RawQuery = query.Encode()
	data, err := utils.RequestBinary(ctx, "Get", u.String(), nil, "API-KEY", apiKey)
	if err != nil {
		return gjson.Result{}, err
	}
	result := gjson.Parse(string(data))
	if !result.Get("success").Bool() {
		return gjson.Result{}, errors.Errorf("status error: %s", result.Get("message").String())
	}
	return result, nil
}

// StandardizeZeroAddress standardizes the zero address.
// If the provided zeroAddress is equal to constant.ZeroAddress, it returns okxZeroAddress; otherwise, it returns the original address.
// Parameter:
// zeroAddress common.Address - The address to be standardized.
// Return:
// common.Address - The standardized okxZeroAddress
func StandardizeZeroAddress(address common.Address) common.Address {
	if address.Cmp(constant.ZeroAddress) == 0 {
		return ZeroAddress
	}
	return address
}
