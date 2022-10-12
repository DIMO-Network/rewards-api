// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package issuance

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

// IssuanceMetaData contains all meta data concerning the Issuance contract.
var IssuanceMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"vehicleNodeId\",\"type\":\"uint256\"}],\"name\":\"TokensTransferred\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"users\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"values\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"vehicleIds\",\"type\":\"uint256[]\"}],\"name\":\"batchTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// IssuanceABI is the input ABI used to generate the binding from.
// Deprecated: Use IssuanceMetaData.ABI instead.
var IssuanceABI = IssuanceMetaData.ABI

// Issuance is an auto generated Go binding around an Ethereum contract.
type Issuance struct {
	IssuanceCaller     // Read-only binding to the contract
	IssuanceTransactor // Write-only binding to the contract
	IssuanceFilterer   // Log filterer for contract events
}

// IssuanceCaller is an auto generated read-only Go binding around an Ethereum contract.
type IssuanceCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IssuanceTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IssuanceTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IssuanceFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IssuanceFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IssuanceSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IssuanceSession struct {
	Contract     *Issuance         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IssuanceCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IssuanceCallerSession struct {
	Contract *IssuanceCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// IssuanceTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IssuanceTransactorSession struct {
	Contract     *IssuanceTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// IssuanceRaw is an auto generated low-level Go binding around an Ethereum contract.
type IssuanceRaw struct {
	Contract *Issuance // Generic contract binding to access the raw methods on
}

// IssuanceCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IssuanceCallerRaw struct {
	Contract *IssuanceCaller // Generic read-only contract binding to access the raw methods on
}

// IssuanceTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IssuanceTransactorRaw struct {
	Contract *IssuanceTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIssuance creates a new instance of Issuance, bound to a specific deployed contract.
func NewIssuance(address common.Address, backend bind.ContractBackend) (*Issuance, error) {
	contract, err := bindIssuance(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Issuance{IssuanceCaller: IssuanceCaller{contract: contract}, IssuanceTransactor: IssuanceTransactor{contract: contract}, IssuanceFilterer: IssuanceFilterer{contract: contract}}, nil
}

// NewIssuanceCaller creates a new read-only instance of Issuance, bound to a specific deployed contract.
func NewIssuanceCaller(address common.Address, caller bind.ContractCaller) (*IssuanceCaller, error) {
	contract, err := bindIssuance(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IssuanceCaller{contract: contract}, nil
}

// NewIssuanceTransactor creates a new write-only instance of Issuance, bound to a specific deployed contract.
func NewIssuanceTransactor(address common.Address, transactor bind.ContractTransactor) (*IssuanceTransactor, error) {
	contract, err := bindIssuance(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IssuanceTransactor{contract: contract}, nil
}

// NewIssuanceFilterer creates a new log filterer instance of Issuance, bound to a specific deployed contract.
func NewIssuanceFilterer(address common.Address, filterer bind.ContractFilterer) (*IssuanceFilterer, error) {
	contract, err := bindIssuance(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IssuanceFilterer{contract: contract}, nil
}

// bindIssuance binds a generic wrapper to an already deployed contract.
func bindIssuance(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(IssuanceABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Issuance *IssuanceRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Issuance.Contract.IssuanceCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Issuance *IssuanceRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Issuance.Contract.IssuanceTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Issuance *IssuanceRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Issuance.Contract.IssuanceTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Issuance *IssuanceCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Issuance.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Issuance *IssuanceTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Issuance.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Issuance *IssuanceTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Issuance.Contract.contract.Transact(opts, method, params...)
}

// BatchTransfer is a paid mutator transaction binding the contract method 0x1ef690c4.
//
// Solidity: function batchTransfer(address[] users, uint256[] values, uint256[] vehicleIds) returns()
func (_Issuance *IssuanceTransactor) BatchTransfer(opts *bind.TransactOpts, users []common.Address, values []*big.Int, vehicleIds []*big.Int) (*types.Transaction, error) {
	return _Issuance.contract.Transact(opts, "batchTransfer", users, values, vehicleIds)
}

// BatchTransfer is a paid mutator transaction binding the contract method 0x1ef690c4.
//
// Solidity: function batchTransfer(address[] users, uint256[] values, uint256[] vehicleIds) returns()
func (_Issuance *IssuanceSession) BatchTransfer(users []common.Address, values []*big.Int, vehicleIds []*big.Int) (*types.Transaction, error) {
	return _Issuance.Contract.BatchTransfer(&_Issuance.TransactOpts, users, values, vehicleIds)
}

// BatchTransfer is a paid mutator transaction binding the contract method 0x1ef690c4.
//
// Solidity: function batchTransfer(address[] users, uint256[] values, uint256[] vehicleIds) returns()
func (_Issuance *IssuanceTransactorSession) BatchTransfer(users []common.Address, values []*big.Int, vehicleIds []*big.Int) (*types.Transaction, error) {
	return _Issuance.Contract.BatchTransfer(&_Issuance.TransactOpts, users, values, vehicleIds)
}

// IssuanceTokensTransferredIterator is returned from FilterTokensTransferred and is used to iterate over the raw logs and unpacked data for TokensTransferred events raised by the Issuance contract.
type IssuanceTokensTransferredIterator struct {
	Event *IssuanceTokensTransferred // Event containing the contract specifics and raw log

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
func (it *IssuanceTokensTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IssuanceTokensTransferred)
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
		it.Event = new(IssuanceTokensTransferred)
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
func (it *IssuanceTokensTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IssuanceTokensTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IssuanceTokensTransferred represents a TokensTransferred event raised by the Issuance contract.
type IssuanceTokensTransferred struct {
	User          common.Address
	Amount        *big.Int
	VehicleNodeId *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterTokensTransferred is a free log retrieval operation binding the contract event 0x57e1000ba5ba7b6ab6670639de9fc3db34d05ef2bbce4a09d60dda560387b0ea.
//
// Solidity: event TokensTransferred(address indexed user, uint256 _amount, uint256 vehicleNodeId)
func (_Issuance *IssuanceFilterer) FilterTokensTransferred(opts *bind.FilterOpts, user []common.Address) (*IssuanceTokensTransferredIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _Issuance.contract.FilterLogs(opts, "TokensTransferred", userRule)
	if err != nil {
		return nil, err
	}
	return &IssuanceTokensTransferredIterator{contract: _Issuance.contract, event: "TokensTransferred", logs: logs, sub: sub}, nil
}

// WatchTokensTransferred is a free log subscription operation binding the contract event 0x57e1000ba5ba7b6ab6670639de9fc3db34d05ef2bbce4a09d60dda560387b0ea.
//
// Solidity: event TokensTransferred(address indexed user, uint256 _amount, uint256 vehicleNodeId)
func (_Issuance *IssuanceFilterer) WatchTokensTransferred(opts *bind.WatchOpts, sink chan<- *IssuanceTokensTransferred, user []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _Issuance.contract.WatchLogs(opts, "TokensTransferred", userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IssuanceTokensTransferred)
				if err := _Issuance.contract.UnpackLog(event, "TokensTransferred", log); err != nil {
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
// Solidity: event TokensTransferred(address indexed user, uint256 _amount, uint256 vehicleNodeId)
func (_Issuance *IssuanceFilterer) ParseTokensTransferred(log types.Log) (*IssuanceTokensTransferred, error) {
	event := new(IssuanceTokensTransferred)
	if err := _Issuance.contract.UnpackLog(event, "TokensTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
