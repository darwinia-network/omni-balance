package provider

import (
	"context"
	"encoding/json"
	"omni-balance/utils/configs"
	"omni-balance/utils/wallets"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"gorm.io/datatypes"
)

type TxStatus string
type TokenInCosts []TokenInCost

const (
	TxStatusSuccess TxStatus = "success"
	TxStatusFailed  TxStatus = "failed"
	TxStatusPending TxStatus = "pending"
)

func (t TxStatus) CanRetry() bool {
	return t == TxStatusFailed || t == ""
}
func (t TxStatus) String() string {
	return string(t)
}

func (t TokenInCosts) GetBest() (TokenInCost, error) {
	if len(t) == 0 {
		return TokenInCost{}, errors.New("no token in costs")
	}
	var bestPrice = t[0]
	for _, v := range t {
		if v.CostAmount.LessThan(bestPrice.CostAmount) {
			continue
		}
		bestPrice = v
	}
	return bestPrice, nil
}

type TokenInCost struct {
	TokenName string
	// CostAmount: the amount of token need to paid
	CostAmount decimal.Decimal
}

type Price interface {
	GetTokenPriceInUSDT(ctx context.Context, tokenName, chainName, tokenAddress string) (decimal.Decimal, error)
}

type Provider interface {

	// GetCost calculates the cost of swapping tokens
	// Params:
	//   ctx: context
	//   args: SwapParams containing the necessary variables for the swap
	//
	// Returns:
	//   decimal.Decimal: the cost of the swap
	//   error: an error if the swap failed to calculate the cost
	GetCost(ctx context.Context, args SwapParams) (TokenInCosts, error)

	CheckToken(ctx context.Context, tokenName, tokenInChainName, tokenOutChainName string, amount decimal.Decimal) (bool, error)

	// Swap token, param:
	// sourceToken: use this token to exchange
	// targetToken: exchange to this token
	// targetReceiver: receive target token wallet address
	// targetChain: target token chain
	// amount: exchange amount
	Swap(ctx context.Context, args SwapParams) (result SwapResult, err error)

	// Help get provider help
	Help() []string

	// Name get provider name
	Name() string

	Type() configs.ProviderType
}

type CheckParams struct {
	Token        string
	SourceChains []string
	Chain        string
	Amount       decimal.Decimal
	Wallet       string
}

type TokenPriceParams struct {
	SourceToken string
	TargetToken string
	TargetChain string
	SourceChain string
	Sender      common.Address
}

type BalanceParams struct {
	Token  string
	Chain  string
	Wallet string
}

type SwapParams struct {
	Order            datatypes.JSON                      `json:"order"`
	OrderId          uint                                `json:"order_id"`
	SourceToken      string                              `json:"source_token"`
	SourceChain      string                              `json:"source_chain"`
	SourceChainNames []string                            `json:"source_chain_names"`
	Sender           wallets.Wallets                     `json:"sender"`
	TargetToken      string                              `json:"target_token"`
	Receiver         string                              `json:"receiver"`
	TargetChain      string                              `json:"target_chain"`
	Amount           decimal.Decimal                     `json:"amount"`
	LastHistory      SwapHistory                         `json:"last_history"`
	Remark           string                              `json:"remark"`
	RecordFn         func(s SwapHistory, errs ...error)  `json:"-"`
	SaveOrderFn      func(update map[string]interface{}) `json:"-"`
	CurrentBalance   decimal.Decimal                     `json:"current_balance"`
}

type SwapResult struct {
	Error            string               `json:"error,omitempty"`
	TokenInChainName string               `json:"source_chain_name,omitempty"`
	TokenInName      string               `json:"token_in_name,omitempty"`
	ProviderType     configs.ProviderType `json:"type,omitempty"`
	ProviderName     string               `json:"provider_name,omitempty"`
	OrderId          string               `json:"provider_order_id,omitempty"`
	Order            interface{}          `json:"order,omitempty"`
	Status           TxStatus             `json:"status,omitempty"`
	CurrentChain     string               `json:"current_chain_name,omitempty"`
	Receiver         string               `json:"receiver,omitempty"`
	// Tx is the transaction hash
	Tx string `json:"tx,omitempty"`
}

func (s *SwapResult) SetReciever(receiver string) *SwapResult {
	s.Receiver = receiver
	return s
}

func (s *SwapResult) SetError(err error) *SwapResult {
	s.Error = err.Error()
	s.Status = TxStatusFailed
	return s
}
func (s *SwapResult) SetTokenInChainName(name string) *SwapResult {
	s.TokenInChainName = name
	return s
}

func (s *SwapResult) SetTokenInName(name string) *SwapResult {
	s.TokenInName = name
	return s
}

func (s *SwapResult) SetProviderType(tp configs.ProviderType) *SwapResult {
	s.ProviderType = tp
	return s
}

func (s *SwapResult) SetProviderName(name string) *SwapResult {
	s.ProviderName = name
	return s
}

func (s *SwapResult) SetOrderId(id string) *SwapResult {
	s.OrderId = id
	return s
}

func (s *SwapResult) SetOrder(order interface{}) *SwapResult {
	s.Order = order
	if _, ok := order.(string); ok {
		s.Order = map[string]interface{}{
			"data": order,
		}
	}
	return s
}

func (s *SwapResult) SetStatus(status TxStatus) *SwapResult {
	s.Status = status
	return s
}

func (s *SwapResult) SetCurrentChain(name string) *SwapResult {
	s.CurrentChain = name
	return s
}

func (s *SwapResult) SetTx(tx string) *SwapResult {
	s.Tx = tx
	return s
}

func (s *SwapResult) Out() SwapResult {
	return *s
}
func (s *SwapResult) OutError(err error) (SwapResult, error) {
	s.Error = err.Error()
	s.Status = TxStatusFailed
	return *s, err
}

func (s SwapResult) Marshal() map[string]interface{} {
	data, _ := json.Marshal(s)
	var result = make(map[string]interface{})
	_ = json.Unmarshal(data, &result)
	return result
}

func (s SwapResult) MarshalOrder() []byte {
	if s.Order == nil {
		return nil
	}
	b, err := json.Marshal(s.Order)
	if err != nil {
		return nil
	}
	return b
}

type SwapHistory struct {
	ProviderName string
	ProviderType string
	Actions      string
	Status       string
	CurrentChain string
	Amount       decimal.Decimal
	Tx           string `json:"tx"`
}

func (s *SwapHistory) Out() SwapHistory {
	return *s
}

func (s *SwapHistory) SetProviderName(providerName string) *SwapHistory {
	s.ProviderName = providerName
	return s
}

func (s *SwapHistory) SetProviderType(providerType configs.ProviderType) *SwapHistory {
	s.ProviderType = string(providerType)
	return s
}

func (s *SwapHistory) SetActions(actions string) *SwapHistory {
	s.Actions = actions
	return s
}

func (s *SwapHistory) SetStatus(status TxStatus) *SwapHistory {
	s.Status = status.String()
	return s
}

func (s *SwapHistory) SetCurrentChain(currentChain string) *SwapHistory {
	s.CurrentChain = currentChain
	return s
}

func (s *SwapHistory) SetAmount(amount decimal.Decimal) *SwapHistory {
	s.Amount = amount
	return s
}

func (s *SwapHistory) SetTx(tx string) *SwapHistory {
	s.Tx = tx
	return s
}

func (s SwapParams) Clone() SwapParams {
	return SwapParams{
		SourceToken: s.SourceToken,
		SourceChain: s.SourceChain,
		Sender:      s.Sender,
		TargetToken: s.TargetToken,
		Receiver:    s.Receiver,
		TargetChain: s.TargetChain,
		Amount:      s.Amount,
		LastHistory: s.LastHistory,
		RecordFn:    s.RecordFn,
	}
}
