// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ens

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
)

// IETHRegistrarControllerReabateInfo is an auto generated low-level Go binding around an user-defined struct.
type IETHRegistrarControllerReabateInfo struct {
	RebateAddress common.Address
	RebateAmount  *big.Int
	RebateName    string
}

// IETHRegistrarControllerRegisterData is an auto generated low-level Go binding around an user-defined struct.
type IETHRegistrarControllerRegisterData struct {
	Name                 string
	Owner                common.Address
	Duration             *big.Int
	Secret               [32]byte
	Resolver             common.Address
	Data                 [][]byte
	ReverseRecord        bool
	OwnerControlledFuses uint16
	RebateName           string
}

// IETHRegistrarControllerRegisterInfo is an auto generated low-level Go binding around an user-defined struct.
type IETHRegistrarControllerRegisterInfo struct {
	Name     string
	Label    [32]byte
	Owner    common.Address
	BaseCost *big.Int
	Premium  *big.Int
	Expires  *big.Int
}

// IPriceOraclePrice is an auto generated low-level Go binding around an user-defined struct.
type IPriceOraclePrice struct {
	Base    *big.Int
	Premium *big.Int
}

// EnsMetaData contains all meta data concerning the Ens contract.
var EnsMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractBaseRegistrarImplementation\",\"name\":\"_base\",\"type\":\"address\"},{\"internalType\":\"contractIPriceOracle\",\"name\":\"_prices\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_minCommitmentAge\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_maxCommitmentAge\",\"type\":\"uint256\"},{\"internalType\":\"contractReverseRegistrar\",\"name\":\"_reverseRegistrar\",\"type\":\"address\"},{\"internalType\":\"contractINameWrapper\",\"name\":\"_nameWrapper\",\"type\":\"address\"},{\"internalType\":\"contractENS\",\"name\":\"_ens\",\"type\":\"address\"},{\"internalType\":\"contractIRebateRegistrar\",\"name\":\"_rebateRegistrar\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_usdtAddress\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"commitment\",\"type\":\"bytes32\"}],\"name\":\"CommitmentTooNew\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"commitment\",\"type\":\"bytes32\"}],\"name\":\"CommitmentTooOld\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientValue\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MaxCommitmentAgeTooHigh\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MaxCommitmentAgeTooLow\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"name\":\"NameNotAvailable\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ResolverRequiredWhenDataSupplied\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"commitment\",\"type\":\"bytes32\"}],\"name\":\"UnexpiredCommitmentExists\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"bytes32\",\"name\":\"label\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"baseCost\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"premium\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expires\",\"type\":\"uint256\"}],\"indexed\":true,\"internalType\":\"structIETHRegistrarController.RegisterInfo\",\"name\":\"info\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"rebateAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"rebateAmount\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"rebateName\",\"type\":\"string\"}],\"indexed\":false,\"internalType\":\"structIETHRegistrarController.ReabateInfo\",\"name\":\"rebateIfno\",\"type\":\"tuple\"}],\"name\":\"NameRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"}],\"name\":\"NewPriceOracle\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"MIN_REGISTRATION_DURATION\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"name\":\"available\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"commitment\",\"type\":\"bytes32\"}],\"name\":\"commit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"commitments\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"duration\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"secret\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"resolver\",\"type\":\"address\"},{\"internalType\":\"bytes[]\",\"name\":\"data\",\"type\":\"bytes[]\"},{\"internalType\":\"bool\",\"name\":\"reverseRecord\",\"type\":\"bool\"},{\"internalType\":\"uint16\",\"name\":\"ownerControlledFuses\",\"type\":\"uint16\"},{\"internalType\":\"string\",\"name\":\"rebateName\",\"type\":\"string\"}],\"internalType\":\"structIETHRegistrarController.RegisterData\",\"name\":\"data\",\"type\":\"tuple\"}],\"name\":\"makeCommitment\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maxCommitmentAge\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minCommitmentAge\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"nameWrapper\",\"outputs\":[{\"internalType\":\"contractINameWrapper\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"prices\",\"outputs\":[{\"internalType\":\"contractIPriceOracle\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rebateRegistrar\",\"outputs\":[{\"internalType\":\"contractIRebateRegistrar\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"recoverFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"duration\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"secret\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"resolver\",\"type\":\"address\"},{\"internalType\":\"bytes[]\",\"name\":\"data\",\"type\":\"bytes[]\"},{\"internalType\":\"bool\",\"name\":\"reverseRecord\",\"type\":\"bool\"},{\"internalType\":\"uint16\",\"name\":\"ownerControlledFuses\",\"type\":\"uint16\"},{\"internalType\":\"string\",\"name\":\"rebateName\",\"type\":\"string\"}],\"internalType\":\"structIETHRegistrarController.RegisterData\",\"name\":\"data\",\"type\":\"tuple\"},{\"internalType\":\"bool\",\"name\":\"isERC20\",\"type\":\"bool\"}],\"name\":\"register\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"duration\",\"type\":\"uint256\"}],\"name\":\"rentPrice\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"base\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"premium\",\"type\":\"uint256\"}],\"internalType\":\"structIPriceOracle.Price\",\"name\":\"price\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"duration\",\"type\":\"uint256\"}],\"name\":\"rentUSDPrice\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"base\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"premium\",\"type\":\"uint256\"}],\"internalType\":\"structIPriceOracle.Price\",\"name\":\"price\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"reverseRegistrar\",\"outputs\":[{\"internalType\":\"contractReverseRegistrar\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIPriceOracle\",\"name\":\"_prices\",\"type\":\"address\"}],\"name\":\"setPriceOracle\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIRebateRegistrar\",\"name\":\"_rebateRegistrar\",\"type\":\"address\"}],\"name\":\"setRebateRegistrar\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"_steps\",\"type\":\"uint256[]\"}],\"name\":\"setSteps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"steps\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceID\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"usdtAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"name\":\"valid\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// EnsABI is the input ABI used to generate the binding from.
// Deprecated: Use EnsMetaData.ABI instead.
var EnsABI = EnsMetaData.ABI

// Ens is an auto generated Go binding around an Ethereum contract.
type Ens struct {
	EnsCaller     // Read-only binding to the contract
	EnsTransactor // Write-only binding to the contract
	EnsFilterer   // Log filterer for contract events
}

// EnsCaller is an auto generated read-only Go binding around an Ethereum contract.
type EnsCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EnsTransactor is an auto generated write-only Go binding around an Ethereum contract.
type EnsTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EnsFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type EnsFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EnsSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type EnsSession struct {
	Contract     *Ens              // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// EnsCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type EnsCallerSession struct {
	Contract *EnsCaller    // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// EnsTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type EnsTransactorSession struct {
	Contract     *EnsTransactor    // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// EnsRaw is an auto generated low-level Go binding around an Ethereum contract.
type EnsRaw struct {
	Contract *Ens // Generic contract binding to access the raw methods on
}

// EnsCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type EnsCallerRaw struct {
	Contract *EnsCaller // Generic read-only contract binding to access the raw methods on
}

// EnsTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type EnsTransactorRaw struct {
	Contract *EnsTransactor // Generic write-only contract binding to access the raw methods on
}

// NewEns creates a new instance of Ens, bound to a specific deployed contract.
func NewEns(address common.Address, backend bind.ContractBackend) (*Ens, error) {
	contract, err := bindEns(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Ens{EnsCaller: EnsCaller{contract: contract}, EnsTransactor: EnsTransactor{contract: contract}, EnsFilterer: EnsFilterer{contract: contract}}, nil
}

// NewEnsCaller creates a new read-only instance of Ens, bound to a specific deployed contract.
func NewEnsCaller(address common.Address, caller bind.ContractCaller) (*EnsCaller, error) {
	contract, err := bindEns(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &EnsCaller{contract: contract}, nil
}

// NewEnsTransactor creates a new write-only instance of Ens, bound to a specific deployed contract.
func NewEnsTransactor(address common.Address, transactor bind.ContractTransactor) (*EnsTransactor, error) {
	contract, err := bindEns(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &EnsTransactor{contract: contract}, nil
}

// NewEnsFilterer creates a new log filterer instance of Ens, bound to a specific deployed contract.
func NewEnsFilterer(address common.Address, filterer bind.ContractFilterer) (*EnsFilterer, error) {
	contract, err := bindEns(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &EnsFilterer{contract: contract}, nil
}

// bindEns binds a generic wrapper to an already deployed contract.
func bindEns(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(EnsABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Ens *EnsRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Ens.Contract.EnsCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Ens *EnsRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Ens.Contract.EnsTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Ens *EnsRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Ens.Contract.EnsTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Ens *EnsCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Ens.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Ens *EnsTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Ens.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Ens *EnsTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Ens.Contract.contract.Transact(opts, method, params...)
}

// MINREGISTRATIONDURATION is a free data retrieval call binding the contract method 0x8a95b09f.
//
// Solidity: function MIN_REGISTRATION_DURATION() view returns(uint256)
func (_Ens *EnsCaller) MINREGISTRATIONDURATION(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Ens.contract.Call(opts, &out, "MIN_REGISTRATION_DURATION")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MINREGISTRATIONDURATION is a free data retrieval call binding the contract method 0x8a95b09f.
//
// Solidity: function MIN_REGISTRATION_DURATION() view returns(uint256)
func (_Ens *EnsSession) MINREGISTRATIONDURATION() (*big.Int, error) {
	return _Ens.Contract.MINREGISTRATIONDURATION(&_Ens.CallOpts)
}

// MINREGISTRATIONDURATION is a free data retrieval call binding the contract method 0x8a95b09f.
//
// Solidity: function MIN_REGISTRATION_DURATION() view returns(uint256)
func (_Ens *EnsCallerSession) MINREGISTRATIONDURATION() (*big.Int, error) {
	return _Ens.Contract.MINREGISTRATIONDURATION(&_Ens.CallOpts)
}

// Available is a free data retrieval call binding the contract method 0xaeb8ce9b.
//
// Solidity: function available(string name) view returns(bool)
func (_Ens *EnsCaller) Available(opts *bind.CallOpts, name string) (bool, error) {
	var out []interface{}
	err := _Ens.contract.Call(opts, &out, "available", name)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Available is a free data retrieval call binding the contract method 0xaeb8ce9b.
//
// Solidity: function available(string name) view returns(bool)
func (_Ens *EnsSession) Available(name string) (bool, error) {
	return _Ens.Contract.Available(&_Ens.CallOpts, name)
}

// Available is a free data retrieval call binding the contract method 0xaeb8ce9b.
//
// Solidity: function available(string name) view returns(bool)
func (_Ens *EnsCallerSession) Available(name string) (bool, error) {
	return _Ens.Contract.Available(&_Ens.CallOpts, name)
}

// Commitments is a free data retrieval call binding the contract method 0x839df945.
//
// Solidity: function commitments(bytes32 ) view returns(uint256)
func (_Ens *EnsCaller) Commitments(opts *bind.CallOpts, arg0 [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _Ens.contract.Call(opts, &out, "commitments", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Commitments is a free data retrieval call binding the contract method 0x839df945.
//
// Solidity: function commitments(bytes32 ) view returns(uint256)
func (_Ens *EnsSession) Commitments(arg0 [32]byte) (*big.Int, error) {
	return _Ens.Contract.Commitments(&_Ens.CallOpts, arg0)
}

// Commitments is a free data retrieval call binding the contract method 0x839df945.
//
// Solidity: function commitments(bytes32 ) view returns(uint256)
func (_Ens *EnsCallerSession) Commitments(arg0 [32]byte) (*big.Int, error) {
	return _Ens.Contract.Commitments(&_Ens.CallOpts, arg0)
}

// MakeCommitment is a free data retrieval call binding the contract method 0x6a5feb59.
//
// Solidity: function makeCommitment((string,address,uint256,bytes32,address,bytes[],bool,uint16,string) data) pure returns(bytes32)
func (_Ens *EnsCaller) MakeCommitment(opts *bind.CallOpts, data IETHRegistrarControllerRegisterData) ([32]byte, error) {
	var out []interface{}
	err := _Ens.contract.Call(opts, &out, "makeCommitment", data)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// MakeCommitment is a free data retrieval call binding the contract method 0x6a5feb59.
//
// Solidity: function makeCommitment((string,address,uint256,bytes32,address,bytes[],bool,uint16,string) data) pure returns(bytes32)
func (_Ens *EnsSession) MakeCommitment(data IETHRegistrarControllerRegisterData) ([32]byte, error) {
	return _Ens.Contract.MakeCommitment(&_Ens.CallOpts, data)
}

// MakeCommitment is a free data retrieval call binding the contract method 0x6a5feb59.
//
// Solidity: function makeCommitment((string,address,uint256,bytes32,address,bytes[],bool,uint16,string) data) pure returns(bytes32)
func (_Ens *EnsCallerSession) MakeCommitment(data IETHRegistrarControllerRegisterData) ([32]byte, error) {
	return _Ens.Contract.MakeCommitment(&_Ens.CallOpts, data)
}

// MaxCommitmentAge is a free data retrieval call binding the contract method 0xce1e09c0.
//
// Solidity: function maxCommitmentAge() view returns(uint256)
func (_Ens *EnsCaller) MaxCommitmentAge(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Ens.contract.Call(opts, &out, "maxCommitmentAge")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MaxCommitmentAge is a free data retrieval call binding the contract method 0xce1e09c0.
//
// Solidity: function maxCommitmentAge() view returns(uint256)
func (_Ens *EnsSession) MaxCommitmentAge() (*big.Int, error) {
	return _Ens.Contract.MaxCommitmentAge(&_Ens.CallOpts)
}

// MaxCommitmentAge is a free data retrieval call binding the contract method 0xce1e09c0.
//
// Solidity: function maxCommitmentAge() view returns(uint256)
func (_Ens *EnsCallerSession) MaxCommitmentAge() (*big.Int, error) {
	return _Ens.Contract.MaxCommitmentAge(&_Ens.CallOpts)
}

// MinCommitmentAge is a free data retrieval call binding the contract method 0x8d839ffe.
//
// Solidity: function minCommitmentAge() view returns(uint256)
func (_Ens *EnsCaller) MinCommitmentAge(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Ens.contract.Call(opts, &out, "minCommitmentAge")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MinCommitmentAge is a free data retrieval call binding the contract method 0x8d839ffe.
//
// Solidity: function minCommitmentAge() view returns(uint256)
func (_Ens *EnsSession) MinCommitmentAge() (*big.Int, error) {
	return _Ens.Contract.MinCommitmentAge(&_Ens.CallOpts)
}

// MinCommitmentAge is a free data retrieval call binding the contract method 0x8d839ffe.
//
// Solidity: function minCommitmentAge() view returns(uint256)
func (_Ens *EnsCallerSession) MinCommitmentAge() (*big.Int, error) {
	return _Ens.Contract.MinCommitmentAge(&_Ens.CallOpts)
}

// NameWrapper is a free data retrieval call binding the contract method 0xa8e5fbc0.
//
// Solidity: function nameWrapper() view returns(address)
func (_Ens *EnsCaller) NameWrapper(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Ens.contract.Call(opts, &out, "nameWrapper")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// NameWrapper is a free data retrieval call binding the contract method 0xa8e5fbc0.
//
// Solidity: function nameWrapper() view returns(address)
func (_Ens *EnsSession) NameWrapper() (common.Address, error) {
	return _Ens.Contract.NameWrapper(&_Ens.CallOpts)
}

// NameWrapper is a free data retrieval call binding the contract method 0xa8e5fbc0.
//
// Solidity: function nameWrapper() view returns(address)
func (_Ens *EnsCallerSession) NameWrapper() (common.Address, error) {
	return _Ens.Contract.NameWrapper(&_Ens.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Ens *EnsCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Ens.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Ens *EnsSession) Owner() (common.Address, error) {
	return _Ens.Contract.Owner(&_Ens.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Ens *EnsCallerSession) Owner() (common.Address, error) {
	return _Ens.Contract.Owner(&_Ens.CallOpts)
}

// Prices is a free data retrieval call binding the contract method 0xd3419bf3.
//
// Solidity: function prices() view returns(address)
func (_Ens *EnsCaller) Prices(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Ens.contract.Call(opts, &out, "prices")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Prices is a free data retrieval call binding the contract method 0xd3419bf3.
//
// Solidity: function prices() view returns(address)
func (_Ens *EnsSession) Prices() (common.Address, error) {
	return _Ens.Contract.Prices(&_Ens.CallOpts)
}

// Prices is a free data retrieval call binding the contract method 0xd3419bf3.
//
// Solidity: function prices() view returns(address)
func (_Ens *EnsCallerSession) Prices() (common.Address, error) {
	return _Ens.Contract.Prices(&_Ens.CallOpts)
}

// RebateRegistrar is a free data retrieval call binding the contract method 0x8dfaf336.
//
// Solidity: function rebateRegistrar() view returns(address)
func (_Ens *EnsCaller) RebateRegistrar(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Ens.contract.Call(opts, &out, "rebateRegistrar")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RebateRegistrar is a free data retrieval call binding the contract method 0x8dfaf336.
//
// Solidity: function rebateRegistrar() view returns(address)
func (_Ens *EnsSession) RebateRegistrar() (common.Address, error) {
	return _Ens.Contract.RebateRegistrar(&_Ens.CallOpts)
}

// RebateRegistrar is a free data retrieval call binding the contract method 0x8dfaf336.
//
// Solidity: function rebateRegistrar() view returns(address)
func (_Ens *EnsCallerSession) RebateRegistrar() (common.Address, error) {
	return _Ens.Contract.RebateRegistrar(&_Ens.CallOpts)
}

// RentPrice is a free data retrieval call binding the contract method 0x83e7f6ff.
//
// Solidity: function rentPrice(string name, uint256 duration) view returns((uint256,uint256) price)
func (_Ens *EnsCaller) RentPrice(opts *bind.CallOpts, name string, duration *big.Int) (IPriceOraclePrice, error) {
	var out []interface{}
	err := _Ens.contract.Call(opts, &out, "rentPrice", name, duration)

	if err != nil {
		return *new(IPriceOraclePrice), err
	}

	out0 := *abi.ConvertType(out[0], new(IPriceOraclePrice)).(*IPriceOraclePrice)

	return out0, err

}

// RentPrice is a free data retrieval call binding the contract method 0x83e7f6ff.
//
// Solidity: function rentPrice(string name, uint256 duration) view returns((uint256,uint256) price)
func (_Ens *EnsSession) RentPrice(name string, duration *big.Int) (IPriceOraclePrice, error) {
	return _Ens.Contract.RentPrice(&_Ens.CallOpts, name, duration)
}

// RentPrice is a free data retrieval call binding the contract method 0x83e7f6ff.
//
// Solidity: function rentPrice(string name, uint256 duration) view returns((uint256,uint256) price)
func (_Ens *EnsCallerSession) RentPrice(name string, duration *big.Int) (IPriceOraclePrice, error) {
	return _Ens.Contract.RentPrice(&_Ens.CallOpts, name, duration)
}

// RentUSDPrice is a free data retrieval call binding the contract method 0x00617a54.
//
// Solidity: function rentUSDPrice(string name, uint256 duration) view returns((uint256,uint256) price)
func (_Ens *EnsCaller) RentUSDPrice(opts *bind.CallOpts, name string, duration *big.Int) (IPriceOraclePrice, error) {
	var out []interface{}
	err := _Ens.contract.Call(opts, &out, "rentUSDPrice", name, duration)

	if err != nil {
		return *new(IPriceOraclePrice), err
	}

	out0 := *abi.ConvertType(out[0], new(IPriceOraclePrice)).(*IPriceOraclePrice)

	return out0, err

}

// RentUSDPrice is a free data retrieval call binding the contract method 0x00617a54.
//
// Solidity: function rentUSDPrice(string name, uint256 duration) view returns((uint256,uint256) price)
func (_Ens *EnsSession) RentUSDPrice(name string, duration *big.Int) (IPriceOraclePrice, error) {
	return _Ens.Contract.RentUSDPrice(&_Ens.CallOpts, name, duration)
}

// RentUSDPrice is a free data retrieval call binding the contract method 0x00617a54.
//
// Solidity: function rentUSDPrice(string name, uint256 duration) view returns((uint256,uint256) price)
func (_Ens *EnsCallerSession) RentUSDPrice(name string, duration *big.Int) (IPriceOraclePrice, error) {
	return _Ens.Contract.RentUSDPrice(&_Ens.CallOpts, name, duration)
}

// ReverseRegistrar is a free data retrieval call binding the contract method 0x80869853.
//
// Solidity: function reverseRegistrar() view returns(address)
func (_Ens *EnsCaller) ReverseRegistrar(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Ens.contract.Call(opts, &out, "reverseRegistrar")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ReverseRegistrar is a free data retrieval call binding the contract method 0x80869853.
//
// Solidity: function reverseRegistrar() view returns(address)
func (_Ens *EnsSession) ReverseRegistrar() (common.Address, error) {
	return _Ens.Contract.ReverseRegistrar(&_Ens.CallOpts)
}

// ReverseRegistrar is a free data retrieval call binding the contract method 0x80869853.
//
// Solidity: function reverseRegistrar() view returns(address)
func (_Ens *EnsCallerSession) ReverseRegistrar() (common.Address, error) {
	return _Ens.Contract.ReverseRegistrar(&_Ens.CallOpts)
}

// Steps is a free data retrieval call binding the contract method 0x7217e0b9.
//
// Solidity: function steps(uint256 ) view returns(uint256)
func (_Ens *EnsCaller) Steps(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Ens.contract.Call(opts, &out, "steps", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Steps is a free data retrieval call binding the contract method 0x7217e0b9.
//
// Solidity: function steps(uint256 ) view returns(uint256)
func (_Ens *EnsSession) Steps(arg0 *big.Int) (*big.Int, error) {
	return _Ens.Contract.Steps(&_Ens.CallOpts, arg0)
}

// Steps is a free data retrieval call binding the contract method 0x7217e0b9.
//
// Solidity: function steps(uint256 ) view returns(uint256)
func (_Ens *EnsCallerSession) Steps(arg0 *big.Int) (*big.Int, error) {
	return _Ens.Contract.Steps(&_Ens.CallOpts, arg0)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceID) pure returns(bool)
func (_Ens *EnsCaller) SupportsInterface(opts *bind.CallOpts, interfaceID [4]byte) (bool, error) {
	var out []interface{}
	err := _Ens.contract.Call(opts, &out, "supportsInterface", interfaceID)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceID) pure returns(bool)
func (_Ens *EnsSession) SupportsInterface(interfaceID [4]byte) (bool, error) {
	return _Ens.Contract.SupportsInterface(&_Ens.CallOpts, interfaceID)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceID) pure returns(bool)
func (_Ens *EnsCallerSession) SupportsInterface(interfaceID [4]byte) (bool, error) {
	return _Ens.Contract.SupportsInterface(&_Ens.CallOpts, interfaceID)
}

// UsdtAddress is a free data retrieval call binding the contract method 0x9ab4a445.
//
// Solidity: function usdtAddress() view returns(address)
func (_Ens *EnsCaller) UsdtAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Ens.contract.Call(opts, &out, "usdtAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// UsdtAddress is a free data retrieval call binding the contract method 0x9ab4a445.
//
// Solidity: function usdtAddress() view returns(address)
func (_Ens *EnsSession) UsdtAddress() (common.Address, error) {
	return _Ens.Contract.UsdtAddress(&_Ens.CallOpts)
}

// UsdtAddress is a free data retrieval call binding the contract method 0x9ab4a445.
//
// Solidity: function usdtAddress() view returns(address)
func (_Ens *EnsCallerSession) UsdtAddress() (common.Address, error) {
	return _Ens.Contract.UsdtAddress(&_Ens.CallOpts)
}

// Valid is a free data retrieval call binding the contract method 0x9791c097.
//
// Solidity: function valid(string name) view returns(bool)
func (_Ens *EnsCaller) Valid(opts *bind.CallOpts, name string) (bool, error) {
	var out []interface{}
	err := _Ens.contract.Call(opts, &out, "valid", name)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Valid is a free data retrieval call binding the contract method 0x9791c097.
//
// Solidity: function valid(string name) view returns(bool)
func (_Ens *EnsSession) Valid(name string) (bool, error) {
	return _Ens.Contract.Valid(&_Ens.CallOpts, name)
}

// Valid is a free data retrieval call binding the contract method 0x9791c097.
//
// Solidity: function valid(string name) view returns(bool)
func (_Ens *EnsCallerSession) Valid(name string) (bool, error) {
	return _Ens.Contract.Valid(&_Ens.CallOpts, name)
}

// Commit is a paid mutator transaction binding the contract method 0xf14fcbc8.
//
// Solidity: function commit(bytes32 commitment) returns()
func (_Ens *EnsTransactor) Commit(opts *bind.TransactOpts, commitment [32]byte) (*types.Transaction, error) {
	return _Ens.contract.Transact(opts, "commit", commitment)
}

// Commit is a paid mutator transaction binding the contract method 0xf14fcbc8.
//
// Solidity: function commit(bytes32 commitment) returns()
func (_Ens *EnsSession) Commit(commitment [32]byte) (*types.Transaction, error) {
	return _Ens.Contract.Commit(&_Ens.TransactOpts, commitment)
}

// Commit is a paid mutator transaction binding the contract method 0xf14fcbc8.
//
// Solidity: function commit(bytes32 commitment) returns()
func (_Ens *EnsTransactorSession) Commit(commitment [32]byte) (*types.Transaction, error) {
	return _Ens.Contract.Commit(&_Ens.TransactOpts, commitment)
}

// RecoverFunds is a paid mutator transaction binding the contract method 0x5d3590d5.
//
// Solidity: function recoverFunds(address _token, address _to, uint256 _amount) returns()
func (_Ens *EnsTransactor) RecoverFunds(opts *bind.TransactOpts, _token common.Address, _to common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Ens.contract.Transact(opts, "recoverFunds", _token, _to, _amount)
}

// RecoverFunds is a paid mutator transaction binding the contract method 0x5d3590d5.
//
// Solidity: function recoverFunds(address _token, address _to, uint256 _amount) returns()
func (_Ens *EnsSession) RecoverFunds(_token common.Address, _to common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Ens.Contract.RecoverFunds(&_Ens.TransactOpts, _token, _to, _amount)
}

// RecoverFunds is a paid mutator transaction binding the contract method 0x5d3590d5.
//
// Solidity: function recoverFunds(address _token, address _to, uint256 _amount) returns()
func (_Ens *EnsTransactorSession) RecoverFunds(_token common.Address, _to common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Ens.Contract.RecoverFunds(&_Ens.TransactOpts, _token, _to, _amount)
}

// Register is a paid mutator transaction binding the contract method 0x001afebf.
//
// Solidity: function register((string,address,uint256,bytes32,address,bytes[],bool,uint16,string) data, bool isERC20) payable returns()
func (_Ens *EnsTransactor) Register(opts *bind.TransactOpts, data IETHRegistrarControllerRegisterData, isERC20 bool) (*types.Transaction, error) {
	return _Ens.contract.Transact(opts, "register", data, isERC20)
}

// Register is a paid mutator transaction binding the contract method 0x001afebf.
//
// Solidity: function register((string,address,uint256,bytes32,address,bytes[],bool,uint16,string) data, bool isERC20) payable returns()
func (_Ens *EnsSession) Register(data IETHRegistrarControllerRegisterData, isERC20 bool) (*types.Transaction, error) {
	return _Ens.Contract.Register(&_Ens.TransactOpts, data, isERC20)
}

// Register is a paid mutator transaction binding the contract method 0x001afebf.
//
// Solidity: function register((string,address,uint256,bytes32,address,bytes[],bool,uint16,string) data, bool isERC20) payable returns()
func (_Ens *EnsTransactorSession) Register(data IETHRegistrarControllerRegisterData, isERC20 bool) (*types.Transaction, error) {
	return _Ens.Contract.Register(&_Ens.TransactOpts, data, isERC20)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Ens *EnsTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Ens.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Ens *EnsSession) RenounceOwnership() (*types.Transaction, error) {
	return _Ens.Contract.RenounceOwnership(&_Ens.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Ens *EnsTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Ens.Contract.RenounceOwnership(&_Ens.TransactOpts)
}

// SetPriceOracle is a paid mutator transaction binding the contract method 0x530e784f.
//
// Solidity: function setPriceOracle(address _prices) returns()
func (_Ens *EnsTransactor) SetPriceOracle(opts *bind.TransactOpts, _prices common.Address) (*types.Transaction, error) {
	return _Ens.contract.Transact(opts, "setPriceOracle", _prices)
}

// SetPriceOracle is a paid mutator transaction binding the contract method 0x530e784f.
//
// Solidity: function setPriceOracle(address _prices) returns()
func (_Ens *EnsSession) SetPriceOracle(_prices common.Address) (*types.Transaction, error) {
	return _Ens.Contract.SetPriceOracle(&_Ens.TransactOpts, _prices)
}

// SetPriceOracle is a paid mutator transaction binding the contract method 0x530e784f.
//
// Solidity: function setPriceOracle(address _prices) returns()
func (_Ens *EnsTransactorSession) SetPriceOracle(_prices common.Address) (*types.Transaction, error) {
	return _Ens.Contract.SetPriceOracle(&_Ens.TransactOpts, _prices)
}

// SetRebateRegistrar is a paid mutator transaction binding the contract method 0x40317b59.
//
// Solidity: function setRebateRegistrar(address _rebateRegistrar) returns()
func (_Ens *EnsTransactor) SetRebateRegistrar(opts *bind.TransactOpts, _rebateRegistrar common.Address) (*types.Transaction, error) {
	return _Ens.contract.Transact(opts, "setRebateRegistrar", _rebateRegistrar)
}

// SetRebateRegistrar is a paid mutator transaction binding the contract method 0x40317b59.
//
// Solidity: function setRebateRegistrar(address _rebateRegistrar) returns()
func (_Ens *EnsSession) SetRebateRegistrar(_rebateRegistrar common.Address) (*types.Transaction, error) {
	return _Ens.Contract.SetRebateRegistrar(&_Ens.TransactOpts, _rebateRegistrar)
}

// SetRebateRegistrar is a paid mutator transaction binding the contract method 0x40317b59.
//
// Solidity: function setRebateRegistrar(address _rebateRegistrar) returns()
func (_Ens *EnsTransactorSession) SetRebateRegistrar(_rebateRegistrar common.Address) (*types.Transaction, error) {
	return _Ens.Contract.SetRebateRegistrar(&_Ens.TransactOpts, _rebateRegistrar)
}

// SetSteps is a paid mutator transaction binding the contract method 0xf1777e7f.
//
// Solidity: function setSteps(uint256[] _steps) returns()
func (_Ens *EnsTransactor) SetSteps(opts *bind.TransactOpts, _steps []*big.Int) (*types.Transaction, error) {
	return _Ens.contract.Transact(opts, "setSteps", _steps)
}

// SetSteps is a paid mutator transaction binding the contract method 0xf1777e7f.
//
// Solidity: function setSteps(uint256[] _steps) returns()
func (_Ens *EnsSession) SetSteps(_steps []*big.Int) (*types.Transaction, error) {
	return _Ens.Contract.SetSteps(&_Ens.TransactOpts, _steps)
}

// SetSteps is a paid mutator transaction binding the contract method 0xf1777e7f.
//
// Solidity: function setSteps(uint256[] _steps) returns()
func (_Ens *EnsTransactorSession) SetSteps(_steps []*big.Int) (*types.Transaction, error) {
	return _Ens.Contract.SetSteps(&_Ens.TransactOpts, _steps)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Ens *EnsTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Ens.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Ens *EnsSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Ens.Contract.TransferOwnership(&_Ens.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Ens *EnsTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Ens.Contract.TransferOwnership(&_Ens.TransactOpts, newOwner)
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns()
func (_Ens *EnsTransactor) Withdraw(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Ens.contract.Transact(opts, "withdraw")
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns()
func (_Ens *EnsSession) Withdraw() (*types.Transaction, error) {
	return _Ens.Contract.Withdraw(&_Ens.TransactOpts)
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns()
func (_Ens *EnsTransactorSession) Withdraw() (*types.Transaction, error) {
	return _Ens.Contract.Withdraw(&_Ens.TransactOpts)
}

// EnsNameRegisteredIterator is returned from FilterNameRegistered and is used to iterate over the raw logs and unpacked data for NameRegistered events raised by the Ens contract.
type EnsNameRegisteredIterator struct {
	Event *EnsNameRegistered // Event containing the contract specifics and raw log

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
func (it *EnsNameRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EnsNameRegistered)
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
		it.Event = new(EnsNameRegistered)
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
func (it *EnsNameRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EnsNameRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EnsNameRegistered represents a NameRegistered event raised by the Ens contract.
type EnsNameRegistered struct {
	Info       IETHRegistrarControllerRegisterInfo
	RebateIfno IETHRegistrarControllerReabateInfo
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterNameRegistered is a free log retrieval operation binding the contract event 0x45b82dd7183546a2b1c0521286176455587e43d887bd7bd833bbe0d24139ee52.
//
// Solidity: event NameRegistered((string,bytes32,address,uint256,uint256,uint256) indexed info, (address,uint256,string) rebateIfno)
func (_Ens *EnsFilterer) FilterNameRegistered(opts *bind.FilterOpts, info []IETHRegistrarControllerRegisterInfo) (*EnsNameRegisteredIterator, error) {

	var infoRule []interface{}
	for _, infoItem := range info {
		infoRule = append(infoRule, infoItem)
	}

	logs, sub, err := _Ens.contract.FilterLogs(opts, "NameRegistered", infoRule)
	if err != nil {
		return nil, err
	}
	return &EnsNameRegisteredIterator{contract: _Ens.contract, event: "NameRegistered", logs: logs, sub: sub}, nil
}

// WatchNameRegistered is a free log subscription operation binding the contract event 0x45b82dd7183546a2b1c0521286176455587e43d887bd7bd833bbe0d24139ee52.
//
// Solidity: event NameRegistered((string,bytes32,address,uint256,uint256,uint256) indexed info, (address,uint256,string) rebateIfno)
func (_Ens *EnsFilterer) WatchNameRegistered(opts *bind.WatchOpts, sink chan<- *EnsNameRegistered, info []IETHRegistrarControllerRegisterInfo) (event.Subscription, error) {

	var infoRule []interface{}
	for _, infoItem := range info {
		infoRule = append(infoRule, infoItem)
	}

	logs, sub, err := _Ens.contract.WatchLogs(opts, "NameRegistered", infoRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EnsNameRegistered)
				if err := _Ens.contract.UnpackLog(event, "NameRegistered", log); err != nil {
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

// ParseNameRegistered is a log parse operation binding the contract event 0x45b82dd7183546a2b1c0521286176455587e43d887bd7bd833bbe0d24139ee52.
//
// Solidity: event NameRegistered((string,bytes32,address,uint256,uint256,uint256) indexed info, (address,uint256,string) rebateIfno)
func (_Ens *EnsFilterer) ParseNameRegistered(log types.Log) (*EnsNameRegistered, error) {
	event := new(EnsNameRegistered)
	if err := _Ens.contract.UnpackLog(event, "NameRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EnsNewPriceOracleIterator is returned from FilterNewPriceOracle and is used to iterate over the raw logs and unpacked data for NewPriceOracle events raised by the Ens contract.
type EnsNewPriceOracleIterator struct {
	Event *EnsNewPriceOracle // Event containing the contract specifics and raw log

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
func (it *EnsNewPriceOracleIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EnsNewPriceOracle)
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
		it.Event = new(EnsNewPriceOracle)
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
func (it *EnsNewPriceOracleIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EnsNewPriceOracleIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EnsNewPriceOracle represents a NewPriceOracle event raised by the Ens contract.
type EnsNewPriceOracle struct {
	Oracle common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterNewPriceOracle is a free log retrieval operation binding the contract event 0xf261845a790fe29bbd6631e2ca4a5bdc83e6eed7c3271d9590d97287e00e9123.
//
// Solidity: event NewPriceOracle(address indexed oracle)
func (_Ens *EnsFilterer) FilterNewPriceOracle(opts *bind.FilterOpts, oracle []common.Address) (*EnsNewPriceOracleIterator, error) {

	var oracleRule []interface{}
	for _, oracleItem := range oracle {
		oracleRule = append(oracleRule, oracleItem)
	}

	logs, sub, err := _Ens.contract.FilterLogs(opts, "NewPriceOracle", oracleRule)
	if err != nil {
		return nil, err
	}
	return &EnsNewPriceOracleIterator{contract: _Ens.contract, event: "NewPriceOracle", logs: logs, sub: sub}, nil
}

// WatchNewPriceOracle is a free log subscription operation binding the contract event 0xf261845a790fe29bbd6631e2ca4a5bdc83e6eed7c3271d9590d97287e00e9123.
//
// Solidity: event NewPriceOracle(address indexed oracle)
func (_Ens *EnsFilterer) WatchNewPriceOracle(opts *bind.WatchOpts, sink chan<- *EnsNewPriceOracle, oracle []common.Address) (event.Subscription, error) {

	var oracleRule []interface{}
	for _, oracleItem := range oracle {
		oracleRule = append(oracleRule, oracleItem)
	}

	logs, sub, err := _Ens.contract.WatchLogs(opts, "NewPriceOracle", oracleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EnsNewPriceOracle)
				if err := _Ens.contract.UnpackLog(event, "NewPriceOracle", log); err != nil {
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

// ParseNewPriceOracle is a log parse operation binding the contract event 0xf261845a790fe29bbd6631e2ca4a5bdc83e6eed7c3271d9590d97287e00e9123.
//
// Solidity: event NewPriceOracle(address indexed oracle)
func (_Ens *EnsFilterer) ParseNewPriceOracle(log types.Log) (*EnsNewPriceOracle, error) {
	event := new(EnsNewPriceOracle)
	if err := _Ens.contract.UnpackLog(event, "NewPriceOracle", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EnsOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Ens contract.
type EnsOwnershipTransferredIterator struct {
	Event *EnsOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *EnsOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EnsOwnershipTransferred)
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
		it.Event = new(EnsOwnershipTransferred)
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
func (it *EnsOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EnsOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EnsOwnershipTransferred represents a OwnershipTransferred event raised by the Ens contract.
type EnsOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Ens *EnsFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*EnsOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Ens.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &EnsOwnershipTransferredIterator{contract: _Ens.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Ens *EnsFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *EnsOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Ens.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EnsOwnershipTransferred)
				if err := _Ens.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_Ens *EnsFilterer) ParseOwnershipTransferred(log types.Log) (*EnsOwnershipTransferred, error) {
	event := new(EnsOwnershipTransferred)
	if err := _Ens.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
