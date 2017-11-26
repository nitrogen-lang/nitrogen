package io

import (
	"bufio"
	"fmt"
	"os"

	"github.com/nitrogen-lang/nitrogen/src/eval"
	"github.com/nitrogen-lang/nitrogen/src/moduleutils"
	"github.com/nitrogen-lang/nitrogen/src/object"
	"github.com/nitrogen-lang/nitrogen/src/vm"
)

func init() {
	eval.RegisterBuiltin("print", printBuiltin)
	eval.RegisterBuiltin("printlnb", printBinaryBuiltin)
	eval.RegisterBuiltin("println", printlnBuiltin)
	eval.RegisterBuiltin("printenv", printEnvBuiltin)

	eval.RegisterBuiltin("readline", readLineBuiltin)

	vm.RegisterBuiltin("print", moduleutils.VMBuiltinWrapper(printBuiltin))
	vm.RegisterBuiltin("printlnb", moduleutils.VMBuiltinWrapper(printBinaryBuiltin))
	vm.RegisterBuiltin("println", moduleutils.VMBuiltinWrapper(printlnBuiltin))
	vm.RegisterBuiltin("printenv", printEnvBuiltinVM)

	vm.RegisterBuiltin("readline", moduleutils.VMBuiltinWrapper(readLineBuiltin))
}

func printBuiltin(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	for _, arg := range args {
		fmt.Fprint(interpreter.GetStdout(), arg.Inspect())
	}
	return object.NullConst
}

func printlnBuiltin(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	printBuiltin(interpreter, env, args...)
	fmt.Fprint(interpreter.GetStdout(), "\n")
	return object.NullConst
}

func printBinaryBuiltin(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	fmt.Fprintf(interpreter.GetStdout(), "%b\n", args[0].(*object.Integer).Value)
	return object.NullConst
}

func printEnvBuiltin(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	env.Print("")
	return object.NullConst
}

func printEnvBuiltinVM(machine object.Interpreter, args ...object.Object) object.Object {
	machine.(*vm.VirtualMachine).CurrentFrame().Env.Print("")
	return object.NullConst
}

func readLineBuiltin(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	if len(args) > 1 {
		return object.NewException("readline only accepts up to one argument. Got %d", len(args))
	}

	if len(args) == 1 {
		prompt, ok := args[0].(*object.String)
		if !ok {
			return object.NewException("readline expects a string for the first arguemnt, got %s", args[0].Type().String())
		}
		fmt.Fprint(interpreter.GetStdout(), prompt.Value)
	}

	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')

	// Return read line without the ending newline byte
	return &object.String{Value: text[:len(text)-1]}
}
