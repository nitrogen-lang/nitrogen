package opcode

const (
	Noop byte = iota
	LoadConst
	LoadFast
	StoreFast
	LoadGlobal
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
	Compare
	Call
	Return
	Pop
	MakeArray
	MakeMap
	PopJumpIfTrue
	PopJumpIfFalse
	JumpIfTrueOrPop
	JumpIfFalseOrPop
	JumpAbsolute
	JumpForward
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

var Names = map[byte]string{
	Noop:             "NOOP",
	LoadConst:        "LOAD_CONST",
	LoadFast:         "LOAD_FAST",
	StoreFast:        "STORE_FAST",
	LoadGlobal:       "LOAD_GLOBAL",
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
	UnaryNeg:         "UNARY_NEG",
	UnaryNot:         "UNARY_NOT",
	Compare:          "COMPARE",
	Call:             "CALL",
	Return:           "RETURN",
	Pop:              "POP",
	MakeArray:        "MAKE_ARRAY",
	MakeMap:          "MAKE_MAP",
	PopJumpIfTrue:    "POP_JUMP_IF_TRUE",
	PopJumpIfFalse:   "POP_JUMP_IF_FALSE",
	JumpIfTrueOrPop:  "JUMP_IF_TRUE_OR_POP",
	JumpIfFalseOrPop: "JUMP_IF_FALSE_OR_POP",
	JumpAbsolute:     "JUMP_ABSOLUTE",
	JumpForward:      "JUMP_FORWARD",
}

var CmpOps = map[byte]string{
	CmpEq:    "==",
	CmpNotEq: "!=",
	CmpLT:    "<",
	CmpGT:    ">",
	CmpLTEq:  "<=",
	CmpGTEq:  ">=",
}
