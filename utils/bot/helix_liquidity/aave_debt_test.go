package helix_liquidity

import (
	"context"
	"omni-balance/utils/chains"
	"omni-balance/utils/chains/chain_mocks"
	"omni-balance/utils/constant"
	"omni-balance/utils/erc20"
	"strconv"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAave_BalanceOf(t *testing.T) {
	mockClient := chain_mocks.NewMockClient(t)

	erc20Abi, err := erc20.TokenMetaData.GetAbi()
	assert.NoError(t, err)
	input, err := erc20Abi.Pack("balanceOf", constant.ZeroAddress)
	assert.NoError(t, err)
	token := common.HexToAddress("0x460b97BD498E1157530AEb3086301d5225b91216")
	mockClient.On("CallContract", context.TODO(), ethereum.CallMsg{To: &token, Data: input}, mock.Anything).Return(
		chains.EthToWei(decimal.RequireFromString("1000"), 6).Bytes(), nil,
	)

	vtoken := common.HexToAddress("0x4fBE3A94C60A5085dA6a2D309965DcF34c36711d")
	mockClient.On("CallContract", context.TODO(), ethereum.CallMsg{To: &vtoken, Data: input}, mock.Anything).Return(
		chains.EthToWei(decimal.RequireFromString("900"), 6).Bytes(), nil,
	)

	balance, err := new(Aave).BalanceOf(context.TODO(), DebtParams{
		Address: constant.ZeroAddress,
		Token:   "USDC",
		Client:  mockClient,
		Chain:   constant.ArbitrumSepolia,
	})
	assert.NoError(t, err)
	assert.Equal(t, balance.String(), strconv.Itoa(1000-900))
}
