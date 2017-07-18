package builtins

import "github.com/nitrogen-lang/nitrogen/src/object"

func checkArgs(name string, expected int, args ...object.Object) *object.Error {
	if len(args) != expected {
		return object.NewError("%s expects %d argument(s). Got %d", name, expected, len(args))
	}
	return nil
}
