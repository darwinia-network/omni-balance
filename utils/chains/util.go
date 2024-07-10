package chains

import (
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient/simulated"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"math/big"
	"omni-balance/utils"
	"omni-balance/utils/constant"
	"omni-balance/utils/erc20"
	"omni-balance/utils/error_types"
	"strings"
	"time"
)

func getDecimals(decimals ...int32) int32 {
	var currentDecimals int32 = 18
	if len(decimals) != 0 && decimals[0] != 0 {
		currentDecimals = decimals[0]
	}
	return currentDecimals
}

func EthToWei(v decimal.Decimal, decimals ...int32) *big.Int {
	decimalStr := v.Mul(decimal.New(1, getDecimals(decimals...)))
	wei, _ := new(big.Int).SetString(decimalStr.Truncate(0).String(), 10)
	return wei
}

func WeiToEth(v *big.Int, decimals ...int32) decimal.Decimal {
	return decimal.NewFromBigInt(v, 0).Div(decimal.New(1, getDecimals(decimals...)))
}

// GetTokenBalance retrieves the token balance of a specified wallet address.
// ctx: Context to control the function's execution, including cancellation and timeouts.
// tokenAddress: Contract address of the token; if "0x0000000000000000000000000000000000000000", it indicates querying the native token (e.g., ETH) balance.
// walletAddress: Address of the wallet whose token balance is being queried.
// tokenDecimal: Optional argument specifying the number of decimal places for the token. Required and must not be zero when querying the native token balance.
// Returns the token balance in its smallest unit and an error if any occurs.
func GetTokenBalance(ctx context.Context, client simulated.Client, tokenAddress, walletAddress string,
	tokenDecimal ...interface{}) (decimal.Decimal, error) {

	var (
		balance *big.Int
		err     error
	)
	if strings.EqualFold(tokenAddress, constant.ZeroAddress.Hex()) {
		if len(tokenDecimal) == 0 && tokenDecimal[0] != 0 {
			return decimal.Zero, errors.Errorf("token decimal is required")
		}
		tokenDecimals := cast.ToInt32(tokenDecimal[0])
		balance, err = client.BalanceAt(ctx, common.HexToAddress(walletAddress), nil)
		if err != nil {
			return decimal.Zero, errors.Wrap(err, "get balance error")
		}
		return WeiToEth(balance, tokenDecimals), nil
	}
	erc20Abi, err := erc20.TokenMetaData.GetAbi()
	if err != nil {
		return decimal.Zero, errors.Wrap(err, "get erc20 abi error")
	}
	input, err := erc20Abi.Pack("balanceOf", common.HexToAddress(walletAddress))
	if err != nil {
		return decimal.Zero, errors.Wrap(err, "pack erc20 balanceOf error")
	}
	to := common.HexToAddress(tokenAddress)
	result, err := client.CallContract(ctx, ethereum.CallMsg{To: &to, Data: input}, nil)
	if err != nil {
		return decimal.Zero, errors.Wrap(err, "call erc20 balanceOf error")
	}

	balance = new(big.Int).SetBytes(result)
	if len(tokenDecimal) == 0 {
		input, err := erc20Abi.Pack("decimals")
		if err != nil {
			return decimal.Zero, errors.Wrap(err, "pack erc20 decimals error")
		}
		tokenDecimals, err := client.CallContract(ctx, ethereum.CallMsg{To: &to, Data: input}, nil)
		if err != nil {
			return decimal.Zero, errors.Wrap(err, "call erc20 decimals error")
		}
		tokenDecimal = append(tokenDecimal, new(big.Int).SetBytes(tokenDecimals).Int64())
	}
	return WeiToEth(balance, cast.ToInt32(tokenDecimal[0])), nil
}

func SignTx(tx *types.Transaction, privateKey string, chainId int64) (*types.Transaction, error) {
	key, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return nil, errors.Wrap(err, "privateKey")
	}
	return types.SignTx(tx, types.NewEIP155Signer(big.NewInt(chainId)), key)
}

func SignMsg(msg []byte, privateKey string) (sig []byte, err error) {
	key, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return crypto.Sign(msg, key)
}

func WaitForTx(ctx context.Context, client simulated.Client, txHash common.Hash) error {
	log := utils.GetLogFromCtx(ctx)
	t := time.NewTicker(time.Second)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return errors.New("context cancelled")
		case <-t.C:
			tx, err := client.TransactionReceipt(ctx, txHash)
			if errors.Is(err, context.Canceled) ||
				errors.Is(err, context.DeadlineExceeded) {
				return err
			}
			if errors.Is(err, ethereum.NotFound) {
				log.Infof("tx not found, txHash: %s, try again later", txHash.Hex())
				continue
			}
			if err != nil {
				log.Debugf("get tx receipt error: %s, try again later", err)
				continue
			}
			if tx.Status == 0 {
				return errors.New("tx failed")
			}
			log.Infof("tx success, txHash: %s", txHash.Hex())
			return nil
		}
	}
}

type TokenApproveParams struct {
	ChainId         int64
	TokenAddress    common.Address
	Owner           common.Address
	SendTransaction func(ctx context.Context, tx *types.LegacyTx, client simulated.Client) (common.Hash, error)
	WaitTransaction func(ctx context.Context, txHash common.Hash, client simulated.Client) error
	Spender         common.Address
	AmountWei       decimal.Decimal
	Client          simulated.Client
}

func TokenApprove(ctx context.Context, args TokenApproveParams) error {
	log := utils.GetLogFromCtx(ctx).WithFields(logrus.Fields{
		"job":          "approve",
		"tokenAddress": args.TokenAddress.Hex(),
		"owner":        args.Owner.Hex(),
		"spender":      args.Spender.Hex(),
		"amountWei":    args.AmountWei.String(),
	})
	contract, err := erc20.NewTokenCaller(args.TokenAddress, args.Client)
	if err != nil {
		return errors.Wrap(err, "erc20 caller")
	}
	allowanceWei, err := contract.Allowance(nil, args.Owner, args.Spender)
	if err != nil {
		return errors.Wrap(err, "erc20 allowance")
	}
	log.Debugf("erc20 allowance: %s", allowanceWei)
	if decimal.NewFromBigInt(allowanceWei, 0).GreaterThanOrEqual(args.AmountWei) {
		log.Debugf("erc20 allowance %s >= %s, skip approve", allowanceWei, args.AmountWei)
		return nil
	}
	erc20Abi, err := erc20.TokenMetaData.GetAbi()
	if err != nil {
		return errors.Wrap(err, "get abi")
	}
	input, err := erc20Abi.Pack("approve", args.Spender, args.AmountWei.BigInt())
	if err != nil {
		return errors.Wrap(err, "abi pack")
	}

	txHash, err := args.SendTransaction(ctx, &types.LegacyTx{
		To:   &args.TokenAddress,
		Data: input,
	}, args.Client)
	if err != nil {
		return errors.Wrap(err, "erc20 approve")
	}
	log.Debugf("erc20 approve txHash: %s", txHash.Hex())
	if args.WaitTransaction != nil {
		return args.WaitTransaction(ctx, txHash, args.Client)
	}
	return WaitForTx(ctx, args.Client, txHash)
}

func MostNewAbiType(t string, internalType string, components []abi.ArgumentMarshaling) abi.Type {
	newType, err := abi.NewType(t, internalType, components)
	if err != nil {
		panic(err)
	}
	return newType
}

type SendTokenParams struct {
	Client        simulated.Client
	Sender        common.Address
	TokenAddress  common.Address
	TokenDecimals int32
	ToAddress     common.Address
	AmountWei     decimal.Decimal
}

func BuildSendToken(ctx context.Context, args SendTokenParams) (*types.LegacyTx, error) {
	var (
		isNativeToken = strings.EqualFold(args.TokenAddress.Hex(), constant.ZeroAddress.Hex())
	)
	balance, err := GetTokenBalance(ctx, args.Client, args.TokenAddress.Hex(), args.Sender.Hex(), args.TokenDecimals)
	if err != nil {
		return nil, errors.Wrap(err, "get balance")
	}
	balanceWei := decimal.NewFromBigInt(EthToWei(balance, args.TokenDecimals), 0)
	if balanceWei.IsZero() || balanceWei.IsNegative() || balanceWei.LessThan(args.AmountWei) {
		return nil, error_types.ErrInsufficientBalance
	}

	if !isNativeToken {
		nativeTokenBalanceWei, err := args.Client.BalanceAt(ctx, args.Sender, nil)
		if err != nil {
			return nil, error_types.ErrNativeTokenInsufficient
		}
		if nativeTokenBalanceWei.Cmp(big.NewInt(0)) < 0 {
			return nil, error_types.ErrNativeTokenInsufficient
		}
	}

	var tx *types.LegacyTx
	if isNativeToken {
		tx = &types.LegacyTx{
			To:    &args.ToAddress,
			Value: args.AmountWei.BigInt(),
		}
	}

	if !isNativeToken {
		erc20Abi, err := erc20.TokenMetaData.GetAbi()
		if err != nil {
			return nil, errors.Wrap(err, "get abi")
		}
		input, err := erc20Abi.Pack("transfer", args.ToAddress, args.AmountWei.BigInt())
		if err != nil {
			return nil, errors.Wrap(err, "abi pack")
		}
		tx = &types.LegacyTx{
			To:   &args.TokenAddress,
			Data: input,
		}
	}
	return tx, nil
}

func SendTransaction(ctx context.Context, client simulated.Client, tx *types.LegacyTx,
	sender common.Address, privateKey string) (common.Hash, error) {
	if tx.Nonce == 0 {
		nonce, err := client.NonceAt(ctx, sender, nil)
		if err != nil {
			return common.Hash{}, errors.Wrap(err, "get nonce")
		}
		tx.Nonce = nonce
	}
	if tx.GasPrice == nil {
		gasPrice, err := client.SuggestGasPrice(ctx)
		if err != nil {
			return common.Hash{}, errors.Wrap(err, "suggest gas price")
		}
		tx.GasPrice = gasPrice
	}

	if tx.Gas == 0 {
		gas, err := client.EstimateGas(ctx, ethereum.CallMsg{
			From:  sender,
			To:    tx.To,
			Value: tx.Value,
			Data:  tx.Data,
		})
		if err != nil {
			return common.Hash{}, errors.Wrap(err, "estimate gas")
		}
		tx.Gas = gas
	}
	chainId, err := client.ChainID(ctx)
	if err != nil {
		return common.Hash{}, errors.Wrap(err, "get chain id")
	}
	transaction, err := SignTx(types.NewTx(tx), privateKey, chainId.Int64())
	if err != nil {
		return common.Hash{}, errors.Wrap(err, "sign tx")
	}
	return transaction.Hash(), client.SendTransaction(ctx, transaction)
}

// GetTransferInfo input ex: a9059cbb0000000000000000000000000350101f2cb6aa65caab7954246a56f906a3f57d0000000000000000000000000000000000000000033b2e3c9fd0803ce8000000
func GetTransferInfo(input []byte) (to common.Address, amount *big.Int, err error) {
	if len(input) < 4 {
		return common.Address{}, nil, errors.New("invalid input")
	}
	erc20Abi, err := erc20.TokenMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, errors.Wrap(err, "get abi")
	}
	args, err := erc20Abi.Methods["transfer"].Inputs.Unpack(input[4:])
	if err != nil {
		return common.Address{}, nil, errors.Wrap(err, "unpack")
	}
	if len(args) != 2 {
		return common.Address{}, nil, errors.New("invalid args")
	}
	return args[0].(common.Address), args[1].(*big.Int), nil
}
