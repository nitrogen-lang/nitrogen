package classes

import (
	"github.com/nitrogen-lang/nitrogen/src/moduleutils"
	"github.com/nitrogen-lang/nitrogen/src/object"
	"github.com/nitrogen-lang/nitrogen/src/vm"
)

func init() {
	vm.RegisterBuiltin("instanceOf", vmInstanceOf)
	vm.RegisterBuiltin("classOf", vmClassOf)
}

func vmInstanceOf(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	if ac := moduleutils.CheckArgs("instanceOf", 2, args...); ac != nil {
		return ac
	}

	instance, ok := args[0].(*vm.VMInstance)
	if !ok {
		return object.NativeBoolToBooleanObj(false)
	}

	switch class := args[1].(type) {
	case *object.String:
		return object.NativeBoolToBooleanObj(vm.InstanceOf(class.String(), instance))
	case *vm.VMClass:
		return object.NativeBoolToBooleanObj(vm.InstanceOf(class.Name, instance))
	}

	return object.NewException("instanceOf expected a class or string for second argument")
}

func vmClassOf(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	if ac := moduleutils.CheckArgs("classOf", 1, args...); ac != nil {
		return ac
	}

	instance, ok := args[0].(*vm.VMInstance)
	if !ok {
		return object.MakeStringObj("")
	}

	return object.MakeStringObj(instance.Class.Name)
}
