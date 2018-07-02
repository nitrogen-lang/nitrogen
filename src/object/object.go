package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"io"
	"strconv"
	"strings"

	"github.com/nitrogen-lang/nitrogen/src/ast"
)

type ObjectType int

func (o ObjectType) String() string {
	return objectTypeNames[o]
}

// These are all the internal object types used in the interpreter
const (
	IntergerObj ObjectType = iota
	FloatObj
	BooleanObj
	NullObj
	ReturnObj
	ExceptionObj
	ErrorObj
	FunctionObj
	StringObj
	BuiltinObj
	ArrayObj
	HashObj
	LoopControlObj
	// ResourceObj can be used by modules to denote a generic resource. The implementation may be module specific.
	ResourceObj
	ModuleObj
	ClassObj
	InstanceObj
	BuiltinMethodObj
	BoundMethodObj
)

var objectTypeNames = map[ObjectType]string{
	IntergerObj:      "INTEGER",
	FloatObj:         "FLOAT",
	BooleanObj:       "BOOLEAN",
	NullObj:          "NULL",
	ReturnObj:        "RETURN",
	ExceptionObj:     "EXCEPTION",
	ErrorObj:         "ERROR",
	FunctionObj:      "FUNCTION",
	StringObj:        "STRING",
	BuiltinObj:       "BUILTIN",
	ArrayObj:         "ARRAY",
	HashObj:          "MAP",
	ResourceObj:      "RESOURCE",
	ModuleObj:        "MODULE",
	ClassObj:         "CLASS",
	InstanceObj:      "INSTANCE",
	BuiltinMethodObj: "BUILTIN METHOD",
	BoundMethodObj:   "BOUND METHOD",
}

// These are all constants in the language that can be represented with a single instance
var (
	NullConst  = &Null{}
	TrueConst  = &Boolean{Value: true}
	FalseConst = &Boolean{Value: false}
)

type Interpreter interface {
	Eval(node ast.Node, env *Environment) Object
	GetCurrentScriptPath() string
	GetStdout() io.Writer
	GetStderr() io.Writer
	GetStdin() io.Reader
}

type BuiltinFunction func(i Interpreter, env *Environment, args ...Object) Object

type BuiltinMethodFunction func(i Interpreter, self *Instance, env *Environment, args ...Object) Object

type Object interface {
	Type() ObjectType
	Inspect() string // Returns the value the object represents
	Dup() Object     // Returns a duplicate of the object
}

type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string  { return strconv.FormatInt(i.Value, 10) }
func (i *Integer) Type() ObjectType { return IntergerObj }
func (i *Integer) Dup() Object      { return &Integer{Value: i.Value} }

func MakeIntObj(v int64) *Integer {
	return &Integer{Value: v}
}

type Float struct {
	Value float64
}

func (f *Float) Inspect() string  { return strconv.FormatFloat(f.Value, 'G', -1, 64) }
func (f *Float) Type() ObjectType { return FloatObj }
func (f *Float) Dup() Object      { return &Float{Value: f.Value} }

func MakeFloatObj(v float64) *Float {
	return &Float{Value: v}
}

type String struct {
	Value []rune
}

func (s *String) Inspect() string  { return string(s.Value) }
func (s *String) String() string   { return string(s.Value) } // Dedicated to the stringified value, no inspection
func (s *String) Type() ObjectType { return StringObj }
func (s *String) Dup() Object      { return &String{Value: s.Value[:]} }

func MakeStringObj(s string) *String {
	return &String{Value: []rune(s)}
}

type Boolean struct {
	Value bool
}

func (b *Boolean) Inspect() string {
	if b.Value {
		return "true"
	}
	return "false"
}
func (b *Boolean) Type() ObjectType { return BooleanObj }
func (b *Boolean) Dup() Object {
	if b.Value {
		return TrueConst
	}
	return FalseConst
}

func NativeBoolToBooleanObj(input bool) *Boolean {
	if input {
		return TrueConst
	}
	return FalseConst
}

type Null struct{}

func (n *Null) Inspect() string  { return "nil" }
func (n *Null) Type() ObjectType { return NullObj }
func (n *Null) Dup() Object      { return NullConst }

type ReturnValue struct {
	Value Object
}

func (r *ReturnValue) Inspect() string  { return r.Value.Inspect() }
func (r *ReturnValue) Type() ObjectType { return ReturnObj }
func (r *ReturnValue) Dup() Object      { return &ReturnValue{Value: r.Value.Dup()} }

// TODO: Expand with line/column numbers and stack trace
type Exception struct {
	Catchable bool
	Message   string
	Caught    bool
}

func (e *Exception) Inspect() string  { return e.Message }
func (e *Exception) Type() ObjectType { return ExceptionObj }
func (e *Exception) Dup() Object      { return &Exception{Message: e.Message} }
func (e *Exception) String() string   { return e.Message }

type Error struct {
	Message string
}

func (e *Error) Inspect() string  { return "Error: " + e.Message }
func (e *Error) Type() ObjectType { return ErrorObj }
func (e *Error) Dup() Object      { return &Error{Message: e.Message} }
func (e *Error) String() string   { return "Error: " + e.Message }

type Function struct {
	Name       string
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
	Instance   *Instance
}

func (f *Function) Inspect() string {
	var out bytes.Buffer
	params := []string{}

	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("func")
	out.WriteByte(' ')
	out.WriteString(f.Name)
	out.WriteByte('(')
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {")
	out.WriteString(f.Body.String())
	out.WriteByte('}')

	return out.String()
}
func (f *Function) Type() ObjectType { return FunctionObj }
func (f *Function) Dup() Object      { return f }
func (f *Function) ClassMethod()     {}

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Inspect() string  { return "builtin function" }
func (b *Builtin) Type() ObjectType { return BuiltinObj }
func (b *Builtin) Dup() Object      { return b }

type BuiltinMethod struct {
	Fn       BuiltinMethodFunction
	Instance *Instance
}

func (b *BuiltinMethod) Inspect() string  { return "builtin method" }
func (b *BuiltinMethod) Type() ObjectType { return BuiltinMethodObj }
func (b *BuiltinMethod) Dup() Object      { return b }
func (b *BuiltinMethod) ClassMethod()     {}

func MakeBuiltinMethod(fn BuiltinMethodFunction) *BuiltinMethod {
	return &BuiltinMethod{Fn: fn}
}

type Array struct {
	Elements []Object
}

func (a *Array) Inspect() string {
	var out bytes.Buffer
	elements := []string{}

	for _, e := range a.Elements {
		if e.Type() == StringObj {
			elements = append(elements, fmt.Sprintf(`"%s"`, e.Inspect()))
		} else {
			elements = append(elements, e.Inspect())
		}
	}

	out.WriteByte('[')
	out.WriteString(strings.Join(elements, ", "))
	out.WriteByte(']')

	return out.String()
}
func (a *Array) Type() ObjectType { return ArrayObj }
func (a *Array) Dup() Object {
	newElements := make([]Object, len(a.Elements))

	for i, element := range a.Elements {
		newElements[i] = element.Dup()
	}

	return &Array{Elements: newElements}
}

func MakeStringArray(e []string) *Array {
	length := len(e)
	newElements := make([]Object, length, length)
	for i, s := range e {
		newElements[i] = MakeStringObj(s)
	}
	return &Array{Elements: newElements}
}

func ArrayToStringSlice(a *Array) []string {
	if a == nil {
		return []string{}
	}

	strs := make([]string, 0, len(a.Elements))
	for _, e := range a.Elements {
		if !(e.Type() == StringObj) {
			continue
		}
		strs = append(strs, string((e.(*String)).Value))
	}
	return strs
}

type Hashable interface {
	HashKey() HashKey
}

type HashKey struct {
	Type  ObjectType
	Value uint64
}

func (k HashKey) Dup() HashKey {
	return HashKey{
		Type:  k.Type,
		Value: k.Value,
	}
}

func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(string(s.Value)))
	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

type HashPair struct {
	Key   Object
	Value Object
}

func (p HashPair) Dup() HashPair {
	return HashPair{
		Key:   p.Key.Dup(),
		Value: p.Value.Dup(),
	}
}

type Hash struct {
	Pairs map[HashKey]HashPair
}

func (h *Hash) Type() ObjectType { return HashObj }
func (h *Hash) Inspect() string {
	var out bytes.Buffer

	pairs := []string{}
	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s",
			pair.Key.Inspect(), pair.Value.Inspect()))
	}

	out.WriteByte('{')
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteByte('}')
	return out.String()
}
func (h *Hash) Dup() Object {
	newElements := make(map[HashKey]HashPair, len(h.Pairs))

	for key, pair := range h.Pairs {
		newElements[key.Dup()] = pair.Dup()
	}

	return &Hash{Pairs: newElements}
}

func StringMapToHash(src map[string]string) *Hash {
	m := &Hash{Pairs: make(map[HashKey]HashPair)}
	for k, v := range src {
		key := MakeStringObj(k)
		m.Pairs[key.HashKey()] = HashPair{
			Key:   key,
			Value: MakeStringObj(v),
		}
	}
	return m
}

type LoopControl struct {
	Continue bool
}

func (lc *LoopControl) Type() ObjectType { return LoopControlObj }
func (lc *LoopControl) Inspect() string {
	if lc.Continue {
		return "continue"
	}
	return "break"
}
func (lc *LoopControl) Dup() Object { return &LoopControl{Continue: lc.Continue} }

type Module struct {
	Name    string
	Methods map[string]BuiltinFunction
	Vars    map[string]Object
}

func (m *Module) Inspect() string  { return fmt.Sprintf("Module %s", m.Name) }
func (m *Module) Type() ObjectType { return ModuleObj }
func (m *Module) Dup() Object      { return NullConst }

func NewException(format string, a ...interface{}) *Exception {
	return &Exception{
		Message:   fmt.Sprintf(format, a...),
		Catchable: true,
	}
}

func NewPanic(format string, a ...interface{}) *Exception {
	return &Exception{Message: fmt.Sprintf(format, a...)}
}

func NewError(format string, a ...interface{}) *Error {
	return &Error{Message: fmt.Sprintf(format, a...)}
}

func ObjectsAre(t ObjectType, o ...Object) bool {
	for _, obj := range o {
		if obj.Type() != t {
			return false
		}
	}
	return true
}

func ObjectIs(o Object, t ...ObjectType) bool {
	if o != nil {
		for _, ot := range t {
			if o.Type() == ot {
				return true
			}
		}
	}
	return false
}
