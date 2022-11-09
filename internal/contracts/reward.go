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

// RewardMetaData contains all meta data concerning the Reward contract.
var RewardMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"vehicleNodeId\",\"type\":\"uint256\"}],\"name\":\"DidntQualify\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"vehicleNodeId\",\"type\":\"uint256\"}],\"name\":\"TokensTransferred\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"users\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"values\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"vehicleIds\",\"type\":\"uint256[]\"}],\"name\":\"batchTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
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
	parsed, err := abi.JSON(strings.NewReader(RewardABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
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

// BatchTransfer is a paid mutator transaction binding the contract method 0x1ef690c4.
//
// Solidity: function batchTransfer(address[] users, uint256[] values, uint256[] vehicleIds) returns()
func (_Reward *RewardTransactor) BatchTransfer(opts *bind.TransactOpts, users []common.Address, values []*big.Int, vehicleIds []*big.Int) (*types.Transaction, error) {
	return _Reward.contract.Transact(opts, "batchTransfer", users, values, vehicleIds)
}

// BatchTransfer is a paid mutator transaction binding the contract method 0x1ef690c4.
//
// Solidity: function batchTransfer(address[] users, uint256[] values, uint256[] vehicleIds) returns()
func (_Reward *RewardSession) BatchTransfer(users []common.Address, values []*big.Int, vehicleIds []*big.Int) (*types.Transaction, error) {
	return _Reward.Contract.BatchTransfer(&_Reward.TransactOpts, users, values, vehicleIds)
}

// BatchTransfer is a paid mutator transaction binding the contract method 0x1ef690c4.
//
// Solidity: function batchTransfer(address[] users, uint256[] values, uint256[] vehicleIds) returns()
func (_Reward *RewardTransactorSession) BatchTransfer(users []common.Address, values []*big.Int, vehicleIds []*big.Int) (*types.Transaction, error) {
	return _Reward.Contract.BatchTransfer(&_Reward.TransactOpts, users, values, vehicleIds)
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
// Solidity: event DidntQualify(address indexed user, uint256 _amount, uint256 vehicleNodeId)
func (_Reward *RewardFilterer) FilterDidntQualify(opts *bind.FilterOpts, user []common.Address) (*RewardDidntQualifyIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _Reward.contract.FilterLogs(opts, "DidntQualify", userRule)
	if err != nil {
		return nil, err
	}
	return &RewardDidntQualifyIterator{contract: _Reward.contract, event: "DidntQualify", logs: logs, sub: sub}, nil
}

// WatchDidntQualify is a free log subscription operation binding the contract event 0xf667cbbd740351d2f58b68dc6a8ee1b1f3a853c609d7ce7eee4a3bdc94d62774.
//
// Solidity: event DidntQualify(address indexed user, uint256 _amount, uint256 vehicleNodeId)
func (_Reward *RewardFilterer) WatchDidntQualify(opts *bind.WatchOpts, sink chan<- *RewardDidntQualify, user []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _Reward.contract.WatchLogs(opts, "DidntQualify", userRule)
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
// Solidity: event DidntQualify(address indexed user, uint256 _amount, uint256 vehicleNodeId)
func (_Reward *RewardFilterer) ParseDidntQualify(log types.Log) (*RewardDidntQualify, error) {
	event := new(RewardDidntQualify)
	if err := _Reward.contract.UnpackLog(event, "DidntQualify", log); err != nil {
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
// Solidity: event TokensTransferred(address indexed user, uint256 _amount, uint256 vehicleNodeId)
func (_Reward *RewardFilterer) FilterTokensTransferred(opts *bind.FilterOpts, user []common.Address) (*RewardTokensTransferredIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _Reward.contract.FilterLogs(opts, "TokensTransferred", userRule)
	if err != nil {
		return nil, err
	}
	return &RewardTokensTransferredIterator{contract: _Reward.contract, event: "TokensTransferred", logs: logs, sub: sub}, nil
}

// WatchTokensTransferred is a free log subscription operation binding the contract event 0x57e1000ba5ba7b6ab6670639de9fc3db34d05ef2bbce4a09d60dda560387b0ea.
//
// Solidity: event TokensTransferred(address indexed user, uint256 _amount, uint256 vehicleNodeId)
func (_Reward *RewardFilterer) WatchTokensTransferred(opts *bind.WatchOpts, sink chan<- *RewardTokensTransferred, user []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _Reward.contract.WatchLogs(opts, "TokensTransferred", userRule)
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
// Solidity: event TokensTransferred(address indexed user, uint256 _amount, uint256 vehicleNodeId)
func (_Reward *RewardFilterer) ParseTokensTransferred(log types.Log) (*RewardTokensTransferred, error) {
	event := new(RewardTokensTransferred)
	if err := _Reward.contract.UnpackLog(event, "TokensTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
