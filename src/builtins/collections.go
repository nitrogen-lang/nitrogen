package eval

import (
	"github.com/nitrogen-lang/nitrogen/src/eval"
	"github.com/nitrogen-lang/nitrogen/src/object"
)

func init() {
	eval.RegisterBuiltin("len", lenBuiltin)
	eval.RegisterBuiltin("first", firstBuiltin)
	eval.RegisterBuiltin("last", lastBuiltin)
	eval.RegisterBuiltin("rest", restBuiltin)
	eval.RegisterBuiltin("push", pushBuiltin)
}

func lenBuiltin(env *object.Environment, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewError("Incorrect number of arguments. Got %d, expected 1", len(args))
	}

	switch arg := args[0].(type) {
	case *object.String:
		return &object.Integer{Value: int64(len(arg.Value))}
	case *object.Array:
		return &object.Integer{Value: int64(len(arg.Elements))}
	}

	return object.NewError("Unsupported type %s", args[0].Type())
}

func firstBuiltin(env *object.Environment, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewError("Incorrect number of arguments. Got %d, expected 1", len(args))
	}
	if args[0].Type() != object.ARRAY_OBJ {
		return object.NewError("Argument to `first` must be ARRAY, got %s", args[0].Type())
	}

	arr := args[0].(*object.Array)
	if len(arr.Elements) > 0 {
		return arr.Elements[0]
	}

	return object.NULL
}

func lastBuiltin(env *object.Environment, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewError("Incorrect number of arguments. Got %d, expected 1", len(args))
	}
	if args[0].Type() != object.ARRAY_OBJ {
		return object.NewError("Argument to `last` must be ARRAY, got %s", args[0].Type())
	}

	arr := args[0].(*object.Array)
	length := len(arr.Elements)
	if length > 0 {
		return arr.Elements[length-1]
	}

	return object.NULL
}

func restBuiltin(env *object.Environment, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewError("Incorrect number of arguments. Got %d, expected 1", len(args))
	}
	if args[0].Type() != object.ARRAY_OBJ {
		return object.NewError("Argument to `rest` must be ARRAY, got %s", args[0].Type())
	}

	arr := args[0].(*object.Array)
	length := len(arr.Elements)
	if length > 0 {
		newElements := make([]object.Object, length-1, length-1)
		copy(newElements, arr.Elements[1:length])
		return &object.Array{Elements: newElements}
	}

	return object.NULL
}

func pushBuiltin(env *object.Environment, args ...object.Object) object.Object {
	if len(args) != 2 {
		return object.NewError("Incorrect number of arguments. Got %d, expected 2", len(args))
	}
	if args[0].Type() != object.ARRAY_OBJ {
		return object.NewError("Argument to `push` must be ARRAY, got %s", args[0].Type())
	}

	arr := args[0].(*object.Array)
	length := len(arr.Elements)
	newElements := make([]object.Object, length+1, length+1)
	copy(newElements, arr.Elements)
	newElements[length] = args[1]

	return &object.Array{Elements: newElements}
}
