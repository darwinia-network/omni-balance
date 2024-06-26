// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package v3SwapRouter

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// ISwapRouterExactInputParams is an auto generated low-level Go binding around an user-defined struct.
type ISwapRouterExactInputParams struct {
	Path             []byte
	Recipient        common.Address
	Deadline         *big.Int
	AmountIn         *big.Int
	AmountOutMinimum *big.Int
}

// ISwapRouterExactInputSingleParams is an auto generated low-level Go binding around an user-defined struct.
type ISwapRouterExactInputSingleParams struct {
	TokenIn           common.Address
	TokenOut          common.Address
	Fee               *big.Int
	Recipient         common.Address
	Deadline          *big.Int
	AmountIn          *big.Int
	AmountOutMinimum  *big.Int
	SqrtPriceLimitX96 *big.Int
}

// ISwapRouterExactOutputParams is an auto generated low-level Go binding around an user-defined struct.
type ISwapRouterExactOutputParams struct {
	Path            []byte
	Recipient       common.Address
	Deadline        *big.Int
	AmountOut       *big.Int
	AmountInMaximum *big.Int
}

// ISwapRouterExactOutputSingleParams is an auto generated low-level Go binding around an user-defined struct.
type ISwapRouterExactOutputSingleParams struct {
	TokenIn           common.Address
	TokenOut          common.Address
	Fee               *big.Int
	Recipient         common.Address
	Deadline          *big.Int
	AmountOut         *big.Int
	AmountInMaximum   *big.Int
	SqrtPriceLimitX96 *big.Int
}

// V3SwapRouterMetaData contains all meta data concerning the V3SwapRouter contract.
var V3SwapRouterMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_factory\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_WETH9\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"WETH9\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"path\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountOutMinimum\",\"type\":\"uint256\"}],\"internalType\":\"structISwapRouter.ExactInputParams\",\"name\":\"params\",\"type\":\"tuple\"}],\"name\":\"exactInput\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountOut\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"tokenIn\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenOut\",\"type\":\"address\"},{\"internalType\":\"uint24\",\"name\":\"fee\",\"type\":\"uint24\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountOutMinimum\",\"type\":\"uint256\"},{\"internalType\":\"uint160\",\"name\":\"sqrtPriceLimitX96\",\"type\":\"uint160\"}],\"internalType\":\"structISwapRouter.ExactInputSingleParams\",\"name\":\"params\",\"type\":\"tuple\"}],\"name\":\"exactInputSingle\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountOut\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"path\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountOut\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountInMaximum\",\"type\":\"uint256\"}],\"internalType\":\"structISwapRouter.ExactOutputParams\",\"name\":\"params\",\"type\":\"tuple\"}],\"name\":\"exactOutput\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"tokenIn\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenOut\",\"type\":\"address\"},{\"internalType\":\"uint24\",\"name\":\"fee\",\"type\":\"uint24\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountOut\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountInMaximum\",\"type\":\"uint256\"},{\"internalType\":\"uint160\",\"name\":\"sqrtPriceLimitX96\",\"type\":\"uint160\"}],\"internalType\":\"structISwapRouter.ExactOutputSingleParams\",\"name\":\"params\",\"type\":\"tuple\"}],\"name\":\"exactOutputSingle\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"factory\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes[]\",\"name\":\"data\",\"type\":\"bytes[]\"}],\"name\":\"multicall\",\"outputs\":[{\"internalType\":\"bytes[]\",\"name\":\"results\",\"type\":\"bytes[]\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"refundETH\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"selfPermit\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expiry\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"selfPermitAllowed\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expiry\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"selfPermitAllowedIfNecessary\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"selfPermitIfNecessary\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amountMinimum\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"}],\"name\":\"sweepToken\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amountMinimum\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"feeBips\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"feeRecipient\",\"type\":\"address\"}],\"name\":\"sweepTokenWithFee\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"amount0Delta\",\"type\":\"int256\"},{\"internalType\":\"int256\",\"name\":\"amount1Delta\",\"type\":\"int256\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"uniswapV3SwapCallback\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amountMinimum\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"}],\"name\":\"unwrapWETH9\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amountMinimum\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"feeBips\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"feeRecipient\",\"type\":\"address\"}],\"name\":\"unwrapWETH9WithFee\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
}

// V3SwapRouterABI is the input ABI used to generate the binding from.
// Deprecated: Use V3SwapRouterMetaData.ABI instead.
var V3SwapRouterABI = V3SwapRouterMetaData.ABI

// V3SwapRouter is an auto generated Go binding around an Ethereum contract.
type V3SwapRouter struct {
	V3SwapRouterCaller     // Read-only binding to the contract
	V3SwapRouterTransactor // Write-only binding to the contract
	V3SwapRouterFilterer   // Log filterer for contract events
}

// V3SwapRouterCaller is an auto generated read-only Go binding around an Ethereum contract.
type V3SwapRouterCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// V3SwapRouterTransactor is an auto generated write-only Go binding around an Ethereum contract.
type V3SwapRouterTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// V3SwapRouterFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type V3SwapRouterFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// V3SwapRouterSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type V3SwapRouterSession struct {
	Contract     *V3SwapRouter     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// V3SwapRouterCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type V3SwapRouterCallerSession struct {
	Contract *V3SwapRouterCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// V3SwapRouterTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type V3SwapRouterTransactorSession struct {
	Contract     *V3SwapRouterTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// V3SwapRouterRaw is an auto generated low-level Go binding around an Ethereum contract.
type V3SwapRouterRaw struct {
	Contract *V3SwapRouter // Generic contract binding to access the raw methods on
}

// V3SwapRouterCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type V3SwapRouterCallerRaw struct {
	Contract *V3SwapRouterCaller // Generic read-only contract binding to access the raw methods on
}

// V3SwapRouterTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type V3SwapRouterTransactorRaw struct {
	Contract *V3SwapRouterTransactor // Generic write-only contract binding to access the raw methods on
}

// NewV3SwapRouter creates a new instance of V3SwapRouter, bound to a specific deployed contract.
func NewV3SwapRouter(address common.Address, backend bind.ContractBackend) (*V3SwapRouter, error) {
	contract, err := bindV3SwapRouter(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &V3SwapRouter{V3SwapRouterCaller: V3SwapRouterCaller{contract: contract}, V3SwapRouterTransactor: V3SwapRouterTransactor{contract: contract}, V3SwapRouterFilterer: V3SwapRouterFilterer{contract: contract}}, nil
}

// NewV3SwapRouterCaller creates a new read-only instance of V3SwapRouter, bound to a specific deployed contract.
func NewV3SwapRouterCaller(address common.Address, caller bind.ContractCaller) (*V3SwapRouterCaller, error) {
	contract, err := bindV3SwapRouter(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &V3SwapRouterCaller{contract: contract}, nil
}

// NewV3SwapRouterTransactor creates a new write-only instance of V3SwapRouter, bound to a specific deployed contract.
func NewV3SwapRouterTransactor(address common.Address, transactor bind.ContractTransactor) (*V3SwapRouterTransactor, error) {
	contract, err := bindV3SwapRouter(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &V3SwapRouterTransactor{contract: contract}, nil
}

// NewV3SwapRouterFilterer creates a new log filterer instance of V3SwapRouter, bound to a specific deployed contract.
func NewV3SwapRouterFilterer(address common.Address, filterer bind.ContractFilterer) (*V3SwapRouterFilterer, error) {
	contract, err := bindV3SwapRouter(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &V3SwapRouterFilterer{contract: contract}, nil
}

// bindV3SwapRouter binds a generic wrapper to an already deployed contract.
func bindV3SwapRouter(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := V3SwapRouterMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_V3SwapRouter *V3SwapRouterRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _V3SwapRouter.Contract.V3SwapRouterCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_V3SwapRouter *V3SwapRouterRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _V3SwapRouter.Contract.V3SwapRouterTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_V3SwapRouter *V3SwapRouterRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _V3SwapRouter.Contract.V3SwapRouterTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_V3SwapRouter *V3SwapRouterCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _V3SwapRouter.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_V3SwapRouter *V3SwapRouterTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _V3SwapRouter.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_V3SwapRouter *V3SwapRouterTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _V3SwapRouter.Contract.contract.Transact(opts, method, params...)
}

// WETH9 is a free data retrieval call binding the contract method 0x4aa4a4fc.
//
// Solidity: function WETH9() view returns(address)
func (_V3SwapRouter *V3SwapRouterCaller) WETH9(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _V3SwapRouter.contract.Call(opts, &out, "WETH9")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// WETH9 is a free data retrieval call binding the contract method 0x4aa4a4fc.
//
// Solidity: function WETH9() view returns(address)
func (_V3SwapRouter *V3SwapRouterSession) WETH9() (common.Address, error) {
	return _V3SwapRouter.Contract.WETH9(&_V3SwapRouter.CallOpts)
}

// WETH9 is a free data retrieval call binding the contract method 0x4aa4a4fc.
//
// Solidity: function WETH9() view returns(address)
func (_V3SwapRouter *V3SwapRouterCallerSession) WETH9() (common.Address, error) {
	return _V3SwapRouter.Contract.WETH9(&_V3SwapRouter.CallOpts)
}

// Factory is a free data retrieval call binding the contract method 0xc45a0155.
//
// Solidity: function factory() view returns(address)
func (_V3SwapRouter *V3SwapRouterCaller) Factory(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _V3SwapRouter.contract.Call(opts, &out, "factory")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Factory is a free data retrieval call binding the contract method 0xc45a0155.
//
// Solidity: function factory() view returns(address)
func (_V3SwapRouter *V3SwapRouterSession) Factory() (common.Address, error) {
	return _V3SwapRouter.Contract.Factory(&_V3SwapRouter.CallOpts)
}

// Factory is a free data retrieval call binding the contract method 0xc45a0155.
//
// Solidity: function factory() view returns(address)
func (_V3SwapRouter *V3SwapRouterCallerSession) Factory() (common.Address, error) {
	return _V3SwapRouter.Contract.Factory(&_V3SwapRouter.CallOpts)
}

// ExactInput is a paid mutator transaction binding the contract method 0xc04b8d59.
//
// Solidity: function exactInput((bytes,address,uint256,uint256,uint256) params) payable returns(uint256 amountOut)
func (_V3SwapRouter *V3SwapRouterTransactor) ExactInput(opts *bind.TransactOpts, params ISwapRouterExactInputParams) (*types.Transaction, error) {
	return _V3SwapRouter.contract.Transact(opts, "exactInput", params)
}

// ExactInput is a paid mutator transaction binding the contract method 0xc04b8d59.
//
// Solidity: function exactInput((bytes,address,uint256,uint256,uint256) params) payable returns(uint256 amountOut)
func (_V3SwapRouter *V3SwapRouterSession) ExactInput(params ISwapRouterExactInputParams) (*types.Transaction, error) {
	return _V3SwapRouter.Contract.ExactInput(&_V3SwapRouter.TransactOpts, params)
}

// ExactInput is a paid mutator transaction binding the contract method 0xc04b8d59.
//
// Solidity: function exactInput((bytes,address,uint256,uint256,uint256) params) payable returns(uint256 amountOut)
func (_V3SwapRouter *V3SwapRouterTransactorSession) ExactInput(params ISwapRouterExactInputParams) (*types.Transaction, error) {
	return _V3SwapRouter.Contract.ExactInput(&_V3SwapRouter.TransactOpts, params)
}

// ExactInputSingle is a paid mutator transaction binding the contract method 0x414bf389.
//
// Solidity: function exactInputSingle((address,address,uint24,address,uint256,uint256,uint256,uint160) params) payable returns(uint256 amountOut)
func (_V3SwapRouter *V3SwapRouterTransactor) ExactInputSingle(opts *bind.TransactOpts, params ISwapRouterExactInputSingleParams) (*types.Transaction, error) {
	return _V3SwapRouter.contract.Transact(opts, "exactInputSingle", params)
}

// ExactInputSingle is a paid mutator transaction binding the contract method 0x414bf389.
//
// Solidity: function exactInputSingle((address,address,uint24,address,uint256,uint256,uint256,uint160) params) payable returns(uint256 amountOut)
func (_V3SwapRouter *V3SwapRouterSession) ExactInputSingle(params ISwapRouterExactInputSingleParams) (*types.Transaction, error) {
	return _V3SwapRouter.Contract.ExactInputSingle(&_V3SwapRouter.TransactOpts, params)
}

// ExactInputSingle is a paid mutator transaction binding the contract method 0x414bf389.
//
// Solidity: function exactInputSingle((address,address,uint24,address,uint256,uint256,uint256,uint160) params) payable returns(uint256 amountOut)
func (_V3SwapRouter *V3SwapRouterTransactorSession) ExactInputSingle(params ISwapRouterExactInputSingleParams) (*types.Transaction, error) {
	return _V3SwapRouter.Contract.ExactInputSingle(&_V3SwapRouter.TransactOpts, params)
}

// ExactOutput is a paid mutator transaction binding the contract method 0xf28c0498.
//
// Solidity: function exactOutput((bytes,address,uint256,uint256,uint256) params) payable returns(uint256 amountIn)
func (_V3SwapRouter *V3SwapRouterTransactor) ExactOutput(opts *bind.TransactOpts, params ISwapRouterExactOutputParams) (*types.Transaction, error) {
	return _V3SwapRouter.contract.Transact(opts, "exactOutput", params)
}

// ExactOutput is a paid mutator transaction binding the contract method 0xf28c0498.
//
// Solidity: function exactOutput((bytes,address,uint256,uint256,uint256) params) payable returns(uint256 amountIn)
func (_V3SwapRouter *V3SwapRouterSession) ExactOutput(params ISwapRouterExactOutputParams) (*types.Transaction, error) {
	return _V3SwapRouter.Contract.ExactOutput(&_V3SwapRouter.TransactOpts, params)
}

// ExactOutput is a paid mutator transaction binding the contract method 0xf28c0498.
//
// Solidity: function exactOutput((bytes,address,uint256,uint256,uint256) params) payable returns(uint256 amountIn)
func (_V3SwapRouter *V3SwapRouterTransactorSession) ExactOutput(params ISwapRouterExactOutputParams) (*types.Transaction, error) {
	return _V3SwapRouter.Contract.ExactOutput(&_V3SwapRouter.TransactOpts, params)
}

// ExactOutputSingle is a paid mutator transaction binding the contract method 0xdb3e2198.
//
// Solidity: function exactOutputSingle((address,address,uint24,address,uint256,uint256,uint256,uint160) params) payable returns(uint256 amountIn)
func (_V3SwapRouter *V3SwapRouterTransactor) ExactOutputSingle(opts *bind.TransactOpts, params ISwapRouterExactOutputSingleParams) (*types.Transaction, error) {
	return _V3SwapRouter.contract.Transact(opts, "exactOutputSingle", params)
}

// ExactOutputSingle is a paid mutator transaction binding the contract method 0xdb3e2198.
//
// Solidity: function exactOutputSingle((address,address,uint24,address,uint256,uint256,uint256,uint160) params) payable returns(uint256 amountIn)
func (_V3SwapRouter *V3SwapRouterSession) ExactOutputSingle(params ISwapRouterExactOutputSingleParams) (*types.Transaction, error) {
	return _V3SwapRouter.Contract.ExactOutputSingle(&_V3SwapRouter.TransactOpts, params)
}

// ExactOutputSingle is a paid mutator transaction binding the contract method 0xdb3e2198.
//
// Solidity: function exactOutputSingle((address,address,uint24,address,uint256,uint256,uint256,uint160) params) payable returns(uint256 amountIn)
func (_V3SwapRouter *V3SwapRouterTransactorSession) ExactOutputSingle(params ISwapRouterExactOutputSingleParams) (*types.Transaction, error) {
	return _V3SwapRouter.Contract.ExactOutputSingle(&_V3SwapRouter.TransactOpts, params)
}

// Multicall is a paid mutator transaction binding the contract method 0xac9650d8.
//
// Solidity: function multicall(bytes[] data) payable returns(bytes[] results)
func (_V3SwapRouter *V3SwapRouterTransactor) Multicall(opts *bind.TransactOpts, data [][]byte) (*types.Transaction, error) {
	return _V3SwapRouter.contract.Transact(opts, "multicall", data)
}

// Multicall is a paid mutator transaction binding the contract method 0xac9650d8.
//
// Solidity: function multicall(bytes[] data) payable returns(bytes[] results)
func (_V3SwapRouter *V3SwapRouterSession) Multicall(data [][]byte) (*types.Transaction, error) {
	return _V3SwapRouter.Contract.Multicall(&_V3SwapRouter.TransactOpts, data)
}

// Multicall is a paid mutator transaction binding the contract method 0xac9650d8.
//
// Solidity: function multicall(bytes[] data) payable returns(bytes[] results)
func (_V3SwapRouter *V3SwapRouterTransactorSession) Multicall(data [][]byte) (*types.Transaction, error) {
	return _V3SwapRouter.Contract.Multicall(&_V3SwapRouter.TransactOpts, data)
}

// RefundETH is a paid mutator transaction binding the contract method 0x12210e8a.
//
// Solidity: function refundETH() payable returns()
func (_V3SwapRouter *V3SwapRouterTransactor) RefundETH(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _V3SwapRouter.contract.Transact(opts, "refundETH")
}

// RefundETH is a paid mutator transaction binding the contract method 0x12210e8a.
//
// Solidity: function refundETH() payable returns()
func (_V3SwapRouter *V3SwapRouterSession) RefundETH() (*types.Transaction, error) {
	return _V3SwapRouter.Contract.RefundETH(&_V3SwapRouter.TransactOpts)
}

// RefundETH is a paid mutator transaction binding the contract method 0x12210e8a.
//
// Solidity: function refundETH() payable returns()
func (_V3SwapRouter *V3SwapRouterTransactorSession) RefundETH() (*types.Transaction, error) {
	return _V3SwapRouter.Contract.RefundETH(&_V3SwapRouter.TransactOpts)
}

// SelfPermit is a paid mutator transaction binding the contract method 0xf3995c67.
//
// Solidity: function selfPermit(address token, uint256 value, uint256 deadline, uint8 v, bytes32 r, bytes32 s) payable returns()
func (_V3SwapRouter *V3SwapRouterTransactor) SelfPermit(opts *bind.TransactOpts, token common.Address, value *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _V3SwapRouter.contract.Transact(opts, "selfPermit", token, value, deadline, v, r, s)
}

// SelfPermit is a paid mutator transaction binding the contract method 0xf3995c67.
//
// Solidity: function selfPermit(address token, uint256 value, uint256 deadline, uint8 v, bytes32 r, bytes32 s) payable returns()
func (_V3SwapRouter *V3SwapRouterSession) SelfPermit(token common.Address, value *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _V3SwapRouter.Contract.SelfPermit(&_V3SwapRouter.TransactOpts, token, value, deadline, v, r, s)
}

// SelfPermit is a paid mutator transaction binding the contract method 0xf3995c67.
//
// Solidity: function selfPermit(address token, uint256 value, uint256 deadline, uint8 v, bytes32 r, bytes32 s) payable returns()
func (_V3SwapRouter *V3SwapRouterTransactorSession) SelfPermit(token common.Address, value *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _V3SwapRouter.Contract.SelfPermit(&_V3SwapRouter.TransactOpts, token, value, deadline, v, r, s)
}

// SelfPermitAllowed is a paid mutator transaction binding the contract method 0x4659a494.
//
// Solidity: function selfPermitAllowed(address token, uint256 nonce, uint256 expiry, uint8 v, bytes32 r, bytes32 s) payable returns()
func (_V3SwapRouter *V3SwapRouterTransactor) SelfPermitAllowed(opts *bind.TransactOpts, token common.Address, nonce *big.Int, expiry *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _V3SwapRouter.contract.Transact(opts, "selfPermitAllowed", token, nonce, expiry, v, r, s)
}

// SelfPermitAllowed is a paid mutator transaction binding the contract method 0x4659a494.
//
// Solidity: function selfPermitAllowed(address token, uint256 nonce, uint256 expiry, uint8 v, bytes32 r, bytes32 s) payable returns()
func (_V3SwapRouter *V3SwapRouterSession) SelfPermitAllowed(token common.Address, nonce *big.Int, expiry *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _V3SwapRouter.Contract.SelfPermitAllowed(&_V3SwapRouter.TransactOpts, token, nonce, expiry, v, r, s)
}

// SelfPermitAllowed is a paid mutator transaction binding the contract method 0x4659a494.
//
// Solidity: function selfPermitAllowed(address token, uint256 nonce, uint256 expiry, uint8 v, bytes32 r, bytes32 s) payable returns()
func (_V3SwapRouter *V3SwapRouterTransactorSession) SelfPermitAllowed(token common.Address, nonce *big.Int, expiry *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _V3SwapRouter.Contract.SelfPermitAllowed(&_V3SwapRouter.TransactOpts, token, nonce, expiry, v, r, s)
}

// SelfPermitAllowedIfNecessary is a paid mutator transaction binding the contract method 0xa4a78f0c.
//
// Solidity: function selfPermitAllowedIfNecessary(address token, uint256 nonce, uint256 expiry, uint8 v, bytes32 r, bytes32 s) payable returns()
func (_V3SwapRouter *V3SwapRouterTransactor) SelfPermitAllowedIfNecessary(opts *bind.TransactOpts, token common.Address, nonce *big.Int, expiry *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _V3SwapRouter.contract.Transact(opts, "selfPermitAllowedIfNecessary", token, nonce, expiry, v, r, s)
}

// SelfPermitAllowedIfNecessary is a paid mutator transaction binding the contract method 0xa4a78f0c.
//
// Solidity: function selfPermitAllowedIfNecessary(address token, uint256 nonce, uint256 expiry, uint8 v, bytes32 r, bytes32 s) payable returns()
func (_V3SwapRouter *V3SwapRouterSession) SelfPermitAllowedIfNecessary(token common.Address, nonce *big.Int, expiry *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _V3SwapRouter.Contract.SelfPermitAllowedIfNecessary(&_V3SwapRouter.TransactOpts, token, nonce, expiry, v, r, s)
}

// SelfPermitAllowedIfNecessary is a paid mutator transaction binding the contract method 0xa4a78f0c.
//
// Solidity: function selfPermitAllowedIfNecessary(address token, uint256 nonce, uint256 expiry, uint8 v, bytes32 r, bytes32 s) payable returns()
func (_V3SwapRouter *V3SwapRouterTransactorSession) SelfPermitAllowedIfNecessary(token common.Address, nonce *big.Int, expiry *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _V3SwapRouter.Contract.SelfPermitAllowedIfNecessary(&_V3SwapRouter.TransactOpts, token, nonce, expiry, v, r, s)
}

// SelfPermitIfNecessary is a paid mutator transaction binding the contract method 0xc2e3140a.
//
// Solidity: function selfPermitIfNecessary(address token, uint256 value, uint256 deadline, uint8 v, bytes32 r, bytes32 s) payable returns()
func (_V3SwapRouter *V3SwapRouterTransactor) SelfPermitIfNecessary(opts *bind.TransactOpts, token common.Address, value *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _V3SwapRouter.contract.Transact(opts, "selfPermitIfNecessary", token, value, deadline, v, r, s)
}

// SelfPermitIfNecessary is a paid mutator transaction binding the contract method 0xc2e3140a.
//
// Solidity: function selfPermitIfNecessary(address token, uint256 value, uint256 deadline, uint8 v, bytes32 r, bytes32 s) payable returns()
func (_V3SwapRouter *V3SwapRouterSession) SelfPermitIfNecessary(token common.Address, value *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _V3SwapRouter.Contract.SelfPermitIfNecessary(&_V3SwapRouter.TransactOpts, token, value, deadline, v, r, s)
}

// SelfPermitIfNecessary is a paid mutator transaction binding the contract method 0xc2e3140a.
//
// Solidity: function selfPermitIfNecessary(address token, uint256 value, uint256 deadline, uint8 v, bytes32 r, bytes32 s) payable returns()
func (_V3SwapRouter *V3SwapRouterTransactorSession) SelfPermitIfNecessary(token common.Address, value *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _V3SwapRouter.Contract.SelfPermitIfNecessary(&_V3SwapRouter.TransactOpts, token, value, deadline, v, r, s)
}

// SweepToken is a paid mutator transaction binding the contract method 0xdf2ab5bb.
//
// Solidity: function sweepToken(address token, uint256 amountMinimum, address recipient) payable returns()
func (_V3SwapRouter *V3SwapRouterTransactor) SweepToken(opts *bind.TransactOpts, token common.Address, amountMinimum *big.Int, recipient common.Address) (*types.Transaction, error) {
	return _V3SwapRouter.contract.Transact(opts, "sweepToken", token, amountMinimum, recipient)
}

// SweepToken is a paid mutator transaction binding the contract method 0xdf2ab5bb.
//
// Solidity: function sweepToken(address token, uint256 amountMinimum, address recipient) payable returns()
func (_V3SwapRouter *V3SwapRouterSession) SweepToken(token common.Address, amountMinimum *big.Int, recipient common.Address) (*types.Transaction, error) {
	return _V3SwapRouter.Contract.SweepToken(&_V3SwapRouter.TransactOpts, token, amountMinimum, recipient)
}

// SweepToken is a paid mutator transaction binding the contract method 0xdf2ab5bb.
//
// Solidity: function sweepToken(address token, uint256 amountMinimum, address recipient) payable returns()
func (_V3SwapRouter *V3SwapRouterTransactorSession) SweepToken(token common.Address, amountMinimum *big.Int, recipient common.Address) (*types.Transaction, error) {
	return _V3SwapRouter.Contract.SweepToken(&_V3SwapRouter.TransactOpts, token, amountMinimum, recipient)
}

// SweepTokenWithFee is a paid mutator transaction binding the contract method 0xe0e189a0.
//
// Solidity: function sweepTokenWithFee(address token, uint256 amountMinimum, address recipient, uint256 feeBips, address feeRecipient) payable returns()
func (_V3SwapRouter *V3SwapRouterTransactor) SweepTokenWithFee(opts *bind.TransactOpts, token common.Address, amountMinimum *big.Int, recipient common.Address, feeBips *big.Int, feeRecipient common.Address) (*types.Transaction, error) {
	return _V3SwapRouter.contract.Transact(opts, "sweepTokenWithFee", token, amountMinimum, recipient, feeBips, feeRecipient)
}

// SweepTokenWithFee is a paid mutator transaction binding the contract method 0xe0e189a0.
//
// Solidity: function sweepTokenWithFee(address token, uint256 amountMinimum, address recipient, uint256 feeBips, address feeRecipient) payable returns()
func (_V3SwapRouter *V3SwapRouterSession) SweepTokenWithFee(token common.Address, amountMinimum *big.Int, recipient common.Address, feeBips *big.Int, feeRecipient common.Address) (*types.Transaction, error) {
	return _V3SwapRouter.Contract.SweepTokenWithFee(&_V3SwapRouter.TransactOpts, token, amountMinimum, recipient, feeBips, feeRecipient)
}

// SweepTokenWithFee is a paid mutator transaction binding the contract method 0xe0e189a0.
//
// Solidity: function sweepTokenWithFee(address token, uint256 amountMinimum, address recipient, uint256 feeBips, address feeRecipient) payable returns()
func (_V3SwapRouter *V3SwapRouterTransactorSession) SweepTokenWithFee(token common.Address, amountMinimum *big.Int, recipient common.Address, feeBips *big.Int, feeRecipient common.Address) (*types.Transaction, error) {
	return _V3SwapRouter.Contract.SweepTokenWithFee(&_V3SwapRouter.TransactOpts, token, amountMinimum, recipient, feeBips, feeRecipient)
}

// UniswapV3SwapCallback is a paid mutator transaction binding the contract method 0xfa461e33.
//
// Solidity: function uniswapV3SwapCallback(int256 amount0Delta, int256 amount1Delta, bytes _data) returns()
func (_V3SwapRouter *V3SwapRouterTransactor) UniswapV3SwapCallback(opts *bind.TransactOpts, amount0Delta *big.Int, amount1Delta *big.Int, _data []byte) (*types.Transaction, error) {
	return _V3SwapRouter.contract.Transact(opts, "uniswapV3SwapCallback", amount0Delta, amount1Delta, _data)
}

// UniswapV3SwapCallback is a paid mutator transaction binding the contract method 0xfa461e33.
//
// Solidity: function uniswapV3SwapCallback(int256 amount0Delta, int256 amount1Delta, bytes _data) returns()
func (_V3SwapRouter *V3SwapRouterSession) UniswapV3SwapCallback(amount0Delta *big.Int, amount1Delta *big.Int, _data []byte) (*types.Transaction, error) {
	return _V3SwapRouter.Contract.UniswapV3SwapCallback(&_V3SwapRouter.TransactOpts, amount0Delta, amount1Delta, _data)
}

// UniswapV3SwapCallback is a paid mutator transaction binding the contract method 0xfa461e33.
//
// Solidity: function uniswapV3SwapCallback(int256 amount0Delta, int256 amount1Delta, bytes _data) returns()
func (_V3SwapRouter *V3SwapRouterTransactorSession) UniswapV3SwapCallback(amount0Delta *big.Int, amount1Delta *big.Int, _data []byte) (*types.Transaction, error) {
	return _V3SwapRouter.Contract.UniswapV3SwapCallback(&_V3SwapRouter.TransactOpts, amount0Delta, amount1Delta, _data)
}

// UnwrapWETH9 is a paid mutator transaction binding the contract method 0x49404b7c.
//
// Solidity: function unwrapWETH9(uint256 amountMinimum, address recipient) payable returns()
func (_V3SwapRouter *V3SwapRouterTransactor) UnwrapWETH9(opts *bind.TransactOpts, amountMinimum *big.Int, recipient common.Address) (*types.Transaction, error) {
	return _V3SwapRouter.contract.Transact(opts, "unwrapWETH9", amountMinimum, recipient)
}

// UnwrapWETH9 is a paid mutator transaction binding the contract method 0x49404b7c.
//
// Solidity: function unwrapWETH9(uint256 amountMinimum, address recipient) payable returns()
func (_V3SwapRouter *V3SwapRouterSession) UnwrapWETH9(amountMinimum *big.Int, recipient common.Address) (*types.Transaction, error) {
	return _V3SwapRouter.Contract.UnwrapWETH9(&_V3SwapRouter.TransactOpts, amountMinimum, recipient)
}

// UnwrapWETH9 is a paid mutator transaction binding the contract method 0x49404b7c.
//
// Solidity: function unwrapWETH9(uint256 amountMinimum, address recipient) payable returns()
func (_V3SwapRouter *V3SwapRouterTransactorSession) UnwrapWETH9(amountMinimum *big.Int, recipient common.Address) (*types.Transaction, error) {
	return _V3SwapRouter.Contract.UnwrapWETH9(&_V3SwapRouter.TransactOpts, amountMinimum, recipient)
}

// UnwrapWETH9WithFee is a paid mutator transaction binding the contract method 0x9b2c0a37.
//
// Solidity: function unwrapWETH9WithFee(uint256 amountMinimum, address recipient, uint256 feeBips, address feeRecipient) payable returns()
func (_V3SwapRouter *V3SwapRouterTransactor) UnwrapWETH9WithFee(opts *bind.TransactOpts, amountMinimum *big.Int, recipient common.Address, feeBips *big.Int, feeRecipient common.Address) (*types.Transaction, error) {
	return _V3SwapRouter.contract.Transact(opts, "unwrapWETH9WithFee", amountMinimum, recipient, feeBips, feeRecipient)
}

// UnwrapWETH9WithFee is a paid mutator transaction binding the contract method 0x9b2c0a37.
//
// Solidity: function unwrapWETH9WithFee(uint256 amountMinimum, address recipient, uint256 feeBips, address feeRecipient) payable returns()
func (_V3SwapRouter *V3SwapRouterSession) UnwrapWETH9WithFee(amountMinimum *big.Int, recipient common.Address, feeBips *big.Int, feeRecipient common.Address) (*types.Transaction, error) {
	return _V3SwapRouter.Contract.UnwrapWETH9WithFee(&_V3SwapRouter.TransactOpts, amountMinimum, recipient, feeBips, feeRecipient)
}

// UnwrapWETH9WithFee is a paid mutator transaction binding the contract method 0x9b2c0a37.
//
// Solidity: function unwrapWETH9WithFee(uint256 amountMinimum, address recipient, uint256 feeBips, address feeRecipient) payable returns()
func (_V3SwapRouter *V3SwapRouterTransactorSession) UnwrapWETH9WithFee(amountMinimum *big.Int, recipient common.Address, feeBips *big.Int, feeRecipient common.Address) (*types.Transaction, error) {
	return _V3SwapRouter.Contract.UnwrapWETH9WithFee(&_V3SwapRouter.TransactOpts, amountMinimum, recipient, feeBips, feeRecipient)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_V3SwapRouter *V3SwapRouterTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _V3SwapRouter.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_V3SwapRouter *V3SwapRouterSession) Receive() (*types.Transaction, error) {
	return _V3SwapRouter.Contract.Receive(&_V3SwapRouter.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_V3SwapRouter *V3SwapRouterTransactorSession) Receive() (*types.Transaction, error) {
	return _V3SwapRouter.Contract.Receive(&_V3SwapRouter.TransactOpts)
}
