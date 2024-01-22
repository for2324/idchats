// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contract

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

// BiBotFeeRewardUserInfo is an auto generated low-level Go binding around an user-defined struct.
type BiBotFeeRewardUserInfo struct {
	RewardDebt *big.Int
}

// BBTTradeRewardMetaData contains all meta data concerning the BBTTradeReward contract.
var BBTTradeRewardMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_bibotPool\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"rewardDebt\",\"type\":\"uint256\"}],\"indexed\":false,\"internalType\":\"structBiBotFeeReward.UserInfo\",\"name\":\"userInfo\",\"type\":\"tuple\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"customData\",\"type\":\"string\"}],\"name\":\"ClaimReward\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newPool\",\"type\":\"address\"}],\"name\":\"UpdateNewPool\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"admins\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"biBotPool\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"reward\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"customData\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"claimReward\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"nonces\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"recoverFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rewardToken\",\"outputs\":[{\"internalType\":\"contractIERC20Metadata\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"_isAdmin\",\"type\":\"bool\"}],\"name\":\"setAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIERC20Metadata\",\"name\":\"_rewardToken\",\"type\":\"address\"}],\"name\":\"setRewardToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_newPool\",\"type\":\"address\"}],\"name\":\"updateBiBotPool\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"userInfo\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"rewardDebt\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawCurrency\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// BBTTradeRewardABI is the input ABI used to generate the binding from.
// Deprecated: Use BBTTradeRewardMetaData.ABI instead.
var BBTTradeRewardABI = BBTTradeRewardMetaData.ABI

// BBTTradeReward is an auto generated Go binding around an Ethereum contract.
type BBTTradeReward struct {
	BBTTradeRewardCaller     // Read-only binding to the contract
	BBTTradeRewardTransactor // Write-only binding to the contract
	BBTTradeRewardFilterer   // Log filterer for contract events
}

// BBTTradeRewardCaller is an auto generated read-only Go binding around an Ethereum contract.
type BBTTradeRewardCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BBTTradeRewardTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BBTTradeRewardTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BBTTradeRewardFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BBTTradeRewardFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BBTTradeRewardSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BBTTradeRewardSession struct {
	Contract     *BBTTradeReward   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BBTTradeRewardCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BBTTradeRewardCallerSession struct {
	Contract *BBTTradeRewardCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// BBTTradeRewardTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BBTTradeRewardTransactorSession struct {
	Contract     *BBTTradeRewardTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// BBTTradeRewardRaw is an auto generated low-level Go binding around an Ethereum contract.
type BBTTradeRewardRaw struct {
	Contract *BBTTradeReward // Generic contract binding to access the raw methods on
}

// BBTTradeRewardCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BBTTradeRewardCallerRaw struct {
	Contract *BBTTradeRewardCaller // Generic read-only contract binding to access the raw methods on
}

// BBTTradeRewardTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BBTTradeRewardTransactorRaw struct {
	Contract *BBTTradeRewardTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBBTTradeReward creates a new instance of BBTTradeReward, bound to a specific deployed contract.
func NewBBTTradeReward(address common.Address, backend bind.ContractBackend) (*BBTTradeReward, error) {
	contract, err := bindBBTTradeReward(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BBTTradeReward{BBTTradeRewardCaller: BBTTradeRewardCaller{contract: contract}, BBTTradeRewardTransactor: BBTTradeRewardTransactor{contract: contract}, BBTTradeRewardFilterer: BBTTradeRewardFilterer{contract: contract}}, nil
}

// NewBBTTradeRewardCaller creates a new read-only instance of BBTTradeReward, bound to a specific deployed contract.
func NewBBTTradeRewardCaller(address common.Address, caller bind.ContractCaller) (*BBTTradeRewardCaller, error) {
	contract, err := bindBBTTradeReward(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BBTTradeRewardCaller{contract: contract}, nil
}

// NewBBTTradeRewardTransactor creates a new write-only instance of BBTTradeReward, bound to a specific deployed contract.
func NewBBTTradeRewardTransactor(address common.Address, transactor bind.ContractTransactor) (*BBTTradeRewardTransactor, error) {
	contract, err := bindBBTTradeReward(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BBTTradeRewardTransactor{contract: contract}, nil
}

// NewBBTTradeRewardFilterer creates a new log filterer instance of BBTTradeReward, bound to a specific deployed contract.
func NewBBTTradeRewardFilterer(address common.Address, filterer bind.ContractFilterer) (*BBTTradeRewardFilterer, error) {
	contract, err := bindBBTTradeReward(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BBTTradeRewardFilterer{contract: contract}, nil
}

// bindBBTTradeReward binds a generic wrapper to an already deployed contract.
func bindBBTTradeReward(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := BBTTradeRewardMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BBTTradeReward *BBTTradeRewardRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BBTTradeReward.Contract.BBTTradeRewardCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BBTTradeReward *BBTTradeRewardRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BBTTradeReward.Contract.BBTTradeRewardTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BBTTradeReward *BBTTradeRewardRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BBTTradeReward.Contract.BBTTradeRewardTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BBTTradeReward *BBTTradeRewardCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BBTTradeReward.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BBTTradeReward *BBTTradeRewardTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BBTTradeReward.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BBTTradeReward *BBTTradeRewardTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BBTTradeReward.Contract.contract.Transact(opts, method, params...)
}

// Admins is a free data retrieval call binding the contract method 0x429b62e5.
//
// Solidity: function admins(address ) view returns(bool)
func (_BBTTradeReward *BBTTradeRewardCaller) Admins(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _BBTTradeReward.contract.Call(opts, &out, "admins", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Admins is a free data retrieval call binding the contract method 0x429b62e5.
//
// Solidity: function admins(address ) view returns(bool)
func (_BBTTradeReward *BBTTradeRewardSession) Admins(arg0 common.Address) (bool, error) {
	return _BBTTradeReward.Contract.Admins(&_BBTTradeReward.CallOpts, arg0)
}

// Admins is a free data retrieval call binding the contract method 0x429b62e5.
//
// Solidity: function admins(address ) view returns(bool)
func (_BBTTradeReward *BBTTradeRewardCallerSession) Admins(arg0 common.Address) (bool, error) {
	return _BBTTradeReward.Contract.Admins(&_BBTTradeReward.CallOpts, arg0)
}

// BiBotPool is a free data retrieval call binding the contract method 0xe1a9681d.
//
// Solidity: function biBotPool() view returns(address)
func (_BBTTradeReward *BBTTradeRewardCaller) BiBotPool(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BBTTradeReward.contract.Call(opts, &out, "biBotPool")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// BiBotPool is a free data retrieval call binding the contract method 0xe1a9681d.
//
// Solidity: function biBotPool() view returns(address)
func (_BBTTradeReward *BBTTradeRewardSession) BiBotPool() (common.Address, error) {
	return _BBTTradeReward.Contract.BiBotPool(&_BBTTradeReward.CallOpts)
}

// BiBotPool is a free data retrieval call binding the contract method 0xe1a9681d.
//
// Solidity: function biBotPool() view returns(address)
func (_BBTTradeReward *BBTTradeRewardCallerSession) BiBotPool() (common.Address, error) {
	return _BBTTradeReward.Contract.BiBotPool(&_BBTTradeReward.CallOpts)
}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address ) view returns(uint256)
func (_BBTTradeReward *BBTTradeRewardCaller) Nonces(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _BBTTradeReward.contract.Call(opts, &out, "nonces", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address ) view returns(uint256)
func (_BBTTradeReward *BBTTradeRewardSession) Nonces(arg0 common.Address) (*big.Int, error) {
	return _BBTTradeReward.Contract.Nonces(&_BBTTradeReward.CallOpts, arg0)
}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address ) view returns(uint256)
func (_BBTTradeReward *BBTTradeRewardCallerSession) Nonces(arg0 common.Address) (*big.Int, error) {
	return _BBTTradeReward.Contract.Nonces(&_BBTTradeReward.CallOpts, arg0)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_BBTTradeReward *BBTTradeRewardCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BBTTradeReward.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_BBTTradeReward *BBTTradeRewardSession) Owner() (common.Address, error) {
	return _BBTTradeReward.Contract.Owner(&_BBTTradeReward.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_BBTTradeReward *BBTTradeRewardCallerSession) Owner() (common.Address, error) {
	return _BBTTradeReward.Contract.Owner(&_BBTTradeReward.CallOpts)
}

// RewardToken is a free data retrieval call binding the contract method 0xf7c618c1.
//
// Solidity: function rewardToken() view returns(address)
func (_BBTTradeReward *BBTTradeRewardCaller) RewardToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BBTTradeReward.contract.Call(opts, &out, "rewardToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RewardToken is a free data retrieval call binding the contract method 0xf7c618c1.
//
// Solidity: function rewardToken() view returns(address)
func (_BBTTradeReward *BBTTradeRewardSession) RewardToken() (common.Address, error) {
	return _BBTTradeReward.Contract.RewardToken(&_BBTTradeReward.CallOpts)
}

// RewardToken is a free data retrieval call binding the contract method 0xf7c618c1.
//
// Solidity: function rewardToken() view returns(address)
func (_BBTTradeReward *BBTTradeRewardCallerSession) RewardToken() (common.Address, error) {
	return _BBTTradeReward.Contract.RewardToken(&_BBTTradeReward.CallOpts)
}

// UserInfo is a free data retrieval call binding the contract method 0x1959a002.
//
// Solidity: function userInfo(address ) view returns(uint256 rewardDebt)
func (_BBTTradeReward *BBTTradeRewardCaller) UserInfo(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _BBTTradeReward.contract.Call(opts, &out, "userInfo", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// UserInfo is a free data retrieval call binding the contract method 0x1959a002.
//
// Solidity: function userInfo(address ) view returns(uint256 rewardDebt)
func (_BBTTradeReward *BBTTradeRewardSession) UserInfo(arg0 common.Address) (*big.Int, error) {
	return _BBTTradeReward.Contract.UserInfo(&_BBTTradeReward.CallOpts, arg0)
}

// UserInfo is a free data retrieval call binding the contract method 0x1959a002.
//
// Solidity: function userInfo(address ) view returns(uint256 rewardDebt)
func (_BBTTradeReward *BBTTradeRewardCallerSession) UserInfo(arg0 common.Address) (*big.Int, error) {
	return _BBTTradeReward.Contract.UserInfo(&_BBTTradeReward.CallOpts, arg0)
}

// ClaimReward is a paid mutator transaction binding the contract method 0xcd08b4fb.
//
// Solidity: function claimReward(uint256 reward, address recipient, uint256 nonce, string customData, uint8 v, bytes32 r, bytes32 s) returns()
func (_BBTTradeReward *BBTTradeRewardTransactor) ClaimReward(opts *bind.TransactOpts, reward *big.Int, recipient common.Address, nonce *big.Int, customData string, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _BBTTradeReward.contract.Transact(opts, "claimReward", reward, recipient, nonce, customData, v, r, s)
}

// ClaimReward is a paid mutator transaction binding the contract method 0xcd08b4fb.
//
// Solidity: function claimReward(uint256 reward, address recipient, uint256 nonce, string customData, uint8 v, bytes32 r, bytes32 s) returns()
func (_BBTTradeReward *BBTTradeRewardSession) ClaimReward(reward *big.Int, recipient common.Address, nonce *big.Int, customData string, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _BBTTradeReward.Contract.ClaimReward(&_BBTTradeReward.TransactOpts, reward, recipient, nonce, customData, v, r, s)
}

// ClaimReward is a paid mutator transaction binding the contract method 0xcd08b4fb.
//
// Solidity: function claimReward(uint256 reward, address recipient, uint256 nonce, string customData, uint8 v, bytes32 r, bytes32 s) returns()
func (_BBTTradeReward *BBTTradeRewardTransactorSession) ClaimReward(reward *big.Int, recipient common.Address, nonce *big.Int, customData string, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _BBTTradeReward.Contract.ClaimReward(&_BBTTradeReward.TransactOpts, reward, recipient, nonce, customData, v, r, s)
}

// RecoverFunds is a paid mutator transaction binding the contract method 0x5d3590d5.
//
// Solidity: function recoverFunds(address _token, address _to, uint256 _amount) returns()
func (_BBTTradeReward *BBTTradeRewardTransactor) RecoverFunds(opts *bind.TransactOpts, _token common.Address, _to common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _BBTTradeReward.contract.Transact(opts, "recoverFunds", _token, _to, _amount)
}

// RecoverFunds is a paid mutator transaction binding the contract method 0x5d3590d5.
//
// Solidity: function recoverFunds(address _token, address _to, uint256 _amount) returns()
func (_BBTTradeReward *BBTTradeRewardSession) RecoverFunds(_token common.Address, _to common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _BBTTradeReward.Contract.RecoverFunds(&_BBTTradeReward.TransactOpts, _token, _to, _amount)
}

// RecoverFunds is a paid mutator transaction binding the contract method 0x5d3590d5.
//
// Solidity: function recoverFunds(address _token, address _to, uint256 _amount) returns()
func (_BBTTradeReward *BBTTradeRewardTransactorSession) RecoverFunds(_token common.Address, _to common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _BBTTradeReward.Contract.RecoverFunds(&_BBTTradeReward.TransactOpts, _token, _to, _amount)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_BBTTradeReward *BBTTradeRewardTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BBTTradeReward.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_BBTTradeReward *BBTTradeRewardSession) RenounceOwnership() (*types.Transaction, error) {
	return _BBTTradeReward.Contract.RenounceOwnership(&_BBTTradeReward.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_BBTTradeReward *BBTTradeRewardTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _BBTTradeReward.Contract.RenounceOwnership(&_BBTTradeReward.TransactOpts)
}

// SetAdmin is a paid mutator transaction binding the contract method 0x4b0bddd2.
//
// Solidity: function setAdmin(address account, bool _isAdmin) returns()
func (_BBTTradeReward *BBTTradeRewardTransactor) SetAdmin(opts *bind.TransactOpts, account common.Address, _isAdmin bool) (*types.Transaction, error) {
	return _BBTTradeReward.contract.Transact(opts, "setAdmin", account, _isAdmin)
}

// SetAdmin is a paid mutator transaction binding the contract method 0x4b0bddd2.
//
// Solidity: function setAdmin(address account, bool _isAdmin) returns()
func (_BBTTradeReward *BBTTradeRewardSession) SetAdmin(account common.Address, _isAdmin bool) (*types.Transaction, error) {
	return _BBTTradeReward.Contract.SetAdmin(&_BBTTradeReward.TransactOpts, account, _isAdmin)
}

// SetAdmin is a paid mutator transaction binding the contract method 0x4b0bddd2.
//
// Solidity: function setAdmin(address account, bool _isAdmin) returns()
func (_BBTTradeReward *BBTTradeRewardTransactorSession) SetAdmin(account common.Address, _isAdmin bool) (*types.Transaction, error) {
	return _BBTTradeReward.Contract.SetAdmin(&_BBTTradeReward.TransactOpts, account, _isAdmin)
}

// SetRewardToken is a paid mutator transaction binding the contract method 0x8aee8127.
//
// Solidity: function setRewardToken(address _rewardToken) returns()
func (_BBTTradeReward *BBTTradeRewardTransactor) SetRewardToken(opts *bind.TransactOpts, _rewardToken common.Address) (*types.Transaction, error) {
	return _BBTTradeReward.contract.Transact(opts, "setRewardToken", _rewardToken)
}

// SetRewardToken is a paid mutator transaction binding the contract method 0x8aee8127.
//
// Solidity: function setRewardToken(address _rewardToken) returns()
func (_BBTTradeReward *BBTTradeRewardSession) SetRewardToken(_rewardToken common.Address) (*types.Transaction, error) {
	return _BBTTradeReward.Contract.SetRewardToken(&_BBTTradeReward.TransactOpts, _rewardToken)
}

// SetRewardToken is a paid mutator transaction binding the contract method 0x8aee8127.
//
// Solidity: function setRewardToken(address _rewardToken) returns()
func (_BBTTradeReward *BBTTradeRewardTransactorSession) SetRewardToken(_rewardToken common.Address) (*types.Transaction, error) {
	return _BBTTradeReward.Contract.SetRewardToken(&_BBTTradeReward.TransactOpts, _rewardToken)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_BBTTradeReward *BBTTradeRewardTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _BBTTradeReward.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_BBTTradeReward *BBTTradeRewardSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _BBTTradeReward.Contract.TransferOwnership(&_BBTTradeReward.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_BBTTradeReward *BBTTradeRewardTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _BBTTradeReward.Contract.TransferOwnership(&_BBTTradeReward.TransactOpts, newOwner)
}

// UpdateBiBotPool is a paid mutator transaction binding the contract method 0x52da1f81.
//
// Solidity: function updateBiBotPool(address _newPool) returns()
func (_BBTTradeReward *BBTTradeRewardTransactor) UpdateBiBotPool(opts *bind.TransactOpts, _newPool common.Address) (*types.Transaction, error) {
	return _BBTTradeReward.contract.Transact(opts, "updateBiBotPool", _newPool)
}

// UpdateBiBotPool is a paid mutator transaction binding the contract method 0x52da1f81.
//
// Solidity: function updateBiBotPool(address _newPool) returns()
func (_BBTTradeReward *BBTTradeRewardSession) UpdateBiBotPool(_newPool common.Address) (*types.Transaction, error) {
	return _BBTTradeReward.Contract.UpdateBiBotPool(&_BBTTradeReward.TransactOpts, _newPool)
}

// UpdateBiBotPool is a paid mutator transaction binding the contract method 0x52da1f81.
//
// Solidity: function updateBiBotPool(address _newPool) returns()
func (_BBTTradeReward *BBTTradeRewardTransactorSession) UpdateBiBotPool(_newPool common.Address) (*types.Transaction, error) {
	return _BBTTradeReward.Contract.UpdateBiBotPool(&_BBTTradeReward.TransactOpts, _newPool)
}

// WithdrawCurrency is a paid mutator transaction binding the contract method 0x65b99f63.
//
// Solidity: function withdrawCurrency() returns()
func (_BBTTradeReward *BBTTradeRewardTransactor) WithdrawCurrency(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BBTTradeReward.contract.Transact(opts, "withdrawCurrency")
}

// WithdrawCurrency is a paid mutator transaction binding the contract method 0x65b99f63.
//
// Solidity: function withdrawCurrency() returns()
func (_BBTTradeReward *BBTTradeRewardSession) WithdrawCurrency() (*types.Transaction, error) {
	return _BBTTradeReward.Contract.WithdrawCurrency(&_BBTTradeReward.TransactOpts)
}

// WithdrawCurrency is a paid mutator transaction binding the contract method 0x65b99f63.
//
// Solidity: function withdrawCurrency() returns()
func (_BBTTradeReward *BBTTradeRewardTransactorSession) WithdrawCurrency() (*types.Transaction, error) {
	return _BBTTradeReward.Contract.WithdrawCurrency(&_BBTTradeReward.TransactOpts)
}

// BBTTradeRewardClaimRewardIterator is returned from FilterClaimReward and is used to iterate over the raw logs and unpacked data for ClaimReward events raised by the BBTTradeReward contract.
type BBTTradeRewardClaimRewardIterator struct {
	Event *BBTTradeRewardClaimReward // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *BBTTradeRewardClaimRewardIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BBTTradeRewardClaimReward)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(BBTTradeRewardClaimReward)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *BBTTradeRewardClaimRewardIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BBTTradeRewardClaimRewardIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BBTTradeRewardClaimReward represents a ClaimReward event raised by the BBTTradeReward contract.
type BBTTradeRewardClaimReward struct {
	User       common.Address
	Amount     *big.Int
	Nonce      *big.Int
	UserInfo   BiBotFeeRewardUserInfo
	CustomData string
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterClaimReward is a free log retrieval operation binding the contract event 0x6ced2948df2f38d8a509f4cf84d2bda0df194d9bf6c15a86473c17151b374598.
//
// Solidity: event ClaimReward(address indexed user, uint256 amount, uint256 nonce, (uint256) userInfo, string customData)
func (_BBTTradeReward *BBTTradeRewardFilterer) FilterClaimReward(opts *bind.FilterOpts, user []common.Address) (*BBTTradeRewardClaimRewardIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _BBTTradeReward.contract.FilterLogs(opts, "ClaimReward", userRule)
	if err != nil {
		return nil, err
	}
	return &BBTTradeRewardClaimRewardIterator{contract: _BBTTradeReward.contract, event: "ClaimReward", logs: logs, sub: sub}, nil
}

// WatchClaimReward is a free log subscription operation binding the contract event 0x6ced2948df2f38d8a509f4cf84d2bda0df194d9bf6c15a86473c17151b374598.
//
// Solidity: event ClaimReward(address indexed user, uint256 amount, uint256 nonce, (uint256) userInfo, string customData)
func (_BBTTradeReward *BBTTradeRewardFilterer) WatchClaimReward(opts *bind.WatchOpts, sink chan<- *BBTTradeRewardClaimReward, user []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _BBTTradeReward.contract.WatchLogs(opts, "ClaimReward", userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BBTTradeRewardClaimReward)
				if err := _BBTTradeReward.contract.UnpackLog(event, "ClaimReward", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseClaimReward is a log parse operation binding the contract event 0x6ced2948df2f38d8a509f4cf84d2bda0df194d9bf6c15a86473c17151b374598.
//
// Solidity: event ClaimReward(address indexed user, uint256 amount, uint256 nonce, (uint256) userInfo, string customData)
func (_BBTTradeReward *BBTTradeRewardFilterer) ParseClaimReward(log types.Log) (*BBTTradeRewardClaimReward, error) {
	event := new(BBTTradeRewardClaimReward)
	if err := _BBTTradeReward.contract.UnpackLog(event, "ClaimReward", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BBTTradeRewardOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the BBTTradeReward contract.
type BBTTradeRewardOwnershipTransferredIterator struct {
	Event *BBTTradeRewardOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *BBTTradeRewardOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BBTTradeRewardOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(BBTTradeRewardOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *BBTTradeRewardOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BBTTradeRewardOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BBTTradeRewardOwnershipTransferred represents a OwnershipTransferred event raised by the BBTTradeReward contract.
type BBTTradeRewardOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_BBTTradeReward *BBTTradeRewardFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*BBTTradeRewardOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _BBTTradeReward.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &BBTTradeRewardOwnershipTransferredIterator{contract: _BBTTradeReward.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_BBTTradeReward *BBTTradeRewardFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *BBTTradeRewardOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _BBTTradeReward.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BBTTradeRewardOwnershipTransferred)
				if err := _BBTTradeReward.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_BBTTradeReward *BBTTradeRewardFilterer) ParseOwnershipTransferred(log types.Log) (*BBTTradeRewardOwnershipTransferred, error) {
	event := new(BBTTradeRewardOwnershipTransferred)
	if err := _BBTTradeReward.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BBTTradeRewardUpdateNewPoolIterator is returned from FilterUpdateNewPool and is used to iterate over the raw logs and unpacked data for UpdateNewPool events raised by the BBTTradeReward contract.
type BBTTradeRewardUpdateNewPoolIterator struct {
	Event *BBTTradeRewardUpdateNewPool // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *BBTTradeRewardUpdateNewPoolIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BBTTradeRewardUpdateNewPool)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(BBTTradeRewardUpdateNewPool)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *BBTTradeRewardUpdateNewPoolIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BBTTradeRewardUpdateNewPoolIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BBTTradeRewardUpdateNewPool represents a UpdateNewPool event raised by the BBTTradeReward contract.
type BBTTradeRewardUpdateNewPool struct {
	NewPool common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUpdateNewPool is a free log retrieval operation binding the contract event 0x745b015da9f343b4dc95e9940b81cdfa98193e1dc3f9d0e16534c6ca5418725b.
//
// Solidity: event UpdateNewPool(address newPool)
func (_BBTTradeReward *BBTTradeRewardFilterer) FilterUpdateNewPool(opts *bind.FilterOpts) (*BBTTradeRewardUpdateNewPoolIterator, error) {

	logs, sub, err := _BBTTradeReward.contract.FilterLogs(opts, "UpdateNewPool")
	if err != nil {
		return nil, err
	}
	return &BBTTradeRewardUpdateNewPoolIterator{contract: _BBTTradeReward.contract, event: "UpdateNewPool", logs: logs, sub: sub}, nil
}

// WatchUpdateNewPool is a free log subscription operation binding the contract event 0x745b015da9f343b4dc95e9940b81cdfa98193e1dc3f9d0e16534c6ca5418725b.
//
// Solidity: event UpdateNewPool(address newPool)
func (_BBTTradeReward *BBTTradeRewardFilterer) WatchUpdateNewPool(opts *bind.WatchOpts, sink chan<- *BBTTradeRewardUpdateNewPool) (event.Subscription, error) {

	logs, sub, err := _BBTTradeReward.contract.WatchLogs(opts, "UpdateNewPool")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BBTTradeRewardUpdateNewPool)
				if err := _BBTTradeReward.contract.UnpackLog(event, "UpdateNewPool", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUpdateNewPool is a log parse operation binding the contract event 0x745b015da9f343b4dc95e9940b81cdfa98193e1dc3f9d0e16534c6ca5418725b.
//
// Solidity: event UpdateNewPool(address newPool)
func (_BBTTradeReward *BBTTradeRewardFilterer) ParseUpdateNewPool(log types.Log) (*BBTTradeRewardUpdateNewPool, error) {
	event := new(BBTTradeRewardUpdateNewPool)
	if err := _BBTTradeReward.contract.UnpackLog(event, "UpdateNewPool", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
