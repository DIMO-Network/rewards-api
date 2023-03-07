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

// ReferralsMetaData contains all meta data concerning the Referrals contract.
var ReferralsMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"InvalidArrayLength\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"previousAdmin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"AdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"destination\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"AdminWithdrawal\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"beacon\",\"type\":\"address\"}],\"name\":\"BeaconUpgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"BonusChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"referred\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"referrer\",\"type\":\"address\"}],\"name\":\"ReferralComplete\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"referred\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"referrer\",\"type\":\"address\"}],\"name\":\"ReferralInvalid\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"ORACLE_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"adminWithdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"bonusAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dimoToken\",\"outputs\":[{\"internalType\":\"contractIDimo\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"bonusAmount_\",\"type\":\"uint256\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"proxiableUUID\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"referreds\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"referrers\",\"type\":\"address[]\"}],\"name\":\"sendReferralBonuses\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSent\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"}],\"name\":\"upgradeTo\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"upgradeToAndCall\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
}

// ReferralsABI is the input ABI used to generate the binding from.
// Deprecated: Use ReferralsMetaData.ABI instead.
var ReferralsABI = ReferralsMetaData.ABI

// Referrals is an auto generated Go binding around an Ethereum contract.
type Referrals struct {
	ReferralsCaller     // Read-only binding to the contract
	ReferralsTransactor // Write-only binding to the contract
	ReferralsFilterer   // Log filterer for contract events
}

// ReferralsCaller is an auto generated read-only Go binding around an Ethereum contract.
type ReferralsCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ReferralsTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ReferralsTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ReferralsFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ReferralsFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ReferralsSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ReferralsSession struct {
	Contract     *Referrals        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ReferralsCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ReferralsCallerSession struct {
	Contract *ReferralsCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// ReferralsTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ReferralsTransactorSession struct {
	Contract     *ReferralsTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// ReferralsRaw is an auto generated low-level Go binding around an Ethereum contract.
type ReferralsRaw struct {
	Contract *Referrals // Generic contract binding to access the raw methods on
}

// ReferralsCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ReferralsCallerRaw struct {
	Contract *ReferralsCaller // Generic read-only contract binding to access the raw methods on
}

// ReferralsTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ReferralsTransactorRaw struct {
	Contract *ReferralsTransactor // Generic write-only contract binding to access the raw methods on
}

// NewReferrals creates a new instance of Referrals, bound to a specific deployed contract.
func NewReferrals(address common.Address, backend bind.ContractBackend) (*Referrals, error) {
	contract, err := bindReferrals(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Referrals{ReferralsCaller: ReferralsCaller{contract: contract}, ReferralsTransactor: ReferralsTransactor{contract: contract}, ReferralsFilterer: ReferralsFilterer{contract: contract}}, nil
}

// NewReferralsCaller creates a new read-only instance of Referrals, bound to a specific deployed contract.
func NewReferralsCaller(address common.Address, caller bind.ContractCaller) (*ReferralsCaller, error) {
	contract, err := bindReferrals(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ReferralsCaller{contract: contract}, nil
}

// NewReferralsTransactor creates a new write-only instance of Referrals, bound to a specific deployed contract.
func NewReferralsTransactor(address common.Address, transactor bind.ContractTransactor) (*ReferralsTransactor, error) {
	contract, err := bindReferrals(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ReferralsTransactor{contract: contract}, nil
}

// NewReferralsFilterer creates a new log filterer instance of Referrals, bound to a specific deployed contract.
func NewReferralsFilterer(address common.Address, filterer bind.ContractFilterer) (*ReferralsFilterer, error) {
	contract, err := bindReferrals(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ReferralsFilterer{contract: contract}, nil
}

// bindReferrals binds a generic wrapper to an already deployed contract.
func bindReferrals(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ReferralsABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Referrals *ReferralsRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Referrals.Contract.ReferralsCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Referrals *ReferralsRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Referrals.Contract.ReferralsTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Referrals *ReferralsRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Referrals.Contract.ReferralsTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Referrals *ReferralsCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Referrals.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Referrals *ReferralsTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Referrals.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Referrals *ReferralsTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Referrals.Contract.contract.Transact(opts, method, params...)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_Referrals *ReferralsCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Referrals.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_Referrals *ReferralsSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _Referrals.Contract.DEFAULTADMINROLE(&_Referrals.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_Referrals *ReferralsCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _Referrals.Contract.DEFAULTADMINROLE(&_Referrals.CallOpts)
}

// ORACLEROLE is a free data retrieval call binding the contract method 0x07e2cea5.
//
// Solidity: function ORACLE_ROLE() view returns(bytes32)
func (_Referrals *ReferralsCaller) ORACLEROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Referrals.contract.Call(opts, &out, "ORACLE_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ORACLEROLE is a free data retrieval call binding the contract method 0x07e2cea5.
//
// Solidity: function ORACLE_ROLE() view returns(bytes32)
func (_Referrals *ReferralsSession) ORACLEROLE() ([32]byte, error) {
	return _Referrals.Contract.ORACLEROLE(&_Referrals.CallOpts)
}

// ORACLEROLE is a free data retrieval call binding the contract method 0x07e2cea5.
//
// Solidity: function ORACLE_ROLE() view returns(bytes32)
func (_Referrals *ReferralsCallerSession) ORACLEROLE() ([32]byte, error) {
	return _Referrals.Contract.ORACLEROLE(&_Referrals.CallOpts)
}

// BonusAmount is a free data retrieval call binding the contract method 0xabadaf9a.
//
// Solidity: function bonusAmount() view returns(uint256)
func (_Referrals *ReferralsCaller) BonusAmount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Referrals.contract.Call(opts, &out, "bonusAmount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BonusAmount is a free data retrieval call binding the contract method 0xabadaf9a.
//
// Solidity: function bonusAmount() view returns(uint256)
func (_Referrals *ReferralsSession) BonusAmount() (*big.Int, error) {
	return _Referrals.Contract.BonusAmount(&_Referrals.CallOpts)
}

// BonusAmount is a free data retrieval call binding the contract method 0xabadaf9a.
//
// Solidity: function bonusAmount() view returns(uint256)
func (_Referrals *ReferralsCallerSession) BonusAmount() (*big.Int, error) {
	return _Referrals.Contract.BonusAmount(&_Referrals.CallOpts)
}

// DimoToken is a free data retrieval call binding the contract method 0x0524f28c.
//
// Solidity: function dimoToken() view returns(address)
func (_Referrals *ReferralsCaller) DimoToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Referrals.contract.Call(opts, &out, "dimoToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// DimoToken is a free data retrieval call binding the contract method 0x0524f28c.
//
// Solidity: function dimoToken() view returns(address)
func (_Referrals *ReferralsSession) DimoToken() (common.Address, error) {
	return _Referrals.Contract.DimoToken(&_Referrals.CallOpts)
}

// DimoToken is a free data retrieval call binding the contract method 0x0524f28c.
//
// Solidity: function dimoToken() view returns(address)
func (_Referrals *ReferralsCallerSession) DimoToken() (common.Address, error) {
	return _Referrals.Contract.DimoToken(&_Referrals.CallOpts)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_Referrals *ReferralsCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _Referrals.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_Referrals *ReferralsSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _Referrals.Contract.GetRoleAdmin(&_Referrals.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_Referrals *ReferralsCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _Referrals.Contract.GetRoleAdmin(&_Referrals.CallOpts, role)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_Referrals *ReferralsCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _Referrals.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_Referrals *ReferralsSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _Referrals.Contract.HasRole(&_Referrals.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_Referrals *ReferralsCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _Referrals.Contract.HasRole(&_Referrals.CallOpts, role, account)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Referrals *ReferralsCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Referrals.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Referrals *ReferralsSession) ProxiableUUID() ([32]byte, error) {
	return _Referrals.Contract.ProxiableUUID(&_Referrals.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Referrals *ReferralsCallerSession) ProxiableUUID() ([32]byte, error) {
	return _Referrals.Contract.ProxiableUUID(&_Referrals.CallOpts)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Referrals *ReferralsCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _Referrals.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Referrals *ReferralsSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _Referrals.Contract.SupportsInterface(&_Referrals.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Referrals *ReferralsCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _Referrals.Contract.SupportsInterface(&_Referrals.CallOpts, interfaceId)
}

// TotalSent is a free data retrieval call binding the contract method 0x46f99063.
//
// Solidity: function totalSent() view returns(uint256)
func (_Referrals *ReferralsCaller) TotalSent(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Referrals.contract.Call(opts, &out, "totalSent")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSent is a free data retrieval call binding the contract method 0x46f99063.
//
// Solidity: function totalSent() view returns(uint256)
func (_Referrals *ReferralsSession) TotalSent() (*big.Int, error) {
	return _Referrals.Contract.TotalSent(&_Referrals.CallOpts)
}

// TotalSent is a free data retrieval call binding the contract method 0x46f99063.
//
// Solidity: function totalSent() view returns(uint256)
func (_Referrals *ReferralsCallerSession) TotalSent() (*big.Int, error) {
	return _Referrals.Contract.TotalSent(&_Referrals.CallOpts)
}

// AdminWithdraw is a paid mutator transaction binding the contract method 0x401d4482.
//
// Solidity: function adminWithdraw(address user, uint256 amount) returns()
func (_Referrals *ReferralsTransactor) AdminWithdraw(opts *bind.TransactOpts, user common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Referrals.contract.Transact(opts, "adminWithdraw", user, amount)
}

// AdminWithdraw is a paid mutator transaction binding the contract method 0x401d4482.
//
// Solidity: function adminWithdraw(address user, uint256 amount) returns()
func (_Referrals *ReferralsSession) AdminWithdraw(user common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Referrals.Contract.AdminWithdraw(&_Referrals.TransactOpts, user, amount)
}

// AdminWithdraw is a paid mutator transaction binding the contract method 0x401d4482.
//
// Solidity: function adminWithdraw(address user, uint256 amount) returns()
func (_Referrals *ReferralsTransactorSession) AdminWithdraw(user common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Referrals.Contract.AdminWithdraw(&_Referrals.TransactOpts, user, amount)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_Referrals *ReferralsTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Referrals.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_Referrals *ReferralsSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Referrals.Contract.GrantRole(&_Referrals.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_Referrals *ReferralsTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Referrals.Contract.GrantRole(&_Referrals.TransactOpts, role, account)
}

// Initialize is a paid mutator transaction binding the contract method 0xcd6dc687.
//
// Solidity: function initialize(address tokenAddress, uint256 bonusAmount_) returns()
func (_Referrals *ReferralsTransactor) Initialize(opts *bind.TransactOpts, tokenAddress common.Address, bonusAmount_ *big.Int) (*types.Transaction, error) {
	return _Referrals.contract.Transact(opts, "initialize", tokenAddress, bonusAmount_)
}

// Initialize is a paid mutator transaction binding the contract method 0xcd6dc687.
//
// Solidity: function initialize(address tokenAddress, uint256 bonusAmount_) returns()
func (_Referrals *ReferralsSession) Initialize(tokenAddress common.Address, bonusAmount_ *big.Int) (*types.Transaction, error) {
	return _Referrals.Contract.Initialize(&_Referrals.TransactOpts, tokenAddress, bonusAmount_)
}

// Initialize is a paid mutator transaction binding the contract method 0xcd6dc687.
//
// Solidity: function initialize(address tokenAddress, uint256 bonusAmount_) returns()
func (_Referrals *ReferralsTransactorSession) Initialize(tokenAddress common.Address, bonusAmount_ *big.Int) (*types.Transaction, error) {
	return _Referrals.Contract.Initialize(&_Referrals.TransactOpts, tokenAddress, bonusAmount_)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_Referrals *ReferralsTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Referrals.contract.Transact(opts, "renounceRole", role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_Referrals *ReferralsSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Referrals.Contract.RenounceRole(&_Referrals.TransactOpts, role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_Referrals *ReferralsTransactorSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Referrals.Contract.RenounceRole(&_Referrals.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_Referrals *ReferralsTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Referrals.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_Referrals *ReferralsSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Referrals.Contract.RevokeRole(&_Referrals.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_Referrals *ReferralsTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Referrals.Contract.RevokeRole(&_Referrals.TransactOpts, role, account)
}

// SendReferralBonuses is a paid mutator transaction binding the contract method 0xa6703aa1.
//
// Solidity: function sendReferralBonuses(address[] referreds, address[] referrers) returns()
func (_Referrals *ReferralsTransactor) SendReferralBonuses(opts *bind.TransactOpts, referreds []common.Address, referrers []common.Address) (*types.Transaction, error) {
	return _Referrals.contract.Transact(opts, "sendReferralBonuses", referreds, referrers)
}

// SendReferralBonuses is a paid mutator transaction binding the contract method 0xa6703aa1.
//
// Solidity: function sendReferralBonuses(address[] referreds, address[] referrers) returns()
func (_Referrals *ReferralsSession) SendReferralBonuses(referreds []common.Address, referrers []common.Address) (*types.Transaction, error) {
	return _Referrals.Contract.SendReferralBonuses(&_Referrals.TransactOpts, referreds, referrers)
}

// SendReferralBonuses is a paid mutator transaction binding the contract method 0xa6703aa1.
//
// Solidity: function sendReferralBonuses(address[] referreds, address[] referrers) returns()
func (_Referrals *ReferralsTransactorSession) SendReferralBonuses(referreds []common.Address, referrers []common.Address) (*types.Transaction, error) {
	return _Referrals.Contract.SendReferralBonuses(&_Referrals.TransactOpts, referreds, referrers)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_Referrals *ReferralsTransactor) UpgradeTo(opts *bind.TransactOpts, newImplementation common.Address) (*types.Transaction, error) {
	return _Referrals.contract.Transact(opts, "upgradeTo", newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_Referrals *ReferralsSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _Referrals.Contract.UpgradeTo(&_Referrals.TransactOpts, newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_Referrals *ReferralsTransactorSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _Referrals.Contract.UpgradeTo(&_Referrals.TransactOpts, newImplementation)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Referrals *ReferralsTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Referrals.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Referrals *ReferralsSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Referrals.Contract.UpgradeToAndCall(&_Referrals.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Referrals *ReferralsTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Referrals.Contract.UpgradeToAndCall(&_Referrals.TransactOpts, newImplementation, data)
}

// ReferralsAdminChangedIterator is returned from FilterAdminChanged and is used to iterate over the raw logs and unpacked data for AdminChanged events raised by the Referrals contract.
type ReferralsAdminChangedIterator struct {
	Event *ReferralsAdminChanged // Event containing the contract specifics and raw log

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
func (it *ReferralsAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ReferralsAdminChanged)
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
		it.Event = new(ReferralsAdminChanged)
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
func (it *ReferralsAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ReferralsAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ReferralsAdminChanged represents a AdminChanged event raised by the Referrals contract.
type ReferralsAdminChanged struct {
	PreviousAdmin common.Address
	NewAdmin      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterAdminChanged is a free log retrieval operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_Referrals *ReferralsFilterer) FilterAdminChanged(opts *bind.FilterOpts) (*ReferralsAdminChangedIterator, error) {

	logs, sub, err := _Referrals.contract.FilterLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return &ReferralsAdminChangedIterator{contract: _Referrals.contract, event: "AdminChanged", logs: logs, sub: sub}, nil
}

// WatchAdminChanged is a free log subscription operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_Referrals *ReferralsFilterer) WatchAdminChanged(opts *bind.WatchOpts, sink chan<- *ReferralsAdminChanged) (event.Subscription, error) {

	logs, sub, err := _Referrals.contract.WatchLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ReferralsAdminChanged)
				if err := _Referrals.contract.UnpackLog(event, "AdminChanged", log); err != nil {
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
func (_Referrals *ReferralsFilterer) ParseAdminChanged(log types.Log) (*ReferralsAdminChanged, error) {
	event := new(ReferralsAdminChanged)
	if err := _Referrals.contract.UnpackLog(event, "AdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ReferralsAdminWithdrawalIterator is returned from FilterAdminWithdrawal and is used to iterate over the raw logs and unpacked data for AdminWithdrawal events raised by the Referrals contract.
type ReferralsAdminWithdrawalIterator struct {
	Event *ReferralsAdminWithdrawal // Event containing the contract specifics and raw log

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
func (it *ReferralsAdminWithdrawalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ReferralsAdminWithdrawal)
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
		it.Event = new(ReferralsAdminWithdrawal)
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
func (it *ReferralsAdminWithdrawalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ReferralsAdminWithdrawalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ReferralsAdminWithdrawal represents a AdminWithdrawal event raised by the Referrals contract.
type ReferralsAdminWithdrawal struct {
	Destination common.Address
	Amount      *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterAdminWithdrawal is a free log retrieval operation binding the contract event 0x1f29bc8239df330207e019f41493b485f9c7d3ce83a795ae64603dde527ada2e.
//
// Solidity: event AdminWithdrawal(address indexed destination, uint256 amount)
func (_Referrals *ReferralsFilterer) FilterAdminWithdrawal(opts *bind.FilterOpts, destination []common.Address) (*ReferralsAdminWithdrawalIterator, error) {

	var destinationRule []interface{}
	for _, destinationItem := range destination {
		destinationRule = append(destinationRule, destinationItem)
	}

	logs, sub, err := _Referrals.contract.FilterLogs(opts, "AdminWithdrawal", destinationRule)
	if err != nil {
		return nil, err
	}
	return &ReferralsAdminWithdrawalIterator{contract: _Referrals.contract, event: "AdminWithdrawal", logs: logs, sub: sub}, nil
}

// WatchAdminWithdrawal is a free log subscription operation binding the contract event 0x1f29bc8239df330207e019f41493b485f9c7d3ce83a795ae64603dde527ada2e.
//
// Solidity: event AdminWithdrawal(address indexed destination, uint256 amount)
func (_Referrals *ReferralsFilterer) WatchAdminWithdrawal(opts *bind.WatchOpts, sink chan<- *ReferralsAdminWithdrawal, destination []common.Address) (event.Subscription, error) {

	var destinationRule []interface{}
	for _, destinationItem := range destination {
		destinationRule = append(destinationRule, destinationItem)
	}

	logs, sub, err := _Referrals.contract.WatchLogs(opts, "AdminWithdrawal", destinationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ReferralsAdminWithdrawal)
				if err := _Referrals.contract.UnpackLog(event, "AdminWithdrawal", log); err != nil {
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
func (_Referrals *ReferralsFilterer) ParseAdminWithdrawal(log types.Log) (*ReferralsAdminWithdrawal, error) {
	event := new(ReferralsAdminWithdrawal)
	if err := _Referrals.contract.UnpackLog(event, "AdminWithdrawal", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ReferralsBeaconUpgradedIterator is returned from FilterBeaconUpgraded and is used to iterate over the raw logs and unpacked data for BeaconUpgraded events raised by the Referrals contract.
type ReferralsBeaconUpgradedIterator struct {
	Event *ReferralsBeaconUpgraded // Event containing the contract specifics and raw log

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
func (it *ReferralsBeaconUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ReferralsBeaconUpgraded)
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
		it.Event = new(ReferralsBeaconUpgraded)
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
func (it *ReferralsBeaconUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ReferralsBeaconUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ReferralsBeaconUpgraded represents a BeaconUpgraded event raised by the Referrals contract.
type ReferralsBeaconUpgraded struct {
	Beacon common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBeaconUpgraded is a free log retrieval operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_Referrals *ReferralsFilterer) FilterBeaconUpgraded(opts *bind.FilterOpts, beacon []common.Address) (*ReferralsBeaconUpgradedIterator, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _Referrals.contract.FilterLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return &ReferralsBeaconUpgradedIterator{contract: _Referrals.contract, event: "BeaconUpgraded", logs: logs, sub: sub}, nil
}

// WatchBeaconUpgraded is a free log subscription operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_Referrals *ReferralsFilterer) WatchBeaconUpgraded(opts *bind.WatchOpts, sink chan<- *ReferralsBeaconUpgraded, beacon []common.Address) (event.Subscription, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _Referrals.contract.WatchLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ReferralsBeaconUpgraded)
				if err := _Referrals.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
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
func (_Referrals *ReferralsFilterer) ParseBeaconUpgraded(log types.Log) (*ReferralsBeaconUpgraded, error) {
	event := new(ReferralsBeaconUpgraded)
	if err := _Referrals.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ReferralsBonusChangedIterator is returned from FilterBonusChanged and is used to iterate over the raw logs and unpacked data for BonusChanged events raised by the Referrals contract.
type ReferralsBonusChangedIterator struct {
	Event *ReferralsBonusChanged // Event containing the contract specifics and raw log

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
func (it *ReferralsBonusChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ReferralsBonusChanged)
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
		it.Event = new(ReferralsBonusChanged)
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
func (it *ReferralsBonusChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ReferralsBonusChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ReferralsBonusChanged represents a BonusChanged event raised by the Referrals contract.
type ReferralsBonusChanged struct {
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBonusChanged is a free log retrieval operation binding the contract event 0x3004fd9893f9a32bdc520978802cadc651d89c5d9a6235f3080819aa599903b8.
//
// Solidity: event BonusChanged(uint256 amount)
func (_Referrals *ReferralsFilterer) FilterBonusChanged(opts *bind.FilterOpts) (*ReferralsBonusChangedIterator, error) {

	logs, sub, err := _Referrals.contract.FilterLogs(opts, "BonusChanged")
	if err != nil {
		return nil, err
	}
	return &ReferralsBonusChangedIterator{contract: _Referrals.contract, event: "BonusChanged", logs: logs, sub: sub}, nil
}

// WatchBonusChanged is a free log subscription operation binding the contract event 0x3004fd9893f9a32bdc520978802cadc651d89c5d9a6235f3080819aa599903b8.
//
// Solidity: event BonusChanged(uint256 amount)
func (_Referrals *ReferralsFilterer) WatchBonusChanged(opts *bind.WatchOpts, sink chan<- *ReferralsBonusChanged) (event.Subscription, error) {

	logs, sub, err := _Referrals.contract.WatchLogs(opts, "BonusChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ReferralsBonusChanged)
				if err := _Referrals.contract.UnpackLog(event, "BonusChanged", log); err != nil {
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
func (_Referrals *ReferralsFilterer) ParseBonusChanged(log types.Log) (*ReferralsBonusChanged, error) {
	event := new(ReferralsBonusChanged)
	if err := _Referrals.contract.UnpackLog(event, "BonusChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ReferralsInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the Referrals contract.
type ReferralsInitializedIterator struct {
	Event *ReferralsInitialized // Event containing the contract specifics and raw log

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
func (it *ReferralsInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ReferralsInitialized)
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
		it.Event = new(ReferralsInitialized)
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
func (it *ReferralsInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ReferralsInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ReferralsInitialized represents a Initialized event raised by the Referrals contract.
type ReferralsInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Referrals *ReferralsFilterer) FilterInitialized(opts *bind.FilterOpts) (*ReferralsInitializedIterator, error) {

	logs, sub, err := _Referrals.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &ReferralsInitializedIterator{contract: _Referrals.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Referrals *ReferralsFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *ReferralsInitialized) (event.Subscription, error) {

	logs, sub, err := _Referrals.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ReferralsInitialized)
				if err := _Referrals.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_Referrals *ReferralsFilterer) ParseInitialized(log types.Log) (*ReferralsInitialized, error) {
	event := new(ReferralsInitialized)
	if err := _Referrals.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ReferralsReferralCompleteIterator is returned from FilterReferralComplete and is used to iterate over the raw logs and unpacked data for ReferralComplete events raised by the Referrals contract.
type ReferralsReferralCompleteIterator struct {
	Event *ReferralsReferralComplete // Event containing the contract specifics and raw log

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
func (it *ReferralsReferralCompleteIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ReferralsReferralComplete)
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
		it.Event = new(ReferralsReferralComplete)
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
func (it *ReferralsReferralCompleteIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ReferralsReferralCompleteIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ReferralsReferralComplete represents a ReferralComplete event raised by the Referrals contract.
type ReferralsReferralComplete struct {
	Referred common.Address
	Referrer common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterReferralComplete is a free log retrieval operation binding the contract event 0x140ec1de14f412a3a1e4295313d7c4190785ddf8930efee41722cd36f00ebb91.
//
// Solidity: event ReferralComplete(address indexed referred, address indexed referrer)
func (_Referrals *ReferralsFilterer) FilterReferralComplete(opts *bind.FilterOpts, referred []common.Address, referrer []common.Address) (*ReferralsReferralCompleteIterator, error) {

	var referredRule []interface{}
	for _, referredItem := range referred {
		referredRule = append(referredRule, referredItem)
	}
	var referrerRule []interface{}
	for _, referrerItem := range referrer {
		referrerRule = append(referrerRule, referrerItem)
	}

	logs, sub, err := _Referrals.contract.FilterLogs(opts, "ReferralComplete", referredRule, referrerRule)
	if err != nil {
		return nil, err
	}
	return &ReferralsReferralCompleteIterator{contract: _Referrals.contract, event: "ReferralComplete", logs: logs, sub: sub}, nil
}

// WatchReferralComplete is a free log subscription operation binding the contract event 0x140ec1de14f412a3a1e4295313d7c4190785ddf8930efee41722cd36f00ebb91.
//
// Solidity: event ReferralComplete(address indexed referred, address indexed referrer)
func (_Referrals *ReferralsFilterer) WatchReferralComplete(opts *bind.WatchOpts, sink chan<- *ReferralsReferralComplete, referred []common.Address, referrer []common.Address) (event.Subscription, error) {

	var referredRule []interface{}
	for _, referredItem := range referred {
		referredRule = append(referredRule, referredItem)
	}
	var referrerRule []interface{}
	for _, referrerItem := range referrer {
		referrerRule = append(referrerRule, referrerItem)
	}

	logs, sub, err := _Referrals.contract.WatchLogs(opts, "ReferralComplete", referredRule, referrerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ReferralsReferralComplete)
				if err := _Referrals.contract.UnpackLog(event, "ReferralComplete", log); err != nil {
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
// Solidity: event ReferralComplete(address indexed referred, address indexed referrer)
func (_Referrals *ReferralsFilterer) ParseReferralComplete(log types.Log) (*ReferralsReferralComplete, error) {
	event := new(ReferralsReferralComplete)
	if err := _Referrals.contract.UnpackLog(event, "ReferralComplete", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ReferralsReferralInvalidIterator is returned from FilterReferralInvalid and is used to iterate over the raw logs and unpacked data for ReferralInvalid events raised by the Referrals contract.
type ReferralsReferralInvalidIterator struct {
	Event *ReferralsReferralInvalid // Event containing the contract specifics and raw log

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
func (it *ReferralsReferralInvalidIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ReferralsReferralInvalid)
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
		it.Event = new(ReferralsReferralInvalid)
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
func (it *ReferralsReferralInvalidIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ReferralsReferralInvalidIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ReferralsReferralInvalid represents a ReferralInvalid event raised by the Referrals contract.
type ReferralsReferralInvalid struct {
	Referred common.Address
	Referrer common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterReferralInvalid is a free log retrieval operation binding the contract event 0xa8129329835c7fa9561f19f5b059d91ada3feffb6a0c10e02a4e3d294fe3a1bc.
//
// Solidity: event ReferralInvalid(address indexed referred, address indexed referrer)
func (_Referrals *ReferralsFilterer) FilterReferralInvalid(opts *bind.FilterOpts, referred []common.Address, referrer []common.Address) (*ReferralsReferralInvalidIterator, error) {

	var referredRule []interface{}
	for _, referredItem := range referred {
		referredRule = append(referredRule, referredItem)
	}
	var referrerRule []interface{}
	for _, referrerItem := range referrer {
		referrerRule = append(referrerRule, referrerItem)
	}

	logs, sub, err := _Referrals.contract.FilterLogs(opts, "ReferralInvalid", referredRule, referrerRule)
	if err != nil {
		return nil, err
	}
	return &ReferralsReferralInvalidIterator{contract: _Referrals.contract, event: "ReferralInvalid", logs: logs, sub: sub}, nil
}

// WatchReferralInvalid is a free log subscription operation binding the contract event 0xa8129329835c7fa9561f19f5b059d91ada3feffb6a0c10e02a4e3d294fe3a1bc.
//
// Solidity: event ReferralInvalid(address indexed referred, address indexed referrer)
func (_Referrals *ReferralsFilterer) WatchReferralInvalid(opts *bind.WatchOpts, sink chan<- *ReferralsReferralInvalid, referred []common.Address, referrer []common.Address) (event.Subscription, error) {

	var referredRule []interface{}
	for _, referredItem := range referred {
		referredRule = append(referredRule, referredItem)
	}
	var referrerRule []interface{}
	for _, referrerItem := range referrer {
		referrerRule = append(referrerRule, referrerItem)
	}

	logs, sub, err := _Referrals.contract.WatchLogs(opts, "ReferralInvalid", referredRule, referrerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ReferralsReferralInvalid)
				if err := _Referrals.contract.UnpackLog(event, "ReferralInvalid", log); err != nil {
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

// ParseReferralInvalid is a log parse operation binding the contract event 0xa8129329835c7fa9561f19f5b059d91ada3feffb6a0c10e02a4e3d294fe3a1bc.
//
// Solidity: event ReferralInvalid(address indexed referred, address indexed referrer)
func (_Referrals *ReferralsFilterer) ParseReferralInvalid(log types.Log) (*ReferralsReferralInvalid, error) {
	event := new(ReferralsReferralInvalid)
	if err := _Referrals.contract.UnpackLog(event, "ReferralInvalid", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ReferralsRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the Referrals contract.
type ReferralsRoleAdminChangedIterator struct {
	Event *ReferralsRoleAdminChanged // Event containing the contract specifics and raw log

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
func (it *ReferralsRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ReferralsRoleAdminChanged)
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
		it.Event = new(ReferralsRoleAdminChanged)
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
func (it *ReferralsRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ReferralsRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ReferralsRoleAdminChanged represents a RoleAdminChanged event raised by the Referrals contract.
type ReferralsRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_Referrals *ReferralsFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*ReferralsRoleAdminChangedIterator, error) {

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

	logs, sub, err := _Referrals.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &ReferralsRoleAdminChangedIterator{contract: _Referrals.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_Referrals *ReferralsFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *ReferralsRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

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

	logs, sub, err := _Referrals.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ReferralsRoleAdminChanged)
				if err := _Referrals.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
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
func (_Referrals *ReferralsFilterer) ParseRoleAdminChanged(log types.Log) (*ReferralsRoleAdminChanged, error) {
	event := new(ReferralsRoleAdminChanged)
	if err := _Referrals.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ReferralsRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the Referrals contract.
type ReferralsRoleGrantedIterator struct {
	Event *ReferralsRoleGranted // Event containing the contract specifics and raw log

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
func (it *ReferralsRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ReferralsRoleGranted)
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
		it.Event = new(ReferralsRoleGranted)
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
func (it *ReferralsRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ReferralsRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ReferralsRoleGranted represents a RoleGranted event raised by the Referrals contract.
type ReferralsRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_Referrals *ReferralsFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*ReferralsRoleGrantedIterator, error) {

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

	logs, sub, err := _Referrals.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &ReferralsRoleGrantedIterator{contract: _Referrals.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_Referrals *ReferralsFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *ReferralsRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _Referrals.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ReferralsRoleGranted)
				if err := _Referrals.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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
func (_Referrals *ReferralsFilterer) ParseRoleGranted(log types.Log) (*ReferralsRoleGranted, error) {
	event := new(ReferralsRoleGranted)
	if err := _Referrals.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ReferralsRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the Referrals contract.
type ReferralsRoleRevokedIterator struct {
	Event *ReferralsRoleRevoked // Event containing the contract specifics and raw log

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
func (it *ReferralsRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ReferralsRoleRevoked)
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
		it.Event = new(ReferralsRoleRevoked)
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
func (it *ReferralsRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ReferralsRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ReferralsRoleRevoked represents a RoleRevoked event raised by the Referrals contract.
type ReferralsRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_Referrals *ReferralsFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*ReferralsRoleRevokedIterator, error) {

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

	logs, sub, err := _Referrals.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &ReferralsRoleRevokedIterator{contract: _Referrals.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_Referrals *ReferralsFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *ReferralsRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _Referrals.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ReferralsRoleRevoked)
				if err := _Referrals.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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
func (_Referrals *ReferralsFilterer) ParseRoleRevoked(log types.Log) (*ReferralsRoleRevoked, error) {
	event := new(ReferralsRoleRevoked)
	if err := _Referrals.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ReferralsUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the Referrals contract.
type ReferralsUpgradedIterator struct {
	Event *ReferralsUpgraded // Event containing the contract specifics and raw log

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
func (it *ReferralsUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ReferralsUpgraded)
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
		it.Event = new(ReferralsUpgraded)
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
func (it *ReferralsUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ReferralsUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ReferralsUpgraded represents a Upgraded event raised by the Referrals contract.
type ReferralsUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_Referrals *ReferralsFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*ReferralsUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _Referrals.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &ReferralsUpgradedIterator{contract: _Referrals.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_Referrals *ReferralsFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *ReferralsUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _Referrals.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ReferralsUpgraded)
				if err := _Referrals.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_Referrals *ReferralsFilterer) ParseUpgraded(log types.Log) (*ReferralsUpgraded, error) {
	event := new(ReferralsUpgraded)
	if err := _Referrals.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
