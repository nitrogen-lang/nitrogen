package eval

import (
	"github.com/nitrogen-lang/nitrogen/src/object"
)

var builtins = map[string]*object.Builtin{}

func RegisterBuiltin(name string, fn object.BuiltinFunction) {
	if _, defined := builtins[name]; defined {
		// Panic because this should NEVER happen when built
		panic("Builtin function " + name + " already defined")
	}

	builtins[name] = &object.Builtin{Fn: fn}
}

func getBuiltin(name string) object.Object {
	if builtin, defined := builtins[name]; defined {
		return builtin
	}
	return nil
} 
