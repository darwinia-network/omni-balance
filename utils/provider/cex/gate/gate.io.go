package gate

import (
	"context"
	"fmt"
	"github.com/antihax/optional"
	"github.com/gateio/gateapi-go/v6"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"omni-balance/utils"
	"omni-balance/utils/configs"
	"omni-balance/utils/error_types"
	"omni-balance/utils/provider"
	"strings"
	"time"
)

var (
	PurchasingAction  = "purchasing"
	PurchasedAction   = "purchased"
	WithdrawingAction = "withdrawing"
	WithdrawnAction   = "withdrawn"
)

func Action2Int(action string) int {
	switch action {
	case PurchasingAction:
		return 1
	case PurchasedAction:
		return 2
	case WithdrawingAction:
		return 3
	case WithdrawnAction:
		return 4
	default:
		return 0
	}
}

type Gate struct {
	authConfig gateConfig
	config     configs.Config
	client     *gateapi.APIClient
}

type gateConfig struct {
	Key    string `json:"key" yaml:"key" help:"API key"`
	Secret string `json:"secret" yaml:"secret" help:"API secret"`
}

func New(conf configs.Config, noInit ...bool) (provider.Provider, error) {
	if len(noInit) > 0 && noInit[0] {
		return &Gate{}, nil
	}
	d := &Gate{
		config: conf,
	}

	err := conf.GetProvidersConfig(d.Name(), configs.CEX, &d.authConfig)
	if err != nil {
		return nil, errors.Wrap(err, "get gate config")
	}
	c := gateapi.NewConfiguration()
	c.Key = d.authConfig.Key
	c.Secret = d.authConfig.Secret
	d.client = gateapi.NewAPIClient(c)
	_, _, err = d.client.AccountApi.GetAccountDetail(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "init gate api")
	}
	return d, nil
}

func (g *Gate) Swap(ctx context.Context, args provider.SwapParams) (provider.SwapResult, error) {
	if args.Remark != "" {
		return g.ProvideLiquidity(ctx, args)
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
		log = logrus.WithFields(logrus.Fields{
			"provider":    g.Name(),
			"amount":      args.Amount,
			"tokenOut":    args.TargetToken,
			"targetChain": args.TargetChain,
		})
		wallet = args.Sender
	)

	isSenderInVerifiedAddress, err := g.IsVerifiedAddress(ctx, wallet.GetAddress(true).Hex(), args.TargetToken, args.TargetChain)
	if err != nil {
		return provider.SwapResult{}, errors.Wrap(err, "check sender address")
	}
	if !isSenderInVerifiedAddress {
		return provider.SwapResult{}, errors.Errorf("%s not in gate.io verified address for %s",
			args.Sender, args.TargetToken)
	}

	actionNumber := Action2Int(args.LastHistory.Actions)

	tokenInCosts, err := g.GetCost(ctx, args)
	if err != nil {
		return provider.SwapResult{}, err
	}
	log.Debugf("tokenInCosts: %+v", tokenInCosts)
	tokenIn, err := tokenInCosts.GetBest()
	if err != nil {
		return provider.SwapResult{}, errors.Wrapf(err, "%s get best token in cost", g.Name())
	}
	log.Debugf("get best token in cost: %s %s", tokenIn.TokenName, tokenIn.CostAmount)

	accounts, _, err := g.client.SpotApi.ListSpotAccounts(context.Background(), &gateapi.ListSpotAccountsOpts{})
	if err != nil {
		return provider.SwapResult{}, errors.Wrap(err, "list spot accounts")
	}
	if len(accounts) == 0 {
		return provider.SwapResult{}, errors.Errorf("no %s in spot account", tokenIn.TokenName)
	}
	var (
		buyAmount       = args.Amount
		tokenInAccount  gateapi.SpotAccount
		tokenOutAccount gateapi.SpotAccount
	)
	if strings.EqualFold(tokenIn.TokenName, args.TargetToken) {
		return provider.SwapResult{}, errors.Errorf("token in and token out are the same")
	}
	for _, account := range accounts {
		log.Debugf("account: %+v, tokenInName: %s, tokenOutName: %s", account, tokenIn.TokenName, args.TargetToken)
		if strings.EqualFold(account.Currency, tokenIn.TokenName) {
			tokenInAccount = account
		}
		if strings.EqualFold(account.Currency, args.TargetToken) {
			tokenOutAccount = account
			continue
		}
	}
	log.Debugf("token in account: %+v, token out account: %+v", tokenInAccount, tokenOutAccount)

	if actionNumber <= 1 && decimal.RequireFromString(tokenInAccount.Available).LessThan(tokenIn.CostAmount) {
		return provider.SwapResult{}, errors.Errorf("not enough %s in spot account, available: %s, need: %s",
			tokenIn.TokenName, tokenInAccount.Available, tokenIn.CostAmount)
	}

	var (
		sr = new(provider.SwapResult).
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
	)

	if buyAmount.GreaterThan(decimal.Zero) && actionNumber <= 1 && !isActionSuccess { // buy token
		recordFn(sh.SetActions(PurchasedAction).SetStatus(provider.TxStatusPending).Out())
		order, err := g.buyToken(ctx, tokenIn, args.TargetToken, buyAmount.Abs(), func(order gateapi.Order) bool {
			logrus.Debugf("wait for order %s filled, the status is %s", order.Id, order.Status)
			return false
		})
		if err != nil {
			recordFn(sh.SetActions(PurchasedAction).SetStatus(provider.TxStatusFailed).Out(), err)
			return sr.SetStatus(provider.TxStatusFailed).SetError(err).Out(), err
		}
		sr = sr.SetOrderId(order.Id).SetOrder(order)
		sh = sh.SetTx(order.Id)
		recordFn(sh.SetActions(PurchasedAction).SetStatus(provider.TxStatusSuccess).Out())
	}

	var withdrawOrderId = args.LastHistory.Tx

	if actionNumber < 4 && !isActionSuccess {
		balance, _ := decimal.NewFromString(tokenOutAccount.Available)
		if balance.LessThanOrEqual(decimal.Zero) {
			err = errors.Errorf("not enough %s in spot account", args.TargetToken)
			return sr.SetStatus(provider.TxStatusFailed).SetError(err).Out(), err
		}
		if balance.LessThanOrEqual(args.Amount) {
			args.Amount = balance
		}
		recordFn(sh.SetActions(WithdrawingAction).SetStatus(provider.TxStatusPending).Out())
		withdrawOrder, _, err := g.client.WithdrawalApi.Withdraw(ctx, gateapi.LedgerRecord{
			Amount:   args.Amount.String(),
			Currency: args.TargetToken,
			Address:  wallet.GetAddress(true).Hex(),
			Chain:    ChainName2GateChainName(args.TargetChain),
		})
		if err != nil {
			recordFn(sh.SetActions(WithdrawingAction).SetStatus(provider.TxStatusFailed).Out(), err)
			return sr.SetStatus(provider.TxStatusFailed).SetError(err).Out(), errors.Wrap(err, "withdraw")
		}
		withdrawOrderId = withdrawOrder.Id

		sr = sr.SetOrder(withdrawOrder).SetOrderId(withdrawOrderId).SetTx(withdrawOrderId)
		sh = sh.SetTx(withdrawOrderId)
		recordFn(sh.SetActions(WithdrawnAction).SetStatus(provider.TxStatusPending).Out())
	}

	for {
		logrus.Debugf("wait for %s withdraw order done", withdrawOrderId)
		withdrawalRecords, _, err := g.client.WalletApi.ListWithdrawals(ctx, &gateapi.ListWithdrawalsOpts{
			Currency: optional.NewString(args.TargetToken),
		})
		if err != nil {
			return sr.SetStatus(provider.TxStatusFailed).SetError(err).Out(), errors.Wrap(err, "list withdraw status")
		}
		if len(withdrawalRecords) == 0 {
			err = errors.Errorf("withdraw order not found")
			return sr.SetStatus(provider.TxStatusFailed).SetError(err).Out(), err
		}
		var withdrawalRecord gateapi.WithdrawalRecord
		for index, v := range withdrawalRecords {
			if withdrawOrderId == "" &&
				strings.EqualFold(v.Chain, ChainName2GateChainName(args.TargetChain)) &&
				strings.EqualFold(v.Currency, args.TargetToken) {
				withdrawalRecord = withdrawalRecords[index]
				break
			}
			if v.Id == withdrawOrderId {
				withdrawalRecord = withdrawalRecords[index]
				break
			}
		}
		if withdrawalRecord.Status == "" {
			err = errors.Errorf("withdraw order not found")
			return sr.SetStatus(provider.TxStatusFailed).SetError(err).Out(), err
		}
		log.Infof("withdraw status: %s, wait for done", withdrawalRecord.Status)
		switch strings.ToUpper(withdrawalRecord.Status) {
		case "DONE":
			recordFn(sh.SetActions(WithdrawingAction).SetStatus(provider.TxStatusSuccess).SetCurrentChain(withdrawalRecord.Chain).Out())
			sr.SetCurrentChain(withdrawalRecord.Chain).
				SetStatus(provider.TxStatusSuccess).
				SetTx(withdrawalRecord.Txid).
				SetOrderId(withdrawalRecord.Id).
				SetOrder(withdrawalRecord)
			log.Infof("withdraw done, txid: %s", withdrawalRecord.Txid)
			return sr.Out(), nil
		case "CANCEL", "FAIL":
			err := errors.Errorf("withdraw failed, status: %s", withdrawalRecord.Status)
			recordFn(sh.SetStatus(provider.TxStatusFailed).SetActions(WithdrawingAction).Out(), err)
			return sr.SetStatus(provider.TxStatusFailed).SetError(err).Out(), err
		case "REQUEST", "BCODE", "EXTPEND", "VERIFY", "PROCES", "PEND", "LOCKED":
			time.Sleep(time.Second)
			continue
		default:
			err = errors.Errorf("unknown status: %+v", withdrawalRecord.Status)
			recordFn(sh.SetStatus(provider.TxStatusFailed).SetActions(WithdrawingAction).Out(), err)
			return sr.SetStatus(provider.TxStatusFailed).SetError(err).Out(), err
		}
	}
}

func (g *Gate) GetCost(ctx context.Context, args provider.SwapParams) (provider.TokenInCosts, error) {
	log := logrus.WithFields(logrus.Fields{
		"tokenOut":    args.TargetToken,
		"targetChain": args.TargetChain,
		"amount":      args.Amount,
		"wallet":      args.Sender,
	})
	verifiedAddress, _, err := g.client.WalletApi.ListSavedAddress(ctx, args.TargetToken, &gateapi.ListSavedAddressOpts{})
	if err != nil {
		return nil, errors.Wrap(err, "list saved address")
	}
	if len(verifiedAddress) == 0 {
		return nil, errors.Errorf("%s not in gate.io verified address for %s", args.Sender, args.TargetToken)
	}
	chains, err := g.GetTokenOutChain(ctx, args.TargetToken)
	if err != nil {
		return nil, errors.Wrap(err, "get token out chains")
	}
	log.Debugf("token %s out chains: %+v", args.TargetToken, chains)
	if !utils.InArray(args.TargetChain, chains) {
		providers, err := provider.GetTokenCrossChainProviders(ctx, provider.GetTokenCrossChainProvidersParams{
			SourceChains: chains,
			TargetChain:  args.TargetChain,
			TokenName:    args.TargetToken,
			Amount:       args.Amount,
			Conf:         g.config,
		})
		if err != nil {
			return nil, errors.Wrap(err, "get token cross chain providers")
		}
		if len(providers) == 0 {
			return nil, error_types.ErrUnsupportedTokenAndChain
		}
	}

	currency, _, err := g.client.SpotApi.GetCurrency(ctx, TokenName2GateTokenName(args.TargetToken))
	if err != nil {
		return nil, errors.Wrap(err, "get currency")
	}
	log.Debugf("currency: %+v", currency)
	if currency.Delisted ||
		currency.WithdrawDisabled {
		return nil, error_types.ErrUnsupportedTokenAndChain
	}

	accounts, _, err := g.client.SpotApi.ListSpotAccounts(ctx, &gateapi.ListSpotAccountsOpts{})
	if err != nil {
		return nil, errors.Wrap(err, "list spot accounts")
	}
	var (
		tokenIns []gateapi.SpotAccount
	)

	for index, v := range accounts {
		log.Debugf("spot %s token balance: %s", v.Currency, v.Available)
		if strings.EqualFold(v.Currency, TokenName2GateTokenName(args.TargetToken)) {
			continue
		}
		for _, tokenIn := range g.config.SourceTokens {
			if !strings.EqualFold(v.Currency, tokenIn.Name) {
				continue
			}
			tokenIns = append(tokenIns, accounts[index])
			continue
		}
	}

	currencyPairs, _, err := g.client.SpotApi.ListCurrencyPairs(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "list currency pairs")
	}
	var (
		tokenInNames []string
		pairs        []string
	)
	for _, v := range currencyPairs {
		if !strings.EqualFold(v.Base, args.TargetToken) {
			continue
		}
		for _, tokenIn := range tokenIns {
			if !strings.EqualFold(tokenIn.Currency, v.Quote) {
				continue
			}
			log.Debugf("currency pair: %s_%s", v.Base, v.Quote)
			tokenInNames = append(tokenInNames, tokenIn.Currency)
			pairs = append(pairs, strings.ToLower(fmt.Sprintf("%s_%s", v.Base, v.Quote)))
		}
	}
	if len(tokenInNames) == 0 {
		return nil, errors.Errorf("no token in")
	}

	pairsList, err := g.Tickers(pairs...)
	if err != nil {
		return nil, errors.Wrap(err, "get market list")
	}
	if len(pairsList) == 0 {
		return nil, errors.Errorf("no market list")
	}
	var providerTokenIns []provider.TokenInCost
	for _, v := range pairsList {
		if strings.EqualFold(v.TokenName, args.TargetToken) {
			providerTokenIns = append(providerTokenIns, provider.TokenInCost{
				TokenName:  strings.ToUpper(v.TokenName),
				CostAmount: args.Amount,
			})
			continue
		}
		providerTokenIns = append(providerTokenIns, provider.TokenInCost{
			TokenName:  strings.ToUpper(v.TokenName),
			CostAmount: v.Price.Mul(args.Amount),
		})
	}
	return providerTokenIns, nil
}

func (g *Gate) CheckToken(ctx context.Context, tokenName, tokenInChainName, tokenOutChainName string, amount decimal.Decimal) (bool, error) {
	// 检查token提现
	token, _, err := g.client.SpotApi.GetCurrency(ctx, TokenName2GateTokenName(tokenName))
	if err != nil {
		return false, errors.Wrap(err, "get currency")
	}
	return !token.WithdrawDisabled, err
}

func (g *Gate) Help() []string {
	var result []string
	for _, v := range utils.ExtractTagFromStruct(&gateConfig{}, "yaml", "help") {
		result = append(result, v["yaml"]+": "+v["help"])
	}
	result = append(result, "api key must be granted with the following permissions: Spot Trade, Wallet, Withdraw, Account")
	result = append(result, "API Key Type must be selected as API v4 Key, and since gate.io has a limitation, it cannot be used with sub-accounts")
	return result
}

func (g *Gate) Name() string {
	return "gate.io"
}

func (g *Gate) Type() configs.ProviderType {
	return configs.CEX
}
