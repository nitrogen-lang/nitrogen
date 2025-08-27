package errors

import (
	"github.com/nitrogen-lang/nitrogen/src/elemental/object"
	"github.com/nitrogen-lang/nitrogen/src/elemental/vm"
	"github.com/nitrogen-lang/nitrogen/src/moduleutils"
)

func init() {
	vm.RegisterNative("std.preamble.main.error", vmMakeError)
}

func vmMakeError(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	if ac := moduleutils.CheckArgs("error", 1, args...); ac != nil {
		return ac
	}

	return &object.Error{Message: args[0].Inspect()}
}
