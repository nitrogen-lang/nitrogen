package vm

import (
	"bytes"
	"regexp"
	"strings"

	"github.com/nitrogen-lang/nitrogen/src/compiler"
	"github.com/nitrogen-lang/nitrogen/src/object"
)

var (
	builtins   = map[string]*object.Builtin{}
	modules    = map[string]*object.Module{}
	identRegex = regexp.MustCompile(`[a-zA-Z_][a-zA-Z0-9_]*`)
)

// RegisterBuiltin allows other packages to register functions for availability in user code
func RegisterBuiltin(name string, fn object.BuiltinFunction) {
	if !validBuiltinIdent(name) {
		panic("Invalid VM builtin function name " + name)
	}

	if _, defined := builtins[name]; defined {
		// Panic because this should NEVER happen when built
		panic("Builtin VM function " + name + " already defined")
	}

	builtins[name] = &object.Builtin{Fn: fn}
}

func RegisterModule(name string, m *object.Module) {
	for k := range m.Methods {
		if !validBuiltinIdent(k) {
			panic("Invalid VM module function name " + name)
		}
	}

	if _, defined := modules[name]; defined {
		// Panic because this should NEVER happen when built
		panic("VM module " + name + " already defined")
	}

	modules[name] = m
}

func validBuiltinIdent(ident string) bool {
	return identRegex.Match([]byte(ident))
}

func getBuiltin(name string) object.Object {
	if builtin, defined := builtins[name]; defined {
		return builtin
	}
	return nil
}

// GetModule returns a Module object is a module with the given name is registered, otherwise nil.
func GetModule(name string) *object.Module {
	if module, defined := modules[name]; defined {
		return module
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
