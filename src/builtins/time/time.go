package time

import (
	"time"

	"github.com/nitrogen-lang/nitrogen/src/moduleutils"
	"github.com/nitrogen-lang/nitrogen/src/object"
)

var moduleName = "std/time"

func Init() object.Object {
	return &object.Module{
		Name: moduleName,
		Methods: map[string]object.BuiltinFunction{
			"now":    timeNowS,
			"now_ms": timeNowMs,
			"now_ns": timeNowNs,
		},
	}
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
