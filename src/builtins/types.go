package eval

import (
	"github.com/nitrogen-lang/nitrogen/src/eval"
	"github.com/nitrogen-lang/nitrogen/src/object"
)

func init() {
	eval.RegisterBuiltin("to_int", toIntBuiltin)
	eval.RegisterBuiltin("to_float", toFloatBuiltin)
}

func toIntBuiltin(env *object.Environment, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewError("Incorrect number of arguments. Got %d, expected 1", len(args))
	}

	switch arg := args[0].(type) {
	case *object.Integer:
		return arg
	case *object.Float:
		return &object.Integer{Value: int64(arg.Value)}
	}

	return object.NewError("Argument to `to_int` must be FLOAT or INT, got %s", args[0].Type())
}

func toFloatBuiltin(env *object.Environment, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewError("Incorrect number of arguments. Got %d, expected 1", len(args))
	}

	switch arg := args[0].(type) {
	case *object.Integer:
		return &object.Float{Value: float64(arg.Value)}
	case *object.Float:
		return arg
	}

	return object.NewError("Argument to `to_float` must be FLOAT or INT, got %s", args[0].Type())
}
