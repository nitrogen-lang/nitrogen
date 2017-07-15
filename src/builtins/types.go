package eval

import (
	"strconv"

	"github.com/nitrogen-lang/nitrogen/src/eval"
	"github.com/nitrogen-lang/nitrogen/src/object"
)

func init() {
	eval.RegisterBuiltin("toInt", toIntBuiltin)
	eval.RegisterBuiltin("toFloat", toFloatBuiltin)
	eval.RegisterBuiltin("toString", toStringBuiltin)

	eval.RegisterBuiltin("isFloat", makeIsTypeBuiltin(object.FLOAT_OBJ))
	eval.RegisterBuiltin("isInt", makeIsTypeBuiltin(object.INTEGER_OBJ))
	eval.RegisterBuiltin("isBool", makeIsTypeBuiltin(object.BOOLEAN_OBJ))
	eval.RegisterBuiltin("isNull", makeIsTypeBuiltin(object.NULL_OBJ))
	eval.RegisterBuiltin("isFunc", makeIsTypeBuiltin(object.FUNCTION_OBJ))
	eval.RegisterBuiltin("isString", makeIsTypeBuiltin(object.STRING_OBJ))
	eval.RegisterBuiltin("isArray", makeIsTypeBuiltin(object.ARRAY_OBJ))
	eval.RegisterBuiltin("isMap", makeIsTypeBuiltin(object.HASH_OBJ))
	// The below function is a placeholder for later
	// eval.RegisterBuiltin("isError", makeIsTypeBuiltin(object.ERROR_OBJ))
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

	return object.NewError("Argument to `toInt` must be FLOAT or INT, got %s", args[0].Type())
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

	return object.NewError("Argument to `toFloat` must be FLOAT or INT, got %s", args[0].Type())
}

func makeIsTypeBuiltin(t object.ObjectType) object.BuiltinFunction {
	return func(env *object.Environment, args ...object.Object) object.Object {
		if len(args) != 1 {
			return object.NewError("Type check requires one argument. Got %d", len(args))
		}

		return object.NativeBoolToBooleanObj(args[0].Type() == t)
	}
}

func toStringBuiltin(env *object.Environment, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewError("to_string expects 1 argument. Got %d", len(args))
	}

	converted := ""

	switch arg := args[0].(type) {
	case *object.String:
		converted = arg.Value
	case *object.Float:
		converted = strconv.FormatFloat(arg.Value, 'G', -1, 64)
	case *object.Integer:
		converted = strconv.FormatInt(arg.Value, 10)
	case *object.Boolean:
		converted = strconv.FormatBool(arg.Value)
	case *object.Null:
		converted = "nil"
	}

	return &object.String{Value: converted}
}
