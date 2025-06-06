package filepath

import (
	"os"
	stdfp "path/filepath"

	"github.com/nitrogen-lang/nitrogen/src/moduleutils"
	"github.com/nitrogen-lang/nitrogen/src/object"
	"github.com/nitrogen-lang/nitrogen/src/vm"
)

func init() {
	vm.RegisterNative("std.filepath.dir", filepathDirectory)
	vm.RegisterNative("std.filepath.basename", filepathBasename)
	vm.RegisterNative("std.filepath.ext", filepathExt)
	vm.RegisterNative("std.filepath.join", filepathJoin)
	vm.RegisterNative("std.filepath.abs", filepathAbs)
	vm.RegisterNative("std.filepath.cwd", filepathCwd)
}

func filepathDirectory(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	if ac := moduleutils.CheckArgs("dir", 1, args...); ac != nil {
		return ac
	}

	path, ok := args[0].(*object.String)
	if !ok {
		return object.NewException("dir expected a string, got %s", args[0].Type().String())
	}

	return object.MakeStringObj(stdfp.Dir(path.String()))
}

func filepathBasename(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	if ac := moduleutils.CheckArgs("basename", 1, args...); ac != nil {
		return ac
	}

	path, ok := args[0].(*object.String)
	if !ok {
		return object.NewException("basename expected a string, got %s", args[0].Type().String())
	}

	return object.MakeStringObj(stdfp.Base(path.String()))
}

func filepathExt(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	if ac := moduleutils.CheckArgs("ext", 1, args...); ac != nil {
		return ac
	}

	path, ok := args[0].(*object.String)
	if !ok {
		return object.NewException("ext expected a string, got %s", args[0].Type().String())
	}

	return object.MakeStringObj(stdfp.Ext(path.String()))
}

func filepathAbs(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	if ac := moduleutils.CheckArgs("abs", 1, args...); ac != nil {
		return ac
	}

	path, ok := args[0].(*object.String)
	if !ok {
		return object.NewException("abs expected a string, got %s", args[0].Type().String())
	}

	abs, _ := stdfp.Abs(path.String())
	return object.MakeStringObj(abs)
}

func filepathCwd(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	cwd, _ := os.Getwd()
	return object.MakeStringObj(cwd)
}

func filepathJoin(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	if ac := moduleutils.CheckMinArgs("join", 1, args...); ac != nil {
		return ac
	}

	pathParts := make([]string, len(args))
	for i, v := range args {
		path, ok := v.(*object.String)
		if !ok {
			return object.NewException("join expected a string, got %s", v.Type().String())
		}
		pathParts[i] = path.String()
	}

	return object.MakeStringObj(stdfp.Join(pathParts...))
}
