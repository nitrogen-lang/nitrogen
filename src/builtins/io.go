package eval

import (
	"fmt"

	"github.com/nitrogen-lang/nitrogen/src/eval"
	"github.com/nitrogen-lang/nitrogen/src/object"
)

func init() {
	eval.RegisterBuiltin("print", printBuiltin)
	eval.RegisterBuiltin("println", printlnBuiltin)
	eval.RegisterBuiltin("printenv", printEnvBuiltin)
}

func printBuiltin(env *object.Environment, args ...object.Object) object.Object {
	for _, arg := range args {
		fmt.Print(arg.Inspect())
	}
	return object.NULL
}

func printlnBuiltin(env *object.Environment, args ...object.Object) object.Object {
	printBuiltin(env, args...)
	fmt.Print("\n")
	return object.NULL
}

func printEnvBuiltin(env *object.Environment, args ...object.Object) object.Object {
	env.Print("")
	return object.NULL
}
