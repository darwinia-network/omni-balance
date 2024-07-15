package chains

import (
	"math/big"
	"omni-balance/utils/constant"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"github.com/stretchr/testify/assert"
)

func TestEncodeForSigning(t *testing.T) {
	// Define a mock TypedData for testing
	typedData := apitypes.TypedData{
		Types: apitypes.Types{
			"EIP712Domain": []apitypes.Type{
				{Name: "name", Type: "string"},
				{Name: "chainId", Type: "uint256"},
				{Name: "verifyingContract", Type: "address"},
			},
			"PermitDetails": []apitypes.Type{
				{Name: "token", Type: "address"},
				{Name: "amount", Type: "uint160"},
				{Name: "expiration", Type: "uint48"},
				{Name: "nonce", Type: "uint48"},
			},
			"PermitSingle": []apitypes.Type{
				{Name: "details", Type: "PermitDetails"},
				{Name: "spender", Type: "address"},
				{Name: "sigDeadline", Type: "uint256"},
			},
		},
		Domain: apitypes.TypedDataDomain{
			Name:              "PermitSingle",
			ChainId:           math.NewHexOrDecimal256(1),
			VerifyingContract: constant.ZeroAddress.Hex(),
		},
		PrimaryType: "PermitSingle",
		Message: map[string]interface{}{
			"details": map[string]interface{}{
				"token":      constant.ZeroAddress.Hex(),
				"amount":     big.NewInt(0),
				"expiration": big.NewInt(1721024951),
				"nonce":      big.NewInt(1),
			},
			"spender":     constant.ZeroAddress.Hex(),
			"sigDeadline": big.NewInt(1721024951),
		},
	}

	// Define the expected hash result
	expectedHash := "0x0197a3976d5059736a317a039e4328fbbbe328c1a330e644ea29fc7beddb7fdb"

	// Call the EncodeForSigning function
	hash, err := EncodeForSigning(typedData)

	// Assert that there is no error
	assert.NoError(t, err)

	// Assert that the returned hash matches the expected hash
	assert.Equal(t, expectedHash, hash.String())
}

func TestSignTypedData(t *testing.T) {
	// Define a mock TypedData for testing
	typedData := apitypes.TypedData{
		Types: apitypes.Types{
			"EIP712Domain": []apitypes.Type{
				{Name: "name", Type: "string"},
				{Name: "chainId", Type: "uint256"},
				{Name: "verifyingContract", Type: "address"},
			},
			"PermitDetails": []apitypes.Type{
				{Name: "token", Type: "address"},
				{Name: "amount", Type: "uint160"},
				{Name: "expiration", Type: "uint48"},
				{Name: "nonce", Type: "uint48"},
			},
			"PermitSingle": []apitypes.Type{
				{Name: "details", Type: "PermitDetails"},
				{Name: "spender", Type: "address"},
				{Name: "sigDeadline", Type: "uint256"},
			},
		},
		Domain: apitypes.TypedDataDomain{
			Name:              "PermitSingle",
			ChainId:           math.NewHexOrDecimal256(1),
			VerifyingContract: constant.ZeroAddress.Hex(),
		},
		PrimaryType: "PermitSingle",
		Message: map[string]interface{}{
			"details": map[string]interface{}{
				"token":      constant.ZeroAddress.Hex(),
				"amount":     big.NewInt(0),
				"expiration": big.NewInt(1721024951),
				"nonce":      big.NewInt(1),
			},
			"spender":     constant.ZeroAddress.Hex(),
			"sigDeadline": big.NewInt(1721024951),
		},
	}

	// Define the expected hash result
	expectedHash := "df31dfa39a6b758368cf9711484926ff65d40eeecb77d51c3c185b73a4239b4a59dc6f4a993b97938b87ffc2269dfc2a7931e2f59e8fe8bfd39138e84e8d13981b"

	// Call the EncodeForSigning function
	hash, err := SignTypedData(typedData, func(msg []byte) (sig []byte, err error) {
		return SignMsg(msg, constant.TestPrivateKey)
	})

	// Assert that there is no error
	assert.NoError(t, err)

	// Assert that the returned hash matches the expected hash
	assert.Equal(t, expectedHash, common.Bytes2Hex(hash))
}
