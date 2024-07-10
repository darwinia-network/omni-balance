package gate

import (
	"context"
	"encoding/json"
	"fmt"

	"net/http"
	"omni-balance/utils"
	"omni-balance/utils/chains"
	"omni-balance/utils/constant"
	"omni-balance/utils/error_types"
	"omni-balance/utils/provider"
	"strings"
	"time"

	"github.com/antihax/optional"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gateio/gateapi-go/v6"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

var (
	gateChain2StandardName = map[string]string{
		"ETH":     constant.Ethereum,
		"ARBEVM":  constant.Arbitrum,
		"BSC":     constant.Bsc,
		"KAVAEVM": constant.Kava,
		"MATIC":   constant.Polygon,
		"OPETH":   constant.Optimism,
		"CELO":    constant.Celo,
		"OPBNB":   constant.OpBNB,
		"BNB":     constant.Bnb,
	}
	gateTokenName2StandardName = map[string]string{
		"ARBEVM": "ARB",
	}
)

type Ticker struct {
	Low                      string  `json:"low"`
	Volume                   string  `json:"volume"`
	Last                     string  `json:"last"`
	Open                     string  `json:"open"`
	Deal                     string  `json:"deal"`
	Close                    string  `json:"close"`
	Change                   string  `json:"change"`
	High                     string  `json:"high"`
	Result                   string  `json:"result"`
	Avg                      float64 `json:"avg"`
	RateChangePercentage     string  `json:"rate_change_percentage"`
	RateChangePercentageUtc0 int     `json:"rate_change_percentage_utc0"`
	RateChangePercentageUtc8 int     `json:"rate_change_percentage_utc8"`
}

type TickerResult struct {
	TokenName string
	Price     decimal.Decimal
}

func TokenName2GateTokenName(tokenName string) string {
	for k, v := range gateTokenName2StandardName {
		if strings.EqualFold(tokenName, v) {
			return k
		}
	}
	return tokenName
}

func ChainName2GateChainName(chainName string) string {
	for k, v := range gateChain2StandardName {
		if strings.EqualFold(chainName, v) {
			return k
		}
	}
	return chainName
}

func GateChainName2StandardName(chainName string) string {
	for k, v := range gateChain2StandardName {
		if strings.EqualFold(chainName, k) {
			return v
		}
	}
	return chainName
}

func (g *Gate) ticker(pairs string) (TickerResult, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://data.gateapi.io/api/1/ticker/%s", pairs), nil)
	if err != nil {
		return TickerResult{}, err
	}
	req.Header.Set("content-type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return TickerResult{}, err
	}
	defer resp.Body.Close()
	var result = new(Ticker)
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return TickerResult{}, err
	}
	if !strings.EqualFold(result.Result, "true") {
		return TickerResult{}, errors.Errorf("pairs %s not found", pairs)
	}
	pairsList := strings.Split(pairs, "_")
	return TickerResult{
		TokenName: pairsList[1],
		Price:     decimal.RequireFromString(result.Last),
	}, nil
}

func (g *Gate) Tickers(pairs ...string) (result []TickerResult, err error) {
	for _, v := range pairs {
		ticker, err := g.ticker(v)
		if err != nil {
			return nil, errors.Wrap(err, "ticker")
		}
		result = append(result, ticker)
	}
	return result, nil
}

func (g *Gate) GetTokenOutChain(ctx context.Context, token string) ([]string, error) {
	chains, _, err := g.client.WalletApi.ListCurrencyChains(ctx, token)
	if err != nil {
		return nil, errors.Wrap(err, "list currency chains")
	}
	var supportedChains []string
	for _, v := range chains {
		if _, ok := gateChain2StandardName[v.Chain]; !ok {
			logrus.Debugf("chain %s not supported in %s", v.Chain, token)
			continue
		}
		if v.IsWithdrawDisabled == 1 {
			continue
		}
		supportedChains = append(supportedChains, gateChain2StandardName[v.Chain])
	}
	if len(supportedChains) == 0 {
		return nil, error_types.ErrUnsupportedTokenAndChain
	}
	return supportedChains, nil
}

func (g *Gate) buyToken(ctx context.Context, tokenIn provider.TokenInCost, targetToken string, amount decimal.Decimal,
	f func(gateapi.Order) bool) (order gateapi.Order, err error) {

	currencyPair := utils.GetCurrencyPair(targetToken, "_", tokenIn.TokenName)
	order = gateapi.Order{
		Text:         "t-ob",
		CurrencyPair: currencyPair,
		Type:         "limit",
		Side:         "buy",
		Price:        tokenIn.CostAmount.Div(amount).String(),
		Amount:       amount.String(),
	}
	logrus.WithFields(logrus.Fields{
		"order": order,
	}).Debugf("buy token")

	// place order
	order, _, err = g.client.SpotApi.CreateOrder(ctx, order)
	if err != nil {
		return order, errors.Wrap(err, "create order")
	}
	if order.Status == "cancelled" {
		return order, errors.Errorf("order status is %s", order.Status)
	}
	// wait for order filled
	for order.Status == "open" {
		order, _, err = g.client.SpotApi.GetOrder(ctx, order.Id, fmt.Sprintf("%s_%s", tokenIn.TokenName, targetToken), nil)
		if err != nil {
			return order, errors.Wrap(err, "get order")
		}
		if order.Status == "closed" {
			break
		}
		if f(order) {
			break
		}
		time.Sleep(time.Second)
	}
	if order.Status == "cancelled" {
		return order, errors.Errorf("order status is %s", order.Status)
	}
	return order, nil
}

// IsVerifiedAddress checks whether the given address is a verified address on Gate.io for a specific token and chain.
// Parameters:
//
//	ctx - The context for the request.
//	address - The address to check.
//	tokenName - The name of the token associated with the address.
//	chainName - The name of the blockchain where the address is located.
//
// Returns:
//
//	A boolean indicating if the address is verified.
//	An error if the operation fails.
func (g *Gate) IsVerifiedAddress(ctx context.Context, address, tokenName, chainName string) (bool, error) {
	verifiedAddress, _, err := g.client.WalletApi.ListSavedAddress(ctx, tokenName,
		&gateapi.ListSavedAddressOpts{Chain: optional.NewString(ChainName2GateChainName(chainName))})
	if err != nil {
		return false, errors.Wrap(err, "list saved address")
	}
	if len(verifiedAddress) == 0 {
		return false, errors.Errorf("%s not in gate.io verified address for %s on %s", address, tokenName, chainName)
	}
	var (
		isSenderInVerifiedAddress bool
	)
	for _, v := range verifiedAddress {
		if v.Verified == "0" {
			continue
		}
		if !strings.EqualFold(v.Address, address) ||
			!strings.EqualFold(v.Chain, ChainName2GateChainName(chainName)) {
			continue
		}
		isSenderInVerifiedAddress = true
		break
	}
	return isSenderInVerifiedAddress, nil
}

// GetDepositAddress retrieves the deposit address for a token.
// This method calls the GateAPI's GetDepositAddress interface, using tokenName and chainName to filter for valid deposit addresses.
// Parameters:
//
//	ctx: Context object for canceling requests or other operations.
//	tokenName: The name of the token.
//	chainName: The name of the blockchain. If chainName is empty, all deposit addresses for the specified token are returned.
//
// Returns:
//
//	[]gateapi.MultiChainAddressItem: A list of deposit addresses that meet the specified criteria.
//	error: If an error occurs while fetching the address, an error message is returned.
func (g *Gate) GetDepositAddress(ctx context.Context, tokenName, chainName string) ([]gateapi.MultiChainAddressItem, error) {
	depositInfo, _, err := g.client.WalletApi.GetDepositAddress(ctx, TokenName2GateTokenName(tokenName))
	if err != nil {
		return nil, errors.Wrap(err, "get deposit address")
	}
	var result []gateapi.MultiChainAddressItem
	for _, v := range depositInfo.MultichainAddresses {
		if v.ObtainFailed != 0 {
			continue
		}
		if chainName != "" && !strings.EqualFold(GateChainName2StandardName(v.Chain), chainName) {
			continue
		}
		result = append(result, gateapi.MultiChainAddressItem{
			Chain:        GateChainName2StandardName(v.Chain),
			Address:      v.Address,
			PaymentId:    v.PaymentId,
			PaymentName:  v.PaymentName,
			ObtainFailed: v.ObtainFailed,
		})
	}
	if len(result) == 0 {
		return nil, errors.Errorf("no deposit address for %s", tokenName)
	}
	return result, nil
}

// CheckDeposit checks if the deposit for a specific token and chain is completed.
// Parameters:
//
//	ctx: Context object for canceling requests or other operations.
//	tokenName: The name of the token.
//	chainName: The name of the chain.
//	tx: The transaction hash to look for.
//
// Returns:
//
//	bool: Indicates if the deposit is done.
//	string: The status of the deposit ("DONE", "FAIL", or others).
//	error: An error if encountered, otherwise nil.
func (g *Gate) CheckDeposit(ctx context.Context, tokenName, chainName string, tx common.Hash) (bool, string, error) {
	deposits, _, err := g.client.WalletApi.ListDeposits(ctx, &gateapi.ListDepositsOpts{
		Currency: optional.NewString(TokenName2GateTokenName(tokenName))})
	if err != nil {
		return false, "", errors.Wrap(err, "list deposits")
	}
	if len(deposits) == 0 {
		return false, "", error_types.ErrNotFound
	}

	for _, v := range deposits {
		if !strings.EqualFold(gateChain2StandardName[v.Chain], chainName) ||
			!strings.EqualFold(tx.Hex(), v.Txid) {
			continue
		}
		switch v.Status {
		case "FAIL":
			return false, v.Status, errors.Errorf("deposit failed status: %s", v.Status)
		}
		return strings.EqualFold(v.Status, "DONE"), v.Status, nil
	}
	return false, "", error_types.ErrNotFound
}

func ProvideLiquidityAction2Int() {

}

func (g *Gate) ProvideLiquidity(ctx context.Context, args provider.SwapParams) (provider.SwapResult, error) {
	if args.SourceChain == "" || args.SourceToken == "" || args.TargetChain == "" || args.TargetToken == "" {
		err := errors.New("source chain, source token, target chain, target token can not be empty")
		return provider.SwapResult{
			Error:        err.Error(),
			ProviderType: g.Type(),
			ProviderName: g.Name(),
			Status:       provider.TxStatusFailed,
		}, err
	}

	var (
		recordFn = func(s provider.SwapHistory, errs ...error) {
			s.ProviderType = string(g.Type())
			s.ProviderName = g.Name()
			s.Amount = args.Amount
			if args.RecordFn == nil {
				return
			}
			args.RecordFn(s, errs...)
		}
		log    = args.GetLogs(g.Name() + "_ProvideLiquidity")
		wallet = args.Sender
		sr     = new(provider.SwapResult).
			SetTokenInName(args.SourceToken).
			SetTokenInChainName(args.SourceChain).
			SetProviderName(g.Name()).
			SetProviderType(g.Type()).
			SetCurrentChain(args.SourceChain).
			SetTx(args.LastHistory.Tx).
			SetReciever(wallet.GetAddress().Hex())
		sh = &provider.SwapHistory{
			ProviderName: g.Name(),
			ProviderType: string(g.Type()),
			Amount:       args.Amount,
			CurrentChain: args.SourceChain,
			Tx:           args.LastHistory.Tx,
		}
		isActionSuccess = args.LastHistory.Status == provider.TxStatusSuccess.String()
		actionNumber    = ProvideLiquidityAction2Number(args.LastHistory.Actions)
		tx              = args.LastHistory.Tx
		gateTokenName   = TokenName2GateTokenName(args.TargetToken)
		gateTargetChain = ChainName2GateChainName(args.TargetChain)
		// gateSourceChain = ChainName2GateChainName(args.SourceChain)
	)

	if args.LastHistory.Actions == ProvideLiquidityActionTargetChainWithdrawComplete {
		return sr.SetStatus(provider.TxStatusSuccess).Out(), nil
	}

	isSenderInVerifiedAddress, err := g.IsVerifiedAddress(ctx, wallet.GetAddress(true).Hex(), gateTokenName, gateTargetChain)
	if err != nil {
		err = errors.Wrap(err, "check sender address")
		return sr.SetError(err).Out(), err
	}
	if !isSenderInVerifiedAddress {
		err = errors.Errorf("%s not in gate.io verified address for %s",
			args.Sender, args.TargetToken)
		return sr.SetError(err).Out(), err
	}
	var (
		sourceChain = g.config.GetChainConfig(args.SourceChain)
		targetChain = g.config.GetChainConfig(args.TargetChain)
	)
	sourceClient, err := chains.NewTryClient(ctx, sourceChain.RpcEndpoints)
	if err != nil {
		err = errors.Wrap(err, "new try client")
		return sr.SetError(err).Out(), err
	}
	defer sourceClient.Close()
	targetClient, err := chains.NewTryClient(ctx, targetChain.RpcEndpoints)
	if err != nil {
		err = errors.Wrap(err, "new try client")
		return sr.SetError(err).Out(), err
	}
	defer targetClient.Close()

	if actionNumber <= 1 && !isActionSuccess {
		items, err := g.GetDepositAddress(ctx, args.SourceToken, args.SourceChain)
		if err != nil {
			err = errors.Wrap(err, "get deposit address")
			return sr.SetError(err).Out(), err
		}
		if len(items) == 0 {
			err = errors.Errorf("no deposit address for %s", args.SourceToken)
			return sr.SetError(err).Out(), err
		}
		depositAddress := common.HexToAddress(items[0].Address)
		if depositAddress.Cmp(constant.ZeroAddress) == 0 {
			err = errors.Errorf("deposit address is empty")
			return sr.SetError(err).Out(), err
		}

		token := g.config.GetTokenInfoOnChain(args.SourceToken, args.SourceChain)
		txData, err := chains.BuildSendToken(ctx, chains.SendTokenParams{
			Client:        sourceClient,
			Sender:        args.Sender.GetAddress(true),
			TokenAddress:  common.HexToAddress(token.ContractAddress),
			TokenDecimals: token.Decimals,
			ToAddress:     depositAddress,
			AmountWei:     decimal.NewFromBigInt(chains.EthToWei(args.Amount, token.Decimals), 0),
		})
		if err != nil {
			err = errors.Wrap(err, "build send token tx")
			return sr.SetError(err).Out(), err
		}
		recordFn(sh.SetActions(ProvideLiquidityActionSourceChainTxsendingAction).SetStatus(provider.TxStatusPending).Out())
		txHash, err := args.Sender.SendTransaction(ctx, txData, sourceClient)
		if err != nil {
			err = errors.Wrap(err, "send transaction")
			recordFn(sh.SetActions(ProvideLiquidityActionSourceChainTxsendingAction).SetStatus(provider.TxStatusFailed).Out(), err)
			return sr.SetError(err).Out(), err
		}
		tx = txHash.Hex()
		sh = sh.SetTx(txHash.Hex())
		recordFn(*sh.SetActions(ProvideLiquidityActionSourceChainTxsendingAction).SetStatus(provider.TxStatusSuccess))
		sr = sr.SetTx(txHash.Hex())
	}

	if tx == "" {
		// back to step 1
		log.Warnf("tx is empty, back to step 2")
		recordFn(sh.SetTx("").SetActions(ProvideLiquidityActionSourceChainTxComplete).SetStatus(provider.TxStatusFailed).Out())
		return sr.SetError(errors.Errorf("tx is empty")).Out(), errors.Errorf("tx is empty")
	}

	sh = sh.SetTx(tx)
	sr = sr.SetTx(tx)

	if actionNumber <= 2 && !isActionSuccess {
		recordFn(sh.SetActions(ProvideLiquidityActionSourceChainTxComplete).SetStatus(provider.TxStatusPending).Out())
		if err := args.Sender.WaitTransaction(ctx, common.HexToHash(tx), sourceClient); err != nil {
			recordFn(sh.SetActions(ProvideLiquidityActionSourceChainTxComplete).SetStatus(provider.TxStatusFailed).Out(), err)
			return sr.SetError(err).Out(), err
		}
		recordFn(sh.SetActions(ProvideLiquidityActionSourceChainTxComplete).SetStatus(provider.TxStatusSuccess).Out())
	}

	if actionNumber <= 3 && !isActionSuccess {
		log.Debugf("start wait for %s deposit from %s to gate, tx is %s", args.SourceToken, args.SourceChain, tx)
		recordFn(sh.SetActions(ProvideLiquidityActionGateRechargeTxWaiting).SetStatus(provider.TxStatusPending).Out())
		err := g.WaitForDeposit(ctx, args.TargetToken, args.SourceChain, common.HexToHash(tx), log)
		if err != nil {
			recordFn(sh.SetActions(ProvideLiquidityActionGateRechargeTxWaiting).SetStatus(provider.TxStatusFailed).Out(), err)
			return sr.SetError(err).Out(), err
		}
		recordFn(sh.SetActions(ProvideLiquidityActionGateRechargeTxComplete).SetStatus(provider.TxStatusSuccess).Out())
	}

	withdrawOrderId := args.LastHistory.Tx

	if actionNumber <= 5 && !isActionSuccess {
		log.Debugf("start %s withdraw from gate to %s", args.SourceToken, args.TargetChain)
		recordFn(sh.SetActions(ProvideLiquidityActionGateWithdrawSending).SetStatus(provider.TxStatusPending).Out())
		withdrawOrder, _, err := g.client.WithdrawalApi.Withdraw(ctx, gateapi.LedgerRecord{
			Amount:   args.Amount.String(),
			Currency: gateTokenName,
			Address:  wallet.GetAddress(true).Hex(),
			Chain:    gateTargetChain,
		})
		if err != nil {
			err = errors.Wrap(err, "withdraw from gate")
			recordFn(sh.SetActions(ProvideLiquidityActionGateWithdrawSending).SetStatus(provider.TxStatusFailed).Out(), err)
			return sr.SetError(err).Out(), err
		}
		withdrawOrderId = withdrawOrder.Id
		recordFn(sh.SetTx(withdrawOrderId).SetActions(ProvideLiquidityActionGateWithdrawSending).SetStatus(provider.TxStatusSuccess).Out())
	}

	if withdrawOrderId == "" {
		// back to step 5
		log.Warnf("withdraw order id is empty, back to step 5")
		err = errors.Errorf("withdraw order id is empty")
		recordFn(sh.SetTx("").SetActions(ProvideLiquidityActionGateWithdrawSending).SetStatus(provider.TxStatusFailed).Out(), err)
		return sr.SetError(err).Out(), err
	}

	sh = sh.SetTx(tx)
	sr = sr.SetTx(tx)

	if actionNumber <= 6 && !isActionSuccess {
		recordFn(sh.SetActions(ProvideLiquidityActionGateWithdrawWaiting).SetStatus(provider.TxStatusPending).Out())
		err := g.WaitForWithdrawal(ctx, withdrawOrderId, args.SourceToken, args.TargetChain, log)
		if err != nil {
			recordFn(sh.SetActions(ProvideLiquidityActionGateWithdrawWaiting).SetStatus(provider.TxStatusFailed).Out(), err)
			return sr.SetError(err).Out(), err
		}
		recordFn(sh.SetActions(ProvideLiquidityActionGateWithdrawComplete).SetStatus(provider.TxStatusSuccess).Out())
	}

	if actionNumber <= 8 && !isActionSuccess {
		recordFn(sh.SetActions(ProvideLiquidityActionTargetChainWithdrawComplete).SetStatus(provider.TxStatusPending).Out())
		token := g.config.GetTokenInfoOnChain(args.TargetToken, args.TargetChain)
		balance, err := chains.GetTokenBalance(ctx, targetClient, token.ContractAddress, args.Sender.GetAddress().Hex(), token.Decimals)
		if err != nil {
			err = errors.Wrap(err, "get token balance")
			recordFn(sh.SetActions(ProvideLiquidityActionTargetChainWithdrawComplete).SetStatus(provider.TxStatusFailed).Out(), err)
			return sr.SetError(err).Out(), err
		}
		if balance.Cmp(decimal.Zero) == 0 {
			err = errors.Errorf("balance is zero")
			recordFn(sh.SetActions(ProvideLiquidityActionTargetChainWithdrawComplete).SetStatus(provider.TxStatusFailed).Out(), err)
			return sr.SetError(err).Out(), err
		}
		recordFn(sh.SetActions(ProvideLiquidityActionTargetChainWithdrawComplete).SetStatus(provider.TxStatusSuccess).Out())
	}
	return sr.SetCurrentChain(args.TargetChain).SetStatus(provider.TxStatusSuccess).SetTx(tx).Out(), nil

}

func (g *Gate) WaitForDeposit(ctx context.Context, tokenName, chainName string, tx common.Hash, log *logrus.Entry) error {
	var (
		t = time.NewTicker(time.Second * 5)
	)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return context.Canceled
		case <-t.C:
			ok, status, err := g.CheckDeposit(ctx, tokenName, chainName, tx)
			if errors.Is(err, error_types.ErrNotFound) {
				log.Debugf("The deposit is not found, waiting it...")
				continue
			}
			if err != nil {
				return errors.Wrap(err, "check deposit error")
			}
			log.Infof("%s on %s deposit to gate status is %s", tokenName, chainName, status)
			if ok {
				return nil
			}
		}
	}
}

func (g *Gate) WaitForWithdrawal(ctx context.Context, withdrawOrderId, tokenName, targetChain string, log *logrus.Entry) error {
	var t = time.NewTicker(time.Second * 5)
	defer t.Stop()
	for {

		select {
		case <-ctx.Done():
			return context.Canceled
		case <-t.C:
			withdrawalRecords, _, err := g.client.WalletApi.ListWithdrawals(ctx, &gateapi.ListWithdrawalsOpts{
				Currency: optional.NewString(TokenName2GateTokenName(tokenName)),
			})
			if err != nil {
				return errors.Wrap(err, "list withdrawals error")
			}
			if len(withdrawalRecords) == 0 {
				return errors.New("withdraw order not found")
			}
			var withdrawalRecord gateapi.WithdrawalRecord
			for index, v := range withdrawalRecords {
				if withdrawOrderId == "" &&
					strings.EqualFold(v.Chain, ChainName2GateChainName(targetChain)) &&
					strings.EqualFold(v.Currency, TokenName2GateTokenName(tokenName)) {
					withdrawalRecord = withdrawalRecords[index]
					break
				}
				if v.Id == withdrawOrderId {
					withdrawalRecord = withdrawalRecords[index]
					break
				}
			}
			if withdrawalRecord.Status == "" {
				return errors.Errorf("withdraw order status is empty")
			}
			log.Debugf("withdraw order %s status: %s", withdrawalRecord.Id, withdrawalRecord.Status)
			switch strings.ToUpper(withdrawalRecord.Status) {
			case "DONE":
				return nil
			}
		}
	}
}
