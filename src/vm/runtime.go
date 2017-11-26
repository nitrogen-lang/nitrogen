package vm

import (
	"fmt"

	"github.com/nitrogen-lang/nitrogen/src/compiler"
	"github.com/nitrogen-lang/nitrogen/src/object"
)

type block struct {
	start, iter, end int
}

type Frame struct {
	lastFrame  *Frame
	code       *compiler.CodeBlock
	stack      []object.Object
	sp         int
	blockStack []block
	bp         int
	Env        *object.Environment
	pc         int
}

func (f *Frame) pushStack(obj object.Object) {
	if f.sp == len(f.stack) {
		panic("Stack overflow")
	}
	f.stack[f.sp] = obj
	f.sp++
}

func (f *Frame) popStack() object.Object {
	if f.sp == 0 {
		panic("Stack exhausted")
	}
	f.sp--
	return f.stack[f.sp]
}

func (f *Frame) getFrontStack() object.Object {
	return f.stack[f.sp-1]
}

func (f *Frame) printStack() {
	for i := f.sp; i >= 0; i-- {
		fmt.Printf(" %d: %s\n", i, f.stack[i].Inspect())
	}
}

func (f *Frame) pushBlock(start, iter, end int) {
	if f.bp == len(f.blockStack) {
		panic("Block stack overflow")
	}
	f.blockStack[f.bp].start = start
	f.blockStack[f.bp].iter = iter
	f.blockStack[f.bp].end = end
	f.bp++
}

// returns start, end
func (f *Frame) popBlock() (int, int, int) {
	if f.bp == 0 {
		panic("Block stack exhausted")
	}
	f.bp--
	return f.blockStack[f.bp].start, f.blockStack[f.bp].iter, f.blockStack[f.bp].end
}

// returns start, end
func (f *Frame) getCurrentBlock() (int, int, int) {
	return f.blockStack[f.bp-1].start, f.blockStack[f.bp-1].iter, f.blockStack[f.bp-1].end
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
