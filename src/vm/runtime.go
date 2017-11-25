package vm

import (
	"bytes"
	"regexp"
	"strings"

	"github.com/nitrogen-lang/nitrogen/src/compiler"
	"github.com/nitrogen-lang/nitrogen/src/object"
)

type frame struct {
	lastFrame *frame
	code      *compiler.CodeBlock
	stack     *object.Stack
	locals    map[string]object.Object
	consts    map[string]object.Object
	outerVars map[string]object.Object
	pc        int
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
	val  *frame
	prev *stackElement
}

func newFrameStack() *frameStack {
	return &frameStack{}
}

func (s *frameStack) Push(val *frame) {
	s.head = &stackElement{
		val:  val,
		prev: s.head,
	}
	s.length++
}

func (s *frameStack) GetFront() *frame {
	if s.head == nil {
		return nil
	}
	return s.head.val
}

func (s *frameStack) Pop() *frame {
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
	Env        map[string]object.Object
	Consts     map[string]object.Object
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
