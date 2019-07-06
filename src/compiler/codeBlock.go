package compiler

import (
	"encoding/binary"
	"fmt"

	"github.com/nitrogen-lang/nitrogen/src/object"
	"github.com/nitrogen-lang/nitrogen/src/vm/opcode"
)

type CodeBlock struct {
	Name         string
	Filename     string
	LocalCount   int
	MaxStackSize int
	MaxBlockSize int
	Constants    []object.Object // Created at compile time
	Locals       []string        // Created by compile time
	Names        []string        // Created at compile time
	Code         []byte
	Native       bool
	ClassMethod  bool
	LineOffsets  []uint16
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
		case opcode.MakeArray, opcode.MakeMap, opcode.StartTry, opcode.BuildClass, opcode.MakeInstance:
			fmt.Printf("\t\t%d", bytesToUint16(cb.Code[offset], cb.Code[offset+1]))
		case opcode.JumpForward:
			target := int(bytesToUint16(cb.Code[offset], cb.Code[offset+1]))
			fmt.Printf("\t\t%d (%d)", target, offset+2+target)
		case opcode.JumpAbsolute:
			target := int(bytesToUint16(cb.Code[offset], cb.Code[offset+1]))
			fmt.Printf("\t\t%d", target)
		case opcode.StartLoop:
			fmt.Printf("\t\t%d %d", bytesToUint16(cb.Code[offset], cb.Code[offset+1]), bytesToUint16(cb.Code[offset+2], cb.Code[offset+3]))
		case opcode.PopJumpIfTrue, opcode.PopJumpIfFalse, opcode.JumpIfTrueOrPop, opcode.JumpIfFalseOrPop:
			fmt.Printf("\t%d", bytesToUint16(cb.Code[offset], cb.Code[offset+1]))
		case opcode.LoadConst, opcode.Import:
			index := bytesToUint16(cb.Code[offset], cb.Code[offset+1])
			fmt.Printf("\t\t%d (%s)", index, cb.Constants[index].Inspect())
		case opcode.LoadFast, opcode.StoreFast, opcode.StoreConst, opcode.DeleteFast:
			index := bytesToUint16(cb.Code[offset], cb.Code[offset+1])
			fmt.Printf("\t\t%d (%s)", index, cb.Locals[index])
		case opcode.Define:
			index := bytesToUint16(cb.Code[offset], cb.Code[offset+1])
			fmt.Printf("\t\t\t%d (%s)", index, cb.Locals[index])
		case opcode.Call:
			params := bytesToUint16(cb.Code[offset], cb.Code[offset+1])
			fmt.Printf("\t\t\t%d (%d positional parameters)", params, params)
		case opcode.LoadGlobal, opcode.StoreGlobal, opcode.LoadAttribute, opcode.StoreAttribute:
			index := bytesToUint16(cb.Code[offset], cb.Code[offset+1])
			fmt.Printf("\t\t%d (%s)", index, cb.Names[index])
		case opcode.Compare:
			fmt.Printf("\t\t\t%d (%s)", cb.Code[offset], opcode.CmpOps[cb.Code[offset]])
		}

		switch {
		case opcode.HasOneByteArg[code]:
			offset++
		case opcode.HasTwoByteArg[code]:
			offset += 2
		case opcode.HasFourByteArg[code]:
			offset += 4
		}

		fmt.Println()
	}
}

func bytesToUint16(a, b byte) uint16 {
	return (uint16(a) << 8) + uint16(b)
}

type codeBlockCompiler struct {
	constants      *constantTable
	locals         *stringTable
	names          *stringTable
	code           *InstSet
	filename, name string
	inLoop         bool
	linenum        uint
}

type constantTable struct {
	table []object.Object
}

func newConstantTable() *constantTable {
	return &constantTable{
		table: make([]object.Object, 0, 5),
	}
}

func (t *constantTable) indexOf(v object.Object) uint16 {
	for i, o := range t.table {
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
		}
	}

	t.table = append(t.table, v)
	return uint16(len(t.table) - 1)
}

type stringTable struct {
	table []string
}

func newStringTable() *stringTable {
	return &stringTable{
		table: make([]string, 0, 5),
	}
}

func newStringTableOffset(offset int) *stringTable {
	if offset < 0 {
		offset = 0
	}
	return &stringTable{
		table: make([]string, offset, offset+5),
	}
}

func (t *stringTable) indexOf(v string) uint16 {
	for i, s := range t.table {
		if s == v {
			return uint16(i)
		}
	}

	t.table = append(t.table, v)
	return uint16(len(t.table) - 1)
}

func (t *stringTable) contains(s string) bool {
	for _, v := range t.table {
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
