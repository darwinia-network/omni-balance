package mantle_safe

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"omni-balance/utils"
	"omni-balance/utils/chains"
	"omni-balance/utils/constant"
	"omni-balance/utils/wallets/safe"
	"omni-balance/utils/wallets/safe/safe_abi"
	"strings"
	"sync"

	log "omni-balance/utils/logging"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient/simulated"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/spf13/cast"
	"github.com/tidwall/gjson"
)

var (
	safeGlobalLocker sync.Mutex
)

type SafeResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ProposeTransaction struct {
	To             string `json:"to"`
	Value          string `json:"value"`
	Data           string `json:"data"`
	Operation      int    `json:"operation"`
	BaseGas        string `json:"baseGas"`
	GasPrice       string `json:"gasPrice"`
	GasToken       string `json:"gasToken"`
	RefundReceiver string `json:"refundReceiver"`
	Nonce          string `json:"nonce"`
	SafeTxGas      string `json:"safeTxGas"`
	SafeTxHash     string `json:"safeTxHash"`
	Sender         string `json:"sender"`
	Signature      string `json:"signature"`
}

type Transaction struct {
	SafeResp
	SafeAddress string `json:"safeAddress"`
	TxID        string `json:"txId"`
	ExecutedAt  *int64 `json:"executedAt"`
	TxStatus    string `json:"txStatus"`
	TxInfo      struct {
		Type             string `json:"type"`
		HumanDescription any    `json:"humanDescription"`
		Sender           struct {
			Value   string `json:"value"`
			Name    any    `json:"name"`
			LogoURI any    `json:"logoUri"`
		} `json:"sender"`
		Recipient struct {
			Value   string `json:"value"`
			Name    any    `json:"name"`
			LogoURI any    `json:"logoUri"`
		} `json:"recipient"`
		Direction    string `json:"direction"`
		TransferInfo struct {
			Type  string `json:"type"`
			Value string `json:"value"`
		} `json:"transferInfo"`
	} `json:"txInfo"`
	TxData struct {
		HexData     any `json:"hexData"`
		DataDecoded any `json:"dataDecoded"`
		To          struct {
			Value   string `json:"value"`
			Name    any    `json:"name"`
			LogoURI any    `json:"logoUri"`
		} `json:"to"`
		Value                     string `json:"value"`
		Operation                 int    `json:"operation"`
		TrustedDelegateCallTarget any    `json:"trustedDelegateCallTarget"`
		AddressInfoIndex          any    `json:"addressInfoIndex"`
	} `json:"txData"`
	TxHash                string `json:"txHash"`
	DetailedExecutionInfo struct {
		Type           string `json:"type"`
		SubmittedAt    int64  `json:"submittedAt"`
		Nonce          int    `json:"nonce"`
		SafeTxGas      string `json:"safeTxGas"`
		BaseGas        string `json:"baseGas"`
		GasPrice       string `json:"gasPrice"`
		GasToken       string `json:"gasToken"`
		RefundReceiver struct {
			Value   string `json:"value"`
			Name    any    `json:"name"`
			LogoURI any    `json:"logoUri"`
		} `json:"refundReceiver"`
		SafeTxHash string `json:"safeTxHash"`
		Executor   struct {
			Value   string `json:"value"`
			Name    any    `json:"name"`
			LogoURI any    `json:"logoUri"`
		} `json:"executor"`
		Signers []struct {
			Value   string `json:"value"`
			Name    any    `json:"name"`
			LogoURI any    `json:"logoUri"`
		} `json:"signers"`
		ConfirmationsRequired int `json:"confirmationsRequired"`
		Confirmations         []struct {
			Signer struct {
				Value   string `json:"value"`
				Name    any    `json:"name"`
				LogoURI any    `json:"logoUri"`
			} `json:"signer"`
			Signature   string `json:"signature"`
			SubmittedAt int64  `json:"submittedAt"`
		} `json:"confirmations"`
		Rejectors    []any `json:"rejectors"`
		GasTokenInfo any   `json:"gasTokenInfo"`
		Trusted      bool  `json:"trusted"`
		Proposer     struct {
			Value   string `json:"value"`
			Name    any    `json:"name"`
			LogoURI any    `json:"logoUri"`
		} `json:"proposer"`
	} `json:"detailedExecutionInfo"`
	SafeAppInfo any `json:"safeAppInfo"`
}

type Info struct {
	SafeResp
	Address struct {
		Value   string `json:"value"`
		Name    any    `json:"name"`
		LogoURI any    `json:"logoUri"`
	} `json:"address"`
	ChainID   string `json:"chainId"`
	Nonce     int    `json:"nonce"`
	Threshold int    `json:"threshold"`
	Owners    []struct {
		Value   string `json:"value"`
		Name    any    `json:"name"`
		LogoURI any    `json:"logoUri"`
	} `json:"owners"`
	Implementation struct {
		Value   string `json:"value"`
		Name    string `json:"name"`
		LogoURI string `json:"logoUri"`
	} `json:"implementation"`
	ImplementationVersionState string `json:"implementationVersionState"`
	CollectiblesTag            any    `json:"collectiblesTag"`
	TxQueuedTag                string `json:"txQueuedTag"`
	TxHistoryTag               string `json:"txHistoryTag"`
	MessagesTag                any    `json:"messagesTag"`
	Modules                    any    `json:"modules"`
	FallbackHandler            struct {
		Value   string `json:"value"`
		Name    string `json:"name"`
		LogoURI string `json:"logoUri"`
	} `json:"fallbackHandler"`
	Guard   any    `json:"guard"`
	Version string `json:"version"`
}

func (s *MantleSafe) GetDomainByCtx(ctx context.Context) string {
	return constant.Mantle
}

func (s *MantleSafe) GetChainIdByCtx(ctx context.Context) int {
	return constant.GetChainId(s.GetDomainByCtx(ctx))
}

func (s *MantleSafe) safeWalletInfo(ctx context.Context) (*Info, error) {
	var address = s.getOperatorSafeAddress().Hex()
	var result = new(Info)
	u := fmt.Sprintf("https://gateway.multisig.mantle.xyz/v1/chains/5000/safes/%s", address)
	if err := utils.Request(ctx, "GET", u, nil, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *MantleSafe) Transfer(ctx context.Context, tx *types.LegacyTx, client simulated.Client) (common.Hash, error) {
	if s.operatorSafe != nil {
		return s.operatorSafe.SendTransaction(ctx, tx, client)
	}
	return chains.SendTransaction(ctx, client, tx, s.GetAddress(true), s.conf.Operator.PrivateKey)
}

func (s *MantleSafe) MultisigTransaction(ctx context.Context, tx *types.LegacyTx, client simulated.Client) (common.Hash, error) {
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
	// if err != nil {
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
	if info.Threshold > 1 {
		return safeTxHash, nil
	}

	txInfo, err := s.GetMultiSigTransaction(ctx, safeTxHash)
	if err != nil {
		return safeTxHash, errors.Wrap(err, "get multisig transaction error")
	}
	if txInfo.ExecutedAt != nil && *txInfo.ExecutedAt > 0 {
		return safeTxHash, errors.New("transaction already executed")
	}
	return safeTxHash, s.ExecTransaction(ctx, txInfo, client)
}

func (s *MantleSafe) getOperatorSafeAddress() common.Address {
	if s.conf.Operator.MultiSignType != "" {
		return s.conf.Operator.Address
	}
	return s.conf.Address
}

func (s *MantleSafe) getOperatorAddress() common.Address {
	if s.conf.Operator.MultiSignType != "" {
		return s.conf.Operator.Operator
	}
	return s.conf.Operator.Address
}

func (s *MantleSafe) nonce(ctx context.Context) (int64, error) {
	u := fmt.Sprintf("https://gateway.multisig.mantle.xyz/v1/chains/5000/safes/%s/nonces", s.getOperatorSafeAddress().Hex())
	var dest = struct {
		RecommendedNonce int64 `json:"recommendedNonce"`
	}{}
	return dest.RecommendedNonce, utils.Request(ctx, "GET", u, nil, &dest)
}

func (s *MantleSafe) proposeTransaction(ctx context.Context, tx *types.LegacyTx) (common.Hash, error) {
	safeGlobalLocker.Lock()
	defer safeGlobalLocker.Unlock()
	nonce, err := s.nonce(ctx)
	if err != nil {
		return common.Hash{}, errors.Wrap(err, "get nonce error")
	}

	t := &safe.Transaction{
		Safe:  s.GetAddress(true),
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
	txData, err := json.Marshal(ProposeTransaction{
		To:             t.To.Hex(),
		Value:          t.Value.String(),
		Data:           "0x" + t.Data,
		BaseGas:        "0",
		GasPrice:       "0",
		GasToken:       constant.ZeroAddress.Hex(),
		RefundReceiver: constant.ZeroAddress.Hex(),
		Nonce:          cast.ToString(t.Nonce),
		SafeTxGas:      "0",
		SafeTxHash:     fmt.Sprintf("0x%s", safeTxHash),
		Sender:         s.getOperatorAddress().Hex(),
		Signature:      "0x" + common.Bytes2Hex(sigData),
	})
	if err != nil {
		return common.Hash{}, errors.Wrap(err, "marshal transaction error")
	}

	data, err := utils.RequestBinary(
		ctx, "POST",
		fmt.Sprintf("https://gateway.multisig.mantle.xyz/v1/chains/5000/transactions/%s/propose", s.GetAddress(true).Hex()),
		bytes.NewReader(txData))
	if err != nil {
		return common.Hash{}, nil
	}
	result := gjson.ParseBytes(data)
	if result.Get("statusCode").Int() != 0 || result.Get("code").Int() != 0 {
		return common.Hash{}, errors.New(result.Get("message").String())
	}

	return common.HexToHash(fmt.Sprintf("0x%s", safeTxHash)), nil
}

func (s *MantleSafe) GetMultiSigTransaction(ctx context.Context, safeTxHash common.Hash) (Transaction, error) {
	u := fmt.Sprintf("https://gateway.multisig.mantle.xyz/v1/chains/5000/transactions/multisig_%s_%s", s.getOperatorSafeAddress().Hex(), safeTxHash.Hex())
	var result Transaction
	if err := utils.Request(ctx, "GET", u, nil, &result); err != nil {
		return result, errors.Wrap(err, "get multisig transaction error")
	}
	if result.Code != 0 {
		return result, errors.New(result.Message)
	}
	return result, nil
}

func (s *MantleSafe) ExecTransaction(ctx context.Context, tx Transaction, client simulated.Client) error {
	abi, err := safe_abi.SafeAbiMetaData.GetAbi()
	if err != nil {
		return errors.Wrap(err, "get abi error")
	}
	nonce, err := client.NonceAt(context.TODO(), s.GetAddress(true), nil)
	if err != nil {
		return errors.Wrap(err, "get nonce error")
	}
	var (
		signatures []string
		hasSelf    = false
	)

	for _, v := range tx.DetailedExecutionInfo.Confirmations {
		if strings.EqualFold(v.Signer.Value, s.GetAddress(true).Hex()) {
			hasSelf = true
		}
		signatures = append(signatures, strings.TrimPrefix(v.Signature, "0x"))
	}

	var signatureData []byte
	if !hasSelf && len(signatures) != 0 {
		address := common.Bytes2Hex(common.LeftPadBytes(s.GetAddress(true).Bytes(), 32))
		confriom := common.Bytes2Hex(common.LeftPadBytes(big.NewInt(1).Bytes(), 33))
		signatureData = common.Hex2Bytes(fmt.Sprintf("%s%s", address, confriom))
	} else {
		proposeTx := &safe.Transaction{
			Safe:  s.GetAddress(true),
			To:    common.HexToAddress(tx.TxData.To.Value),
			Data:  "0x",
			Nonce: int(tx.DetailedExecutionInfo.Nonce),
		}
		if tx.TxData.HexData != nil {
			proposeTx.Data = cast.ToString(tx.TxData.HexData)
		}
		signatureData, err = chains.SignTypedData(s.eip712(ctx, *proposeTx), func(msg []byte) (sig []byte, err error) {
			return chains.SignMsg(msg, s.conf.PrivateKey)
		})
		if err != nil {
			return errors.Wrap(err, "sign error")
		}
	}

	if len(signatures) < tx.DetailedExecutionInfo.ConfirmationsRequired && hasSelf {
		return errors.New("not enough signatures")
	}

	if len(signatures) < tx.DetailedExecutionInfo.ConfirmationsRequired {
		signatures = append(signatures, common.Bytes2Hex(signatureData))
	}

	if len(signatures) < tx.DetailedExecutionInfo.ConfirmationsRequired {
		return errors.New("not enough signatures")
	}
	data := "0x"
	if tx.TxData.HexData != nil {
		data = cast.ToString(tx.TxData.HexData)
	}
	input, err := abi.Pack("execTransaction",
		tx.TxData.To.Value,
		decimal.RequireFromString(tx.TxData.Value).IntPart(),
		common.Hex2Bytes(strings.TrimPrefix(data, "0x")),
		uint8(tx.TxData.Operation),
		cast.ToInt64(tx.DetailedExecutionInfo.SafeTxGas),
		cast.ToInt64(tx.DetailedExecutionInfo.BaseGas),
		cast.ToInt64(tx.DetailedExecutionInfo.GasPrice),
		tx.DetailedExecutionInfo.GasToken,
		tx.DetailedExecutionInfo.RefundReceiver.Value,
		common.Hex2Bytes(strings.Join(signatures, "")),
	)
	if err != nil {
		return errors.Wrap(err, "pack error")
	}
	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return errors.Wrap(err, "suggest gas price error")
	}
	to := common.HexToAddress(tx.TxData.To.Value)
	gas, err := client.EstimateGas(context.TODO(), ethereum.CallMsg{
		From:     s.GetAddress(true),
		To:       &to,
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
		To:       &to,
		Data:     input,
	})

	signTx, err := chains.SignTx(chainTransaction, s.conf.PrivateKey, int64(s.GetChainIdByCtx(ctx)))
	if err != nil {
		return errors.Wrap(err, "sign tx error")
	}
	return client.SendTransaction(ctx, signTx)
}

func (s *MantleSafe) eip712(ctx context.Context, t safe.Transaction) apitypes.TypedData {
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
			"operation":      big.NewInt(int64(t.Operation)),
			"baseGas":        t.BaseGas.BigInt(),
			"gasPrice":       t.GasPrice.BigInt(),
			"gasToken":       t.GasToken.Hex(),
			"refundReceiver": t.RefundReceiver.Hex(),
			"nonce":          big.NewInt(int64(t.Nonce)),
			"safeTxGas":      t.SafeTxGas.BigInt(),
		},
	}
}
