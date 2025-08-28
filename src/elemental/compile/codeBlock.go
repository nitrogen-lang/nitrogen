package compile

import (
	"encoding/binary"
	"fmt"

	"github.com/nitrogen-lang/nitrogen/src/elemental/object"
	"github.com/nitrogen-lang/nitrogen/src/elemental/vm/opcode"
)

type CodeBlock struct {
	Name       string // Fully qualified name
	Filename   string // Source file
	LocalCount int    // len(.Locals), this is here for convience

	// Maximum stack and block size can be calculated based on the instructions
	// used. This is a memory optimzation measure and can be larger than the
	// real maximum required to run this block. Each VM frame as a separate
	// stack.
	MaxStackSize int
	MaxBlockSize int

	Constants []object.Object // Constant VM objects used in the code
	Locals    []string        // Identifiers for local variables
	Names     []string        // Identifiers for non-local variables
	Code      []byte

	// Line numbers for each instruction are encoded as pairs of code offsets
	// and line numbers. See InstSet.Assemble for how it's created and .LineNum
	// for how it's decoded.
	LineOffsets []uint16

	// This CodeBlock represents a native-implemented function. If this is true, len(Code) == 0.
	Native      bool
	ClassMethod bool
}

// Implement object.Object interface

func (cb *CodeBlock) Type() object.ObjectType { return object.ResourceObj }
func (cb *CodeBlock) Inspect() string         { return "<codeblock>" }
func (cb *CodeBlock) Dup() object.Object      { return object.NullConst }
func (cb *CodeBlock) Print(indent string) {
	offset := 0
	lineOffsetIdx := 0

	for offset < len(cb.Code) {
		code := opcode.Opcode(cb.Code[offset])
		if lineOffsetIdx < len(cb.LineOffsets) && int(cb.LineOffsets[lineOffsetIdx]) == offset {
			fmt.Printf("%s%d\t%s%d:\t%s", indent, cb.LineOffsets[lineOffsetIdx+1], indent, offset, opcode.Names[code])
			lineOffsetIdx += 2
		} else {
			fmt.Printf("%s\t%s%d:\t%s", indent, indent, offset, opcode.Names[code])
		}
		offset++

		switch code {
		case opcode.MakeArray, opcode.MakeMap, opcode.Recover, opcode.BuildClass, opcode.MakeInstance:
			fmt.Printf("\t\t%d", bytesToUint16(cb.Code[offset], cb.Code[offset+1]))
		case opcode.JumpForward:
			target := int(bytesToUint16(cb.Code[offset], cb.Code[offset+1]))
			fmt.Printf("\t%d (%d)", target, offset+2+target)
		case opcode.JumpAbsolute:
			target := int(bytesToUint16(cb.Code[offset], cb.Code[offset+1]))
			fmt.Printf("\t%d", target)
		case opcode.StartLoop:
			fmt.Printf("\t%d %d", bytesToUint16(cb.Code[offset], cb.Code[offset+1]), bytesToUint16(cb.Code[offset+2], cb.Code[offset+3]))
		case opcode.PopJumpIfTrue, opcode.PopJumpIfFalse, opcode.JumpIfTrueOrPop, opcode.JumpIfFalseOrPop:
			fmt.Printf("\t%d", bytesToUint16(cb.Code[offset], cb.Code[offset+1]))
		case opcode.LoadConst, opcode.Import:
			index := bytesToUint16(cb.Code[offset], cb.Code[offset+1])
			fmt.Printf("\t%d (%s)", index, cb.Constants[index].Inspect())
		case opcode.LoadFast, opcode.StoreFast, opcode.DeleteFast:
			index := bytesToUint16(cb.Code[offset], cb.Code[offset+1])
			fmt.Printf("\t%d (%s)", index, cb.Locals[index])
		case opcode.Define:
			index := bytesToUint16(cb.Code[offset], cb.Code[offset+1])
			flags := cb.Code[offset+2]
			fmt.Printf("\t\t%d (%s) (%#.2x)", index, cb.Locals[index], flags)
		case opcode.Call:
			params := bytesToUint16(cb.Code[offset], cb.Code[offset+1])
			fmt.Printf("\t\t%d (%d positional parameters)", params, params)
		case opcode.LoadGlobal, opcode.StoreGlobal, opcode.LoadAttribute, opcode.StoreAttribute:
			index := bytesToUint16(cb.Code[offset], cb.Code[offset+1])
			fmt.Printf("\t%d (%s)", index, cb.Names[index])
		case opcode.Compare:
			fmt.Printf("\t%d (%s)", cb.Code[offset], opcode.CmpOps[cb.Code[offset]])
		}

		switch {
		case opcode.HasOneByteArg[code]:
			offset++
		case opcode.HasTwoByteArg[code]:
			offset += 2
		case opcode.HasThreeByteArg[code]:
			offset += 3
		case opcode.HasFourByteArg[code]:
			offset += 4
		}

		fmt.Println()
	}
}

func (cb *CodeBlock) LineNum(pc int) uint16 {
	offset := 0
	lineOffsetIdx := 0
	var line uint16

	for offset < pc {
		code := opcode.Opcode(cb.Code[offset])
		if lineOffsetIdx < len(cb.LineOffsets) && int(cb.LineOffsets[lineOffsetIdx]) == offset {
			line = cb.LineOffsets[lineOffsetIdx+1]
			lineOffsetIdx += 2
		}
		offset++

		switch {
		case opcode.HasOneByteArg[code]:
			offset++
		case opcode.HasTwoByteArg[code]:
			offset += 2
		case opcode.HasThreeByteArg[code]:
			offset += 3
		case opcode.HasFourByteArg[code]:
			offset += 4
		}
	}

	return line
}

func bytesToUint16(a, b byte) uint16 {
	return (uint16(a) << 8) + uint16(b)
}

type CodeBlockCompiler struct {
	Constants      *ConstantTable // Constant VM objects used in the code
	Locals         *StringTable   // Identifiers for local variables
	Names          *StringTable   // Identifiers for non-local variables
	Code           *InstSet
	Filename, Name string
	InLoop         bool
	Linenum        uint // Absolute line number for the first line of this block in the source file
}

type ConstantTable struct {
	Table []object.Object
}

func NewConstantTable() *ConstantTable {
	return &ConstantTable{
		Table: make([]object.Object, 0, 5),
	}
}

func (t *ConstantTable) IndexOf(v object.Object) uint16 {
	for i, o := range t.Table {
		if o.Type() != v.Type() {
			continue
		}

		switch node := v.(type) {
		case *object.Null:
			return uint16(i)
		case *object.Integer:
			if node.Value == o.(*object.Integer).Value {
				return uint16(i)
			}
		case *object.String:
			if node.Inspect() == o.(*object.String).Inspect() {
				return uint16(i)
			}
		case *object.Float:
			if node.Value == o.(*object.Float).Value {
				return uint16(i)
			}
		case *object.Boolean:
			if node.Value == o.(*object.Boolean).Value {
				return uint16(i)
			}
		case *object.Interface:
			if node.Name == o.(*object.Interface).Name {
				return uint16(i)
			}
		}
	}

	t.Table = append(t.Table, v)
	return uint16(len(t.Table) - 1)
}

type StringTable struct {
	Table []string
}

func NewStringTable() *StringTable {
	return &StringTable{
		Table: make([]string, 0, 5),
	}
}

func NewStringTableOffset(offset int) *StringTable {
	if offset < 0 {
		offset = 0
	}
	return &StringTable{
		Table: make([]string, offset, offset+5),
	}
}

func (t *StringTable) IndexOf(v string) uint16 {
	for i, s := range t.Table {
		if s == v {
			return uint16(i)
		}
	}

	t.Table = append(t.Table, v)
	return uint16(len(t.Table) - 1)
}

func (t *StringTable) Contains(s string) bool {
	for _, v := range t.Table {
		if v == s {
			return true
		}
	}
	return false
}

func uint16ToBytes(i uint16) []byte {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, i)
	return b
}
