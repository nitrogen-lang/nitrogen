package eval

import (
	"github.com/nitrogen-lang/nitrogen/src/object"
)

func (i *Interpreter) applyFunction(fn object.Object, args []object.Object, env *object.Environment) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		if len(args) < len(fn.Parameters) {
			return object.NewException("Not enough parameters to call function %s", fn.Name)
		}
		extendedEnv := i.extendFunctionEnv(fn, fn.Env, args)
		oldInstance := i.currentInstance
		i.currentInstance = fn.Instance
		evaled := i.Eval(fn.Body, extendedEnv)
		i.currentInstance = oldInstance
		return unwrapReturnValue(evaled)
	case *object.Builtin:
		return fn.Fn(i, env, args...)
	case *object.BuiltinMethod:
		oldInstance := i.currentInstance
		i.currentInstance = fn.Instance
		result := fn.Fn(i, i.currentInstance, env, args...)
		i.currentInstance = oldInstance
		return result
	case *object.Class: // Class init function
		init := fn.GetMethod("init")
		if init == nil {
			return object.NullConst
		}

		if initBuiltin, ok := init.(*object.BuiltinMethod); ok {
			return initBuiltin.Fn(i, i.currentInstance, env, args...)
		}

		initFn := init.(*object.Function)

		if len(args) < len(initFn.Parameters) {
			return object.NewException("Not enough parameters to call class initializer %s", fn.Name)
		}
		extendedEnv := i.extendFunctionEnv(initFn, initFn.Env, args)
		extendedEnv.SetParent(env)
		if fn.Parent != nil {
			extendedEnv.CreateConst("parent", fn.Parent)
		}
		evaled := i.Eval(initFn.Body, extendedEnv)
		return unwrapReturnValue(evaled)
	}

	return object.NewException("%s is not a function", fn.Type())
}

func (i *Interpreter) applyFunctionDirect(fn object.Object, args []object.Object, env *object.Environment) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		if len(args) < len(fn.Parameters) {
			return object.NewException("Not enough parameters to call function %s", fn.Name)
		}
		extendedEnv := i.extendFunctionEnv(fn, env, args)
		evaled := i.Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaled)
	case *object.Builtin:
		return fn.Fn(i, env, args...)
	case *object.BuiltinMethod:
		return fn.Fn(i, i.currentInstance, env, args...)
	}

	return object.NewException("%s is not a function", fn.Type())
}

func (i *Interpreter) extendFunctionEnv(fn *object.Function, outer *object.Environment, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnv(outer)

	for i, param := range fn.Parameters {
		env.Create(param.Value, args[i])
	}

	// The "args" variable hold all extra parameters beyond those defined
	// by the function at runtime. "args[0]" is the first EXTRA parameter
	// after those that were defined have been bound.

	// Although the elements of the args variable could be reassigned,
	// I'm trying to discourage it by at least making the variable itself
	// a constant. Trying to indicate "please don't mess with it". Mainly
	// this is so the variable isn't overwritten accidentally.
	if len(args) > len(fn.Parameters) {
		env.CreateConst("args", &object.Array{
			Elements: args[len(fn.Parameters):],
		})
	} else {
		// The idea is for functions to call "len(args)" to check for
		// anything extra. "len(nil)" returns 0.
		env.CreateConst("args", object.NullConst)
	}

	return env
}
