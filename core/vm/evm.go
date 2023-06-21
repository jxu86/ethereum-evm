// Copyright 2014 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package vm

import (
	// "ethereum-evm/params"
	"fmt"
	"math/big"

	// "time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/crypto"
)

type codeAndHash struct {
	code []byte
	hash common.Hash
}

// BlockContext provides the EVM with auxiliary information. Once provided
// it shouldn't be modified.
type BlockContext struct {
	// // CanTransfer returns whether the account contains
	// // sufficient ether to transfer the value
	// CanTransfer CanTransferFunc
	// // Transfer transfers ether from one account to the other
	// Transfer TransferFunc
	// // GetHash returns the hash corresponding to n
	// GetHash GetHashFunc

	// Block information
	Coinbase    common.Address // Provides information for COINBASE
	GasLimit    uint64         // Provides information for GASLIMIT
	BlockNumber *big.Int       // Provides information for NUMBER
	Time        *big.Int       // Provides information for TIME
	Difficulty  *big.Int       // Provides information for DIFFICULTY
	BaseFee     *big.Int       // Provides information for BASEFEE
}

type EVM struct {
	// Context provides auxiliary blockchain related information
	Context BlockContext
	// TxContext
	// // StateDB gives access to the underlying state
	StateDB StateDB
	// // Depth is the current call stack
	depth int

	// // chainConfig contains information about the current chain
	// chainConfig *params.ChainConfig
	// // chain rules contains the chain rules for the current epoch
	// chainRules params.Rules
	// // virtual machine configuration options used to initialise the
	// // evm.
	// Config Config
	// // global (to this context) ethereum virtual machine
	// // used throughout the execution of the tx.
	interpreter *EVMInterpreter
	// // abort is used to abort the EVM calling operations
	// // NOTE: must be set atomically
	abort int32
	// // callGasTemp holds the gas available for the current call. This is needed because the
	// // available gas is calculated in gasCall* according to the 63/64 rule and later
	// // applied in opCall*.
	callGasTemp uint64
}

func NewEVM(statedb StateDB) *EVM {
	evm := &EVM{
		depth:   10,
		StateDB: statedb,
	}
	evm.interpreter = NewEVMInterpreter(evm, Config{})
	return evm
}

// Interpreter returns the current interpreter
func (evm *EVM) Interpreter() *EVMInterpreter {
	return evm.interpreter
}

// create creates a new contract using code as deployment code.
func (evm *EVM) create(caller ContractRef, codeAndHash *codeAndHash, gas uint64, value *big.Int, address common.Address) ([]byte, common.Address, uint64, error) {
	// Depth check execution. Fail if we're trying to execute above the
	// limit.
	// 检查合约创建的递归调用次数，防止在合约中创建合约的递归次数
	// if evm.depth > int(params.CallCreateDepth) {
	// 	return nil, common.Address{}, gas, ErrDepth
	// }
	// 检查合约创建者是否有足够的以太币
	// if !evm.Context.CanTransfer(evm.StateDB, caller.Address(), value) {
	// 	return nil, common.Address{}, gas, ErrInsufficientBalance
	// }
	// 增加合约创建者的 Nonce 值
	// fmt.Println("###caller.Address()=>", caller.Address())
	nonce := evm.StateDB.GetNonce(caller.Address())
	fmt.Println("###nonce1=>", nonce)
	// evm.StateDB.SetNonce(caller.Address(), nonce+1)
	// nonce = evm.StateDB.GetNonce(caller.Address())
	// fmt.Println("###nonce=>", nonce)
	// We add this to the access list _before_ taking a snapshot. Even if the creation fails,
	// the access-list change should not be rolled back
	// if evm.chainRules.IsBerlin {
	// 	evm.StateDB.AddAddressToAccessList(address)
	// }
	// Ensure there's no existing contract already at the designated address
	// 确认合约账户地址不存在，否则返回错误
	contractHash := evm.StateDB.GetCodeHash(address)
	// if evm.StateDB.GetNonce(address) != 0 || (contractHash != (common.Hash{}) && contractHash != emptyCodeHash) {
	if evm.StateDB.GetNonce(address) != 0 || (contractHash != (common.Hash{}) && contractHash != crypto.Keccak256Hash(nil)) {
		return nil, common.Address{}, 0, ErrContractAddressCollision
	}
	// Create a new account on the state
	// snapshot := evm.StateDB.Snapshot()
	// evm.StateDB.CreateAccount(address)
	// if evm.chainRules.IsEIP158 {
	// 	evm.StateDB.SetNonce(address, 1)
	// }
	// 把以太币(如果需要)转账到这个新建的合约地址上
	// evm.Context.Transfer(evm.StateDB, caller.Address(), address, value)

	// Initialise a new contract and set the code that is to be used by the EVM.
	// The contract is a scoped environment for this execution context only.
	// 创建一个Contract对象
	contract := NewContract(caller, AccountRef(address), value, gas)
	contract.SetCodeOptionalHash(&address, codeAndHash)

	// if evm.Config.NoRecursion && evm.depth > 0 {
	// 	return nil, address, gas, nil
	// }

	// if evm.Config.Debug && evm.depth == 0 {
	// 	evm.Config.Tracer.CaptureStart(evm, caller.Address(), address, true, codeAndHash.code, gas, value)
	// }
	// start := time.Now()
	// 运行合约代码，应该是说运行部署合约的代码，真正合约的代码是返回的ret
	ret, err := evm.interpreter.Run(contract, nil, false)
	fmt.Println("#######run ret=>", ret)
	// Check whether the max code size has been exceeded, assign err if the case.
	// 检查合约代码长度是否超过限制
	// if err == nil && evm.chainRules.IsEIP158 && len(ret) > params.MaxCodeSize {
	// 	err = ErrMaxCodeSizeExceeded
	// }

	// Reject code starting with 0xEF if EIP-3541 is enabled.
	// if err == nil && len(ret) >= 1 && ret[0] == 0xEF && evm.chainRules.IsLondon {
	// 	err = ErrInvalidCode
	// }

	// if the contract creation ran successfully and no errors were returned
	// calculate the gas required to store the code. If the code could not
	// be stored due to not enough gas set an error and let it be handled
	// by the error checking condition below.
	if err == nil {
		// createDataGas := uint64(len(ret)) * params.CreateDataGas
		// if contract.UseGas(createDataGas) {
		// 	// 将合约代码存储到以太坊状态数据库的合约账户中,需要消耗一定的gas
		// 	evm.StateDB.SetCode(address, ret)
		// } else {
		// 	err = ErrCodeStoreOutOfGas
		// }
		evm.StateDB.SetCode(address, ret)

		evm.StateDB.SetNonce(caller.Address(), nonce+1)
		nonce = evm.StateDB.GetNonce(caller.Address())
		fmt.Println("###nonce=>", nonce)
	}

	// When an error was returned by the EVM or when setting the creation code
	// above we revert to the snapshot and consume any gas remaining. Additionally
	// when we're in homestead this also counts for code storage gas errors.
	// if err != nil && (evm.chainRules.IsHomestead || err != ErrCodeStoreOutOfGas) {
	// 	evm.StateDB.RevertToSnapshot(snapshot)
	// 	if err != ErrExecutionReverted {
	// 		contract.UseGas(contract.Gas)
	// 	}
	// }

	// if evm.Config.Debug && evm.depth == 0 {
	// 	evm.Config.Tracer.CaptureEnd(ret, gas-contract.Gas, time.Since(start), err)
	// }
	return ret, address, contract.Gas, err
}

// Create creates a new contract using code as deployment code.
func (evm *EVM) Create(caller ContractRef, code []byte, gas uint64, value *big.Int) (ret []byte, contractAddr common.Address, leftOverGas uint64, err error) {
	// contractAddr = crypto.CreateAddress(caller.Address(), evm.StateDB.GetNonce(caller.Address()))
	contractAddr = crypto.CreateAddress(caller.Address(), evm.StateDB.GetNonce(caller.Address()))
	fmt.Println("###contractAddr: ", contractAddr)
	return evm.create(caller, &codeAndHash{code: code}, gas, value, contractAddr)
}

func (evm *EVM) Call(caller ContractRef, addr common.Address, input []byte, gas uint64, value *big.Int) (ret []byte, leftOverGas uint64, err error) {
	// if evm.Config.NoRecursion && evm.depth > 0 {
	// 	return nil, gas, nil
	// }
	// Fail if we're trying to execute above the call depth limit
	// if evm.depth > int(params.CallCreateDepth) {
	// 	return nil, gas, ErrDepth
	// }
	// Fail if we're trying to transfer more than the available balance
	// if value.Sign() != 0 && !evm.Context.CanTransfer(evm.StateDB, caller.Address(), value) {
	// 	return nil, gas, ErrInsufficientBalance
	// }
	// snapshot := evm.StateDB.Snapshot()
	// p, isPrecompile := evm.precompile(addr)
	isPrecompile := false

	if !evm.StateDB.Exist(addr) {
		// if !isPrecompile && evm.chainRules.IsEIP158 && value.Sign() == 0 {
		// 	// Calling a non existing account, don't do anything, but ping the tracer
		// 	if evm.Config.Debug && evm.depth == 0 {
		// 		evm.Config.Tracer.CaptureStart(evm, caller.Address(), addr, false, input, gas, value)
		// 		evm.Config.Tracer.CaptureEnd(ret, 0, 0, nil)
		// 	}
		// 	return nil, gas, nil
		// }
		evm.StateDB.CreateAccount(addr)
	}
	// evm.Context.Transfer(evm.StateDB, caller.Address(), addr, value)

	// Capture the tracer start/end events in debug mode
	// if evm.Config.Debug && evm.depth == 0 {
	// 	evm.Config.Tracer.CaptureStart(evm, caller.Address(), addr, false, input, gas, value)
	// 	defer func(startGas uint64, startTime time.Time) { // Lazy evaluation of the parameters
	// 		evm.Config.Tracer.CaptureEnd(ret, startGas-gas, time.Since(startTime), err)
	// 	}(gas, time.Now())
	// }

	if isPrecompile {
		// ret, gas, err = RunPrecompiledContract(p, input, gas)
	} else {
		// Initialise a new contract and set the code that is to be used by the EVM.
		// The contract is a scoped environment for this execution context only.
		fmt.Println("##@@contract addr: ", addr)
		code := evm.StateDB.GetCode(addr)
		fmt.Printf("##@@contract code: %x\n", code)
		if len(code) == 0 { // 没有合约代码，普通转账
			ret, err = nil, nil // gas is unchanged
		} else {
			addrCopy := addr
			// If the account has no code, we can abort here
			// The depth-check is already done, and precompiles handled above
			contract := NewContract(caller, AccountRef(addrCopy), value, gas)
			contract.SetCallCode(&addrCopy, evm.StateDB.GetCodeHash(addrCopy), code)
			ret, err = evm.interpreter.Run(contract, input, false)
			fmt.Printf("#############ret============> %s\n", ret)
			// gas = contract.Gas
		}
	}
	// When an error was returned by the EVM or when setting the creation code
	// above we revert to the snapshot and consume any gas remaining. Additionally
	// when we're in homestead this also counts for code storage gas errors.
	// if err != nil {
	// 	evm.StateDB.RevertToSnapshot(snapshot)
	// 	if err != ErrExecutionReverted {
	// 		gas = 0
	// 	}
	// 	// TODO: consider clearing up unused snapshots:
	// 	//} else {
	// 	//	evm.StateDB.DiscardSnapshot(snapshot)
	// }
	return ret, gas, err
}
