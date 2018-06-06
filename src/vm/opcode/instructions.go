package opcode

import "encoding/binary"

type Instruction struct {
	instr Opcode
	args  []uint16
}

func NoopInstruction() *Instruction {
	return &Instruction{instr: Noop}
}

func (i Instruction) String() string {
	return Names[i.instr]
}

func (i Instruction) IsLoad() bool {
	switch i.instr {
	case LoadConst, LoadFast, LoadGlobal, LoadIndex, LoadAttribute:
		return true
	}
	return false
}

func (i Instruction) Is(code Opcode) bool {
	return i.instr == code
}

func InstructionLen(is []*Instruction) int64 {
	size := int64(0)
	for _, i := range is {
		size++

		if HasOneByteArg[i.instr] {
			size++
		} else if HasTwoByteArg[i.instr] {
			size += 2
		} else if HasFourByteArg[i.instr] {
			size += 4
		}
	}
	return size
}

func AssembleInstructions(is []*Instruction) []byte {
	bytes := make([]byte, InstructionLen(is))
	offset := 0

	for _, i := range is {
		bytes[offset] = i.instr.ToByte()
		offset++

		if HasOneByteArg[i.instr] {
			bytes[offset] = byte(i.args[0])
			offset++
		} else if HasTwoByteArg[i.instr] {
			arg := uint16ToBytes(i.args[0])
			bytes[offset] = arg[0]
			bytes[offset+1] = arg[1]
			offset += 2
		} else if HasFourByteArg[i.instr] {
			arg := uint16ToBytes(i.args[0])
			bytes[offset] = arg[0]
			bytes[offset+1] = arg[1]

			arg = uint16ToBytes(i.args[1])
			bytes[offset+2] = arg[0]
			bytes[offset+3] = arg[1]
			offset += 4
		}
	}

	return bytes
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

	curr := Opcode(c.code[c.i])

	i := &Instruction{
		instr: curr,
		args:  make([]uint16, 0),
	}
	c.i++

	if HasOneByteArg[curr] {
		i.args = []uint16{uint16(c.code[c.i])}
		c.i++
	} else if HasTwoByteArg[curr] {
		i.args = []uint16{bytesToUint16(c.code[c.i], c.code[c.i+1])}
		c.i += 2
	} else if HasFourByteArg[curr] {
		i.args = []uint16{
			bytesToUint16(c.code[c.i], c.code[c.i+1]),
			bytesToUint16(c.code[c.i+2], c.code[c.i+3]),
		}
		c.i += 4
	}

	return i
}

func uint16ToBytes(i uint16) []byte {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, i)
	return b
}

func bytesToUint16(a, b byte) uint16 {
	return uint16(b) + uint16(a<<4)
}
