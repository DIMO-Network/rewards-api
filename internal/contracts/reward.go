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

// RewardTransferInfo is an auto generated low-level Go binding around an user-defined struct.
type RewardTransferInfo struct {
	User                       common.Address
	VehicleId                  *big.Int
	AftermarketDeviceId        *big.Int
	ValueFromAftermarketDevice *big.Int
	SyntheticDeviceId          *big.Int
	ValueFromSyntheticDevice   *big.Int
	ConnectionStreak           *big.Int
	ValueFromStreak            *big.Int
}

// RewardMetaData contains all meta data concerning the Reward contract.
var RewardMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"TokenTransferFailed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"WeeklyLimitExceeded\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"previousAdmin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"AdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"AdminWithdrawal\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"beacon\",\"type\":\"address\"}],\"name\":\"BeaconUpgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"vehicleNodeId\",\"type\":\"uint256\"}],\"name\":\"DidntQualify\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"vehicleNodeId\",\"type\":\"uint256\"}],\"name\":\"TokensTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"week\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"vehicleNodeId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"connectionStreak\",\"type\":\"uint256\"}],\"name\":\"TokensTransferredForConnectionStreak\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"week\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"vehicleNodeId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"deviceNftProxy\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"deviceNode\",\"type\":\"uint256\"}],\"name\":\"TokensTransferredForDevice\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newWeek\",\"type\":\"uint256\"}],\"name\":\"WeekAdvanced\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"ORACLE_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"adminWithdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"aftermarketIdContractAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"vehicleId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"aftermarketDeviceId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"valueFromAftermarketDevice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"syntheticDeviceId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"valueFromSyntheticDevice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"connectionStreak\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"valueFromStreak\",\"type\":\"uint256\"}],\"internalType\":\"structReward.TransferInfo[]\",\"name\":\"transferInfos\",\"type\":\"tuple[]\"}],\"name\":\"batchTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"calculateWeeklyLimit\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"maxRewards\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"currentWeek\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"currentWeekLimit\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"currentWeekSpent\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dimoToken\",\"outputs\":[{\"internalType\":\"contractIDimo\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dimoTotalSentOutByContract\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialWeeklyLimit\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"vehicleIdProxyAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"sanctionsContractAddress\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"name\":\"manuallySetRewardsGenesisTime\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minimumTimeForRewards\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"proxiableUUID\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"registry\",\"outputs\":[{\"internalType\":\"contractIRegistry\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"resetRewardsGenesisTime\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rewardsGenesisTime\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"sanctionsList\",\"outputs\":[{\"internalType\":\"contractISanctionsList\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"aftermarketProxyAddress\",\"type\":\"address\"}],\"name\":\"setAftermarketProxyAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"newTimeInSeconds\",\"type\":\"uint256\"}],\"name\":\"setMinimumTimeForRewards\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"registryContractAddress\",\"type\":\"address\"}],\"name\":\"setRegistryContractAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"syntheticProxyAddress\",\"type\":\"address\"}],\"name\":\"setSyntheticProxyAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"syntheticIdContractAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"}],\"name\":\"upgradeTo\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"upgradeToAndCall\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"vehicleIdLastRewardTime\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"vehicleIdProxy\",\"outputs\":[{\"internalType\":\"contractINft\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// RewardABI is the input ABI used to generate the binding from.
// Deprecated: Use RewardMetaData.ABI instead.
var RewardABI = RewardMetaData.ABI

// Reward is an auto generated Go binding around an Ethereum contract.
type Reward struct {
	RewardCaller     // Read-only binding to the contract
	RewardTransactor // Write-only binding to the contract
	RewardFilterer   // Log filterer for contract events
}

// RewardCaller is an auto generated read-only Go binding around an Ethereum contract.
type RewardCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RewardTransactor is an auto generated write-only Go binding around an Ethereum contract.
type RewardTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RewardFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type RewardFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RewardSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type RewardSession struct {
	Contract     *Reward           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RewardCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type RewardCallerSession struct {
	Contract *RewardCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// RewardTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type RewardTransactorSession struct {
	Contract     *RewardTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RewardRaw is an auto generated low-level Go binding around an Ethereum contract.
type RewardRaw struct {
	Contract *Reward // Generic contract binding to access the raw methods on
}

// RewardCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type RewardCallerRaw struct {
	Contract *RewardCaller // Generic read-only contract binding to access the raw methods on
}

// RewardTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type RewardTransactorRaw struct {
	Contract *RewardTransactor // Generic write-only contract binding to access the raw methods on
}

// NewReward creates a new instance of Reward, bound to a specific deployed contract.
func NewReward(address common.Address, backend bind.ContractBackend) (*Reward, error) {
	contract, err := bindReward(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Reward{RewardCaller: RewardCaller{contract: contract}, RewardTransactor: RewardTransactor{contract: contract}, RewardFilterer: RewardFilterer{contract: contract}}, nil
}

// NewRewardCaller creates a new read-only instance of Reward, bound to a specific deployed contract.
func NewRewardCaller(address common.Address, caller bind.ContractCaller) (*RewardCaller, error) {
	contract, err := bindReward(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RewardCaller{contract: contract}, nil
}

// NewRewardTransactor creates a new write-only instance of Reward, bound to a specific deployed contract.
func NewRewardTransactor(address common.Address, transactor bind.ContractTransactor) (*RewardTransactor, error) {
	contract, err := bindReward(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RewardTransactor{contract: contract}, nil
}

// NewRewardFilterer creates a new log filterer instance of Reward, bound to a specific deployed contract.
func NewRewardFilterer(address common.Address, filterer bind.ContractFilterer) (*RewardFilterer, error) {
	contract, err := bindReward(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RewardFilterer{contract: contract}, nil
}

// bindReward binds a generic wrapper to an already deployed contract.
func bindReward(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := RewardMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Reward *RewardRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Reward.Contract.RewardCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Reward *RewardRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Reward.Contract.RewardTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Reward *RewardRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Reward.Contract.RewardTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Reward *RewardCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Reward.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Reward *RewardTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Reward.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Reward *RewardTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Reward.Contract.contract.Transact(opts, method, params...)
}

// ADMINROLE is a free data retrieval call binding the contract method 0x75b238fc.
//
// Solidity: function ADMIN_ROLE() view returns(bytes32)
func (_Reward *RewardCaller) ADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Reward.contract.Call(opts, &out, "ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ADMINROLE is a free data retrieval call binding the contract method 0x75b238fc.
//
// Solidity: function ADMIN_ROLE() view returns(bytes32)
func (_Reward *RewardSession) ADMINROLE() ([32]byte, error) {
	return _Reward.Contract.ADMINROLE(&_Reward.CallOpts)
}

// ADMINROLE is a free data retrieval call binding the contract method 0x75b238fc.
//
// Solidity: function ADMIN_ROLE() view returns(bytes32)
func (_Reward *RewardCallerSession) ADMINROLE() ([32]byte, error) {
	return _Reward.Contract.ADMINROLE(&_Reward.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_Reward *RewardCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Reward.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_Reward *RewardSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _Reward.Contract.DEFAULTADMINROLE(&_Reward.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_Reward *RewardCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _Reward.Contract.DEFAULTADMINROLE(&_Reward.CallOpts)
}

// ORACLEROLE is a free data retrieval call binding the contract method 0x07e2cea5.
//
// Solidity: function ORACLE_ROLE() view returns(bytes32)
func (_Reward *RewardCaller) ORACLEROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Reward.contract.Call(opts, &out, "ORACLE_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ORACLEROLE is a free data retrieval call binding the contract method 0x07e2cea5.
//
// Solidity: function ORACLE_ROLE() view returns(bytes32)
func (_Reward *RewardSession) ORACLEROLE() ([32]byte, error) {
	return _Reward.Contract.ORACLEROLE(&_Reward.CallOpts)
}

// ORACLEROLE is a free data retrieval call binding the contract method 0x07e2cea5.
//
// Solidity: function ORACLE_ROLE() view returns(bytes32)
func (_Reward *RewardCallerSession) ORACLEROLE() ([32]byte, error) {
	return _Reward.Contract.ORACLEROLE(&_Reward.CallOpts)
}

// AftermarketIdContractAddress is a free data retrieval call binding the contract method 0xd849d652.
//
// Solidity: function aftermarketIdContractAddress() view returns(address)
func (_Reward *RewardCaller) AftermarketIdContractAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Reward.contract.Call(opts, &out, "aftermarketIdContractAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// AftermarketIdContractAddress is a free data retrieval call binding the contract method 0xd849d652.
//
// Solidity: function aftermarketIdContractAddress() view returns(address)
func (_Reward *RewardSession) AftermarketIdContractAddress() (common.Address, error) {
	return _Reward.Contract.AftermarketIdContractAddress(&_Reward.CallOpts)
}

// AftermarketIdContractAddress is a free data retrieval call binding the contract method 0xd849d652.
//
// Solidity: function aftermarketIdContractAddress() view returns(address)
func (_Reward *RewardCallerSession) AftermarketIdContractAddress() (common.Address, error) {
	return _Reward.Contract.AftermarketIdContractAddress(&_Reward.CallOpts)
}

// CalculateWeeklyLimit is a free data retrieval call binding the contract method 0x60e1f139.
//
// Solidity: function calculateWeeklyLimit() view returns(uint256 maxRewards)
func (_Reward *RewardCaller) CalculateWeeklyLimit(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Reward.contract.Call(opts, &out, "calculateWeeklyLimit")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CalculateWeeklyLimit is a free data retrieval call binding the contract method 0x60e1f139.
//
// Solidity: function calculateWeeklyLimit() view returns(uint256 maxRewards)
func (_Reward *RewardSession) CalculateWeeklyLimit() (*big.Int, error) {
	return _Reward.Contract.CalculateWeeklyLimit(&_Reward.CallOpts)
}

// CalculateWeeklyLimit is a free data retrieval call binding the contract method 0x60e1f139.
//
// Solidity: function calculateWeeklyLimit() view returns(uint256 maxRewards)
func (_Reward *RewardCallerSession) CalculateWeeklyLimit() (*big.Int, error) {
	return _Reward.Contract.CalculateWeeklyLimit(&_Reward.CallOpts)
}

// CurrentWeek is a free data retrieval call binding the contract method 0x06575c89.
//
// Solidity: function currentWeek() view returns(uint256)
func (_Reward *RewardCaller) CurrentWeek(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Reward.contract.Call(opts, &out, "currentWeek")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CurrentWeek is a free data retrieval call binding the contract method 0x06575c89.
//
// Solidity: function currentWeek() view returns(uint256)
func (_Reward *RewardSession) CurrentWeek() (*big.Int, error) {
	return _Reward.Contract.CurrentWeek(&_Reward.CallOpts)
}

// CurrentWeek is a free data retrieval call binding the contract method 0x06575c89.
//
// Solidity: function currentWeek() view returns(uint256)
func (_Reward *RewardCallerSession) CurrentWeek() (*big.Int, error) {
	return _Reward.Contract.CurrentWeek(&_Reward.CallOpts)
}

// CurrentWeekLimit is a free data retrieval call binding the contract method 0x9cb12ca2.
//
// Solidity: function currentWeekLimit() view returns(uint256)
func (_Reward *RewardCaller) CurrentWeekLimit(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Reward.contract.Call(opts, &out, "currentWeekLimit")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CurrentWeekLimit is a free data retrieval call binding the contract method 0x9cb12ca2.
//
// Solidity: function currentWeekLimit() view returns(uint256)
func (_Reward *RewardSession) CurrentWeekLimit() (*big.Int, error) {
	return _Reward.Contract.CurrentWeekLimit(&_Reward.CallOpts)
}

// CurrentWeekLimit is a free data retrieval call binding the contract method 0x9cb12ca2.
//
// Solidity: function currentWeekLimit() view returns(uint256)
func (_Reward *RewardCallerSession) CurrentWeekLimit() (*big.Int, error) {
	return _Reward.Contract.CurrentWeekLimit(&_Reward.CallOpts)
}

// CurrentWeekSpent is a free data retrieval call binding the contract method 0x7802674d.
//
// Solidity: function currentWeekSpent() view returns(uint256)
func (_Reward *RewardCaller) CurrentWeekSpent(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Reward.contract.Call(opts, &out, "currentWeekSpent")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CurrentWeekSpent is a free data retrieval call binding the contract method 0x7802674d.
//
// Solidity: function currentWeekSpent() view returns(uint256)
func (_Reward *RewardSession) CurrentWeekSpent() (*big.Int, error) {
	return _Reward.Contract.CurrentWeekSpent(&_Reward.CallOpts)
}

// CurrentWeekSpent is a free data retrieval call binding the contract method 0x7802674d.
//
// Solidity: function currentWeekSpent() view returns(uint256)
func (_Reward *RewardCallerSession) CurrentWeekSpent() (*big.Int, error) {
	return _Reward.Contract.CurrentWeekSpent(&_Reward.CallOpts)
}

// DimoToken is a free data retrieval call binding the contract method 0x0524f28c.
//
// Solidity: function dimoToken() view returns(address)
func (_Reward *RewardCaller) DimoToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Reward.contract.Call(opts, &out, "dimoToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// DimoToken is a free data retrieval call binding the contract method 0x0524f28c.
//
// Solidity: function dimoToken() view returns(address)
func (_Reward *RewardSession) DimoToken() (common.Address, error) {
	return _Reward.Contract.DimoToken(&_Reward.CallOpts)
}

// DimoToken is a free data retrieval call binding the contract method 0x0524f28c.
//
// Solidity: function dimoToken() view returns(address)
func (_Reward *RewardCallerSession) DimoToken() (common.Address, error) {
	return _Reward.Contract.DimoToken(&_Reward.CallOpts)
}

// DimoTotalSentOutByContract is a free data retrieval call binding the contract method 0x7f350fea.
//
// Solidity: function dimoTotalSentOutByContract() view returns(uint256)
func (_Reward *RewardCaller) DimoTotalSentOutByContract(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Reward.contract.Call(opts, &out, "dimoTotalSentOutByContract")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DimoTotalSentOutByContract is a free data retrieval call binding the contract method 0x7f350fea.
//
// Solidity: function dimoTotalSentOutByContract() view returns(uint256)
func (_Reward *RewardSession) DimoTotalSentOutByContract() (*big.Int, error) {
	return _Reward.Contract.DimoTotalSentOutByContract(&_Reward.CallOpts)
}

// DimoTotalSentOutByContract is a free data retrieval call binding the contract method 0x7f350fea.
//
// Solidity: function dimoTotalSentOutByContract() view returns(uint256)
func (_Reward *RewardCallerSession) DimoTotalSentOutByContract() (*big.Int, error) {
	return _Reward.Contract.DimoTotalSentOutByContract(&_Reward.CallOpts)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_Reward *RewardCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _Reward.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_Reward *RewardSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _Reward.Contract.GetRoleAdmin(&_Reward.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_Reward *RewardCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _Reward.Contract.GetRoleAdmin(&_Reward.CallOpts, role)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_Reward *RewardCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _Reward.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_Reward *RewardSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _Reward.Contract.HasRole(&_Reward.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_Reward *RewardCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _Reward.Contract.HasRole(&_Reward.CallOpts, role, account)
}

// InitialWeeklyLimit is a free data retrieval call binding the contract method 0x611fbc96.
//
// Solidity: function initialWeeklyLimit() view returns(uint256)
func (_Reward *RewardCaller) InitialWeeklyLimit(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Reward.contract.Call(opts, &out, "initialWeeklyLimit")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// InitialWeeklyLimit is a free data retrieval call binding the contract method 0x611fbc96.
//
// Solidity: function initialWeeklyLimit() view returns(uint256)
func (_Reward *RewardSession) InitialWeeklyLimit() (*big.Int, error) {
	return _Reward.Contract.InitialWeeklyLimit(&_Reward.CallOpts)
}

// InitialWeeklyLimit is a free data retrieval call binding the contract method 0x611fbc96.
//
// Solidity: function initialWeeklyLimit() view returns(uint256)
func (_Reward *RewardCallerSession) InitialWeeklyLimit() (*big.Int, error) {
	return _Reward.Contract.InitialWeeklyLimit(&_Reward.CallOpts)
}

// MinimumTimeForRewards is a free data retrieval call binding the contract method 0x9084025e.
//
// Solidity: function minimumTimeForRewards() view returns(uint256)
func (_Reward *RewardCaller) MinimumTimeForRewards(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Reward.contract.Call(opts, &out, "minimumTimeForRewards")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MinimumTimeForRewards is a free data retrieval call binding the contract method 0x9084025e.
//
// Solidity: function minimumTimeForRewards() view returns(uint256)
func (_Reward *RewardSession) MinimumTimeForRewards() (*big.Int, error) {
	return _Reward.Contract.MinimumTimeForRewards(&_Reward.CallOpts)
}

// MinimumTimeForRewards is a free data retrieval call binding the contract method 0x9084025e.
//
// Solidity: function minimumTimeForRewards() view returns(uint256)
func (_Reward *RewardCallerSession) MinimumTimeForRewards() (*big.Int, error) {
	return _Reward.Contract.MinimumTimeForRewards(&_Reward.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Reward *RewardCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Reward.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Reward *RewardSession) ProxiableUUID() ([32]byte, error) {
	return _Reward.Contract.ProxiableUUID(&_Reward.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Reward *RewardCallerSession) ProxiableUUID() ([32]byte, error) {
	return _Reward.Contract.ProxiableUUID(&_Reward.CallOpts)
}

// Registry is a free data retrieval call binding the contract method 0x7b103999.
//
// Solidity: function registry() view returns(address)
func (_Reward *RewardCaller) Registry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Reward.contract.Call(opts, &out, "registry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Registry is a free data retrieval call binding the contract method 0x7b103999.
//
// Solidity: function registry() view returns(address)
func (_Reward *RewardSession) Registry() (common.Address, error) {
	return _Reward.Contract.Registry(&_Reward.CallOpts)
}

// Registry is a free data retrieval call binding the contract method 0x7b103999.
//
// Solidity: function registry() view returns(address)
func (_Reward *RewardCallerSession) Registry() (common.Address, error) {
	return _Reward.Contract.Registry(&_Reward.CallOpts)
}

// RewardsGenesisTime is a free data retrieval call binding the contract method 0x7f2842b4.
//
// Solidity: function rewardsGenesisTime() view returns(uint256)
func (_Reward *RewardCaller) RewardsGenesisTime(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Reward.contract.Call(opts, &out, "rewardsGenesisTime")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// RewardsGenesisTime is a free data retrieval call binding the contract method 0x7f2842b4.
//
// Solidity: function rewardsGenesisTime() view returns(uint256)
func (_Reward *RewardSession) RewardsGenesisTime() (*big.Int, error) {
	return _Reward.Contract.RewardsGenesisTime(&_Reward.CallOpts)
}

// RewardsGenesisTime is a free data retrieval call binding the contract method 0x7f2842b4.
//
// Solidity: function rewardsGenesisTime() view returns(uint256)
func (_Reward *RewardCallerSession) RewardsGenesisTime() (*big.Int, error) {
	return _Reward.Contract.RewardsGenesisTime(&_Reward.CallOpts)
}

// SanctionsList is a free data retrieval call binding the contract method 0xec571c6a.
//
// Solidity: function sanctionsList() view returns(address)
func (_Reward *RewardCaller) SanctionsList(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Reward.contract.Call(opts, &out, "sanctionsList")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// SanctionsList is a free data retrieval call binding the contract method 0xec571c6a.
//
// Solidity: function sanctionsList() view returns(address)
func (_Reward *RewardSession) SanctionsList() (common.Address, error) {
	return _Reward.Contract.SanctionsList(&_Reward.CallOpts)
}

// SanctionsList is a free data retrieval call binding the contract method 0xec571c6a.
//
// Solidity: function sanctionsList() view returns(address)
func (_Reward *RewardCallerSession) SanctionsList() (common.Address, error) {
	return _Reward.Contract.SanctionsList(&_Reward.CallOpts)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Reward *RewardCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _Reward.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Reward *RewardSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _Reward.Contract.SupportsInterface(&_Reward.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Reward *RewardCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _Reward.Contract.SupportsInterface(&_Reward.CallOpts, interfaceId)
}

// SyntheticIdContractAddress is a free data retrieval call binding the contract method 0x3fddd852.
//
// Solidity: function syntheticIdContractAddress() view returns(address)
func (_Reward *RewardCaller) SyntheticIdContractAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Reward.contract.Call(opts, &out, "syntheticIdContractAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// SyntheticIdContractAddress is a free data retrieval call binding the contract method 0x3fddd852.
//
// Solidity: function syntheticIdContractAddress() view returns(address)
func (_Reward *RewardSession) SyntheticIdContractAddress() (common.Address, error) {
	return _Reward.Contract.SyntheticIdContractAddress(&_Reward.CallOpts)
}

// SyntheticIdContractAddress is a free data retrieval call binding the contract method 0x3fddd852.
//
// Solidity: function syntheticIdContractAddress() view returns(address)
func (_Reward *RewardCallerSession) SyntheticIdContractAddress() (common.Address, error) {
	return _Reward.Contract.SyntheticIdContractAddress(&_Reward.CallOpts)
}

// VehicleIdLastRewardTime is a free data retrieval call binding the contract method 0xdf959d01.
//
// Solidity: function vehicleIdLastRewardTime(uint256 ) view returns(uint256)
func (_Reward *RewardCaller) VehicleIdLastRewardTime(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Reward.contract.Call(opts, &out, "vehicleIdLastRewardTime", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// VehicleIdLastRewardTime is a free data retrieval call binding the contract method 0xdf959d01.
//
// Solidity: function vehicleIdLastRewardTime(uint256 ) view returns(uint256)
func (_Reward *RewardSession) VehicleIdLastRewardTime(arg0 *big.Int) (*big.Int, error) {
	return _Reward.Contract.VehicleIdLastRewardTime(&_Reward.CallOpts, arg0)
}

// VehicleIdLastRewardTime is a free data retrieval call binding the contract method 0xdf959d01.
//
// Solidity: function vehicleIdLastRewardTime(uint256 ) view returns(uint256)
func (_Reward *RewardCallerSession) VehicleIdLastRewardTime(arg0 *big.Int) (*big.Int, error) {
	return _Reward.Contract.VehicleIdLastRewardTime(&_Reward.CallOpts, arg0)
}

// VehicleIdProxy is a free data retrieval call binding the contract method 0x40c05794.
//
// Solidity: function vehicleIdProxy() view returns(address)
func (_Reward *RewardCaller) VehicleIdProxy(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Reward.contract.Call(opts, &out, "vehicleIdProxy")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// VehicleIdProxy is a free data retrieval call binding the contract method 0x40c05794.
//
// Solidity: function vehicleIdProxy() view returns(address)
func (_Reward *RewardSession) VehicleIdProxy() (common.Address, error) {
	return _Reward.Contract.VehicleIdProxy(&_Reward.CallOpts)
}

// VehicleIdProxy is a free data retrieval call binding the contract method 0x40c05794.
//
// Solidity: function vehicleIdProxy() view returns(address)
func (_Reward *RewardCallerSession) VehicleIdProxy() (common.Address, error) {
	return _Reward.Contract.VehicleIdProxy(&_Reward.CallOpts)
}

// AdminWithdraw is a paid mutator transaction binding the contract method 0x401d4482.
//
// Solidity: function adminWithdraw(address user, uint256 amount) returns()
func (_Reward *RewardTransactor) AdminWithdraw(opts *bind.TransactOpts, user common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Reward.contract.Transact(opts, "adminWithdraw", user, amount)
}

// AdminWithdraw is a paid mutator transaction binding the contract method 0x401d4482.
//
// Solidity: function adminWithdraw(address user, uint256 amount) returns()
func (_Reward *RewardSession) AdminWithdraw(user common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Reward.Contract.AdminWithdraw(&_Reward.TransactOpts, user, amount)
}

// AdminWithdraw is a paid mutator transaction binding the contract method 0x401d4482.
//
// Solidity: function adminWithdraw(address user, uint256 amount) returns()
func (_Reward *RewardTransactorSession) AdminWithdraw(user common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Reward.Contract.AdminWithdraw(&_Reward.TransactOpts, user, amount)
}

// BatchTransfer is a paid mutator transaction binding the contract method 0x13cc10c9.
//
// Solidity: function batchTransfer((address,uint256,uint256,uint256,uint256,uint256,uint256,uint256)[] transferInfos) returns()
func (_Reward *RewardTransactor) BatchTransfer(opts *bind.TransactOpts, transferInfos []RewardTransferInfo) (*types.Transaction, error) {
	return _Reward.contract.Transact(opts, "batchTransfer", transferInfos)
}

// BatchTransfer is a paid mutator transaction binding the contract method 0x13cc10c9.
//
// Solidity: function batchTransfer((address,uint256,uint256,uint256,uint256,uint256,uint256,uint256)[] transferInfos) returns()
func (_Reward *RewardSession) BatchTransfer(transferInfos []RewardTransferInfo) (*types.Transaction, error) {
	return _Reward.Contract.BatchTransfer(&_Reward.TransactOpts, transferInfos)
}

// BatchTransfer is a paid mutator transaction binding the contract method 0x13cc10c9.
//
// Solidity: function batchTransfer((address,uint256,uint256,uint256,uint256,uint256,uint256,uint256)[] transferInfos) returns()
func (_Reward *RewardTransactorSession) BatchTransfer(transferInfos []RewardTransferInfo) (*types.Transaction, error) {
	return _Reward.Contract.BatchTransfer(&_Reward.TransactOpts, transferInfos)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_Reward *RewardTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Reward.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_Reward *RewardSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Reward.Contract.GrantRole(&_Reward.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_Reward *RewardTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Reward.Contract.GrantRole(&_Reward.TransactOpts, role, account)
}

// Initialize is a paid mutator transaction binding the contract method 0xc0c53b8b.
//
// Solidity: function initialize(address tokenAddress, address vehicleIdProxyAddress, address sanctionsContractAddress) returns()
func (_Reward *RewardTransactor) Initialize(opts *bind.TransactOpts, tokenAddress common.Address, vehicleIdProxyAddress common.Address, sanctionsContractAddress common.Address) (*types.Transaction, error) {
	return _Reward.contract.Transact(opts, "initialize", tokenAddress, vehicleIdProxyAddress, sanctionsContractAddress)
}

// Initialize is a paid mutator transaction binding the contract method 0xc0c53b8b.
//
// Solidity: function initialize(address tokenAddress, address vehicleIdProxyAddress, address sanctionsContractAddress) returns()
func (_Reward *RewardSession) Initialize(tokenAddress common.Address, vehicleIdProxyAddress common.Address, sanctionsContractAddress common.Address) (*types.Transaction, error) {
	return _Reward.Contract.Initialize(&_Reward.TransactOpts, tokenAddress, vehicleIdProxyAddress, sanctionsContractAddress)
}

// Initialize is a paid mutator transaction binding the contract method 0xc0c53b8b.
//
// Solidity: function initialize(address tokenAddress, address vehicleIdProxyAddress, address sanctionsContractAddress) returns()
func (_Reward *RewardTransactorSession) Initialize(tokenAddress common.Address, vehicleIdProxyAddress common.Address, sanctionsContractAddress common.Address) (*types.Transaction, error) {
	return _Reward.Contract.Initialize(&_Reward.TransactOpts, tokenAddress, vehicleIdProxyAddress, sanctionsContractAddress)
}

// ManuallySetRewardsGenesisTime is a paid mutator transaction binding the contract method 0x132153c7.
//
// Solidity: function manuallySetRewardsGenesisTime(uint256 timestamp) returns()
func (_Reward *RewardTransactor) ManuallySetRewardsGenesisTime(opts *bind.TransactOpts, timestamp *big.Int) (*types.Transaction, error) {
	return _Reward.contract.Transact(opts, "manuallySetRewardsGenesisTime", timestamp)
}

// ManuallySetRewardsGenesisTime is a paid mutator transaction binding the contract method 0x132153c7.
//
// Solidity: function manuallySetRewardsGenesisTime(uint256 timestamp) returns()
func (_Reward *RewardSession) ManuallySetRewardsGenesisTime(timestamp *big.Int) (*types.Transaction, error) {
	return _Reward.Contract.ManuallySetRewardsGenesisTime(&_Reward.TransactOpts, timestamp)
}

// ManuallySetRewardsGenesisTime is a paid mutator transaction binding the contract method 0x132153c7.
//
// Solidity: function manuallySetRewardsGenesisTime(uint256 timestamp) returns()
func (_Reward *RewardTransactorSession) ManuallySetRewardsGenesisTime(timestamp *big.Int) (*types.Transaction, error) {
	return _Reward.Contract.ManuallySetRewardsGenesisTime(&_Reward.TransactOpts, timestamp)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_Reward *RewardTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Reward.contract.Transact(opts, "renounceRole", role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_Reward *RewardSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Reward.Contract.RenounceRole(&_Reward.TransactOpts, role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_Reward *RewardTransactorSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Reward.Contract.RenounceRole(&_Reward.TransactOpts, role, account)
}

// ResetRewardsGenesisTime is a paid mutator transaction binding the contract method 0x1caa945d.
//
// Solidity: function resetRewardsGenesisTime() returns()
func (_Reward *RewardTransactor) ResetRewardsGenesisTime(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Reward.contract.Transact(opts, "resetRewardsGenesisTime")
}

// ResetRewardsGenesisTime is a paid mutator transaction binding the contract method 0x1caa945d.
//
// Solidity: function resetRewardsGenesisTime() returns()
func (_Reward *RewardSession) ResetRewardsGenesisTime() (*types.Transaction, error) {
	return _Reward.Contract.ResetRewardsGenesisTime(&_Reward.TransactOpts)
}

// ResetRewardsGenesisTime is a paid mutator transaction binding the contract method 0x1caa945d.
//
// Solidity: function resetRewardsGenesisTime() returns()
func (_Reward *RewardTransactorSession) ResetRewardsGenesisTime() (*types.Transaction, error) {
	return _Reward.Contract.ResetRewardsGenesisTime(&_Reward.TransactOpts)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_Reward *RewardTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Reward.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_Reward *RewardSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Reward.Contract.RevokeRole(&_Reward.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_Reward *RewardTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Reward.Contract.RevokeRole(&_Reward.TransactOpts, role, account)
}

// SetAftermarketProxyAddress is a paid mutator transaction binding the contract method 0x56ba7adb.
//
// Solidity: function setAftermarketProxyAddress(address aftermarketProxyAddress) returns()
func (_Reward *RewardTransactor) SetAftermarketProxyAddress(opts *bind.TransactOpts, aftermarketProxyAddress common.Address) (*types.Transaction, error) {
	return _Reward.contract.Transact(opts, "setAftermarketProxyAddress", aftermarketProxyAddress)
}

// SetAftermarketProxyAddress is a paid mutator transaction binding the contract method 0x56ba7adb.
//
// Solidity: function setAftermarketProxyAddress(address aftermarketProxyAddress) returns()
func (_Reward *RewardSession) SetAftermarketProxyAddress(aftermarketProxyAddress common.Address) (*types.Transaction, error) {
	return _Reward.Contract.SetAftermarketProxyAddress(&_Reward.TransactOpts, aftermarketProxyAddress)
}

// SetAftermarketProxyAddress is a paid mutator transaction binding the contract method 0x56ba7adb.
//
// Solidity: function setAftermarketProxyAddress(address aftermarketProxyAddress) returns()
func (_Reward *RewardTransactorSession) SetAftermarketProxyAddress(aftermarketProxyAddress common.Address) (*types.Transaction, error) {
	return _Reward.Contract.SetAftermarketProxyAddress(&_Reward.TransactOpts, aftermarketProxyAddress)
}

// SetMinimumTimeForRewards is a paid mutator transaction binding the contract method 0x066404a7.
//
// Solidity: function setMinimumTimeForRewards(uint256 newTimeInSeconds) returns()
func (_Reward *RewardTransactor) SetMinimumTimeForRewards(opts *bind.TransactOpts, newTimeInSeconds *big.Int) (*types.Transaction, error) {
	return _Reward.contract.Transact(opts, "setMinimumTimeForRewards", newTimeInSeconds)
}

// SetMinimumTimeForRewards is a paid mutator transaction binding the contract method 0x066404a7.
//
// Solidity: function setMinimumTimeForRewards(uint256 newTimeInSeconds) returns()
func (_Reward *RewardSession) SetMinimumTimeForRewards(newTimeInSeconds *big.Int) (*types.Transaction, error) {
	return _Reward.Contract.SetMinimumTimeForRewards(&_Reward.TransactOpts, newTimeInSeconds)
}

// SetMinimumTimeForRewards is a paid mutator transaction binding the contract method 0x066404a7.
//
// Solidity: function setMinimumTimeForRewards(uint256 newTimeInSeconds) returns()
func (_Reward *RewardTransactorSession) SetMinimumTimeForRewards(newTimeInSeconds *big.Int) (*types.Transaction, error) {
	return _Reward.Contract.SetMinimumTimeForRewards(&_Reward.TransactOpts, newTimeInSeconds)
}

// SetRegistryContractAddress is a paid mutator transaction binding the contract method 0xd20838f4.
//
// Solidity: function setRegistryContractAddress(address registryContractAddress) returns()
func (_Reward *RewardTransactor) SetRegistryContractAddress(opts *bind.TransactOpts, registryContractAddress common.Address) (*types.Transaction, error) {
	return _Reward.contract.Transact(opts, "setRegistryContractAddress", registryContractAddress)
}

// SetRegistryContractAddress is a paid mutator transaction binding the contract method 0xd20838f4.
//
// Solidity: function setRegistryContractAddress(address registryContractAddress) returns()
func (_Reward *RewardSession) SetRegistryContractAddress(registryContractAddress common.Address) (*types.Transaction, error) {
	return _Reward.Contract.SetRegistryContractAddress(&_Reward.TransactOpts, registryContractAddress)
}

// SetRegistryContractAddress is a paid mutator transaction binding the contract method 0xd20838f4.
//
// Solidity: function setRegistryContractAddress(address registryContractAddress) returns()
func (_Reward *RewardTransactorSession) SetRegistryContractAddress(registryContractAddress common.Address) (*types.Transaction, error) {
	return _Reward.Contract.SetRegistryContractAddress(&_Reward.TransactOpts, registryContractAddress)
}

// SetSyntheticProxyAddress is a paid mutator transaction binding the contract method 0xee670696.
//
// Solidity: function setSyntheticProxyAddress(address syntheticProxyAddress) returns()
func (_Reward *RewardTransactor) SetSyntheticProxyAddress(opts *bind.TransactOpts, syntheticProxyAddress common.Address) (*types.Transaction, error) {
	return _Reward.contract.Transact(opts, "setSyntheticProxyAddress", syntheticProxyAddress)
}

// SetSyntheticProxyAddress is a paid mutator transaction binding the contract method 0xee670696.
//
// Solidity: function setSyntheticProxyAddress(address syntheticProxyAddress) returns()
func (_Reward *RewardSession) SetSyntheticProxyAddress(syntheticProxyAddress common.Address) (*types.Transaction, error) {
	return _Reward.Contract.SetSyntheticProxyAddress(&_Reward.TransactOpts, syntheticProxyAddress)
}

// SetSyntheticProxyAddress is a paid mutator transaction binding the contract method 0xee670696.
//
// Solidity: function setSyntheticProxyAddress(address syntheticProxyAddress) returns()
func (_Reward *RewardTransactorSession) SetSyntheticProxyAddress(syntheticProxyAddress common.Address) (*types.Transaction, error) {
	return _Reward.Contract.SetSyntheticProxyAddress(&_Reward.TransactOpts, syntheticProxyAddress)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_Reward *RewardTransactor) UpgradeTo(opts *bind.TransactOpts, newImplementation common.Address) (*types.Transaction, error) {
	return _Reward.contract.Transact(opts, "upgradeTo", newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_Reward *RewardSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _Reward.Contract.UpgradeTo(&_Reward.TransactOpts, newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_Reward *RewardTransactorSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _Reward.Contract.UpgradeTo(&_Reward.TransactOpts, newImplementation)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Reward *RewardTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Reward.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Reward *RewardSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Reward.Contract.UpgradeToAndCall(&_Reward.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Reward *RewardTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Reward.Contract.UpgradeToAndCall(&_Reward.TransactOpts, newImplementation, data)
}

// RewardAdminChangedIterator is returned from FilterAdminChanged and is used to iterate over the raw logs and unpacked data for AdminChanged events raised by the Reward contract.
type RewardAdminChangedIterator struct {
	Event *RewardAdminChanged // Event containing the contract specifics and raw log

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
func (it *RewardAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RewardAdminChanged)
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
		it.Event = new(RewardAdminChanged)
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
func (it *RewardAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RewardAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RewardAdminChanged represents a AdminChanged event raised by the Reward contract.
type RewardAdminChanged struct {
	PreviousAdmin common.Address
	NewAdmin      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterAdminChanged is a free log retrieval operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_Reward *RewardFilterer) FilterAdminChanged(opts *bind.FilterOpts) (*RewardAdminChangedIterator, error) {

	logs, sub, err := _Reward.contract.FilterLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return &RewardAdminChangedIterator{contract: _Reward.contract, event: "AdminChanged", logs: logs, sub: sub}, nil
}

// WatchAdminChanged is a free log subscription operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_Reward *RewardFilterer) WatchAdminChanged(opts *bind.WatchOpts, sink chan<- *RewardAdminChanged) (event.Subscription, error) {

	logs, sub, err := _Reward.contract.WatchLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RewardAdminChanged)
				if err := _Reward.contract.UnpackLog(event, "AdminChanged", log); err != nil {
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
func (_Reward *RewardFilterer) ParseAdminChanged(log types.Log) (*RewardAdminChanged, error) {
	event := new(RewardAdminChanged)
	if err := _Reward.contract.UnpackLog(event, "AdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RewardAdminWithdrawalIterator is returned from FilterAdminWithdrawal and is used to iterate over the raw logs and unpacked data for AdminWithdrawal events raised by the Reward contract.
type RewardAdminWithdrawalIterator struct {
	Event *RewardAdminWithdrawal // Event containing the contract specifics and raw log

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
func (it *RewardAdminWithdrawalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RewardAdminWithdrawal)
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
		it.Event = new(RewardAdminWithdrawal)
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
func (it *RewardAdminWithdrawalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RewardAdminWithdrawalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RewardAdminWithdrawal represents a AdminWithdrawal event raised by the Reward contract.
type RewardAdminWithdrawal struct {
	User   common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterAdminWithdrawal is a free log retrieval operation binding the contract event 0x1f29bc8239df330207e019f41493b485f9c7d3ce83a795ae64603dde527ada2e.
//
// Solidity: event AdminWithdrawal(address indexed user, uint256 _amount)
func (_Reward *RewardFilterer) FilterAdminWithdrawal(opts *bind.FilterOpts, user []common.Address) (*RewardAdminWithdrawalIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _Reward.contract.FilterLogs(opts, "AdminWithdrawal", userRule)
	if err != nil {
		return nil, err
	}
	return &RewardAdminWithdrawalIterator{contract: _Reward.contract, event: "AdminWithdrawal", logs: logs, sub: sub}, nil
}

// WatchAdminWithdrawal is a free log subscription operation binding the contract event 0x1f29bc8239df330207e019f41493b485f9c7d3ce83a795ae64603dde527ada2e.
//
// Solidity: event AdminWithdrawal(address indexed user, uint256 _amount)
func (_Reward *RewardFilterer) WatchAdminWithdrawal(opts *bind.WatchOpts, sink chan<- *RewardAdminWithdrawal, user []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _Reward.contract.WatchLogs(opts, "AdminWithdrawal", userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RewardAdminWithdrawal)
				if err := _Reward.contract.UnpackLog(event, "AdminWithdrawal", log); err != nil {
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
// Solidity: event AdminWithdrawal(address indexed user, uint256 _amount)
func (_Reward *RewardFilterer) ParseAdminWithdrawal(log types.Log) (*RewardAdminWithdrawal, error) {
	event := new(RewardAdminWithdrawal)
	if err := _Reward.contract.UnpackLog(event, "AdminWithdrawal", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RewardBeaconUpgradedIterator is returned from FilterBeaconUpgraded and is used to iterate over the raw logs and unpacked data for BeaconUpgraded events raised by the Reward contract.
type RewardBeaconUpgradedIterator struct {
	Event *RewardBeaconUpgraded // Event containing the contract specifics and raw log

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
func (it *RewardBeaconUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RewardBeaconUpgraded)
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
		it.Event = new(RewardBeaconUpgraded)
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
func (it *RewardBeaconUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RewardBeaconUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RewardBeaconUpgraded represents a BeaconUpgraded event raised by the Reward contract.
type RewardBeaconUpgraded struct {
	Beacon common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBeaconUpgraded is a free log retrieval operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_Reward *RewardFilterer) FilterBeaconUpgraded(opts *bind.FilterOpts, beacon []common.Address) (*RewardBeaconUpgradedIterator, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _Reward.contract.FilterLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return &RewardBeaconUpgradedIterator{contract: _Reward.contract, event: "BeaconUpgraded", logs: logs, sub: sub}, nil
}

// WatchBeaconUpgraded is a free log subscription operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_Reward *RewardFilterer) WatchBeaconUpgraded(opts *bind.WatchOpts, sink chan<- *RewardBeaconUpgraded, beacon []common.Address) (event.Subscription, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _Reward.contract.WatchLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RewardBeaconUpgraded)
				if err := _Reward.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
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
func (_Reward *RewardFilterer) ParseBeaconUpgraded(log types.Log) (*RewardBeaconUpgraded, error) {
	event := new(RewardBeaconUpgraded)
	if err := _Reward.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RewardDidntQualifyIterator is returned from FilterDidntQualify and is used to iterate over the raw logs and unpacked data for DidntQualify events raised by the Reward contract.
type RewardDidntQualifyIterator struct {
	Event *RewardDidntQualify // Event containing the contract specifics and raw log

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
func (it *RewardDidntQualifyIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RewardDidntQualify)
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
		it.Event = new(RewardDidntQualify)
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
func (it *RewardDidntQualifyIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RewardDidntQualifyIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RewardDidntQualify represents a DidntQualify event raised by the Reward contract.
type RewardDidntQualify struct {
	User          common.Address
	Amount        *big.Int
	VehicleNodeId *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterDidntQualify is a free log retrieval operation binding the contract event 0xf667cbbd740351d2f58b68dc6a8ee1b1f3a853c609d7ce7eee4a3bdc94d62774.
//
// Solidity: event DidntQualify(address indexed user, uint256 _amount, uint256 indexed vehicleNodeId)
func (_Reward *RewardFilterer) FilterDidntQualify(opts *bind.FilterOpts, user []common.Address, vehicleNodeId []*big.Int) (*RewardDidntQualifyIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	var vehicleNodeIdRule []interface{}
	for _, vehicleNodeIdItem := range vehicleNodeId {
		vehicleNodeIdRule = append(vehicleNodeIdRule, vehicleNodeIdItem)
	}

	logs, sub, err := _Reward.contract.FilterLogs(opts, "DidntQualify", userRule, vehicleNodeIdRule)
	if err != nil {
		return nil, err
	}
	return &RewardDidntQualifyIterator{contract: _Reward.contract, event: "DidntQualify", logs: logs, sub: sub}, nil
}

// WatchDidntQualify is a free log subscription operation binding the contract event 0xf667cbbd740351d2f58b68dc6a8ee1b1f3a853c609d7ce7eee4a3bdc94d62774.
//
// Solidity: event DidntQualify(address indexed user, uint256 _amount, uint256 indexed vehicleNodeId)
func (_Reward *RewardFilterer) WatchDidntQualify(opts *bind.WatchOpts, sink chan<- *RewardDidntQualify, user []common.Address, vehicleNodeId []*big.Int) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	var vehicleNodeIdRule []interface{}
	for _, vehicleNodeIdItem := range vehicleNodeId {
		vehicleNodeIdRule = append(vehicleNodeIdRule, vehicleNodeIdItem)
	}

	logs, sub, err := _Reward.contract.WatchLogs(opts, "DidntQualify", userRule, vehicleNodeIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RewardDidntQualify)
				if err := _Reward.contract.UnpackLog(event, "DidntQualify", log); err != nil {
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

// ParseDidntQualify is a log parse operation binding the contract event 0xf667cbbd740351d2f58b68dc6a8ee1b1f3a853c609d7ce7eee4a3bdc94d62774.
//
// Solidity: event DidntQualify(address indexed user, uint256 _amount, uint256 indexed vehicleNodeId)
func (_Reward *RewardFilterer) ParseDidntQualify(log types.Log) (*RewardDidntQualify, error) {
	event := new(RewardDidntQualify)
	if err := _Reward.contract.UnpackLog(event, "DidntQualify", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RewardInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the Reward contract.
type RewardInitializedIterator struct {
	Event *RewardInitialized // Event containing the contract specifics and raw log

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
func (it *RewardInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RewardInitialized)
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
		it.Event = new(RewardInitialized)
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
func (it *RewardInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RewardInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RewardInitialized represents a Initialized event raised by the Reward contract.
type RewardInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Reward *RewardFilterer) FilterInitialized(opts *bind.FilterOpts) (*RewardInitializedIterator, error) {

	logs, sub, err := _Reward.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &RewardInitializedIterator{contract: _Reward.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Reward *RewardFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *RewardInitialized) (event.Subscription, error) {

	logs, sub, err := _Reward.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RewardInitialized)
				if err := _Reward.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_Reward *RewardFilterer) ParseInitialized(log types.Log) (*RewardInitialized, error) {
	event := new(RewardInitialized)
	if err := _Reward.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RewardRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the Reward contract.
type RewardRoleAdminChangedIterator struct {
	Event *RewardRoleAdminChanged // Event containing the contract specifics and raw log

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
func (it *RewardRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RewardRoleAdminChanged)
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
		it.Event = new(RewardRoleAdminChanged)
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
func (it *RewardRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RewardRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RewardRoleAdminChanged represents a RoleAdminChanged event raised by the Reward contract.
type RewardRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_Reward *RewardFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*RewardRoleAdminChangedIterator, error) {

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

	logs, sub, err := _Reward.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &RewardRoleAdminChangedIterator{contract: _Reward.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_Reward *RewardFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *RewardRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

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

	logs, sub, err := _Reward.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RewardRoleAdminChanged)
				if err := _Reward.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
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
func (_Reward *RewardFilterer) ParseRoleAdminChanged(log types.Log) (*RewardRoleAdminChanged, error) {
	event := new(RewardRoleAdminChanged)
	if err := _Reward.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RewardRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the Reward contract.
type RewardRoleGrantedIterator struct {
	Event *RewardRoleGranted // Event containing the contract specifics and raw log

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
func (it *RewardRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RewardRoleGranted)
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
		it.Event = new(RewardRoleGranted)
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
func (it *RewardRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RewardRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RewardRoleGranted represents a RoleGranted event raised by the Reward contract.
type RewardRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_Reward *RewardFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*RewardRoleGrantedIterator, error) {

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

	logs, sub, err := _Reward.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &RewardRoleGrantedIterator{contract: _Reward.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_Reward *RewardFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *RewardRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _Reward.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RewardRoleGranted)
				if err := _Reward.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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
func (_Reward *RewardFilterer) ParseRoleGranted(log types.Log) (*RewardRoleGranted, error) {
	event := new(RewardRoleGranted)
	if err := _Reward.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RewardRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the Reward contract.
type RewardRoleRevokedIterator struct {
	Event *RewardRoleRevoked // Event containing the contract specifics and raw log

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
func (it *RewardRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RewardRoleRevoked)
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
		it.Event = new(RewardRoleRevoked)
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
func (it *RewardRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RewardRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RewardRoleRevoked represents a RoleRevoked event raised by the Reward contract.
type RewardRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_Reward *RewardFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*RewardRoleRevokedIterator, error) {

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

	logs, sub, err := _Reward.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &RewardRoleRevokedIterator{contract: _Reward.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_Reward *RewardFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *RewardRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _Reward.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RewardRoleRevoked)
				if err := _Reward.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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
func (_Reward *RewardFilterer) ParseRoleRevoked(log types.Log) (*RewardRoleRevoked, error) {
	event := new(RewardRoleRevoked)
	if err := _Reward.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RewardTokensTransferredIterator is returned from FilterTokensTransferred and is used to iterate over the raw logs and unpacked data for TokensTransferred events raised by the Reward contract.
type RewardTokensTransferredIterator struct {
	Event *RewardTokensTransferred // Event containing the contract specifics and raw log

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
func (it *RewardTokensTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RewardTokensTransferred)
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
		it.Event = new(RewardTokensTransferred)
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
func (it *RewardTokensTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RewardTokensTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RewardTokensTransferred represents a TokensTransferred event raised by the Reward contract.
type RewardTokensTransferred struct {
	User          common.Address
	Amount        *big.Int
	VehicleNodeId *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterTokensTransferred is a free log retrieval operation binding the contract event 0x57e1000ba5ba7b6ab6670639de9fc3db34d05ef2bbce4a09d60dda560387b0ea.
//
// Solidity: event TokensTransferred(address indexed user, uint256 _amount, uint256 indexed vehicleNodeId)
func (_Reward *RewardFilterer) FilterTokensTransferred(opts *bind.FilterOpts, user []common.Address, vehicleNodeId []*big.Int) (*RewardTokensTransferredIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	var vehicleNodeIdRule []interface{}
	for _, vehicleNodeIdItem := range vehicleNodeId {
		vehicleNodeIdRule = append(vehicleNodeIdRule, vehicleNodeIdItem)
	}

	logs, sub, err := _Reward.contract.FilterLogs(opts, "TokensTransferred", userRule, vehicleNodeIdRule)
	if err != nil {
		return nil, err
	}
	return &RewardTokensTransferredIterator{contract: _Reward.contract, event: "TokensTransferred", logs: logs, sub: sub}, nil
}

// WatchTokensTransferred is a free log subscription operation binding the contract event 0x57e1000ba5ba7b6ab6670639de9fc3db34d05ef2bbce4a09d60dda560387b0ea.
//
// Solidity: event TokensTransferred(address indexed user, uint256 _amount, uint256 indexed vehicleNodeId)
func (_Reward *RewardFilterer) WatchTokensTransferred(opts *bind.WatchOpts, sink chan<- *RewardTokensTransferred, user []common.Address, vehicleNodeId []*big.Int) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	var vehicleNodeIdRule []interface{}
	for _, vehicleNodeIdItem := range vehicleNodeId {
		vehicleNodeIdRule = append(vehicleNodeIdRule, vehicleNodeIdItem)
	}

	logs, sub, err := _Reward.contract.WatchLogs(opts, "TokensTransferred", userRule, vehicleNodeIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RewardTokensTransferred)
				if err := _Reward.contract.UnpackLog(event, "TokensTransferred", log); err != nil {
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

// ParseTokensTransferred is a log parse operation binding the contract event 0x57e1000ba5ba7b6ab6670639de9fc3db34d05ef2bbce4a09d60dda560387b0ea.
//
// Solidity: event TokensTransferred(address indexed user, uint256 _amount, uint256 indexed vehicleNodeId)
func (_Reward *RewardFilterer) ParseTokensTransferred(log types.Log) (*RewardTokensTransferred, error) {
	event := new(RewardTokensTransferred)
	if err := _Reward.contract.UnpackLog(event, "TokensTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RewardTokensTransferredForConnectionStreakIterator is returned from FilterTokensTransferredForConnectionStreak and is used to iterate over the raw logs and unpacked data for TokensTransferredForConnectionStreak events raised by the Reward contract.
type RewardTokensTransferredForConnectionStreakIterator struct {
	Event *RewardTokensTransferredForConnectionStreak // Event containing the contract specifics and raw log

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
func (it *RewardTokensTransferredForConnectionStreakIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RewardTokensTransferredForConnectionStreak)
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
		it.Event = new(RewardTokensTransferredForConnectionStreak)
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
func (it *RewardTokensTransferredForConnectionStreakIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RewardTokensTransferredForConnectionStreakIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RewardTokensTransferredForConnectionStreak represents a TokensTransferredForConnectionStreak event raised by the Reward contract.
type RewardTokensTransferredForConnectionStreak struct {
	Week             *big.Int
	User             common.Address
	Amount           *big.Int
	VehicleNodeId    *big.Int
	ConnectionStreak *big.Int
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterTokensTransferredForConnectionStreak is a free log retrieval operation binding the contract event 0x4608830ab3ec9e091c65f293254fe1aef6fcd37ed152914a665e6ad38de6a992.
//
// Solidity: event TokensTransferredForConnectionStreak(uint256 week, address indexed user, uint256 _amount, uint256 indexed vehicleNodeId, uint256 connectionStreak)
func (_Reward *RewardFilterer) FilterTokensTransferredForConnectionStreak(opts *bind.FilterOpts, user []common.Address, vehicleNodeId []*big.Int) (*RewardTokensTransferredForConnectionStreakIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	var vehicleNodeIdRule []interface{}
	for _, vehicleNodeIdItem := range vehicleNodeId {
		vehicleNodeIdRule = append(vehicleNodeIdRule, vehicleNodeIdItem)
	}

	logs, sub, err := _Reward.contract.FilterLogs(opts, "TokensTransferredForConnectionStreak", userRule, vehicleNodeIdRule)
	if err != nil {
		return nil, err
	}
	return &RewardTokensTransferredForConnectionStreakIterator{contract: _Reward.contract, event: "TokensTransferredForConnectionStreak", logs: logs, sub: sub}, nil
}

// WatchTokensTransferredForConnectionStreak is a free log subscription operation binding the contract event 0x4608830ab3ec9e091c65f293254fe1aef6fcd37ed152914a665e6ad38de6a992.
//
// Solidity: event TokensTransferredForConnectionStreak(uint256 week, address indexed user, uint256 _amount, uint256 indexed vehicleNodeId, uint256 connectionStreak)
func (_Reward *RewardFilterer) WatchTokensTransferredForConnectionStreak(opts *bind.WatchOpts, sink chan<- *RewardTokensTransferredForConnectionStreak, user []common.Address, vehicleNodeId []*big.Int) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	var vehicleNodeIdRule []interface{}
	for _, vehicleNodeIdItem := range vehicleNodeId {
		vehicleNodeIdRule = append(vehicleNodeIdRule, vehicleNodeIdItem)
	}

	logs, sub, err := _Reward.contract.WatchLogs(opts, "TokensTransferredForConnectionStreak", userRule, vehicleNodeIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RewardTokensTransferredForConnectionStreak)
				if err := _Reward.contract.UnpackLog(event, "TokensTransferredForConnectionStreak", log); err != nil {
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

// ParseTokensTransferredForConnectionStreak is a log parse operation binding the contract event 0x4608830ab3ec9e091c65f293254fe1aef6fcd37ed152914a665e6ad38de6a992.
//
// Solidity: event TokensTransferredForConnectionStreak(uint256 week, address indexed user, uint256 _amount, uint256 indexed vehicleNodeId, uint256 connectionStreak)
func (_Reward *RewardFilterer) ParseTokensTransferredForConnectionStreak(log types.Log) (*RewardTokensTransferredForConnectionStreak, error) {
	event := new(RewardTokensTransferredForConnectionStreak)
	if err := _Reward.contract.UnpackLog(event, "TokensTransferredForConnectionStreak", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RewardTokensTransferredForDeviceIterator is returned from FilterTokensTransferredForDevice and is used to iterate over the raw logs and unpacked data for TokensTransferredForDevice events raised by the Reward contract.
type RewardTokensTransferredForDeviceIterator struct {
	Event *RewardTokensTransferredForDevice // Event containing the contract specifics and raw log

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
func (it *RewardTokensTransferredForDeviceIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RewardTokensTransferredForDevice)
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
		it.Event = new(RewardTokensTransferredForDevice)
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
func (it *RewardTokensTransferredForDeviceIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RewardTokensTransferredForDeviceIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RewardTokensTransferredForDevice represents a TokensTransferredForDevice event raised by the Reward contract.
type RewardTokensTransferredForDevice struct {
	Week           *big.Int
	User           common.Address
	Amount         *big.Int
	VehicleNodeId  *big.Int
	DeviceNftProxy common.Address
	DeviceNode     *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterTokensTransferredForDevice is a free log retrieval operation binding the contract event 0x3ab5458e0b305f1e637781c67c5be0f43b7ee1a4427584e3b229524929184bc4.
//
// Solidity: event TokensTransferredForDevice(uint256 week, address indexed user, uint256 _amount, uint256 indexed vehicleNodeId, address deviceNftProxy, uint256 indexed deviceNode)
func (_Reward *RewardFilterer) FilterTokensTransferredForDevice(opts *bind.FilterOpts, user []common.Address, vehicleNodeId []*big.Int, deviceNode []*big.Int) (*RewardTokensTransferredForDeviceIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	var vehicleNodeIdRule []interface{}
	for _, vehicleNodeIdItem := range vehicleNodeId {
		vehicleNodeIdRule = append(vehicleNodeIdRule, vehicleNodeIdItem)
	}

	var deviceNodeRule []interface{}
	for _, deviceNodeItem := range deviceNode {
		deviceNodeRule = append(deviceNodeRule, deviceNodeItem)
	}

	logs, sub, err := _Reward.contract.FilterLogs(opts, "TokensTransferredForDevice", userRule, vehicleNodeIdRule, deviceNodeRule)
	if err != nil {
		return nil, err
	}
	return &RewardTokensTransferredForDeviceIterator{contract: _Reward.contract, event: "TokensTransferredForDevice", logs: logs, sub: sub}, nil
}

// WatchTokensTransferredForDevice is a free log subscription operation binding the contract event 0x3ab5458e0b305f1e637781c67c5be0f43b7ee1a4427584e3b229524929184bc4.
//
// Solidity: event TokensTransferredForDevice(uint256 week, address indexed user, uint256 _amount, uint256 indexed vehicleNodeId, address deviceNftProxy, uint256 indexed deviceNode)
func (_Reward *RewardFilterer) WatchTokensTransferredForDevice(opts *bind.WatchOpts, sink chan<- *RewardTokensTransferredForDevice, user []common.Address, vehicleNodeId []*big.Int, deviceNode []*big.Int) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	var vehicleNodeIdRule []interface{}
	for _, vehicleNodeIdItem := range vehicleNodeId {
		vehicleNodeIdRule = append(vehicleNodeIdRule, vehicleNodeIdItem)
	}

	var deviceNodeRule []interface{}
	for _, deviceNodeItem := range deviceNode {
		deviceNodeRule = append(deviceNodeRule, deviceNodeItem)
	}

	logs, sub, err := _Reward.contract.WatchLogs(opts, "TokensTransferredForDevice", userRule, vehicleNodeIdRule, deviceNodeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RewardTokensTransferredForDevice)
				if err := _Reward.contract.UnpackLog(event, "TokensTransferredForDevice", log); err != nil {
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

// ParseTokensTransferredForDevice is a log parse operation binding the contract event 0x3ab5458e0b305f1e637781c67c5be0f43b7ee1a4427584e3b229524929184bc4.
//
// Solidity: event TokensTransferredForDevice(uint256 week, address indexed user, uint256 _amount, uint256 indexed vehicleNodeId, address deviceNftProxy, uint256 indexed deviceNode)
func (_Reward *RewardFilterer) ParseTokensTransferredForDevice(log types.Log) (*RewardTokensTransferredForDevice, error) {
	event := new(RewardTokensTransferredForDevice)
	if err := _Reward.contract.UnpackLog(event, "TokensTransferredForDevice", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RewardUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the Reward contract.
type RewardUpgradedIterator struct {
	Event *RewardUpgraded // Event containing the contract specifics and raw log

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
func (it *RewardUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RewardUpgraded)
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
		it.Event = new(RewardUpgraded)
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
func (it *RewardUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RewardUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RewardUpgraded represents a Upgraded event raised by the Reward contract.
type RewardUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_Reward *RewardFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*RewardUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _Reward.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &RewardUpgradedIterator{contract: _Reward.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_Reward *RewardFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *RewardUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _Reward.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RewardUpgraded)
				if err := _Reward.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_Reward *RewardFilterer) ParseUpgraded(log types.Log) (*RewardUpgraded, error) {
	event := new(RewardUpgraded)
	if err := _Reward.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RewardWeekAdvancedIterator is returned from FilterWeekAdvanced and is used to iterate over the raw logs and unpacked data for WeekAdvanced events raised by the Reward contract.
type RewardWeekAdvancedIterator struct {
	Event *RewardWeekAdvanced // Event containing the contract specifics and raw log

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
func (it *RewardWeekAdvancedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RewardWeekAdvanced)
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
		it.Event = new(RewardWeekAdvanced)
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
func (it *RewardWeekAdvancedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RewardWeekAdvancedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RewardWeekAdvanced represents a WeekAdvanced event raised by the Reward contract.
type RewardWeekAdvanced struct {
	NewWeek *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterWeekAdvanced is a free log retrieval operation binding the contract event 0x54ab1f5ea80e88ee480c10edc5b60056a0e392c6125954a9f47ce8677acc59d0.
//
// Solidity: event WeekAdvanced(uint256 newWeek)
func (_Reward *RewardFilterer) FilterWeekAdvanced(opts *bind.FilterOpts) (*RewardWeekAdvancedIterator, error) {

	logs, sub, err := _Reward.contract.FilterLogs(opts, "WeekAdvanced")
	if err != nil {
		return nil, err
	}
	return &RewardWeekAdvancedIterator{contract: _Reward.contract, event: "WeekAdvanced", logs: logs, sub: sub}, nil
}

// WatchWeekAdvanced is a free log subscription operation binding the contract event 0x54ab1f5ea80e88ee480c10edc5b60056a0e392c6125954a9f47ce8677acc59d0.
//
// Solidity: event WeekAdvanced(uint256 newWeek)
func (_Reward *RewardFilterer) WatchWeekAdvanced(opts *bind.WatchOpts, sink chan<- *RewardWeekAdvanced) (event.Subscription, error) {

	logs, sub, err := _Reward.contract.WatchLogs(opts, "WeekAdvanced")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RewardWeekAdvanced)
				if err := _Reward.contract.UnpackLog(event, "WeekAdvanced", log); err != nil {
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

// ParseWeekAdvanced is a log parse operation binding the contract event 0x54ab1f5ea80e88ee480c10edc5b60056a0e392c6125954a9f47ce8677acc59d0.
//
// Solidity: event WeekAdvanced(uint256 newWeek)
func (_Reward *RewardFilterer) ParseWeekAdvanced(log types.Log) (*RewardWeekAdvanced, error) {
	event := new(RewardWeekAdvanced)
	if err := _Reward.contract.UnpackLog(event, "WeekAdvanced", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
