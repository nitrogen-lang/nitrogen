package eval

import (
	"github.com/nitrogen-lang/nitrogen/src/object"
)

func applyFunction(fn object.Object, args []object.Object, env *object.Environment) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		if len(fn.Parameters) != len(args) {
			return object.NewError("Not enough parameters to call function %s", fn.Name)
		}
		extendedEnv := extendFunctionEnv(fn, args)
		evaled := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaled)
	case *object.Builtin:
		return fn.Fn(env, args...)
	}

	return object.NewError("%s is not a function", fn.Type())
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnv(fn.Env)

	for i, param := range fn.Parameters {
		env.Create(param.Value, args[i])
	}

	return env
}
