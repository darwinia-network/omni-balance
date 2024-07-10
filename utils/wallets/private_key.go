package wallets

import (
	"context"
	"encoding/json"
	"log"
	"omni-balance/utils/chains"
	"omni-balance/utils/constant"
	"omni-balance/utils/erc20"
	"omni-balance/utils/error_types"
	"omni-balance/utils/locks"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient/simulated"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

func init() {
	Register("", func(conf WalletConfig) Wallets {
		return NewPrivateKeyWallet(conf)
	})
}

type PrivateKeyWallet struct {
	conf WalletConfig
}

func NewPrivateKeyWallet(conf WalletConfig) *PrivateKeyWallet {
	p := &PrivateKeyWallet{
		conf: conf,
	}

	if p.conf.PrivateKey == "" &&
		conf.Operator.Address.Cmp(constant.ZeroAddress) == 0 &&
		conf.Operator.PrivateKey == "" {
		log.Fatalln("privateKey and operator.privateKey can not be empty")
	}
	if p.conf.Operator.Address.Cmp(constant.ZeroAddress) == 0 && p.conf.PrivateKey == "" {
		log.Fatalln("privateKey can not be empty when operator.address is not empty")
	}

	if p.conf.Operator.Operator.Cmp(constant.ZeroAddress) != 0 && p.conf.Operator.PrivateKey == "" {
		log.Fatalln("operator.privateKey can not be empty when operator.address is not empty")
	}
	checkKey := func(privateKey string, address common.Address) error {
		key, err := crypto.HexToECDSA(privateKey)
		if err != nil {
			return errors.Wrap(err, "crypto.HexToECDSA")
		}
		if !strings.EqualFold(address.Hex(), crypto.PubkeyToAddress(key.PublicKey).Hex()) {
			return error_types.ErrPrivateKeyNotMatch
		}
		return nil
	}

	if conf.Operator.Address.Cmp(constant.ZeroAddress) == 0 && conf.PrivateKey != "" && checkKey(conf.PrivateKey, conf.Address) != nil {
		log.Fatalln("privateKey not match address")
	}

	if conf.Operator.Address.Cmp(constant.ZeroAddress) != 0 &&
		conf.Operator.Operator.Cmp(constant.ZeroAddress) == 0 &&
		conf.PrivateKey != "" && checkKey(conf.PrivateKey, conf.Operator.Address) != nil {
		log.Fatalln("privateKey not match operator address")
	}

	if conf.Operator.Address.Cmp(constant.ZeroAddress) != 0 &&
		conf.Operator.Operator.Cmp(constant.ZeroAddress) != 0 &&
		conf.PrivateKey != "" && checkKey(conf.PrivateKey, conf.Operator.Operator) != nil {
		log.Fatalln("privateKey not match operator address")
	}
	return p
}

func (p *PrivateKeyWallet) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"address": p.GetAddress(),
		"operator": map[string]interface{}{
			"address":  p.conf.Operator.Address,
			"operator": p.conf.Operator.Operator,
		},
	})
}

func (p *PrivateKeyWallet) GetAddress(isReal ...bool) common.Address {
	if len(isReal) == 0 ||
		!isReal[0] ||
		p.conf.Operator.Address.Cmp(constant.ZeroAddress) == 0 {
		return p.conf.Address
	}
	return p.conf.Operator.Address
}

func (p *PrivateKeyWallet) IsSupportEip712() bool {
	return true
}

func (p *PrivateKeyWallet) Init(conf WalletConfig) error {
	p.conf = conf

	return nil
}

func (p *PrivateKeyWallet) IsDifferentAddress() bool {
	return p.GetAddress().Cmp(p.GetAddress(true)) != 0
}

func (p *PrivateKeyWallet) CheckFullAccess(_ context.Context) error {
	return nil
}

func (p *PrivateKeyWallet) GetExternalBalance(ctx context.Context, tokenAddress common.Address, decimals int32,
	client simulated.Client) (decimal.Decimal, error) {
	return p.balanceOf(ctx, p.GetAddress(), tokenAddress, decimals, client)
}

func (p *PrivateKeyWallet) GetBalance(ctx context.Context, tokenAddress common.Address, decimals int32,
	client simulated.Client) (decimal.Decimal, error) {
	return p.balanceOf(ctx, p.GetAddress(true), tokenAddress, decimals, client)
}

func (p *PrivateKeyWallet) balanceOf(ctx context.Context, wallet, tokenAddress common.Address, decimals int32,
	client simulated.Client) (decimal.Decimal, error) {
	if tokenAddress.Cmp(constant.ZeroAddress) == 0 {
		balance, err := client.BalanceAt(ctx, wallet, nil)
		if err != nil {
			return decimal.Zero, errors.Wrap(err, "client.BalanceAt")
		}
		return chains.WeiToEth(balance, decimals), nil
	}
	token, err := erc20.NewTokenCaller(tokenAddress, client)
	if err != nil {
		return decimal.Zero, errors.Wrap(err, "erc20.NewToken")
	}
	balance, err := token.BalanceOf(&bind.CallOpts{Context: ctx}, wallet)
	if err != nil {
		return decimal.Zero, errors.Wrap(err, "token.BalanceOf")
	}
	return chains.WeiToEth(balance, decimals), nil
}

func (p *PrivateKeyWallet) GetNonce(ctx context.Context, client simulated.Client) (uint64, error) {
	return client.NonceAt(ctx, p.GetAddress(true), nil)
}

func (p *PrivateKeyWallet) getPrivateKey() string {
	if p.conf.Operator.Address.Cmp(constant.ZeroAddress) != 0 {
		return p.conf.Operator.PrivateKey
	}
	return p.conf.PrivateKey
}

func (p *PrivateKeyWallet) SendTransaction(ctx context.Context, tx *types.LegacyTx,
	client simulated.Client) (common.Hash, error) {
	if tx.GasPrice == nil {
		gasPrice, err := client.SuggestGasPrice(ctx)
		if err != nil {
			return common.Hash{}, errors.Wrap(err, "get suggest gas")
		}
		tx.GasPrice = gasPrice
	}
	if tx.Gas == 0 {
		gas, err := client.EstimateGas(ctx, ethereum.CallMsg{
			From:  p.GetAddress(true),
			To:    tx.To,
			Value: tx.Value,
			Data:  tx.Data,
		})
		if err != nil {
			return common.Hash{}, errors.Wrap(err, "get estimate gas")
		}
		tx.Gas = gas
	}
	if tx.Nonce == 0 {
		nonce, err := p.GetNonce(ctx, client)
		if err != nil {
			return common.Hash{}, errors.Wrap(err, "get nonce")
		}
		tx.Nonce = nonce
	}

	chainId, err := client.ChainID(ctx)
	if err != nil {
		return common.Hash{}, errors.Wrap(err, "get chain id")
	}

	txData := types.NewTx(tx)
	txData, err = chains.SignTx(txData, p.getPrivateKey(), chainId.Int64())
	if err != nil {
		return common.Hash{}, errors.Wrap(err, "sign tx")
	}
	key := locks.LockKey(p.GetAddress(true).Hex(), chainId.String())
	locks.LockWithKey(ctx, key)
	defer locks.UnlockWithKey(ctx, key)
	return txData.Hash(), client.SendTransaction(ctx, txData)
}

func (p *PrivateKeyWallet) WaitTransaction(ctx context.Context, txHash common.Hash, client simulated.Client) error {
	return chains.WaitForTx(ctx, client, txHash)
}

func (p *PrivateKeyWallet) SignRawMessage(msg []byte) (sig []byte, err error) {
	return chains.SignMsg(msg, p.conf.PrivateKey)
}
func (p *PrivateKeyWallet) GetRealHash(_ context.Context, txHash common.Hash, _ simulated.Client) (common.Hash, error) {
	return txHash, nil
}
