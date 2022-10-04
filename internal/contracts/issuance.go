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

// AbiMetaData contains all meta data concerning the Abi contract.
var AbiMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"vehicleNodeId\",\"type\":\"uint256\"}],\"name\":\"TokensTransferred\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"users\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"values\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"vehicleIds\",\"type\":\"uint256[]\"}],\"name\":\"batchTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// AbiABI is the input ABI used to generate the binding from.
// Deprecated: Use AbiMetaData.ABI instead.
var AbiABI = AbiMetaData.ABI

// Abi is an auto generated Go binding around an Ethereum contract.
type Abi struct {
	AbiCaller     // Read-only binding to the contract
	AbiTransactor // Write-only binding to the contract
	AbiFilterer   // Log filterer for contract events
}

// AbiCaller is an auto generated read-only Go binding around an Ethereum contract.
type AbiCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AbiTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AbiTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AbiFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AbiFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AbiSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AbiSession struct {
	Contract     *Abi              // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AbiCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AbiCallerSession struct {
	Contract *AbiCaller    // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// AbiTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AbiTransactorSession struct {
	Contract     *AbiTransactor    // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AbiRaw is an auto generated low-level Go binding around an Ethereum contract.
type AbiRaw struct {
	Contract *Abi // Generic contract binding to access the raw methods on
}

// AbiCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AbiCallerRaw struct {
	Contract *AbiCaller // Generic read-only contract binding to access the raw methods on
}

// AbiTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AbiTransactorRaw struct {
	Contract *AbiTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAbi creates a new instance of Abi, bound to a specific deployed contract.
func NewAbi(address common.Address, backend bind.ContractBackend) (*Abi, error) {
	contract, err := bindAbi(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Abi{AbiCaller: AbiCaller{contract: contract}, AbiTransactor: AbiTransactor{contract: contract}, AbiFilterer: AbiFilterer{contract: contract}}, nil
}

// NewAbiCaller creates a new read-only instance of Abi, bound to a specific deployed contract.
func NewAbiCaller(address common.Address, caller bind.ContractCaller) (*AbiCaller, error) {
	contract, err := bindAbi(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AbiCaller{contract: contract}, nil
}

// NewAbiTransactor creates a new write-only instance of Abi, bound to a specific deployed contract.
func NewAbiTransactor(address common.Address, transactor bind.ContractTransactor) (*AbiTransactor, error) {
	contract, err := bindAbi(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AbiTransactor{contract: contract}, nil
}

// NewAbiFilterer creates a new log filterer instance of Abi, bound to a specific deployed contract.
func NewAbiFilterer(address common.Address, filterer bind.ContractFilterer) (*AbiFilterer, error) {
	contract, err := bindAbi(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AbiFilterer{contract: contract}, nil
}

// bindAbi binds a generic wrapper to an already deployed contract.
func bindAbi(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(AbiABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Abi *AbiRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Abi.Contract.AbiCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Abi *AbiRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Abi.Contract.AbiTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Abi *AbiRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Abi.Contract.AbiTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Abi *AbiCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Abi.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Abi *AbiTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Abi.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Abi *AbiTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Abi.Contract.contract.Transact(opts, method, params...)
}

// BatchTransfer is a paid mutator transaction binding the contract method 0x1ef690c4.
//
// Solidity: function batchTransfer(address[] users, uint256[] values, uint256[] vehicleIds) returns()
func (_Abi *AbiTransactor) BatchTransfer(opts *bind.TransactOpts, users []common.Address, values []*big.Int, vehicleIds []*big.Int) (*types.Transaction, error) {
	return _Abi.contract.Transact(opts, "batchTransfer", users, values, vehicleIds)
}

// BatchTransfer is a paid mutator transaction binding the contract method 0x1ef690c4.
//
// Solidity: function batchTransfer(address[] users, uint256[] values, uint256[] vehicleIds) returns()
func (_Abi *AbiSession) BatchTransfer(users []common.Address, values []*big.Int, vehicleIds []*big.Int) (*types.Transaction, error) {
	return _Abi.Contract.BatchTransfer(&_Abi.TransactOpts, users, values, vehicleIds)
}

// BatchTransfer is a paid mutator transaction binding the contract method 0x1ef690c4.
//
// Solidity: function batchTransfer(address[] users, uint256[] values, uint256[] vehicleIds) returns()
func (_Abi *AbiTransactorSession) BatchTransfer(users []common.Address, values []*big.Int, vehicleIds []*big.Int) (*types.Transaction, error) {
	return _Abi.Contract.BatchTransfer(&_Abi.TransactOpts, users, values, vehicleIds)
}

// AbiTokensTransferredIterator is returned from FilterTokensTransferred and is used to iterate over the raw logs and unpacked data for TokensTransferred events raised by the Abi contract.
type AbiTokensTransferredIterator struct {
	Event *AbiTokensTransferred // Event containing the contract specifics and raw log

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
func (it *AbiTokensTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AbiTokensTransferred)
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
		it.Event = new(AbiTokensTransferred)
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
func (it *AbiTokensTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AbiTokensTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AbiTokensTransferred represents a TokensTransferred event raised by the Abi contract.
type AbiTokensTransferred struct {
	User          common.Address
	Amount        *big.Int
	VehicleNodeId *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterTokensTransferred is a free log retrieval operation binding the contract event 0x57e1000ba5ba7b6ab6670639de9fc3db34d05ef2bbce4a09d60dda560387b0ea.
//
// Solidity: event TokensTransferred(address indexed user, uint256 _amount, uint256 vehicleNodeId)
func (_Abi *AbiFilterer) FilterTokensTransferred(opts *bind.FilterOpts, user []common.Address) (*AbiTokensTransferredIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _Abi.contract.FilterLogs(opts, "TokensTransferred", userRule)
	if err != nil {
		return nil, err
	}
	return &AbiTokensTransferredIterator{contract: _Abi.contract, event: "TokensTransferred", logs: logs, sub: sub}, nil
}

// WatchTokensTransferred is a free log subscription operation binding the contract event 0x57e1000ba5ba7b6ab6670639de9fc3db34d05ef2bbce4a09d60dda560387b0ea.
//
// Solidity: event TokensTransferred(address indexed user, uint256 _amount, uint256 vehicleNodeId)
func (_Abi *AbiFilterer) WatchTokensTransferred(opts *bind.WatchOpts, sink chan<- *AbiTokensTransferred, user []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _Abi.contract.WatchLogs(opts, "TokensTransferred", userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AbiTokensTransferred)
				if err := _Abi.contract.UnpackLog(event, "TokensTransferred", log); err != nil {
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
func (_Abi *AbiFilterer) ParseTokensTransferred(log types.Log) (*AbiTokensTransferred, error) {
	event := new(AbiTokensTransferred)
	if err := _Abi.contract.UnpackLog(event, "TokensTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
