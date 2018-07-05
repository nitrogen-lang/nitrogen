package string

import (
	"bytes"
	"io"

	"github.com/nitrogen-lang/nitrogen/src/object"
	"github.com/nitrogen-lang/nitrogen/src/vm"
)

var (
	moduleName = "stdlib/opbuf"
	oldOut     io.Writer
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
	if oldOut != nil {
		return object.NewException("Output buffering is already started")
	}

	theVM, _ := interpreter.(*vm.VirtualMachine)
	oldOut = theVM.Settings.Stdout
	theVM.Settings.Stdout = &bytes.Buffer{}

	return nil
}

func stopAndGet(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	if oldOut == nil {
		return object.NewException("Output buffering is not started")
	}

	theVM := interpreter.(*vm.VirtualMachine)
	buf := theVM.Settings.Stdout.(*bytes.Buffer)
	theVM.Settings.Stdout = oldOut
	oldOut = nil

	return object.MakeStringObj(buf.String())
}

func stop(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	if oldOut == nil {
		return object.NewException("Output buffering is not started")
	}

	theVM := interpreter.(*vm.VirtualMachine)
	theVM.Settings.Stdout = oldOut
	oldOut = nil

	return nil
}

func clear(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	if oldOut == nil {
		return object.NewException("Output buffering is not started")
	}

	theVM := interpreter.(*vm.VirtualMachine)
	buf := theVM.Settings.Stdout.(*bytes.Buffer)
	buf.Reset()

	return nil
}

func flush(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	if oldOut == nil {
		return object.NewException("Output buffering is not started")
	}

	theVM := interpreter.(*vm.VirtualMachine)
	buf := theVM.Settings.Stdout.(*bytes.Buffer)
	_, err := io.Copy(oldOut, buf)
	if err != nil {
		return object.NewException(err.Error())
	}

	return nil
}

func get(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	if oldOut == nil {
		return object.NewException("Output buffering is not started")
	}

	theVM := interpreter.(*vm.VirtualMachine)
	buf := theVM.Settings.Stdout.(*bytes.Buffer)

	return object.MakeStringObj(buf.String())
}

func isStarted(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	return object.NativeBoolToBooleanObj(oldOut != nil)
}
