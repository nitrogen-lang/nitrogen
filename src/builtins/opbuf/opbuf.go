package string

import (
	"bytes"
	"io"

	"github.com/nitrogen-lang/nitrogen/src/object"
	"github.com/nitrogen-lang/nitrogen/src/vm"
)

const (
	moduleName      = "std/opbuf"
	instanceVarName = "opbuf-oldout"
)

func init() {
	vm.RegisterModule(moduleName, &object.Module{
		Name: moduleName,
		Methods: map[string]object.BuiltinFunction{
			"start":      start,
			"clear":      clear,
			"flush":      flush,
			"get":        get,
			"stopAndGet": stopAndGet,
			"stop":       stop,
			"isStarted":  isStarted,
		},
		Vars: map[string]object.Object{},
	})
}

func start(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	theVM := interpreter.(*vm.VirtualMachine)

	if theVM.HasInstanceVar(instanceVarName) {
		return object.NewException("Output buffering is already started")
	}

	theVM.SetInstanceVar(instanceVarName, theVM.Settings.Stdout)
	theVM.Settings.Stdout = &bytes.Buffer{}

	return nil
}

func stopAndGet(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	theVM := interpreter.(*vm.VirtualMachine)

	if !theVM.HasInstanceVar(instanceVarName) {
		return object.NewException("Output buffering is not started")
	}

	buf := theVM.Settings.Stdout.(*bytes.Buffer)
	theVM.Settings.Stdout = theVM.GetInstanceVar(instanceVarName).(io.Writer)
	theVM.RemoveInstanceVar(instanceVarName)

	return object.MakeStringObj(buf.String())
}

func stop(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	theVM := interpreter.(*vm.VirtualMachine)

	if !theVM.HasInstanceVar(instanceVarName) {
		return object.NewException("Output buffering is not started")
	}

	theVM.Settings.Stdout = theVM.GetInstanceVar(instanceVarName).(io.Writer)
	theVM.RemoveInstanceVar(instanceVarName)

	return nil
}

func clear(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	theVM := interpreter.(*vm.VirtualMachine)

	if !theVM.HasInstanceVar(instanceVarName) {
		return object.NewException("Output buffering is not started")
	}

	buf := theVM.Settings.Stdout.(*bytes.Buffer)
	buf.Reset()

	return nil
}

func flush(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	theVM := interpreter.(*vm.VirtualMachine)

	if !theVM.HasInstanceVar(instanceVarName) {
		return object.NewException("Output buffering is not started")
	}

	buf := theVM.Settings.Stdout.(*bytes.Buffer)
	oldOut := theVM.GetInstanceVar(instanceVarName).(io.Writer)
	_, err := io.Copy(oldOut, buf)
	if err != nil {
		return object.NewException(err.Error())
	}

	return nil
}

func get(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	theVM := interpreter.(*vm.VirtualMachine)

	if !theVM.HasInstanceVar(instanceVarName) {
		return object.NewException("Output buffering is not started")
	}

	buf := theVM.Settings.Stdout.(*bytes.Buffer)

	return object.MakeStringObj(buf.String())
}

func isStarted(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	theVM := interpreter.(*vm.VirtualMachine)
	return object.NativeBoolToBooleanObj(theVM.HasInstanceVar(instanceVarName))
}
