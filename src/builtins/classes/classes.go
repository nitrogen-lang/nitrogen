package classes

import (
	"github.com/nitrogen-lang/nitrogen/src/elemental/object"
	"github.com/nitrogen-lang/nitrogen/src/elemental/vm"
	"github.com/nitrogen-lang/nitrogen/src/moduleutils"
)

func init() {
	vm.RegisterNative("std.preamble.main.instanceOf", vmInstanceOf)
	vm.RegisterNative("std.preamble.main.classOf", vmClassOf)
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
