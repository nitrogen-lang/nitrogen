package typing

import (
	"strconv"

	"github.com/nitrogen-lang/nitrogen/src/moduleutils"
	"github.com/nitrogen-lang/nitrogen/src/object"
	"github.com/nitrogen-lang/nitrogen/src/vm"
)

func init() {
	// Register with virual machine
	vm.RegisterBuiltin("toInt", toIntBuiltin)
	vm.RegisterBuiltin("toFloat", toFloatBuiltin)
	vm.RegisterBuiltin("toString", toStringBuiltin)

	vm.RegisterBuiltin("parseInt", parseIntBuiltin)
	vm.RegisterBuiltin("parseFloat", parseFloatBuiltin)

	vm.RegisterBuiltin("varType", varTypeBuiltin)
	vm.RegisterBuiltin("isDefined", isDefinedBuiltin)
	vm.RegisterBuiltin("isFloat", makeIsTypeBuiltin(object.FloatObj))
	vm.RegisterBuiltin("isInt", makeIsTypeBuiltin(object.IntergerObj))
	vm.RegisterBuiltin("isBool", makeIsTypeBuiltin(object.BooleanObj))
	vm.RegisterBuiltin("isNull", makeIsTypeBuiltin(object.NullObj))
	vm.RegisterBuiltin("isNil", makeIsTypeBuiltin(object.NullObj))
	vm.RegisterBuiltin("isFunc", makeIsTypeBuiltin(object.FunctionObj))
	vm.RegisterBuiltin("isString", makeIsTypeBuiltin(object.StringObj))
	vm.RegisterBuiltin("isArray", makeIsTypeBuiltin(object.ArrayObj))
	vm.RegisterBuiltin("isMap", makeIsTypeBuiltin(object.HashObj))
	vm.RegisterBuiltin("isError", makeIsTypeBuiltin(object.ErrorObj))
	vm.RegisterBuiltin("isResource", makeIsTypeBuiltin(object.ResourceObj))
	vm.RegisterBuiltin("isClass", makeIsTypeBuiltin(object.ClassObj))
	vm.RegisterBuiltin("isInstance", makeIsTypeBuiltin(object.InstanceObj))

	vm.RegisterBuiltin("errorVal", getErrorVal)
	vm.RegisterBuiltin("resourceID", getResourceID)
}

func toIntBuiltin(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewException("Incorrect number of arguments. Got %d, expected 1", len(args))
	}

	switch arg := args[0].(type) {
	case *object.Integer:
		return arg
	case *object.Float:
		return object.MakeIntObj(int64(arg.Value))
	}

	return object.NewException("Argument to `toInt` must be FLOAT or INT, got %s", args[0].Type())
}

func toFloatBuiltin(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewException("Incorrect number of arguments. Got %d, expected 1", len(args))
	}

	switch arg := args[0].(type) {
	case *object.Integer:
		return &object.Float{Value: float64(arg.Value)}
	case *object.Float:
		return arg
	}

	return object.NewException("Argument to `toFloat` must be FLOAT or INT, got %s", args[0].Type())
}

func makeIsTypeBuiltin(t object.ObjectType) object.BuiltinFunction {
	return func(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
		if len(args) != 1 {
			return object.NewException("Type check requires one argument. Got %d", len(args))
		}

		return object.NativeBoolToBooleanObj(args[0].Type() == t)
	}
}

func varTypeBuiltin(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	return object.MakeStringObj(args[0].Type().String())
}

func getErrorVal(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	if ac := moduleutils.CheckArgs("errorVal", 1, args...); ac != nil {
		return ac
	}

	switch arg := args[0].(type) {
	case *object.Error:
		return object.MakeStringObj(arg.Message)
	case *object.Exception:
		return object.MakeStringObj(arg.Message)
	}

	return object.MakeStringObj("")
}

type resource interface {
	ResourceID() string
}

func getResourceID(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	if ac := moduleutils.CheckArgs("resourceID", 1, args...); ac != nil {
		return ac
	}

	if !object.ObjectIs(args[0], object.ResourceObj) {
		return object.NewException("cannot retrieve resource ID from non-resource object")
	}

	arg, ok := args[0].(resource)
	if !ok {
		return object.NewPanic("object marked as a ResourceObj doesn't implement resource ID interface")
	}

	return object.MakeStringObj(arg.ResourceID())
}

func toStringBuiltin(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewException("toString expects 1 argument. Got %d", len(args))
	}

	converted := ""

	switch arg := args[0].(type) {
	case *object.String:
		converted = arg.String()
	case *object.Float:
		converted = strconv.FormatFloat(arg.Value, 'G', -1, 64)
	case *object.Integer:
		converted = strconv.FormatInt(arg.Value, 10)
	case *object.Boolean:
		converted = strconv.FormatBool(arg.Value)
	case *object.Null:
		converted = "nil"
	default:
		converted = arg.Inspect()
	}

	return object.MakeStringObj(converted)
}

func parseIntBuiltin(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewException("parseInt expects 1 argument. Got %d", len(args))
	}

	str, ok := args[0].(*object.String)
	if !ok {
		return object.NewException("parseInt expected a string, got %s", args[0].Type().String())
	}

	i, err := strconv.ParseInt(str.String(), 10, 64)
	if err != nil {
		return object.NullConst
	}

	return object.MakeIntObj(i)
}

func parseFloatBuiltin(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewException("parseFloat expects 1 argument. Got %d", len(args))
	}

	str, ok := args[0].(*object.String)
	if !ok {
		return object.NewException("parseFloat expected a string, got %s", args[0].Type().String())
	}

	f, err := strconv.ParseFloat(str.String(), 64)
	if err != nil {
		return object.NullConst
	}

	return &object.Float{Value: f}
}

func isDefinedBuiltin(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	if ac := moduleutils.CheckArgs("isDefined", 1, args...); ac != nil {
		return ac
	}

	ident, ok := args[0].(*object.String)
	if !ok {
		return object.NewException("isDefined expects a string, got %s", args[0].Type().String())
	}

	_, ok = env.Get(ident.String())
	return object.NativeBoolToBooleanObj(ok)
}
