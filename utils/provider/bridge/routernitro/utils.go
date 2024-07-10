package routernitro

import (
	"bytes"
	"context"
	"encoding/json"
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

func (r Routernitro) Quote(ctx context.Context, args QuoteParams) (Quote, error) {
	params := &url.Values{}
	params.Set("fromTokenAddress", StandardizeZeroAddress(args.FromTokenAddress).Hex())
	params.Set("toTokenAddress", StandardizeZeroAddress(args.ToTokenAddress).Hex())
	params.Set("amount", args.AmountWei.String())
	params.Set("fromTokenChainId", strconv.Itoa(args.FromTokenChainId))
	params.Set("toTokenChainId", strconv.Itoa(args.ToTokenChainId))
	params.Set("partnerId", "1")
	params.Set("slippageTolerance", "2")
	params.Set("destFuel", "0")
	var result Quote
	u, _ := url.Parse("https://api-beta.pathfinder.routerprotocol.com/api/v2/quote")
	u.RawQuery = params.Encode()
	if err := utils.Request(ctx, "GET", u.String(), nil, &result); err != nil {
		return result, err
	}
	if err := result.Error(); err != nil {
		return result, err
	}
	return result, nil
}

func (r Routernitro) BuildTx(ctx context.Context, quote Quote, sender, receiver common.Address) (Txn, error) {
	quoteData := utils.ToMap(quote)
	quoteData["receiverAddress"] = receiver.Hex()
	if receiver.Cmp(constant.ZeroAddress) == 0 {
		quoteData["receiverAddress"] = sender.Hex()
	}
	quoteData["senderAddress"] = sender.Hex()
	if quote.BridgeFee.Symbol == "" {
		delete(quoteData, "bridgeFee")
		quoteData["bridgeFee"] = map[string]string{}
	}

	if quote.Source.StableReserveAsset.Symbol == "" {
		data := quoteData["source"].(map[string]interface{})
		delete(data, "stableReserveAsset")
		quoteData["source"] = data
	}

	if quote.Destination.StableReserveAsset.Symbol == "" {
		data := quoteData["destination"].(map[string]interface{})
		delete(data, "stableReserveAsset")
		quoteData["destination"] = data
	}

	var body = bytes.NewBuffer(nil)
	_ = json.NewEncoder(body).Encode(quoteData)
	var result Txn
	err := utils.Request(ctx, "POST",
		"https://api-beta.pathfinder.routerprotocol.com/api/v2/transaction", body, &result)
	if err != nil {
		return Txn{}, err
	}
	if err := result.Error(); err != nil {
		return result, err
	}
	return result, nil
}

func (r Routernitro) GetBestQuote(ctx context.Context, args provider.SwapParams) (tokenInName, tokenInChainName string,
	tokenInAmount decimal.Decimal, quote Quote, err error) {
	log := utils.GetLogFromCtx(ctx).WithField("name", r.Name())
	if args.TargetToken == "" || args.TargetChain == "" {
		return tokenInName, tokenInChainName, tokenInAmount, Quote{}, errors.New("target token or target chain is empty")
	}

	var (
		tokenOut = r.conf.GetTokenInfoOnChain(args.TargetToken, args.TargetChain)
	)

	getQuote := func(chainName, token string) error {
		sourceToken := r.conf.GetTokenInfoOnChain(token, chainName)
		currentLog := log.WithField("TokenIn", sourceToken.Name).WithField("sourceChain", chainName).
			WithField("TargetToken", args.TargetToken).WithField("TargetChain", args.TargetChain)
		currentLog.Debug("start check tokenIn")
		chain := r.conf.GetChainConfig(chainName)
		tokenIn := r.conf.GetTokenInfoOnChain(sourceToken.Name, chainName)
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
		quoteData, err := r.Quote(ctx, QuoteParams{
			FromTokenAddress: common.HexToAddress(tokenIn.ContractAddress),
			ToTokenAddress:   common.HexToAddress(tokenOut.ContractAddress),
			AmountWei:        tokenInTestBalanceWei,
			FromTokenChainId: chain.Id,
			ToTokenChainId:   constant.GetChainId(args.TargetChain),
		})
		if err != nil {
			currentLog.Debugf("get quote error: %s", err)
			return errors.Wrap(err, "get quote")
		}
		currentLog = currentLog.WithField("quote", utils.ToMap(quoteData))

		minimumReceived := chains.WeiToEth(quoteData.Destination.TokenAmount.BigInt(), tokenOut.Decimals)
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

	for _, sourceToken := range r.conf.SourceTokens {
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
	// tokenInAmount 只保留tokenIn.Decimals位小数

	tokenIn := r.conf.GetTokenInfoOnChain(tokenInName, tokenInChainName)
	quoteData, err := r.Quote(ctx, QuoteParams{
		FromTokenAddress: common.HexToAddress(tokenIn.ContractAddress),
		ToTokenAddress:   common.HexToAddress(tokenOut.ContractAddress),
		AmountWei:        decimal.NewFromBigInt(chains.EthToWei(tokenInAmount, tokenIn.Decimals), 0),
		FromTokenChainId: constant.GetChainId(tokenInChainName),
		ToTokenChainId:   constant.GetChainId(args.TargetChain),
	})
	if err != nil {
		return tokenInName, tokenInChainName, tokenInAmount, Quote{}, err
	}
	quote = quoteData
	log.Debugf("best tokenInName: %s, tokenInChainName: %s, tokenInAmount: %s", tokenInName, tokenInChainName, tokenInAmount)
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
	log := utils.GetLogFromCtx(ctx)
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
			log.WithFields(utils.ToMap(status)).Infof("tx %s status: %s", hash, status.Data.FindNitroTransactionByFilter.Status)
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
