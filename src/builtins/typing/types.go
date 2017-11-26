package typing

import (
	"strconv"

	"github.com/nitrogen-lang/nitrogen/src/eval"
	"github.com/nitrogen-lang/nitrogen/src/moduleutils"
	"github.com/nitrogen-lang/nitrogen/src/object"
	"github.com/nitrogen-lang/nitrogen/src/vm"
)

func init() {
	// Register with interpreter
	eval.RegisterBuiltin("toInt", toIntBuiltin)
	eval.RegisterBuiltin("toFloat", toFloatBuiltin)
	eval.RegisterBuiltin("toString", toStringBuiltin)

	eval.RegisterBuiltin("parseInt", parseIntBuiltin)
	eval.RegisterBuiltin("parseFloat", parseFloatBuiltin)

	eval.RegisterBuiltin("varType", varTypeBuiltin)
	eval.RegisterBuiltin("isDefined", isDefinedBuiltin)
	eval.RegisterBuiltin("isFloat", makeIsTypeBuiltin(object.FloatObj))
	eval.RegisterBuiltin("isInt", makeIsTypeBuiltin(object.IntergerObj))
	eval.RegisterBuiltin("isBool", makeIsTypeBuiltin(object.BooleanObj))
	eval.RegisterBuiltin("isNull", makeIsTypeBuiltin(object.NullObj))
	eval.RegisterBuiltin("isFunc", makeIsTypeBuiltin(object.FunctionObj))
	eval.RegisterBuiltin("isString", makeIsTypeBuiltin(object.StringObj))
	eval.RegisterBuiltin("isArray", makeIsTypeBuiltin(object.ArrayObj))
	eval.RegisterBuiltin("isMap", makeIsTypeBuiltin(object.HashObj))
	eval.RegisterBuiltin("isError", makeIsTypeBuiltin(object.ErrorObj))
	eval.RegisterBuiltin("isClass", makeIsTypeBuiltin(object.ClassObj))
	eval.RegisterBuiltin("isInstance", makeIsTypeBuiltin(object.InstanceObj))

	eval.RegisterBuiltin("errorVal", getErrorVal)

	// Register with virual machine
	vm.RegisterBuiltin("toInt", moduleutils.VMBuiltinWrapper(toIntBuiltin))
	vm.RegisterBuiltin("toFloat", moduleutils.VMBuiltinWrapper(toFloatBuiltin))
	vm.RegisterBuiltin("toString", moduleutils.VMBuiltinWrapper(toStringBuiltin))

	vm.RegisterBuiltin("parseInt", moduleutils.VMBuiltinWrapper(parseIntBuiltin))
	vm.RegisterBuiltin("parseFloat", moduleutils.VMBuiltinWrapper(parseFloatBuiltin))

	vm.RegisterBuiltin("varType", moduleutils.VMBuiltinWrapper(varTypeBuiltin))
	vm.RegisterBuiltin("isDefined", isDefinedBuiltinVM)
	vm.RegisterBuiltin("isFloat", moduleutils.VMBuiltinWrapper(makeIsTypeBuiltin(object.FloatObj)))
	vm.RegisterBuiltin("isInt", moduleutils.VMBuiltinWrapper(makeIsTypeBuiltin(object.IntergerObj)))
	vm.RegisterBuiltin("isBool", moduleutils.VMBuiltinWrapper(makeIsTypeBuiltin(object.BooleanObj)))
	vm.RegisterBuiltin("isNull", moduleutils.VMBuiltinWrapper(makeIsTypeBuiltin(object.NullObj)))
	vm.RegisterBuiltin("isFunc", moduleutils.VMBuiltinWrapper(makeIsTypeBuiltin(object.FunctionObj)))
	vm.RegisterBuiltin("isString", moduleutils.VMBuiltinWrapper(makeIsTypeBuiltin(object.StringObj)))
	vm.RegisterBuiltin("isArray", moduleutils.VMBuiltinWrapper(makeIsTypeBuiltin(object.ArrayObj)))
	vm.RegisterBuiltin("isMap", moduleutils.VMBuiltinWrapper(makeIsTypeBuiltin(object.HashObj)))
	vm.RegisterBuiltin("isError", moduleutils.VMBuiltinWrapper(makeIsTypeBuiltin(object.ErrorObj)))
	vm.RegisterBuiltin("isClass", moduleutils.VMBuiltinWrapper(makeIsTypeBuiltin(object.ClassObj)))
	vm.RegisterBuiltin("isInstance", moduleutils.VMBuiltinWrapper(makeIsTypeBuiltin(object.InstanceObj)))

	vm.RegisterBuiltin("errorVal", moduleutils.VMBuiltinWrapper(getErrorVal))
}

func toIntBuiltin(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewException("Incorrect number of arguments. Got %d, expected 1", len(args))
	}

	switch arg := args[0].(type) {
	case *object.Integer:
		return arg
	case *object.Float:
		return &object.Integer{Value: int64(arg.Value)}
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
	return &object.String{Value: args[0].Type().String()}
}

func getErrorVal(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	if ac := moduleutils.CheckArgs("errorVal", 1, args...); ac != nil {
		return ac
	}

	switch arg := args[0].(type) {
	case *object.Error:
		return &object.String{Value: arg.Message}
	case *object.Exception:
		return &object.String{Value: arg.Message}
	}

	return &object.String{Value: ""}
}

func toStringBuiltin(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewException("toString expects 1 argument. Got %d", len(args))
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

func parseIntBuiltin(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewException("parseInt expects 1 argument. Got %d", len(args))
	}

	str, ok := args[0].(*object.String)
	if !ok {
		return object.NewException("parseInt expected a string, got %s", args[0].Type().String())
	}

	i, err := strconv.ParseInt(str.Value, 10, 64)
	if err != nil {
		return object.NullConst
	}

	return &object.Integer{Value: i}
}

func parseFloatBuiltin(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewException("parseFloat expects 1 argument. Got %d", len(args))
	}

	str, ok := args[0].(*object.String)
	if !ok {
		return object.NewException("parseFloat expected a string, got %s", args[0].Type().String())
	}

	f, err := strconv.ParseFloat(str.Value, 64)
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

	_, ok = env.Get(ident.Value)
	return object.NativeBoolToBooleanObj(ok)
}

func isDefinedBuiltinVM(machine object.Interpreter, args ...object.Object) object.Object {
	if ac := moduleutils.CheckArgs("isDefined", 1, args...); ac != nil {
		return ac
	}

	ident, ok := args[0].(*object.String)
	if !ok {
		return object.NewException("isDefined expects a string, got %s", args[0].Type().String())
	}

	_, ok = machine.(*vm.VirtualMachine).CurrentFrame().Env.Get(ident.Value)
	return object.NativeBoolToBooleanObj(ok)
}
