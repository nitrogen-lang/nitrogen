package eval

import (
	"github.com/nitrogen-lang/nitrogen/src/eval"
	"github.com/nitrogen-lang/nitrogen/src/object"
)

func init() {
	eval.RegisterBuiltin("to_int", toIntBuiltin)
	eval.RegisterBuiltin("to_float", toFloatBuiltin)

	eval.RegisterBuiltin("isFloat", makeIsTypeBuiltin(object.FLOAT_OBJ))
	eval.RegisterBuiltin("isInt", makeIsTypeBuiltin(object.INTEGER_OBJ))
	eval.RegisterBuiltin("isBool", makeIsTypeBuiltin(object.BOOLEAN_OBJ))
	eval.RegisterBuiltin("isNull", makeIsTypeBuiltin(object.NULL_OBJ))
	eval.RegisterBuiltin("isFunc", makeIsTypeBuiltin(object.FUNCTION_OBJ))
	eval.RegisterBuiltin("isString", makeIsTypeBuiltin(object.STRING_OBJ))
	eval.RegisterBuiltin("isArray", makeIsTypeBuiltin(object.ARRAY_OBJ))
	eval.RegisterBuiltin("isMap", makeIsTypeBuiltin(object.HASH_OBJ))
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

func makeIsTypeBuiltin(t object.ObjectType) object.BuiltinFunction {
	return func(env *object.Environment, args ...object.Object) object.Object {
		if len(args) != 1 {
			return object.NewError("Type check requires one argument. Got %d", len(args))
		}

		return object.NativeBoolToBooleanObj(args[0].Type() == t)
	}
}
