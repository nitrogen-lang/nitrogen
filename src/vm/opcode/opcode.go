package opcode

const (
	Noop byte = iota
	LoadConst
	StoreConst
	LoadFast
	StoreFast
	Define
	LoadGlobal
	StoreGlobal
	LoadIndex
	StoreIndex
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
	MakeFunction
	PopJumpIfTrue
	PopJumpIfFalse
	JumpIfTrueOrPop
	JumpIfFalseOrPop
	JumpAbsolute
	JumpForward
	PrepareBlock
	EndBlock
	StartLoop
	Continue
	NextIter
	Break
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

var HasFourByteArg = map[byte]bool{
	StartLoop: true,
}

var HasTwoByteArg = map[byte]bool{
	LoadConst:        true,
	StoreConst:       true,
	LoadFast:         true,
	StoreFast:        true,
	Define:           true,
	LoadGlobal:       true,
	StoreGlobal:      true,
	Call:             true,
	MakeArray:        true,
	MakeMap:          true,
	PopJumpIfTrue:    true,
	PopJumpIfFalse:   true,
	JumpIfTrueOrPop:  true,
	JumpIfFalseOrPop: true,
	JumpAbsolute:     true,
	JumpForward:      true,
}

var HasOneByteArg = map[byte]bool{
	Compare: true,
}

var HasNoArg = map[byte]bool{
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
	UnaryNeg:     true,
	UnaryNot:     true,
	Return:       true,
	Pop:          true,
	MakeFunction: true,
	PrepareBlock: true,
	EndBlock:     true,
	Continue:     true,
	NextIter:     true,
	Break:        true,
}

var Names = map[byte]string{
	Noop:             "NOOP",
	LoadConst:        "LOAD_CONST",
	StoreConst:       "STORE_CONST",
	LoadFast:         "LOAD_FAST",
	StoreFast:        "STORE_FAST",
	Define:           "DEFINE",
	LoadGlobal:       "LOAD_GLOBAL",
	StoreGlobal:      "STORE_GLOBAL",
	LoadIndex:        "LOAD_INDEX",
	StoreIndex:       "STORE_INDEX",
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
	MakeFunction:     "MAKE_FUNCTION",
	PopJumpIfTrue:    "POP_JUMP_IF_TRUE",
	PopJumpIfFalse:   "POP_JUMP_IF_FALSE",
	JumpIfTrueOrPop:  "JUMP_IF_TRUE_OR_POP",
	JumpIfFalseOrPop: "JUMP_IF_FALSE_OR_POP",
	JumpAbsolute:     "JUMP_ABSOLUTE",
	JumpForward:      "JUMP_FORWARD",
	PrepareBlock:     "PREPARE_BLOCK",
	EndBlock:         "END_BLOCK",
	StartLoop:        "START_LOOP",
	Continue:         "CONTINUE",
	NextIter:         "NEXT_ITER",
	Break:            "BREAK",
}

var CmpOps = map[byte]string{
	CmpEq:    "==",
	CmpNotEq: "!=",
	CmpLT:    "<",
	CmpGT:    ">",
	CmpLTEq:  "<=",
	CmpGTEq:  ">=",
}
