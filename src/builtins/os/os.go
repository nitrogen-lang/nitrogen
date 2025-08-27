package os

import (
	"bytes"
	"os/exec"

	"github.com/nitrogen-lang/nitrogen/src/elemental/object"
	"github.com/nitrogen-lang/nitrogen/src/elemental/vm"
	"github.com/nitrogen-lang/nitrogen/src/moduleutils"
)

var (
	commandArgs *object.Array
)

func init() {
	vm.RegisterNative("std.os.env", env)
	vm.RegisterNative("std.os.argv", argv)
	vm.RegisterNative("std.os.system", runSystem)
	vm.RegisterNative("std.os.exec", runSystemPT)
}

// SetCmdArgs sets the command line arguments array.
func SetCmdArgs(args *object.Array) {
	commandArgs = args
}

func env(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	return (interpreter.(*vm.VirtualMachine)).GetInstanceVar("os.env")
}

func argv(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	return commandArgs
}

func runSystem(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	if ac := moduleutils.CheckMinArgs("system", 1, args...); ac != nil {
		return ac
	}

	cmdName, ok := args[0].(*object.String)
	if !ok {
		return object.NewException("system expected a string, got %s", args[0].Type().String())
	}

	var cmdArgs []string
	if len(args) > 1 {
		cmdArgsArray, ok := args[1].(*object.Array)
		if !ok {
			return object.NewException("system expected an array, got %s", args[0].Type().String())
		}

		cmdArgs = make([]string, len(cmdArgsArray.Elements))
		for i, element := range cmdArgsArray.Elements {
			arg, ok := element.(*object.String)
			if !ok {
				return object.NewException("system arguments must be a string %s", element.Inspect())
			}
			cmdArgs[i] = arg.String()
		}
	}

	var (
		stdout bytes.Buffer
		stderr bytes.Buffer
	)

	cmd := exec.Command(cmdName.String(), cmdArgs...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return object.NewError("Error executing command %s", err.Error())
	}

	return &object.Array{
		Elements: []object.Object{
			object.MakeStringObj(stdout.String()),
			object.MakeStringObj(stderr.String()),
		},
	}
}

func runSystemPT(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	if ac := moduleutils.CheckMinArgs("system", 1, args...); ac != nil {
		return ac
	}

	cmdName, ok := args[0].(*object.String)
	if !ok {
		return object.NewException("system expected a string, got %s", args[0].Type().String())
	}

	var cmdArgs []string
	if len(args) > 1 {
		cmdArgsArray, ok := args[1].(*object.Array)
		if !ok {
			return object.NewException("system expected an array, got %s", args[0].Type().String())
		}

		cmdArgs = make([]string, len(cmdArgsArray.Elements))
		for i, element := range cmdArgsArray.Elements {
			arg, ok := element.(*object.String)
			if !ok {
				return object.NewException("system arguments must be a string %s", element.Inspect())
			}
			cmdArgs[i] = arg.String()
		}
	}

	cmd := exec.Command(cmdName.String(), cmdArgs...)
	cmd.Stdin = interpreter.GetStdin()
	cmd.Stdout = interpreter.GetStdout()
	cmd.Stderr = interpreter.GetStderr()

	if err := cmd.Run(); err != nil {
		return object.NewError("Error executing command %s", err.Error())
	}
	return object.NullConst
}
