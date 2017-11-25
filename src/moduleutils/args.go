package moduleutils

import (
	"github.com/nitrogen-lang/nitrogen/src/object"
	"github.com/nitrogen-lang/nitrogen/src/vm"
)

// CheckArgs is a convience function to check if the correct number of arguments were supplied
func CheckArgs(name string, expected int, args ...object.Object) *object.Exception {
	if len(args) != expected {
		return object.NewException("%s expects %d argument(s). Got %d", name, expected, len(args))
	}
	return nil
}

// CheckMinArgs is a convience function to check if the minimum number of arguments were supplied
func CheckMinArgs(name string, expected int, args ...object.Object) *object.Exception {
	if len(args) < expected {
		return object.NewException("%s expects %d argument(s). Got %d", name, expected, len(args))
	}
	return nil
}

func VMBuiltinWrapper(fn object.BuiltinFunction) vm.VMBuiltinFunc {
	return func(vm object.Interpreter, args ...object.Object) object.Object {
		return fn(vm, nil, args...)
	}
}
