package main

import (
	"bytes"
	"os/exec"

	"github.com/nitrogen-lang/nitrogen/src/eval"

	"github.com/nitrogen-lang/nitrogen/src/moduleutils"
	"github.com/nitrogen-lang/nitrogen/src/object"
)

func init() {
	eval.RegisterModule("os", &object.Module{
		Name: "os",
		Methods: map[string]object.BuiltinFunction{
			"system": runSystem,
			"exec":   runSystemPT,
		},
		Vars: map[string]object.Object{
			"name": object.MakeStringObj(ModuleName),
		},
	})
}

func main() {}

var ModuleName = "os"

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
			cmdArgs[i] = arg.Value
		}
	}

	var (
		stdout bytes.Buffer
		stderr bytes.Buffer
	)

	cmd := exec.Command(cmdName.Value, cmdArgs...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return object.NewError("Error executing command %s", err.Error())
	}

	return &object.Array{
		Elements: []object.Object{
			&object.String{Value: stdout.String()},
			&object.String{Value: stderr.String()},
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
			cmdArgs[i] = arg.Value
		}
	}

	cmd := exec.Command(cmdName.Value, cmdArgs...)
	cmd.Stdin = interpreter.GetStdin()
	cmd.Stdout = interpreter.GetStdout()
	cmd.Stderr = interpreter.GetStderr()

	if err := cmd.Run(); err != nil {
		return object.NewError("Error executing command %s", err.Error())
	}
	return object.NullConst
}
