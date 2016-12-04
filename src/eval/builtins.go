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
	registerBuiltin("first", firstBuiltin)
	registerBuiltin("last", lastBuiltin)
	registerBuiltin("rest", restBuiltin)
	registerBuiltin("push", pushBuiltin)
}

func lenBuiltin(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("Incorrect number of arguments. Got %d, expected 1", len(args))
	}

	switch arg := args[0].(type) {
	case *object.String:
		return &object.Integer{Value: int64(len(arg.Value))}
	case *object.Array:
		return &object.Integer{Value: int64(len(arg.Elements))}
	}

	return newError("Unsupported type %s", args[0].Type())
}

func firstBuiltin(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("Incorrect number of arguments. Got %d, expected 1", len(args))
	}
	if args[0].Type() != object.ARRAY_OBJ {
		return newError("Argument to `first` must be ARRAY, got %s", args[0].Type())
	}

	arr := args[0].(*object.Array)
	if len(arr.Elements) > 0 {
		return arr.Elements[0]
	}

	return NULL
}

func lastBuiltin(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("Incorrect number of arguments. Got %d, expected 1", len(args))
	}
	if args[0].Type() != object.ARRAY_OBJ {
		return newError("Argument to `last` must be ARRAY, got %s", args[0].Type())
	}

	arr := args[0].(*object.Array)
	length := len(arr.Elements)
	if length > 0 {
		return arr.Elements[length-1]
	}

	return NULL
}

func restBuiltin(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("Incorrect number of arguments. Got %d, expected 1", len(args))
	}
	if args[0].Type() != object.ARRAY_OBJ {
		return newError("Argument to `rest` must be ARRAY, got %s", args[0].Type())
	}

	arr := args[0].(*object.Array)
	length := len(arr.Elements)
	if length > 0 {
		newElements := make([]object.Object, length-1, length-1)
		copy(newElements, arr.Elements[1:length])
		return &object.Array{Elements: newElements}
	}

	return NULL
}

func pushBuiltin(args ...object.Object) object.Object {
	if len(args) != 2 {
		return newError("Incorrect number of arguments. Got %d, expected 2", len(args))
	}
	if args[0].Type() != object.ARRAY_OBJ {
		return newError("Argument to `push` must be ARRAY, got %s", args[0].Type())
	}

	arr := args[0].(*object.Array)
	length := len(arr.Elements)
	newElements := make([]object.Object, length+1, length+1)
	copy(newElements, arr.Elements)
	newElements[length] = args[1]

	return &object.Array{Elements: newElements}
}
