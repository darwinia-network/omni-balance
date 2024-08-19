package bungee

import (
	"omni-balance/utils"
	"omni-balance/utils/constant"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

var (
	ZeroAddress = common.HexToAddress("0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee")
	// see https://docs.bungee.exchange/socket-api/support/chains
	SupportedChain = map[string]struct{}{
		constant.Ethereum:     {},
		constant.Bsc:          {},
		constant.Fantom:       {},
		constant.Avalanche:    {},
		constant.Optimism:     {},
		constant.Arbitrum:     {},
		constant.Gnosis:       {},
		constant.Polygon:      {},
		constant.Aurora:       {},
		constant.Zksync:       {},
		constant.PolygonZkEvm: {},
		constant.Base:         {},
		constant.Linea:        {},
		constant.Scroll:       {},
		constant.Blast:        {},
		constant.Mantle:       {},
	}
	// see https://docs.bungee.exchange/socket-api/introduction#want-bungee-to-integrate-your-protocol
	apiKey = utils.GetEnv("BUNGEE_API_KEY", "1b2fd225-062f-41aa-8c63-d1fef19945e7")
)

var (
	ApproveTransactionAction = "ApproveTransaction"
	SourceChainSendingAction = "sourceChainSending"
	WaitForTxAction          = "WaitForTx"
)

type BaseResp struct {
	ErrorMsg  string `json:"error,omitempty"`
	ErrorCode string `json:"error_code,omitempty"`
}

type Quote struct {
	BaseResp
	FlowType       string      `json:"flowType"`
	IsTransfer     interface{} `json:"isTransfer"`
	IsWrappedToken bool        `json:"isWrappedToken"`
	AllowanceTo    string      `json:"allowanceTo"`
	BridgeFee      struct {
		Amount   decimal.Decimal `json:"amount"`
		Decimals int             `json:"decimals"`
		Symbol   string          `json:"symbol"`
		Address  string          `json:"address"`
	} `json:"bridgeFee"`
	FuelTransfer     interface{}    `json:"fuelTransfer"`
	FromTokenAddress common.Address `json:"fromTokenAddress"`
	ToTokenAddress   common.Address `json:"toTokenAddress"`
	Source           struct {
		ChainId   string `json:"chainId"`
		ChainType string `json:"chainType"`
		Asset     struct {
			Decimals       int    `json:"decimals"`
			Symbol         string `json:"symbol"`
			Name           string `json:"name"`
			ChainId        string `json:"chainId"`
			Address        string `json:"address"`
			ResourceID     string `json:"resourceID"`
			IsMintable     bool   `json:"isMintable"`
			IsWrappedAsset bool   `json:"isWrappedAsset"`
			IsReserveAsset bool   `json:"isReserveAsset"`
			TokenInstance  struct {
				Decimals int    `json:"decimals"`
				Symbol   string `json:"symbol"`
				Name     string `json:"name"`
				ChainId  int    `json:"chainId"`
				Address  string `json:"address"`
			} `json:"tokenInstance"`
		} `json:"asset"`
		StableReserveAsset struct {
			Decimals       int    `json:"decimals"`
			Symbol         string `json:"symbol"`
			Name           string `json:"name"`
			ChainId        string `json:"chainId"`
			Address        string `json:"address"`
			ResourceID     string `json:"resourceID"`
			IsMintable     *bool  `json:"isMintable,omitempty"`
			IsWrappedAsset *bool  `json:"isWrappedAsset,omitempty"`
			IsReserveAsset *bool  `json:"isReserveAsset,omitempty"`
		} `json:"stableReserveAsset"`
		TokenAmount         decimal.Decimal `json:"tokenAmount"`
		StableReserveAmount string          `json:"stableReserveAmount"`
		Path                []interface{}   `json:"path"`
		Flags               []interface{}   `json:"flags"`
		PriceImpact         string          `json:"priceImpact"`
		TokenPath           string          `json:"tokenPath"`
		DataTx              []string        `json:"dataTx"`
	} `json:"source"`
	Destination struct {
		ChainId string `json:"chainId"`
		Asset   struct {
			Decimals       int    `json:"decimals"`
			Symbol         string `json:"symbol"`
			Name           string `json:"name"`
			ChainId        string `json:"chainId"`
			Address        string `json:"address"`
			ResourceID     string `json:"resourceID"`
			IsMintable     bool   `json:"isMintable"`
			IsWrappedAsset bool   `json:"isWrappedAsset"`
			IsReserveAsset bool   `json:"isReserveAsset"`
			TokenInstance  struct {
				Decimals int    `json:"decimals"`
				Symbol   string `json:"symbol"`
				Name     string `json:"name"`
				ChainId  int    `json:"chainId"`
				Address  string `json:"address"`
			} `json:"tokenInstance"`
		} `json:"asset"`
		StableReserveAsset struct {
			Decimals       int    `json:"decimals"`
			Symbol         string `json:"symbol"`
			Name           string `json:"name"`
			ChainId        string `json:"chainId"`
			Address        string `json:"address"`
			ResourceID     string `json:"resourceID"`
			IsMintable     bool   `json:"isMintable"`
			IsWrappedAsset bool   `json:"isWrappedAsset"`
			IsReserveAsset bool   `json:"isReserveAsset"`
		} `json:"stableReserveAsset"`
		TokenAmount decimal.Decimal `json:"tokenAmount"`
	} `json:"destination"`
	PartnerId         interface{} `json:"partnerId"`
	SlippageTolerance interface{} `json:"slippageTolerance"`
	EstimatedTime     interface{} `json:"estimatedTime"`
}

type Txn struct {
	Quote
	Txn struct {
		From     common.Address  `json:"from"`
		To       common.Address  `json:"to"`
		Data     string          `json:"data"`
		Value    string          `json:"value"`
		GasPrice decimal.Decimal `json:"gasPrice"`
		GasLimit decimal.Decimal `json:"gasLimit"`
	} `json:"txn"`
}

func (b BaseResp) Error() error {
	if b.ErrorCode != "" {
		return errors.New(b.ErrorMsg)
	}
	return nil
}
func (q Quote) Error() error {
	if err := q.BaseResp.Error(); err != nil {
		return err
	}
	if q.FlowType == "" {
		return errors.New("flowType is empty")
	}
	if q.Destination.Asset.Symbol == "" {
		return errors.New("destination.asset.symbol is empty")
	}
	return nil
}

func Action2Int(action string) int {
	switch action {
	case ApproveTransactionAction:
		return 1
	case SourceChainSendingAction:
		return 2
	case WaitForTxAction:
		return 3
	default:
		return 0
	}
}
