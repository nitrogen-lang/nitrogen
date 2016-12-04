package eval

import "github.com/lfkeitel/nitrogen/src/object"

var builtins = map[string]*object.Builtin{}

func registerBuiltin(name string, fn object.BuiltinFunction) {
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

func init() {
	registerBuiltin("len", lenBuiltin)
}

func lenBuiltin(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("Incorrect number of arguments. Got %d, expected 1", len(args))
	}

	switch arg := args[0].(type) {
	case *object.String:
		return &object.Integer{Value: int64(len(arg.Value))}
	}

	return newError("Unsupported type %s", args[0].Type())
}
