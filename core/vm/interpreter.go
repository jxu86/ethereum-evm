package vm

import (
	"bufio"
	"encoding/hex"
	"ethereum-evm/common"
	"ethereum-evm/common/math"
	"fmt"
	"hash"
	"math/big"
	"os"

	"github.com/cloudflare/cfssl/log"
	"github.com/ethereum/go-ethereum/crypto"
)

type Config struct {
	Debug bool // Enables debugging
	// Tracer                  Tracer // Opcode logger
	NoRecursion             bool // Disables call, callcode, delegate call and create
	NoBaseFee               bool // Forces the EIP-1559 baseFee to 0 (needed for 0 price calls)
	EnablePreimageRecording bool // Enables recording of SHA3/keccak preimages

	JumpTable [256]*operation // EVM instruction table, automatically populated if unset

	ExtraEips []int // Additional EIPS that are to be enabled
}

// ScopeContext contains the things that are per-call, such as stack and memory,
// but not transients like pc and gas
type ScopeContext struct {
	Memory   *Memory
	Stack    *Stack
	Contract *Contract
}

type keccakState interface {
	hash.Hash
	Read([]byte) (int, error)
}

type EVMInterpreter struct {
	evm *EVM
	cfg Config

	hasher    keccakState // Keccak256 hasher instance shared across opcodes
	hasherBuf common.Hash // Keccak256 hasher result array shared aross opcodes

	readOnly   bool   // Whether to throw on stateful modifications
	returnData []byte // Last CALL's return data for subsequent reuse
}

// NewEVMInterpreter returns a new instance of the Interpreter.
func NewEVMInterpreter(evm *EVM, cfg Config) *EVMInterpreter {
	// func NewEVMInterpreter() *EVMInterpreter {
	// We use the STOP instruction whether to see
	// the jump table was initialised. If it was not
	// we'll set the default jump table.
	if cfg.JumpTable[STOP] == nil {
		var jt JumpTable
		// switch {
		// case evm.chainRules.IsLondon:
		// 	jt = londonInstructionSet
		// case evm.chainRules.IsBerlin:
		// 	jt = berlinInstructionSet
		// case evm.chainRules.IsIstanbul:
		// 	jt = istanbulInstructionSet
		// case evm.chainRules.IsConstantinople:
		// 	jt = constantinopleInstructionSet
		// case evm.chainRules.IsByzantium:
		// 	jt = byzantiumInstructionSet
		// case evm.chainRules.IsEIP158:
		// 	jt = spuriousDragonInstructionSet
		// case evm.chainRules.IsEIP150:
		// 	jt = tangerineWhistleInstructionSet
		// case evm.chainRules.IsHomestead:
		// 	jt = homesteadInstructionSet
		// default:
		// 	jt = frontierInstructionSet
		// }
		jt = londonInstructionSet
		for i, eip := range cfg.ExtraEips {
			if err := EnableEIP(eip, &jt); err != nil {
				// Disable it, so caller can check if it's activated or not
				cfg.ExtraEips = append(cfg.ExtraEips[:i], cfg.ExtraEips[i+1:]...)
				log.Error("EIP activation failed", "eip", eip, "error", err)
			}
		}
		cfg.JumpTable = jt
	}

	return &EVMInterpreter{
		evm: evm,
		cfg: cfg,
	}
}

func (in *EVMInterpreter) Run(contract *Contract, input []byte, readOnly bool) (ret []byte, err error) {

	if len(contract.Code) == 0 {
		return nil, nil
	}
	var (
		op          OpCode        // current opcode 当前操作码
		mem         = NewMemory() // bound memory 内存
		stack       = newstack()  // local stack 栈
		callContext = &ScopeContext{
			Memory:   mem,
			Stack:    stack,
			Contract: contract,
		}
		// For optimisation reason we're using uint64 as the program counter.
		// It's theoretically possible to go above 2^64. The YP defines the PC
		// to be uint256. Practically much less so feasible.
		pc = uint64(0) // program counter
		// cost uint64      // gas花费
		// // copies used by tracer
		// // debug使用
		// pcCopy uint64 // needed for the deferred Tracer
		// // debug使用
		// gasCopy uint64 // for Tracer to log gas remaining before execution
		// // debug使用
		// logged bool // deferred Tracer should ignore already logged steps
		// // 当前操作码执行函数的返回值
		res []byte // result of the opcode execution function
	)
	contract.Input = input
	codeLen := uint64(len(contract.Code))

	scanner := bufio.NewScanner(os.Stdin)
	for {
		scanner.Scan()
		if pc >= codeLen {
			break
		}
		fmt.Printf("##pc==>%04x\n", pc)
		op = contract.GetOp(pc)
		operation := in.cfg.JumpTable[op]
		fmt.Println(op)
		if operation == nil {
			return nil, &ErrInvalidOpCode{opcode: op}
		}

		// mem.Resize(1024)

		var memorySize uint64
		if operation.memorySize != nil {
			// stack.Print()
			memSize, overflow := operation.memorySize(stack)
			if overflow {
				return nil, ErrGasUintOverflow
			}
			// memory is expanded in words of 32 bytes. Gas
			// is also calculated in words.
			if memorySize, overflow = math.SafeMul(toWordSize(memSize), 32); overflow {
				return nil, ErrGasUintOverflow
			}
		}
		if memorySize > 0 {
			fmt.Println("##memorySize: ", memorySize)
			mem.Resize(memorySize)
		}

		res, err = operation.execute(&pc, in, callContext)
		stack.PrintReverse()
		if mem.Len() > 0 {
			mem.Print()
		}
		in.evm.StateDB.PrintAccount(contract.Address())
		if operation.returns {
			in.returnData = res
		}

		switch {
		case err != nil: // 报错
			return nil, err
		case operation.reverts: // 出错回滚
			return res, ErrExecutionReverted
		case operation.halts: // 停止
			return res, nil
		case !operation.jumps: // 跳转
			pc++
		}
	}

	return nil, nil

}

func (in *EVMInterpreter) Disassembler(contractCode string) {
	// byteCodeStr := "6080604052348015600f57600080fd5b506004361060285760003560e01c8063ef5fb05b14602d575b600080fd5b604080518082018252600a8152691a195b1b1bdddbdc9b1960b21b60208201529051605791906060565b60405180910390f35b600060208083528351808285015260005b81811015608b578581018301518582016040015282016071565b506000604082860101526040601f19601f830116850101925050509291505056fea2646970667358221220ed09cf55ba7895ca4cb49da8268e332340bcb8de2af198510c97fd3a2a06f36864736f6c63430008110033"
	byteCodes, err := hex.DecodeString(contractCode)
	if err != nil {
		fmt.Println("转换失败:", err)
		return
	}
	pc := uint64(0)
	caller := common.HexToAddress("0x0000000000000000000000000000000000000000")
	contract := NewContract(AccountRef(caller), AccountRef(caller), big.NewInt(0), 0)
	contractAddr := crypto.CreateAddress(AccountRef(caller).Address(), 0)
	contract.SetCodeOptionalHash(&contractAddr, &codeAndHash{code: byteCodes})
	codeLen := uint64(len(contract.Code))
	for {
		if pc >= codeLen {
			break
		}
		sPc := pc
		opCode := contract.GetOp(pc)
		opCodeStr := opCodeToString[opCode]
		if opCodeStr == "" {
			opCodeStr = "INVALID"
		}
		holdiInstruction := fmt.Sprintf("%04x    %s", sPc, opCodeStr)

		skipPc := opCodeInfoList[opCodeStr].opCodeCount
		if skipPc > 0 {
			opCodeByteArray := []byte{}
			for skipPc > 0 {
				pc += 1
				skipPc -= 1
				opCodeByte := contract.GetByte(pc)
				opCodeByteArray = append(opCodeByteArray, opCodeByte)
			}
			holdiInstruction = fmt.Sprintf("%04x    %-20s 0x%x", sPc, opCodeStr, opCodeByteArray)
		}

		fmt.Println(holdiInstruction)
		pc += 1
	}
	fmt.Printf("over ...\n")
}
