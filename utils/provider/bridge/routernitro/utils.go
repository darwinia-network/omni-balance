package routernitro

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/tidwall/gjson"
)

var (
	pendingStatus   = "pending"
	completedStatus = "completed"
)

type QuoteParams struct {
	FromTokenAddress common.Address  `json:"fromTokenAddress"`
	ToTokenAddress   common.Address  `json:"toTokenAddress"`
	AmountWei        decimal.Decimal `json:"amount"`
	FromTokenChainId int             `json:"fromTokenChainId"`
	ToTokenChainId   int             `json:"toTokenChainId"`
}

func (r Routernitro) Quote(ctx context.Context, args QuoteParams) (gjson.Result, error) {
	params := &url.Values{}
	params.Set("fromTokenAddress", StandardizeZeroAddress(args.FromTokenAddress).Hex())
	params.Set("toTokenAddress", StandardizeZeroAddress(args.ToTokenAddress).Hex())
	params.Set("amount", args.AmountWei.String())
	params.Set("fromTokenChainId", strconv.Itoa(args.FromTokenChainId))
	params.Set("toTokenChainId", strconv.Itoa(args.ToTokenChainId))
	params.Set("partnerId", "1")
	params.Set("slippageTolerance", "2")
	params.Set("destFuel", "0")
	u, _ := url.Parse("https://api-beta.pathfinder.routerprotocol.com/api/v2/quote")
	u.RawQuery = params.Encode()
	data, err := utils.RequestBinary(ctx, "GET", u.String(), nil)
	if err != nil {
		return gjson.Result{}, err
	}
	return gjson.Parse(string(data)), nil
}

func (r Routernitro) BuildTx(ctx context.Context, quote gjson.Result, sender, receiver common.Address) (Txn, error) {
	var (
		result Txn
		body   = make(map[string]interface{})
	)
	if err := json.Unmarshal([]byte(quote.String()), &body); err != nil {
		return Txn{}, err
	}
	body["receiverAddress"] = receiver.Hex()
	body["senderAddress"] = sender.Hex()
	data, err := json.Marshal(body)
	if err != nil {
		return Txn{}, err
	}

	err = utils.Request(ctx, "POST",
		"https://api-beta.pathfinder.routerprotocol.com/api/v2/transaction", bytes.NewReader(data), &result)
	if err != nil {
		return Txn{}, err
	}
	if err := result.Error(); err != nil {
		return result, err
	}
	return result, nil
}

func (r Routernitro) GetBestQuote(ctx context.Context, args provider.SwapParams) (tokenInName, tokenInChainName string,
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

		tokenAmount := quoteData.Get("destination").Get("tokenAmount").String()
		if tokenAmount == "" {
			log.Warnf("#%d %s %s get quote error: destination token amount is empty: %s", args.OrderId, msg, tokenIn.Name, quoteData.Raw)
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

	wallet := r.conf.GetWalletConfig(args.Sender.GetAddress().Hex())
	for _, sourceToken := range r.conf.SourceTokens {
		if args.SourceToken != "" && sourceToken.Name != args.SourceToken {
			continue
		}
		if wallet.Mode.IsBalance() && !strings.EqualFold(sourceToken.Name, args.TargetToken) {
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

type Status struct {
	Data struct {
		FindNitroTransactionByFilter struct {
			UsdcValue         string      `json:"usdc_value"`
			UpdateTxHash      interface{} `json:"update_tx_hash"`
			Status            string      `json:"status"`
			SrcTxHash         string      `json:"src_tx_hash"`
			SrcTimestamp      int         `json:"src_timestamp"`
			SrcSymbol         interface{} `json:"src_symbol"`
			SrcChainId        string      `json:"src_chain_id"`
			SrcAmount         string      `json:"src_amount"`
			SrcAddress        string      `json:"src_address"`
			SrcStableSymbol   string      `json:"src_stable_symbol"`
			SrcStableAmount   string      `json:"src_stable_amount"`
			SenderAddress     string      `json:"sender_address"`
			RefundOutboundId  string      `json:"refund_outbound_id"`
			ReceiverAddress   string      `json:"receiver_address"`
			Message           interface{} `json:"message"`
			HasMessage        bool        `json:"has_message"`
			FlowType          string      `json:"flow_type"`
			FundPaidConfirmed bool        `json:"fund_paid_confirmed"`
			FeeSymbol         string      `json:"fee_symbol"`
			FeeAmount         string      `json:"fee_amount"`
			FeeAddress        string      `json:"fee_address"`
			DepositId         string      `json:"deposit_id"`
			DepositInfoType   string      `json:"deposit_info_type"`
			DestTxHash        string      `json:"dest_tx_hash"`
			DestTimestamp     int         `json:"dest_timestamp"`
			DestSymbol        string      `json:"dest_symbol"`
			DestChainId       string      `json:"dest_chain_id"`
			DestAmount        string      `json:"dest_amount"`
			DestStableSymbol  string      `json:"dest_stable_symbol"`
			DestStableAmount  string      `json:"dest_stable_amount"`
			DestAddress       string      `json:"dest_address"`
			NativeTokenAmount string      `json:"native_token_amount"`
			NativeTokenSymbol string      `json:"native_token_symbol"`
		} `json:"findNitroTransactionByFilter"`
	} `json:"data"`
}

func (r Routernitro) WaitForTx(ctx context.Context, hash common.Hash) error {
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
			status, err := r.Status(ctx, hash)
			if errors.Is(err, error_types.ErrNotFound) {
				count++
				log.Debugf("tx %s not found, count: %d", hash, count)
				continue
			}
			log.Infof("tx %s status: %s", hash, status.Data.FindNitroTransactionByFilter.Status)
			switch status.Data.FindNitroTransactionByFilter.Status {
			case pendingStatus:
				count = 0
				continue
			case completedStatus:
				return nil
			default:
				count++
				continue
			}
		}
	}
	return errors.New("wait for tx timeout")
}

func (r Routernitro) Status(ctx context.Context, hash common.Hash) (Status, error) {
	var query = map[string]interface{}{
		"query": "query Transaction($hash: String!) {\n  findNitroTransactionByFilter(hash: $hash) {\n    usdc_value\n    update_tx_hash\n    status\n    src_tx_hash\n    src_timestamp\n    src_symbol\n    src_chain_id\n    src_amount\n    src_address\n    src_stable_symbol\n    src_stable_amount\n    sender_address\n    refund_outbound_id\n    receiver_address\n    message\n    has_message\n    flow_type\n    fund_paid_confirmed\n    fee_symbol\n    fee_amount\n    fee_address\n    deposit_id\n    deposit_info_type\n    dest_tx_hash\n    dest_timestamp\n    dest_symbol\n    dest_chain_id\n    dest_amount\n    dest_stable_symbol\n    dest_stable_amount\n    dest_address\n    native_token_amount\n    native_token_symbol\n  }\n}",
		"variables": map[string]string{
			"hash": hash.Hex(),
		},
	}
	var body = &bytes.Buffer{}
	if err := json.NewEncoder(body).Encode(query); err != nil {
		return Status{}, err
	}
	var status Status
	err := utils.Request(ctx, "POST", "https://api.pro-nitro-explorer.routernitro.com/graphql", body, &status)
	if err != nil {
		return Status{}, err
	}
	if status.Data.FindNitroTransactionByFilter.SrcTxHash == "" {
		return Status{}, error_types.ErrNotFound
	}
	return status, nil
}

// StandardizeZeroAddress standardizes the zero address.
// If the provided zeroAddress is equal to constant.ZeroAddress, it returns okxZeroAddress; otherwise, it returns the original address.
// Parameter:
// zeroAddress common.Address - The address to be standardized.
// Return:
// common.Address - The standardized okxZeroAddress
func StandardizeZeroAddress(address common.Address) common.Address {
	if address.Cmp(constant.ZeroAddress) == 0 {
		return routernitroZeroAddress
	}
	return address
}
