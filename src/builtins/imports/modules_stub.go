// +build !linux,!darwin !cgo

package imports

import (
	"github.com/nitrogen-lang/nitrogen/src/object"
)

func moduleSupport(i object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	return object.NativeBoolToBooleanObj(false)
}

func importSharedModule(scriptPath string, required bool, interpreter object.Interpreter, env *object.Environment) object.Object {
	if required {
		return object.NewException("Shared object modules are not supported in this build")
	}
	return object.NewError("Shared object modules are not supported in this build")
}
