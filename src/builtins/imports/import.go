package imports

import (
	"github.com/nitrogen-lang/nitrogen/src/elemental/vm"
)

func init() {
	vm.RegisterNative("std.preamble.main.modulesSupported", moduleSupport)
}

// func evalScript(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
// 	cleanEnv := object.NewEnvironment()

// 	envvar, _ := env.Get("_ARGV")
// 	cleanEnv.CreateConst("_ARGV", envvar.Dup())

// 	envvar, _ = env.Get("_ENV")
// 	cleanEnv.CreateConst("_ENV", envvar.Dup())

// 	return commonInclude(false, false, interpreter, cleanEnv, args...)
// }
