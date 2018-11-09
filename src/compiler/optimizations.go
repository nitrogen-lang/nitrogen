package compiler

import (
	"github.com/nitrogen-lang/nitrogen/src/object"
	"github.com/nitrogen-lang/nitrogen/src/vm/opcode"
)

func init() {
	AddOptimizer(optimizeLoadPop)
	AddOptimizer(optimizeNegativeNums)
}

// optimizeLoadPop removes the pattern LOAD_ followed by POP.
// A Load Pop doesn't do anything since the value isn't being stored.
// This function does not account for Labels which would require
// calculating jump targets and blocks.
func optimizeLoadPop(i *InstSet, ccb *codeBlockCompiler) {
	var last *Instruction
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
func optimizeNegativeNums(i *InstSet, ccb *codeBlockCompiler) {
	curr := i.Head

	for curr != nil && curr.Next != nil {
		if curr.Is(opcode.LoadConst) && curr.Next.Is(opcode.UnaryNeg) {
			numObj, ok := ccb.constants.table[curr.Args[0]].(*object.Integer)
			if ok {
				curr.Args[0] = ccb.constants.indexOf(object.MakeIntObj(-numObj.Value))
				curr.Next = curr.Next.Next
			}
		}

		curr = curr.Next
	}
}
