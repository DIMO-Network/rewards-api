// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contracts

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

// ReferralMetaData contains all meta data concerning the Referral contract.
var ReferralMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"InvalidArrayLength\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"previousAdmin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"AdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"destination\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"AdminWithdrawal\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"beacon\",\"type\":\"address\"}],\"name\":\"BeaconUpgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"BonusChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"referreds\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"referrers\",\"type\":\"address\"}],\"name\":\"ReferralComplete\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"ORACLE_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"adminWithdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"bonusAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dimoToken\",\"outputs\":[{\"internalType\":\"contractIDimo\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"bonusAmount_\",\"type\":\"uint256\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"proxiableUUID\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"referreds\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"referrers\",\"type\":\"address[]\"}],\"name\":\"sendReferralBonuses\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSent\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"}],\"name\":\"upgradeTo\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"upgradeToAndCall\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
}

// ReferralABI is the input ABI used to generate the binding from.
// Deprecated: Use ReferralMetaData.ABI instead.
var ReferralABI = ReferralMetaData.ABI

// Referral is an auto generated Go binding around an Ethereum contract.
type Referral struct {
	ReferralCaller     // Read-only binding to the contract
	ReferralTransactor // Write-only binding to the contract
	ReferralFilterer   // Log filterer for contract events
}

// ReferralCaller is an auto generated read-only Go binding around an Ethereum contract.
type ReferralCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ReferralTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ReferralTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ReferralFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ReferralFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ReferralSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ReferralSession struct {
	Contract     *Referral         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ReferralCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ReferralCallerSession struct {
	Contract *ReferralCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// ReferralTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ReferralTransactorSession struct {
	Contract     *ReferralTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// ReferralRaw is an auto generated low-level Go binding around an Ethereum contract.
type ReferralRaw struct {
	Contract *Referral // Generic contract binding to access the raw methods on
}

// ReferralCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ReferralCallerRaw struct {
	Contract *ReferralCaller // Generic read-only contract binding to access the raw methods on
}

// ReferralTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ReferralTransactorRaw struct {
	Contract *ReferralTransactor // Generic write-only contract binding to access the raw methods on
}

// NewReferral creates a new instance of Referral, bound to a specific deployed contract.
func NewReferral(address common.Address, backend bind.ContractBackend) (*Referral, error) {
	contract, err := bindReferral(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Referral{ReferralCaller: ReferralCaller{contract: contract}, ReferralTransactor: ReferralTransactor{contract: contract}, ReferralFilterer: ReferralFilterer{contract: contract}}, nil
}

// NewReferralCaller creates a new read-only instance of Referral, bound to a specific deployed contract.
func NewReferralCaller(address common.Address, caller bind.ContractCaller) (*ReferralCaller, error) {
	contract, err := bindReferral(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ReferralCaller{contract: contract}, nil
}

// NewReferralTransactor creates a new write-only instance of Referral, bound to a specific deployed contract.
func NewReferralTransactor(address common.Address, transactor bind.ContractTransactor) (*ReferralTransactor, error) {
	contract, err := bindReferral(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ReferralTransactor{contract: contract}, nil
}

// NewReferralFilterer creates a new log filterer instance of Referral, bound to a specific deployed contract.
func NewReferralFilterer(address common.Address, filterer bind.ContractFilterer) (*ReferralFilterer, error) {
	contract, err := bindReferral(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ReferralFilterer{contract: contract}, nil
}

// bindReferral binds a generic wrapper to an already deployed contract.
func bindReferral(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ReferralABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Referral *ReferralRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Referral.Contract.ReferralCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Referral *ReferralRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Referral.Contract.ReferralTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Referral *ReferralRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Referral.Contract.ReferralTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Referral *ReferralCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Referral.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Referral *ReferralTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Referral.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Referral *ReferralTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Referral.Contract.contract.Transact(opts, method, params...)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_Referral *ReferralCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Referral.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_Referral *ReferralSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _Referral.Contract.DEFAULTADMINROLE(&_Referral.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_Referral *ReferralCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _Referral.Contract.DEFAULTADMINROLE(&_Referral.CallOpts)
}

// ORACLEROLE is a free data retrieval call binding the contract method 0x07e2cea5.
//
// Solidity: function ORACLE_ROLE() view returns(bytes32)
func (_Referral *ReferralCaller) ORACLEROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Referral.contract.Call(opts, &out, "ORACLE_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ORACLEROLE is a free data retrieval call binding the contract method 0x07e2cea5.
//
// Solidity: function ORACLE_ROLE() view returns(bytes32)
func (_Referral *ReferralSession) ORACLEROLE() ([32]byte, error) {
	return _Referral.Contract.ORACLEROLE(&_Referral.CallOpts)
}

// ORACLEROLE is a free data retrieval call binding the contract method 0x07e2cea5.
//
// Solidity: function ORACLE_ROLE() view returns(bytes32)
func (_Referral *ReferralCallerSession) ORACLEROLE() ([32]byte, error) {
	return _Referral.Contract.ORACLEROLE(&_Referral.CallOpts)
}

// BonusAmount is a free data retrieval call binding the contract method 0xabadaf9a.
//
// Solidity: function bonusAmount() view returns(uint256)
func (_Referral *ReferralCaller) BonusAmount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Referral.contract.Call(opts, &out, "bonusAmount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BonusAmount is a free data retrieval call binding the contract method 0xabadaf9a.
//
// Solidity: function bonusAmount() view returns(uint256)
func (_Referral *ReferralSession) BonusAmount() (*big.Int, error) {
	return _Referral.Contract.BonusAmount(&_Referral.CallOpts)
}

// BonusAmount is a free data retrieval call binding the contract method 0xabadaf9a.
//
// Solidity: function bonusAmount() view returns(uint256)
func (_Referral *ReferralCallerSession) BonusAmount() (*big.Int, error) {
	return _Referral.Contract.BonusAmount(&_Referral.CallOpts)
}

// DimoToken is a free data retrieval call binding the contract method 0x0524f28c.
//
// Solidity: function dimoToken() view returns(address)
func (_Referral *ReferralCaller) DimoToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Referral.contract.Call(opts, &out, "dimoToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// DimoToken is a free data retrieval call binding the contract method 0x0524f28c.
//
// Solidity: function dimoToken() view returns(address)
func (_Referral *ReferralSession) DimoToken() (common.Address, error) {
	return _Referral.Contract.DimoToken(&_Referral.CallOpts)
}

// DimoToken is a free data retrieval call binding the contract method 0x0524f28c.
//
// Solidity: function dimoToken() view returns(address)
func (_Referral *ReferralCallerSession) DimoToken() (common.Address, error) {
	return _Referral.Contract.DimoToken(&_Referral.CallOpts)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_Referral *ReferralCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _Referral.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_Referral *ReferralSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _Referral.Contract.GetRoleAdmin(&_Referral.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_Referral *ReferralCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _Referral.Contract.GetRoleAdmin(&_Referral.CallOpts, role)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_Referral *ReferralCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _Referral.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_Referral *ReferralSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _Referral.Contract.HasRole(&_Referral.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_Referral *ReferralCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _Referral.Contract.HasRole(&_Referral.CallOpts, role, account)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Referral *ReferralCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Referral.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Referral *ReferralSession) ProxiableUUID() ([32]byte, error) {
	return _Referral.Contract.ProxiableUUID(&_Referral.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Referral *ReferralCallerSession) ProxiableUUID() ([32]byte, error) {
	return _Referral.Contract.ProxiableUUID(&_Referral.CallOpts)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Referral *ReferralCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _Referral.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Referral *ReferralSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _Referral.Contract.SupportsInterface(&_Referral.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Referral *ReferralCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _Referral.Contract.SupportsInterface(&_Referral.CallOpts, interfaceId)
}

// TotalSent is a free data retrieval call binding the contract method 0x46f99063.
//
// Solidity: function totalSent() view returns(uint256)
func (_Referral *ReferralCaller) TotalSent(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Referral.contract.Call(opts, &out, "totalSent")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSent is a free data retrieval call binding the contract method 0x46f99063.
//
// Solidity: function totalSent() view returns(uint256)
func (_Referral *ReferralSession) TotalSent() (*big.Int, error) {
	return _Referral.Contract.TotalSent(&_Referral.CallOpts)
}

// TotalSent is a free data retrieval call binding the contract method 0x46f99063.
//
// Solidity: function totalSent() view returns(uint256)
func (_Referral *ReferralCallerSession) TotalSent() (*big.Int, error) {
	return _Referral.Contract.TotalSent(&_Referral.CallOpts)
}

// AdminWithdraw is a paid mutator transaction binding the contract method 0x401d4482.
//
// Solidity: function adminWithdraw(address user, uint256 amount) returns()
func (_Referral *ReferralTransactor) AdminWithdraw(opts *bind.TransactOpts, user common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Referral.contract.Transact(opts, "adminWithdraw", user, amount)
}

// AdminWithdraw is a paid mutator transaction binding the contract method 0x401d4482.
//
// Solidity: function adminWithdraw(address user, uint256 amount) returns()
func (_Referral *ReferralSession) AdminWithdraw(user common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Referral.Contract.AdminWithdraw(&_Referral.TransactOpts, user, amount)
}

// AdminWithdraw is a paid mutator transaction binding the contract method 0x401d4482.
//
// Solidity: function adminWithdraw(address user, uint256 amount) returns()
func (_Referral *ReferralTransactorSession) AdminWithdraw(user common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Referral.Contract.AdminWithdraw(&_Referral.TransactOpts, user, amount)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_Referral *ReferralTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Referral.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_Referral *ReferralSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Referral.Contract.GrantRole(&_Referral.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_Referral *ReferralTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Referral.Contract.GrantRole(&_Referral.TransactOpts, role, account)
}

// Initialize is a paid mutator transaction binding the contract method 0xcd6dc687.
//
// Solidity: function initialize(address tokenAddress, uint256 bonusAmount_) returns()
func (_Referral *ReferralTransactor) Initialize(opts *bind.TransactOpts, tokenAddress common.Address, bonusAmount_ *big.Int) (*types.Transaction, error) {
	return _Referral.contract.Transact(opts, "initialize", tokenAddress, bonusAmount_)
}

// Initialize is a paid mutator transaction binding the contract method 0xcd6dc687.
//
// Solidity: function initialize(address tokenAddress, uint256 bonusAmount_) returns()
func (_Referral *ReferralSession) Initialize(tokenAddress common.Address, bonusAmount_ *big.Int) (*types.Transaction, error) {
	return _Referral.Contract.Initialize(&_Referral.TransactOpts, tokenAddress, bonusAmount_)
}

// Initialize is a paid mutator transaction binding the contract method 0xcd6dc687.
//
// Solidity: function initialize(address tokenAddress, uint256 bonusAmount_) returns()
func (_Referral *ReferralTransactorSession) Initialize(tokenAddress common.Address, bonusAmount_ *big.Int) (*types.Transaction, error) {
	return _Referral.Contract.Initialize(&_Referral.TransactOpts, tokenAddress, bonusAmount_)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_Referral *ReferralTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Referral.contract.Transact(opts, "renounceRole", role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_Referral *ReferralSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Referral.Contract.RenounceRole(&_Referral.TransactOpts, role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_Referral *ReferralTransactorSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Referral.Contract.RenounceRole(&_Referral.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_Referral *ReferralTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Referral.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_Referral *ReferralSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Referral.Contract.RevokeRole(&_Referral.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_Referral *ReferralTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Referral.Contract.RevokeRole(&_Referral.TransactOpts, role, account)
}

// SendReferralBonuses is a paid mutator transaction binding the contract method 0xa6703aa1.
//
// Solidity: function sendReferralBonuses(address[] referreds, address[] referrers) returns()
func (_Referral *ReferralTransactor) SendReferralBonuses(opts *bind.TransactOpts, referreds []common.Address, referrers []common.Address) (*types.Transaction, error) {
	return _Referral.contract.Transact(opts, "sendReferralBonuses", referreds, referrers)
}

// SendReferralBonuses is a paid mutator transaction binding the contract method 0xa6703aa1.
//
// Solidity: function sendReferralBonuses(address[] referreds, address[] referrers) returns()
func (_Referral *ReferralSession) SendReferralBonuses(referreds []common.Address, referrers []common.Address) (*types.Transaction, error) {
	return _Referral.Contract.SendReferralBonuses(&_Referral.TransactOpts, referreds, referrers)
}

// SendReferralBonuses is a paid mutator transaction binding the contract method 0xa6703aa1.
//
// Solidity: function sendReferralBonuses(address[] referreds, address[] referrers) returns()
func (_Referral *ReferralTransactorSession) SendReferralBonuses(referreds []common.Address, referrers []common.Address) (*types.Transaction, error) {
	return _Referral.Contract.SendReferralBonuses(&_Referral.TransactOpts, referreds, referrers)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_Referral *ReferralTransactor) UpgradeTo(opts *bind.TransactOpts, newImplementation common.Address) (*types.Transaction, error) {
	return _Referral.contract.Transact(opts, "upgradeTo", newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_Referral *ReferralSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _Referral.Contract.UpgradeTo(&_Referral.TransactOpts, newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_Referral *ReferralTransactorSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _Referral.Contract.UpgradeTo(&_Referral.TransactOpts, newImplementation)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Referral *ReferralTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Referral.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Referral *ReferralSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Referral.Contract.UpgradeToAndCall(&_Referral.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Referral *ReferralTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Referral.Contract.UpgradeToAndCall(&_Referral.TransactOpts, newImplementation, data)
}

// ReferralAdminChangedIterator is returned from FilterAdminChanged and is used to iterate over the raw logs and unpacked data for AdminChanged events raised by the Referral contract.
type ReferralAdminChangedIterator struct {
	Event *ReferralAdminChanged // Event containing the contract specifics and raw log

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
func (it *ReferralAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ReferralAdminChanged)
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
		it.Event = new(ReferralAdminChanged)
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
func (it *ReferralAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ReferralAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ReferralAdminChanged represents a AdminChanged event raised by the Referral contract.
type ReferralAdminChanged struct {
	PreviousAdmin common.Address
	NewAdmin      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterAdminChanged is a free log retrieval operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_Referral *ReferralFilterer) FilterAdminChanged(opts *bind.FilterOpts) (*ReferralAdminChangedIterator, error) {

	logs, sub, err := _Referral.contract.FilterLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return &ReferralAdminChangedIterator{contract: _Referral.contract, event: "AdminChanged", logs: logs, sub: sub}, nil
}

// WatchAdminChanged is a free log subscription operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_Referral *ReferralFilterer) WatchAdminChanged(opts *bind.WatchOpts, sink chan<- *ReferralAdminChanged) (event.Subscription, error) {

	logs, sub, err := _Referral.contract.WatchLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ReferralAdminChanged)
				if err := _Referral.contract.UnpackLog(event, "AdminChanged", log); err != nil {
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

// ParseAdminChanged is a log parse operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_Referral *ReferralFilterer) ParseAdminChanged(log types.Log) (*ReferralAdminChanged, error) {
	event := new(ReferralAdminChanged)
	if err := _Referral.contract.UnpackLog(event, "AdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ReferralAdminWithdrawalIterator is returned from FilterAdminWithdrawal and is used to iterate over the raw logs and unpacked data for AdminWithdrawal events raised by the Referral contract.
type ReferralAdminWithdrawalIterator struct {
	Event *ReferralAdminWithdrawal // Event containing the contract specifics and raw log

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
func (it *ReferralAdminWithdrawalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ReferralAdminWithdrawal)
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
		it.Event = new(ReferralAdminWithdrawal)
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
func (it *ReferralAdminWithdrawalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ReferralAdminWithdrawalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ReferralAdminWithdrawal represents a AdminWithdrawal event raised by the Referral contract.
type ReferralAdminWithdrawal struct {
	Destination common.Address
	Amount      *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterAdminWithdrawal is a free log retrieval operation binding the contract event 0x1f29bc8239df330207e019f41493b485f9c7d3ce83a795ae64603dde527ada2e.
//
// Solidity: event AdminWithdrawal(address indexed destination, uint256 amount)
func (_Referral *ReferralFilterer) FilterAdminWithdrawal(opts *bind.FilterOpts, destination []common.Address) (*ReferralAdminWithdrawalIterator, error) {

	var destinationRule []interface{}
	for _, destinationItem := range destination {
		destinationRule = append(destinationRule, destinationItem)
	}

	logs, sub, err := _Referral.contract.FilterLogs(opts, "AdminWithdrawal", destinationRule)
	if err != nil {
		return nil, err
	}
	return &ReferralAdminWithdrawalIterator{contract: _Referral.contract, event: "AdminWithdrawal", logs: logs, sub: sub}, nil
}

// WatchAdminWithdrawal is a free log subscription operation binding the contract event 0x1f29bc8239df330207e019f41493b485f9c7d3ce83a795ae64603dde527ada2e.
//
// Solidity: event AdminWithdrawal(address indexed destination, uint256 amount)
func (_Referral *ReferralFilterer) WatchAdminWithdrawal(opts *bind.WatchOpts, sink chan<- *ReferralAdminWithdrawal, destination []common.Address) (event.Subscription, error) {

	var destinationRule []interface{}
	for _, destinationItem := range destination {
		destinationRule = append(destinationRule, destinationItem)
	}

	logs, sub, err := _Referral.contract.WatchLogs(opts, "AdminWithdrawal", destinationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ReferralAdminWithdrawal)
				if err := _Referral.contract.UnpackLog(event, "AdminWithdrawal", log); err != nil {
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

// ParseAdminWithdrawal is a log parse operation binding the contract event 0x1f29bc8239df330207e019f41493b485f9c7d3ce83a795ae64603dde527ada2e.
//
// Solidity: event AdminWithdrawal(address indexed destination, uint256 amount)
func (_Referral *ReferralFilterer) ParseAdminWithdrawal(log types.Log) (*ReferralAdminWithdrawal, error) {
	event := new(ReferralAdminWithdrawal)
	if err := _Referral.contract.UnpackLog(event, "AdminWithdrawal", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ReferralBeaconUpgradedIterator is returned from FilterBeaconUpgraded and is used to iterate over the raw logs and unpacked data for BeaconUpgraded events raised by the Referral contract.
type ReferralBeaconUpgradedIterator struct {
	Event *ReferralBeaconUpgraded // Event containing the contract specifics and raw log

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
func (it *ReferralBeaconUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ReferralBeaconUpgraded)
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
		it.Event = new(ReferralBeaconUpgraded)
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
func (it *ReferralBeaconUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ReferralBeaconUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ReferralBeaconUpgraded represents a BeaconUpgraded event raised by the Referral contract.
type ReferralBeaconUpgraded struct {
	Beacon common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBeaconUpgraded is a free log retrieval operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_Referral *ReferralFilterer) FilterBeaconUpgraded(opts *bind.FilterOpts, beacon []common.Address) (*ReferralBeaconUpgradedIterator, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _Referral.contract.FilterLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return &ReferralBeaconUpgradedIterator{contract: _Referral.contract, event: "BeaconUpgraded", logs: logs, sub: sub}, nil
}

// WatchBeaconUpgraded is a free log subscription operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_Referral *ReferralFilterer) WatchBeaconUpgraded(opts *bind.WatchOpts, sink chan<- *ReferralBeaconUpgraded, beacon []common.Address) (event.Subscription, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _Referral.contract.WatchLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ReferralBeaconUpgraded)
				if err := _Referral.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
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

// ParseBeaconUpgraded is a log parse operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_Referral *ReferralFilterer) ParseBeaconUpgraded(log types.Log) (*ReferralBeaconUpgraded, error) {
	event := new(ReferralBeaconUpgraded)
	if err := _Referral.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ReferralBonusChangedIterator is returned from FilterBonusChanged and is used to iterate over the raw logs and unpacked data for BonusChanged events raised by the Referral contract.
type ReferralBonusChangedIterator struct {
	Event *ReferralBonusChanged // Event containing the contract specifics and raw log

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
func (it *ReferralBonusChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ReferralBonusChanged)
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
		it.Event = new(ReferralBonusChanged)
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
func (it *ReferralBonusChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ReferralBonusChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ReferralBonusChanged represents a BonusChanged event raised by the Referral contract.
type ReferralBonusChanged struct {
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBonusChanged is a free log retrieval operation binding the contract event 0x3004fd9893f9a32bdc520978802cadc651d89c5d9a6235f3080819aa599903b8.
//
// Solidity: event BonusChanged(uint256 amount)
func (_Referral *ReferralFilterer) FilterBonusChanged(opts *bind.FilterOpts) (*ReferralBonusChangedIterator, error) {

	logs, sub, err := _Referral.contract.FilterLogs(opts, "BonusChanged")
	if err != nil {
		return nil, err
	}
	return &ReferralBonusChangedIterator{contract: _Referral.contract, event: "BonusChanged", logs: logs, sub: sub}, nil
}

// WatchBonusChanged is a free log subscription operation binding the contract event 0x3004fd9893f9a32bdc520978802cadc651d89c5d9a6235f3080819aa599903b8.
//
// Solidity: event BonusChanged(uint256 amount)
func (_Referral *ReferralFilterer) WatchBonusChanged(opts *bind.WatchOpts, sink chan<- *ReferralBonusChanged) (event.Subscription, error) {

	logs, sub, err := _Referral.contract.WatchLogs(opts, "BonusChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ReferralBonusChanged)
				if err := _Referral.contract.UnpackLog(event, "BonusChanged", log); err != nil {
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

// ParseBonusChanged is a log parse operation binding the contract event 0x3004fd9893f9a32bdc520978802cadc651d89c5d9a6235f3080819aa599903b8.
//
// Solidity: event BonusChanged(uint256 amount)
func (_Referral *ReferralFilterer) ParseBonusChanged(log types.Log) (*ReferralBonusChanged, error) {
	event := new(ReferralBonusChanged)
	if err := _Referral.contract.UnpackLog(event, "BonusChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ReferralInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the Referral contract.
type ReferralInitializedIterator struct {
	Event *ReferralInitialized // Event containing the contract specifics and raw log

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
func (it *ReferralInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ReferralInitialized)
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
		it.Event = new(ReferralInitialized)
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
func (it *ReferralInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ReferralInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ReferralInitialized represents a Initialized event raised by the Referral contract.
type ReferralInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Referral *ReferralFilterer) FilterInitialized(opts *bind.FilterOpts) (*ReferralInitializedIterator, error) {

	logs, sub, err := _Referral.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &ReferralInitializedIterator{contract: _Referral.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Referral *ReferralFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *ReferralInitialized) (event.Subscription, error) {

	logs, sub, err := _Referral.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ReferralInitialized)
				if err := _Referral.contract.UnpackLog(event, "Initialized", log); err != nil {
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

// ParseInitialized is a log parse operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Referral *ReferralFilterer) ParseInitialized(log types.Log) (*ReferralInitialized, error) {
	event := new(ReferralInitialized)
	if err := _Referral.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ReferralReferralCompleteIterator is returned from FilterReferralComplete and is used to iterate over the raw logs and unpacked data for ReferralComplete events raised by the Referral contract.
type ReferralReferralCompleteIterator struct {
	Event *ReferralReferralComplete // Event containing the contract specifics and raw log

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
func (it *ReferralReferralCompleteIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ReferralReferralComplete)
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
		it.Event = new(ReferralReferralComplete)
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
func (it *ReferralReferralCompleteIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ReferralReferralCompleteIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ReferralReferralComplete represents a ReferralComplete event raised by the Referral contract.
type ReferralReferralComplete struct {
	Referreds common.Address
	Referrers common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterReferralComplete is a free log retrieval operation binding the contract event 0x140ec1de14f412a3a1e4295313d7c4190785ddf8930efee41722cd36f00ebb91.
//
// Solidity: event ReferralComplete(address indexed referreds, address indexed referrers)
func (_Referral *ReferralFilterer) FilterReferralComplete(opts *bind.FilterOpts, referreds []common.Address, referrers []common.Address) (*ReferralReferralCompleteIterator, error) {

	var referredsRule []interface{}
	for _, referredsItem := range referreds {
		referredsRule = append(referredsRule, referredsItem)
	}
	var referrersRule []interface{}
	for _, referrersItem := range referrers {
		referrersRule = append(referrersRule, referrersItem)
	}

	logs, sub, err := _Referral.contract.FilterLogs(opts, "ReferralComplete", referredsRule, referrersRule)
	if err != nil {
		return nil, err
	}
	return &ReferralReferralCompleteIterator{contract: _Referral.contract, event: "ReferralComplete", logs: logs, sub: sub}, nil
}

// WatchReferralComplete is a free log subscription operation binding the contract event 0x140ec1de14f412a3a1e4295313d7c4190785ddf8930efee41722cd36f00ebb91.
//
// Solidity: event ReferralComplete(address indexed referreds, address indexed referrers)
func (_Referral *ReferralFilterer) WatchReferralComplete(opts *bind.WatchOpts, sink chan<- *ReferralReferralComplete, referreds []common.Address, referrers []common.Address) (event.Subscription, error) {

	var referredsRule []interface{}
	for _, referredsItem := range referreds {
		referredsRule = append(referredsRule, referredsItem)
	}
	var referrersRule []interface{}
	for _, referrersItem := range referrers {
		referrersRule = append(referrersRule, referrersItem)
	}

	logs, sub, err := _Referral.contract.WatchLogs(opts, "ReferralComplete", referredsRule, referrersRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ReferralReferralComplete)
				if err := _Referral.contract.UnpackLog(event, "ReferralComplete", log); err != nil {
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

// ParseReferralComplete is a log parse operation binding the contract event 0x140ec1de14f412a3a1e4295313d7c4190785ddf8930efee41722cd36f00ebb91.
//
// Solidity: event ReferralComplete(address indexed referreds, address indexed referrers)
func (_Referral *ReferralFilterer) ParseReferralComplete(log types.Log) (*ReferralReferralComplete, error) {
	event := new(ReferralReferralComplete)
	if err := _Referral.contract.UnpackLog(event, "ReferralComplete", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ReferralRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the Referral contract.
type ReferralRoleAdminChangedIterator struct {
	Event *ReferralRoleAdminChanged // Event containing the contract specifics and raw log

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
func (it *ReferralRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ReferralRoleAdminChanged)
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
		it.Event = new(ReferralRoleAdminChanged)
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
func (it *ReferralRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ReferralRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ReferralRoleAdminChanged represents a RoleAdminChanged event raised by the Referral contract.
type ReferralRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_Referral *ReferralFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*ReferralRoleAdminChangedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _Referral.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &ReferralRoleAdminChangedIterator{contract: _Referral.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_Referral *ReferralFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *ReferralRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _Referral.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ReferralRoleAdminChanged)
				if err := _Referral.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
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

// ParseRoleAdminChanged is a log parse operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_Referral *ReferralFilterer) ParseRoleAdminChanged(log types.Log) (*ReferralRoleAdminChanged, error) {
	event := new(ReferralRoleAdminChanged)
	if err := _Referral.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ReferralRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the Referral contract.
type ReferralRoleGrantedIterator struct {
	Event *ReferralRoleGranted // Event containing the contract specifics and raw log

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
func (it *ReferralRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ReferralRoleGranted)
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
		it.Event = new(ReferralRoleGranted)
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
func (it *ReferralRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ReferralRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ReferralRoleGranted represents a RoleGranted event raised by the Referral contract.
type ReferralRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_Referral *ReferralFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*ReferralRoleGrantedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _Referral.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &ReferralRoleGrantedIterator{contract: _Referral.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_Referral *ReferralFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *ReferralRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _Referral.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ReferralRoleGranted)
				if err := _Referral.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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

// ParseRoleGranted is a log parse operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_Referral *ReferralFilterer) ParseRoleGranted(log types.Log) (*ReferralRoleGranted, error) {
	event := new(ReferralRoleGranted)
	if err := _Referral.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ReferralRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the Referral contract.
type ReferralRoleRevokedIterator struct {
	Event *ReferralRoleRevoked // Event containing the contract specifics and raw log

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
func (it *ReferralRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ReferralRoleRevoked)
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
		it.Event = new(ReferralRoleRevoked)
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
func (it *ReferralRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ReferralRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ReferralRoleRevoked represents a RoleRevoked event raised by the Referral contract.
type ReferralRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_Referral *ReferralFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*ReferralRoleRevokedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _Referral.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &ReferralRoleRevokedIterator{contract: _Referral.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_Referral *ReferralFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *ReferralRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _Referral.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ReferralRoleRevoked)
				if err := _Referral.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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

// ParseRoleRevoked is a log parse operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_Referral *ReferralFilterer) ParseRoleRevoked(log types.Log) (*ReferralRoleRevoked, error) {
	event := new(ReferralRoleRevoked)
	if err := _Referral.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ReferralUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the Referral contract.
type ReferralUpgradedIterator struct {
	Event *ReferralUpgraded // Event containing the contract specifics and raw log

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
func (it *ReferralUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ReferralUpgraded)
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
		it.Event = new(ReferralUpgraded)
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
func (it *ReferralUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ReferralUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ReferralUpgraded represents a Upgraded event raised by the Referral contract.
type ReferralUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_Referral *ReferralFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*ReferralUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _Referral.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &ReferralUpgradedIterator{contract: _Referral.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_Referral *ReferralFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *ReferralUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _Referral.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ReferralUpgraded)
				if err := _Referral.contract.UnpackLog(event, "Upgraded", log); err != nil {
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

// ParseUpgraded is a log parse operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_Referral *ReferralFilterer) ParseUpgraded(log types.Log) (*ReferralUpgraded, error) {
	event := new(ReferralUpgraded)
	if err := _Referral.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
