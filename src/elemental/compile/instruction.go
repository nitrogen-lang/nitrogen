package compile

import (
	"fmt"

	"github.com/nitrogen-lang/nitrogen/src/elemental/vm/opcode"
)

type Instruction struct {
	Instr     opcode.Opcode
	Args      []uint16 // len = 1 or 2
	ArgLabels []string // len = 1 or 2, name of label for corresponding argument, prefix "~" means relative
	Label     string   // Label names this instruction for linking later
	Prev      *Instruction
	Next      *Instruction
	Line      uint
}

func (i *Instruction) String() string {
	return opcode.Names[i.Instr]
}

func (i *Instruction) Is(code opcode.Opcode) bool {
	if i == nil {
		return false
	}
	return i.Instr == code
}

func (i *Instruction) IsLoad() bool {
	switch i.Instr {
	case opcode.LoadConst, opcode.LoadFast, opcode.LoadGlobal, opcode.LoadIndex, opcode.LoadAttribute:
		return true
	default:
		return false
	}
}

func (i *Instruction) Size() uint16 {
	switch {
	case i.Instr == opcode.Label:
		return 0
	case opcode.HasOneByteArg[i.Instr]:
		return 2
	case opcode.HasTwoByteArg[i.Instr]:
		return 3
	case opcode.HasThreeByteArg[i.Instr]:
		return 4
	case opcode.HasFourByteArg[i.Instr]:
		return 5
	}
	return 1
}

type InstSet struct {
	Head *Instruction
	Tail *Instruction
}

func NewInstSet() *InstSet {
	return &InstSet{}
}

func (i *InstSet) Last() *Instruction {
	if i.Tail == nil {
		return &Instruction{}
	}
	return i.Tail
}

func (i *InstSet) AddInst(code opcode.Opcode, line uint, args ...uint16) {
	checkArgLength(code, len(args))
	inst := &Instruction{
		Instr: code,
		Args:  args,
		Line:  line,
		Prev:  i.Tail,
	}

	if i.Head == nil {
		i.Head = inst
		i.Tail = inst
	} else {
		i.Tail.Next = inst
		i.Tail = inst
	}
}

func (i *InstSet) AddLabel(label string, line uint) {
	inst := &Instruction{
		Instr: opcode.Label,
		Label: label,
		Line:  line,
		Prev:  i.Tail,
	}

	if i.Head == nil {
		i.Head = inst
		i.Tail = inst
	} else {
		i.Tail.Next = inst
		i.Tail = inst
	}
}

func (i *InstSet) AddLabeledArgs(code opcode.Opcode, line uint, argLabels ...string) {
	checkArgLength(code, len(argLabels))
	inst := &Instruction{
		Instr:     code,
		Args:      make([]uint16, len(argLabels)),
		ArgLabels: argLabels,
		Line:      line,
		Prev:      i.Tail,
	}

	if i.Head == nil {
		i.Head = inst
		i.Tail = inst
	} else {
		i.Tail.Next = inst
		i.Tail = inst
	}
}

func (i *InstSet) Merge(j *InstSet) {
	i.Tail.Next = j.Head
	i.Tail = j.Tail
}

type Optimization func(*InstSet, *CodeBlockCompiler)

var optimizations = []Optimization{}

func AddOptimizer(o Optimization) {
	optimizations = append(optimizations, o)
}

func checkArgLength(code opcode.Opcode, argLen int) {
	if (opcode.HasOneByteArg[code] || opcode.HasTwoByteArg[code]) && argLen != 1 {
		panic(fmt.Sprintf("opcode %s requires 1 16-bit argument, given %d", code.String(), argLen))
	} else if opcode.HasThreeByteArg[code] && argLen != 2 {
		panic(fmt.Sprintf("opcode %s requires 1 16-bit and 1 8-bit argument, given %d", code.String(), argLen))
	} else if opcode.HasFourByteArg[code] && argLen != 2 {
		panic(fmt.Sprintf("opcode %s requires 2 16-bit argument, given %d", code.String(), argLen))
	}
}

func (i *InstSet) Len() uint16 {
	var size uint16
	in := i.Head
	for in != nil {
		size += in.Size()
		in = in.Next
	}
	return size
}

func (i *InstSet) Link() {
	labels := make(map[string]uint16)

	var offset uint16

	in := i.Head
	for in != nil {
		if in.Is(opcode.Label) {
			labels[in.Label] = offset
			in = in.Next
			continue
		}
		offset += in.Size()
		in = in.Next
	}

	offset = 0
	in = i.Head
	for in != nil {
		if in.ArgLabels != nil {
			for arg, lbl := range in.ArgLabels {
				if lbl != "" {
					in.Args[arg] = labels[lbl]
				}
			}
		}
		offset += in.Size()
		in = in.Next
	}
}

func (i *InstSet) Assemble(ccb *CodeBlockCompiler) ([]byte, []uint16) {
	for _, o := range optimizations {
		o(i, ccb)
	}

	i.Link()

	size := i.Len()
	bytes := make([]byte, size)
	offsetMap := make([]uint16, 0, 100)
	var lastLine uint
	var offset uint16

	in := i.Head
	for in != nil {
		if in.Is(opcode.Label) {
			in = in.Next
			continue
		}

		if in.Line != lastLine {
			offsetMap = append(offsetMap, offset, uint16(in.Line))
			lastLine = in.Line
		}
		bytes[offset] = in.Instr.ToByte()
		offset++

		switch {
		case opcode.HasOneByteArg[in.Instr]:
			bytes[offset] = byte(in.Args[0])
			offset++
		case opcode.HasTwoByteArg[in.Instr]:
			arg := uint16ToBytes(in.Args[0])
			bytes[offset] = arg[0]
			bytes[offset+1] = arg[1]
			offset += 2
		case opcode.HasThreeByteArg[in.Instr]:
			arg := uint16ToBytes(in.Args[0])
			bytes[offset] = arg[0]
			bytes[offset+1] = arg[1]

			bytes[offset+2] = byte(in.Args[1])

			offset += 3
		case opcode.HasFourByteArg[in.Instr]:
			arg := uint16ToBytes(in.Args[0])
			bytes[offset] = arg[0]
			bytes[offset+1] = arg[1]

			arg = uint16ToBytes(in.Args[1])
			bytes[offset+2] = arg[0]
			bytes[offset+3] = arg[1]
			offset += 4
		}
		in = in.Next
	}

	return bytes, offsetMap
}

type Code struct {
	code []byte
	i    int64
}

func NewCode(code []byte) *Code {
	return &Code{
		code: code,
	}
}

func (c *Code) NextInstruction() *Instruction {
	if c.i >= int64(len(c.code)) {
		return nil
	}

	curr := opcode.Opcode(c.code[c.i])

	i := &Instruction{
		Instr: curr,
		Args:  make([]uint16, 0),
	}
	c.i++

	switch {
	case opcode.HasOneByteArg[curr]:
		i.Args = []uint16{uint16(c.code[c.i])}
		c.i++
	case opcode.HasTwoByteArg[curr]:
		i.Args = []uint16{bytesToUint16(c.code[c.i], c.code[c.i+1])}
		c.i += 2
	case opcode.HasThreeByteArg[curr]:
		i.Args = []uint16{
			bytesToUint16(c.code[c.i], c.code[c.i+1]),
			uint16(c.code[c.i+2]),
		}
		c.i += 3
	case opcode.HasFourByteArg[curr]:
		i.Args = []uint16{
			bytesToUint16(c.code[c.i], c.code[c.i+1]),
			bytesToUint16(c.code[c.i+2], c.code[c.i+3]),
		}
		c.i += 4
	}

	return i
}
