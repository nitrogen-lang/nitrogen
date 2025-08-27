package vm

import (
	"fmt"

	"github.com/nitrogen-lang/nitrogen/src/elemental/compile"
	"github.com/nitrogen-lang/nitrogen/src/elemental/object"
)

type blockType byte

const (
	loopBlockT blockType = iota
	tryBlockT
	doBlockT
)

type block interface {
	blockType() blockType
}

type forLoopBlock struct {
	start, iter, end int
}

func (b *forLoopBlock) blockType() blockType { return loopBlockT }

type recoverBlock struct {
	pc, sp int
	caught bool
}

func (b *recoverBlock) blockType() blockType { return tryBlockT }

type doBlock struct{}

func (b *doBlock) blockType() blockType { return doBlockT }

type Frame struct {
	module     string
	lastFrame  *Frame
	code       *compile.CodeBlock
	stack      []object.Object
	sp         int
	blockStack []block
	bp         int
	env        *object.Environment
	pc         int
	unwind     bool
}

func (f *Frame) lineno() uint {
	pc := uint16(f.pc)
	var linenum uint16

	for i := 0; i < len(f.code.LineOffsets); i += 2 {
		addr := f.code.LineOffsets[i]
		line := f.code.LineOffsets[i+1]

		if addr >= pc {
			return uint(linenum)
		}

		linenum = line
	}
	return uint(linenum)
}

func (f *Frame) pushStack(obj object.Object) {
	if f == nil {
		return
	}

	if f.sp == len(f.stack) {
		panic(fmt.Sprintf("VM stack exhausted (%d) - %s (%d)", len(f.stack), f.code.Name, f.pc))
	}
	f.stack[f.sp] = obj
	f.sp++
}

func (f *Frame) popStack() object.Object {
	if f.sp == 0 {
		panic(fmt.Sprintf("VM stack exhausted - %s (%d)", f.code.Name, f.pc))
	}
	f.sp--
	return f.stack[f.sp]
}

func (f *Frame) getFrontStack() object.Object {
	return f.stack[f.sp-1]
}

func (f *Frame) printStack() {
	for i := f.sp - 1; i >= 0; i-- {
		fmt.Printf(" %d: %s\n", i, f.stack[i].Inspect())
	}
}

func (f *Frame) pushBlock(b block) {
	if f.bp == len(f.blockStack) {
		panic("Block stack overflow")
	}
	f.blockStack[f.bp] = b
	f.bp++
}

func (f *Frame) popBlock() block {
	if f.bp == 0 {
		panic("Block stack exhausted")
	}
	f.bp--
	return f.blockStack[f.bp]
}

func (f *Frame) popBlockUntil(bt blockType) block {
	if f.bp == 0 {
		return nil
	}

	for f.blockStack[f.bp-1].blockType() != bt {
		f.bp--
		if f.bp == 0 {
			return nil
		}
	}
	return f.blockStack[f.bp-1]
}

func (f *Frame) getCurrentBlock() block {
	return f.blockStack[f.bp-1]
}

type frameStack struct {
	head   *frameStackElement
	length int
}

type frameStackElement struct {
	val  *Frame
	prev *frameStackElement
}

func newFrameStack() *frameStack {
	return &frameStack{}
}

func (s *frameStack) Push(val *Frame) {
	s.head = &frameStackElement{
		val:  val,
		prev: s.head,
	}
	s.length++
}

func (s *frameStack) GetFront() *Frame {
	if s.head == nil {
		return nil
	}
	return s.head.val
}

func (s *frameStack) Pop() *Frame {
	if s.head == nil {
		return nil
	}
	r := s.head.val
	s.head = s.head.prev
	s.length--
	return r
}

func (s *frameStack) Len() int {
	return s.length
}

func (s *frameStack) forEach(fn func(*Frame)) {
	f := s.head

	for {
		if f == nil {
			return
		}
		fn(f.val)
		f = f.prev
	}
}
