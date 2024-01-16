package opcode

type Opcode byte

func (o Opcode) ToByte() byte {
	return byte(o)
}

func (o Opcode) String() string {
	return Names[o]
}

/*
When adding a new opcode, make sure to check and make any needed changes to the following places:

- Add the opcode to the appropiate arg count map below.
- Add a string representation of the opcode below.
- If the opcode changes the stack in any way, edit the calculateStackSize function in the compiler.
- If the opcode changes the block stack in any way, edit the calculateBlockSize function in the compiler.
- If the opcode takes any arguments, edit the compiler.CodeBlock.Print() method to print the correct output.
- And obviously, implement it in the virtual machine.
*/
const (
	Noop Opcode = iota
	LoadConst
	StoreConst
	LoadFast
	StoreFast
	DeleteFast
	Define
	LoadGlobal
	StoreGlobal
	LoadIndex
	StoreIndex
	LoadAttribute
	StoreAttribute
	BinaryAdd
	BinarySub
	BinaryMul
	BinaryDivide
	BinaryMod
	BinaryShiftL
	BinaryShiftR
	BinaryAnd
	BinaryOr
	BinaryNot
	BinaryAndNot
	UnaryNeg
	UnaryNot
	Implements
	Compare
	Call
	Return
	Pop
	MakeArray
	MakeMap
	MakeFunction
	PopJumpIfTrue
	PopJumpIfFalse
	JumpIfTrueOrPop
	JumpIfFalseOrPop
	JumpAbsolute
	JumpForward
	StartBlock
	EndBlock
	StartLoop
	Continue
	NextIter
	Break
	Recover
	BuildClass
	MakeInstance
	Import
	Dup
	GetIter

	MaxOpcode // Not a real opcode, just used to denote the maximum value of a valid opcode
	Label
)

const (
	CmpEq byte = iota
	CmpNotEq
	CmpLT
	CmpGT
	CmpLTEq
	CmpGTEq
	MaxCmpCodes
)

// 2 16-bit arguments
var HasFourByteArg = map[Opcode]bool{
	StartLoop: true,
}

// 1 16-bit argument
var HasTwoByteArg = map[Opcode]bool{
	LoadConst:        true,
	StoreConst:       true,
	LoadFast:         true,
	StoreFast:        true,
	DeleteFast:       true,
	Define:           true,
	LoadGlobal:       true,
	StoreGlobal:      true,
	LoadAttribute:    true,
	StoreAttribute:   true,
	Call:             true,
	MakeArray:        true,
	MakeMap:          true,
	PopJumpIfTrue:    true,
	PopJumpIfFalse:   true,
	JumpIfTrueOrPop:  true,
	JumpIfFalseOrPop: true,
	JumpAbsolute:     true,
	JumpForward:      true,
	Recover:          true,
	BuildClass:       true,
	MakeInstance:     true,
	Import:           true,
}

// 1 8-bit argument
var HasOneByteArg = map[Opcode]bool{
	Compare: true,
}

var HasNoArg = map[Opcode]bool{
	Noop:         true,
	LoadIndex:    true,
	StoreIndex:   true,
	BinaryAdd:    true,
	BinarySub:    true,
	BinaryMul:    true,
	BinaryDivide: true,
	BinaryMod:    true,
	BinaryShiftL: true,
	BinaryShiftR: true,
	BinaryAnd:    true,
	BinaryOr:     true,
	BinaryNot:    true,
	BinaryAndNot: true,
	Implements:   true,
	UnaryNeg:     true,
	UnaryNot:     true,
	Return:       true,
	Pop:          true,
	MakeFunction: true,
	StartBlock:   true,
	EndBlock:     true,
	Continue:     true,
	NextIter:     true,
	Break:        true,
	Dup:          true,
	GetIter:      true,
}

var Names = map[Opcode]string{
	Noop:             "NOOP",
	LoadConst:        "LOAD_CONST",
	StoreConst:       "STORE_CONST",
	LoadFast:         "LOAD_FAST",
	StoreFast:        "STORE_FAST",
	DeleteFast:       "DELETE_FAST",
	Define:           "DEFINE",
	LoadGlobal:       "LOAD_GLOBAL",
	StoreGlobal:      "STORE_GLOBAL",
	LoadIndex:        "LOAD_INDEX",
	StoreIndex:       "STORE_INDEX",
	LoadAttribute:    "LOAD_ATTRIBUTE",
	StoreAttribute:   "STORE_ATTRIBUTE",
	BinaryAdd:        "BINARY_ADD",
	BinarySub:        "BINARY_SUB",
	BinaryMul:        "BINARY_MUL",
	BinaryDivide:     "BINARY_DIVIDE",
	BinaryMod:        "BINARY_MOD",
	BinaryShiftL:     "BINARY_SHIFTL",
	BinaryShiftR:     "BINARY_SHIFTR",
	BinaryAnd:        "BINARY_AND",
	BinaryOr:         "BINARY_OR",
	BinaryNot:        "BINARY_NOT",
	BinaryAndNot:     "BINARY_ANDNOT",
	Implements:       "IMPLEMENTS",
	UnaryNeg:         "UNARY_NEG",
	UnaryNot:         "UNARY_NOT",
	Compare:          "COMPARE",
	Call:             "CALL",
	Return:           "RETURN",
	Pop:              "POP",
	MakeArray:        "MAKE_ARRAY",
	MakeMap:          "MAKE_MAP",
	MakeFunction:     "MAKE_FUNCTION",
	PopJumpIfTrue:    "POP_JUMP_IF_TRUE",
	PopJumpIfFalse:   "POP_JUMP_IF_FALSE",
	JumpIfTrueOrPop:  "JUMP_IF_TRUE_OR_POP",
	JumpIfFalseOrPop: "JUMP_IF_FALSE_OR_POP",
	JumpAbsolute:     "JUMP_ABSOLUTE",
	JumpForward:      "JUMP_FORWARD",
	StartBlock:       "START_BLOCK",
	EndBlock:         "END_BLOCK",
	StartLoop:        "START_LOOP",
	Continue:         "CONTINUE",
	NextIter:         "NEXT_ITER",
	Break:            "BREAK",
	Recover:          "RECOVER",
	BuildClass:       "BUILD_CLASS",
	MakeInstance:     "MAKE_INSTANCE",
	Import:           "IMPORT",
	Dup:              "DUP",
	GetIter:          "GET_ITER",
}

var CmpOps = map[byte]string{
	CmpEq:    "==",
	CmpNotEq: "!=",
	CmpLT:    "<",
	CmpGT:    ">",
	CmpLTEq:  "<=",
	CmpGTEq:  ">=",
}
