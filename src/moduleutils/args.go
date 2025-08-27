package moduleutils

import (
	"github.com/nitrogen-lang/nitrogen/src/elemental/object"
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
