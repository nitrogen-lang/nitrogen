package main

import (
	"os/exec"

	"github.com/nitrogen-lang/nitrogen/src/eval"
	"github.com/nitrogen-lang/nitrogen/src/moduleutils"
	"github.com/nitrogen-lang/nitrogen/src/object"
)

func init() {
	eval.RegisterBuiltin("system", runSystem)
}

func main() {}

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

	out, err := exec.Command(cmdName.Value, cmdArgs...).Output()
	if err != nil {
		return object.NewException("Error executing command %s", err.Error())
	}
	return &object.String{Value: string(out)}
}
