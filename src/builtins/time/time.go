package time

import (
	"time"

	"github.com/nitrogen-lang/nitrogen/src/moduleutils"
	"github.com/nitrogen-lang/nitrogen/src/object"
	"github.com/nitrogen-lang/nitrogen/src/vm"
)

func init() {
	vm.RegisterNative("std.time.now", timeNowS)
	vm.RegisterNative("std.time.now_ms", timeNowMs)
	vm.RegisterNative("std.time.now_ns", timeNowNs)
}

func timeNowS(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	if ac := moduleutils.CheckArgs("now", 0, args...); ac != nil {
		return ac
	}

	return object.MakeIntObj(time.Now().Unix())
}

func timeNowMs(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	if ac := moduleutils.CheckArgs("now", 0, args...); ac != nil {
		return ac
	}

	return object.MakeIntObj(time.Now().UnixNano() / 1e6)
}

func timeNowNs(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	if ac := moduleutils.CheckArgs("now", 0, args...); ac != nil {
		return ac
	}

	return object.MakeIntObj(time.Now().UnixNano())
}
