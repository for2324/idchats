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

// StakingStakingInfo is an auto generated low-level Go binding around an user-defined struct.
type StakingStakingInfo struct {
	Id            *big.Int
	StakedAmount  *big.Int
	StakingPeriod *big.Int
	StartTime     *big.Int
}

// BBTPledgePoolMetaData contains all meta data concerning the BBTPledgePool contract.
var BBTPledgePoolMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"periodForMinimumRate_\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"stakingPeriod1_\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"stakingPeriod2_\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"stakingPeriod3_\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"minimumRate_\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"rewardsRate1_\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"rewardsRate2_\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"rewardsRate3_\",\"type\":\"uint64\"},{\"internalType\":\"contractIERC20\",\"name\":\"token_\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"stakedAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"stakingPeriod\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startTime\",\"type\":\"uint256\"}],\"name\":\"Stake\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"stakedAmount\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"withdrawAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"stakingPeriod\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startTime\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"endTime\",\"type\":\"uint256\"}],\"name\":\"Withdraw\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"adminWithdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"ratio\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"principal\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"n\",\"type\":\"uint256\"}],\"name\":\"calculateCompound\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"stakingPeriod\",\"type\":\"uint256\"}],\"name\":\"extendStakingPeriod\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"getUserStakingInfo\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"stakedAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"stakingPeriod\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"startTime\",\"type\":\"uint256\"}],\"internalType\":\"structStaking.StakingInfo[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"idCounter\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minimumRate\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rewardsRate1\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rewardsRate2\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rewardsRate3\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"_minimumRate\",\"type\":\"uint64\"}],\"name\":\"setMinimumRate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"_rewardsRate1\",\"type\":\"uint64\"}],\"name\":\"setRewardsRate1\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"_rewardsRate2\",\"type\":\"uint64\"}],\"name\":\"setRewardsRate2\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"_rewardsRate3\",\"type\":\"uint64\"}],\"name\":\"setRewardsRate3\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"stakingPeriod\",\"type\":\"uint256\"}],\"name\":\"stake\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"amounts\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"stakingPeriods\",\"type\":\"uint256[]\"}],\"name\":\"stakeBatch\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"token\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalLocked\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"withdrawAll\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"ids\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"amounts\",\"type\":\"uint256[]\"}],\"name\":\"withdrawBatch\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"withdrawableAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// BBTPledgePoolABI is the input ABI used to generate the binding from.
// Deprecated: Use BBTPledgePoolMetaData.ABI instead.
var BBTPledgePoolABI = BBTPledgePoolMetaData.ABI

// BBTPledgePool is an auto generated Go binding around an Ethereum contract.
type BBTPledgePool struct {
	BBTPledgePoolCaller     // Read-only binding to the contract
	BBTPledgePoolTransactor // Write-only binding to the contract
	BBTPledgePoolFilterer   // Log filterer for contract events
}

// BBTPledgePoolCaller is an auto generated read-only Go binding around an Ethereum contract.
type BBTPledgePoolCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BBTPledgePoolTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BBTPledgePoolTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BBTPledgePoolFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BBTPledgePoolFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BBTPledgePoolSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BBTPledgePoolSession struct {
	Contract     *BBTPledgePool    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BBTPledgePoolCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BBTPledgePoolCallerSession struct {
	Contract *BBTPledgePoolCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// BBTPledgePoolTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BBTPledgePoolTransactorSession struct {
	Contract     *BBTPledgePoolTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// BBTPledgePoolRaw is an auto generated low-level Go binding around an Ethereum contract.
type BBTPledgePoolRaw struct {
	Contract *BBTPledgePool // Generic contract binding to access the raw methods on
}

// BBTPledgePoolCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BBTPledgePoolCallerRaw struct {
	Contract *BBTPledgePoolCaller // Generic read-only contract binding to access the raw methods on
}

// BBTPledgePoolTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BBTPledgePoolTransactorRaw struct {
	Contract *BBTPledgePoolTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBBTPledgePool creates a new instance of BBTPledgePool, bound to a specific deployed contract.
func NewBBTPledgePool(address common.Address, backend bind.ContractBackend) (*BBTPledgePool, error) {
	contract, err := bindBBTPledgePool(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BBTPledgePool{BBTPledgePoolCaller: BBTPledgePoolCaller{contract: contract}, BBTPledgePoolTransactor: BBTPledgePoolTransactor{contract: contract}, BBTPledgePoolFilterer: BBTPledgePoolFilterer{contract: contract}}, nil
}

// NewBBTPledgePoolCaller creates a new read-only instance of BBTPledgePool, bound to a specific deployed contract.
func NewBBTPledgePoolCaller(address common.Address, caller bind.ContractCaller) (*BBTPledgePoolCaller, error) {
	contract, err := bindBBTPledgePool(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BBTPledgePoolCaller{contract: contract}, nil
}

// NewBBTPledgePoolTransactor creates a new write-only instance of BBTPledgePool, bound to a specific deployed contract.
func NewBBTPledgePoolTransactor(address common.Address, transactor bind.ContractTransactor) (*BBTPledgePoolTransactor, error) {
	contract, err := bindBBTPledgePool(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BBTPledgePoolTransactor{contract: contract}, nil
}

// NewBBTPledgePoolFilterer creates a new log filterer instance of BBTPledgePool, bound to a specific deployed contract.
func NewBBTPledgePoolFilterer(address common.Address, filterer bind.ContractFilterer) (*BBTPledgePoolFilterer, error) {
	contract, err := bindBBTPledgePool(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BBTPledgePoolFilterer{contract: contract}, nil
}

// bindBBTPledgePool binds a generic wrapper to an already deployed contract.
func bindBBTPledgePool(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := BBTPledgePoolMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BBTPledgePool *BBTPledgePoolRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BBTPledgePool.Contract.BBTPledgePoolCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BBTPledgePool *BBTPledgePoolRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BBTPledgePool.Contract.BBTPledgePoolTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BBTPledgePool *BBTPledgePoolRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BBTPledgePool.Contract.BBTPledgePoolTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BBTPledgePool *BBTPledgePoolCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BBTPledgePool.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BBTPledgePool *BBTPledgePoolTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BBTPledgePool.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BBTPledgePool *BBTPledgePoolTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BBTPledgePool.Contract.contract.Transact(opts, method, params...)
}

// CalculateCompound is a free data retrieval call binding the contract method 0xc6928f19.
//
// Solidity: function calculateCompound(uint256 ratio, uint256 principal, uint256 n) pure returns(uint256)
func (_BBTPledgePool *BBTPledgePoolCaller) CalculateCompound(opts *bind.CallOpts, ratio *big.Int, principal *big.Int, n *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _BBTPledgePool.contract.Call(opts, &out, "calculateCompound", ratio, principal, n)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CalculateCompound is a free data retrieval call binding the contract method 0xc6928f19.
//
// Solidity: function calculateCompound(uint256 ratio, uint256 principal, uint256 n) pure returns(uint256)
func (_BBTPledgePool *BBTPledgePoolSession) CalculateCompound(ratio *big.Int, principal *big.Int, n *big.Int) (*big.Int, error) {
	return _BBTPledgePool.Contract.CalculateCompound(&_BBTPledgePool.CallOpts, ratio, principal, n)
}

// CalculateCompound is a free data retrieval call binding the contract method 0xc6928f19.
//
// Solidity: function calculateCompound(uint256 ratio, uint256 principal, uint256 n) pure returns(uint256)
func (_BBTPledgePool *BBTPledgePoolCallerSession) CalculateCompound(ratio *big.Int, principal *big.Int, n *big.Int) (*big.Int, error) {
	return _BBTPledgePool.Contract.CalculateCompound(&_BBTPledgePool.CallOpts, ratio, principal, n)
}

// GetUserStakingInfo is a free data retrieval call binding the contract method 0x80933608.
//
// Solidity: function getUserStakingInfo(address user) view returns((uint256,uint256,uint256,uint256)[])
func (_BBTPledgePool *BBTPledgePoolCaller) GetUserStakingInfo(opts *bind.CallOpts, user common.Address) ([]StakingStakingInfo, error) {
	var out []interface{}
	err := _BBTPledgePool.contract.Call(opts, &out, "getUserStakingInfo", user)

	if err != nil {
		return *new([]StakingStakingInfo), err
	}

	out0 := *abi.ConvertType(out[0], new([]StakingStakingInfo)).(*[]StakingStakingInfo)

	return out0, err

}

// GetUserStakingInfo is a free data retrieval call binding the contract method 0x80933608.
//
// Solidity: function getUserStakingInfo(address user) view returns((uint256,uint256,uint256,uint256)[])
func (_BBTPledgePool *BBTPledgePoolSession) GetUserStakingInfo(user common.Address) ([]StakingStakingInfo, error) {
	return _BBTPledgePool.Contract.GetUserStakingInfo(&_BBTPledgePool.CallOpts, user)
}

// GetUserStakingInfo is a free data retrieval call binding the contract method 0x80933608.
//
// Solidity: function getUserStakingInfo(address user) view returns((uint256,uint256,uint256,uint256)[])
func (_BBTPledgePool *BBTPledgePoolCallerSession) GetUserStakingInfo(user common.Address) ([]StakingStakingInfo, error) {
	return _BBTPledgePool.Contract.GetUserStakingInfo(&_BBTPledgePool.CallOpts, user)
}

// IdCounter is a free data retrieval call binding the contract method 0xeb08ab28.
//
// Solidity: function idCounter() view returns(uint256)
func (_BBTPledgePool *BBTPledgePoolCaller) IdCounter(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _BBTPledgePool.contract.Call(opts, &out, "idCounter")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// IdCounter is a free data retrieval call binding the contract method 0xeb08ab28.
//
// Solidity: function idCounter() view returns(uint256)
func (_BBTPledgePool *BBTPledgePoolSession) IdCounter() (*big.Int, error) {
	return _BBTPledgePool.Contract.IdCounter(&_BBTPledgePool.CallOpts)
}

// IdCounter is a free data retrieval call binding the contract method 0xeb08ab28.
//
// Solidity: function idCounter() view returns(uint256)
func (_BBTPledgePool *BBTPledgePoolCallerSession) IdCounter() (*big.Int, error) {
	return _BBTPledgePool.Contract.IdCounter(&_BBTPledgePool.CallOpts)
}

// MinimumRate is a free data retrieval call binding the contract method 0xfae6dbe3.
//
// Solidity: function minimumRate() view returns(uint64)
func (_BBTPledgePool *BBTPledgePoolCaller) MinimumRate(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _BBTPledgePool.contract.Call(opts, &out, "minimumRate")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// MinimumRate is a free data retrieval call binding the contract method 0xfae6dbe3.
//
// Solidity: function minimumRate() view returns(uint64)
func (_BBTPledgePool *BBTPledgePoolSession) MinimumRate() (uint64, error) {
	return _BBTPledgePool.Contract.MinimumRate(&_BBTPledgePool.CallOpts)
}

// MinimumRate is a free data retrieval call binding the contract method 0xfae6dbe3.
//
// Solidity: function minimumRate() view returns(uint64)
func (_BBTPledgePool *BBTPledgePoolCallerSession) MinimumRate() (uint64, error) {
	return _BBTPledgePool.Contract.MinimumRate(&_BBTPledgePool.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_BBTPledgePool *BBTPledgePoolCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BBTPledgePool.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_BBTPledgePool *BBTPledgePoolSession) Owner() (common.Address, error) {
	return _BBTPledgePool.Contract.Owner(&_BBTPledgePool.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_BBTPledgePool *BBTPledgePoolCallerSession) Owner() (common.Address, error) {
	return _BBTPledgePool.Contract.Owner(&_BBTPledgePool.CallOpts)
}

// RewardsRate1 is a free data retrieval call binding the contract method 0x9da750be.
//
// Solidity: function rewardsRate1() view returns(uint64)
func (_BBTPledgePool *BBTPledgePoolCaller) RewardsRate1(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _BBTPledgePool.contract.Call(opts, &out, "rewardsRate1")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// RewardsRate1 is a free data retrieval call binding the contract method 0x9da750be.
//
// Solidity: function rewardsRate1() view returns(uint64)
func (_BBTPledgePool *BBTPledgePoolSession) RewardsRate1() (uint64, error) {
	return _BBTPledgePool.Contract.RewardsRate1(&_BBTPledgePool.CallOpts)
}

// RewardsRate1 is a free data retrieval call binding the contract method 0x9da750be.
//
// Solidity: function rewardsRate1() view returns(uint64)
func (_BBTPledgePool *BBTPledgePoolCallerSession) RewardsRate1() (uint64, error) {
	return _BBTPledgePool.Contract.RewardsRate1(&_BBTPledgePool.CallOpts)
}

// RewardsRate2 is a free data retrieval call binding the contract method 0x5dbd6e11.
//
// Solidity: function rewardsRate2() view returns(uint64)
func (_BBTPledgePool *BBTPledgePoolCaller) RewardsRate2(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _BBTPledgePool.contract.Call(opts, &out, "rewardsRate2")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// RewardsRate2 is a free data retrieval call binding the contract method 0x5dbd6e11.
//
// Solidity: function rewardsRate2() view returns(uint64)
func (_BBTPledgePool *BBTPledgePoolSession) RewardsRate2() (uint64, error) {
	return _BBTPledgePool.Contract.RewardsRate2(&_BBTPledgePool.CallOpts)
}

// RewardsRate2 is a free data retrieval call binding the contract method 0x5dbd6e11.
//
// Solidity: function rewardsRate2() view returns(uint64)
func (_BBTPledgePool *BBTPledgePoolCallerSession) RewardsRate2() (uint64, error) {
	return _BBTPledgePool.Contract.RewardsRate2(&_BBTPledgePool.CallOpts)
}

// RewardsRate3 is a free data retrieval call binding the contract method 0x8662c533.
//
// Solidity: function rewardsRate3() view returns(uint64)
func (_BBTPledgePool *BBTPledgePoolCaller) RewardsRate3(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _BBTPledgePool.contract.Call(opts, &out, "rewardsRate3")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// RewardsRate3 is a free data retrieval call binding the contract method 0x8662c533.
//
// Solidity: function rewardsRate3() view returns(uint64)
func (_BBTPledgePool *BBTPledgePoolSession) RewardsRate3() (uint64, error) {
	return _BBTPledgePool.Contract.RewardsRate3(&_BBTPledgePool.CallOpts)
}

// RewardsRate3 is a free data retrieval call binding the contract method 0x8662c533.
//
// Solidity: function rewardsRate3() view returns(uint64)
func (_BBTPledgePool *BBTPledgePoolCallerSession) RewardsRate3() (uint64, error) {
	return _BBTPledgePool.Contract.RewardsRate3(&_BBTPledgePool.CallOpts)
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() view returns(address)
func (_BBTPledgePool *BBTPledgePoolCaller) Token(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BBTPledgePool.contract.Call(opts, &out, "token")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() view returns(address)
func (_BBTPledgePool *BBTPledgePoolSession) Token() (common.Address, error) {
	return _BBTPledgePool.Contract.Token(&_BBTPledgePool.CallOpts)
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() view returns(address)
func (_BBTPledgePool *BBTPledgePoolCallerSession) Token() (common.Address, error) {
	return _BBTPledgePool.Contract.Token(&_BBTPledgePool.CallOpts)
}

// TotalLocked is a free data retrieval call binding the contract method 0x56891412.
//
// Solidity: function totalLocked() view returns(uint256)
func (_BBTPledgePool *BBTPledgePoolCaller) TotalLocked(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _BBTPledgePool.contract.Call(opts, &out, "totalLocked")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalLocked is a free data retrieval call binding the contract method 0x56891412.
//
// Solidity: function totalLocked() view returns(uint256)
func (_BBTPledgePool *BBTPledgePoolSession) TotalLocked() (*big.Int, error) {
	return _BBTPledgePool.Contract.TotalLocked(&_BBTPledgePool.CallOpts)
}

// TotalLocked is a free data retrieval call binding the contract method 0x56891412.
//
// Solidity: function totalLocked() view returns(uint256)
func (_BBTPledgePool *BBTPledgePoolCallerSession) TotalLocked() (*big.Int, error) {
	return _BBTPledgePool.Contract.TotalLocked(&_BBTPledgePool.CallOpts)
}

// WithdrawableAmount is a free data retrieval call binding the contract method 0x7831a74f.
//
// Solidity: function withdrawableAmount(uint256 id, address user) view returns(uint256)
func (_BBTPledgePool *BBTPledgePoolCaller) WithdrawableAmount(opts *bind.CallOpts, id *big.Int, user common.Address) (*big.Int, error) {
	var out []interface{}
	err := _BBTPledgePool.contract.Call(opts, &out, "withdrawableAmount", id, user)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// WithdrawableAmount is a free data retrieval call binding the contract method 0x7831a74f.
//
// Solidity: function withdrawableAmount(uint256 id, address user) view returns(uint256)
func (_BBTPledgePool *BBTPledgePoolSession) WithdrawableAmount(id *big.Int, user common.Address) (*big.Int, error) {
	return _BBTPledgePool.Contract.WithdrawableAmount(&_BBTPledgePool.CallOpts, id, user)
}

// WithdrawableAmount is a free data retrieval call binding the contract method 0x7831a74f.
//
// Solidity: function withdrawableAmount(uint256 id, address user) view returns(uint256)
func (_BBTPledgePool *BBTPledgePoolCallerSession) WithdrawableAmount(id *big.Int, user common.Address) (*big.Int, error) {
	return _BBTPledgePool.Contract.WithdrawableAmount(&_BBTPledgePool.CallOpts, id, user)
}

// AdminWithdraw is a paid mutator transaction binding the contract method 0x7c5b4a37.
//
// Solidity: function adminWithdraw(uint256 amount) returns()
func (_BBTPledgePool *BBTPledgePoolTransactor) AdminWithdraw(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _BBTPledgePool.contract.Transact(opts, "adminWithdraw", amount)
}

// AdminWithdraw is a paid mutator transaction binding the contract method 0x7c5b4a37.
//
// Solidity: function adminWithdraw(uint256 amount) returns()
func (_BBTPledgePool *BBTPledgePoolSession) AdminWithdraw(amount *big.Int) (*types.Transaction, error) {
	return _BBTPledgePool.Contract.AdminWithdraw(&_BBTPledgePool.TransactOpts, amount)
}

// AdminWithdraw is a paid mutator transaction binding the contract method 0x7c5b4a37.
//
// Solidity: function adminWithdraw(uint256 amount) returns()
func (_BBTPledgePool *BBTPledgePoolTransactorSession) AdminWithdraw(amount *big.Int) (*types.Transaction, error) {
	return _BBTPledgePool.Contract.AdminWithdraw(&_BBTPledgePool.TransactOpts, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address _spender, uint256 _amount) returns()
func (_BBTPledgePool *BBTPledgePoolTransactor) Approve(opts *bind.TransactOpts, _spender common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _BBTPledgePool.contract.Transact(opts, "approve", _spender, _amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address _spender, uint256 _amount) returns()
func (_BBTPledgePool *BBTPledgePoolSession) Approve(_spender common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _BBTPledgePool.Contract.Approve(&_BBTPledgePool.TransactOpts, _spender, _amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address _spender, uint256 _amount) returns()
func (_BBTPledgePool *BBTPledgePoolTransactorSession) Approve(_spender common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _BBTPledgePool.Contract.Approve(&_BBTPledgePool.TransactOpts, _spender, _amount)
}

// ExtendStakingPeriod is a paid mutator transaction binding the contract method 0xabdae543.
//
// Solidity: function extendStakingPeriod(uint256 id, uint256 stakingPeriod) returns()
func (_BBTPledgePool *BBTPledgePoolTransactor) ExtendStakingPeriod(opts *bind.TransactOpts, id *big.Int, stakingPeriod *big.Int) (*types.Transaction, error) {
	return _BBTPledgePool.contract.Transact(opts, "extendStakingPeriod", id, stakingPeriod)
}

// ExtendStakingPeriod is a paid mutator transaction binding the contract method 0xabdae543.
//
// Solidity: function extendStakingPeriod(uint256 id, uint256 stakingPeriod) returns()
func (_BBTPledgePool *BBTPledgePoolSession) ExtendStakingPeriod(id *big.Int, stakingPeriod *big.Int) (*types.Transaction, error) {
	return _BBTPledgePool.Contract.ExtendStakingPeriod(&_BBTPledgePool.TransactOpts, id, stakingPeriod)
}

// ExtendStakingPeriod is a paid mutator transaction binding the contract method 0xabdae543.
//
// Solidity: function extendStakingPeriod(uint256 id, uint256 stakingPeriod) returns()
func (_BBTPledgePool *BBTPledgePoolTransactorSession) ExtendStakingPeriod(id *big.Int, stakingPeriod *big.Int) (*types.Transaction, error) {
	return _BBTPledgePool.Contract.ExtendStakingPeriod(&_BBTPledgePool.TransactOpts, id, stakingPeriod)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_BBTPledgePool *BBTPledgePoolTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BBTPledgePool.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_BBTPledgePool *BBTPledgePoolSession) RenounceOwnership() (*types.Transaction, error) {
	return _BBTPledgePool.Contract.RenounceOwnership(&_BBTPledgePool.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_BBTPledgePool *BBTPledgePoolTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _BBTPledgePool.Contract.RenounceOwnership(&_BBTPledgePool.TransactOpts)
}

// SetMinimumRate is a paid mutator transaction binding the contract method 0x25ac0992.
//
// Solidity: function setMinimumRate(uint64 _minimumRate) returns()
func (_BBTPledgePool *BBTPledgePoolTransactor) SetMinimumRate(opts *bind.TransactOpts, _minimumRate uint64) (*types.Transaction, error) {
	return _BBTPledgePool.contract.Transact(opts, "setMinimumRate", _minimumRate)
}

// SetMinimumRate is a paid mutator transaction binding the contract method 0x25ac0992.
//
// Solidity: function setMinimumRate(uint64 _minimumRate) returns()
func (_BBTPledgePool *BBTPledgePoolSession) SetMinimumRate(_minimumRate uint64) (*types.Transaction, error) {
	return _BBTPledgePool.Contract.SetMinimumRate(&_BBTPledgePool.TransactOpts, _minimumRate)
}

// SetMinimumRate is a paid mutator transaction binding the contract method 0x25ac0992.
//
// Solidity: function setMinimumRate(uint64 _minimumRate) returns()
func (_BBTPledgePool *BBTPledgePoolTransactorSession) SetMinimumRate(_minimumRate uint64) (*types.Transaction, error) {
	return _BBTPledgePool.Contract.SetMinimumRate(&_BBTPledgePool.TransactOpts, _minimumRate)
}

// SetRewardsRate1 is a paid mutator transaction binding the contract method 0x5d4076ee.
//
// Solidity: function setRewardsRate1(uint64 _rewardsRate1) returns()
func (_BBTPledgePool *BBTPledgePoolTransactor) SetRewardsRate1(opts *bind.TransactOpts, _rewardsRate1 uint64) (*types.Transaction, error) {
	return _BBTPledgePool.contract.Transact(opts, "setRewardsRate1", _rewardsRate1)
}

// SetRewardsRate1 is a paid mutator transaction binding the contract method 0x5d4076ee.
//
// Solidity: function setRewardsRate1(uint64 _rewardsRate1) returns()
func (_BBTPledgePool *BBTPledgePoolSession) SetRewardsRate1(_rewardsRate1 uint64) (*types.Transaction, error) {
	return _BBTPledgePool.Contract.SetRewardsRate1(&_BBTPledgePool.TransactOpts, _rewardsRate1)
}

// SetRewardsRate1 is a paid mutator transaction binding the contract method 0x5d4076ee.
//
// Solidity: function setRewardsRate1(uint64 _rewardsRate1) returns()
func (_BBTPledgePool *BBTPledgePoolTransactorSession) SetRewardsRate1(_rewardsRate1 uint64) (*types.Transaction, error) {
	return _BBTPledgePool.Contract.SetRewardsRate1(&_BBTPledgePool.TransactOpts, _rewardsRate1)
}

// SetRewardsRate2 is a paid mutator transaction binding the contract method 0xf8aac963.
//
// Solidity: function setRewardsRate2(uint64 _rewardsRate2) returns()
func (_BBTPledgePool *BBTPledgePoolTransactor) SetRewardsRate2(opts *bind.TransactOpts, _rewardsRate2 uint64) (*types.Transaction, error) {
	return _BBTPledgePool.contract.Transact(opts, "setRewardsRate2", _rewardsRate2)
}

// SetRewardsRate2 is a paid mutator transaction binding the contract method 0xf8aac963.
//
// Solidity: function setRewardsRate2(uint64 _rewardsRate2) returns()
func (_BBTPledgePool *BBTPledgePoolSession) SetRewardsRate2(_rewardsRate2 uint64) (*types.Transaction, error) {
	return _BBTPledgePool.Contract.SetRewardsRate2(&_BBTPledgePool.TransactOpts, _rewardsRate2)
}

// SetRewardsRate2 is a paid mutator transaction binding the contract method 0xf8aac963.
//
// Solidity: function setRewardsRate2(uint64 _rewardsRate2) returns()
func (_BBTPledgePool *BBTPledgePoolTransactorSession) SetRewardsRate2(_rewardsRate2 uint64) (*types.Transaction, error) {
	return _BBTPledgePool.Contract.SetRewardsRate2(&_BBTPledgePool.TransactOpts, _rewardsRate2)
}

// SetRewardsRate3 is a paid mutator transaction binding the contract method 0x13575261.
//
// Solidity: function setRewardsRate3(uint64 _rewardsRate3) returns()
func (_BBTPledgePool *BBTPledgePoolTransactor) SetRewardsRate3(opts *bind.TransactOpts, _rewardsRate3 uint64) (*types.Transaction, error) {
	return _BBTPledgePool.contract.Transact(opts, "setRewardsRate3", _rewardsRate3)
}

// SetRewardsRate3 is a paid mutator transaction binding the contract method 0x13575261.
//
// Solidity: function setRewardsRate3(uint64 _rewardsRate3) returns()
func (_BBTPledgePool *BBTPledgePoolSession) SetRewardsRate3(_rewardsRate3 uint64) (*types.Transaction, error) {
	return _BBTPledgePool.Contract.SetRewardsRate3(&_BBTPledgePool.TransactOpts, _rewardsRate3)
}

// SetRewardsRate3 is a paid mutator transaction binding the contract method 0x13575261.
//
// Solidity: function setRewardsRate3(uint64 _rewardsRate3) returns()
func (_BBTPledgePool *BBTPledgePoolTransactorSession) SetRewardsRate3(_rewardsRate3 uint64) (*types.Transaction, error) {
	return _BBTPledgePool.Contract.SetRewardsRate3(&_BBTPledgePool.TransactOpts, _rewardsRate3)
}

// Stake is a paid mutator transaction binding the contract method 0x7b0472f0.
//
// Solidity: function stake(uint256 amount, uint256 stakingPeriod) returns()
func (_BBTPledgePool *BBTPledgePoolTransactor) Stake(opts *bind.TransactOpts, amount *big.Int, stakingPeriod *big.Int) (*types.Transaction, error) {
	return _BBTPledgePool.contract.Transact(opts, "stake", amount, stakingPeriod)
}

// Stake is a paid mutator transaction binding the contract method 0x7b0472f0.
//
// Solidity: function stake(uint256 amount, uint256 stakingPeriod) returns()
func (_BBTPledgePool *BBTPledgePoolSession) Stake(amount *big.Int, stakingPeriod *big.Int) (*types.Transaction, error) {
	return _BBTPledgePool.Contract.Stake(&_BBTPledgePool.TransactOpts, amount, stakingPeriod)
}

// Stake is a paid mutator transaction binding the contract method 0x7b0472f0.
//
// Solidity: function stake(uint256 amount, uint256 stakingPeriod) returns()
func (_BBTPledgePool *BBTPledgePoolTransactorSession) Stake(amount *big.Int, stakingPeriod *big.Int) (*types.Transaction, error) {
	return _BBTPledgePool.Contract.Stake(&_BBTPledgePool.TransactOpts, amount, stakingPeriod)
}

// StakeBatch is a paid mutator transaction binding the contract method 0xd0f7d24f.
//
// Solidity: function stakeBatch(uint256[] amounts, uint256[] stakingPeriods) returns()
func (_BBTPledgePool *BBTPledgePoolTransactor) StakeBatch(opts *bind.TransactOpts, amounts []*big.Int, stakingPeriods []*big.Int) (*types.Transaction, error) {
	return _BBTPledgePool.contract.Transact(opts, "stakeBatch", amounts, stakingPeriods)
}

// StakeBatch is a paid mutator transaction binding the contract method 0xd0f7d24f.
//
// Solidity: function stakeBatch(uint256[] amounts, uint256[] stakingPeriods) returns()
func (_BBTPledgePool *BBTPledgePoolSession) StakeBatch(amounts []*big.Int, stakingPeriods []*big.Int) (*types.Transaction, error) {
	return _BBTPledgePool.Contract.StakeBatch(&_BBTPledgePool.TransactOpts, amounts, stakingPeriods)
}

// StakeBatch is a paid mutator transaction binding the contract method 0xd0f7d24f.
//
// Solidity: function stakeBatch(uint256[] amounts, uint256[] stakingPeriods) returns()
func (_BBTPledgePool *BBTPledgePoolTransactorSession) StakeBatch(amounts []*big.Int, stakingPeriods []*big.Int) (*types.Transaction, error) {
	return _BBTPledgePool.Contract.StakeBatch(&_BBTPledgePool.TransactOpts, amounts, stakingPeriods)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_BBTPledgePool *BBTPledgePoolTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _BBTPledgePool.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_BBTPledgePool *BBTPledgePoolSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _BBTPledgePool.Contract.TransferOwnership(&_BBTPledgePool.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_BBTPledgePool *BBTPledgePoolTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _BBTPledgePool.Contract.TransferOwnership(&_BBTPledgePool.TransactOpts, newOwner)
}

// WithdrawAll is a paid mutator transaction binding the contract method 0x958e2d31.
//
// Solidity: function withdrawAll(uint256 id) returns()
func (_BBTPledgePool *BBTPledgePoolTransactor) WithdrawAll(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _BBTPledgePool.contract.Transact(opts, "withdrawAll", id)
}

// WithdrawAll is a paid mutator transaction binding the contract method 0x958e2d31.
//
// Solidity: function withdrawAll(uint256 id) returns()
func (_BBTPledgePool *BBTPledgePoolSession) WithdrawAll(id *big.Int) (*types.Transaction, error) {
	return _BBTPledgePool.Contract.WithdrawAll(&_BBTPledgePool.TransactOpts, id)
}

// WithdrawAll is a paid mutator transaction binding the contract method 0x958e2d31.
//
// Solidity: function withdrawAll(uint256 id) returns()
func (_BBTPledgePool *BBTPledgePoolTransactorSession) WithdrawAll(id *big.Int) (*types.Transaction, error) {
	return _BBTPledgePool.Contract.WithdrawAll(&_BBTPledgePool.TransactOpts, id)
}

// WithdrawBatch is a paid mutator transaction binding the contract method 0x5c5fb521.
//
// Solidity: function withdrawBatch(uint256[] ids, uint256[] amounts) returns()
func (_BBTPledgePool *BBTPledgePoolTransactor) WithdrawBatch(opts *bind.TransactOpts, ids []*big.Int, amounts []*big.Int) (*types.Transaction, error) {
	return _BBTPledgePool.contract.Transact(opts, "withdrawBatch", ids, amounts)
}

// WithdrawBatch is a paid mutator transaction binding the contract method 0x5c5fb521.
//
// Solidity: function withdrawBatch(uint256[] ids, uint256[] amounts) returns()
func (_BBTPledgePool *BBTPledgePoolSession) WithdrawBatch(ids []*big.Int, amounts []*big.Int) (*types.Transaction, error) {
	return _BBTPledgePool.Contract.WithdrawBatch(&_BBTPledgePool.TransactOpts, ids, amounts)
}

// WithdrawBatch is a paid mutator transaction binding the contract method 0x5c5fb521.
//
// Solidity: function withdrawBatch(uint256[] ids, uint256[] amounts) returns()
func (_BBTPledgePool *BBTPledgePoolTransactorSession) WithdrawBatch(ids []*big.Int, amounts []*big.Int) (*types.Transaction, error) {
	return _BBTPledgePool.Contract.WithdrawBatch(&_BBTPledgePool.TransactOpts, ids, amounts)
}

// BBTPledgePoolOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the BBTPledgePool contract.
type BBTPledgePoolOwnershipTransferredIterator struct {
	Event *BBTPledgePoolOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *BBTPledgePoolOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BBTPledgePoolOwnershipTransferred)
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
		it.Event = new(BBTPledgePoolOwnershipTransferred)
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
func (it *BBTPledgePoolOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BBTPledgePoolOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BBTPledgePoolOwnershipTransferred represents a OwnershipTransferred event raised by the BBTPledgePool contract.
type BBTPledgePoolOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_BBTPledgePool *BBTPledgePoolFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*BBTPledgePoolOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _BBTPledgePool.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &BBTPledgePoolOwnershipTransferredIterator{contract: _BBTPledgePool.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_BBTPledgePool *BBTPledgePoolFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *BBTPledgePoolOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _BBTPledgePool.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BBTPledgePoolOwnershipTransferred)
				if err := _BBTPledgePool.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_BBTPledgePool *BBTPledgePoolFilterer) ParseOwnershipTransferred(log types.Log) (*BBTPledgePoolOwnershipTransferred, error) {
	event := new(BBTPledgePoolOwnershipTransferred)
	if err := _BBTPledgePool.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BBTPledgePoolStakeIterator is returned from FilterStake and is used to iterate over the raw logs and unpacked data for Stake events raised by the BBTPledgePool contract.
type BBTPledgePoolStakeIterator struct {
	Event *BBTPledgePoolStake // Event containing the contract specifics and raw log

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
func (it *BBTPledgePoolStakeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BBTPledgePoolStake)
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
		it.Event = new(BBTPledgePoolStake)
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
func (it *BBTPledgePoolStakeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BBTPledgePoolStakeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BBTPledgePoolStake represents a Stake event raised by the BBTPledgePool contract.
type BBTPledgePoolStake struct {
	User          common.Address
	Id            *big.Int
	StakedAmount  *big.Int
	StakingPeriod *big.Int
	StartTime     *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterStake is a free log retrieval operation binding the contract event 0x2720efa4b2dd4f3f8a347da3cbd290a522e9432da9072c5b8e6300496fdde282.
//
// Solidity: event Stake(address indexed user, uint256 indexed id, uint256 indexed stakedAmount, uint256 stakingPeriod, uint256 startTime)
func (_BBTPledgePool *BBTPledgePoolFilterer) FilterStake(opts *bind.FilterOpts, user []common.Address, id []*big.Int, stakedAmount []*big.Int) (*BBTPledgePoolStakeIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var stakedAmountRule []interface{}
	for _, stakedAmountItem := range stakedAmount {
		stakedAmountRule = append(stakedAmountRule, stakedAmountItem)
	}

	logs, sub, err := _BBTPledgePool.contract.FilterLogs(opts, "Stake", userRule, idRule, stakedAmountRule)
	if err != nil {
		return nil, err
	}
	return &BBTPledgePoolStakeIterator{contract: _BBTPledgePool.contract, event: "Stake", logs: logs, sub: sub}, nil
}

// WatchStake is a free log subscription operation binding the contract event 0x2720efa4b2dd4f3f8a347da3cbd290a522e9432da9072c5b8e6300496fdde282.
//
// Solidity: event Stake(address indexed user, uint256 indexed id, uint256 indexed stakedAmount, uint256 stakingPeriod, uint256 startTime)
func (_BBTPledgePool *BBTPledgePoolFilterer) WatchStake(opts *bind.WatchOpts, sink chan<- *BBTPledgePoolStake, user []common.Address, id []*big.Int, stakedAmount []*big.Int) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var stakedAmountRule []interface{}
	for _, stakedAmountItem := range stakedAmount {
		stakedAmountRule = append(stakedAmountRule, stakedAmountItem)
	}

	logs, sub, err := _BBTPledgePool.contract.WatchLogs(opts, "Stake", userRule, idRule, stakedAmountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BBTPledgePoolStake)
				if err := _BBTPledgePool.contract.UnpackLog(event, "Stake", log); err != nil {
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

// ParseStake is a log parse operation binding the contract event 0x2720efa4b2dd4f3f8a347da3cbd290a522e9432da9072c5b8e6300496fdde282.
//
// Solidity: event Stake(address indexed user, uint256 indexed id, uint256 indexed stakedAmount, uint256 stakingPeriod, uint256 startTime)
func (_BBTPledgePool *BBTPledgePoolFilterer) ParseStake(log types.Log) (*BBTPledgePoolStake, error) {
	event := new(BBTPledgePoolStake)
	if err := _BBTPledgePool.contract.UnpackLog(event, "Stake", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BBTPledgePoolWithdrawIterator is returned from FilterWithdraw and is used to iterate over the raw logs and unpacked data for Withdraw events raised by the BBTPledgePool contract.
type BBTPledgePoolWithdrawIterator struct {
	Event *BBTPledgePoolWithdraw // Event containing the contract specifics and raw log

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
func (it *BBTPledgePoolWithdrawIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BBTPledgePoolWithdraw)
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
		it.Event = new(BBTPledgePoolWithdraw)
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
func (it *BBTPledgePoolWithdrawIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BBTPledgePoolWithdrawIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BBTPledgePoolWithdraw represents a Withdraw event raised by the BBTPledgePool contract.
type BBTPledgePoolWithdraw struct {
	User           common.Address
	Id             *big.Int
	StakedAmount   *big.Int
	WithdrawAmount *big.Int
	StakingPeriod  *big.Int
	StartTime      *big.Int
	EndTime        *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterWithdraw is a free log retrieval operation binding the contract event 0x194c8d0132d20112211dfa71bb87a92766fde4f4318e08efd2cc4a6e188e509d.
//
// Solidity: event Withdraw(address indexed user, uint256 indexed id, uint256 stakedAmount, uint256 indexed withdrawAmount, uint256 stakingPeriod, uint256 startTime, uint256 endTime)
func (_BBTPledgePool *BBTPledgePoolFilterer) FilterWithdraw(opts *bind.FilterOpts, user []common.Address, id []*big.Int, withdrawAmount []*big.Int) (*BBTPledgePoolWithdrawIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	var withdrawAmountRule []interface{}
	for _, withdrawAmountItem := range withdrawAmount {
		withdrawAmountRule = append(withdrawAmountRule, withdrawAmountItem)
	}

	logs, sub, err := _BBTPledgePool.contract.FilterLogs(opts, "Withdraw", userRule, idRule, withdrawAmountRule)
	if err != nil {
		return nil, err
	}
	return &BBTPledgePoolWithdrawIterator{contract: _BBTPledgePool.contract, event: "Withdraw", logs: logs, sub: sub}, nil
}

// WatchWithdraw is a free log subscription operation binding the contract event 0x194c8d0132d20112211dfa71bb87a92766fde4f4318e08efd2cc4a6e188e509d.
//
// Solidity: event Withdraw(address indexed user, uint256 indexed id, uint256 stakedAmount, uint256 indexed withdrawAmount, uint256 stakingPeriod, uint256 startTime, uint256 endTime)
func (_BBTPledgePool *BBTPledgePoolFilterer) WatchWithdraw(opts *bind.WatchOpts, sink chan<- *BBTPledgePoolWithdraw, user []common.Address, id []*big.Int, withdrawAmount []*big.Int) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	var withdrawAmountRule []interface{}
	for _, withdrawAmountItem := range withdrawAmount {
		withdrawAmountRule = append(withdrawAmountRule, withdrawAmountItem)
	}

	logs, sub, err := _BBTPledgePool.contract.WatchLogs(opts, "Withdraw", userRule, idRule, withdrawAmountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BBTPledgePoolWithdraw)
				if err := _BBTPledgePool.contract.UnpackLog(event, "Withdraw", log); err != nil {
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

// ParseWithdraw is a log parse operation binding the contract event 0x194c8d0132d20112211dfa71bb87a92766fde4f4318e08efd2cc4a6e188e509d.
//
// Solidity: event Withdraw(address indexed user, uint256 indexed id, uint256 stakedAmount, uint256 indexed withdrawAmount, uint256 stakingPeriod, uint256 startTime, uint256 endTime)
func (_BBTPledgePool *BBTPledgePoolFilterer) ParseWithdraw(log types.Log) (*BBTPledgePoolWithdraw, error) {
	event := new(BBTPledgePoolWithdraw)
	if err := _BBTPledgePool.contract.UnpackLog(event, "Withdraw", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
