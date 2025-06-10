package binance

import (
	"context"
	"encoding/json"
	"fmt"
	"omni-balance/utils"
	"omni-balance/utils/chains"
	"omni-balance/utils/configs"
	"omni-balance/utils/constant"
	"omni-balance/utils/provider"
	"strings"
	"time"

	log "omni-balance/utils/logging"

	binance_connector "github.com/binance/binance-connector-go"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/spf13/cast"
)

var (
	// chain to binance action
	// send tx
	SendTxAction = "send_tx"
	// waiting tx
	WaitingTxAction = "waiting_tx"
	// deposit success
	DepositSuccessAction = "deposit_success"

	// binance to chain action
	// withdraw
	WithdrawAction = "withdraw"
	// waiting withdraw
	WaitingWithdrawAction = "waiting_withdraw"
	// withdraw success
	WithdrawSuccessAction = "withdraw_success"

	// swap action
	TobinanceSwapAction = "swap"
	// swap success
	TobinanceSwapSuccessAction = "swap_success"

	// toChain
	ToChainSwapAction = "to_chain_swap"
	// toBinance
	ToChainSwapSuccessAction = "to_chain_swap_success"

	ChainToBinanceRebalanceDirection = "chain_to_binance"
	BinanceToChainRebalanceDirection = "binance_to_chain"

	Chain2BinanceNetwork = map[string]string{
		constant.Ethereum: "ETH",
	}
)

func Action2Int(action string) int {
	switch action {
	case ToChainSwapAction:
		return 1
	case ToChainSwapSuccessAction:
		return 2
	case SendTxAction:
		return 1
	case WaitingTxAction:
		return 2
	case DepositSuccessAction:
		return 3
	case WithdrawAction:
		return 4
	case WaitingWithdrawAction:
		return 5
	case WithdrawSuccessAction:
		return 6
	case TobinanceSwapAction:
		return 7
	case TobinanceSwapSuccessAction:
		return 8
	default:
		return 0
	}
}

type Binance struct {
	authConfig binanceConfig
	config     configs.Config
	client     *binance_connector.Client
}

type binanceConfig struct {
	Key    string `json:"key" yaml:"key" help:"API key"`
	Secret string `json:"secret" yaml:"secret" help:"API secret"`
}

func New(conf configs.Config, noInit ...bool) (provider.Provider, error) {
	if len(noInit) > 0 && noInit[0] {
		return &Binance{}, nil
	}
	d := &Binance{
		config: conf,
	}

	err := conf.GetProvidersConfig(d.Name(), configs.CEX, &d.authConfig)
	if err != nil {
		return nil, errors.Wrap(err, "get gate config")
	}
	d.client = binance_connector.NewClient(d.authConfig.Key, d.authConfig.Secret)
	_, err = d.client.NewAccountInfoService().Do(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "init binance api")
	}
	return d, nil
}

func (g *Binance) Swap(ctx context.Context, args provider.SwapParams) (provider.SwapResult, error) {

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
		wallet = args.Sender
	)

	botConfigRaw := g.config.GetBotConfigUnderWallet(args.Receiver, g.Name())
	var (
		botConfig = struct {
			// tokenName: chainName
			OnlyWayTokens map[string]string `json:"onlyOneWayTokens" yaml:"onlyOneWayTokens" help:"Only way tokens"`
			// sourceToken: targetToken
			ChangeTokens map[string]string `json:"changeTokens" yaml:"changeTokens" help:"Change tokens"`
		}{}
	)
	if botConfigRaw != nil {
		data, _ := json.Marshal(botConfigRaw)
		_ = json.Unmarshal(data, &botConfig)
	}
	needSwap := len(botConfig.ChangeTokens) != 0 && botConfig.ChangeTokens[args.TargetToken] != ""

	if args.Remark == "" {
		return provider.SwapResult{}, errors.Errorf("binance swap remark must be set, like %s or %s", ChainToBinanceRebalanceDirection, BinanceToChainRebalanceDirection)
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

	if args.TargetChain == "" {
		return sr.OutError(errors.Errorf("binance swap target chain must be set"))
	}
	if args.SourceChain == "" {
		for _, v := range g.config.SourceTokens {
			if v.Name != args.SourceToken {
				continue
			}
			for _, chain := range v.Chains {
				if args.TargetChain == chain {
					continue
				}
				chainConfig := g.config.GetChainConfig(chain)
				if chainConfig.Name == "" || len(chainConfig.RpcEndpoints) == 0 {
					continue
				}
				ethClient, err := chains.NewTryClient(ctx, chainConfig.RpcEndpoints)
				if err != nil {
					return sr.OutError(errors.Wrap(err, "get balance"))
				}
				tokenName := args.SourceToken
				token := g.config.GetTokenInfoOnChain(tokenName, chainConfig.Name)
				if token.Name == "" {
					ethClient.Close()
					return sr.OutError(errors.Errorf("token %s not found in chain %s", args.SourceToken, chainConfig.Name))
				}
				balance, err := args.Sender.GetExternalBalance(ctx, common.HexToAddress(token.ContractAddress), token.Decimals, ethClient)
				if err != nil {
					ethClient.Close()
					return sr.OutError(errors.Wrap(err, "get balance"))
				}
				ethClient.Close()
				log.Infof("try to find source chain, chain: %s, balance: %s, amount: %s", chainConfig.Name, balance.String(), args.Amount.String())
				if balance.LessThan(args.Amount) {
					continue
				}
				args.SourceToken = tokenName
				args.SourceChain = chain
				break
			}
		}
		args.SaveOrderFn(map[string]interface{}{
			"source_chain_name": args.SourceChain,
		})
	}
	if args.SourceChain == "" {
		log.Infof("can't find source chain in binance config, source token: %s, target chain: %s", args.SourceToken, args.TargetChain)
		return sr.OutError(errors.Errorf("can't find source chain, source token: %s, target chain: %s", args.SourceToken, args.TargetChain))
	}
	var (
		ethClient *chains.Client
		err       error
	)
	if args.SourceChain != "binance" {
		chain := g.config.GetChainConfig(args.SourceChain)
		ethClient, err = chains.NewTryClient(ctx, chain.RpcEndpoints)
		if err != nil {
			return sr.OutError(errors.New("dial rpc"))
		}
		defer ethClient.Close()
	}

	actionNumber := Action2Int(args.LastHistory.Actions)

	if strings.EqualFold(args.TargetChain, "binance") { // send token from chain to binance
		ctx = context.WithValue(ctx, constant.ChainNameKeyInCtx, args.SourceChain)
		// get deposit address
		depositAddress, err := g.client.NewDepositAddressListService().
			Network(Chain2BinanceNetwork[args.SourceChain]).
			Coin(constant.GetBinanceTokenName(args.TargetToken)).Do(ctx)
		if err != nil {
			return sr.OutError(errors.Wrap(err, "get deposit address"))
		}
		if len(depositAddress) == 0 {
			return sr.OutError(errors.Errorf("no deposit address found"))
		}
		var deposit = depositAddress[0].Address
		if deposit == "" { // not found
			return sr.OutError(errors.Errorf("no deposit address found"))
		}
		token := g.config.GetTokenInfoOnChain(args.SourceToken, args.SourceChain)
		if actionNumber <= Action2Int(SendTxAction) && !isActionSuccess {
			amountWei := decimal.NewFromBigInt(chains.EthToWei(args.Amount, token.Decimals), 0)
			log.Infof("send tx to %s, source chain: %s, token: %s, amount: %s, deposit address: %s, sender: %s token contract: %s", args.SourceChain, args.SourceChain, token.Name, amountWei, deposit, args.Sender.GetAddress(true).Hex(), token.ContractAddress)
			// send tx
			tx, err := chains.BuildSendToken(ctx, chains.SendTokenParams{
				Client:        ethClient,
				Sender:        args.Sender.GetAddress(true),
				TokenAddress:  common.HexToAddress(token.ContractAddress),
				TokenDecimals: token.Decimals,
				ToAddress:     common.HexToAddress(deposit),
				AmountWei:     amountWei,
			})
			if err != nil {
				return sr.OutError(errors.Wrap(err, "build send tx"))
			}
			recordFn(sh.SetActions(SendTxAction).SetStatus(provider.TxStatusPending).Out())
			txHash, err := args.Sender.SendTransaction(ctx, tx, ethClient)
			if err != nil {
				return sr.OutError(errors.Wrap(err, "send tx"))
			}
			recordFn(sh.SetTx(txHash.Hex()).SetActions(WaitingTxAction).SetStatus(provider.TxStatusPending).Out())
			args.LastHistory.Tx = txHash.Hex()
		}
		if actionNumber <= Action2Int(WaitingTxAction) && !isActionSuccess {
			err = args.Sender.WaitTransaction(ctx, common.HexToHash(args.LastHistory.Tx), ethClient)
			if err != nil {
				return sr.OutError(errors.Wrap(err, "wait tx"))
			}
			recordFn(sh.SetActions(DepositSuccessAction).SetStatus(provider.TxStatusPending).Out())
		}

		if actionNumber <= Action2Int(DepositSuccessAction) && !isActionSuccess {
			// wait deposit success from binance
			status2text := map[int]string{
				0: "pending",
				6: "credited but cannot withdraw",
				7: "Wrong Deposit",
				8: "Waiting User confirm",
				1: "success",
				2: "rejected",
			}
			realHash, err := args.Sender.GetRealHash(ctx, common.HexToHash(args.LastHistory.Tx), ethClient)
			if err != nil {
				return sr.OutError(errors.Wrap(err, "get real hash"))
			}
		WAIT_DEPOSIT_SUCCESS:
			for {
				select {
				case <-ctx.Done():
					return sr.OutError(errors.New("context done"))
				case <-time.After(time.Second * 10):
					depositHistory, err := g.client.NewDepositHistoryService().TxId(realHash.Hex()).Do(ctx)
					if err != nil {
						return sr.OutError(errors.Wrap(err, "get deposit history"))
					}
					if len(depositHistory) == 0 {
						log.Infof("wait for deposit success, txid: %s ...", realHash.Hex())
						continue
					}
					for _, v := range depositHistory {
						if !strings.EqualFold(v.TxId, realHash.Hex()) {
							continue
						}
						switch v.Status { // 0(0:pending, 6:credited but cannot withdraw, 7:Wrong Deposit, 8:Waiting User confirm, 1:success, 2:rejected)
						case 1:
							recordFn(sh.SetActions(TobinanceSwapAction).SetStatus(provider.TxStatusPending).Out())
							log.Infof("deposit success, txid: %s, orderid: %s", v.TxId, v.Id)
							break WAIT_DEPOSIT_SUCCESS
						case 7, 8, 2: //
							sr.SetTx(v.TxId).SetOrderId(v.Id).SetStatus(provider.TxStatusPending)
							log.Infof("deposit not success, The deposit status is %s", status2text[v.Status])
							return sr.OutError(errors.Errorf("deposit not success, The deposit status is %s", status2text[v.Status]))
						default:
							log.Infof("binance deposit status code: %d, status text: %s, txid: %s, orderid: %s", v.Status, status2text[v.Status], v.TxId, v.Id)
							continue
						}
					}
				}
			}
		}

		botConfigRaw := g.config.GetBotConfigUnderWallet(args.Receiver, g.Name())
		var (
			botConfig = struct {
				// tokenName: chainName
				OnlyWayTokens map[string]string `json:"onlyOneWayTokens" yaml:"onlyOneWayTokens" help:"Only way tokens"`
				// sourceToken: targetToken
				ChangeTokens map[string]string `json:"changeTokens" yaml:"changeTokens" help:"Change tokens"`
			}{}
		)
		if botConfigRaw != nil {
			data, _ := json.Marshal(botConfigRaw)
			_ = json.Unmarshal(data, &botConfig)
		}
		if len(botConfig.ChangeTokens) == 0 || botConfig.ChangeTokens[args.TargetToken] == "" {
			recordFn(sh.SetActions(TobinanceSwapSuccessAction).SetStatus(provider.TxStatusSuccess).Out())
			return sr.SetStatus(provider.TxStatusSuccess).SetCurrentChain(args.TargetChain).Out(), nil
		}
		if actionNumber <= Action2Int(TobinanceSwapAction) && !isActionSuccess {

			info, err := g.client.NewExchangeInfoService().Do(context.TODO())
			if err != nil {
				return sr.OutError(errors.Wrap(err, "get exchange info"))
			}
			symbol := ""
			for _, v := range info.Symbols {
				if strings.Contains(v.Symbol, args.TargetToken) && strings.Contains(v.Symbol, botConfig.ChangeTokens[args.TargetToken]) {
					symbol = v.Symbol
					break
				}
			}
			side := ""
			if strings.EqualFold(botConfig.ChangeTokens[args.TargetToken]+args.TargetToken, symbol) {
				side = "SELL"
			} else {
				side = "BUY"
			}
			log.Infof("start binance swap change token: %s -> %s, side: %s, symbol: %s, amount: %s", botConfig.ChangeTokens[args.TargetToken], args.TargetToken, side, symbol, args.Amount.Truncate(0).String())
			resp, err := g.client.NewCreateOrderService().
				Side(side).Type("MARKET").
				Symbol(symbol).
				Quantity(args.Amount.Truncate(0).InexactFloat64()).
				Do(ctx)
			if err != nil {
				log.Errorf("binance swap change token: %s -> %s, side: %s, symbol: %s, amount: %s, error: %s", botConfig.ChangeTokens[args.TargetToken], args.TargetToken, side, symbol, args.Amount.String(), err)
				return sr.OutError(errors.Wrap(err, "swap"))
			}
			fullResp := resp.(*binance_connector.CreateOrderResponseFULL)
			log.Infof("binance swap change token: %s -> %s success, orderid: %d", botConfig.ChangeTokens[args.TargetToken], args.TargetToken, fullResp.OrderId)
			recordFn(sh.SetTx(cast.ToString(fullResp.OrderId)).SetActions(TobinanceSwapSuccessAction).SetStatus(provider.TxStatusPending).Out())
		}

		if actionNumber <= Action2Int(TobinanceSwapSuccessAction) && !isActionSuccess {
			for {
				select {
				case <-ctx.Done():
					return sr.OutError(errors.New("context done"))
				case <-time.After(time.Second * 2):
					orderStatus, err := g.client.NewGetOrderService().Symbol(botConfig.ChangeTokens[args.TargetToken] + args.TargetToken).OrderId(cast.ToInt64(args.LastHistory.Tx)).Do(ctx)
					if err != nil {
						return sr.OutError(errors.Wrap(err, "get order status"))
					}
					if orderStatus.Status == "FILLED" {
						recordFn(sh.SetActions(TobinanceSwapSuccessAction).SetStatus(provider.TxStatusSuccess).Out())
						return sr.SetStatus(provider.TxStatusSuccess).SetCurrentChain(args.TargetChain).Out(), nil
					}
					log.Infof("binance swap status: %s, txid: %s, orderid: %s", orderStatus.Status, args.LastHistory.Tx, args.LastHistory.Tx)
				}
			}
		}
		return sr.SetStatus(provider.TxStatusSuccess).SetCurrentChain(args.TargetChain).Out(), nil
	}

	if !strings.EqualFold(args.Remark, BinanceToChainRebalanceDirection) {
		log.Infof("binance swap remark must be set, like %s or %s", ChainToBinanceRebalanceDirection, BinanceToChainRebalanceDirection)
		return sr.OutError(errors.Errorf("binance swap remark must be set, like %s or %s", ChainToBinanceRebalanceDirection, BinanceToChainRebalanceDirection))
	}

	if actionNumber <= Action2Int(ToChainSwapAction) && !isActionSuccess && needSwap {

		info, err := g.client.NewExchangeInfoService().Do(context.TODO())
		if err != nil {
			return sr.OutError(errors.Wrap(err, "get exchange info"))
		}
		symbol := ""
		for _, v := range info.Symbols {
			if strings.Contains(v.Symbol, args.TargetToken) && strings.Contains(v.Symbol, botConfig.ChangeTokens[args.TargetToken]) {
				symbol = v.Symbol
				break
			}
		}
		side := ""
		if strings.EqualFold(botConfig.ChangeTokens[args.TargetToken]+args.TargetToken, symbol) {
			side = "SELL"
		} else {
			side = "BUY"
		}
		log.Infof("start binance swap change token: %s -> %s, side: %s, symbol: %s, amount: %s", botConfig.ChangeTokens[args.TargetToken], args.TargetToken, side, symbol, args.Amount.Truncate(0).String())
		resp, err := g.client.NewCreateOrderService().
			Side(side).Type("MARKET").
			Symbol(symbol).
			Quantity(args.Amount.InexactFloat64()).
			Do(ctx)
		if err != nil {
			return sr.OutError(errors.Wrap(err, "swap"))
		}
		fullResp := resp.(*binance_connector.CreateOrderResponseFULL)
		log.Infof("binance swap change token success: %s -> %s, orderid: %d", args.TargetToken, botConfig.ChangeTokens[args.TargetToken], fullResp.OrderId)
		args.LastHistory.Tx = cast.ToString(fullResp.OrderId)
		recordFn(sh.SetTx(cast.ToString(fullResp.OrderId)).SetActions(TobinanceSwapSuccessAction).SetStatus(provider.TxStatusPending).Out())
	}

	if actionNumber <= Action2Int(TobinanceSwapSuccessAction) && !isActionSuccess && needSwap {
	WAIT_ORDER_SUCCESS:
		for {
			select {
			case <-ctx.Done():
				return sr.OutError(errors.New("context done"))
			case <-time.After(time.Second * 2):
				orderStatus, err := g.client.NewGetOrderService().Symbol(args.TargetToken + botConfig.ChangeTokens[args.TargetToken]).OrderId(cast.ToInt64(args.LastHistory.Tx)).Do(ctx)
				if err != nil {
					return sr.OutError(errors.Wrap(err, "get order status"))
				}
				if orderStatus.Status != "FILLED" {
					log.Infof("binance swap status: %s, txid: %s, orderid: %d", orderStatus.Status, args.LastHistory.Tx, orderStatus.OrderId)
					continue
				}
				recordFn(sh.SetActions(WithdrawAction).SetStatus(provider.TxStatusPending).Out())
				log.Infof("binance swap success, txid: %s, orderid: %d status: %s", args.LastHistory.Tx, orderStatus.OrderId, orderStatus.Status)
				break WAIT_ORDER_SUCCESS
			}
		}
	}
	if actionNumber <= Action2Int(WithdrawAction) && !isActionSuccess {
		// withdraw from binance
		withdrawOrderId := fmt.Sprintf("omniBalanceTo%s%s%s%s", args.SourceChain, args.TargetChain, args.TargetToken, uuid.New().String())
		log.Infof("binance withdraw, orderid: %s, amount: %s, token: %s, chain: %s",
			withdrawOrderId, args.Amount.String(), args.TargetToken, args.TargetChain)
		_, err = g.client.NewWithdrawService().
			Address(args.Receiver).
			Amount(args.Amount.Truncate(0).InexactFloat64()).
			Coin(args.TargetToken).
			Network(args.TargetChain).
			WithdrawOrderId(withdrawOrderId).
			Do(ctx)
		if err != nil {
			return sr.OutError(errors.Wrap(err, "withdraw"))
		}
		recordFn(sh.SetTx(withdrawOrderId).SetActions(WaitingWithdrawAction).SetStatus(provider.TxStatusPending).Out())
		args.LastHistory.Tx = withdrawOrderId
	}
	status2text := map[int]string{
		0: "send confirm email",
		2: "waiting confirm",
		3: "rejected",
		4: "processing",
		6: "withdraw success",
	}
	if actionNumber <= Action2Int(WaitingWithdrawAction) && !isActionSuccess {
		withdrawOrderId := cast.ToString(args.LastHistory.Tx)
		// wait withdraw success
	LOOP:
		for {
			select {
			case <-ctx.Done():
				return sr.OutError(errors.New("context done"))
			case <-time.After(time.Second * 10):
				withdrawHistory, err := g.client.NewWithdrawHistoryService().WithdrawOrderId(withdrawOrderId).Do(ctx)
				if err != nil {
					return sr.OutError(errors.Wrap(err, "get withdraw history"))
				}
				if len(withdrawHistory) == 0 {
					log.Infof("wait for withdraw success, txid: %s ...", withdrawOrderId)
					continue
				}
				for _, v := range withdrawHistory {
					if !strings.EqualFold(v.WithdrawOrderId, withdrawOrderId) {
						continue
					}
					log.Infof("binance withdraw status: %s, txid: %s, orderid: %s", status2text[v.Status], v.WithdrawOrderId, v.Id)
					switch v.Status {
					case 6:
						recordFn(sh.SetActions(WithdrawSuccessAction).SetStatus(provider.TxStatusSuccess).Out())
						sr.SetTx(v.WithdrawOrderId).SetOrderId(v.Id).SetStatus(provider.TxStatusSuccess).SetCurrentChain(args.TargetChain)
						break LOOP
					default:
						log.Infof("binance withdraw status: %s, txid: %s, orderid: %s", status2text[v.Status], v.WithdrawOrderId, v.Id)
						continue
					}
				}
			}
		}
	}
	return sr.SetCurrentChain(args.TargetChain).Out(), nil
}

func (g *Binance) GetCost(ctx context.Context, args provider.SwapParams) (provider.TokenInCosts, error) {
	if strings.Contains(args.Remark, "binance") {
		return provider.TokenInCosts{
			{
				TokenName:  args.TargetToken,
				CostAmount: decimal.Zero,
			},
		}, nil
	}
	return provider.TokenInCosts{}, nil
}

func (g *Binance) CheckToken(ctx context.Context, tokenName, tokenInChainName, tokenOutChainName string, amount decimal.Decimal) (bool, error) {
	return true, nil
}

func (g *Binance) Help() []string {
	var result []string
	for _, v := range utils.ExtractTagFromStruct(&binanceConfig{}, "yaml", "help") {
		result = append(result, v["yaml"]+": "+v["help"])
	}
	result = append(result, "api key must be granted with the following permissions: Withdraw")
	return result
}

func (g *Binance) Name() string {
	return "binance"
}

func (g *Binance) Type() configs.ProviderType {
	return configs.CEX
}
