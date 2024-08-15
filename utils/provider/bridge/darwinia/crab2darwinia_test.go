package darwinia

import (
	"context"
	"omni-balance/utils"
	"omni-balance/utils/chains"
	"omni-balance/utils/constant"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestCrab2Darwinia(t *testing.T) {
	var (
		args   = InitVar(t)
		amount = decimal.NewFromBigInt(chains.EthToWei(decimal.RequireFromString("0")), 0)
	)
	args.ctx = context.WithValue(args.ctx, constant.FeeTestKeyInCtx, decimal.RequireFromString("1"))
	tx, err := Crab2Darwinia(args.ctx, SwapParams{
		Sender:    args.sender,
		TokenName: "CRAB",
		Amount:    amount,
		Nonce:     1716966768607,
		Client:    args.client,
	})
	assert.NoError(t, err)
	utils.AssertEqualFold(t, common.Bytes2Hex(tx.Data), "7104AAD5000000000000000000000000000000000000000000000000000000000000002E00000000000000000000000043EF13E84D9992D1461A1F90CAC4653658CEA4FD00000000000000000000000043EF13E84D9992D1461A1F90CAC4653658CEA4FD00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000018FC331BFDF00000000000000000000000000000000000000000000000000000000000000E0000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000080000000000000000000000000000000000000000000000000000000000005461400000000000000000000000043EF13E84D9992D1461A1F90CAC4653658CEA4FD00000000000000000000000000000000000000000000000000000000000000600000000000000000000000000000000000000000000000000000000000000000")

	assert.Equal(t, tx.To.Hex(), "0x004D0dE211BC148c3Ce696C51Cbc85BD421727E9")

	tx, err = Crab2Darwinia(args.ctx, SwapParams{
		Sender:    args.sender,
		TokenName: "RING",
		Amount:    amount,
		Nonce:     1716966714726,
		Client:    args.client,
	})
	assert.NoError(t, err)
	utils.AssertEqualFold(t, common.Bytes2Hex(tx.Data), "57BF0985000000000000000000000000273131F7CB50AC002BDD08CA721988731F7E1092000000000000000000000000A8D0E9A45249EC839C397FA0F371F5F64ECAB7F700000000000000000000000043EF13E84D9992D1461A1F90CAC4653658CEA4FD00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000018FC330ED6600000000000000000000000000000000000000000000000000000000000000E00000000000000000000000000000000000000000000000000000000000000120000000000000000000000000000000000000000000000000000000000000001443EF13E84D9992D1461A1F90CAC4653658CEA4FD000000000000000000000000000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000552C200000000000000000000000043EF13E84D9992D1461A1F90CAC4653658CEA4FD00000000000000000000000000000000000000000000000000000000000000600000000000000000000000000000000000000000000000000000000000000000")
	assert.Equal(t, tx.To.Hex(), "0xf6372ab2d35B32156A19F2d2F23FA6dDeFBE58bd")
}
