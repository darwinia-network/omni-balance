package chains

import (
	"context"
	"math/big"
	"omni-balance/utils/chains/chain_mocks"
	"omni-balance/utils/constant"
	"omni-balance/utils/erc20"
	"omni-balance/utils/error_types"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestBuildSendToken tests the BuildSendToken function
func TestBuildSendToken(t *testing.T) {

	var (
		mockClient = new(chain_mocks.MockClient)
		ctx        = context.Background()
	)

	mockClient.On("BalanceAt", ctx, constant.ZeroAddress, mock.Anything).Return(EthToWei(decimal.RequireFromString("1000"), 18), nil)
	tx, err := BuildSendToken(ctx, SendTokenParams{
		Client:        mockClient,
		Sender:        constant.ZeroAddress,
		TokenAddress:  constant.ZeroAddress,
		TokenDecimals: 18,
		ToAddress:     constant.ZeroAddress,
		AmountWei:     decimal.NewFromBigInt(EthToWei(decimal.RequireFromString("1000"), 18), 0),
	})
	assert.NoError(t, err)
	assert.NotNil(t, tx)
	assert.Equal(t, constant.ZeroAddress.Hex(), tx.To.Hex())
	assert.Equal(t, EthToWei(decimal.RequireFromString("1000"), 18).String(), tx.Value.String())
	assert.Equal(t, "", common.Bytes2Hex(tx.Data))

	mockClient = new(chain_mocks.MockClient)
	mockClient.On("BalanceAt", ctx, constant.ZeroAddress, mock.Anything).Return(EthToWei(decimal.RequireFromString("0"), 18), nil)
	tx, err = BuildSendToken(ctx, SendTokenParams{
		Client:        mockClient,
		Sender:        constant.ZeroAddress,
		TokenAddress:  constant.ZeroAddress,
		TokenDecimals: 18,
		ToAddress:     constant.ZeroAddress,
		AmountWei:     decimal.NewFromBigInt(EthToWei(decimal.RequireFromString("1000"), 18), 0),
	})
	assert.Error(t, err)
	assert.Equal(t, err.Error(), error_types.ErrInsufficientBalance.Error())
	assert.Nil(t, tx)

	// erc20 token
	mockClient = new(chain_mocks.MockClient)
	mockClient.On("BalanceAt", ctx, constant.ZeroAddress, mock.Anything).Return(EthToWei(decimal.RequireFromString("1000"), 6), nil)
	erc20Abi, err := erc20.TokenMetaData.GetAbi()
	assert.NoError(t, err)
	input, err := erc20Abi.Pack("balanceOf", constant.ZeroAddress)
	assert.NoError(t, err)
	token := common.HexToAddress("0x75faf114eafb1BDbe2F0316DF893fd58CE46AA4d")
	mockClient.On("CallContract", ctx, ethereum.CallMsg{To: &token, Data: input}, mock.Anything).Return(big.NewInt(1000000000000).Bytes(), nil)

	tx, err = BuildSendToken(ctx, SendTokenParams{
		Client:        mockClient,
		Sender:        constant.ZeroAddress,
		TokenAddress:  token,
		TokenDecimals: 18,
		ToAddress:     constant.ZeroAddress,
		AmountWei:     decimal.NewFromBigInt(EthToWei(decimal.RequireFromString("1000"), 6), 0),
	})
	assert.NoError(t, err)
	assert.NotNil(t, tx)
	assert.Equal(t, token.Hex(), tx.To.Hex())
	assert.Nil(t, tx.Value)
	assert.Equal(t, "a9059cbb0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000003b9aca00", common.Bytes2Hex(tx.Data))
}
