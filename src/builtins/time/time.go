package time

import (
	"time"

	"github.com/nitrogen-lang/nitrogen/src/elemental/object"
	"github.com/nitrogen-lang/nitrogen/src/elemental/vm"
	"github.com/nitrogen-lang/nitrogen/src/moduleutils"
)

func init() {
	vm.RegisterNative("std.time.now", timeNow)
	vm.RegisterNative("std.time.now_unix", timeNowS)
	vm.RegisterNative("std.time.now_ms", timeNowMs)
	vm.RegisterNative("std.time.now_ns", timeNowNs)

	vm.RegisterNativeMethod("std.time.Time.format", timeCFormat, 1)
}

func timeNowS(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	if ac := moduleutils.CheckArgs("now_unix", 0, args...); ac != nil {
		return ac
	}

	return object.MakeIntObj(time.Now().Unix())
}

func timeNowMs(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	if ac := moduleutils.CheckArgs("now_ms", 0, args...); ac != nil {
		return ac
	}

	return object.MakeIntObj(time.Now().UnixNano() / 1e6)
}

func timeNowNs(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	if ac := moduleutils.CheckArgs("now_ns", 0, args...); ac != nil {
		return ac
	}

	return object.MakeIntObj(time.Now().UnixNano())
}

func timeNow(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	if ac := moduleutils.CheckArgs("now", 0, args...); ac != nil {
		return ac
	}

	theVM := interpreter.(*vm.VirtualMachine)
	timeMod := theVM.GetLoadedModule("std.time")
	if timeMod == nil {
		return object.NewException("std.time library not loaded")
	}

	timeClass := timeMod.Vars["Time"]
	theVM.MakeInstance(0, timeClass)
	timeInstanceObj := theVM.PopStack()
	if timeInstanceObj.Type() != object.InstanceObj {
		return object.NewException("Failed to create instance of std.time.Time")
	}
	timeInstance := timeInstanceObj.(*vm.VMInstance)
	timeInstance.Fields.SetForce("u", object.MakeIntObj(time.Now().Unix()), true)

	return timeInstance
}

func timeCFormat(interpreter *vm.VirtualMachine, self *vm.VMInstance, env *object.Environment, args ...object.Object) object.Object {
	if ac := moduleutils.CheckArgs("format", 1, args...); ac != nil {
		return ac
	}

	selfUnix, _ := self.Fields.Get("u")
	target, ok := selfUnix.(*object.Integer)
	if !ok {
		return object.NewException("u field is not an int")
	}

	f, ok := args[0].(*object.String)
	if !ok {
		return object.NewException("format expected a string, got %s", args[0].Type().String())
	}

	formatStr := string(f.Value)
	if formatStr == "" {
		formatStr = time.RFC1123
	}

	return object.MakeStringObj(time.Unix(target.Value, 0).Format(formatStr))
}
