package darwinia

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient/simulated"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"omni-balance/utils"
	"omni-balance/utils/chains"
	"omni-balance/utils/configs"
	"omni-balance/utils/constant"
	"omni-balance/utils/error_types"
	"omni-balance/utils/provider"
	"omni-balance/utils/wallets"
	"time"
)

var (
	// action
	sourceChainSendingAction  = "sourceChainSending"
	sourceChainWaitAction     = "sourceChainWait"
	targetChainSendingAction  = "targetChainSending"
	targetChainReceivedAction = "targetChainReceived"
)

func Action2Int(action string) int {
	switch action {
	case sourceChainSendingAction:
		return 1
	case sourceChainWaitAction:
		return 2
	case targetChainSendingAction:
		return 3
	case targetChainReceivedAction:
		return 4
	default:
		return 0
	}
}

type Bridge struct {
	config configs.Config
}

var graphOlQuery = "query GetHistory($sender: String, $page: Int, $row: Int) {\n  historyRecords(sender: $sender, page: $page, row: $row) {\n    total\n    records {\n      requestTxHash\n      responseTxHash\n      fromChain\n      toChain\n      startTime\n      sendToken\n      sendAmount\n      result\n      id\n      __typename\n    }\n    __typename\n  }\n}"

type SwapParams struct {
	Sender    wallets.Wallets
	TokenName string
	Amount    decimal.Decimal
	Nonce     int64
	Client    simulated.Client
}

type GraphQLRequest struct {
	OperationName string                 `json:"operationName"`
	Variables     map[string]interface{} `json:"variables"`
	Query         string                 `json:"query"`
}

type HistoryRecordResponse struct {
	Data struct {
		HistoryRecords struct {
			Total    int             `json:"total"`
			Records  []HistoryRecord `json:"records"`
			Typename string          `json:"__typename"`
		} `json:"historyRecords"`
	} `json:"data"`
}

type HistoryRecord struct {
	RequestTxHash  string `json:"requestTxHash"`
	ResponseTxHash string `json:"responseTxHash"`
	FromChain      string `json:"fromChain"`
	ToChain        string `json:"toChain"`
	StartTime      int    `json:"startTime"`
	SendToken      string `json:"sendToken"`
	SendAmount     string `json:"sendAmount"`
	Result         int    `json:"result"`
	Id             string `json:"id"`
	Typename       string `json:"__typename"`
}

func New(conf configs.Config, noInit ...bool) (provider.Provider, error) {
	if len(noInit) > 0 && noInit[0] {
		return &Bridge{}, nil
	}
	d := &Bridge{
		config: conf,
	}
	return d, nil
}

func (b *Bridge) CheckToken(_ context.Context, tokenName, tokenInChainName, tokenOutChainName string,
	_ decimal.Decimal) (bool, error) {

	sourceChains := GetSourceChains(int64(constant.GetChainId(tokenOutChainName)), tokenName)
	if len(sourceChains) == 0 {
		return false, nil
	}
	return utils.InArray(int64(constant.GetChainId(tokenInChainName)), sourceChains), nil
}

func (b *Bridge) GetCost(ctx context.Context, args provider.SwapParams) (provider.TokenInCosts, error) {
	if args.TargetToken == "" || args.TargetChain == "" {
		return nil, errors.Errorf("target token or target chain is empty")
	}
	if args.SourceChain == "" {
		sourceChains := b.FindValidSourceChains(ctx, constant.GetChainId(args.TargetChain), args.TargetToken,
			args.Sender.GetAddress(true).Hex(), args.Amount)
		if len(sourceChains) == 0 {
			return nil, error_types.ErrUnsupportedTokenAndChain
		}
		args.SourceChain = constant.GetChainName(utils.Choose(sourceChains))
	}

	if args.SourceToken == "" {
		args.SourceToken = b.config.GetTokenInfoOnChain(args.TargetToken, args.SourceChain).Name
	}

	return provider.TokenInCosts{
		provider.TokenInCost{
			TokenName:  args.SourceToken,
			CostAmount: args.Amount,
		},
	}, nil
}

func (b *Bridge) Swap(ctx context.Context, args provider.SwapParams) (result provider.SwapResult, err error) {
	var (
		history  = args.LastHistory
		log      = args.GetLogs("darwinia")
		recordFn = func(s provider.SwapHistory, errs ...error) {
			s.ProviderType = string(b.Type())
			s.ProviderName = b.Name()
			s.Amount = args.Amount
			if args.RecordFn == nil {
				return
			}
			args.RecordFn(s, errs...)
		}
	)
	if history.Actions == targetChainReceivedAction && history.Status == string(provider.TxStatusSuccess) {
		log.Debugf("target chain received, order id: %s", history.Tx)
		return provider.SwapResult{
			ProviderType: b.Type(),
			ProviderName: b.Name(),
			OrderId:      history.Tx,
			Status:       provider.TxStatusSuccess,
			CurrentChain: args.TargetChain,
			Tx:           history.Tx,
		}, nil
	}
	if args.SourceChain == args.TargetChain && history.Status == string(provider.TxStatusSuccess) {
		log.Debugf("source chain %s and target chain %s is same", args.SourceChain, args.TargetChain)
		return provider.SwapResult{}, errors.Errorf("source chain %s and target chain %s is same", args.SourceChain, args.TargetChain)
	}
	if history.CurrentChain != "" && history.CurrentChain != args.TargetChain {
		args.SourceChain = history.CurrentChain
	}
	actionNumber := Action2Int(history.Actions)
	ctx = context.WithValue(ctx, constant.ChainNameKeyInCtx, args.SourceChain)
	if actionNumber == 0 && args.SourceChain == "" {
		validSourceChain := b.GetValidSourceChain(ctx, constant.GetChainId(args.TargetChain),
			args.TargetToken, args.Sender.GetAddress(true).Hex(), args.Amount)
		if validSourceChain == 0 {
			return result, errors.Errorf("can not find source chain for %s %s", args.TargetToken, args.TargetChain)
		}
		args.SourceChain = constant.GetChainName(validSourceChain)
	}
	var sourceChainConf = b.config.GetChainConfig(args.SourceChain)
	if args.SourceToken == "" {
		args.SourceToken = b.config.GetTokenInfoOnChain(args.TargetToken, args.SourceChain).Name
	}
	sourceToken := b.config.GetTokenInfoOnChain(args.SourceToken, args.SourceChain)
	targetChainConf := b.config.GetChainConfig(args.TargetChain)

	sr := new(provider.SwapResult).
		SetTokenInName(args.SourceToken).
		SetTokenInChainName(args.SourceChain).
		SetProviderName(b.Name()).
		SetProviderType(b.Type()).
		SetCurrentChain(args.SourceChain).
		SetTx(args.LastHistory.Tx)
	sh := &provider.SwapHistory{
		ProviderName: b.Name(),
		ProviderType: string(b.Type()),
		Amount:       args.Amount,
		CurrentChain: args.SourceChain,
		Tx:           history.Tx,
	}
	isActionSuccess := history.Status == string(provider.TxStatusSuccess)

	fn := BRIDGE_DIRECTION[int64(sourceChainConf.Id)][int64(targetChainConf.Id)]
	if fn == nil {
		err = errors.Errorf("not support swap %s to %s", args.SourceChain, args.TargetChain)
		return sr.SetStatus(provider.TxStatusFailed).SetError(err).Out(), err
	}
	wallet := args.Sender
	ethClient, err := chains.NewTryClient(ctx, sourceChainConf.RpcEndpoints)
	if err != nil {
		return sr.SetStatus(provider.TxStatusFailed).SetError(err).Out(), errors.Wrap(err, "dial rpc")
	}
	defer ethClient.Close()

	log.Debugf("start transfer %s from %s to %s", args.SourceToken, args.SourceChain, args.TargetChain)
	tx, err := fn(ctx, SwapParams{
		Sender:    wallet,
		TokenName: args.SourceToken,
		Amount:    decimal.NewFromBigInt(chains.EthToWei(args.Amount, sourceToken.Decimals), 0),
		Nonce:     time.Now().UnixMilli(),
		Client:    ethClient,
	})
	if err != nil {
		return sr.SetStatus(provider.TxStatusFailed).SetError(err).Out(), errors.Wrap(err, "swap")
	}

	if tx.Gas != 0 {
		tx.Gas = tx.Gas * 2
	}

	if actionNumber <= 1 && !isActionSuccess {
		recordFn(sh.SetActions(sourceChainSendingAction).SetStatus(provider.TxStatusPending).Out())
		sr = sr.SetReciever(args.Receiver)
		ctx = provider.WithNotify(ctx, provider.WithNotifyParams{
			OrderId:         args.OrderId,
			Receiver:        common.HexToAddress(args.Receiver),
			TokenIn:         args.SourceToken,
			TokenOut:        args.TargetToken,
			TokenInChain:    args.SourceChain,
			TokenOutChain:   args.TargetChain,
			ProviderName:    b.Name(),
			TokenInAmount:   args.Amount,
			TokenOutAmount:  args.Amount,
			TransactionType: provider.TransferTransactionAction,
		})

		txHash, err := args.Sender.SendTransaction(ctx, tx, ethClient)
		if err != nil {
			recordFn(sh.SetActions(sourceChainSendingAction).SetStatus(provider.TxStatusFailed).Out(), err)
			return sr.SetStatus(provider.TxStatusFailed).SetError(err).Out(), errors.Wrap(err, "send signed transaction")
		}
		recordFn(sh.SetActions(sourceChainSendingAction).SetStatus(provider.TxStatusSuccess).Out())
		sh = sh.SetTx(txHash.Hex())
		sr = sr.SetTx(txHash.Hex()).SetOrderId(txHash.Hex())
	}

	if actionNumber <= 2 && !isActionSuccess {
		recordFn(sh.SetActions(sourceChainWaitAction).SetStatus(provider.TxStatusPending).Out())
		err = wallet.WaitTransaction(ctx, common.HexToHash(sr.Tx), ethClient)
		if err != nil {
			recordFn(sh.SetActions(sourceChainWaitAction).SetStatus(provider.TxStatusPending).Out(), err)
			return sr.SetStatus(provider.TxStatusFailed).SetError(err).Out(), errors.Wrap(err, "wait for tx")
		}
		recordFn(sh.SetActions(sourceChainWaitAction).SetStatus(provider.TxStatusSuccess).Out())
	}

	if actionNumber <= 4 && !isActionSuccess {
		recordFn(sh.SetActions(targetChainSendingAction).SetStatus(provider.TxStatusPending).Out())
		log.Debugf("waiting for bridge success")
		tx, err := wallet.GetRealHash(ctx, common.HexToHash(sr.Tx), ethClient)
		if err != nil {
			recordFn(sh.SetActions(targetChainSendingAction).SetStatus(provider.TxStatusFailed).Out(), err)
			return sr.SetStatus(provider.TxStatusFailed).SetError(err).Out(), errors.Wrap(err, "get real tx hash")
		}

		record, err := b.WaitForBridgeSuccess(ctx, tx.Hex(), args.Sender.GetAddress(true).Hex())
		if err != nil {
			recordFn(sh.SetActions(targetChainSendingAction).SetStatus(provider.TxStatusFailed).Out(), err)
			return sr.SetStatus(provider.TxStatusFailed).SetError(err).Out(), errors.Wrap(err, "wait for bridge success")
		}
		sr.SetOrder(record)

		if record.Result == 3 {
			recordFn(sh.SetActions(targetChainReceivedAction).SetStatus(provider.TxStatusSuccess).SetCurrentChain(args.TargetToken).Out())
			sr = sr.SetStatus(provider.TxStatusSuccess).SetCurrentChain(args.TargetChain)
		}
		if record.Result != 3 {
			log.Debugf("bridge failed: %+v", record)
			err := errors.Errorf("bridge failed: %d", record.Result)
			recordFn(sh.SetActions(targetChainReceivedAction).SetStatus(provider.TxStatusSuccess).Out(), err)
			return sr.SetStatus(provider.TxStatusFailed).SetError(err).Out(), err
		}
	}
	return sr.SetStatus(provider.TxStatusSuccess).SetCurrentChain(args.TargetChain).Out(), nil
}

func (b *Bridge) Check(ctx context.Context, args provider.CheckParams) (bool, error) {
	if args.Chain == "" {
		return false, errors.Errorf("target chain is empty")
	}
	if err := b.config.GetWallet(args.Wallet).CheckFullAccess(ctx); err != nil {
		return false, errors.Wrap(err, "check full access")
	}

	targetChain := b.config.GetChainConfig(args.Chain)
	sourceChains := b.FindValidSourceChains(ctx, targetChain.Id, args.Token, args.Wallet, args.Amount)
	if len(sourceChains) == 0 {
		return false, nil
	}
	return true, nil
}

func (b *Bridge) Help() []string {
	return []string{"See https://bridge.darwinia.network"}
}

func (b *Bridge) Name() string {
	return "darwinia-bridge"
}

func (b *Bridge) Type() configs.ProviderType {
	return configs.Bridge
}
