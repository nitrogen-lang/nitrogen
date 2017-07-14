package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"strconv"
	"strings"

	"github.com/nitrogen-lang/nitrogen/src/ast"
)

type ObjectType int

func (o ObjectType) String() string {
	return objectTypeNames[o]
}

const (
	INTEGER_OBJ ObjectType = iota
	FLOAT_OBJ
	BOOLEAN_OBJ
	NULL_OBJ
	RETURN_OBJ
	ERROR_OBJ
	FUNCTION_OBJ
	STRING_OBJ
	BUILTIN_OBJ
	ARRAY_OBJ
	HASH_OBJ
)

var objectTypeNames = map[ObjectType]string{
	INTEGER_OBJ:  "INTEGER",
	FLOAT_OBJ:    "FLOAT",
	BOOLEAN_OBJ:  "BOOLEAN",
	NULL_OBJ:     "NULL",
	RETURN_OBJ:   "RETURN",
	ERROR_OBJ:    "ERROR",
	FUNCTION_OBJ: "FUNCTION",
	STRING_OBJ:   "STRING",
	BUILTIN_OBJ:  "BUILTIN",
	ARRAY_OBJ:    "ARRAY",
	HASH_OBJ:     "MAP",
}

var (
	NULL  = &Null{}
	TRUE  = &Boolean{Value: true}
	FALSE = &Boolean{Value: false}
)

type BuiltinFunction func(env *Environment, args ...Object) Object

type Object interface {
	Type() ObjectType
	Inspect() string // Returns the value the object represents
}

type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string  { return strconv.FormatInt(i.Value, 10) }
func (i *Integer) Type() ObjectType { return INTEGER_OBJ }

type Float struct {
	Value float64
}

func (f *Float) Inspect() string  { return strconv.FormatFloat(f.Value, 'G', -1, 64) }
func (f *Float) Type() ObjectType { return FLOAT_OBJ }

type String struct {
	Value string
}

func (s *String) Inspect() string  { return s.Value }
func (s *String) Type() ObjectType { return STRING_OBJ }

type Boolean struct {
	Value bool
}

func (b *Boolean) Inspect() string {
	if b.Value {
		return "true"
	}
	return "false"
}
func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }

func NativeBoolToBooleanObj(input bool) *Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

type Null struct{}

func (n *Null) Inspect() string  { return "nil" }
func (n *Null) Type() ObjectType { return NULL_OBJ }

type ReturnValue struct {
	Value Object
}

func (r *ReturnValue) Inspect() string  { return r.Value.Inspect() }
func (r *ReturnValue) Type() ObjectType { return RETURN_OBJ }

// TODO: Expand with line/column numbers and stack trace
type Error struct {
	Message string
}

func (e *Error) Inspect() string  { return "ERROR: " + e.Message }
func (e *Error) Type() ObjectType { return ERROR_OBJ }

type Function struct {
	Name       string
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
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
func (f *Function) Type() ObjectType { return FUNCTION_OBJ }

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Inspect() string  { return "builtin function" }
func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }

type Array struct {
	Elements []Object
}

func (a *Array) Inspect() string {
	var out bytes.Buffer
	elements := []string{}

	for _, e := range a.Elements {
		elements = append(elements, e.Inspect())
	}

	out.WriteByte('[')
	out.WriteString(strings.Join(elements, ", "))
	out.WriteByte(']')

	return out.String()
}
func (a *Array) Type() ObjectType { return ARRAY_OBJ }

type Hashable interface {
	HashKey() HashKey
}

type HashKey struct {
	Type  ObjectType
	Value uint64
}

func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))
	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

type HashPair struct {
	Key   Object
	Value Object
}

type Hash struct {
	Pairs map[HashKey]HashPair
}

func (h *Hash) Type() ObjectType { return HASH_OBJ }
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

func NewError(format string, a ...interface{}) *Error {
	return &Error{Message: fmt.Sprintf(format, a...)}
}
