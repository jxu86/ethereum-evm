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
	"fmt"
)

// OpCode is an EVM opcode
type OpCode byte

// IsPush specifies if an opcode is a PUSH opcode.
func (op OpCode) IsPush() bool {
	switch op {
	case PUSH1, PUSH2, PUSH3, PUSH4, PUSH5, PUSH6, PUSH7, PUSH8, PUSH9, PUSH10, PUSH11, PUSH12, PUSH13, PUSH14, PUSH15, PUSH16, PUSH17, PUSH18, PUSH19, PUSH20, PUSH21, PUSH22, PUSH23, PUSH24, PUSH25, PUSH26, PUSH27, PUSH28, PUSH29, PUSH30, PUSH31, PUSH32:
		return true
	}
	return false
}

// IsStaticJump specifies if an opcode is JUMP.
func (op OpCode) IsStaticJump() bool {
	return op == JUMP
}

// 0x0 range - arithmetic ops.
const (
	STOP OpCode = iota
	ADD
	MUL
	SUB
	DIV
	SDIV
	MOD
	SMOD
	ADDMOD
	MULMOD
	EXP
	SIGNEXTEND
)

// 0x10 range - comparison ops.
const (
	LT OpCode = iota + 0x10
	GT
	SLT
	SGT
	EQ
	ISZERO
	AND
	OR
	XOR
	NOT
	BYTE
	SHL
	SHR
	SAR

	SHA3 OpCode = 0x20
)

// 0x30 range - closure state.
const (
	ADDRESS OpCode = 0x30 + iota
	BALANCE
	ORIGIN
	CALLER
	CALLVALUE
	CALLDATALOAD
	CALLDATASIZE
	CALLDATACOPY
	CODESIZE
	CODECOPY
	GASPRICE
	EXTCODESIZE
	EXTCODECOPY
	RETURNDATASIZE
	RETURNDATACOPY
	EXTCODEHASH
)

// 0x40 range - block operations.
const (
	BLOCKHASH OpCode = 0x40 + iota
	COINBASE
	TIMESTAMP
	NUMBER
	DIFFICULTY
	GASLIMIT
	CHAINID     OpCode = 0x46
	SELFBALANCE OpCode = 0x47
	BASEFEE     OpCode = 0x48
)

// 0x50 range - 'storage' and execution.
const (
	POP      OpCode = 0x50
	MLOAD    OpCode = 0x51
	MSTORE   OpCode = 0x52
	MSTORE8  OpCode = 0x53
	SLOAD    OpCode = 0x54
	SSTORE   OpCode = 0x55
	JUMP     OpCode = 0x56
	JUMPI    OpCode = 0x57
	PC       OpCode = 0x58
	MSIZE    OpCode = 0x59
	GAS      OpCode = 0x5a
	JUMPDEST OpCode = 0x5b
)

// 0x60 range.
const (
	PUSH1 OpCode = 0x60 + iota
	PUSH2
	PUSH3
	PUSH4
	PUSH5
	PUSH6
	PUSH7
	PUSH8
	PUSH9
	PUSH10
	PUSH11
	PUSH12
	PUSH13
	PUSH14
	PUSH15
	PUSH16
	PUSH17
	PUSH18
	PUSH19
	PUSH20
	PUSH21
	PUSH22
	PUSH23
	PUSH24
	PUSH25
	PUSH26
	PUSH27
	PUSH28
	PUSH29
	PUSH30
	PUSH31
	PUSH32
	DUP1
	DUP2
	DUP3
	DUP4
	DUP5
	DUP6
	DUP7
	DUP8
	DUP9
	DUP10
	DUP11
	DUP12
	DUP13
	DUP14
	DUP15
	DUP16
	SWAP1
	SWAP2
	SWAP3
	SWAP4
	SWAP5
	SWAP6
	SWAP7
	SWAP8
	SWAP9
	SWAP10
	SWAP11
	SWAP12
	SWAP13
	SWAP14
	SWAP15
	SWAP16
)

// 0xa0 range - logging ops.
const (
	LOG0 OpCode = 0xa0 + iota
	LOG1
	LOG2
	LOG3
	LOG4
)

// unofficial opcodes used for parsing.
const (
	PUSH OpCode = 0xb0 + iota
	DUP
	SWAP
)

// 0xf0 range - closures.
const (
	CREATE OpCode = 0xf0 + iota
	CALL
	CALLCODE
	RETURN
	DELEGATECALL
	CREATE2
	STATICCALL   OpCode = 0xfa
	REVERT       OpCode = 0xfd
	SELFDESTRUCT OpCode = 0xff
)

// Since the opcodes aren't all in order we can't use a regular slice.
var opCodeToString = map[OpCode]string{
	// 0x0 range - arithmetic ops.
	STOP:       "STOP",
	ADD:        "ADD",
	MUL:        "MUL",
	SUB:        "SUB",
	DIV:        "DIV",
	SDIV:       "SDIV",
	MOD:        "MOD",
	SMOD:       "SMOD",
	EXP:        "EXP",
	NOT:        "NOT",
	LT:         "LT",
	GT:         "GT",
	SLT:        "SLT",
	SGT:        "SGT",
	EQ:         "EQ",
	ISZERO:     "ISZERO",
	SIGNEXTEND: "SIGNEXTEND",

	// 0x10 range - bit ops.
	AND:    "AND",
	OR:     "OR",
	XOR:    "XOR",
	BYTE:   "BYTE",
	SHL:    "SHL",
	SHR:    "SHR",
	SAR:    "SAR",
	ADDMOD: "ADDMOD",
	MULMOD: "MULMOD",

	// 0x20 range - crypto.
	SHA3: "SHA3",

	// 0x30 range - closure state.
	ADDRESS:        "ADDRESS",
	BALANCE:        "BALANCE",
	ORIGIN:         "ORIGIN",
	CALLER:         "CALLER",
	CALLVALUE:      "CALLVALUE",
	CALLDATALOAD:   "CALLDATALOAD",
	CALLDATASIZE:   "CALLDATASIZE",
	CALLDATACOPY:   "CALLDATACOPY",
	CODESIZE:       "CODESIZE",
	CODECOPY:       "CODECOPY",
	GASPRICE:       "GASPRICE",
	EXTCODESIZE:    "EXTCODESIZE",
	EXTCODECOPY:    "EXTCODECOPY",
	RETURNDATASIZE: "RETURNDATASIZE",
	RETURNDATACOPY: "RETURNDATACOPY",
	EXTCODEHASH:    "EXTCODEHASH",

	// 0x40 range - block operations.
	BLOCKHASH:   "BLOCKHASH",
	COINBASE:    "COINBASE",
	TIMESTAMP:   "TIMESTAMP",
	NUMBER:      "NUMBER",
	DIFFICULTY:  "DIFFICULTY",
	GASLIMIT:    "GASLIMIT",
	CHAINID:     "CHAINID",
	SELFBALANCE: "SELFBALANCE",
	BASEFEE:     "BASEFEE",

	// 0x50 range - 'storage' and execution.
	POP: "POP",
	//DUP:     "DUP",
	//SWAP:    "SWAP",
	MLOAD:    "MLOAD",
	MSTORE:   "MSTORE",
	MSTORE8:  "MSTORE8",
	SLOAD:    "SLOAD",
	SSTORE:   "SSTORE",
	JUMP:     "JUMP",
	JUMPI:    "JUMPI",
	PC:       "PC",
	MSIZE:    "MSIZE",
	GAS:      "GAS",
	JUMPDEST: "JUMPDEST",

	// 0x60 range - push.
	PUSH1:  "PUSH1",
	PUSH2:  "PUSH2",
	PUSH3:  "PUSH3",
	PUSH4:  "PUSH4",
	PUSH5:  "PUSH5",
	PUSH6:  "PUSH6",
	PUSH7:  "PUSH7",
	PUSH8:  "PUSH8",
	PUSH9:  "PUSH9",
	PUSH10: "PUSH10",
	PUSH11: "PUSH11",
	PUSH12: "PUSH12",
	PUSH13: "PUSH13",
	PUSH14: "PUSH14",
	PUSH15: "PUSH15",
	PUSH16: "PUSH16",
	PUSH17: "PUSH17",
	PUSH18: "PUSH18",
	PUSH19: "PUSH19",
	PUSH20: "PUSH20",
	PUSH21: "PUSH21",
	PUSH22: "PUSH22",
	PUSH23: "PUSH23",
	PUSH24: "PUSH24",
	PUSH25: "PUSH25",
	PUSH26: "PUSH26",
	PUSH27: "PUSH27",
	PUSH28: "PUSH28",
	PUSH29: "PUSH29",
	PUSH30: "PUSH30",
	PUSH31: "PUSH31",
	PUSH32: "PUSH32",

	DUP1:  "DUP1",
	DUP2:  "DUP2",
	DUP3:  "DUP3",
	DUP4:  "DUP4",
	DUP5:  "DUP5",
	DUP6:  "DUP6",
	DUP7:  "DUP7",
	DUP8:  "DUP8",
	DUP9:  "DUP9",
	DUP10: "DUP10",
	DUP11: "DUP11",
	DUP12: "DUP12",
	DUP13: "DUP13",
	DUP14: "DUP14",
	DUP15: "DUP15",
	DUP16: "DUP16",

	SWAP1:  "SWAP1",
	SWAP2:  "SWAP2",
	SWAP3:  "SWAP3",
	SWAP4:  "SWAP4",
	SWAP5:  "SWAP5",
	SWAP6:  "SWAP6",
	SWAP7:  "SWAP7",
	SWAP8:  "SWAP8",
	SWAP9:  "SWAP9",
	SWAP10: "SWAP10",
	SWAP11: "SWAP11",
	SWAP12: "SWAP12",
	SWAP13: "SWAP13",
	SWAP14: "SWAP14",
	SWAP15: "SWAP15",
	SWAP16: "SWAP16",
	LOG0:   "LOG0",
	LOG1:   "LOG1",
	LOG2:   "LOG2",
	LOG3:   "LOG3",
	LOG4:   "LOG4",

	// 0xf0 range.
	CREATE:       "CREATE",
	CALL:         "CALL",
	RETURN:       "RETURN",
	CALLCODE:     "CALLCODE",
	DELEGATECALL: "DELEGATECALL",
	CREATE2:      "CREATE2",
	STATICCALL:   "STATICCALL",
	REVERT:       "REVERT",
	SELFDESTRUCT: "SELFDESTRUCT",

	PUSH: "PUSH",
	DUP:  "DUP",
	SWAP: "SWAP",
}

func (op OpCode) String() string {
	str := opCodeToString[op]
	if len(str) == 0 {
		return fmt.Sprintf("opcode 0x%x not defined", int(op))
	}

	return str
}

var stringToOp = map[string]OpCode{
	"STOP":           STOP,
	"ADD":            ADD,
	"MUL":            MUL,
	"SUB":            SUB,
	"DIV":            DIV,
	"SDIV":           SDIV,
	"MOD":            MOD,
	"SMOD":           SMOD,
	"EXP":            EXP,
	"NOT":            NOT,
	"LT":             LT,
	"GT":             GT,
	"SLT":            SLT,
	"SGT":            SGT,
	"EQ":             EQ,
	"ISZERO":         ISZERO,
	"SIGNEXTEND":     SIGNEXTEND,
	"AND":            AND,
	"OR":             OR,
	"XOR":            XOR,
	"BYTE":           BYTE,
	"SHL":            SHL,
	"SHR":            SHR,
	"SAR":            SAR,
	"ADDMOD":         ADDMOD,
	"MULMOD":         MULMOD,
	"SHA3":           SHA3,
	"ADDRESS":        ADDRESS,
	"BALANCE":        BALANCE,
	"ORIGIN":         ORIGIN,
	"CALLER":         CALLER,
	"CALLVALUE":      CALLVALUE,
	"CALLDATALOAD":   CALLDATALOAD,
	"CALLDATASIZE":   CALLDATASIZE,
	"CALLDATACOPY":   CALLDATACOPY,
	"CHAINID":        CHAINID,
	"BASEFEE":        BASEFEE,
	"DELEGATECALL":   DELEGATECALL,
	"STATICCALL":     STATICCALL,
	"CODESIZE":       CODESIZE,
	"CODECOPY":       CODECOPY,
	"GASPRICE":       GASPRICE,
	"EXTCODESIZE":    EXTCODESIZE,
	"EXTCODECOPY":    EXTCODECOPY,
	"RETURNDATASIZE": RETURNDATASIZE,
	"RETURNDATACOPY": RETURNDATACOPY,
	"EXTCODEHASH":    EXTCODEHASH,
	"BLOCKHASH":      BLOCKHASH,
	"COINBASE":       COINBASE,
	"TIMESTAMP":      TIMESTAMP,
	"NUMBER":         NUMBER,
	"DIFFICULTY":     DIFFICULTY,
	"GASLIMIT":       GASLIMIT,
	"SELFBALANCE":    SELFBALANCE,
	"POP":            POP,
	"MLOAD":          MLOAD,
	"MSTORE":         MSTORE,
	"MSTORE8":        MSTORE8,
	"SLOAD":          SLOAD,
	"SSTORE":         SSTORE,
	"JUMP":           JUMP,
	"JUMPI":          JUMPI,
	"PC":             PC,
	"MSIZE":          MSIZE,
	"GAS":            GAS,
	"JUMPDEST":       JUMPDEST,
	"PUSH1":          PUSH1,
	"PUSH2":          PUSH2,
	"PUSH3":          PUSH3,
	"PUSH4":          PUSH4,
	"PUSH5":          PUSH5,
	"PUSH6":          PUSH6,
	"PUSH7":          PUSH7,
	"PUSH8":          PUSH8,
	"PUSH9":          PUSH9,
	"PUSH10":         PUSH10,
	"PUSH11":         PUSH11,
	"PUSH12":         PUSH12,
	"PUSH13":         PUSH13,
	"PUSH14":         PUSH14,
	"PUSH15":         PUSH15,
	"PUSH16":         PUSH16,
	"PUSH17":         PUSH17,
	"PUSH18":         PUSH18,
	"PUSH19":         PUSH19,
	"PUSH20":         PUSH20,
	"PUSH21":         PUSH21,
	"PUSH22":         PUSH22,
	"PUSH23":         PUSH23,
	"PUSH24":         PUSH24,
	"PUSH25":         PUSH25,
	"PUSH26":         PUSH26,
	"PUSH27":         PUSH27,
	"PUSH28":         PUSH28,
	"PUSH29":         PUSH29,
	"PUSH30":         PUSH30,
	"PUSH31":         PUSH31,
	"PUSH32":         PUSH32,
	"DUP1":           DUP1,
	"DUP2":           DUP2,
	"DUP3":           DUP3,
	"DUP4":           DUP4,
	"DUP5":           DUP5,
	"DUP6":           DUP6,
	"DUP7":           DUP7,
	"DUP8":           DUP8,
	"DUP9":           DUP9,
	"DUP10":          DUP10,
	"DUP11":          DUP11,
	"DUP12":          DUP12,
	"DUP13":          DUP13,
	"DUP14":          DUP14,
	"DUP15":          DUP15,
	"DUP16":          DUP16,
	"SWAP1":          SWAP1,
	"SWAP2":          SWAP2,
	"SWAP3":          SWAP3,
	"SWAP4":          SWAP4,
	"SWAP5":          SWAP5,
	"SWAP6":          SWAP6,
	"SWAP7":          SWAP7,
	"SWAP8":          SWAP8,
	"SWAP9":          SWAP9,
	"SWAP10":         SWAP10,
	"SWAP11":         SWAP11,
	"SWAP12":         SWAP12,
	"SWAP13":         SWAP13,
	"SWAP14":         SWAP14,
	"SWAP15":         SWAP15,
	"SWAP16":         SWAP16,
	"LOG0":           LOG0,
	"LOG1":           LOG1,
	"LOG2":           LOG2,
	"LOG3":           LOG3,
	"LOG4":           LOG4,
	"CREATE":         CREATE,
	"CREATE2":        CREATE2,
	"CALL":           CALL,
	"RETURN":         RETURN,
	"CALLCODE":       CALLCODE,
	"REVERT":         REVERT,
	"SELFDESTRUCT":   SELFDESTRUCT,
}

// StringToOp finds the opcode whose name is stored in `str`.
func StringToOp(str string) OpCode {
	return stringToOp[str]
}

type OpCodeInfo struct {
	opCodeCount      uint8
	stackInputCount  uint8
	stackOutputCount uint8
	gasCount         uint8
}

var opCodeInfoList = map[string]OpCodeInfo{
	"STOP":           OpCodeInfo{0, 0, 0, 0},
	"ADD":            OpCodeInfo{0, 0, 0, 0},
	"MUL":            OpCodeInfo{0, 0, 0, 0},
	"SUB":            OpCodeInfo{0, 0, 0, 0},
	"DIV":            OpCodeInfo{0, 0, 0, 0},
	"SDIV":           OpCodeInfo{0, 0, 0, 0},
	"MOD":            OpCodeInfo{0, 0, 0, 0},
	"SMOD":           OpCodeInfo{0, 0, 0, 0},
	"EXP":            OpCodeInfo{0, 0, 0, 0},
	"NOT":            OpCodeInfo{0, 0, 0, 0},
	"LT":             OpCodeInfo{0, 0, 0, 0},
	"GT":             OpCodeInfo{0, 0, 0, 0},
	"SLT":            OpCodeInfo{0, 0, 0, 0},
	"SGT":            OpCodeInfo{0, 0, 0, 0},
	"EQ":             OpCodeInfo{0, 0, 0, 0},
	"ISZERO":         OpCodeInfo{0, 0, 0, 0},
	"SIGNEXTEND":     OpCodeInfo{0, 0, 0, 0},
	"AND":            OpCodeInfo{0, 0, 0, 0},
	"OR":             OpCodeInfo{0, 0, 0, 0},
	"XOR":            OpCodeInfo{0, 0, 0, 0},
	"BYTE":           OpCodeInfo{0, 0, 0, 0},
	"SHL":            OpCodeInfo{0, 0, 0, 0},
	"SHR":            OpCodeInfo{0, 0, 0, 0},
	"SAR":            OpCodeInfo{0, 0, 0, 0},
	"ADDMOD":         OpCodeInfo{0, 0, 0, 0},
	"MULMOD":         OpCodeInfo{0, 0, 0, 0},
	"SHA3":           OpCodeInfo{0, 0, 0, 0},
	"ADDRESS":        OpCodeInfo{0, 0, 0, 0},
	"BALANCE":        OpCodeInfo{0, 0, 0, 0},
	"ORIGIN":         OpCodeInfo{0, 0, 0, 0},
	"CALLER":         OpCodeInfo{0, 0, 0, 0},
	"CALLVALUE":      OpCodeInfo{0, 0, 0, 0},
	"CALLDATALOAD":   OpCodeInfo{0, 0, 0, 0},
	"CALLDATASIZE":   OpCodeInfo{0, 0, 0, 0},
	"CALLDATACOPY":   OpCodeInfo{0, 0, 0, 0},
	"CHAINID":        OpCodeInfo{0, 0, 0, 0},
	"BASEFEE":        OpCodeInfo{0, 0, 0, 0},
	"DELEGATECALL":   OpCodeInfo{0, 0, 0, 0},
	"STATICCALL":     OpCodeInfo{0, 0, 0, 0},
	"CODESIZE":       OpCodeInfo{0, 0, 0, 0},
	"CODECOPY":       OpCodeInfo{0, 0, 0, 0},
	"GASPRICE":       OpCodeInfo{0, 0, 0, 0},
	"EXTCODESIZE":    OpCodeInfo{0, 0, 0, 0},
	"EXTCODECOPY":    OpCodeInfo{0, 0, 0, 0},
	"RETURNDATASIZE": OpCodeInfo{0, 0, 0, 0},
	"RETURNDATACOPY": OpCodeInfo{0, 0, 0, 0},
	"EXTCODEHASH":    OpCodeInfo{0, 0, 0, 0},
	"BLOCKHASH":      OpCodeInfo{0, 0, 0, 0},
	"COINBASE":       OpCodeInfo{0, 0, 0, 0},
	"TIMESTAMP":      OpCodeInfo{0, 0, 0, 0},
	"NUMBER":         OpCodeInfo{0, 0, 0, 0},
	"DIFFICULTY":     OpCodeInfo{0, 0, 0, 0},
	"GASLIMIT":       OpCodeInfo{0, 0, 0, 0},
	"SELFBALANCE":    OpCodeInfo{0, 0, 0, 0},
	"POP":            OpCodeInfo{0, 0, 0, 0},
	"MLOAD":          OpCodeInfo{0, 0, 0, 0},
	"MSTORE":         OpCodeInfo{0, 0, 0, 0},
	"MSTORE8":        OpCodeInfo{0, 0, 0, 0},
	"SLOAD":          OpCodeInfo{0, 0, 0, 0},
	"SSTORE":         OpCodeInfo{0, 0, 0, 0},
	"JUMP":           OpCodeInfo{0, 0, 0, 0},
	"JUMPI":          OpCodeInfo{0, 0, 0, 0},
	"PC":             OpCodeInfo{0, 0, 0, 0},
	"MSIZE":          OpCodeInfo{0, 0, 0, 0},
	"GAS":            OpCodeInfo{0, 0, 0, 0},
	"JUMPDEST":       OpCodeInfo{0, 0, 0, 0},
	"PUSH1":          OpCodeInfo{1, 0, 0, 0},
	"PUSH2":          OpCodeInfo{2, 0, 0, 0},
	"PUSH3":          OpCodeInfo{3, 0, 0, 0},
	"PUSH4":          OpCodeInfo{4, 0, 0, 0},
	"PUSH5":          OpCodeInfo{5, 0, 0, 0},
	"PUSH6":          OpCodeInfo{6, 0, 0, 0},
	"PUSH7":          OpCodeInfo{7, 0, 0, 0},
	"PUSH8":          OpCodeInfo{8, 0, 0, 0},
	"PUSH9":          OpCodeInfo{9, 0, 0, 0},
	"PUSH10":         OpCodeInfo{10, 0, 0, 0},
	"PUSH11":         OpCodeInfo{11, 0, 0, 0},
	"PUSH12":         OpCodeInfo{12, 0, 0, 0},
	"PUSH13":         OpCodeInfo{13, 0, 0, 0},
	"PUSH14":         OpCodeInfo{14, 0, 0, 0},
	"PUSH15":         OpCodeInfo{15, 0, 0, 0},
	"PUSH16":         OpCodeInfo{16, 0, 0, 0},
	"PUSH17":         OpCodeInfo{17, 0, 0, 0},
	"PUSH18":         OpCodeInfo{18, 0, 0, 0},
	"PUSH19":         OpCodeInfo{19, 0, 0, 0},
	"PUSH20":         OpCodeInfo{20, 0, 0, 0},
	"PUSH21":         OpCodeInfo{21, 0, 0, 0},
	"PUSH22":         OpCodeInfo{22, 0, 0, 0},
	"PUSH23":         OpCodeInfo{23, 0, 0, 0},
	"PUSH24":         OpCodeInfo{24, 0, 0, 0},
	"PUSH25":         OpCodeInfo{25, 0, 0, 0},
	"PUSH26":         OpCodeInfo{26, 0, 0, 0},
	"PUSH27":         OpCodeInfo{27, 0, 0, 0},
	"PUSH28":         OpCodeInfo{28, 0, 0, 0},
	"PUSH29":         OpCodeInfo{29, 0, 0, 0},
	"PUSH30":         OpCodeInfo{30, 0, 0, 0},
	"PUSH31":         OpCodeInfo{31, 0, 0, 0},
	"PUSH32":         OpCodeInfo{32, 0, 0, 0},
	"DUP1":           OpCodeInfo{0, 0, 0, 0},
	"DUP2":           OpCodeInfo{0, 0, 0, 0},
	"DUP3":           OpCodeInfo{0, 0, 0, 0},
	"DUP4":           OpCodeInfo{0, 0, 0, 0},
	"DUP5":           OpCodeInfo{0, 0, 0, 0},
	"DUP6":           OpCodeInfo{0, 0, 0, 0},
	"DUP7":           OpCodeInfo{0, 0, 0, 0},
	"DUP8":           OpCodeInfo{0, 0, 0, 0},
	"DUP9":           OpCodeInfo{0, 0, 0, 0},
	"DUP10":          OpCodeInfo{0, 0, 0, 0},
	"DUP11":          OpCodeInfo{0, 0, 0, 0},
	"DUP12":          OpCodeInfo{0, 0, 0, 0},
	"DUP13":          OpCodeInfo{0, 0, 0, 0},
	"DUP14":          OpCodeInfo{0, 0, 0, 0},
	"DUP15":          OpCodeInfo{0, 0, 0, 0},
	"DUP16":          OpCodeInfo{0, 0, 0, 0},
	"SWAP1":          OpCodeInfo{0, 0, 0, 0},
	"SWAP2":          OpCodeInfo{0, 0, 0, 0},
	"SWAP3":          OpCodeInfo{0, 0, 0, 0},
	"SWAP4":          OpCodeInfo{0, 0, 0, 0},
	"SWAP5":          OpCodeInfo{0, 0, 0, 0},
	"SWAP6":          OpCodeInfo{0, 0, 0, 0},
	"SWAP7":          OpCodeInfo{0, 0, 0, 0},
	"SWAP8":          OpCodeInfo{0, 0, 0, 0},
	"SWAP9":          OpCodeInfo{0, 0, 0, 0},
	"SWAP10":         OpCodeInfo{0, 0, 0, 0},
	"SWAP11":         OpCodeInfo{0, 0, 0, 0},
	"SWAP12":         OpCodeInfo{0, 0, 0, 0},
	"SWAP13":         OpCodeInfo{0, 0, 0, 0},
	"SWAP14":         OpCodeInfo{0, 0, 0, 0},
	"SWAP15":         OpCodeInfo{0, 0, 0, 0},
	"SWAP16":         OpCodeInfo{0, 0, 0, 0},
	"LOG0":           OpCodeInfo{0, 0, 0, 0},
	"LOG1":           OpCodeInfo{0, 0, 0, 0},
	"LOG2":           OpCodeInfo{0, 0, 0, 0},
	"LOG3":           OpCodeInfo{0, 0, 0, 0},
	"LOG4":           OpCodeInfo{0, 0, 0, 0},
	"CREATE":         OpCodeInfo{0, 0, 0, 0},
	"CREATE2":        OpCodeInfo{0, 0, 0, 0},
	"CALL":           OpCodeInfo{0, 0, 0, 0},
	"RETURN":         OpCodeInfo{0, 0, 0, 0},
	"CALLCODE":       OpCodeInfo{0, 0, 0, 0},
	"REVERT":         OpCodeInfo{0, 0, 0, 0},
	"SELFDESTRUCT":   OpCodeInfo{0, 0, 0, 0},
	"INVALID":        OpCodeInfo{0, 0, 0, 0},
}
