package compiler

import "github.com/nitrogen-lang/nitrogen/src/vm/opcode"

type maxsizer struct {
	max, current int
}

func (s *maxsizer) sub(delta int) {
	s.current -= delta
	if s.current > s.max { // Delta can be negative which would add to the size
		s.max = s.current
	}
}
func (s *maxsizer) add(delta int) {
	s.current += delta
	if s.current > s.max {
		s.max = s.current
	}
}

func calculateStackSize(c *InstSet) int {
	stackSize := &maxsizer{}

	i := c.Head
	for i != nil {
		switch i.Instr {
		case opcode.LoadConst, opcode.LoadFast, opcode.LoadGlobal, opcode.StartTry, opcode.Import, opcode.Dup:
			stackSize.add(1)
		case opcode.StoreIndex:
			stackSize.sub(3)
		case opcode.BinaryAdd, opcode.BinarySub, opcode.BinaryMul, opcode.BinaryDivide, opcode.BinaryMod, opcode.BinaryShiftL,
			opcode.BinaryShiftR, opcode.BinaryAnd, opcode.BinaryOr, opcode.BinaryNot, opcode.BinaryAndNot,
			opcode.StoreConst, opcode.StoreFast, opcode.Define, opcode.StoreGlobal, opcode.LoadIndex, opcode.Compare,
			opcode.Return, opcode.Pop, opcode.PopJumpIfTrue, opcode.PopJumpIfFalse, opcode.Throw:
			stackSize.sub(1)
		case opcode.Call:
			stackSize.sub(int(i.Args[0]))
		case opcode.MakeArray:
			stackSize.sub(int(i.Args[0]) - 1)
		case opcode.BuildClass:
			stackSize.sub(int(i.Args[0]) + 2)
		case opcode.MakeMap:
			stackSize.sub(int(i.Args[0])*2 - 1)
		case opcode.MakeFunction, opcode.StoreAttribute:
			stackSize.sub(2)
		}
		i = i.Next
	}

	return stackSize.max
}

func calculateBlockSize(c *InstSet) int {
	blockLen := &maxsizer{}

	i := c.Head
	for i != nil {
		switch i.Instr {
		case opcode.StartLoop, opcode.StartTry:
			blockLen.add(1)
		case opcode.EndBlock:
			blockLen.sub(1)
		}
		i = i.Next
	}

	return blockLen.max
}
