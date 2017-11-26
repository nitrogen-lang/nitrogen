// +build !linux,!darwin !cgo

package imports

import (
	"github.com/nitrogen-lang/nitrogen/src/eval"
	"github.com/nitrogen-lang/nitrogen/src/moduleutils"
	"github.com/nitrogen-lang/nitrogen/src/object"
	"github.com/nitrogen-lang/nitrogen/src/vm"
)

func init() {
	eval.RegisterBuiltin("module", importModule)
	eval.RegisterBuiltin("modulesSupported", moduleSupport)

	vm.RegisterBuiltin("module", moduleutils.VMBuiltinWrapper(importModule))
	vm.RegisterBuiltin("modulesSupported", moduleutils.VMBuiltinWrapper(moduleSupport))
}

func importModule(i object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	return object.NewException("Shared object modules are not supported in this build")
}

func moduleSupport(i object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	return object.NativeBoolToBooleanObj(false)
}
