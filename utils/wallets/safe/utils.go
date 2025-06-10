package safe

import (
	"context"
	"fmt"
	"math/big"
	"net/url"
	"omni-balance/utils"
	"omni-balance/utils/chains"
	"omni-balance/utils/constant"
	"omni-balance/utils/safe_api"
	apiclient "omni-balance/utils/safe_api/client"
	"omni-balance/utils/safe_api/client/safes"
	"omni-balance/utils/safe_api/client/transactions"
	"omni-balance/utils/wallets/safe/safe_abi"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	log "omni-balance/utils/logging"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient/simulated"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/spf13/cast"
)

var (
	safeDomain = map[string]string{
		constant.Arbitrum:     "safe-transaction-arbitrum.safe.global",
		constant.Avalanche:    "safe-transaction-avalanche.safe.global",
		constant.Base:         "safe-transaction-base.safe.global",
		constant.Bnb:          "safe-transaction-bsc.safe.global",
		constant.Celo:         "safe-transaction-celo.safe.global",
		constant.Ethereum:     "safe-transaction-mainnet.safe.global",
		constant.Gnosis:       "safe-transaction-gnosis-chain.safe.global",
		constant.Optimism:     "safe-transaction-optimism.safe.global",
		constant.Polygon:      "safe-transaction-polygon.safe.global",
		constant.PolygonZkEvm: "safe-transaction-zkevm.safe.global",
		constant.Zksync:       "safe-transaction-zkevm.safe.global",
		constant.Sepolia:      "safe-transaction-sepolia.safe.global",
		constant.Scroll:       "safe-transaction-scroll.safe.global",
	}
	safeGlobalLocker sync.Mutex
)

type SafeResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Info struct {
	SafeResp
	Address         string        `json:"address"`
	Nonce           int           `json:"nonce"`
	Threshold       int           `json:"threshold"`
	Owners          []string      `json:"owners"`
	MasterCopy      string        `json:"masterCopy"`
	Modules         []interface{} `json:"modules"`
	FallbackHandler string        `json:"fallbackHandler"`
	Guard           string        `json:"guard"`
	Version         string        `json:"version"`
}

type Transaction struct {
	Detail                string          `json:"detail"`
	Safe                  common.Address  `json:"safe"`
	To                    common.Address  `json:"to"`
	Value                 decimal.Decimal `json:"value"`
	Data                  string          `json:"data"`
	Operation             int             `json:"operation"`
	GasToken              common.Address  `json:"gasToken"`
	SafeTxGas             decimal.Decimal `json:"safeTxGas"`
	BaseGas               decimal.Decimal `json:"baseGas"`
	GasPrice              decimal.Decimal `json:"gasPrice"`
	RefundReceiver        common.Address  `json:"refundReceiver"`
	Nonce                 int             `json:"nonce"`
	TransactionHash       interface{}     `json:"transactionHash"`
	SafeTxHash            common.Hash     `json:"safeTxHash"`
	Proposer              string          `json:"proposer"`
	Executor              interface{}     `json:"executor"`
	IsExecuted            bool            `json:"isExecuted"`
	IsSuccessful          bool            `json:"isSuccessful"`
	EthGasPrice           interface{}     `json:"ethGasPrice"`
	MaxFeePerGas          interface{}     `json:"maxFeePerGas"`
	MaxPriorityFeePerGas  interface{}     `json:"maxPriorityFeePerGas"`
	GasUsed               interface{}     `json:"gasUsed"`
	Fee                   interface{}     `json:"fee"`
	Origin                string          `json:"origin"`
	DataDecoded           interface{}     `json:"dataDecoded"`
	ConfirmationsRequired int             `json:"confirmationsRequired"`
	Confirmations         []struct {
		Owner           string      `json:"owner"`
		SubmissionDate  time.Time   `json:"submissionDate"`
		TransactionHash interface{} `json:"transactionHash"`
		Signature       string      `json:"signature"`
		SignatureType   string      `json:"signatureType"`
	} `json:"confirmations"`
	Trusted    bool        `json:"trusted"`
	Signatures interface{} `json:"signatures"`
}

func (s *Safe) GetDomainByCtx(ctx context.Context) string {
	chainName := cast.ToString(ctx.Value(constant.ChainNameKeyInCtx))
	if chainName == "" {
		debug.PrintStack()
		log.Fatalf("chain name not found in context")
	}
	if _, ok := safeDomain[chainName]; ok {
		return safeDomain[chainName]
	}
	debug.PrintStack()
	log.Fatalf("chain name %s not found in safe domain", chainName)
	return ""
}

func (s *Safe) GetChainIdByCtx(ctx context.Context) int {
	chainId := constant.GetChainId(cast.ToString(ctx.Value(constant.ChainNameKeyInCtx)))
	if chainId == 0 {
		debug.PrintStack()
		log.Fatalf("chain name not found in context")
	}
	return chainId
}

func (s *Safe) safeWalletInfo(ctx context.Context) (*safe_api.SafeInfoResponse, error) {
	var address = s.getOperatorSafeAddress().Hex()

	resp, err := s.Client(ctx).Safes.V1SafesRead(&safes.V1SafesReadParams{Address: address}, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "get %s safe info error", address)
	}
	return resp.Payload, nil
}

func (s *Safe) Client(ctx context.Context) *apiclient.SafeTransactionServiceAPI {
	c := httptransport.New(
		s.GetDomainByCtx(ctx), apiclient.DefaultBasePath, apiclient.DefaultSchemes)
	c.Context = ctx
	return apiclient.New(c, strfmt.Default)
}

func (s *Safe) Transfer(ctx context.Context, tx *types.LegacyTx, client simulated.Client) (common.Hash, error) {
	if s.operatorSafe != nil {
		return s.operatorSafe.SendTransaction(ctx, tx, client)
	}
	return chains.SendTransaction(ctx, client, tx, s.GetAddress(true), s.conf.Operator.PrivateKey)
}

func (s *Safe) MultisigTransaction(ctx context.Context, tx *types.LegacyTx, client simulated.Client) (common.Hash, error) {
	if tx.Nonce == 0 {
		nonce, err := s.nonce(ctx)
		if err != nil {
			return common.Hash{}, errors.Wrap(err, "get nonce error")
		}
		tx.Nonce = uint64(nonce)
		log.Debugf("%s nonce: %d", s.getOperatorSafeAddress(), tx.Nonce)
	}

	to, _, err := chains.GetTransferInfo(tx.Data)
	if err == nil && s.GetAddress().Cmp(to) == 0 {
		log.Debugf("transfer to self")
		return s.Transfer(ctx, tx, client)
	}

	// _, err = client.EstimateGas(ctx, ethereum.CallMsg{
	// 	From:  s.GetAddress(),
	// 	To:    tx.To,
	// 	Value: tx.Value,
	// 	Data:  tx.Data,
	// })
	// if err != nil && !strings.Contains(err.Error(), "insufficient funds for gas * price + value") {
	// 	log.Debugf("estimate gas error: %s", err.Error())
	// 	return common.Hash{}, errors.Wrap(err, "estimate gas error")
	// }

	info, err := s.safeWalletInfo(ctx)
	if err != nil {
		return common.Hash{}, err
	}
	safeTxHash, err := s.proposeTransaction(ctx, tx)
	if err != nil {
		return common.Hash{}, errors.Wrap(err, "propose transaction error")
	}
	log.Debugf("safe tx hash: %s", safeTxHash)
	if *info.Threshold > 1 {
		return safeTxHash, nil
	}

	txInfo, err := s.GetMultiSigTransaction(ctx, safeTxHash)
	if err != nil {
		return safeTxHash, errors.Wrap(err, "get multisig transaction error")
	}
	if txInfo.IsExecuted {
		return safeTxHash, errors.New("transaction already executed")
	}
	return safeTxHash, s.ExecTransaction(ctx, txInfo, client)
}

func (s *Safe) getOperatorSafeAddress() common.Address {
	if s.conf.Operator.MultiSignType != "" {
		return s.conf.Operator.Address
	}
	return s.conf.Address
}

func (s *Safe) getOperatorAddress() common.Address {
	if s.conf.Operator.MultiSignType != "" {
		return s.conf.Operator.Operator
	}
	return s.conf.Operator.Address
}

func (s *Safe) nonce(ctx context.Context) (int64, error) {
	u := fmt.Sprintf("https://safe-client.safe.global/v1/chains/%d/safes/%s/nonces", s.GetChainIdByCtx(ctx), s.getOperatorSafeAddress().Hex())
	var dest = struct {
		RecommendedNonce int64 `json:"recommendedNonce"`
	}{}
	return dest.RecommendedNonce, utils.Request(ctx, "GET", u, nil, &dest)
}

type TxList struct {
	Count    int         `json:"count"`
	Next     interface{} `json:"next"`
	Previous interface{} `json:"previous"`
	Results  []struct {
		Safe                 string      `json:"safe"`
		To                   string      `json:"to"`
		Value                string      `json:"value"`
		Data                 string      `json:"data"`
		Operation            int         `json:"operation"`
		GasToken             string      `json:"gasToken"`
		SafeTxGas            int         `json:"safeTxGas"`
		BaseGas              int         `json:"baseGas"`
		GasPrice             string      `json:"gasPrice"`
		RefundReceiver       string      `json:"refundReceiver"`
		Nonce                int         `json:"nonce"`
		ExecutionDate        interface{} `json:"executionDate"`
		SubmissionDate       time.Time   `json:"submissionDate"`
		Modified             time.Time   `json:"modified"`
		BlockNumber          interface{} `json:"blockNumber"`
		TransactionHash      interface{} `json:"transactionHash"`
		SafeTxHash           string      `json:"safeTxHash"`
		Proposer             string      `json:"proposer"`
		ProposedByDelegate   interface{} `json:"proposedByDelegate"`
		Executor             interface{} `json:"executor"`
		IsExecuted           bool        `json:"isExecuted"`
		IsSuccessful         interface{} `json:"isSuccessful"`
		EthGasPrice          interface{} `json:"ethGasPrice"`
		MaxFeePerGas         interface{} `json:"maxFeePerGas"`
		MaxPriorityFeePerGas interface{} `json:"maxPriorityFeePerGas"`
		GasUsed              interface{} `json:"gasUsed"`
		Fee                  interface{} `json:"fee"`
		Origin               string      `json:"origin"`
		DataDecoded          struct {
			Method     string `json:"method"`
			Parameters []struct {
				Name  string `json:"name"`
				Type  string `json:"type"`
				Value string `json:"value"`
			} `json:"parameters"`
		} `json:"dataDecoded"`
		ConfirmationsRequired int `json:"confirmationsRequired"`
		Confirmations         []struct {
			Owner           string      `json:"owner"`
			SubmissionDate  time.Time   `json:"submissionDate"`
			TransactionHash interface{} `json:"transactionHash"`
			Signature       string      `json:"signature"`
			SignatureType   string      `json:"signatureType"`
		} `json:"confirmations"`
		Trusted    bool        `json:"trusted"`
		Signatures interface{} `json:"signatures"`
	} `json:"results"`
	CountUniqueNonce int `json:"countUniqueNonce"`
}

func (s *Safe) TxIsInQueue(ctx context.Context, tx *types.LegacyTx) (string, error) {
	u := fmt.Sprintf("https://%s/api/v1/safes/%s/multisig-transactions/?limit=100&ordering=-timestamp&executed=false",
		s.GetDomainByCtx(ctx), s.getOperatorSafeAddress().Hex())
	for {
		var resp TxList
		if err := utils.Request(ctx, "GET", u, nil, &resp); err != nil {
			return "", errors.Wrapf(err, "get tx list error")
		}
		for _, v := range resp.Results {
			if strings.EqualFold(v.Data, "0x"+common.Bytes2Hex(tx.Data)) &&
				v.Value == tx.Value.String() {
				return v.SafeTxHash, nil
			}
		}
		if resp.Next != nil || cast.ToString(resp.Next) != "" {
			u = cast.ToString(resp.Next)
			continue
		}
		if resp.Next == nil || cast.ToString(resp.Next) == "" {
			return "", nil
		}
	}
}

func (s *Safe) proposeTransaction(ctx context.Context, tx *types.LegacyTx) (common.Hash, error) {
	safeGlobalLocker.Lock()
	defer safeGlobalLocker.Unlock()
	safeHash, err := s.TxIsInQueue(ctx, tx)
	if err != nil {
		return common.Hash{}, errors.Wrapf(err, "check tx in queue error")
	}
	if safeHash != "" {
		return common.HexToHash(safeHash), nil
	}
	nonce, err := s.nonce(ctx)
	if err != nil {
		return common.Hash{}, errors.Wrap(err, "get nonce error")
	}
	t := &Transaction{
		Safe:  s.getOperatorSafeAddress(),
		To:    *tx.To,
		Data:  common.Bytes2Hex(tx.Data),
		Nonce: int(nonce),
	}
	if tx.Value != nil {
		t.Value = decimal.NewFromBigInt(tx.Value, 0)
	}

	typedData := s.eip712(ctx, *t)
	var safeTxHash string

	sigData, err := chains.SignTypedData(typedData, func(msg []byte) (sig []byte, err error) {
		safeTxHash = common.Bytes2Hex(msg)
		privateKey := s.conf.PrivateKey
		if s.conf.Operator.PrivateKey != "" {
			privateKey = s.conf.Operator.PrivateKey
		}
		return chains.SignMsg(msg, privateKey)
	})
	if err != nil {
		return common.Hash{}, errors.Wrap(err, "sign typed data error")
	}
	params := &transactions.V1SafesMultisigTransactionsCreateParams{
		Address: s.getOperatorSafeAddress().Hex(),
		Data: &safe_api.SafeMultisigTransaction{
			Safe:                    utils.String(s.getOperatorSafeAddress().Hex()),
			BaseGas:                 utils.Number(int64(0)),
			ContractTransactionHash: utils.String(fmt.Sprintf("0x%s", safeTxHash)),
			Data:                    utils.String(t.Data),
			GasPrice:                utils.Number(int64(0)),
			GasToken:                utils.String(constant.ZeroAddress.Hex()),
			Nonce:                   utils.Number(int64(t.Nonce)),
			Operation:               utils.Number(int64(0)),
			RefundReceiver:          utils.String(constant.ZeroAddress.Hex()),
			SafeTxGas:               utils.Number(int64(0)),
			Sender:                  utils.String(s.getOperatorAddress().Hex()),
			Signature:               utils.String(fmt.Sprintf("0x%s", common.Bytes2Hex(sigData))),
			To:                      utils.String(t.To.Hex()),
			Value:                   utils.Number(t.Value.IntPart()),
		},
		Context: ctx,
	}

	_, err = s.Client(ctx).Transactions.V1SafesMultisigTransactionsCreate(params, nil)
	if err != nil {
		return common.Hash{}, errors.Wrap(err, "create multisig transaction error")
	}
	return common.HexToHash(fmt.Sprintf("0x%s", safeTxHash)), nil
}

func (s *Safe) GetMultiSigTransaction(ctx context.Context, safeTxHash common.Hash) (Transaction, error) {
	u := fmt.Sprintf("https://%s/api/v1/multisig-transactions/%s/", s.GetDomainByCtx(ctx), safeTxHash.Hex())
	var result Transaction
	if err := utils.Request(ctx, "GET", u, nil, &result); err != nil {
		return result, errors.Wrap(err, "get multisig transaction error")
	}
	if result.Detail != "" {
		return result, errors.New(result.Detail)
	}
	return result, nil
}

func (s *Safe) ListMultiSigTransactions(ctx context.Context, filters ...string) ([]Transaction, error) {
	// executed=false&ordering=nonce
	u, err := url.Parse(fmt.Sprintf("https://%s/api/v1/safes/%s/multisig-transactions/", s.GetDomainByCtx(ctx), "0x0350101f2cB6aA65caaB7954246a56f906A3F57D"))
	if err != nil {
		return nil, errors.Wrap(err, "parse url error")
	}
	query := u.Query()
	var (
		key string
	)
	for index, v := range filters {
		if index%2 == 0 {
			key = v
			continue
		}
		query.Add(key, v)
	}
	u.RawQuery = query.Encode()
	var result = struct {
		Count  int           `json:"count"`
		Result []Transaction `json:"results"`
	}{}

	if err := utils.Request(ctx, "GET", u.String(), nil, &result); err != nil {
		return nil, errors.Wrap(err, "get multisig transactions error")
	}
	return result.Result, nil
}

func (s *Safe) ExecTransaction(ctx context.Context, tx Transaction, client simulated.Client) error {
	abi, err := safe_abi.SafeAbiMetaData.GetAbi()
	if err != nil {
		return errors.Wrap(err, "get abi error")
	}
	nonce, err := client.NonceAt(context.TODO(), s.getOperatorAddress(), nil)
	if err != nil {
		return errors.Wrap(err, "get nonce error")
	}
	var (
		signatures []string
		hasSelf    = false
	)

	for _, v := range tx.Confirmations {
		if strings.EqualFold(v.Owner, s.getOperatorAddress().Hex()) {
			hasSelf = true
		}
		signatures = append(signatures, strings.TrimPrefix(v.Signature, "0x"))
	}

	var signatureData []byte
	if !hasSelf && len(signatures) != 0 {
		address := common.Bytes2Hex(common.LeftPadBytes(s.GetAddress(true).Bytes(), 32))
		confriom := common.Bytes2Hex(common.LeftPadBytes(big.NewInt(1).Bytes(), 33))
		signatureData = common.Hex2Bytes(fmt.Sprintf("%s%s", address, confriom))
	} else if len(signatures) == 0 {
		signatureData, err = chains.SignTypedData(s.eip712(ctx, tx), func(msg []byte) (sig []byte, err error) {
			return chains.SignMsg(msg, s.conf.Operator.PrivateKey)
		})
		if err != nil {
			return errors.Wrap(err, "sign error")
		}
	}

	if len(signatures) < tx.ConfirmationsRequired && hasSelf {
		return errors.New("not enough signatures")
	}

	if len(signatures) < tx.ConfirmationsRequired && len(signatureData) > 0 {
		signatures = append(signatures, common.Bytes2Hex(signatureData))
	}

	if len(signatures) < tx.ConfirmationsRequired {
		return errors.New("not enough signatures")
	}
	input, err := abi.Pack("execTransaction",
		tx.To,
		tx.Value.BigInt(),
		common.Hex2Bytes(strings.TrimPrefix(tx.Data, "0x")),
		uint8(tx.Operation),
		tx.SafeTxGas.BigInt(),
		tx.BaseGas.BigInt(),
		tx.GasPrice.BigInt(),
		tx.GasToken,
		tx.RefundReceiver,
		common.Hex2Bytes(strings.Join(signatures, "")),
	)
	if err != nil {
		return errors.Wrap(err, "pack error")
	}
	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return errors.Wrap(err, "suggest gas price error")
	}
	gas, err := client.EstimateGas(context.TODO(), ethereum.CallMsg{
		From:     s.getOperatorAddress(),
		To:       &tx.Safe,
		GasPrice: gasPrice,
		Data:     input,
	})
	if err != nil {
		return errors.Wrap(err, "estimate gas error")
	}
	chainTransaction := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		GasPrice: gasPrice,
		Gas:      gas,
		To:       &tx.Safe,
		Data:     input,
	})

	signTx, err := chains.SignTx(chainTransaction, s.conf.Operator.PrivateKey, int64(s.GetChainIdByCtx(ctx)))
	if err != nil {
		return errors.Wrap(err, "sign tx error")
	}
	err = client.SendTransaction(ctx, signTx)
	if err != nil {
		return errors.Wrap(err, "send tx error")
	}
	return nil
}

func (s *Safe) eip712(ctx context.Context, t Transaction) apitypes.TypedData {
	return apitypes.TypedData{
		Types: apitypes.Types{
			"EIP712Domain": []apitypes.Type{
				{Name: "chainId", Type: "uint256"},
				{Name: "verifyingContract", Type: "address"},
			},
			"SafeTx": []apitypes.Type{
				{Type: "address", Name: "to"},
				{Type: "uint256", Name: "value"},
				{Type: "bytes", Name: "data"},
				{Type: "uint8", Name: "operation"},
				{Type: "uint256", Name: "safeTxGas"},
				{Type: "uint256", Name: "baseGas"},
				{Type: "uint256", Name: "gasPrice"},
				{Type: "address", Name: "gasToken"},
				{Type: "address", Name: "refundReceiver"},
				{Type: "uint256", Name: "nonce"},
			},
		},
		PrimaryType: "SafeTx",
		Domain: apitypes.TypedDataDomain{
			ChainId:           math.NewHexOrDecimal256(int64(s.GetChainIdByCtx(ctx))),
			VerifyingContract: t.Safe.Hex(),
		},
		Message: apitypes.TypedDataMessage{
			"to":             t.To.Hex(),
			"value":          t.Value.BigInt(),
			"data":           common.Hex2Bytes(t.Data),
			"operation":      math.NewHexOrDecimal256(int64(t.Operation)),
			"baseGas":        math.NewHexOrDecimal256(t.BaseGas.IntPart()),
			"gasPrice":       math.NewHexOrDecimal256(t.GasPrice.IntPart()),
			"gasToken":       t.GasToken.Hex(),
			"refundReceiver": t.RefundReceiver.Hex(),
			"nonce":          math.NewHexOrDecimal256(int64(t.Nonce)),
			"safeTxGas":      math.NewHexOrDecimal256(t.SafeTxGas.IntPart()),
		},
	}
}
