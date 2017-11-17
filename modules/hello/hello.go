package main

import (
	"fmt"

	"github.com/nitrogen-lang/nitrogen/src/eval"
	"github.com/nitrogen-lang/nitrogen/src/object"
)

func init() {
	eval.RegisterBuiltin("hello", printHello)
}

func main() {}

func printHello(env *object.Environment, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewError("Incorrect number of arguments. Got %d, expected 1", len(args))
	}

	if args[0].Type() != object.STRING_OBJ {
		return object.NewError("Argument to `hello` must be a string, got %s", args[0].Type())
	}

	fmt.Printf("Hello %s!", args[0].Inspect())
	return object.NULL
}
