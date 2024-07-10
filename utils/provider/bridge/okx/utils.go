package okx

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"io"
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

type QuoteParams struct {
	Amount           decimal.Decimal
	FormChainId      int
	ToChainId        int
	ToTokenAddress   common.Address
	FromTokenAddress common.Address
}

func (o *OKX) Quote(ctx context.Context, args QuoteParams) (QuoteData, error) {
	var params = url.Values{}
	params.Set("amount", args.Amount.String())
	params.Set("fromChainId", strconv.Itoa(args.FormChainId))
	params.Set("toChainId", strconv.Itoa(args.ToChainId))
	params.Set("toTokenAddress", StandardizeZeroAddress(args.ToTokenAddress).Hex())
	params.Set("fromTokenAddress", StandardizeZeroAddress(args.FromTokenAddress).Hex())
	// see https://www.okx.com/zh-hans/web3/build/docs/waas/dex-get-route-information
	params.Set("slippage", "0.01")
	var result OksResp
	if err := o.request(ctx, "GET", "/api/v5/dex/cross-chain/quote", params, nil, &result); err != nil {
		return QuoteData{}, errors.Wrap(err, "quote")
	}
	return result.Quote()
}

func (o *OKX) request(ctx context.Context, method string, path string, params url.Values, body map[string]interface{}, dest *OksResp) error {
	signature, timestamp := createSignature(method, path, params, body, o.Key.SecretKey)
	headers := map[string]string{
		"OK-ACCESS-KEY":        o.Key.ApiKey,
		"OK-ACCESS-SIGN":       signature,
		"OK-ACCESS-TIMESTAMP":  timestamp,
		"OK-ACCESS-PASSPHRASE": o.Key.Passphrase,
		"OK-ACCESS-PROJECT":    o.Key.Project,
		"content-type":         "application/json",
	}
	var headersList []string
	for k, v := range headers {
		headersList = append(headersList, k, v)
	}
	var bodyReader io.Reader
	if len(body) != 0 {
		data, _ := json.Marshal(body)
		bodyReader = bytes.NewReader(data)
	}
	u, _ := url.Parse("https://www.okx.com")
	u.RawPath = path
	u.Path = path
	u.RawQuery = params.Encode()
	utils.GetLogFromCtx(ctx).Debugf("request url: %s", u.String())
	if err := utils.Request(ctx, method, u.String(), bodyReader, dest, headersList...); err != nil {
		return err
	}
	return dest.Error()
}

func (o *OKX) GetBestTokenInChain(ctx context.Context, args provider.SwapParams) (tokenInName, tokenInChainName string, tokenInAmount decimal.Decimal, err error) {
	log := utils.GetLogFromCtx(ctx)
	if args.TargetToken == "" || args.TargetChain == "" {
		return "", "", tokenInAmount, errors.New("target token or target chain is empty")
	}
	var (
		tokenOut = o.conf.GetTokenInfoOnChain(args.TargetToken, args.TargetChain)
	)
	getQuote := func(chainName, tokenName string) error {
		sourceToken := o.conf.GetTokenInfoOnChain(tokenName, chainName)
		currentLog := log.WithField("TokenIn", sourceToken.Name).WithField("sourceChain", chainName).
			WithField("TargetToken", args.TargetToken).WithField("TargetChain", args.TargetChain)
		currentLog.Debug("start check tokenIn")
		chain := o.conf.GetChainConfig(chainName)
		tokenIn := o.conf.GetTokenInfoOnChain(sourceToken.Name, chainName)
		if tokenIn.ContractAddress == "" {
			return errors.Errorf("token %s on chain %s is not supported", sourceToken.Name, chainName)
		}
		client, err := chains.NewTryClient(ctx, chain.RpcEndpoints)
		if err != nil {
			currentLog.Debugf("new client error: %s", err)
			return err
		}

		tokenInTestBalance := decimal.RequireFromString("1")
		tokenInTestBalanceWei := decimal.NewFromBigInt(chains.EthToWei(tokenInTestBalance, tokenIn.Decimals), 0)
		quote, err := o.Quote(ctx, QuoteParams{
			Amount:           tokenInTestBalanceWei,
			FormChainId:      chain.Id,
			ToChainId:        constant.GetChainId(args.TargetChain),
			ToTokenAddress:   common.HexToAddress(tokenOut.ContractAddress),
			FromTokenAddress: common.HexToAddress(tokenIn.ContractAddress),
		})
		if err != nil {
			currentLog.Debugf("get quote error: %s", err)
			return err
		}
		if len(quote.RouterList) == 0 {
			currentLog.Debugf("no router list")
			return errors.Errorf("no router list")
		}
		currentLog = currentLog.WithField("quote", utils.ToMap(quote))
		currentLog.Debug("get quote success")

		minimumReceived := chains.WeiToEth(quote.RouterList[0].MinimumReceived.BigInt(), tokenOut.Decimals)
		needBalance := tokenInTestBalance.Div(minimumReceived).Mul(args.Amount)

		balance, err := chains.GetTokenBalance(ctx, client, tokenIn.ContractAddress,
			args.Sender.GetAddress(true).Hex(), tokenIn.Decimals)
		if err != nil {
			currentLog.Debugf("get balance error: %s", err)
			return err
		}

		log.Debugf("need %s balance: %s, wallet %s balance: %s on %s",
			tokenIn.Name, needBalance, tokenIn.Name, balance, chainName)
		if needBalance.GreaterThan(balance) {
			currentLog.Debugf("%s need balance: %s, balance: %s", tokenIn.Name, needBalance, balance)
			return errors.Errorf("need %s balance: %s, balance: %s", tokenIn.Name, needBalance, balance)
		}
		if tokenInAmount.Equal(decimal.Zero) {
			tokenInAmount = needBalance
		}
		if tokenInAmount.GreaterThan(needBalance) {
			currentLog.Debugf("need balance: %s, balance: %s", needBalance, balance)
			return errors.Errorf("need %s balance: %s, balance: %s", tokenIn.Name, needBalance, balance)
		}
		tokenInAmount = needBalance
		tokenInName = sourceToken.Name
		tokenInChainName = chainName
		return nil
	}
	if args.SourceChain != "" && args.SourceToken != "" {
		if err := getQuote(args.SourceChain, args.SourceToken); err != nil {
			return "", "", tokenInAmount, err
		}
	}

	for _, sourceToken := range o.conf.SourceTokens {
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
		return "", "", tokenInAmount, error_types.ErrUnsupportedTokenAndChain
	}
	log.Debugf("best tokenInName: %s, tokenInChainName: %s, tokenInAmount: %s", tokenInName, tokenInChainName, tokenInAmount)
	return

}

func (o *OKX) approveTransaction(ctx context.Context, chainId int, tokenAddress common.Address, amountWei decimal.Decimal) (ApproveTransaction, error) {
	params := &url.Values{}
	params.Set("chainId", strconv.Itoa(chainId))
	params.Set("tokenContractAddress", tokenAddress.Hex())
	params.Set("approveAmount", amountWei.String())
	var result OksResp
	if err := o.request(ctx, "GET", "/api/v5/dex/aggregator/approve-transaction", *params, nil, &result); err != nil {
		return ApproveTransaction{}, err
	}
	return result.ApproveTransaction()
}

func (o *OKX) buildTx(ctx context.Context, args provider.SwapParams) (BuildTxData, error) {
	var (
		tokenIn       = o.conf.GetTokenInfoOnChain(args.SourceToken, args.SourceChain)
		tokenOut      = o.conf.GetTokenInfoOnChain(args.TargetToken, args.TargetChain)
		tokenInChain  = o.conf.GetChainConfig(args.SourceChain)
		tokenOutChain = o.conf.GetChainConfig(args.TargetChain)
		amountWei     = chains.EthToWei(args.Amount, tokenIn.Decimals)
	)
	if tokenIn.Name == "" ||
		tokenOut.Name == "" ||
		tokenInChain.Name == "" ||
		tokenOutChain.Name == "" ||
		tokenIn.ContractAddress == "" ||
		tokenOut.ContractAddress == "" {
		return BuildTxData{}, errors.New("token info is empty")
	}
	params := &url.Values{}
	params.Set("fromChainId", strconv.Itoa(tokenInChain.Id))
	params.Set("toChainId", strconv.Itoa(tokenOutChain.Id))
	params.Set("fromTokenAddress", StandardizeZeroAddress(common.HexToAddress(tokenIn.ContractAddress)).Hex())
	params.Set("toTokenAddress", StandardizeZeroAddress(common.HexToAddress(tokenOut.ContractAddress)).Hex())
	params.Set("amount", amountWei.String())
	params.Set("slippage", "0.01")
	params.Set("sort", "1")
	params.Set("userWalletAddress", args.Sender.GetAddress(true).Hex())
	if args.Receiver != "" && !strings.EqualFold(args.Receiver, args.Sender.GetAddress(true).Hex()) {
		params.Set("receiveAddress", args.Receiver)
	}
	var tx OksResp
	if err := o.request(ctx, "GET", "/api/v5/dex/cross-chain/build-tx", *params, nil, &tx); err != nil {
		return BuildTxData{}, err
	}
	return tx.BuildTx()
}

func (o *OKX) TxStatus(ctx context.Context, hash common.Hash, sourceChainId int) (Status, error) {
	params := &url.Values{}
	params.Set("hash", hash.Hex())
	params.Set("chainId", strconv.Itoa(sourceChainId))
	var txStatus OksResp
	if err := o.request(ctx, "GET", "/api/v5/dex/cross-chain/status", *params, nil, &txStatus); err != nil {
		return Status{}, err
	}
	return txStatus.TxStatus()
}

func (o *OKX) WaitForTx(ctx context.Context, hash common.Hash, sourceChainId int) error {
	log := utils.GetLogFromCtx(ctx).WithFields(logrus.Fields{
		"chain_id": sourceChainId,
		"hash":     hash,
	})
	var (
		t     = time.NewTicker(time.Second * 2)
		count = 0
	)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return errors.New("context done")
		case <-t.C:
			s, err := o.TxStatus(ctx, hash, sourceChainId)
			if err != nil {
				return err
			}
			log = log.WithField("okx_ts_status", s.ToMap())
			log.Infof("tx status: %s, detail status: %s", s.Status, okxWaitForTxStatus[s.DetailStatus])
			switch s.Status {
			case "PENDING":
				count = 0
				continue
			case "SUCCESS":
				return nil
			case "FAILURE", "REFUND":
				return errors.New(okxWaitForTxStatus[s.DetailStatus])
			default:
				if count >= 600 {
					return errors.Errorf("unknown status: %s, detail status: %s", s.Status, s.DetailStatus)
				}
				count++
			}
		}
	}
}

func sign(message, secretKey string) string {
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func createSignature(method, requestPath string, params url.Values, body map[string]interface{}, secretKey string) (string, string) {
	timestamp := time.Now().UTC().Format(time.RFC3339)
	message := preHash(timestamp, method, requestPath, params, body)
	signature := sign(message, secretKey)
	return signature, timestamp
}

func preHash(timestamp, method, requestPath string, params url.Values, body map[string]interface{}) string {
	var queryString string
	if params != nil {
		queryString = "?" + params.Encode()
	}
	if body != nil {
		data, _ := json.Marshal(body)
		queryString = string(data)
	}
	return timestamp + method + requestPath + queryString
}
