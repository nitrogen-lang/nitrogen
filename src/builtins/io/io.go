package io

import (
	"bufio"
	"fmt"
	"os"

	"github.com/nitrogen-lang/nitrogen/src/elemental/object"
	"github.com/nitrogen-lang/nitrogen/src/elemental/vm"
	"github.com/nitrogen-lang/nitrogen/src/moduleutils"
)

func init() {
	vm.RegisterNative("std.preamble.io.print", printBuiltin)
	vm.RegisterNative("std.preamble.io.printlnb", printBinaryBuiltin)
	vm.RegisterNative("std.preamble.io.println", printlnBuiltin)
	vm.RegisterNative("std.preamble.io.printerr", printerrBuiltin)
	vm.RegisterNative("std.preamble.io.printerrln", printerrlnBuiltin)
	vm.RegisterNative("std.preamble.io.printenv", printEnvBuiltin)
	vm.RegisterNative("std.preamble.io.varDump", varDump)
	vm.RegisterNative("std.preamble.os.exit", exitScript)

	vm.RegisterNative("std.preamble.io.readline", readLineBuiltin)
}

func varDump(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	if ac := moduleutils.CheckMinArgs("varDump", 1, args...); ac != nil {
		return ac
	}

	return printBuiltin(interpreter, env, args...)
}

func exitScript(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	code := 0

	if len(args) > 0 {
		c, ok := args[0].(*object.Integer)
		if !ok {
			return object.NewException("exit expected an int. Got %s", args[0].Type().String())
		}
		code = int(c.Value)
	}

	machine := interpreter.(*vm.VirtualMachine)
	machine.Exit(code)
	return object.NullConst
}

func printBuiltin(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	out := interpreter.GetStdout()
	for _, arg := range args {
		if instance, ok := arg.(*vm.VMInstance); ok {
			machine := interpreter.(*vm.VirtualMachine)
			toString := instance.GetBoundMethod("toString")
			if toString != nil {
				machine.CallFunction(0, toString, true, nil, false)
				printBuiltin(interpreter, env, machine.PopStack())
				continue
			}
		}

		fmt.Fprint(out, arg.Inspect())
	}
	return object.NullConst
}

func printlnBuiltin(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	printBuiltin(interpreter, env, args...)
	fmt.Fprint(interpreter.GetStdout(), "\n")
	return object.NullConst
}

func printerrBuiltin(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	for _, arg := range args {
		if instance, ok := arg.(*vm.VMInstance); ok {
			machine := interpreter.(*vm.VirtualMachine)
			toString := instance.GetBoundMethod("toString")
			if toString != nil {
				machine.CallFunction(0, toString, true, nil, false)
				printBuiltin(interpreter, env, machine.PopStack())
				continue
			}
		}

		fmt.Fprint(interpreter.GetStderr(), arg.Inspect())
	}
	return object.NullConst
}

func printerrlnBuiltin(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	printerrBuiltin(interpreter, env, args...)
	fmt.Fprint(interpreter.GetStderr(), "\n")
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
	return object.MakeStringObj(text[:len(text)-1])
}
