package eval

import (
	"regexp"

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
		panic("Invalid builtin function name " + name)
	}

	if _, defined := builtins[name]; defined {
		// Panic because this should NEVER happen when built
		panic("Builtin function " + name + " already defined")
	}

	builtins[name] = &object.Builtin{Fn: fn}
}

// RegisterModule allows other packages to register a Module object for availability in user code
func RegisterModule(name string, m *object.Module) {
	for k := range m.Methods {
		if !validBuiltinIdent(k) {
			panic("Invalid module function name " + name)
		}
	}

	if _, defined := modules[name]; defined {
		// Panic because this should NEVER happen when built
		panic("Module " + name + " already defined")
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
