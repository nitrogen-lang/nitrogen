package object

import (
	"bytes"
	"strconv"
	"strings"

	"github.com/lfkeitel/nitrogen/src/ast"
)

type ObjectType int

const (
	INTEGER_OBJ ObjectType = iota
	BOOLEAN_OBJ
	NULL_OBJ
	RETURN_OBJ
	ERROR_OBJ
	FUNCTION_OBJ
)

var objectTypeNames = map[ObjectType]string{
	INTEGER_OBJ:  "INTEGER",
	BOOLEAN_OBJ:  "BOOLEAN",
	NULL_OBJ:     "NULL",
	RETURN_OBJ:   "RETURN",
	ERROR_OBJ:    "ERROR",
	FUNCTION_OBJ: "FUNCTION",
}

type Object interface {
	Type() ObjectType
	Inspect() string // Returns the value the object represents
	String() string  // Returns a string of the Object type
}

type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string  { return strconv.FormatInt(i.Value, 10) }
func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) String() string   { return objectTypeNames[INTEGER_OBJ] }

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
func (b *Boolean) String() string   { return objectTypeNames[BOOLEAN_OBJ] }

type Null struct{}

func (n *Null) Inspect() string  { return "null" }
func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) String() string   { return objectTypeNames[NULL_OBJ] }

type ReturnValue struct {
	Value Object
}

func (r *ReturnValue) Inspect() string  { return r.Value.Inspect() }
func (r *ReturnValue) Type() ObjectType { return RETURN_OBJ }
func (r *ReturnValue) String() string   { return objectTypeNames[RETURN_OBJ] }

// TODO: Expand with line/column numbers and stack trace
type Error struct {
	Message string
}

func (e *Error) Inspect() string  { return "ERROR: " + e.Message }
func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) String() string   { return objectTypeNames[ERROR_OBJ] }

type Function struct {
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

	out.WriteString("fn")
	out.WriteByte('(')
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(f.Body.String())

	return out.String()
}
func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) String() string   { return objectTypeNames[FUNCTION_OBJ] }
