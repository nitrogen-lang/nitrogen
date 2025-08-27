package compiler

import (
	"github.com/nitrogen-lang/nitrogen/src/elemental/compile"
	"github.com/nitrogen-lang/nitrogen/src/elemental/object"
	"github.com/nitrogen-lang/nitrogen/src/elemental/vm/opcode"
)

func init() {
	compile.AddOptimizer(optimizeLoadPop)
	compile.AddOptimizer(optimizeNegativeNums)
	compile.AddOptimizer(optimizeDefineLoadFast)
}

// optimizeLoadPop removes the pattern LOAD_ followed by POP.
// A Load Pop doesn't do anything since the value isn't being stored.
// This function does not account for Labels which would require
// calculating jump targets and blocks.
func optimizeLoadPop(i *compile.InstSet, ccb *compile.CodeBlockCompiler) {
	var last *compile.Instruction
	curr := i.Head
	for curr != nil && curr.Next != nil {
		if curr.IsLoad() && curr.Next.Is(opcode.Pop) && last != nil {
			last.Next = curr.Next.Next
			curr = curr.Next.Next
			continue
		}

		last = curr
		curr = curr.Next
	}
}

// optimizeNegativeNums replaces a LoadConst -> UnaryNeg sequence with
// a single LoadConst of the negative number.
func optimizeNegativeNums(i *compile.InstSet, ccb *compile.CodeBlockCompiler) {
	curr := i.Head

	for curr != nil && curr.Next != nil {
		if curr.Is(opcode.LoadConst) && curr.Next.Is(opcode.UnaryNeg) {
			numObj, ok := ccb.Constants.Table[curr.Args[0]].(*object.Integer)
			if ok {
				curr.Args[0] = ccb.Constants.IndexOf(object.MakeIntObj(-numObj.Value))
				curr.Next = curr.Next.Next
			}
		}

		curr = curr.Next
	}
}

func optimizeDefineLoadFast(i *compile.InstSet, ccb *compile.CodeBlockCompiler) {
	curr := i.Head

	for curr != nil && curr.Next != nil {
		if curr.Is(opcode.Define) && curr.Next.Is(opcode.LoadFast) {
			if curr.Args[0] == curr.Next.Args[0] {
				def := curr

				dup := &compile.Instruction{
					Instr: opcode.Dup,
					Line:  curr.Line,
					Next:  def,
				}

				curr.Prev.Next = dup
				def.Next = curr.Next.Next
			}
		}

		curr = curr.Next
	}
}
