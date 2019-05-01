// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contracts

import (
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
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = abi.U256
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// PuzzleABI is the input ABI used to generate the binding from.
const PuzzleABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"player\",\"type\":\"address\"}],\"name\":\"endGame\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"manager\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"player\",\"type\":\"address\"},{\"name\":\"level\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"string\"}],\"name\":\"payout\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"player\",\"type\":\"address\"},{\"name\":\"level\",\"type\":\"uint256\"}],\"name\":\"setLevel\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getPlayers\",\"outputs\":[{\"name\":\"\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"play\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"player\",\"type\":\"address\"}],\"name\":\"resetPlayer\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"reset\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"players\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"}]"

// PuzzleBin is the compiled bytecode used for deploying new contracts.
const PuzzleBin = `0x608060405234801561001057600080fd5b50600480546001600160a01b03191633179055610977806100326000396000f3fe6080604052600436106100865760003560e01c80638b5b9ccc116100595780638b5b9ccc146101e557806393e84cd91461024a578063c95e090914610252578063d826f88f14610285578063f71d96cb1461029a57610086565b80632a035b6c1461008b578063481c6a75146100c057806352bcd7c8146100f1578063722dcd8f146101ac575b600080fd5b34801561009757600080fd5b506100be600480360360208110156100ae57600080fd5b50356001600160a01b03166102c4565b005b3480156100cc57600080fd5b506100d56103ad565b604080516001600160a01b039092168252519081900360200190f35b6100be6004803603606081101561010757600080fd5b6001600160a01b038235169160208101359181019060608101604082013564010000000081111561013757600080fd5b82018360208201111561014957600080fd5b8035906020019184600183028401116401000000008311171561016b57600080fd5b91908080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152509295506103bc945050505050565b3480156101b857600080fd5b506100be600480360360408110156101cf57600080fd5b506001600160a01b038135169060200135610569565b3480156101f157600080fd5b506101fa610585565b60408051602080825283518183015283519192839290830191858101910280838360005b8381101561023657818101518382015260200161021e565b505050509050019250505060405180910390f35b6100be6105e8565b34801561025e57600080fd5b506100be6004803603602081101561027557600080fd5b50356001600160a01b0316610722565b34801561029157600080fd5b506100be6107f6565b3480156102a657600080fd5b506100d5600480360360208110156102bd57600080fd5b50356108dc565b6004546040805180820190915260138152600160681b72556e617574686f72697a656420416363657373026020820152906001600160a01b0316331461038b57604051600160e51b62461bcd0281526004018080602001828103825283818151815260200191508051906020019080838360005b83811015610350578181015183820152602001610338565b50505050905090810190601f16801561037d5780820380516001836020036101000a031916815260200191505b509250505060405180910390fd5b506001600160a01b03166000908152600260205260409020805460ff19169055565b6004546001600160a01b031681565b6004546040805180820190915260138152600160681b72556e617574686f72697a656420416363657373026020820152906001600160a01b0316331461044757604051600160e51b62461bcd02815260040180806020018281038252838181518152602001915080519060200190808383600083811015610350578181015183820152602001610338565b506001600160a01b038316600090815260026020908152604091829020548251808401909352601f83527f506c61796572206973206e6f7420696e20616e206163746976652067616d65009183019190915260ff1615156001146104f057604051600160e51b62461bcd02815260040180806020018281038252838181518152602001915080519060200190808383600083811015610350578181015183820152602001610338565b506001600160a01b03831660009081526020819052604090205482036105168484610569565b6001600160a01b03841660008181526001602052604080822054905160059091048402929183156108fc02918491818181858888f19350505050158015610561573d6000803e3d6000fd5b505050505050565b6001600160a01b03909116600090815260208190526040902055565b606060058054806020026020016040519081016040528092919081815260200182805480156105dd57602002820191906000526020600020905b81546001600160a01b031681526001909101906020018083116105bf575b505050505090505b90565b60408051808201909152601181527f496e73756666696369656e742046756e6400000000000000000000000000000060208201526801158e460913d0000034101561067857604051600160e51b62461bcd02815260040180806020018281038252838181518152602001915080519060200190808383600083811015610350578181015183820152602001610338565b503360009081526003602052604090205460ff1615156106ef57336000818152600360205260408120805460ff191660019081179091556005805491820181559091527f036b6384b5eca791c62761152d0c79bb0604c104a5fb6f4eb0703f3154bb3db00180546001600160a01b03191690911790555b3360009081526020818152604080832083905560018083528184203490556002909252909120805460ff19169091179055565b6004546040805180820190915260138152600160681b72556e617574686f72697a656420416363657373026020820152906001600160a01b031633146107ad57604051600160e51b62461bcd02815260040180806020018281038252838181518152602001915080519060200190808383600083811015610350578181015183820152602001610338565b506001600160a01b031660009081526020818152604080832083905560028252808320805460ff1990811690915560018352818420849055600390925290912080549091169055565b6004546040805180820190915260138152600160681b72556e617574686f72697a656420416363657373026020820152906001600160a01b0316331461088157604051600160e51b62461bcd02815260040180806020018281038252838181518152602001915080519060200190808383600083811015610350578181015183820152602001610338565b5060055460005b818110156108ca5760006005828154811015156108a157fe5b6000918252602090912001546001600160a01b031690506108c181610722565b50600101610888565b5060006108d8600582610904565b5050565b60058054829081106108ea57fe5b6000918252602090912001546001600160a01b0316905081565b8154818355818111156109285760008381526020902061092891810190830161092d565b505050565b6105e591905b808211156109475760008155600101610933565b509056fea165627a7a72305820d77629ee8b56472dfdca022f06934c4f4dfcf700ce46dccd736c4eef81d380bc0029`

// DeployPuzzle deploys a new Ethereum contract, binding an instance of Puzzle to it.
func DeployPuzzle(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Puzzle, error) {
	parsed, err := abi.JSON(strings.NewReader(PuzzleABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(PuzzleBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Puzzle{PuzzleCaller: PuzzleCaller{contract: contract}, PuzzleTransactor: PuzzleTransactor{contract: contract}, PuzzleFilterer: PuzzleFilterer{contract: contract}}, nil
}

// Puzzle is an auto generated Go binding around an Ethereum contract.
type Puzzle struct {
	PuzzleCaller     // Read-only binding to the contract
	PuzzleTransactor // Write-only binding to the contract
	PuzzleFilterer   // Log filterer for contract events
}

// PuzzleCaller is an auto generated read-only Go binding around an Ethereum contract.
type PuzzleCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PuzzleTransactor is an auto generated write-only Go binding around an Ethereum contract.
type PuzzleTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PuzzleFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type PuzzleFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PuzzleSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type PuzzleSession struct {
	Contract     *Puzzle           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// PuzzleCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type PuzzleCallerSession struct {
	Contract *PuzzleCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// PuzzleTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type PuzzleTransactorSession struct {
	Contract     *PuzzleTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// PuzzleRaw is an auto generated low-level Go binding around an Ethereum contract.
type PuzzleRaw struct {
	Contract *Puzzle // Generic contract binding to access the raw methods on
}

// PuzzleCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type PuzzleCallerRaw struct {
	Contract *PuzzleCaller // Generic read-only contract binding to access the raw methods on
}

// PuzzleTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type PuzzleTransactorRaw struct {
	Contract *PuzzleTransactor // Generic write-only contract binding to access the raw methods on
}

// NewPuzzle creates a new instance of Puzzle, bound to a specific deployed contract.
func NewPuzzle(address common.Address, backend bind.ContractBackend) (*Puzzle, error) {
	contract, err := bindPuzzle(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Puzzle{PuzzleCaller: PuzzleCaller{contract: contract}, PuzzleTransactor: PuzzleTransactor{contract: contract}, PuzzleFilterer: PuzzleFilterer{contract: contract}}, nil
}

// NewPuzzleCaller creates a new read-only instance of Puzzle, bound to a specific deployed contract.
func NewPuzzleCaller(address common.Address, caller bind.ContractCaller) (*PuzzleCaller, error) {
	contract, err := bindPuzzle(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PuzzleCaller{contract: contract}, nil
}

// NewPuzzleTransactor creates a new write-only instance of Puzzle, bound to a specific deployed contract.
func NewPuzzleTransactor(address common.Address, transactor bind.ContractTransactor) (*PuzzleTransactor, error) {
	contract, err := bindPuzzle(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PuzzleTransactor{contract: contract}, nil
}

// NewPuzzleFilterer creates a new log filterer instance of Puzzle, bound to a specific deployed contract.
func NewPuzzleFilterer(address common.Address, filterer bind.ContractFilterer) (*PuzzleFilterer, error) {
	contract, err := bindPuzzle(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PuzzleFilterer{contract: contract}, nil
}

// bindPuzzle binds a generic wrapper to an already deployed contract.
func bindPuzzle(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(PuzzleABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Puzzle *PuzzleRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Puzzle.Contract.PuzzleCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Puzzle *PuzzleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Puzzle.Contract.PuzzleTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Puzzle *PuzzleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Puzzle.Contract.PuzzleTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Puzzle *PuzzleCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Puzzle.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Puzzle *PuzzleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Puzzle.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Puzzle *PuzzleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Puzzle.Contract.contract.Transact(opts, method, params...)
}

// GetPlayers is a free data retrieval call binding the contract method 0x8b5b9ccc.
//
// Solidity: function getPlayers() constant returns(address[])
func (_Puzzle *PuzzleCaller) GetPlayers(opts *bind.CallOpts) ([]common.Address, error) {
	var (
		ret0 = new([]common.Address)
	)
	out := ret0
	err := _Puzzle.contract.Call(opts, out, "getPlayers")
	return *ret0, err
}

// GetPlayers is a free data retrieval call binding the contract method 0x8b5b9ccc.
//
// Solidity: function getPlayers() constant returns(address[])
func (_Puzzle *PuzzleSession) GetPlayers() ([]common.Address, error) {
	return _Puzzle.Contract.GetPlayers(&_Puzzle.CallOpts)
}

// GetPlayers is a free data retrieval call binding the contract method 0x8b5b9ccc.
//
// Solidity: function getPlayers() constant returns(address[])
func (_Puzzle *PuzzleCallerSession) GetPlayers() ([]common.Address, error) {
	return _Puzzle.Contract.GetPlayers(&_Puzzle.CallOpts)
}

// Manager is a free data retrieval call binding the contract method 0x481c6a75.
//
// Solidity: function manager() constant returns(address)
func (_Puzzle *PuzzleCaller) Manager(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Puzzle.contract.Call(opts, out, "manager")
	return *ret0, err
}

// Manager is a free data retrieval call binding the contract method 0x481c6a75.
//
// Solidity: function manager() constant returns(address)
func (_Puzzle *PuzzleSession) Manager() (common.Address, error) {
	return _Puzzle.Contract.Manager(&_Puzzle.CallOpts)
}

// Manager is a free data retrieval call binding the contract method 0x481c6a75.
//
// Solidity: function manager() constant returns(address)
func (_Puzzle *PuzzleCallerSession) Manager() (common.Address, error) {
	return _Puzzle.Contract.Manager(&_Puzzle.CallOpts)
}

// Players is a free data retrieval call binding the contract method 0xf71d96cb.
//
// Solidity: function players(uint256 ) constant returns(address)
func (_Puzzle *PuzzleCaller) Players(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Puzzle.contract.Call(opts, out, "players", arg0)
	return *ret0, err
}

// Players is a free data retrieval call binding the contract method 0xf71d96cb.
//
// Solidity: function players(uint256 ) constant returns(address)
func (_Puzzle *PuzzleSession) Players(arg0 *big.Int) (common.Address, error) {
	return _Puzzle.Contract.Players(&_Puzzle.CallOpts, arg0)
}

// Players is a free data retrieval call binding the contract method 0xf71d96cb.
//
// Solidity: function players(uint256 ) constant returns(address)
func (_Puzzle *PuzzleCallerSession) Players(arg0 *big.Int) (common.Address, error) {
	return _Puzzle.Contract.Players(&_Puzzle.CallOpts, arg0)
}

// EndGame is a paid mutator transaction binding the contract method 0x2a035b6c.
//
// Solidity: function endGame(address player) returns()
func (_Puzzle *PuzzleTransactor) EndGame(opts *bind.TransactOpts, player common.Address) (*types.Transaction, error) {
	return _Puzzle.contract.Transact(opts, "endGame", player)
}

// EndGame is a paid mutator transaction binding the contract method 0x2a035b6c.
//
// Solidity: function endGame(address player) returns()
func (_Puzzle *PuzzleSession) EndGame(player common.Address) (*types.Transaction, error) {
	return _Puzzle.Contract.EndGame(&_Puzzle.TransactOpts, player)
}

// EndGame is a paid mutator transaction binding the contract method 0x2a035b6c.
//
// Solidity: function endGame(address player) returns()
func (_Puzzle *PuzzleTransactorSession) EndGame(player common.Address) (*types.Transaction, error) {
	return _Puzzle.Contract.EndGame(&_Puzzle.TransactOpts, player)
}

// Payout is a paid mutator transaction binding the contract method 0x52bcd7c8.
//
// Solidity: function payout(address player, uint256 level, string ) returns()
func (_Puzzle *PuzzleTransactor) Payout(opts *bind.TransactOpts, player common.Address, level *big.Int, arg2 string) (*types.Transaction, error) {
	return _Puzzle.contract.Transact(opts, "payout", player, level, arg2)
}

// Payout is a paid mutator transaction binding the contract method 0x52bcd7c8.
//
// Solidity: function payout(address player, uint256 level, string ) returns()
func (_Puzzle *PuzzleSession) Payout(player common.Address, level *big.Int, arg2 string) (*types.Transaction, error) {
	return _Puzzle.Contract.Payout(&_Puzzle.TransactOpts, player, level, arg2)
}

// Payout is a paid mutator transaction binding the contract method 0x52bcd7c8.
//
// Solidity: function payout(address player, uint256 level, string ) returns()
func (_Puzzle *PuzzleTransactorSession) Payout(player common.Address, level *big.Int, arg2 string) (*types.Transaction, error) {
	return _Puzzle.Contract.Payout(&_Puzzle.TransactOpts, player, level, arg2)
}

// Play is a paid mutator transaction binding the contract method 0x93e84cd9.
//
// Solidity: function play() returns()
func (_Puzzle *PuzzleTransactor) Play(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Puzzle.contract.Transact(opts, "play")
}

// Play is a paid mutator transaction binding the contract method 0x93e84cd9.
//
// Solidity: function play() returns()
func (_Puzzle *PuzzleSession) Play() (*types.Transaction, error) {
	return _Puzzle.Contract.Play(&_Puzzle.TransactOpts)
}

// Play is a paid mutator transaction binding the contract method 0x93e84cd9.
//
// Solidity: function play() returns()
func (_Puzzle *PuzzleTransactorSession) Play() (*types.Transaction, error) {
	return _Puzzle.Contract.Play(&_Puzzle.TransactOpts)
}

// Reset is a paid mutator transaction binding the contract method 0xd826f88f.
//
// Solidity: function reset() returns()
func (_Puzzle *PuzzleTransactor) Reset(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Puzzle.contract.Transact(opts, "reset")
}

// Reset is a paid mutator transaction binding the contract method 0xd826f88f.
//
// Solidity: function reset() returns()
func (_Puzzle *PuzzleSession) Reset() (*types.Transaction, error) {
	return _Puzzle.Contract.Reset(&_Puzzle.TransactOpts)
}

// Reset is a paid mutator transaction binding the contract method 0xd826f88f.
//
// Solidity: function reset() returns()
func (_Puzzle *PuzzleTransactorSession) Reset() (*types.Transaction, error) {
	return _Puzzle.Contract.Reset(&_Puzzle.TransactOpts)
}

// ResetPlayer is a paid mutator transaction binding the contract method 0xc95e0909.
//
// Solidity: function resetPlayer(address player) returns()
func (_Puzzle *PuzzleTransactor) ResetPlayer(opts *bind.TransactOpts, player common.Address) (*types.Transaction, error) {
	return _Puzzle.contract.Transact(opts, "resetPlayer", player)
}

// ResetPlayer is a paid mutator transaction binding the contract method 0xc95e0909.
//
// Solidity: function resetPlayer(address player) returns()
func (_Puzzle *PuzzleSession) ResetPlayer(player common.Address) (*types.Transaction, error) {
	return _Puzzle.Contract.ResetPlayer(&_Puzzle.TransactOpts, player)
}

// ResetPlayer is a paid mutator transaction binding the contract method 0xc95e0909.
//
// Solidity: function resetPlayer(address player) returns()
func (_Puzzle *PuzzleTransactorSession) ResetPlayer(player common.Address) (*types.Transaction, error) {
	return _Puzzle.Contract.ResetPlayer(&_Puzzle.TransactOpts, player)
}

// SetLevel is a paid mutator transaction binding the contract method 0x722dcd8f.
//
// Solidity: function setLevel(address player, uint256 level) returns()
func (_Puzzle *PuzzleTransactor) SetLevel(opts *bind.TransactOpts, player common.Address, level *big.Int) (*types.Transaction, error) {
	return _Puzzle.contract.Transact(opts, "setLevel", player, level)
}

// SetLevel is a paid mutator transaction binding the contract method 0x722dcd8f.
//
// Solidity: function setLevel(address player, uint256 level) returns()
func (_Puzzle *PuzzleSession) SetLevel(player common.Address, level *big.Int) (*types.Transaction, error) {
	return _Puzzle.Contract.SetLevel(&_Puzzle.TransactOpts, player, level)
}

// SetLevel is a paid mutator transaction binding the contract method 0x722dcd8f.
//
// Solidity: function setLevel(address player, uint256 level) returns()
func (_Puzzle *PuzzleTransactorSession) SetLevel(player common.Address, level *big.Int) (*types.Transaction, error) {
	return _Puzzle.Contract.SetLevel(&_Puzzle.TransactOpts, player, level)
}
