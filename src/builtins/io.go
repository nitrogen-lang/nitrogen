package builtins

import (
	"bufio"
	"fmt"
	"os"

	"github.com/nitrogen-lang/nitrogen/src/eval"
	"github.com/nitrogen-lang/nitrogen/src/object"
)

func init() {
	eval.RegisterBuiltin("print", printBuiltin)
	eval.RegisterBuiltin("println", printlnBuiltin)
	eval.RegisterBuiltin("printenv", printEnvBuiltin)

	eval.RegisterBuiltin("readline", readLineBuiltin)
}

func printBuiltin(env *object.Environment, args ...object.Object) object.Object {
	for _, arg := range args {
		fmt.Print(arg.Inspect())
	}
	return object.NullConst
}

func printlnBuiltin(env *object.Environment, args ...object.Object) object.Object {
	printBuiltin(env, args...)
	fmt.Print("\n")
	return object.NullConst
}

func printEnvBuiltin(env *object.Environment, args ...object.Object) object.Object {
	env.Print("")
	return object.NullConst
}

func readLineBuiltin(env *object.Environment, args ...object.Object) object.Object {
	if len(args) > 1 {
		return object.NewException("readline only accepts up to one argument. Got %d", len(args))
	}

	if len(args) == 1 {
		prompt, ok := args[0].(*object.String)
		if !ok {
			return object.NewException("readline expects a string for the first arguemnt, got %s", args[0].Type().String())
		}
		fmt.Print(prompt.Value)
	}

	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')

	// Return read line without the ending newline byte
	return &object.String{Value: text[:len(text)-1]}
}
