package compiler

import (
	"github.com/nitrogen-lang/nitrogen/src/vm/opcode"
)

func init() {
	AddOptimizer(optimizeLoadPop)
}

// optimizeLoadPop removes the pattern LOAD_ followed by POP.
// A Load Pop doesn't do anything since the value isn't being stored.
// This function does not account for Labels which would require
// calculating jump targets and blocks.
func optimizeLoadPop(i *InstSet) {
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
