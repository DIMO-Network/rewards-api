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
	_ = abi.ConvertType
)

// MerkleDistributorMetaData contains all meta data concerning the MerkleDistributor contract.
var MerkleDistributorMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AllocationExceeded\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"AlreadyClaimed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ArrayLengthMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientPoolBalance\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidProof\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotPoolAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotRootSetter\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PoolDoesNotExist\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RootImmutableAfterClaims\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RootNotSet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"Sanctioned\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"SweepExceedsUncommitted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"WeeklyLimitExceeded\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddress\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAllocation\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAmount\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"previousAdmin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"AdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"beacon\",\"type\":\"address\"}],\"name\":\"BeaconUpgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"poolId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"week\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Claimed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"poolId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Funded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"poolId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"weeklyLimit\",\"type\":\"uint256\"}],\"name\":\"PoolCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"poolId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"week\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"root\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"allocation\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"proofsURI\",\"type\":\"string\"}],\"name\":\"RootSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"poolId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"setter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"}],\"name\":\"RootSetterSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"poolId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newBalance\",\"type\":\"uint256\"}],\"name\":\"Swept\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"poolId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"limit\",\"type\":\"uint256\"}],\"name\":\"WeeklyLimitSet\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"UPGRADER_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"allocations\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"poolId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"week\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes32[]\",\"name\":\"proof\",\"type\":\"bytes32[]\"}],\"name\":\"claim\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"poolId\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"weeks_\",\"type\":\"uint256[]\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256[]\",\"name\":\"amounts\",\"type\":\"uint256[]\"},{\"internalType\":\"bytes32[][]\",\"name\":\"proofs\",\"type\":\"bytes32[][]\"}],\"name\":\"claimBatch\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"claimed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"claimedTotals\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"weeklyLimit\",\"type\":\"uint256\"}],\"name\":\"createPool\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"poolId\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"poolId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"fund\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sanctionsList_\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"poolId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256[]\",\"name\":\"weeks_\",\"type\":\"uint256[]\"}],\"name\":\"isClaimedBatch\",\"outputs\":[{\"internalType\":\"bool[]\",\"name\":\"statuses\",\"type\":\"bool[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"isRootSetter\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"poolCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"pools\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"weeklyLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"committed\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"proofsURIs\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"proxiableUUID\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"roots\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"sanctionsList\",\"outputs\":[{\"internalType\":\"contractISanctionsList\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"poolId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"week\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"root\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"allocation\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"proofsURI\",\"type\":\"string\"}],\"name\":\"setRoot\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"poolId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"setter\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"}],\"name\":\"setRootSetter\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"poolId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"limit\",\"type\":\"uint256\"}],\"name\":\"setWeeklyLimit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"poolId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"sweep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"}],\"name\":\"upgradeTo\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"upgradeToAndCall\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
}

// MerkleDistributorABI is the input ABI used to generate the binding from.
// Deprecated: Use MerkleDistributorMetaData.ABI instead.
var MerkleDistributorABI = MerkleDistributorMetaData.ABI

// MerkleDistributor is an auto generated Go binding around an Ethereum contract.
type MerkleDistributor struct {
	MerkleDistributorCaller     // Read-only binding to the contract
	MerkleDistributorTransactor // Write-only binding to the contract
	MerkleDistributorFilterer   // Log filterer for contract events
}

// MerkleDistributorCaller is an auto generated read-only Go binding around an Ethereum contract.
type MerkleDistributorCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MerkleDistributorTransactor is an auto generated write-only Go binding around an Ethereum contract.
type MerkleDistributorTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MerkleDistributorFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type MerkleDistributorFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MerkleDistributorSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type MerkleDistributorSession struct {
	Contract     *MerkleDistributor // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// MerkleDistributorCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type MerkleDistributorCallerSession struct {
	Contract *MerkleDistributorCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// MerkleDistributorTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type MerkleDistributorTransactorSession struct {
	Contract     *MerkleDistributorTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// MerkleDistributorRaw is an auto generated low-level Go binding around an Ethereum contract.
type MerkleDistributorRaw struct {
	Contract *MerkleDistributor // Generic contract binding to access the raw methods on
}

// MerkleDistributorCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type MerkleDistributorCallerRaw struct {
	Contract *MerkleDistributorCaller // Generic read-only contract binding to access the raw methods on
}

// MerkleDistributorTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type MerkleDistributorTransactorRaw struct {
	Contract *MerkleDistributorTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMerkleDistributor creates a new instance of MerkleDistributor, bound to a specific deployed contract.
func NewMerkleDistributor(address common.Address, backend bind.ContractBackend) (*MerkleDistributor, error) {
	contract, err := bindMerkleDistributor(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MerkleDistributor{MerkleDistributorCaller: MerkleDistributorCaller{contract: contract}, MerkleDistributorTransactor: MerkleDistributorTransactor{contract: contract}, MerkleDistributorFilterer: MerkleDistributorFilterer{contract: contract}}, nil
}

// NewMerkleDistributorCaller creates a new read-only instance of MerkleDistributor, bound to a specific deployed contract.
func NewMerkleDistributorCaller(address common.Address, caller bind.ContractCaller) (*MerkleDistributorCaller, error) {
	contract, err := bindMerkleDistributor(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MerkleDistributorCaller{contract: contract}, nil
}

// NewMerkleDistributorTransactor creates a new write-only instance of MerkleDistributor, bound to a specific deployed contract.
func NewMerkleDistributorTransactor(address common.Address, transactor bind.ContractTransactor) (*MerkleDistributorTransactor, error) {
	contract, err := bindMerkleDistributor(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MerkleDistributorTransactor{contract: contract}, nil
}

// NewMerkleDistributorFilterer creates a new log filterer instance of MerkleDistributor, bound to a specific deployed contract.
func NewMerkleDistributorFilterer(address common.Address, filterer bind.ContractFilterer) (*MerkleDistributorFilterer, error) {
	contract, err := bindMerkleDistributor(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MerkleDistributorFilterer{contract: contract}, nil
}

// bindMerkleDistributor binds a generic wrapper to an already deployed contract.
func bindMerkleDistributor(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := MerkleDistributorMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MerkleDistributor *MerkleDistributorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MerkleDistributor.Contract.MerkleDistributorCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MerkleDistributor *MerkleDistributorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MerkleDistributor.Contract.MerkleDistributorTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MerkleDistributor *MerkleDistributorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MerkleDistributor.Contract.MerkleDistributorTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MerkleDistributor *MerkleDistributorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MerkleDistributor.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MerkleDistributor *MerkleDistributorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MerkleDistributor.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MerkleDistributor *MerkleDistributorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MerkleDistributor.Contract.contract.Transact(opts, method, params...)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_MerkleDistributor *MerkleDistributorCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _MerkleDistributor.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_MerkleDistributor *MerkleDistributorSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _MerkleDistributor.Contract.DEFAULTADMINROLE(&_MerkleDistributor.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_MerkleDistributor *MerkleDistributorCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _MerkleDistributor.Contract.DEFAULTADMINROLE(&_MerkleDistributor.CallOpts)
}

// UPGRADERROLE is a free data retrieval call binding the contract method 0xf72c0d8b.
//
// Solidity: function UPGRADER_ROLE() view returns(bytes32)
func (_MerkleDistributor *MerkleDistributorCaller) UPGRADERROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _MerkleDistributor.contract.Call(opts, &out, "UPGRADER_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// UPGRADERROLE is a free data retrieval call binding the contract method 0xf72c0d8b.
//
// Solidity: function UPGRADER_ROLE() view returns(bytes32)
func (_MerkleDistributor *MerkleDistributorSession) UPGRADERROLE() ([32]byte, error) {
	return _MerkleDistributor.Contract.UPGRADERROLE(&_MerkleDistributor.CallOpts)
}

// UPGRADERROLE is a free data retrieval call binding the contract method 0xf72c0d8b.
//
// Solidity: function UPGRADER_ROLE() view returns(bytes32)
func (_MerkleDistributor *MerkleDistributorCallerSession) UPGRADERROLE() ([32]byte, error) {
	return _MerkleDistributor.Contract.UPGRADERROLE(&_MerkleDistributor.CallOpts)
}

// Allocations is a free data retrieval call binding the contract method 0x4eb71d73.
//
// Solidity: function allocations(uint256 , uint256 ) view returns(uint256)
func (_MerkleDistributor *MerkleDistributorCaller) Allocations(opts *bind.CallOpts, arg0 *big.Int, arg1 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _MerkleDistributor.contract.Call(opts, &out, "allocations", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allocations is a free data retrieval call binding the contract method 0x4eb71d73.
//
// Solidity: function allocations(uint256 , uint256 ) view returns(uint256)
func (_MerkleDistributor *MerkleDistributorSession) Allocations(arg0 *big.Int, arg1 *big.Int) (*big.Int, error) {
	return _MerkleDistributor.Contract.Allocations(&_MerkleDistributor.CallOpts, arg0, arg1)
}

// Allocations is a free data retrieval call binding the contract method 0x4eb71d73.
//
// Solidity: function allocations(uint256 , uint256 ) view returns(uint256)
func (_MerkleDistributor *MerkleDistributorCallerSession) Allocations(arg0 *big.Int, arg1 *big.Int) (*big.Int, error) {
	return _MerkleDistributor.Contract.Allocations(&_MerkleDistributor.CallOpts, arg0, arg1)
}

// Claimed is a free data retrieval call binding the contract method 0x47f5af8d.
//
// Solidity: function claimed(uint256 , uint256 , address ) view returns(bool)
func (_MerkleDistributor *MerkleDistributorCaller) Claimed(opts *bind.CallOpts, arg0 *big.Int, arg1 *big.Int, arg2 common.Address) (bool, error) {
	var out []interface{}
	err := _MerkleDistributor.contract.Call(opts, &out, "claimed", arg0, arg1, arg2)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Claimed is a free data retrieval call binding the contract method 0x47f5af8d.
//
// Solidity: function claimed(uint256 , uint256 , address ) view returns(bool)
func (_MerkleDistributor *MerkleDistributorSession) Claimed(arg0 *big.Int, arg1 *big.Int, arg2 common.Address) (bool, error) {
	return _MerkleDistributor.Contract.Claimed(&_MerkleDistributor.CallOpts, arg0, arg1, arg2)
}

// Claimed is a free data retrieval call binding the contract method 0x47f5af8d.
//
// Solidity: function claimed(uint256 , uint256 , address ) view returns(bool)
func (_MerkleDistributor *MerkleDistributorCallerSession) Claimed(arg0 *big.Int, arg1 *big.Int, arg2 common.Address) (bool, error) {
	return _MerkleDistributor.Contract.Claimed(&_MerkleDistributor.CallOpts, arg0, arg1, arg2)
}

// ClaimedTotals is a free data retrieval call binding the contract method 0x653bec5f.
//
// Solidity: function claimedTotals(uint256 , uint256 ) view returns(uint256)
func (_MerkleDistributor *MerkleDistributorCaller) ClaimedTotals(opts *bind.CallOpts, arg0 *big.Int, arg1 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _MerkleDistributor.contract.Call(opts, &out, "claimedTotals", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ClaimedTotals is a free data retrieval call binding the contract method 0x653bec5f.
//
// Solidity: function claimedTotals(uint256 , uint256 ) view returns(uint256)
func (_MerkleDistributor *MerkleDistributorSession) ClaimedTotals(arg0 *big.Int, arg1 *big.Int) (*big.Int, error) {
	return _MerkleDistributor.Contract.ClaimedTotals(&_MerkleDistributor.CallOpts, arg0, arg1)
}

// ClaimedTotals is a free data retrieval call binding the contract method 0x653bec5f.
//
// Solidity: function claimedTotals(uint256 , uint256 ) view returns(uint256)
func (_MerkleDistributor *MerkleDistributorCallerSession) ClaimedTotals(arg0 *big.Int, arg1 *big.Int) (*big.Int, error) {
	return _MerkleDistributor.Contract.ClaimedTotals(&_MerkleDistributor.CallOpts, arg0, arg1)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_MerkleDistributor *MerkleDistributorCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _MerkleDistributor.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_MerkleDistributor *MerkleDistributorSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _MerkleDistributor.Contract.GetRoleAdmin(&_MerkleDistributor.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_MerkleDistributor *MerkleDistributorCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _MerkleDistributor.Contract.GetRoleAdmin(&_MerkleDistributor.CallOpts, role)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_MerkleDistributor *MerkleDistributorCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _MerkleDistributor.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_MerkleDistributor *MerkleDistributorSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _MerkleDistributor.Contract.HasRole(&_MerkleDistributor.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_MerkleDistributor *MerkleDistributorCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _MerkleDistributor.Contract.HasRole(&_MerkleDistributor.CallOpts, role, account)
}

// IsClaimedBatch is a free data retrieval call binding the contract method 0xd883a185.
//
// Solidity: function isClaimedBatch(uint256 poolId, address account, uint256[] weeks_) view returns(bool[] statuses)
func (_MerkleDistributor *MerkleDistributorCaller) IsClaimedBatch(opts *bind.CallOpts, poolId *big.Int, account common.Address, weeks_ []*big.Int) ([]bool, error) {
	var out []interface{}
	err := _MerkleDistributor.contract.Call(opts, &out, "isClaimedBatch", poolId, account, weeks_)

	if err != nil {
		return *new([]bool), err
	}

	out0 := *abi.ConvertType(out[0], new([]bool)).(*[]bool)

	return out0, err

}

// IsClaimedBatch is a free data retrieval call binding the contract method 0xd883a185.
//
// Solidity: function isClaimedBatch(uint256 poolId, address account, uint256[] weeks_) view returns(bool[] statuses)
func (_MerkleDistributor *MerkleDistributorSession) IsClaimedBatch(poolId *big.Int, account common.Address, weeks_ []*big.Int) ([]bool, error) {
	return _MerkleDistributor.Contract.IsClaimedBatch(&_MerkleDistributor.CallOpts, poolId, account, weeks_)
}

// IsClaimedBatch is a free data retrieval call binding the contract method 0xd883a185.
//
// Solidity: function isClaimedBatch(uint256 poolId, address account, uint256[] weeks_) view returns(bool[] statuses)
func (_MerkleDistributor *MerkleDistributorCallerSession) IsClaimedBatch(poolId *big.Int, account common.Address, weeks_ []*big.Int) ([]bool, error) {
	return _MerkleDistributor.Contract.IsClaimedBatch(&_MerkleDistributor.CallOpts, poolId, account, weeks_)
}

// IsRootSetter is a free data retrieval call binding the contract method 0xf7bcba6a.
//
// Solidity: function isRootSetter(uint256 , address ) view returns(bool)
func (_MerkleDistributor *MerkleDistributorCaller) IsRootSetter(opts *bind.CallOpts, arg0 *big.Int, arg1 common.Address) (bool, error) {
	var out []interface{}
	err := _MerkleDistributor.contract.Call(opts, &out, "isRootSetter", arg0, arg1)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsRootSetter is a free data retrieval call binding the contract method 0xf7bcba6a.
//
// Solidity: function isRootSetter(uint256 , address ) view returns(bool)
func (_MerkleDistributor *MerkleDistributorSession) IsRootSetter(arg0 *big.Int, arg1 common.Address) (bool, error) {
	return _MerkleDistributor.Contract.IsRootSetter(&_MerkleDistributor.CallOpts, arg0, arg1)
}

// IsRootSetter is a free data retrieval call binding the contract method 0xf7bcba6a.
//
// Solidity: function isRootSetter(uint256 , address ) view returns(bool)
func (_MerkleDistributor *MerkleDistributorCallerSession) IsRootSetter(arg0 *big.Int, arg1 common.Address) (bool, error) {
	return _MerkleDistributor.Contract.IsRootSetter(&_MerkleDistributor.CallOpts, arg0, arg1)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_MerkleDistributor *MerkleDistributorCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _MerkleDistributor.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_MerkleDistributor *MerkleDistributorSession) Paused() (bool, error) {
	return _MerkleDistributor.Contract.Paused(&_MerkleDistributor.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_MerkleDistributor *MerkleDistributorCallerSession) Paused() (bool, error) {
	return _MerkleDistributor.Contract.Paused(&_MerkleDistributor.CallOpts)
}

// PoolCount is a free data retrieval call binding the contract method 0xf525cb68.
//
// Solidity: function poolCount() view returns(uint256)
func (_MerkleDistributor *MerkleDistributorCaller) PoolCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _MerkleDistributor.contract.Call(opts, &out, "poolCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PoolCount is a free data retrieval call binding the contract method 0xf525cb68.
//
// Solidity: function poolCount() view returns(uint256)
func (_MerkleDistributor *MerkleDistributorSession) PoolCount() (*big.Int, error) {
	return _MerkleDistributor.Contract.PoolCount(&_MerkleDistributor.CallOpts)
}

// PoolCount is a free data retrieval call binding the contract method 0xf525cb68.
//
// Solidity: function poolCount() view returns(uint256)
func (_MerkleDistributor *MerkleDistributorCallerSession) PoolCount() (*big.Int, error) {
	return _MerkleDistributor.Contract.PoolCount(&_MerkleDistributor.CallOpts)
}

// Pools is a free data retrieval call binding the contract method 0xac4afa38.
//
// Solidity: function pools(uint256 ) view returns(address token, address admin, uint256 weeklyLimit, uint256 balance, uint256 committed)
func (_MerkleDistributor *MerkleDistributorCaller) Pools(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Token       common.Address
	Admin       common.Address
	WeeklyLimit *big.Int
	Balance     *big.Int
	Committed   *big.Int
}, error) {
	var out []interface{}
	err := _MerkleDistributor.contract.Call(opts, &out, "pools", arg0)

	outstruct := new(struct {
		Token       common.Address
		Admin       common.Address
		WeeklyLimit *big.Int
		Balance     *big.Int
		Committed   *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Token = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Admin = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.WeeklyLimit = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.Balance = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.Committed = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// Pools is a free data retrieval call binding the contract method 0xac4afa38.
//
// Solidity: function pools(uint256 ) view returns(address token, address admin, uint256 weeklyLimit, uint256 balance, uint256 committed)
func (_MerkleDistributor *MerkleDistributorSession) Pools(arg0 *big.Int) (struct {
	Token       common.Address
	Admin       common.Address
	WeeklyLimit *big.Int
	Balance     *big.Int
	Committed   *big.Int
}, error) {
	return _MerkleDistributor.Contract.Pools(&_MerkleDistributor.CallOpts, arg0)
}

// Pools is a free data retrieval call binding the contract method 0xac4afa38.
//
// Solidity: function pools(uint256 ) view returns(address token, address admin, uint256 weeklyLimit, uint256 balance, uint256 committed)
func (_MerkleDistributor *MerkleDistributorCallerSession) Pools(arg0 *big.Int) (struct {
	Token       common.Address
	Admin       common.Address
	WeeklyLimit *big.Int
	Balance     *big.Int
	Committed   *big.Int
}, error) {
	return _MerkleDistributor.Contract.Pools(&_MerkleDistributor.CallOpts, arg0)
}

// ProofsURIs is a free data retrieval call binding the contract method 0x67e84e94.
//
// Solidity: function proofsURIs(uint256 , uint256 ) view returns(string)
func (_MerkleDistributor *MerkleDistributorCaller) ProofsURIs(opts *bind.CallOpts, arg0 *big.Int, arg1 *big.Int) (string, error) {
	var out []interface{}
	err := _MerkleDistributor.contract.Call(opts, &out, "proofsURIs", arg0, arg1)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// ProofsURIs is a free data retrieval call binding the contract method 0x67e84e94.
//
// Solidity: function proofsURIs(uint256 , uint256 ) view returns(string)
func (_MerkleDistributor *MerkleDistributorSession) ProofsURIs(arg0 *big.Int, arg1 *big.Int) (string, error) {
	return _MerkleDistributor.Contract.ProofsURIs(&_MerkleDistributor.CallOpts, arg0, arg1)
}

// ProofsURIs is a free data retrieval call binding the contract method 0x67e84e94.
//
// Solidity: function proofsURIs(uint256 , uint256 ) view returns(string)
func (_MerkleDistributor *MerkleDistributorCallerSession) ProofsURIs(arg0 *big.Int, arg1 *big.Int) (string, error) {
	return _MerkleDistributor.Contract.ProofsURIs(&_MerkleDistributor.CallOpts, arg0, arg1)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_MerkleDistributor *MerkleDistributorCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _MerkleDistributor.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_MerkleDistributor *MerkleDistributorSession) ProxiableUUID() ([32]byte, error) {
	return _MerkleDistributor.Contract.ProxiableUUID(&_MerkleDistributor.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_MerkleDistributor *MerkleDistributorCallerSession) ProxiableUUID() ([32]byte, error) {
	return _MerkleDistributor.Contract.ProxiableUUID(&_MerkleDistributor.CallOpts)
}

// Roots is a free data retrieval call binding the contract method 0x0fc3cec0.
//
// Solidity: function roots(uint256 , uint256 ) view returns(bytes32)
func (_MerkleDistributor *MerkleDistributorCaller) Roots(opts *bind.CallOpts, arg0 *big.Int, arg1 *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _MerkleDistributor.contract.Call(opts, &out, "roots", arg0, arg1)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// Roots is a free data retrieval call binding the contract method 0x0fc3cec0.
//
// Solidity: function roots(uint256 , uint256 ) view returns(bytes32)
func (_MerkleDistributor *MerkleDistributorSession) Roots(arg0 *big.Int, arg1 *big.Int) ([32]byte, error) {
	return _MerkleDistributor.Contract.Roots(&_MerkleDistributor.CallOpts, arg0, arg1)
}

// Roots is a free data retrieval call binding the contract method 0x0fc3cec0.
//
// Solidity: function roots(uint256 , uint256 ) view returns(bytes32)
func (_MerkleDistributor *MerkleDistributorCallerSession) Roots(arg0 *big.Int, arg1 *big.Int) ([32]byte, error) {
	return _MerkleDistributor.Contract.Roots(&_MerkleDistributor.CallOpts, arg0, arg1)
}

// SanctionsList is a free data retrieval call binding the contract method 0xec571c6a.
//
// Solidity: function sanctionsList() view returns(address)
func (_MerkleDistributor *MerkleDistributorCaller) SanctionsList(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _MerkleDistributor.contract.Call(opts, &out, "sanctionsList")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// SanctionsList is a free data retrieval call binding the contract method 0xec571c6a.
//
// Solidity: function sanctionsList() view returns(address)
func (_MerkleDistributor *MerkleDistributorSession) SanctionsList() (common.Address, error) {
	return _MerkleDistributor.Contract.SanctionsList(&_MerkleDistributor.CallOpts)
}

// SanctionsList is a free data retrieval call binding the contract method 0xec571c6a.
//
// Solidity: function sanctionsList() view returns(address)
func (_MerkleDistributor *MerkleDistributorCallerSession) SanctionsList() (common.Address, error) {
	return _MerkleDistributor.Contract.SanctionsList(&_MerkleDistributor.CallOpts)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_MerkleDistributor *MerkleDistributorCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _MerkleDistributor.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_MerkleDistributor *MerkleDistributorSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _MerkleDistributor.Contract.SupportsInterface(&_MerkleDistributor.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_MerkleDistributor *MerkleDistributorCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _MerkleDistributor.Contract.SupportsInterface(&_MerkleDistributor.CallOpts, interfaceId)
}

// Claim is a paid mutator transaction binding the contract method 0x5d4df3bf.
//
// Solidity: function claim(uint256 poolId, uint256 week, address account, uint256 amount, bytes32[] proof) returns()
func (_MerkleDistributor *MerkleDistributorTransactor) Claim(opts *bind.TransactOpts, poolId *big.Int, week *big.Int, account common.Address, amount *big.Int, proof [][32]byte) (*types.Transaction, error) {
	return _MerkleDistributor.contract.Transact(opts, "claim", poolId, week, account, amount, proof)
}

// Claim is a paid mutator transaction binding the contract method 0x5d4df3bf.
//
// Solidity: function claim(uint256 poolId, uint256 week, address account, uint256 amount, bytes32[] proof) returns()
func (_MerkleDistributor *MerkleDistributorSession) Claim(poolId *big.Int, week *big.Int, account common.Address, amount *big.Int, proof [][32]byte) (*types.Transaction, error) {
	return _MerkleDistributor.Contract.Claim(&_MerkleDistributor.TransactOpts, poolId, week, account, amount, proof)
}

// Claim is a paid mutator transaction binding the contract method 0x5d4df3bf.
//
// Solidity: function claim(uint256 poolId, uint256 week, address account, uint256 amount, bytes32[] proof) returns()
func (_MerkleDistributor *MerkleDistributorTransactorSession) Claim(poolId *big.Int, week *big.Int, account common.Address, amount *big.Int, proof [][32]byte) (*types.Transaction, error) {
	return _MerkleDistributor.Contract.Claim(&_MerkleDistributor.TransactOpts, poolId, week, account, amount, proof)
}

// ClaimBatch is a paid mutator transaction binding the contract method 0xe5be943b.
//
// Solidity: function claimBatch(uint256 poolId, uint256[] weeks_, address account, uint256[] amounts, bytes32[][] proofs) returns()
func (_MerkleDistributor *MerkleDistributorTransactor) ClaimBatch(opts *bind.TransactOpts, poolId *big.Int, weeks_ []*big.Int, account common.Address, amounts []*big.Int, proofs [][][32]byte) (*types.Transaction, error) {
	return _MerkleDistributor.contract.Transact(opts, "claimBatch", poolId, weeks_, account, amounts, proofs)
}

// ClaimBatch is a paid mutator transaction binding the contract method 0xe5be943b.
//
// Solidity: function claimBatch(uint256 poolId, uint256[] weeks_, address account, uint256[] amounts, bytes32[][] proofs) returns()
func (_MerkleDistributor *MerkleDistributorSession) ClaimBatch(poolId *big.Int, weeks_ []*big.Int, account common.Address, amounts []*big.Int, proofs [][][32]byte) (*types.Transaction, error) {
	return _MerkleDistributor.Contract.ClaimBatch(&_MerkleDistributor.TransactOpts, poolId, weeks_, account, amounts, proofs)
}

// ClaimBatch is a paid mutator transaction binding the contract method 0xe5be943b.
//
// Solidity: function claimBatch(uint256 poolId, uint256[] weeks_, address account, uint256[] amounts, bytes32[][] proofs) returns()
func (_MerkleDistributor *MerkleDistributorTransactorSession) ClaimBatch(poolId *big.Int, weeks_ []*big.Int, account common.Address, amounts []*big.Int, proofs [][][32]byte) (*types.Transaction, error) {
	return _MerkleDistributor.Contract.ClaimBatch(&_MerkleDistributor.TransactOpts, poolId, weeks_, account, amounts, proofs)
}

// CreatePool is a paid mutator transaction binding the contract method 0x51810fb5.
//
// Solidity: function createPool(address token, address admin, uint256 weeklyLimit) returns(uint256 poolId)
func (_MerkleDistributor *MerkleDistributorTransactor) CreatePool(opts *bind.TransactOpts, token common.Address, admin common.Address, weeklyLimit *big.Int) (*types.Transaction, error) {
	return _MerkleDistributor.contract.Transact(opts, "createPool", token, admin, weeklyLimit)
}

// CreatePool is a paid mutator transaction binding the contract method 0x51810fb5.
//
// Solidity: function createPool(address token, address admin, uint256 weeklyLimit) returns(uint256 poolId)
func (_MerkleDistributor *MerkleDistributorSession) CreatePool(token common.Address, admin common.Address, weeklyLimit *big.Int) (*types.Transaction, error) {
	return _MerkleDistributor.Contract.CreatePool(&_MerkleDistributor.TransactOpts, token, admin, weeklyLimit)
}

// CreatePool is a paid mutator transaction binding the contract method 0x51810fb5.
//
// Solidity: function createPool(address token, address admin, uint256 weeklyLimit) returns(uint256 poolId)
func (_MerkleDistributor *MerkleDistributorTransactorSession) CreatePool(token common.Address, admin common.Address, weeklyLimit *big.Int) (*types.Transaction, error) {
	return _MerkleDistributor.Contract.CreatePool(&_MerkleDistributor.TransactOpts, token, admin, weeklyLimit)
}

// Fund is a paid mutator transaction binding the contract method 0xa65e2cfd.
//
// Solidity: function fund(uint256 poolId, uint256 amount) returns()
func (_MerkleDistributor *MerkleDistributorTransactor) Fund(opts *bind.TransactOpts, poolId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _MerkleDistributor.contract.Transact(opts, "fund", poolId, amount)
}

// Fund is a paid mutator transaction binding the contract method 0xa65e2cfd.
//
// Solidity: function fund(uint256 poolId, uint256 amount) returns()
func (_MerkleDistributor *MerkleDistributorSession) Fund(poolId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _MerkleDistributor.Contract.Fund(&_MerkleDistributor.TransactOpts, poolId, amount)
}

// Fund is a paid mutator transaction binding the contract method 0xa65e2cfd.
//
// Solidity: function fund(uint256 poolId, uint256 amount) returns()
func (_MerkleDistributor *MerkleDistributorTransactorSession) Fund(poolId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _MerkleDistributor.Contract.Fund(&_MerkleDistributor.TransactOpts, poolId, amount)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_MerkleDistributor *MerkleDistributorTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _MerkleDistributor.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_MerkleDistributor *MerkleDistributorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _MerkleDistributor.Contract.GrantRole(&_MerkleDistributor.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_MerkleDistributor *MerkleDistributorTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _MerkleDistributor.Contract.GrantRole(&_MerkleDistributor.TransactOpts, role, account)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address sanctionsList_) returns()
func (_MerkleDistributor *MerkleDistributorTransactor) Initialize(opts *bind.TransactOpts, sanctionsList_ common.Address) (*types.Transaction, error) {
	return _MerkleDistributor.contract.Transact(opts, "initialize", sanctionsList_)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address sanctionsList_) returns()
func (_MerkleDistributor *MerkleDistributorSession) Initialize(sanctionsList_ common.Address) (*types.Transaction, error) {
	return _MerkleDistributor.Contract.Initialize(&_MerkleDistributor.TransactOpts, sanctionsList_)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address sanctionsList_) returns()
func (_MerkleDistributor *MerkleDistributorTransactorSession) Initialize(sanctionsList_ common.Address) (*types.Transaction, error) {
	return _MerkleDistributor.Contract.Initialize(&_MerkleDistributor.TransactOpts, sanctionsList_)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_MerkleDistributor *MerkleDistributorTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MerkleDistributor.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_MerkleDistributor *MerkleDistributorSession) Pause() (*types.Transaction, error) {
	return _MerkleDistributor.Contract.Pause(&_MerkleDistributor.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_MerkleDistributor *MerkleDistributorTransactorSession) Pause() (*types.Transaction, error) {
	return _MerkleDistributor.Contract.Pause(&_MerkleDistributor.TransactOpts)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_MerkleDistributor *MerkleDistributorTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _MerkleDistributor.contract.Transact(opts, "renounceRole", role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_MerkleDistributor *MerkleDistributorSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _MerkleDistributor.Contract.RenounceRole(&_MerkleDistributor.TransactOpts, role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_MerkleDistributor *MerkleDistributorTransactorSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _MerkleDistributor.Contract.RenounceRole(&_MerkleDistributor.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_MerkleDistributor *MerkleDistributorTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _MerkleDistributor.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_MerkleDistributor *MerkleDistributorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _MerkleDistributor.Contract.RevokeRole(&_MerkleDistributor.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_MerkleDistributor *MerkleDistributorTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _MerkleDistributor.Contract.RevokeRole(&_MerkleDistributor.TransactOpts, role, account)
}

// SetRoot is a paid mutator transaction binding the contract method 0x2e3f4fc2.
//
// Solidity: function setRoot(uint256 poolId, uint256 week, bytes32 root, uint256 allocation, string proofsURI) returns()
func (_MerkleDistributor *MerkleDistributorTransactor) SetRoot(opts *bind.TransactOpts, poolId *big.Int, week *big.Int, root [32]byte, allocation *big.Int, proofsURI string) (*types.Transaction, error) {
	return _MerkleDistributor.contract.Transact(opts, "setRoot", poolId, week, root, allocation, proofsURI)
}

// SetRoot is a paid mutator transaction binding the contract method 0x2e3f4fc2.
//
// Solidity: function setRoot(uint256 poolId, uint256 week, bytes32 root, uint256 allocation, string proofsURI) returns()
func (_MerkleDistributor *MerkleDistributorSession) SetRoot(poolId *big.Int, week *big.Int, root [32]byte, allocation *big.Int, proofsURI string) (*types.Transaction, error) {
	return _MerkleDistributor.Contract.SetRoot(&_MerkleDistributor.TransactOpts, poolId, week, root, allocation, proofsURI)
}

// SetRoot is a paid mutator transaction binding the contract method 0x2e3f4fc2.
//
// Solidity: function setRoot(uint256 poolId, uint256 week, bytes32 root, uint256 allocation, string proofsURI) returns()
func (_MerkleDistributor *MerkleDistributorTransactorSession) SetRoot(poolId *big.Int, week *big.Int, root [32]byte, allocation *big.Int, proofsURI string) (*types.Transaction, error) {
	return _MerkleDistributor.Contract.SetRoot(&_MerkleDistributor.TransactOpts, poolId, week, root, allocation, proofsURI)
}

// SetRootSetter is a paid mutator transaction binding the contract method 0x10fd63cb.
//
// Solidity: function setRootSetter(uint256 poolId, address setter, bool enabled) returns()
func (_MerkleDistributor *MerkleDistributorTransactor) SetRootSetter(opts *bind.TransactOpts, poolId *big.Int, setter common.Address, enabled bool) (*types.Transaction, error) {
	return _MerkleDistributor.contract.Transact(opts, "setRootSetter", poolId, setter, enabled)
}

// SetRootSetter is a paid mutator transaction binding the contract method 0x10fd63cb.
//
// Solidity: function setRootSetter(uint256 poolId, address setter, bool enabled) returns()
func (_MerkleDistributor *MerkleDistributorSession) SetRootSetter(poolId *big.Int, setter common.Address, enabled bool) (*types.Transaction, error) {
	return _MerkleDistributor.Contract.SetRootSetter(&_MerkleDistributor.TransactOpts, poolId, setter, enabled)
}

// SetRootSetter is a paid mutator transaction binding the contract method 0x10fd63cb.
//
// Solidity: function setRootSetter(uint256 poolId, address setter, bool enabled) returns()
func (_MerkleDistributor *MerkleDistributorTransactorSession) SetRootSetter(poolId *big.Int, setter common.Address, enabled bool) (*types.Transaction, error) {
	return _MerkleDistributor.Contract.SetRootSetter(&_MerkleDistributor.TransactOpts, poolId, setter, enabled)
}

// SetWeeklyLimit is a paid mutator transaction binding the contract method 0xe8080159.
//
// Solidity: function setWeeklyLimit(uint256 poolId, uint256 limit) returns()
func (_MerkleDistributor *MerkleDistributorTransactor) SetWeeklyLimit(opts *bind.TransactOpts, poolId *big.Int, limit *big.Int) (*types.Transaction, error) {
	return _MerkleDistributor.contract.Transact(opts, "setWeeklyLimit", poolId, limit)
}

// SetWeeklyLimit is a paid mutator transaction binding the contract method 0xe8080159.
//
// Solidity: function setWeeklyLimit(uint256 poolId, uint256 limit) returns()
func (_MerkleDistributor *MerkleDistributorSession) SetWeeklyLimit(poolId *big.Int, limit *big.Int) (*types.Transaction, error) {
	return _MerkleDistributor.Contract.SetWeeklyLimit(&_MerkleDistributor.TransactOpts, poolId, limit)
}

// SetWeeklyLimit is a paid mutator transaction binding the contract method 0xe8080159.
//
// Solidity: function setWeeklyLimit(uint256 poolId, uint256 limit) returns()
func (_MerkleDistributor *MerkleDistributorTransactorSession) SetWeeklyLimit(poolId *big.Int, limit *big.Int) (*types.Transaction, error) {
	return _MerkleDistributor.Contract.SetWeeklyLimit(&_MerkleDistributor.TransactOpts, poolId, limit)
}

// Sweep is a paid mutator transaction binding the contract method 0xba01ed79.
//
// Solidity: function sweep(uint256 poolId, address to, uint256 amount) returns()
func (_MerkleDistributor *MerkleDistributorTransactor) Sweep(opts *bind.TransactOpts, poolId *big.Int, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _MerkleDistributor.contract.Transact(opts, "sweep", poolId, to, amount)
}

// Sweep is a paid mutator transaction binding the contract method 0xba01ed79.
//
// Solidity: function sweep(uint256 poolId, address to, uint256 amount) returns()
func (_MerkleDistributor *MerkleDistributorSession) Sweep(poolId *big.Int, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _MerkleDistributor.Contract.Sweep(&_MerkleDistributor.TransactOpts, poolId, to, amount)
}

// Sweep is a paid mutator transaction binding the contract method 0xba01ed79.
//
// Solidity: function sweep(uint256 poolId, address to, uint256 amount) returns()
func (_MerkleDistributor *MerkleDistributorTransactorSession) Sweep(poolId *big.Int, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _MerkleDistributor.Contract.Sweep(&_MerkleDistributor.TransactOpts, poolId, to, amount)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_MerkleDistributor *MerkleDistributorTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MerkleDistributor.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_MerkleDistributor *MerkleDistributorSession) Unpause() (*types.Transaction, error) {
	return _MerkleDistributor.Contract.Unpause(&_MerkleDistributor.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_MerkleDistributor *MerkleDistributorTransactorSession) Unpause() (*types.Transaction, error) {
	return _MerkleDistributor.Contract.Unpause(&_MerkleDistributor.TransactOpts)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_MerkleDistributor *MerkleDistributorTransactor) UpgradeTo(opts *bind.TransactOpts, newImplementation common.Address) (*types.Transaction, error) {
	return _MerkleDistributor.contract.Transact(opts, "upgradeTo", newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_MerkleDistributor *MerkleDistributorSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _MerkleDistributor.Contract.UpgradeTo(&_MerkleDistributor.TransactOpts, newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_MerkleDistributor *MerkleDistributorTransactorSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _MerkleDistributor.Contract.UpgradeTo(&_MerkleDistributor.TransactOpts, newImplementation)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_MerkleDistributor *MerkleDistributorTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _MerkleDistributor.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_MerkleDistributor *MerkleDistributorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _MerkleDistributor.Contract.UpgradeToAndCall(&_MerkleDistributor.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_MerkleDistributor *MerkleDistributorTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _MerkleDistributor.Contract.UpgradeToAndCall(&_MerkleDistributor.TransactOpts, newImplementation, data)
}

// MerkleDistributorAdminChangedIterator is returned from FilterAdminChanged and is used to iterate over the raw logs and unpacked data for AdminChanged events raised by the MerkleDistributor contract.
type MerkleDistributorAdminChangedIterator struct {
	Event *MerkleDistributorAdminChanged // Event containing the contract specifics and raw log

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
func (it *MerkleDistributorAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MerkleDistributorAdminChanged)
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
		it.Event = new(MerkleDistributorAdminChanged)
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
func (it *MerkleDistributorAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MerkleDistributorAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MerkleDistributorAdminChanged represents a AdminChanged event raised by the MerkleDistributor contract.
type MerkleDistributorAdminChanged struct {
	PreviousAdmin common.Address
	NewAdmin      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterAdminChanged is a free log retrieval operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_MerkleDistributor *MerkleDistributorFilterer) FilterAdminChanged(opts *bind.FilterOpts) (*MerkleDistributorAdminChangedIterator, error) {

	logs, sub, err := _MerkleDistributor.contract.FilterLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return &MerkleDistributorAdminChangedIterator{contract: _MerkleDistributor.contract, event: "AdminChanged", logs: logs, sub: sub}, nil
}

// WatchAdminChanged is a free log subscription operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_MerkleDistributor *MerkleDistributorFilterer) WatchAdminChanged(opts *bind.WatchOpts, sink chan<- *MerkleDistributorAdminChanged) (event.Subscription, error) {

	logs, sub, err := _MerkleDistributor.contract.WatchLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MerkleDistributorAdminChanged)
				if err := _MerkleDistributor.contract.UnpackLog(event, "AdminChanged", log); err != nil {
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
func (_MerkleDistributor *MerkleDistributorFilterer) ParseAdminChanged(log types.Log) (*MerkleDistributorAdminChanged, error) {
	event := new(MerkleDistributorAdminChanged)
	if err := _MerkleDistributor.contract.UnpackLog(event, "AdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MerkleDistributorBeaconUpgradedIterator is returned from FilterBeaconUpgraded and is used to iterate over the raw logs and unpacked data for BeaconUpgraded events raised by the MerkleDistributor contract.
type MerkleDistributorBeaconUpgradedIterator struct {
	Event *MerkleDistributorBeaconUpgraded // Event containing the contract specifics and raw log

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
func (it *MerkleDistributorBeaconUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MerkleDistributorBeaconUpgraded)
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
		it.Event = new(MerkleDistributorBeaconUpgraded)
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
func (it *MerkleDistributorBeaconUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MerkleDistributorBeaconUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MerkleDistributorBeaconUpgraded represents a BeaconUpgraded event raised by the MerkleDistributor contract.
type MerkleDistributorBeaconUpgraded struct {
	Beacon common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBeaconUpgraded is a free log retrieval operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_MerkleDistributor *MerkleDistributorFilterer) FilterBeaconUpgraded(opts *bind.FilterOpts, beacon []common.Address) (*MerkleDistributorBeaconUpgradedIterator, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _MerkleDistributor.contract.FilterLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return &MerkleDistributorBeaconUpgradedIterator{contract: _MerkleDistributor.contract, event: "BeaconUpgraded", logs: logs, sub: sub}, nil
}

// WatchBeaconUpgraded is a free log subscription operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_MerkleDistributor *MerkleDistributorFilterer) WatchBeaconUpgraded(opts *bind.WatchOpts, sink chan<- *MerkleDistributorBeaconUpgraded, beacon []common.Address) (event.Subscription, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _MerkleDistributor.contract.WatchLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MerkleDistributorBeaconUpgraded)
				if err := _MerkleDistributor.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
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
func (_MerkleDistributor *MerkleDistributorFilterer) ParseBeaconUpgraded(log types.Log) (*MerkleDistributorBeaconUpgraded, error) {
	event := new(MerkleDistributorBeaconUpgraded)
	if err := _MerkleDistributor.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MerkleDistributorClaimedIterator is returned from FilterClaimed and is used to iterate over the raw logs and unpacked data for Claimed events raised by the MerkleDistributor contract.
type MerkleDistributorClaimedIterator struct {
	Event *MerkleDistributorClaimed // Event containing the contract specifics and raw log

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
func (it *MerkleDistributorClaimedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MerkleDistributorClaimed)
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
		it.Event = new(MerkleDistributorClaimed)
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
func (it *MerkleDistributorClaimedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MerkleDistributorClaimedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MerkleDistributorClaimed represents a Claimed event raised by the MerkleDistributor contract.
type MerkleDistributorClaimed struct {
	PoolId  *big.Int
	Week    *big.Int
	Account common.Address
	Amount  *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterClaimed is a free log retrieval operation binding the contract event 0xb94bf7f9302edf52a596286915a69b4b0685574cffdedd0712e3c62f2550f0ba.
//
// Solidity: event Claimed(uint256 indexed poolId, uint256 indexed week, address indexed account, uint256 amount)
func (_MerkleDistributor *MerkleDistributorFilterer) FilterClaimed(opts *bind.FilterOpts, poolId []*big.Int, week []*big.Int, account []common.Address) (*MerkleDistributorClaimedIterator, error) {

	var poolIdRule []interface{}
	for _, poolIdItem := range poolId {
		poolIdRule = append(poolIdRule, poolIdItem)
	}
	var weekRule []interface{}
	for _, weekItem := range week {
		weekRule = append(weekRule, weekItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _MerkleDistributor.contract.FilterLogs(opts, "Claimed", poolIdRule, weekRule, accountRule)
	if err != nil {
		return nil, err
	}
	return &MerkleDistributorClaimedIterator{contract: _MerkleDistributor.contract, event: "Claimed", logs: logs, sub: sub}, nil
}

// WatchClaimed is a free log subscription operation binding the contract event 0xb94bf7f9302edf52a596286915a69b4b0685574cffdedd0712e3c62f2550f0ba.
//
// Solidity: event Claimed(uint256 indexed poolId, uint256 indexed week, address indexed account, uint256 amount)
func (_MerkleDistributor *MerkleDistributorFilterer) WatchClaimed(opts *bind.WatchOpts, sink chan<- *MerkleDistributorClaimed, poolId []*big.Int, week []*big.Int, account []common.Address) (event.Subscription, error) {

	var poolIdRule []interface{}
	for _, poolIdItem := range poolId {
		poolIdRule = append(poolIdRule, poolIdItem)
	}
	var weekRule []interface{}
	for _, weekItem := range week {
		weekRule = append(weekRule, weekItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _MerkleDistributor.contract.WatchLogs(opts, "Claimed", poolIdRule, weekRule, accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MerkleDistributorClaimed)
				if err := _MerkleDistributor.contract.UnpackLog(event, "Claimed", log); err != nil {
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

// ParseClaimed is a log parse operation binding the contract event 0xb94bf7f9302edf52a596286915a69b4b0685574cffdedd0712e3c62f2550f0ba.
//
// Solidity: event Claimed(uint256 indexed poolId, uint256 indexed week, address indexed account, uint256 amount)
func (_MerkleDistributor *MerkleDistributorFilterer) ParseClaimed(log types.Log) (*MerkleDistributorClaimed, error) {
	event := new(MerkleDistributorClaimed)
	if err := _MerkleDistributor.contract.UnpackLog(event, "Claimed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MerkleDistributorFundedIterator is returned from FilterFunded and is used to iterate over the raw logs and unpacked data for Funded events raised by the MerkleDistributor contract.
type MerkleDistributorFundedIterator struct {
	Event *MerkleDistributorFunded // Event containing the contract specifics and raw log

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
func (it *MerkleDistributorFundedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MerkleDistributorFunded)
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
		it.Event = new(MerkleDistributorFunded)
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
func (it *MerkleDistributorFundedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MerkleDistributorFundedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MerkleDistributorFunded represents a Funded event raised by the MerkleDistributor contract.
type MerkleDistributorFunded struct {
	PoolId *big.Int
	From   common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterFunded is a free log retrieval operation binding the contract event 0x38c48552690c96ec2872092ac1db6c19fb59f5a8c5b49bbf41ed4886d0ca6926.
//
// Solidity: event Funded(uint256 indexed poolId, address indexed from, uint256 amount)
func (_MerkleDistributor *MerkleDistributorFilterer) FilterFunded(opts *bind.FilterOpts, poolId []*big.Int, from []common.Address) (*MerkleDistributorFundedIterator, error) {

	var poolIdRule []interface{}
	for _, poolIdItem := range poolId {
		poolIdRule = append(poolIdRule, poolIdItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _MerkleDistributor.contract.FilterLogs(opts, "Funded", poolIdRule, fromRule)
	if err != nil {
		return nil, err
	}
	return &MerkleDistributorFundedIterator{contract: _MerkleDistributor.contract, event: "Funded", logs: logs, sub: sub}, nil
}

// WatchFunded is a free log subscription operation binding the contract event 0x38c48552690c96ec2872092ac1db6c19fb59f5a8c5b49bbf41ed4886d0ca6926.
//
// Solidity: event Funded(uint256 indexed poolId, address indexed from, uint256 amount)
func (_MerkleDistributor *MerkleDistributorFilterer) WatchFunded(opts *bind.WatchOpts, sink chan<- *MerkleDistributorFunded, poolId []*big.Int, from []common.Address) (event.Subscription, error) {

	var poolIdRule []interface{}
	for _, poolIdItem := range poolId {
		poolIdRule = append(poolIdRule, poolIdItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _MerkleDistributor.contract.WatchLogs(opts, "Funded", poolIdRule, fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MerkleDistributorFunded)
				if err := _MerkleDistributor.contract.UnpackLog(event, "Funded", log); err != nil {
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

// ParseFunded is a log parse operation binding the contract event 0x38c48552690c96ec2872092ac1db6c19fb59f5a8c5b49bbf41ed4886d0ca6926.
//
// Solidity: event Funded(uint256 indexed poolId, address indexed from, uint256 amount)
func (_MerkleDistributor *MerkleDistributorFilterer) ParseFunded(log types.Log) (*MerkleDistributorFunded, error) {
	event := new(MerkleDistributorFunded)
	if err := _MerkleDistributor.contract.UnpackLog(event, "Funded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MerkleDistributorInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the MerkleDistributor contract.
type MerkleDistributorInitializedIterator struct {
	Event *MerkleDistributorInitialized // Event containing the contract specifics and raw log

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
func (it *MerkleDistributorInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MerkleDistributorInitialized)
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
		it.Event = new(MerkleDistributorInitialized)
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
func (it *MerkleDistributorInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MerkleDistributorInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MerkleDistributorInitialized represents a Initialized event raised by the MerkleDistributor contract.
type MerkleDistributorInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_MerkleDistributor *MerkleDistributorFilterer) FilterInitialized(opts *bind.FilterOpts) (*MerkleDistributorInitializedIterator, error) {

	logs, sub, err := _MerkleDistributor.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &MerkleDistributorInitializedIterator{contract: _MerkleDistributor.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_MerkleDistributor *MerkleDistributorFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *MerkleDistributorInitialized) (event.Subscription, error) {

	logs, sub, err := _MerkleDistributor.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MerkleDistributorInitialized)
				if err := _MerkleDistributor.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_MerkleDistributor *MerkleDistributorFilterer) ParseInitialized(log types.Log) (*MerkleDistributorInitialized, error) {
	event := new(MerkleDistributorInitialized)
	if err := _MerkleDistributor.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MerkleDistributorPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the MerkleDistributor contract.
type MerkleDistributorPausedIterator struct {
	Event *MerkleDistributorPaused // Event containing the contract specifics and raw log

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
func (it *MerkleDistributorPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MerkleDistributorPaused)
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
		it.Event = new(MerkleDistributorPaused)
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
func (it *MerkleDistributorPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MerkleDistributorPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MerkleDistributorPaused represents a Paused event raised by the MerkleDistributor contract.
type MerkleDistributorPaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_MerkleDistributor *MerkleDistributorFilterer) FilterPaused(opts *bind.FilterOpts) (*MerkleDistributorPausedIterator, error) {

	logs, sub, err := _MerkleDistributor.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &MerkleDistributorPausedIterator{contract: _MerkleDistributor.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_MerkleDistributor *MerkleDistributorFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *MerkleDistributorPaused) (event.Subscription, error) {

	logs, sub, err := _MerkleDistributor.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MerkleDistributorPaused)
				if err := _MerkleDistributor.contract.UnpackLog(event, "Paused", log); err != nil {
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

// ParsePaused is a log parse operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_MerkleDistributor *MerkleDistributorFilterer) ParsePaused(log types.Log) (*MerkleDistributorPaused, error) {
	event := new(MerkleDistributorPaused)
	if err := _MerkleDistributor.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MerkleDistributorPoolCreatedIterator is returned from FilterPoolCreated and is used to iterate over the raw logs and unpacked data for PoolCreated events raised by the MerkleDistributor contract.
type MerkleDistributorPoolCreatedIterator struct {
	Event *MerkleDistributorPoolCreated // Event containing the contract specifics and raw log

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
func (it *MerkleDistributorPoolCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MerkleDistributorPoolCreated)
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
		it.Event = new(MerkleDistributorPoolCreated)
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
func (it *MerkleDistributorPoolCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MerkleDistributorPoolCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MerkleDistributorPoolCreated represents a PoolCreated event raised by the MerkleDistributor contract.
type MerkleDistributorPoolCreated struct {
	PoolId      *big.Int
	Token       common.Address
	Admin       common.Address
	WeeklyLimit *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterPoolCreated is a free log retrieval operation binding the contract event 0xf756876565d755443228c555dcee8573d6405d9be751ac121eb5e386ed86913a.
//
// Solidity: event PoolCreated(uint256 indexed poolId, address indexed token, address indexed admin, uint256 weeklyLimit)
func (_MerkleDistributor *MerkleDistributorFilterer) FilterPoolCreated(opts *bind.FilterOpts, poolId []*big.Int, token []common.Address, admin []common.Address) (*MerkleDistributorPoolCreatedIterator, error) {

	var poolIdRule []interface{}
	for _, poolIdItem := range poolId {
		poolIdRule = append(poolIdRule, poolIdItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var adminRule []interface{}
	for _, adminItem := range admin {
		adminRule = append(adminRule, adminItem)
	}

	logs, sub, err := _MerkleDistributor.contract.FilterLogs(opts, "PoolCreated", poolIdRule, tokenRule, adminRule)
	if err != nil {
		return nil, err
	}
	return &MerkleDistributorPoolCreatedIterator{contract: _MerkleDistributor.contract, event: "PoolCreated", logs: logs, sub: sub}, nil
}

// WatchPoolCreated is a free log subscription operation binding the contract event 0xf756876565d755443228c555dcee8573d6405d9be751ac121eb5e386ed86913a.
//
// Solidity: event PoolCreated(uint256 indexed poolId, address indexed token, address indexed admin, uint256 weeklyLimit)
func (_MerkleDistributor *MerkleDistributorFilterer) WatchPoolCreated(opts *bind.WatchOpts, sink chan<- *MerkleDistributorPoolCreated, poolId []*big.Int, token []common.Address, admin []common.Address) (event.Subscription, error) {

	var poolIdRule []interface{}
	for _, poolIdItem := range poolId {
		poolIdRule = append(poolIdRule, poolIdItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var adminRule []interface{}
	for _, adminItem := range admin {
		adminRule = append(adminRule, adminItem)
	}

	logs, sub, err := _MerkleDistributor.contract.WatchLogs(opts, "PoolCreated", poolIdRule, tokenRule, adminRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MerkleDistributorPoolCreated)
				if err := _MerkleDistributor.contract.UnpackLog(event, "PoolCreated", log); err != nil {
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

// ParsePoolCreated is a log parse operation binding the contract event 0xf756876565d755443228c555dcee8573d6405d9be751ac121eb5e386ed86913a.
//
// Solidity: event PoolCreated(uint256 indexed poolId, address indexed token, address indexed admin, uint256 weeklyLimit)
func (_MerkleDistributor *MerkleDistributorFilterer) ParsePoolCreated(log types.Log) (*MerkleDistributorPoolCreated, error) {
	event := new(MerkleDistributorPoolCreated)
	if err := _MerkleDistributor.contract.UnpackLog(event, "PoolCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MerkleDistributorRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the MerkleDistributor contract.
type MerkleDistributorRoleAdminChangedIterator struct {
	Event *MerkleDistributorRoleAdminChanged // Event containing the contract specifics and raw log

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
func (it *MerkleDistributorRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MerkleDistributorRoleAdminChanged)
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
		it.Event = new(MerkleDistributorRoleAdminChanged)
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
func (it *MerkleDistributorRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MerkleDistributorRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MerkleDistributorRoleAdminChanged represents a RoleAdminChanged event raised by the MerkleDistributor contract.
type MerkleDistributorRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_MerkleDistributor *MerkleDistributorFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*MerkleDistributorRoleAdminChangedIterator, error) {

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

	logs, sub, err := _MerkleDistributor.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &MerkleDistributorRoleAdminChangedIterator{contract: _MerkleDistributor.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_MerkleDistributor *MerkleDistributorFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *MerkleDistributorRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

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

	logs, sub, err := _MerkleDistributor.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MerkleDistributorRoleAdminChanged)
				if err := _MerkleDistributor.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
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
func (_MerkleDistributor *MerkleDistributorFilterer) ParseRoleAdminChanged(log types.Log) (*MerkleDistributorRoleAdminChanged, error) {
	event := new(MerkleDistributorRoleAdminChanged)
	if err := _MerkleDistributor.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MerkleDistributorRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the MerkleDistributor contract.
type MerkleDistributorRoleGrantedIterator struct {
	Event *MerkleDistributorRoleGranted // Event containing the contract specifics and raw log

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
func (it *MerkleDistributorRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MerkleDistributorRoleGranted)
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
		it.Event = new(MerkleDistributorRoleGranted)
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
func (it *MerkleDistributorRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MerkleDistributorRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MerkleDistributorRoleGranted represents a RoleGranted event raised by the MerkleDistributor contract.
type MerkleDistributorRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_MerkleDistributor *MerkleDistributorFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*MerkleDistributorRoleGrantedIterator, error) {

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

	logs, sub, err := _MerkleDistributor.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &MerkleDistributorRoleGrantedIterator{contract: _MerkleDistributor.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_MerkleDistributor *MerkleDistributorFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *MerkleDistributorRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _MerkleDistributor.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MerkleDistributorRoleGranted)
				if err := _MerkleDistributor.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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
func (_MerkleDistributor *MerkleDistributorFilterer) ParseRoleGranted(log types.Log) (*MerkleDistributorRoleGranted, error) {
	event := new(MerkleDistributorRoleGranted)
	if err := _MerkleDistributor.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MerkleDistributorRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the MerkleDistributor contract.
type MerkleDistributorRoleRevokedIterator struct {
	Event *MerkleDistributorRoleRevoked // Event containing the contract specifics and raw log

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
func (it *MerkleDistributorRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MerkleDistributorRoleRevoked)
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
		it.Event = new(MerkleDistributorRoleRevoked)
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
func (it *MerkleDistributorRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MerkleDistributorRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MerkleDistributorRoleRevoked represents a RoleRevoked event raised by the MerkleDistributor contract.
type MerkleDistributorRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_MerkleDistributor *MerkleDistributorFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*MerkleDistributorRoleRevokedIterator, error) {

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

	logs, sub, err := _MerkleDistributor.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &MerkleDistributorRoleRevokedIterator{contract: _MerkleDistributor.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_MerkleDistributor *MerkleDistributorFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *MerkleDistributorRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _MerkleDistributor.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MerkleDistributorRoleRevoked)
				if err := _MerkleDistributor.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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
func (_MerkleDistributor *MerkleDistributorFilterer) ParseRoleRevoked(log types.Log) (*MerkleDistributorRoleRevoked, error) {
	event := new(MerkleDistributorRoleRevoked)
	if err := _MerkleDistributor.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MerkleDistributorRootSetIterator is returned from FilterRootSet and is used to iterate over the raw logs and unpacked data for RootSet events raised by the MerkleDistributor contract.
type MerkleDistributorRootSetIterator struct {
	Event *MerkleDistributorRootSet // Event containing the contract specifics and raw log

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
func (it *MerkleDistributorRootSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MerkleDistributorRootSet)
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
		it.Event = new(MerkleDistributorRootSet)
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
func (it *MerkleDistributorRootSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MerkleDistributorRootSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MerkleDistributorRootSet represents a RootSet event raised by the MerkleDistributor contract.
type MerkleDistributorRootSet struct {
	PoolId     *big.Int
	Week       *big.Int
	Root       [32]byte
	Allocation *big.Int
	ProofsURI  string
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterRootSet is a free log retrieval operation binding the contract event 0x121a417ff6acf9dc80400d2b74677426b59add2a9d014e26fa79fddd2fe44120.
//
// Solidity: event RootSet(uint256 indexed poolId, uint256 indexed week, bytes32 root, uint256 allocation, string proofsURI)
func (_MerkleDistributor *MerkleDistributorFilterer) FilterRootSet(opts *bind.FilterOpts, poolId []*big.Int, week []*big.Int) (*MerkleDistributorRootSetIterator, error) {

	var poolIdRule []interface{}
	for _, poolIdItem := range poolId {
		poolIdRule = append(poolIdRule, poolIdItem)
	}
	var weekRule []interface{}
	for _, weekItem := range week {
		weekRule = append(weekRule, weekItem)
	}

	logs, sub, err := _MerkleDistributor.contract.FilterLogs(opts, "RootSet", poolIdRule, weekRule)
	if err != nil {
		return nil, err
	}
	return &MerkleDistributorRootSetIterator{contract: _MerkleDistributor.contract, event: "RootSet", logs: logs, sub: sub}, nil
}

// WatchRootSet is a free log subscription operation binding the contract event 0x121a417ff6acf9dc80400d2b74677426b59add2a9d014e26fa79fddd2fe44120.
//
// Solidity: event RootSet(uint256 indexed poolId, uint256 indexed week, bytes32 root, uint256 allocation, string proofsURI)
func (_MerkleDistributor *MerkleDistributorFilterer) WatchRootSet(opts *bind.WatchOpts, sink chan<- *MerkleDistributorRootSet, poolId []*big.Int, week []*big.Int) (event.Subscription, error) {

	var poolIdRule []interface{}
	for _, poolIdItem := range poolId {
		poolIdRule = append(poolIdRule, poolIdItem)
	}
	var weekRule []interface{}
	for _, weekItem := range week {
		weekRule = append(weekRule, weekItem)
	}

	logs, sub, err := _MerkleDistributor.contract.WatchLogs(opts, "RootSet", poolIdRule, weekRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MerkleDistributorRootSet)
				if err := _MerkleDistributor.contract.UnpackLog(event, "RootSet", log); err != nil {
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

// ParseRootSet is a log parse operation binding the contract event 0x121a417ff6acf9dc80400d2b74677426b59add2a9d014e26fa79fddd2fe44120.
//
// Solidity: event RootSet(uint256 indexed poolId, uint256 indexed week, bytes32 root, uint256 allocation, string proofsURI)
func (_MerkleDistributor *MerkleDistributorFilterer) ParseRootSet(log types.Log) (*MerkleDistributorRootSet, error) {
	event := new(MerkleDistributorRootSet)
	if err := _MerkleDistributor.contract.UnpackLog(event, "RootSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MerkleDistributorRootSetterSetIterator is returned from FilterRootSetterSet and is used to iterate over the raw logs and unpacked data for RootSetterSet events raised by the MerkleDistributor contract.
type MerkleDistributorRootSetterSetIterator struct {
	Event *MerkleDistributorRootSetterSet // Event containing the contract specifics and raw log

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
func (it *MerkleDistributorRootSetterSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MerkleDistributorRootSetterSet)
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
		it.Event = new(MerkleDistributorRootSetterSet)
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
func (it *MerkleDistributorRootSetterSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MerkleDistributorRootSetterSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MerkleDistributorRootSetterSet represents a RootSetterSet event raised by the MerkleDistributor contract.
type MerkleDistributorRootSetterSet struct {
	PoolId  *big.Int
	Setter  common.Address
	Enabled bool
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRootSetterSet is a free log retrieval operation binding the contract event 0x502dff0c4a37599b52ef85eab1e0c76e990f573eb54603a1aa58020e813bbc88.
//
// Solidity: event RootSetterSet(uint256 indexed poolId, address indexed setter, bool enabled)
func (_MerkleDistributor *MerkleDistributorFilterer) FilterRootSetterSet(opts *bind.FilterOpts, poolId []*big.Int, setter []common.Address) (*MerkleDistributorRootSetterSetIterator, error) {

	var poolIdRule []interface{}
	for _, poolIdItem := range poolId {
		poolIdRule = append(poolIdRule, poolIdItem)
	}
	var setterRule []interface{}
	for _, setterItem := range setter {
		setterRule = append(setterRule, setterItem)
	}

	logs, sub, err := _MerkleDistributor.contract.FilterLogs(opts, "RootSetterSet", poolIdRule, setterRule)
	if err != nil {
		return nil, err
	}
	return &MerkleDistributorRootSetterSetIterator{contract: _MerkleDistributor.contract, event: "RootSetterSet", logs: logs, sub: sub}, nil
}

// WatchRootSetterSet is a free log subscription operation binding the contract event 0x502dff0c4a37599b52ef85eab1e0c76e990f573eb54603a1aa58020e813bbc88.
//
// Solidity: event RootSetterSet(uint256 indexed poolId, address indexed setter, bool enabled)
func (_MerkleDistributor *MerkleDistributorFilterer) WatchRootSetterSet(opts *bind.WatchOpts, sink chan<- *MerkleDistributorRootSetterSet, poolId []*big.Int, setter []common.Address) (event.Subscription, error) {

	var poolIdRule []interface{}
	for _, poolIdItem := range poolId {
		poolIdRule = append(poolIdRule, poolIdItem)
	}
	var setterRule []interface{}
	for _, setterItem := range setter {
		setterRule = append(setterRule, setterItem)
	}

	logs, sub, err := _MerkleDistributor.contract.WatchLogs(opts, "RootSetterSet", poolIdRule, setterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MerkleDistributorRootSetterSet)
				if err := _MerkleDistributor.contract.UnpackLog(event, "RootSetterSet", log); err != nil {
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

// ParseRootSetterSet is a log parse operation binding the contract event 0x502dff0c4a37599b52ef85eab1e0c76e990f573eb54603a1aa58020e813bbc88.
//
// Solidity: event RootSetterSet(uint256 indexed poolId, address indexed setter, bool enabled)
func (_MerkleDistributor *MerkleDistributorFilterer) ParseRootSetterSet(log types.Log) (*MerkleDistributorRootSetterSet, error) {
	event := new(MerkleDistributorRootSetterSet)
	if err := _MerkleDistributor.contract.UnpackLog(event, "RootSetterSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MerkleDistributorSweptIterator is returned from FilterSwept and is used to iterate over the raw logs and unpacked data for Swept events raised by the MerkleDistributor contract.
type MerkleDistributorSweptIterator struct {
	Event *MerkleDistributorSwept // Event containing the contract specifics and raw log

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
func (it *MerkleDistributorSweptIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MerkleDistributorSwept)
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
		it.Event = new(MerkleDistributorSwept)
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
func (it *MerkleDistributorSweptIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MerkleDistributorSweptIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MerkleDistributorSwept represents a Swept event raised by the MerkleDistributor contract.
type MerkleDistributorSwept struct {
	PoolId     *big.Int
	To         common.Address
	Amount     *big.Int
	NewBalance *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterSwept is a free log retrieval operation binding the contract event 0xebadffe1424633afd93dae56bfe8d73282a6fd80c8171ad9bb4982b6c40bf228.
//
// Solidity: event Swept(uint256 indexed poolId, address indexed to, uint256 amount, uint256 newBalance)
func (_MerkleDistributor *MerkleDistributorFilterer) FilterSwept(opts *bind.FilterOpts, poolId []*big.Int, to []common.Address) (*MerkleDistributorSweptIterator, error) {

	var poolIdRule []interface{}
	for _, poolIdItem := range poolId {
		poolIdRule = append(poolIdRule, poolIdItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _MerkleDistributor.contract.FilterLogs(opts, "Swept", poolIdRule, toRule)
	if err != nil {
		return nil, err
	}
	return &MerkleDistributorSweptIterator{contract: _MerkleDistributor.contract, event: "Swept", logs: logs, sub: sub}, nil
}

// WatchSwept is a free log subscription operation binding the contract event 0xebadffe1424633afd93dae56bfe8d73282a6fd80c8171ad9bb4982b6c40bf228.
//
// Solidity: event Swept(uint256 indexed poolId, address indexed to, uint256 amount, uint256 newBalance)
func (_MerkleDistributor *MerkleDistributorFilterer) WatchSwept(opts *bind.WatchOpts, sink chan<- *MerkleDistributorSwept, poolId []*big.Int, to []common.Address) (event.Subscription, error) {

	var poolIdRule []interface{}
	for _, poolIdItem := range poolId {
		poolIdRule = append(poolIdRule, poolIdItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _MerkleDistributor.contract.WatchLogs(opts, "Swept", poolIdRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MerkleDistributorSwept)
				if err := _MerkleDistributor.contract.UnpackLog(event, "Swept", log); err != nil {
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

// ParseSwept is a log parse operation binding the contract event 0xebadffe1424633afd93dae56bfe8d73282a6fd80c8171ad9bb4982b6c40bf228.
//
// Solidity: event Swept(uint256 indexed poolId, address indexed to, uint256 amount, uint256 newBalance)
func (_MerkleDistributor *MerkleDistributorFilterer) ParseSwept(log types.Log) (*MerkleDistributorSwept, error) {
	event := new(MerkleDistributorSwept)
	if err := _MerkleDistributor.contract.UnpackLog(event, "Swept", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MerkleDistributorUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the MerkleDistributor contract.
type MerkleDistributorUnpausedIterator struct {
	Event *MerkleDistributorUnpaused // Event containing the contract specifics and raw log

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
func (it *MerkleDistributorUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MerkleDistributorUnpaused)
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
		it.Event = new(MerkleDistributorUnpaused)
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
func (it *MerkleDistributorUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MerkleDistributorUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MerkleDistributorUnpaused represents a Unpaused event raised by the MerkleDistributor contract.
type MerkleDistributorUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_MerkleDistributor *MerkleDistributorFilterer) FilterUnpaused(opts *bind.FilterOpts) (*MerkleDistributorUnpausedIterator, error) {

	logs, sub, err := _MerkleDistributor.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &MerkleDistributorUnpausedIterator{contract: _MerkleDistributor.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_MerkleDistributor *MerkleDistributorFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *MerkleDistributorUnpaused) (event.Subscription, error) {

	logs, sub, err := _MerkleDistributor.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MerkleDistributorUnpaused)
				if err := _MerkleDistributor.contract.UnpackLog(event, "Unpaused", log); err != nil {
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

// ParseUnpaused is a log parse operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_MerkleDistributor *MerkleDistributorFilterer) ParseUnpaused(log types.Log) (*MerkleDistributorUnpaused, error) {
	event := new(MerkleDistributorUnpaused)
	if err := _MerkleDistributor.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MerkleDistributorUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the MerkleDistributor contract.
type MerkleDistributorUpgradedIterator struct {
	Event *MerkleDistributorUpgraded // Event containing the contract specifics and raw log

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
func (it *MerkleDistributorUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MerkleDistributorUpgraded)
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
		it.Event = new(MerkleDistributorUpgraded)
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
func (it *MerkleDistributorUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MerkleDistributorUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MerkleDistributorUpgraded represents a Upgraded event raised by the MerkleDistributor contract.
type MerkleDistributorUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_MerkleDistributor *MerkleDistributorFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*MerkleDistributorUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _MerkleDistributor.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &MerkleDistributorUpgradedIterator{contract: _MerkleDistributor.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_MerkleDistributor *MerkleDistributorFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *MerkleDistributorUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _MerkleDistributor.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MerkleDistributorUpgraded)
				if err := _MerkleDistributor.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_MerkleDistributor *MerkleDistributorFilterer) ParseUpgraded(log types.Log) (*MerkleDistributorUpgraded, error) {
	event := new(MerkleDistributorUpgraded)
	if err := _MerkleDistributor.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MerkleDistributorWeeklyLimitSetIterator is returned from FilterWeeklyLimitSet and is used to iterate over the raw logs and unpacked data for WeeklyLimitSet events raised by the MerkleDistributor contract.
type MerkleDistributorWeeklyLimitSetIterator struct {
	Event *MerkleDistributorWeeklyLimitSet // Event containing the contract specifics and raw log

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
func (it *MerkleDistributorWeeklyLimitSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MerkleDistributorWeeklyLimitSet)
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
		it.Event = new(MerkleDistributorWeeklyLimitSet)
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
func (it *MerkleDistributorWeeklyLimitSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MerkleDistributorWeeklyLimitSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MerkleDistributorWeeklyLimitSet represents a WeeklyLimitSet event raised by the MerkleDistributor contract.
type MerkleDistributorWeeklyLimitSet struct {
	PoolId *big.Int
	Limit  *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterWeeklyLimitSet is a free log retrieval operation binding the contract event 0xf6bdd6278a6ce93210cad414aba14b52a80ae25fcc64b9b2a91b6c1a19bd8420.
//
// Solidity: event WeeklyLimitSet(uint256 indexed poolId, uint256 limit)
func (_MerkleDistributor *MerkleDistributorFilterer) FilterWeeklyLimitSet(opts *bind.FilterOpts, poolId []*big.Int) (*MerkleDistributorWeeklyLimitSetIterator, error) {

	var poolIdRule []interface{}
	for _, poolIdItem := range poolId {
		poolIdRule = append(poolIdRule, poolIdItem)
	}

	logs, sub, err := _MerkleDistributor.contract.FilterLogs(opts, "WeeklyLimitSet", poolIdRule)
	if err != nil {
		return nil, err
	}
	return &MerkleDistributorWeeklyLimitSetIterator{contract: _MerkleDistributor.contract, event: "WeeklyLimitSet", logs: logs, sub: sub}, nil
}

// WatchWeeklyLimitSet is a free log subscription operation binding the contract event 0xf6bdd6278a6ce93210cad414aba14b52a80ae25fcc64b9b2a91b6c1a19bd8420.
//
// Solidity: event WeeklyLimitSet(uint256 indexed poolId, uint256 limit)
func (_MerkleDistributor *MerkleDistributorFilterer) WatchWeeklyLimitSet(opts *bind.WatchOpts, sink chan<- *MerkleDistributorWeeklyLimitSet, poolId []*big.Int) (event.Subscription, error) {

	var poolIdRule []interface{}
	for _, poolIdItem := range poolId {
		poolIdRule = append(poolIdRule, poolIdItem)
	}

	logs, sub, err := _MerkleDistributor.contract.WatchLogs(opts, "WeeklyLimitSet", poolIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MerkleDistributorWeeklyLimitSet)
				if err := _MerkleDistributor.contract.UnpackLog(event, "WeeklyLimitSet", log); err != nil {
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

// ParseWeeklyLimitSet is a log parse operation binding the contract event 0xf6bdd6278a6ce93210cad414aba14b52a80ae25fcc64b9b2a91b6c1a19bd8420.
//
// Solidity: event WeeklyLimitSet(uint256 indexed poolId, uint256 limit)
func (_MerkleDistributor *MerkleDistributorFilterer) ParseWeeklyLimitSet(log types.Log) (*MerkleDistributorWeeklyLimitSet, error) {
	event := new(MerkleDistributorWeeklyLimitSet)
	if err := _MerkleDistributor.contract.UnpackLog(event, "WeeklyLimitSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
