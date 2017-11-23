package classes

import (
	"github.com/nitrogen-lang/nitrogen/src/eval"
	"github.com/nitrogen-lang/nitrogen/src/moduleutils"
	"github.com/nitrogen-lang/nitrogen/src/object"
)

func init() {
	eval.RegisterBuiltin("is_a", instanceOf)
	eval.RegisterBuiltin("classOf", classOf)
}

func instanceOf(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	if ac := moduleutils.CheckArgs("is_a", 2, args...); ac != nil {
		return ac
	}

	instance, ok := args[0].(*object.Instance)
	if !ok {
		return object.NativeBoolToBooleanObj(false)
	}

	switch class := args[1].(type) {
	case *object.String:
		return object.NativeBoolToBooleanObj(object.InstanceOf(class.Value, instance))
	case *object.Class:
		return object.NativeBoolToBooleanObj(object.InstanceOf(class.Name, instance))
	}

	return object.NewException("is_a expected a class or string for second argument")
}

func classOf(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	if ac := moduleutils.CheckArgs("classOf", 1, args...); ac != nil {
		return ac
	}

	instance, ok := args[0].(*object.Instance)
	if !ok {
		return object.MakeStringObj("")
	}

	return object.MakeStringObj(instance.Class.Name)
}
