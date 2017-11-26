package vm

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

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
	head   *stackElement
	length int
}

type VMBuiltinFunc func(vm object.Interpreter, args ...object.Object) object.Object

type vmBuiltin struct {
	fn VMBuiltinFunc
}

func (vb *vmBuiltin) Type() object.ObjectType { return object.ResourceObj }
func (vb *vmBuiltin) Inspect() string         { return "<vmBuiltin>" }
func (vb *vmBuiltin) Dup() object.Object      { return object.NullConst }

type stackElement struct {
	val  *Frame
	prev *stackElement
}

func newFrameStack() *frameStack {
	return &frameStack{}
}

func (s *frameStack) Push(val *Frame) {
	s.head = &stackElement{
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

var (
	builtins   = map[string]*vmBuiltin{}
	identRegex = regexp.MustCompile(`[a-zA-Z_][a-zA-Z0-9_]*`)
)

// RegisterBuiltin allows other packages to register functions for availability in user code
func RegisterBuiltin(name string, fn VMBuiltinFunc) {
	if !validBuiltinIdent(name) {
		panic("Invalid VM builtin function name " + name)
	}

	if _, defined := builtins[name]; defined {
		// Panic because this should NEVER happen when built
		panic("Builtin VM function " + name + " already defined")
	}

	builtins[name] = &vmBuiltin{fn: fn}
}

func validBuiltinIdent(ident string) bool {
	return identRegex.Match([]byte(ident))
}

func getBuiltin(name string) *vmBuiltin {
	if builtin, defined := builtins[name]; defined {
		return builtin
	}
	return nil
}

type VMFunction struct {
	Name       string
	Parameters []string
	Body       *compiler.CodeBlock
	Env        *object.Environment
}

func (f *VMFunction) Inspect() string {
	var out bytes.Buffer

	out.WriteString("func")
	out.WriteByte(' ')
	out.WriteString(f.Name)
	out.WriteByte('(')
	out.WriteString(strings.Join(f.Parameters, ", "))
	out.WriteString(") {...}")

	return out.String()
}
func (f *VMFunction) Type() object.ObjectType { return object.FunctionObj }
func (f *VMFunction) Dup() object.Object      { return f }
