// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package atomic_swap

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// AtomicSwapEtherABI is the input ABI used to generate the binding from.
const AtomicSwapEtherABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"}],\"name\":\"getState\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"},{\"name\":\"_secretKey\",\"type\":\"bytes\"}],\"name\":\"close\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"}],\"name\":\"check\",\"outputs\":[{\"name\":\"timeRemaining\",\"type\":\"uint256\"},{\"name\":\"value\",\"type\":\"uint256\"},{\"name\":\"withdrawTrader\",\"type\":\"address\"},{\"name\":\"secretLock\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"},{\"name\":\"_withdrawTrader\",\"type\":\"address\"},{\"name\":\"_secretLock\",\"type\":\"bytes32\"}],\"name\":\"open\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"}],\"name\":\"expire\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"}],\"name\":\"checkSecretKey\",\"outputs\":[{\"name\":\"secretKey\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_swapID\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"_withdrawTrader\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_secretLock\",\"type\":\"bytes32\"}],\"name\":\"Open\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_swapID\",\"type\":\"bytes32\"}],\"name\":\"Expire\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_swapID\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"_secretKey\",\"type\":\"bytes\"}],\"name\":\"Close\",\"type\":\"event\"}]"

// AtomicSwapEtherBin is the compiled bytecode used for deploying new contracts.
const AtomicSwapEtherBin = `0x6060604052341561000f57600080fd5b610c4a8061001e6000396000f3006060604052600436106100775763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166309648a9d811461007c5780631a26720c146100b6578063399e07921461010e5780633ee94b1314610159578063c644179814610173578063f200e40414610189575b600080fd5b341561008757600080fd5b610092600435610216565b604051808260038111156100a257fe5b60ff16815260200191505060405180910390f35b34156100c157600080fd5b61010c600480359060446024803590810190830135806020601f8201819004810201604051908101604052818152929190602084018383808284375094965061022b95505050505050565b005b341561011957600080fd5b61012460043561051f565b6040519384526020840192909252600160a060020a031660408084019190915260608301919091526080909101905180910390f35b61010c600435600160a060020a0360243516604435610656565b341561017e57600080fd5b61010c60043561080e565b341561019457600080fd5b61019f6004356109ed565b60405160208082528190810183818151815260200191508051906020019080838360005b838110156101db5780820151838201526020016101c3565b50505050905090810190601f1680156102085780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b60009081526001602052604090205460ff1690565b610233610b35565b82600160008281526001602052604090205460ff16600381111561025357fe5b14156105195783836002816000604051602001526040518082805190602001908083835b602083106102965780518252601f199092019160209182019101610277565b6001836020036101000a03801982511681845116808217855250505050505090500191505060206040518083038160008661646e5a03f115156102d857600080fd5b505060405180516000848152602081905260409020600401541415905061051657600086815260208190526040908190209060c09051908101604090815282548252600180840154602080850191909152600280860154600160a060020a039081168587015260038701541660608601526004860154608086015260058601805495969560a088019591948116156101000260001901169190910491601f8301819004810201905190810160405280929190818152602001828054600181600116156101000203166002900480156103f15780601f106103c6576101008083540402835291602001916103f1565b820191906000526020600020905b8154815290600101906020018083116103d457829003601f168201915b505050919092525050506000878152602081905260409020909450600501858051610420929160200190610b71565b506000868152600160205260409020805460ff191660021790556060840151600160a060020a03166108fc85602001519081150290604051600060405180830381858888f19350505050151561047557600080fd5b7f692fd10a275135b9a2a2f5819db3d9965a5129ea2ad3640a0156dbce2fc81bdd868660405182815260406020820181815290820183818151815260200191508051906020019080838360005b838110156104da5780820151838201526020016104c2565b50505050905090810190601f1680156105075780820380516001836020036101000a031916815260200191505b50935050505060405180910390a15b50505b50505050565b60008060008061052d610b35565b600086815260208190526040908190209060c09051908101604090815282548252600180840154602080850191909152600280860154600160a060020a039081168587015260038701541660608601526004860154608086015260058601805495969560a088019591948116156101000260001901169190910491601f8301819004810201905190810160405280929190818152602001828054600181600116156101000203166002900480156106255780601f106105fa57610100808354040283529160200191610625565b820191906000526020600020905b81548152906001019060200180831161060857829003601f168201915b5050505050815250509050428160000151038160200151826060015183608001519450945094509450509193509193565b61065e610b35565b836000808281526001602052604090205460ff16600381111561067d57fe5b14156108075760c06040519081016040528042815260200134815260200133600160a060020a0316815260200185600160a060020a031681526020018460001916815260200160006040518059106106d25750595b818152601f19601f830116810160200160405290509052600086815260208190526040902090925082908151815560208201518160010155604082015160028201805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a0392909216919091179055606082015160038201805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a03929092169190911790556080820151600482015560a082015181600501908051610796929160200190610b71565b5050506000858152600160208190526040909120805460ff1916828002179055507f6ed79a08bf5c8a7d4a330df315e4ac386627ecafbe5d2bfd6654237d967b24f3858585604051928352600160a060020a0390911660208301526040808301919091526060909101905180910390a15b5050505050565b610816610b35565b81600160008281526001602052604090205460ff16600381111561083657fe5b14156109e857600083815260208190526040902054839062015180429091031061051957600084815260208190526040908190209060c09051908101604090815282548252600180840154602080850191909152600280860154600160a060020a039081168587015260038701541660608601526004860154608086015260058601805495969560a088019591948116156101000260001901169190910491601f8301819004810201905190810160405280929190818152602001828054600181600116156101000203166002900480156109525780601f1061092757610100808354040283529160200191610952565b820191906000526020600020905b81548152906001019060200180831161093557829003601f168201915b5050509190925250505060008581526001602052604090819020805460ff19166003179055909350830151600160a060020a03166108fc84602001519081150290604051600060405180830381858888f1935050505015156109b357600080fd5b7fbddd9b693ea862fad6ecf78fd51c065be26fda94d1f3cad3a7d691453a38a7358460405190815260200160405180910390a1505b505050565b6109f5610bef565b6109fd610b35565b82600260008281526001602052604090205460ff166003811115610a1d57fe5b1415610b2e57600084815260208190526040908190209060c09051908101604090815282548252600180840154602080850191909152600280860154600160a060020a039081168587015260038701541660608601526004860154608086015260058601805495969560a088019591948116156101000260001901169190910491601f830181900481020190519081016040528092919081815260200182805460018160011615610100020316600290048015610b1b5780601f10610af057610100808354040283529160200191610b1b565b820191906000526020600020905b815481529060010190602001808311610afe57829003601f168201915b50505050508152505091508160a0015192505b5050919050565b60c0604051908101604090815260008083526020830181905290820181905260608201819052608082015260a08101610b6c610bef565b905290565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f10610bb257805160ff1916838001178555610bdf565b82800160010185558215610bdf579182015b82811115610bdf578251825591602001919060010190610bc4565b50610beb929150610c01565b5090565b60206040519081016040526000815290565b610c1b91905b80821115610beb5760008155600101610c07565b905600a165627a7a7230582015a0a09844c821d556e3391bcd53b192b5921ff4bfe68b5eb33b559309f5e3c30029`

// DeployAtomicSwapEther deploys a new Ethereum contract, binding an instance of AtomicSwapEther to it.
func DeployAtomicSwapEther(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *AtomicSwapEther, error) {
	parsed, err := abi.JSON(strings.NewReader(AtomicSwapEtherABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(AtomicSwapEtherBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &AtomicSwapEther{AtomicSwapEtherCaller: AtomicSwapEtherCaller{contract: contract}, AtomicSwapEtherTransactor: AtomicSwapEtherTransactor{contract: contract}}, nil
}

// AtomicSwapEther is an auto generated Go binding around an Ethereum contract.
type AtomicSwapEther struct {
	AtomicSwapEtherCaller     // Read-only binding to the contract
	AtomicSwapEtherTransactor // Write-only binding to the contract
}

// AtomicSwapEtherCaller is an auto generated read-only Go binding around an Ethereum contract.
type AtomicSwapEtherCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AtomicSwapEtherTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AtomicSwapEtherTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AtomicSwapEtherSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AtomicSwapEtherSession struct {
	Contract     *AtomicSwapEther  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AtomicSwapEtherCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AtomicSwapEtherCallerSession struct {
	Contract *AtomicSwapEtherCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// AtomicSwapEtherTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AtomicSwapEtherTransactorSession struct {
	Contract     *AtomicSwapEtherTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// AtomicSwapEtherRaw is an auto generated low-level Go binding around an Ethereum contract.
type AtomicSwapEtherRaw struct {
	Contract *AtomicSwapEther // Generic contract binding to access the raw methods on
}

// AtomicSwapEtherCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AtomicSwapEtherCallerRaw struct {
	Contract *AtomicSwapEtherCaller // Generic read-only contract binding to access the raw methods on
}

// AtomicSwapEtherTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AtomicSwapEtherTransactorRaw struct {
	Contract *AtomicSwapEtherTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAtomicSwapEther creates a new instance of AtomicSwapEther, bound to a specific deployed contract.
func NewAtomicSwapEther(address common.Address, backend bind.ContractBackend) (*AtomicSwapEther, error) {
	contract, err := bindAtomicSwapEther(address, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AtomicSwapEther{AtomicSwapEtherCaller: AtomicSwapEtherCaller{contract: contract}, AtomicSwapEtherTransactor: AtomicSwapEtherTransactor{contract: contract}}, nil
}

// NewAtomicSwapEtherCaller creates a new read-only instance of AtomicSwapEther, bound to a specific deployed contract.
func NewAtomicSwapEtherCaller(address common.Address, caller bind.ContractCaller) (*AtomicSwapEtherCaller, error) {
	contract, err := bindAtomicSwapEther(address, caller, nil)
	if err != nil {
		return nil, err
	}
	return &AtomicSwapEtherCaller{contract: contract}, nil
}

// NewAtomicSwapEtherTransactor creates a new write-only instance of AtomicSwapEther, bound to a specific deployed contract.
func NewAtomicSwapEtherTransactor(address common.Address, transactor bind.ContractTransactor) (*AtomicSwapEtherTransactor, error) {
	contract, err := bindAtomicSwapEther(address, nil, transactor)
	if err != nil {
		return nil, err
	}
	return &AtomicSwapEtherTransactor{contract: contract}, nil
}

// bindAtomicSwapEther binds a generic wrapper to an already deployed contract.
func bindAtomicSwapEther(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(AtomicSwapEtherABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AtomicSwapEther *AtomicSwapEtherRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _AtomicSwapEther.Contract.AtomicSwapEtherCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AtomicSwapEther *AtomicSwapEtherRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AtomicSwapEther.Contract.AtomicSwapEtherTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AtomicSwapEther *AtomicSwapEtherRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AtomicSwapEther.Contract.AtomicSwapEtherTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AtomicSwapEther *AtomicSwapEtherCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _AtomicSwapEther.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AtomicSwapEther *AtomicSwapEtherTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AtomicSwapEther.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AtomicSwapEther *AtomicSwapEtherTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AtomicSwapEther.Contract.contract.Transact(opts, method, params...)
}

// Check is a free data retrieval call binding the contract method 0x399e0792.
//
// Solidity: function check(_swapID bytes32) constant returns(timeRemaining uint256, value uint256, withdrawTrader address, secretLock bytes32)
func (_AtomicSwapEther *AtomicSwapEtherCaller) Check(opts *bind.CallOpts, _swapID [32]byte) (struct {
	TimeRemaining  *big.Int
	Value          *big.Int
	WithdrawTrader common.Address
	SecretLock     [32]byte
}, error) {
	ret := new(struct {
		TimeRemaining  *big.Int
		Value          *big.Int
		WithdrawTrader common.Address
		SecretLock     [32]byte
	})
	out := ret
	err := _AtomicSwapEther.contract.Call(opts, out, "check", _swapID)
	return *ret, err
}

// Check is a free data retrieval call binding the contract method 0x399e0792.
//
// Solidity: function check(_swapID bytes32) constant returns(timeRemaining uint256, value uint256, withdrawTrader address, secretLock bytes32)
func (_AtomicSwapEther *AtomicSwapEtherSession) Check(_swapID [32]byte) (struct {
	TimeRemaining  *big.Int
	Value          *big.Int
	WithdrawTrader common.Address
	SecretLock     [32]byte
}, error) {
	return _AtomicSwapEther.Contract.Check(&_AtomicSwapEther.CallOpts, _swapID)
}

// Check is a free data retrieval call binding the contract method 0x399e0792.
//
// Solidity: function check(_swapID bytes32) constant returns(timeRemaining uint256, value uint256, withdrawTrader address, secretLock bytes32)
func (_AtomicSwapEther *AtomicSwapEtherCallerSession) Check(_swapID [32]byte) (struct {
	TimeRemaining  *big.Int
	Value          *big.Int
	WithdrawTrader common.Address
	SecretLock     [32]byte
}, error) {
	return _AtomicSwapEther.Contract.Check(&_AtomicSwapEther.CallOpts, _swapID)
}

// CheckSecretKey is a free data retrieval call binding the contract method 0xf200e404.
//
// Solidity: function checkSecretKey(_swapID bytes32) constant returns(secretKey bytes)
func (_AtomicSwapEther *AtomicSwapEtherCaller) CheckSecretKey(opts *bind.CallOpts, _swapID [32]byte) ([]byte, error) {
	var (
		ret0 = new([]byte)
	)
	out := ret0
	err := _AtomicSwapEther.contract.Call(opts, out, "checkSecretKey", _swapID)
	return *ret0, err
}

// CheckSecretKey is a free data retrieval call binding the contract method 0xf200e404.
//
// Solidity: function checkSecretKey(_swapID bytes32) constant returns(secretKey bytes)
func (_AtomicSwapEther *AtomicSwapEtherSession) CheckSecretKey(_swapID [32]byte) ([]byte, error) {
	return _AtomicSwapEther.Contract.CheckSecretKey(&_AtomicSwapEther.CallOpts, _swapID)
}

// CheckSecretKey is a free data retrieval call binding the contract method 0xf200e404.
//
// Solidity: function checkSecretKey(_swapID bytes32) constant returns(secretKey bytes)
func (_AtomicSwapEther *AtomicSwapEtherCallerSession) CheckSecretKey(_swapID [32]byte) ([]byte, error) {
	return _AtomicSwapEther.Contract.CheckSecretKey(&_AtomicSwapEther.CallOpts, _swapID)
}

// GetState is a free data retrieval call binding the contract method 0x09648a9d.
//
// Solidity: function getState(_swapID bytes32) constant returns(uint8)
func (_AtomicSwapEther *AtomicSwapEtherCaller) GetState(opts *bind.CallOpts, _swapID [32]byte) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _AtomicSwapEther.contract.Call(opts, out, "getState", _swapID)
	return *ret0, err
}

// GetState is a free data retrieval call binding the contract method 0x09648a9d.
//
// Solidity: function getState(_swapID bytes32) constant returns(uint8)
func (_AtomicSwapEther *AtomicSwapEtherSession) GetState(_swapID [32]byte) (uint8, error) {
	return _AtomicSwapEther.Contract.GetState(&_AtomicSwapEther.CallOpts, _swapID)
}

// GetState is a free data retrieval call binding the contract method 0x09648a9d.
//
// Solidity: function getState(_swapID bytes32) constant returns(uint8)
func (_AtomicSwapEther *AtomicSwapEtherCallerSession) GetState(_swapID [32]byte) (uint8, error) {
	return _AtomicSwapEther.Contract.GetState(&_AtomicSwapEther.CallOpts, _swapID)
}

// Close is a paid mutator transaction binding the contract method 0x1a26720c.
//
// Solidity: function close(_swapID bytes32, _secretKey bytes) returns()
func (_AtomicSwapEther *AtomicSwapEtherTransactor) Close(opts *bind.TransactOpts, _swapID [32]byte, _secretKey []byte) (*types.Transaction, error) {
	return _AtomicSwapEther.contract.Transact(opts, "close", _swapID, _secretKey)
}

// Close is a paid mutator transaction binding the contract method 0x1a26720c.
//
// Solidity: function close(_swapID bytes32, _secretKey bytes) returns()
func (_AtomicSwapEther *AtomicSwapEtherSession) Close(_swapID [32]byte, _secretKey []byte) (*types.Transaction, error) {
	return _AtomicSwapEther.Contract.Close(&_AtomicSwapEther.TransactOpts, _swapID, _secretKey)
}

// Close is a paid mutator transaction binding the contract method 0x1a26720c.
//
// Solidity: function close(_swapID bytes32, _secretKey bytes) returns()
func (_AtomicSwapEther *AtomicSwapEtherTransactorSession) Close(_swapID [32]byte, _secretKey []byte) (*types.Transaction, error) {
	return _AtomicSwapEther.Contract.Close(&_AtomicSwapEther.TransactOpts, _swapID, _secretKey)
}

// Expire is a paid mutator transaction binding the contract method 0xc6441798.
//
// Solidity: function expire(_swapID bytes32) returns()
func (_AtomicSwapEther *AtomicSwapEtherTransactor) Expire(opts *bind.TransactOpts, _swapID [32]byte) (*types.Transaction, error) {
	return _AtomicSwapEther.contract.Transact(opts, "expire", _swapID)
}

// Expire is a paid mutator transaction binding the contract method 0xc6441798.
//
// Solidity: function expire(_swapID bytes32) returns()
func (_AtomicSwapEther *AtomicSwapEtherSession) Expire(_swapID [32]byte) (*types.Transaction, error) {
	return _AtomicSwapEther.Contract.Expire(&_AtomicSwapEther.TransactOpts, _swapID)
}

// Expire is a paid mutator transaction binding the contract method 0xc6441798.
//
// Solidity: function expire(_swapID bytes32) returns()
func (_AtomicSwapEther *AtomicSwapEtherTransactorSession) Expire(_swapID [32]byte) (*types.Transaction, error) {
	return _AtomicSwapEther.Contract.Expire(&_AtomicSwapEther.TransactOpts, _swapID)
}

// Open is a paid mutator transaction binding the contract method 0x3ee94b13.
//
// Solidity: function open(_swapID bytes32, _withdrawTrader address, _secretLock bytes32) returns()
func (_AtomicSwapEther *AtomicSwapEtherTransactor) Open(opts *bind.TransactOpts, _swapID [32]byte, _withdrawTrader common.Address, _secretLock [32]byte) (*types.Transaction, error) {
	return _AtomicSwapEther.contract.Transact(opts, "open", _swapID, _withdrawTrader, _secretLock)
}

// Open is a paid mutator transaction binding the contract method 0x3ee94b13.
//
// Solidity: function open(_swapID bytes32, _withdrawTrader address, _secretLock bytes32) returns()
func (_AtomicSwapEther *AtomicSwapEtherSession) Open(_swapID [32]byte, _withdrawTrader common.Address, _secretLock [32]byte) (*types.Transaction, error) {
	return _AtomicSwapEther.Contract.Open(&_AtomicSwapEther.TransactOpts, _swapID, _withdrawTrader, _secretLock)
}

// Open is a paid mutator transaction binding the contract method 0x3ee94b13.
//
// Solidity: function open(_swapID bytes32, _withdrawTrader address, _secretLock bytes32) returns()
func (_AtomicSwapEther *AtomicSwapEtherTransactorSession) Open(_swapID [32]byte, _withdrawTrader common.Address, _secretLock [32]byte) (*types.Transaction, error) {
	return _AtomicSwapEther.Contract.Open(&_AtomicSwapEther.TransactOpts, _swapID, _withdrawTrader, _secretLock)
}
