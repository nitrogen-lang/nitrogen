package opbuf

import (
	"bytes"
	"io"

	"github.com/nitrogen-lang/nitrogen/src/object"
	"github.com/nitrogen-lang/nitrogen/src/vm"
)

var (
	oldWriter io.Writer
)

func init() {
	vm.RegisterNative("std.opbuf.start", start)
	vm.RegisterNative("std.opbuf.clear", clear)
	vm.RegisterNative("std.opbuf.flush", flush)
	vm.RegisterNative("std.opbuf.get", get)
	vm.RegisterNative("std.opbuf.stopAndGet", stopAndGet)
	vm.RegisterNative("std.opbuf.stop", stop)
	vm.RegisterNative("std.opbuf.isStarted", isStarted)
}

func start(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	theVM := interpreter.(*vm.VirtualMachine)

	if oldWriter != nil {
		return object.NewException("Output buffering is already started")
	}

	oldWriter = theVM.Settings.Stdout
	theVM.Settings.Stdout = &bytes.Buffer{}

	return nil
}

func stopAndGet(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	theVM := interpreter.(*vm.VirtualMachine)

	if oldWriter == nil {
		return object.NewException("Output buffering is not started")
	}

	buf := theVM.Settings.Stdout.(*bytes.Buffer)
	theVM.Settings.Stdout = oldWriter
	oldWriter = nil

	return object.MakeStringObj(buf.String())
}

func stop(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	theVM := interpreter.(*vm.VirtualMachine)

	if oldWriter == nil {
		return object.NewException("Output buffering is not started")
	}

	theVM.Settings.Stdout = oldWriter
	oldWriter = nil

	return nil
}

func clear(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	theVM := interpreter.(*vm.VirtualMachine)

	if oldWriter == nil {
		return object.NewException("Output buffering is not started")
	}

	buf := theVM.Settings.Stdout.(*bytes.Buffer)
	buf.Reset()

	return nil
}

func flush(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	theVM := interpreter.(*vm.VirtualMachine)

	if oldWriter == nil {
		return object.NewException("Output buffering is not started")
	}

	buf := theVM.Settings.Stdout.(*bytes.Buffer)
	oldOut := oldWriter
	_, err := io.Copy(oldOut, buf)
	if err != nil {
		return object.NewException(err.Error())
	}

	return nil
}

func get(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	theVM := interpreter.(*vm.VirtualMachine)

	if oldWriter == nil {
		return object.NewException("Output buffering is not started")
	}

	buf := theVM.Settings.Stdout.(*bytes.Buffer)

	return object.MakeStringObj(buf.String())
}

func isStarted(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	return object.NativeBoolToBooleanObj(oldWriter != nil)
}
