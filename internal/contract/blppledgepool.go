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

// MiningV2UserInfo is an auto generated low-level Go binding around an user-defined struct.
type MiningV2UserInfo struct {
	Id         *big.Int
	Amount     *big.Int
	RewardDebt *big.Int
	LockPeriod *big.Int
	StartTime  *big.Int
	Power      *big.Int
}

// BLPPledgePoolMetaData contains all meta data concerning the BLPPledgePool contract.
var BLPPledgePoolMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"_bbt\",\"type\":\"address\"},{\"internalType\":\"contractIERC20\",\"name\":\"_bbtlp\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_startBlock\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"pid\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"lockPeriod\",\"type\":\"uint256\"}],\"name\":\"Deposit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"pid\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"EmergencyWithdraw\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"pid\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Harvest\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"pid\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Withdraw\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"BONUS_MULTIPLIER\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_allocPoint\",\"type\":\"uint256\"},{\"internalType\":\"contractIERC20\",\"name\":\"_lpToken\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"_powerRatio1\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"_powerRatio2\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"_powerRatio3\",\"type\":\"uint8\"},{\"internalType\":\"bool\",\"name\":\"_withUpdate\",\"type\":\"bool\"}],\"name\":\"add\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"bbt\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_pid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"lockPeriod\",\"type\":\"uint256\"}],\"name\":\"deposit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_pid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_user\",\"type\":\"address\"}],\"name\":\"emergencyWithdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"erc20Recover\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_from\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_to\",\"type\":\"uint256\"}],\"name\":\"getMultiplier\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_pid\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"getUserInfos\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"rewardDebt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"lockPeriod\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"startTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"power\",\"type\":\"uint256\"}],\"internalType\":\"structMiningV2.UserInfo[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_pid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"isBatch\",\"type\":\"bool\"}],\"name\":\"harvest\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"harvests\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"idCounter\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"massUpdatePools\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_pid\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_user\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"pendingBBT\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"poolInfo\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"lpToken\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"allocPoint\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"lastRewardBlock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"accBBTPerShare\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalStakedPower\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalLocked\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"powerRatio1\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"powerRatio2\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"powerRatio3\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"poolLength\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_pid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_allocPoint\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"_withUpdate\",\"type\":\"bool\"}],\"name\":\"set\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"setLastBalance\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"startBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalAllocPoint\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"multiplierNumber\",\"type\":\"uint256\"}],\"name\":\"updateMultiplier\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_pid\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"_isBatch\",\"type\":\"bool\"}],\"name\":\"updatePool\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"userInfo\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"rewardDebt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"lockPeriod\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"startTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"power\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_pid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// BLPPledgePoolABI is the input ABI used to generate the binding from.
// Deprecated: Use BLPPledgePoolMetaData.ABI instead.
var BLPPledgePoolABI = BLPPledgePoolMetaData.ABI

// BLPPledgePool is an auto generated Go binding around an Ethereum contract.
type BLPPledgePool struct {
	BLPPledgePoolCaller     // Read-only binding to the contract
	BLPPledgePoolTransactor // Write-only binding to the contract
	BLPPledgePoolFilterer   // Log filterer for contract events
}

// BLPPledgePoolCaller is an auto generated read-only Go binding around an Ethereum contract.
type BLPPledgePoolCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BLPPledgePoolTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BLPPledgePoolTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BLPPledgePoolFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BLPPledgePoolFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BLPPledgePoolSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BLPPledgePoolSession struct {
	Contract     *BLPPledgePool    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BLPPledgePoolCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BLPPledgePoolCallerSession struct {
	Contract *BLPPledgePoolCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// BLPPledgePoolTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BLPPledgePoolTransactorSession struct {
	Contract     *BLPPledgePoolTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// BLPPledgePoolRaw is an auto generated low-level Go binding around an Ethereum contract.
type BLPPledgePoolRaw struct {
	Contract *BLPPledgePool // Generic contract binding to access the raw methods on
}

// BLPPledgePoolCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BLPPledgePoolCallerRaw struct {
	Contract *BLPPledgePoolCaller // Generic read-only contract binding to access the raw methods on
}

// BLPPledgePoolTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BLPPledgePoolTransactorRaw struct {
	Contract *BLPPledgePoolTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBLPPledgePool creates a new instance of BLPPledgePool, bound to a specific deployed contract.
func NewBLPPledgePool(address common.Address, backend bind.ContractBackend) (*BLPPledgePool, error) {
	contract, err := bindBLPPledgePool(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BLPPledgePool{BLPPledgePoolCaller: BLPPledgePoolCaller{contract: contract}, BLPPledgePoolTransactor: BLPPledgePoolTransactor{contract: contract}, BLPPledgePoolFilterer: BLPPledgePoolFilterer{contract: contract}}, nil
}

// NewBLPPledgePoolCaller creates a new read-only instance of BLPPledgePool, bound to a specific deployed contract.
func NewBLPPledgePoolCaller(address common.Address, caller bind.ContractCaller) (*BLPPledgePoolCaller, error) {
	contract, err := bindBLPPledgePool(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BLPPledgePoolCaller{contract: contract}, nil
}

// NewBLPPledgePoolTransactor creates a new write-only instance of BLPPledgePool, bound to a specific deployed contract.
func NewBLPPledgePoolTransactor(address common.Address, transactor bind.ContractTransactor) (*BLPPledgePoolTransactor, error) {
	contract, err := bindBLPPledgePool(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BLPPledgePoolTransactor{contract: contract}, nil
}

// NewBLPPledgePoolFilterer creates a new log filterer instance of BLPPledgePool, bound to a specific deployed contract.
func NewBLPPledgePoolFilterer(address common.Address, filterer bind.ContractFilterer) (*BLPPledgePoolFilterer, error) {
	contract, err := bindBLPPledgePool(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BLPPledgePoolFilterer{contract: contract}, nil
}

// bindBLPPledgePool binds a generic wrapper to an already deployed contract.
func bindBLPPledgePool(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := BLPPledgePoolMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BLPPledgePool *BLPPledgePoolRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BLPPledgePool.Contract.BLPPledgePoolCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BLPPledgePool *BLPPledgePoolRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BLPPledgePool.Contract.BLPPledgePoolTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BLPPledgePool *BLPPledgePoolRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BLPPledgePool.Contract.BLPPledgePoolTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BLPPledgePool *BLPPledgePoolCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BLPPledgePool.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BLPPledgePool *BLPPledgePoolTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BLPPledgePool.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BLPPledgePool *BLPPledgePoolTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BLPPledgePool.Contract.contract.Transact(opts, method, params...)
}

// BONUSMULTIPLIER is a free data retrieval call binding the contract method 0x8aa28550.
//
// Solidity: function BONUS_MULTIPLIER() view returns(uint256)
func (_BLPPledgePool *BLPPledgePoolCaller) BONUSMULTIPLIER(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _BLPPledgePool.contract.Call(opts, &out, "BONUS_MULTIPLIER")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BONUSMULTIPLIER is a free data retrieval call binding the contract method 0x8aa28550.
//
// Solidity: function BONUS_MULTIPLIER() view returns(uint256)
func (_BLPPledgePool *BLPPledgePoolSession) BONUSMULTIPLIER() (*big.Int, error) {
	return _BLPPledgePool.Contract.BONUSMULTIPLIER(&_BLPPledgePool.CallOpts)
}

// BONUSMULTIPLIER is a free data retrieval call binding the contract method 0x8aa28550.
//
// Solidity: function BONUS_MULTIPLIER() view returns(uint256)
func (_BLPPledgePool *BLPPledgePoolCallerSession) BONUSMULTIPLIER() (*big.Int, error) {
	return _BLPPledgePool.Contract.BONUSMULTIPLIER(&_BLPPledgePool.CallOpts)
}

// Bbt is a free data retrieval call binding the contract method 0x72927452.
//
// Solidity: function bbt() view returns(address)
func (_BLPPledgePool *BLPPledgePoolCaller) Bbt(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BLPPledgePool.contract.Call(opts, &out, "bbt")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Bbt is a free data retrieval call binding the contract method 0x72927452.
//
// Solidity: function bbt() view returns(address)
func (_BLPPledgePool *BLPPledgePoolSession) Bbt() (common.Address, error) {
	return _BLPPledgePool.Contract.Bbt(&_BLPPledgePool.CallOpts)
}

// Bbt is a free data retrieval call binding the contract method 0x72927452.
//
// Solidity: function bbt() view returns(address)
func (_BLPPledgePool *BLPPledgePoolCallerSession) Bbt() (common.Address, error) {
	return _BLPPledgePool.Contract.Bbt(&_BLPPledgePool.CallOpts)
}

// GetMultiplier is a free data retrieval call binding the contract method 0x8dbb1e3a.
//
// Solidity: function getMultiplier(uint256 _from, uint256 _to) view returns(uint256)
func (_BLPPledgePool *BLPPledgePoolCaller) GetMultiplier(opts *bind.CallOpts, _from *big.Int, _to *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _BLPPledgePool.contract.Call(opts, &out, "getMultiplier", _from, _to)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetMultiplier is a free data retrieval call binding the contract method 0x8dbb1e3a.
//
// Solidity: function getMultiplier(uint256 _from, uint256 _to) view returns(uint256)
func (_BLPPledgePool *BLPPledgePoolSession) GetMultiplier(_from *big.Int, _to *big.Int) (*big.Int, error) {
	return _BLPPledgePool.Contract.GetMultiplier(&_BLPPledgePool.CallOpts, _from, _to)
}

// GetMultiplier is a free data retrieval call binding the contract method 0x8dbb1e3a.
//
// Solidity: function getMultiplier(uint256 _from, uint256 _to) view returns(uint256)
func (_BLPPledgePool *BLPPledgePoolCallerSession) GetMultiplier(_from *big.Int, _to *big.Int) (*big.Int, error) {
	return _BLPPledgePool.Contract.GetMultiplier(&_BLPPledgePool.CallOpts, _from, _to)
}

// GetUserInfos is a free data retrieval call binding the contract method 0x0900f44b.
//
// Solidity: function getUserInfos(uint256 _pid, address _addr) view returns((uint256,uint256,uint256,uint256,uint256,uint256)[])
func (_BLPPledgePool *BLPPledgePoolCaller) GetUserInfos(opts *bind.CallOpts, _pid *big.Int, _addr common.Address) ([]MiningV2UserInfo, error) {
	var out []interface{}
	err := _BLPPledgePool.contract.Call(opts, &out, "getUserInfos", _pid, _addr)

	if err != nil {
		return *new([]MiningV2UserInfo), err
	}

	out0 := *abi.ConvertType(out[0], new([]MiningV2UserInfo)).(*[]MiningV2UserInfo)

	return out0, err

}

// GetUserInfos is a free data retrieval call binding the contract method 0x0900f44b.
//
// Solidity: function getUserInfos(uint256 _pid, address _addr) view returns((uint256,uint256,uint256,uint256,uint256,uint256)[])
func (_BLPPledgePool *BLPPledgePoolSession) GetUserInfos(_pid *big.Int, _addr common.Address) ([]MiningV2UserInfo, error) {
	return _BLPPledgePool.Contract.GetUserInfos(&_BLPPledgePool.CallOpts, _pid, _addr)
}

// GetUserInfos is a free data retrieval call binding the contract method 0x0900f44b.
//
// Solidity: function getUserInfos(uint256 _pid, address _addr) view returns((uint256,uint256,uint256,uint256,uint256,uint256)[])
func (_BLPPledgePool *BLPPledgePoolCallerSession) GetUserInfos(_pid *big.Int, _addr common.Address) ([]MiningV2UserInfo, error) {
	return _BLPPledgePool.Contract.GetUserInfos(&_BLPPledgePool.CallOpts, _pid, _addr)
}

// IdCounter is a free data retrieval call binding the contract method 0xeb08ab28.
//
// Solidity: function idCounter() view returns(uint256)
func (_BLPPledgePool *BLPPledgePoolCaller) IdCounter(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _BLPPledgePool.contract.Call(opts, &out, "idCounter")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// IdCounter is a free data retrieval call binding the contract method 0xeb08ab28.
//
// Solidity: function idCounter() view returns(uint256)
func (_BLPPledgePool *BLPPledgePoolSession) IdCounter() (*big.Int, error) {
	return _BLPPledgePool.Contract.IdCounter(&_BLPPledgePool.CallOpts)
}

// IdCounter is a free data retrieval call binding the contract method 0xeb08ab28.
//
// Solidity: function idCounter() view returns(uint256)
func (_BLPPledgePool *BLPPledgePoolCallerSession) IdCounter() (*big.Int, error) {
	return _BLPPledgePool.Contract.IdCounter(&_BLPPledgePool.CallOpts)
}

// LastBalance is a free data retrieval call binding the contract method 0x8f1c56bd.
//
// Solidity: function lastBalance() view returns(uint256)
func (_BLPPledgePool *BLPPledgePoolCaller) LastBalance(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _BLPPledgePool.contract.Call(opts, &out, "lastBalance")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LastBalance is a free data retrieval call binding the contract method 0x8f1c56bd.
//
// Solidity: function lastBalance() view returns(uint256)
func (_BLPPledgePool *BLPPledgePoolSession) LastBalance() (*big.Int, error) {
	return _BLPPledgePool.Contract.LastBalance(&_BLPPledgePool.CallOpts)
}

// LastBalance is a free data retrieval call binding the contract method 0x8f1c56bd.
//
// Solidity: function lastBalance() view returns(uint256)
func (_BLPPledgePool *BLPPledgePoolCallerSession) LastBalance() (*big.Int, error) {
	return _BLPPledgePool.Contract.LastBalance(&_BLPPledgePool.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_BLPPledgePool *BLPPledgePoolCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BLPPledgePool.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_BLPPledgePool *BLPPledgePoolSession) Owner() (common.Address, error) {
	return _BLPPledgePool.Contract.Owner(&_BLPPledgePool.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_BLPPledgePool *BLPPledgePoolCallerSession) Owner() (common.Address, error) {
	return _BLPPledgePool.Contract.Owner(&_BLPPledgePool.CallOpts)
}

// PendingBBT is a free data retrieval call binding the contract method 0x1bc06893.
//
// Solidity: function pendingBBT(uint256 _pid, address _user, uint256 id) view returns(uint256)
func (_BLPPledgePool *BLPPledgePoolCaller) PendingBBT(opts *bind.CallOpts, _pid *big.Int, _user common.Address, id *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _BLPPledgePool.contract.Call(opts, &out, "pendingBBT", _pid, _user, id)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PendingBBT is a free data retrieval call binding the contract method 0x1bc06893.
//
// Solidity: function pendingBBT(uint256 _pid, address _user, uint256 id) view returns(uint256)
func (_BLPPledgePool *BLPPledgePoolSession) PendingBBT(_pid *big.Int, _user common.Address, id *big.Int) (*big.Int, error) {
	return _BLPPledgePool.Contract.PendingBBT(&_BLPPledgePool.CallOpts, _pid, _user, id)
}

// PendingBBT is a free data retrieval call binding the contract method 0x1bc06893.
//
// Solidity: function pendingBBT(uint256 _pid, address _user, uint256 id) view returns(uint256)
func (_BLPPledgePool *BLPPledgePoolCallerSession) PendingBBT(_pid *big.Int, _user common.Address, id *big.Int) (*big.Int, error) {
	return _BLPPledgePool.Contract.PendingBBT(&_BLPPledgePool.CallOpts, _pid, _user, id)
}

// PoolInfo is a free data retrieval call binding the contract method 0x1526fe27.
//
// Solidity: function poolInfo(uint256 ) view returns(address lpToken, uint256 allocPoint, uint256 lastRewardBlock, uint256 accBBTPerShare, uint256 totalStakedPower, uint256 totalLocked, uint8 powerRatio1, uint8 powerRatio2, uint8 powerRatio3)
func (_BLPPledgePool *BLPPledgePoolCaller) PoolInfo(opts *bind.CallOpts, arg0 *big.Int) (struct {
	LpToken          common.Address
	AllocPoint       *big.Int
	LastRewardBlock  *big.Int
	AccBBTPerShare   *big.Int
	TotalStakedPower *big.Int
	TotalLocked      *big.Int
	PowerRatio1      uint8
	PowerRatio2      uint8
	PowerRatio3      uint8
}, error) {
	var out []interface{}
	err := _BLPPledgePool.contract.Call(opts, &out, "poolInfo", arg0)

	outstruct := new(struct {
		LpToken          common.Address
		AllocPoint       *big.Int
		LastRewardBlock  *big.Int
		AccBBTPerShare   *big.Int
		TotalStakedPower *big.Int
		TotalLocked      *big.Int
		PowerRatio1      uint8
		PowerRatio2      uint8
		PowerRatio3      uint8
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.LpToken = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.AllocPoint = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.LastRewardBlock = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.AccBBTPerShare = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.TotalStakedPower = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.TotalLocked = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)
	outstruct.PowerRatio1 = *abi.ConvertType(out[6], new(uint8)).(*uint8)
	outstruct.PowerRatio2 = *abi.ConvertType(out[7], new(uint8)).(*uint8)
	outstruct.PowerRatio3 = *abi.ConvertType(out[8], new(uint8)).(*uint8)

	return *outstruct, err

}

// PoolInfo is a free data retrieval call binding the contract method 0x1526fe27.
//
// Solidity: function poolInfo(uint256 ) view returns(address lpToken, uint256 allocPoint, uint256 lastRewardBlock, uint256 accBBTPerShare, uint256 totalStakedPower, uint256 totalLocked, uint8 powerRatio1, uint8 powerRatio2, uint8 powerRatio3)
func (_BLPPledgePool *BLPPledgePoolSession) PoolInfo(arg0 *big.Int) (struct {
	LpToken          common.Address
	AllocPoint       *big.Int
	LastRewardBlock  *big.Int
	AccBBTPerShare   *big.Int
	TotalStakedPower *big.Int
	TotalLocked      *big.Int
	PowerRatio1      uint8
	PowerRatio2      uint8
	PowerRatio3      uint8
}, error) {
	return _BLPPledgePool.Contract.PoolInfo(&_BLPPledgePool.CallOpts, arg0)
}

// PoolInfo is a free data retrieval call binding the contract method 0x1526fe27.
//
// Solidity: function poolInfo(uint256 ) view returns(address lpToken, uint256 allocPoint, uint256 lastRewardBlock, uint256 accBBTPerShare, uint256 totalStakedPower, uint256 totalLocked, uint8 powerRatio1, uint8 powerRatio2, uint8 powerRatio3)
func (_BLPPledgePool *BLPPledgePoolCallerSession) PoolInfo(arg0 *big.Int) (struct {
	LpToken          common.Address
	AllocPoint       *big.Int
	LastRewardBlock  *big.Int
	AccBBTPerShare   *big.Int
	TotalStakedPower *big.Int
	TotalLocked      *big.Int
	PowerRatio1      uint8
	PowerRatio2      uint8
	PowerRatio3      uint8
}, error) {
	return _BLPPledgePool.Contract.PoolInfo(&_BLPPledgePool.CallOpts, arg0)
}

// PoolLength is a free data retrieval call binding the contract method 0x081e3eda.
//
// Solidity: function poolLength() view returns(uint256)
func (_BLPPledgePool *BLPPledgePoolCaller) PoolLength(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _BLPPledgePool.contract.Call(opts, &out, "poolLength")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PoolLength is a free data retrieval call binding the contract method 0x081e3eda.
//
// Solidity: function poolLength() view returns(uint256)
func (_BLPPledgePool *BLPPledgePoolSession) PoolLength() (*big.Int, error) {
	return _BLPPledgePool.Contract.PoolLength(&_BLPPledgePool.CallOpts)
}

// PoolLength is a free data retrieval call binding the contract method 0x081e3eda.
//
// Solidity: function poolLength() view returns(uint256)
func (_BLPPledgePool *BLPPledgePoolCallerSession) PoolLength() (*big.Int, error) {
	return _BLPPledgePool.Contract.PoolLength(&_BLPPledgePool.CallOpts)
}

// StartBlock is a free data retrieval call binding the contract method 0x48cd4cb1.
//
// Solidity: function startBlock() view returns(uint256)
func (_BLPPledgePool *BLPPledgePoolCaller) StartBlock(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _BLPPledgePool.contract.Call(opts, &out, "startBlock")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// StartBlock is a free data retrieval call binding the contract method 0x48cd4cb1.
//
// Solidity: function startBlock() view returns(uint256)
func (_BLPPledgePool *BLPPledgePoolSession) StartBlock() (*big.Int, error) {
	return _BLPPledgePool.Contract.StartBlock(&_BLPPledgePool.CallOpts)
}

// StartBlock is a free data retrieval call binding the contract method 0x48cd4cb1.
//
// Solidity: function startBlock() view returns(uint256)
func (_BLPPledgePool *BLPPledgePoolCallerSession) StartBlock() (*big.Int, error) {
	return _BLPPledgePool.Contract.StartBlock(&_BLPPledgePool.CallOpts)
}

// TotalAllocPoint is a free data retrieval call binding the contract method 0x17caf6f1.
//
// Solidity: function totalAllocPoint() view returns(uint256)
func (_BLPPledgePool *BLPPledgePoolCaller) TotalAllocPoint(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _BLPPledgePool.contract.Call(opts, &out, "totalAllocPoint")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalAllocPoint is a free data retrieval call binding the contract method 0x17caf6f1.
//
// Solidity: function totalAllocPoint() view returns(uint256)
func (_BLPPledgePool *BLPPledgePoolSession) TotalAllocPoint() (*big.Int, error) {
	return _BLPPledgePool.Contract.TotalAllocPoint(&_BLPPledgePool.CallOpts)
}

// TotalAllocPoint is a free data retrieval call binding the contract method 0x17caf6f1.
//
// Solidity: function totalAllocPoint() view returns(uint256)
func (_BLPPledgePool *BLPPledgePoolCallerSession) TotalAllocPoint() (*big.Int, error) {
	return _BLPPledgePool.Contract.TotalAllocPoint(&_BLPPledgePool.CallOpts)
}

// UserInfo is a free data retrieval call binding the contract method 0xdeb019bd.
//
// Solidity: function userInfo(uint256 , address , uint256 ) view returns(uint256 id, uint256 amount, uint256 rewardDebt, uint256 lockPeriod, uint256 startTime, uint256 power)
func (_BLPPledgePool *BLPPledgePoolCaller) UserInfo(opts *bind.CallOpts, arg0 *big.Int, arg1 common.Address, arg2 *big.Int) (struct {
	Id         *big.Int
	Amount     *big.Int
	RewardDebt *big.Int
	LockPeriod *big.Int
	StartTime  *big.Int
	Power      *big.Int
}, error) {
	var out []interface{}
	err := _BLPPledgePool.contract.Call(opts, &out, "userInfo", arg0, arg1, arg2)

	outstruct := new(struct {
		Id         *big.Int
		Amount     *big.Int
		RewardDebt *big.Int
		LockPeriod *big.Int
		StartTime  *big.Int
		Power      *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Id = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Amount = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.RewardDebt = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.LockPeriod = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.StartTime = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.Power = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// UserInfo is a free data retrieval call binding the contract method 0xdeb019bd.
//
// Solidity: function userInfo(uint256 , address , uint256 ) view returns(uint256 id, uint256 amount, uint256 rewardDebt, uint256 lockPeriod, uint256 startTime, uint256 power)
func (_BLPPledgePool *BLPPledgePoolSession) UserInfo(arg0 *big.Int, arg1 common.Address, arg2 *big.Int) (struct {
	Id         *big.Int
	Amount     *big.Int
	RewardDebt *big.Int
	LockPeriod *big.Int
	StartTime  *big.Int
	Power      *big.Int
}, error) {
	return _BLPPledgePool.Contract.UserInfo(&_BLPPledgePool.CallOpts, arg0, arg1, arg2)
}

// UserInfo is a free data retrieval call binding the contract method 0xdeb019bd.
//
// Solidity: function userInfo(uint256 , address , uint256 ) view returns(uint256 id, uint256 amount, uint256 rewardDebt, uint256 lockPeriod, uint256 startTime, uint256 power)
func (_BLPPledgePool *BLPPledgePoolCallerSession) UserInfo(arg0 *big.Int, arg1 common.Address, arg2 *big.Int) (struct {
	Id         *big.Int
	Amount     *big.Int
	RewardDebt *big.Int
	LockPeriod *big.Int
	StartTime  *big.Int
	Power      *big.Int
}, error) {
	return _BLPPledgePool.Contract.UserInfo(&_BLPPledgePool.CallOpts, arg0, arg1, arg2)
}

// Add is a paid mutator transaction binding the contract method 0xf28d37f3.
//
// Solidity: function add(uint256 _allocPoint, address _lpToken, uint8 _powerRatio1, uint8 _powerRatio2, uint8 _powerRatio3, bool _withUpdate) returns()
func (_BLPPledgePool *BLPPledgePoolTransactor) Add(opts *bind.TransactOpts, _allocPoint *big.Int, _lpToken common.Address, _powerRatio1 uint8, _powerRatio2 uint8, _powerRatio3 uint8, _withUpdate bool) (*types.Transaction, error) {
	return _BLPPledgePool.contract.Transact(opts, "add", _allocPoint, _lpToken, _powerRatio1, _powerRatio2, _powerRatio3, _withUpdate)
}

// Add is a paid mutator transaction binding the contract method 0xf28d37f3.
//
// Solidity: function add(uint256 _allocPoint, address _lpToken, uint8 _powerRatio1, uint8 _powerRatio2, uint8 _powerRatio3, bool _withUpdate) returns()
func (_BLPPledgePool *BLPPledgePoolSession) Add(_allocPoint *big.Int, _lpToken common.Address, _powerRatio1 uint8, _powerRatio2 uint8, _powerRatio3 uint8, _withUpdate bool) (*types.Transaction, error) {
	return _BLPPledgePool.Contract.Add(&_BLPPledgePool.TransactOpts, _allocPoint, _lpToken, _powerRatio1, _powerRatio2, _powerRatio3, _withUpdate)
}

// Add is a paid mutator transaction binding the contract method 0xf28d37f3.
//
// Solidity: function add(uint256 _allocPoint, address _lpToken, uint8 _powerRatio1, uint8 _powerRatio2, uint8 _powerRatio3, bool _withUpdate) returns()
func (_BLPPledgePool *BLPPledgePoolTransactorSession) Add(_allocPoint *big.Int, _lpToken common.Address, _powerRatio1 uint8, _powerRatio2 uint8, _powerRatio3 uint8, _withUpdate bool) (*types.Transaction, error) {
	return _BLPPledgePool.Contract.Add(&_BLPPledgePool.TransactOpts, _allocPoint, _lpToken, _powerRatio1, _powerRatio2, _powerRatio3, _withUpdate)
}

// Deposit is a paid mutator transaction binding the contract method 0x00aeef8a.
//
// Solidity: function deposit(uint256 _pid, uint256 _amount, uint256 lockPeriod) returns()
func (_BLPPledgePool *BLPPledgePoolTransactor) Deposit(opts *bind.TransactOpts, _pid *big.Int, _amount *big.Int, lockPeriod *big.Int) (*types.Transaction, error) {
	return _BLPPledgePool.contract.Transact(opts, "deposit", _pid, _amount, lockPeriod)
}

// Deposit is a paid mutator transaction binding the contract method 0x00aeef8a.
//
// Solidity: function deposit(uint256 _pid, uint256 _amount, uint256 lockPeriod) returns()
func (_BLPPledgePool *BLPPledgePoolSession) Deposit(_pid *big.Int, _amount *big.Int, lockPeriod *big.Int) (*types.Transaction, error) {
	return _BLPPledgePool.Contract.Deposit(&_BLPPledgePool.TransactOpts, _pid, _amount, lockPeriod)
}

// Deposit is a paid mutator transaction binding the contract method 0x00aeef8a.
//
// Solidity: function deposit(uint256 _pid, uint256 _amount, uint256 lockPeriod) returns()
func (_BLPPledgePool *BLPPledgePoolTransactorSession) Deposit(_pid *big.Int, _amount *big.Int, lockPeriod *big.Int) (*types.Transaction, error) {
	return _BLPPledgePool.Contract.Deposit(&_BLPPledgePool.TransactOpts, _pid, _amount, lockPeriod)
}

// EmergencyWithdraw is a paid mutator transaction binding the contract method 0xf72b2a4c.
//
// Solidity: function emergencyWithdraw(uint256 _pid, uint256 id, address _user) returns()
func (_BLPPledgePool *BLPPledgePoolTransactor) EmergencyWithdraw(opts *bind.TransactOpts, _pid *big.Int, id *big.Int, _user common.Address) (*types.Transaction, error) {
	return _BLPPledgePool.contract.Transact(opts, "emergencyWithdraw", _pid, id, _user)
}

// EmergencyWithdraw is a paid mutator transaction binding the contract method 0xf72b2a4c.
//
// Solidity: function emergencyWithdraw(uint256 _pid, uint256 id, address _user) returns()
func (_BLPPledgePool *BLPPledgePoolSession) EmergencyWithdraw(_pid *big.Int, id *big.Int, _user common.Address) (*types.Transaction, error) {
	return _BLPPledgePool.Contract.EmergencyWithdraw(&_BLPPledgePool.TransactOpts, _pid, id, _user)
}

// EmergencyWithdraw is a paid mutator transaction binding the contract method 0xf72b2a4c.
//
// Solidity: function emergencyWithdraw(uint256 _pid, uint256 id, address _user) returns()
func (_BLPPledgePool *BLPPledgePoolTransactorSession) EmergencyWithdraw(_pid *big.Int, id *big.Int, _user common.Address) (*types.Transaction, error) {
	return _BLPPledgePool.Contract.EmergencyWithdraw(&_BLPPledgePool.TransactOpts, _pid, id, _user)
}

// Erc20Recover is a paid mutator transaction binding the contract method 0x68477bca.
//
// Solidity: function erc20Recover(address _token, address _to, uint256 _amount) returns()
func (_BLPPledgePool *BLPPledgePoolTransactor) Erc20Recover(opts *bind.TransactOpts, _token common.Address, _to common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _BLPPledgePool.contract.Transact(opts, "erc20Recover", _token, _to, _amount)
}

// Erc20Recover is a paid mutator transaction binding the contract method 0x68477bca.
//
// Solidity: function erc20Recover(address _token, address _to, uint256 _amount) returns()
func (_BLPPledgePool *BLPPledgePoolSession) Erc20Recover(_token common.Address, _to common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _BLPPledgePool.Contract.Erc20Recover(&_BLPPledgePool.TransactOpts, _token, _to, _amount)
}

// Erc20Recover is a paid mutator transaction binding the contract method 0x68477bca.
//
// Solidity: function erc20Recover(address _token, address _to, uint256 _amount) returns()
func (_BLPPledgePool *BLPPledgePoolTransactorSession) Erc20Recover(_token common.Address, _to common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _BLPPledgePool.Contract.Erc20Recover(&_BLPPledgePool.TransactOpts, _token, _to, _amount)
}

// Harvest is a paid mutator transaction binding the contract method 0x89eee90b.
//
// Solidity: function harvest(uint256 _pid, uint256 id, bool isBatch) returns()
func (_BLPPledgePool *BLPPledgePoolTransactor) Harvest(opts *bind.TransactOpts, _pid *big.Int, id *big.Int, isBatch bool) (*types.Transaction, error) {
	return _BLPPledgePool.contract.Transact(opts, "harvest", _pid, id, isBatch)
}

// Harvest is a paid mutator transaction binding the contract method 0x89eee90b.
//
// Solidity: function harvest(uint256 _pid, uint256 id, bool isBatch) returns()
func (_BLPPledgePool *BLPPledgePoolSession) Harvest(_pid *big.Int, id *big.Int, isBatch bool) (*types.Transaction, error) {
	return _BLPPledgePool.Contract.Harvest(&_BLPPledgePool.TransactOpts, _pid, id, isBatch)
}

// Harvest is a paid mutator transaction binding the contract method 0x89eee90b.
//
// Solidity: function harvest(uint256 _pid, uint256 id, bool isBatch) returns()
func (_BLPPledgePool *BLPPledgePoolTransactorSession) Harvest(_pid *big.Int, id *big.Int, isBatch bool) (*types.Transaction, error) {
	return _BLPPledgePool.Contract.Harvest(&_BLPPledgePool.TransactOpts, _pid, id, isBatch)
}

// Harvests is a paid mutator transaction binding the contract method 0x2a9bd994.
//
// Solidity: function harvests() returns()
func (_BLPPledgePool *BLPPledgePoolTransactor) Harvests(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BLPPledgePool.contract.Transact(opts, "harvests")
}

// Harvests is a paid mutator transaction binding the contract method 0x2a9bd994.
//
// Solidity: function harvests() returns()
func (_BLPPledgePool *BLPPledgePoolSession) Harvests() (*types.Transaction, error) {
	return _BLPPledgePool.Contract.Harvests(&_BLPPledgePool.TransactOpts)
}

// Harvests is a paid mutator transaction binding the contract method 0x2a9bd994.
//
// Solidity: function harvests() returns()
func (_BLPPledgePool *BLPPledgePoolTransactorSession) Harvests() (*types.Transaction, error) {
	return _BLPPledgePool.Contract.Harvests(&_BLPPledgePool.TransactOpts)
}

// MassUpdatePools is a paid mutator transaction binding the contract method 0x630b5ba1.
//
// Solidity: function massUpdatePools() returns()
func (_BLPPledgePool *BLPPledgePoolTransactor) MassUpdatePools(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BLPPledgePool.contract.Transact(opts, "massUpdatePools")
}

// MassUpdatePools is a paid mutator transaction binding the contract method 0x630b5ba1.
//
// Solidity: function massUpdatePools() returns()
func (_BLPPledgePool *BLPPledgePoolSession) MassUpdatePools() (*types.Transaction, error) {
	return _BLPPledgePool.Contract.MassUpdatePools(&_BLPPledgePool.TransactOpts)
}

// MassUpdatePools is a paid mutator transaction binding the contract method 0x630b5ba1.
//
// Solidity: function massUpdatePools() returns()
func (_BLPPledgePool *BLPPledgePoolTransactorSession) MassUpdatePools() (*types.Transaction, error) {
	return _BLPPledgePool.Contract.MassUpdatePools(&_BLPPledgePool.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_BLPPledgePool *BLPPledgePoolTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BLPPledgePool.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_BLPPledgePool *BLPPledgePoolSession) RenounceOwnership() (*types.Transaction, error) {
	return _BLPPledgePool.Contract.RenounceOwnership(&_BLPPledgePool.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_BLPPledgePool *BLPPledgePoolTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _BLPPledgePool.Contract.RenounceOwnership(&_BLPPledgePool.TransactOpts)
}

// Set is a paid mutator transaction binding the contract method 0x64482f79.
//
// Solidity: function set(uint256 _pid, uint256 _allocPoint, bool _withUpdate) returns()
func (_BLPPledgePool *BLPPledgePoolTransactor) Set(opts *bind.TransactOpts, _pid *big.Int, _allocPoint *big.Int, _withUpdate bool) (*types.Transaction, error) {
	return _BLPPledgePool.contract.Transact(opts, "set", _pid, _allocPoint, _withUpdate)
}

// Set is a paid mutator transaction binding the contract method 0x64482f79.
//
// Solidity: function set(uint256 _pid, uint256 _allocPoint, bool _withUpdate) returns()
func (_BLPPledgePool *BLPPledgePoolSession) Set(_pid *big.Int, _allocPoint *big.Int, _withUpdate bool) (*types.Transaction, error) {
	return _BLPPledgePool.Contract.Set(&_BLPPledgePool.TransactOpts, _pid, _allocPoint, _withUpdate)
}

// Set is a paid mutator transaction binding the contract method 0x64482f79.
//
// Solidity: function set(uint256 _pid, uint256 _allocPoint, bool _withUpdate) returns()
func (_BLPPledgePool *BLPPledgePoolTransactorSession) Set(_pid *big.Int, _allocPoint *big.Int, _withUpdate bool) (*types.Transaction, error) {
	return _BLPPledgePool.Contract.Set(&_BLPPledgePool.TransactOpts, _pid, _allocPoint, _withUpdate)
}

// SetLastBalance is a paid mutator transaction binding the contract method 0x9751fb0f.
//
// Solidity: function setLastBalance(uint256 _amount) returns()
func (_BLPPledgePool *BLPPledgePoolTransactor) SetLastBalance(opts *bind.TransactOpts, _amount *big.Int) (*types.Transaction, error) {
	return _BLPPledgePool.contract.Transact(opts, "setLastBalance", _amount)
}

// SetLastBalance is a paid mutator transaction binding the contract method 0x9751fb0f.
//
// Solidity: function setLastBalance(uint256 _amount) returns()
func (_BLPPledgePool *BLPPledgePoolSession) SetLastBalance(_amount *big.Int) (*types.Transaction, error) {
	return _BLPPledgePool.Contract.SetLastBalance(&_BLPPledgePool.TransactOpts, _amount)
}

// SetLastBalance is a paid mutator transaction binding the contract method 0x9751fb0f.
//
// Solidity: function setLastBalance(uint256 _amount) returns()
func (_BLPPledgePool *BLPPledgePoolTransactorSession) SetLastBalance(_amount *big.Int) (*types.Transaction, error) {
	return _BLPPledgePool.Contract.SetLastBalance(&_BLPPledgePool.TransactOpts, _amount)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_BLPPledgePool *BLPPledgePoolTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _BLPPledgePool.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_BLPPledgePool *BLPPledgePoolSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _BLPPledgePool.Contract.TransferOwnership(&_BLPPledgePool.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_BLPPledgePool *BLPPledgePoolTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _BLPPledgePool.Contract.TransferOwnership(&_BLPPledgePool.TransactOpts, newOwner)
}

// UpdateMultiplier is a paid mutator transaction binding the contract method 0x5ffe6146.
//
// Solidity: function updateMultiplier(uint256 multiplierNumber) returns()
func (_BLPPledgePool *BLPPledgePoolTransactor) UpdateMultiplier(opts *bind.TransactOpts, multiplierNumber *big.Int) (*types.Transaction, error) {
	return _BLPPledgePool.contract.Transact(opts, "updateMultiplier", multiplierNumber)
}

// UpdateMultiplier is a paid mutator transaction binding the contract method 0x5ffe6146.
//
// Solidity: function updateMultiplier(uint256 multiplierNumber) returns()
func (_BLPPledgePool *BLPPledgePoolSession) UpdateMultiplier(multiplierNumber *big.Int) (*types.Transaction, error) {
	return _BLPPledgePool.Contract.UpdateMultiplier(&_BLPPledgePool.TransactOpts, multiplierNumber)
}

// UpdateMultiplier is a paid mutator transaction binding the contract method 0x5ffe6146.
//
// Solidity: function updateMultiplier(uint256 multiplierNumber) returns()
func (_BLPPledgePool *BLPPledgePoolTransactorSession) UpdateMultiplier(multiplierNumber *big.Int) (*types.Transaction, error) {
	return _BLPPledgePool.Contract.UpdateMultiplier(&_BLPPledgePool.TransactOpts, multiplierNumber)
}

// UpdatePool is a paid mutator transaction binding the contract method 0x63b8d33a.
//
// Solidity: function updatePool(uint256 _pid, bool _isBatch) returns()
func (_BLPPledgePool *BLPPledgePoolTransactor) UpdatePool(opts *bind.TransactOpts, _pid *big.Int, _isBatch bool) (*types.Transaction, error) {
	return _BLPPledgePool.contract.Transact(opts, "updatePool", _pid, _isBatch)
}

// UpdatePool is a paid mutator transaction binding the contract method 0x63b8d33a.
//
// Solidity: function updatePool(uint256 _pid, bool _isBatch) returns()
func (_BLPPledgePool *BLPPledgePoolSession) UpdatePool(_pid *big.Int, _isBatch bool) (*types.Transaction, error) {
	return _BLPPledgePool.Contract.UpdatePool(&_BLPPledgePool.TransactOpts, _pid, _isBatch)
}

// UpdatePool is a paid mutator transaction binding the contract method 0x63b8d33a.
//
// Solidity: function updatePool(uint256 _pid, bool _isBatch) returns()
func (_BLPPledgePool *BLPPledgePoolTransactorSession) UpdatePool(_pid *big.Int, _isBatch bool) (*types.Transaction, error) {
	return _BLPPledgePool.Contract.UpdatePool(&_BLPPledgePool.TransactOpts, _pid, _isBatch)
}

// Withdraw is a paid mutator transaction binding the contract method 0xa41fe49f.
//
// Solidity: function withdraw(uint256 _pid, uint256 _amount, uint256 id) returns()
func (_BLPPledgePool *BLPPledgePoolTransactor) Withdraw(opts *bind.TransactOpts, _pid *big.Int, _amount *big.Int, id *big.Int) (*types.Transaction, error) {
	return _BLPPledgePool.contract.Transact(opts, "withdraw", _pid, _amount, id)
}

// Withdraw is a paid mutator transaction binding the contract method 0xa41fe49f.
//
// Solidity: function withdraw(uint256 _pid, uint256 _amount, uint256 id) returns()
func (_BLPPledgePool *BLPPledgePoolSession) Withdraw(_pid *big.Int, _amount *big.Int, id *big.Int) (*types.Transaction, error) {
	return _BLPPledgePool.Contract.Withdraw(&_BLPPledgePool.TransactOpts, _pid, _amount, id)
}

// Withdraw is a paid mutator transaction binding the contract method 0xa41fe49f.
//
// Solidity: function withdraw(uint256 _pid, uint256 _amount, uint256 id) returns()
func (_BLPPledgePool *BLPPledgePoolTransactorSession) Withdraw(_pid *big.Int, _amount *big.Int, id *big.Int) (*types.Transaction, error) {
	return _BLPPledgePool.Contract.Withdraw(&_BLPPledgePool.TransactOpts, _pid, _amount, id)
}

// BLPPledgePoolDepositIterator is returned from FilterDeposit and is used to iterate over the raw logs and unpacked data for Deposit events raised by the BLPPledgePool contract.
type BLPPledgePoolDepositIterator struct {
	Event *BLPPledgePoolDeposit // Event containing the contract specifics and raw log

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
func (it *BLPPledgePoolDepositIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BLPPledgePoolDeposit)
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
		it.Event = new(BLPPledgePoolDeposit)
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
func (it *BLPPledgePoolDepositIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BLPPledgePoolDepositIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BLPPledgePoolDeposit represents a Deposit event raised by the BLPPledgePool contract.
type BLPPledgePoolDeposit struct {
	User       common.Address
	Pid        *big.Int
	Id         *big.Int
	Amount     *big.Int
	LockPeriod *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterDeposit is a free log retrieval operation binding the contract event 0x7162984403f6c73c8639375d45a9187dfd04602231bd8e587c415718b5f7e5f9.
//
// Solidity: event Deposit(address indexed user, uint256 indexed pid, uint256 id, uint256 amount, uint256 lockPeriod)
func (_BLPPledgePool *BLPPledgePoolFilterer) FilterDeposit(opts *bind.FilterOpts, user []common.Address, pid []*big.Int) (*BLPPledgePoolDepositIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var pidRule []interface{}
	for _, pidItem := range pid {
		pidRule = append(pidRule, pidItem)
	}

	logs, sub, err := _BLPPledgePool.contract.FilterLogs(opts, "Deposit", userRule, pidRule)
	if err != nil {
		return nil, err
	}
	return &BLPPledgePoolDepositIterator{contract: _BLPPledgePool.contract, event: "Deposit", logs: logs, sub: sub}, nil
}

// WatchDeposit is a free log subscription operation binding the contract event 0x7162984403f6c73c8639375d45a9187dfd04602231bd8e587c415718b5f7e5f9.
//
// Solidity: event Deposit(address indexed user, uint256 indexed pid, uint256 id, uint256 amount, uint256 lockPeriod)
func (_BLPPledgePool *BLPPledgePoolFilterer) WatchDeposit(opts *bind.WatchOpts, sink chan<- *BLPPledgePoolDeposit, user []common.Address, pid []*big.Int) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var pidRule []interface{}
	for _, pidItem := range pid {
		pidRule = append(pidRule, pidItem)
	}

	logs, sub, err := _BLPPledgePool.contract.WatchLogs(opts, "Deposit", userRule, pidRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BLPPledgePoolDeposit)
				if err := _BLPPledgePool.contract.UnpackLog(event, "Deposit", log); err != nil {
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

// ParseDeposit is a log parse operation binding the contract event 0x7162984403f6c73c8639375d45a9187dfd04602231bd8e587c415718b5f7e5f9.
//
// Solidity: event Deposit(address indexed user, uint256 indexed pid, uint256 id, uint256 amount, uint256 lockPeriod)
func (_BLPPledgePool *BLPPledgePoolFilterer) ParseDeposit(log types.Log) (*BLPPledgePoolDeposit, error) {
	event := new(BLPPledgePoolDeposit)
	if err := _BLPPledgePool.contract.UnpackLog(event, "Deposit", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BLPPledgePoolEmergencyWithdrawIterator is returned from FilterEmergencyWithdraw and is used to iterate over the raw logs and unpacked data for EmergencyWithdraw events raised by the BLPPledgePool contract.
type BLPPledgePoolEmergencyWithdrawIterator struct {
	Event *BLPPledgePoolEmergencyWithdraw // Event containing the contract specifics and raw log

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
func (it *BLPPledgePoolEmergencyWithdrawIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BLPPledgePoolEmergencyWithdraw)
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
		it.Event = new(BLPPledgePoolEmergencyWithdraw)
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
func (it *BLPPledgePoolEmergencyWithdrawIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BLPPledgePoolEmergencyWithdrawIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BLPPledgePoolEmergencyWithdraw represents a EmergencyWithdraw event raised by the BLPPledgePool contract.
type BLPPledgePoolEmergencyWithdraw struct {
	User   common.Address
	Pid    *big.Int
	Id     *big.Int
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterEmergencyWithdraw is a free log retrieval operation binding the contract event 0x2369db1bafee945aee5630782f4a170682e3f8188d8dc247a4c73eb8c9e692d2.
//
// Solidity: event EmergencyWithdraw(address indexed user, uint256 indexed pid, uint256 id, uint256 amount)
func (_BLPPledgePool *BLPPledgePoolFilterer) FilterEmergencyWithdraw(opts *bind.FilterOpts, user []common.Address, pid []*big.Int) (*BLPPledgePoolEmergencyWithdrawIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var pidRule []interface{}
	for _, pidItem := range pid {
		pidRule = append(pidRule, pidItem)
	}

	logs, sub, err := _BLPPledgePool.contract.FilterLogs(opts, "EmergencyWithdraw", userRule, pidRule)
	if err != nil {
		return nil, err
	}
	return &BLPPledgePoolEmergencyWithdrawIterator{contract: _BLPPledgePool.contract, event: "EmergencyWithdraw", logs: logs, sub: sub}, nil
}

// WatchEmergencyWithdraw is a free log subscription operation binding the contract event 0x2369db1bafee945aee5630782f4a170682e3f8188d8dc247a4c73eb8c9e692d2.
//
// Solidity: event EmergencyWithdraw(address indexed user, uint256 indexed pid, uint256 id, uint256 amount)
func (_BLPPledgePool *BLPPledgePoolFilterer) WatchEmergencyWithdraw(opts *bind.WatchOpts, sink chan<- *BLPPledgePoolEmergencyWithdraw, user []common.Address, pid []*big.Int) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var pidRule []interface{}
	for _, pidItem := range pid {
		pidRule = append(pidRule, pidItem)
	}

	logs, sub, err := _BLPPledgePool.contract.WatchLogs(opts, "EmergencyWithdraw", userRule, pidRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BLPPledgePoolEmergencyWithdraw)
				if err := _BLPPledgePool.contract.UnpackLog(event, "EmergencyWithdraw", log); err != nil {
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

// ParseEmergencyWithdraw is a log parse operation binding the contract event 0x2369db1bafee945aee5630782f4a170682e3f8188d8dc247a4c73eb8c9e692d2.
//
// Solidity: event EmergencyWithdraw(address indexed user, uint256 indexed pid, uint256 id, uint256 amount)
func (_BLPPledgePool *BLPPledgePoolFilterer) ParseEmergencyWithdraw(log types.Log) (*BLPPledgePoolEmergencyWithdraw, error) {
	event := new(BLPPledgePoolEmergencyWithdraw)
	if err := _BLPPledgePool.contract.UnpackLog(event, "EmergencyWithdraw", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BLPPledgePoolHarvestIterator is returned from FilterHarvest and is used to iterate over the raw logs and unpacked data for Harvest events raised by the BLPPledgePool contract.
type BLPPledgePoolHarvestIterator struct {
	Event *BLPPledgePoolHarvest // Event containing the contract specifics and raw log

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
func (it *BLPPledgePoolHarvestIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BLPPledgePoolHarvest)
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
		it.Event = new(BLPPledgePoolHarvest)
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
func (it *BLPPledgePoolHarvestIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BLPPledgePoolHarvestIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BLPPledgePoolHarvest represents a Harvest event raised by the BLPPledgePool contract.
type BLPPledgePoolHarvest struct {
	User   common.Address
	Pid    *big.Int
	Id     *big.Int
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterHarvest is a free log retrieval operation binding the contract event 0x4534f107610758c3931de9ad1e176476fcfb8c74adf920167e1d54ee84fcfe76.
//
// Solidity: event Harvest(address indexed user, uint256 indexed pid, uint256 id, uint256 amount)
func (_BLPPledgePool *BLPPledgePoolFilterer) FilterHarvest(opts *bind.FilterOpts, user []common.Address, pid []*big.Int) (*BLPPledgePoolHarvestIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var pidRule []interface{}
	for _, pidItem := range pid {
		pidRule = append(pidRule, pidItem)
	}

	logs, sub, err := _BLPPledgePool.contract.FilterLogs(opts, "Harvest", userRule, pidRule)
	if err != nil {
		return nil, err
	}
	return &BLPPledgePoolHarvestIterator{contract: _BLPPledgePool.contract, event: "Harvest", logs: logs, sub: sub}, nil
}

// WatchHarvest is a free log subscription operation binding the contract event 0x4534f107610758c3931de9ad1e176476fcfb8c74adf920167e1d54ee84fcfe76.
//
// Solidity: event Harvest(address indexed user, uint256 indexed pid, uint256 id, uint256 amount)
func (_BLPPledgePool *BLPPledgePoolFilterer) WatchHarvest(opts *bind.WatchOpts, sink chan<- *BLPPledgePoolHarvest, user []common.Address, pid []*big.Int) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var pidRule []interface{}
	for _, pidItem := range pid {
		pidRule = append(pidRule, pidItem)
	}

	logs, sub, err := _BLPPledgePool.contract.WatchLogs(opts, "Harvest", userRule, pidRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BLPPledgePoolHarvest)
				if err := _BLPPledgePool.contract.UnpackLog(event, "Harvest", log); err != nil {
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

// ParseHarvest is a log parse operation binding the contract event 0x4534f107610758c3931de9ad1e176476fcfb8c74adf920167e1d54ee84fcfe76.
//
// Solidity: event Harvest(address indexed user, uint256 indexed pid, uint256 id, uint256 amount)
func (_BLPPledgePool *BLPPledgePoolFilterer) ParseHarvest(log types.Log) (*BLPPledgePoolHarvest, error) {
	event := new(BLPPledgePoolHarvest)
	if err := _BLPPledgePool.contract.UnpackLog(event, "Harvest", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BLPPledgePoolOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the BLPPledgePool contract.
type BLPPledgePoolOwnershipTransferredIterator struct {
	Event *BLPPledgePoolOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *BLPPledgePoolOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BLPPledgePoolOwnershipTransferred)
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
		it.Event = new(BLPPledgePoolOwnershipTransferred)
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
func (it *BLPPledgePoolOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BLPPledgePoolOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BLPPledgePoolOwnershipTransferred represents a OwnershipTransferred event raised by the BLPPledgePool contract.
type BLPPledgePoolOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_BLPPledgePool *BLPPledgePoolFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*BLPPledgePoolOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _BLPPledgePool.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &BLPPledgePoolOwnershipTransferredIterator{contract: _BLPPledgePool.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_BLPPledgePool *BLPPledgePoolFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *BLPPledgePoolOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _BLPPledgePool.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BLPPledgePoolOwnershipTransferred)
				if err := _BLPPledgePool.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_BLPPledgePool *BLPPledgePoolFilterer) ParseOwnershipTransferred(log types.Log) (*BLPPledgePoolOwnershipTransferred, error) {
	event := new(BLPPledgePoolOwnershipTransferred)
	if err := _BLPPledgePool.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BLPPledgePoolWithdrawIterator is returned from FilterWithdraw and is used to iterate over the raw logs and unpacked data for Withdraw events raised by the BLPPledgePool contract.
type BLPPledgePoolWithdrawIterator struct {
	Event *BLPPledgePoolWithdraw // Event containing the contract specifics and raw log

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
func (it *BLPPledgePoolWithdrawIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BLPPledgePoolWithdraw)
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
		it.Event = new(BLPPledgePoolWithdraw)
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
func (it *BLPPledgePoolWithdrawIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BLPPledgePoolWithdrawIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BLPPledgePoolWithdraw represents a Withdraw event raised by the BLPPledgePool contract.
type BLPPledgePoolWithdraw struct {
	User   common.Address
	Pid    *big.Int
	Id     *big.Int
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterWithdraw is a free log retrieval operation binding the contract event 0x02f25270a4d87bea75db541cdfe559334a275b4a233520ed6c0a2429667cca94.
//
// Solidity: event Withdraw(address indexed user, uint256 indexed pid, uint256 id, uint256 amount)
func (_BLPPledgePool *BLPPledgePoolFilterer) FilterWithdraw(opts *bind.FilterOpts, user []common.Address, pid []*big.Int) (*BLPPledgePoolWithdrawIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var pidRule []interface{}
	for _, pidItem := range pid {
		pidRule = append(pidRule, pidItem)
	}

	logs, sub, err := _BLPPledgePool.contract.FilterLogs(opts, "Withdraw", userRule, pidRule)
	if err != nil {
		return nil, err
	}
	return &BLPPledgePoolWithdrawIterator{contract: _BLPPledgePool.contract, event: "Withdraw", logs: logs, sub: sub}, nil
}

// WatchWithdraw is a free log subscription operation binding the contract event 0x02f25270a4d87bea75db541cdfe559334a275b4a233520ed6c0a2429667cca94.
//
// Solidity: event Withdraw(address indexed user, uint256 indexed pid, uint256 id, uint256 amount)
func (_BLPPledgePool *BLPPledgePoolFilterer) WatchWithdraw(opts *bind.WatchOpts, sink chan<- *BLPPledgePoolWithdraw, user []common.Address, pid []*big.Int) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var pidRule []interface{}
	for _, pidItem := range pid {
		pidRule = append(pidRule, pidItem)
	}

	logs, sub, err := _BLPPledgePool.contract.WatchLogs(opts, "Withdraw", userRule, pidRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BLPPledgePoolWithdraw)
				if err := _BLPPledgePool.contract.UnpackLog(event, "Withdraw", log); err != nil {
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

// ParseWithdraw is a log parse operation binding the contract event 0x02f25270a4d87bea75db541cdfe559334a275b4a233520ed6c0a2429667cca94.
//
// Solidity: event Withdraw(address indexed user, uint256 indexed pid, uint256 id, uint256 amount)
func (_BLPPledgePool *BLPPledgePoolFilterer) ParseWithdraw(log types.Log) (*BLPPledgePoolWithdraw, error) {
	event := new(BLPPledgePoolWithdraw)
	if err := _BLPPledgePool.contract.UnpackLog(event, "Withdraw", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
