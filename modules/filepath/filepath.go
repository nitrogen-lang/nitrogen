package main

import (
	"os"
	"path/filepath"

	"github.com/nitrogen-lang/nitrogen/src/moduleutils"
	"github.com/nitrogen-lang/nitrogen/src/object"
	"github.com/nitrogen-lang/nitrogen/src/vm"
)

func init() {
	filepathModule := &object.Module{
		Name: ModuleName,
		Methods: map[string]object.BuiltinFunction{
			"dir":      filepathDirectory,
			"basename": filepathBasename,
			"ext":      filepathExt,
			"join":     filepathJoin,
			"abs":      filepathAbs,
			"cwd":      filepathCwd,
		},
		Vars: map[string]object.Object{
			"name": object.MakeStringObj(ModuleName),
		},
	}

	vm.RegisterModule(ModuleName, filepathModule)
}

func main() {}

var ModuleName = "filepath"

func filepathDirectory(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	if ac := moduleutils.CheckArgs("dir", 1, args...); ac != nil {
		return ac
	}

	path, ok := args[0].(*object.String)
	if !ok {
		return object.NewException("dir expected a string, got %s", args[0].Type().String())
	}

	return object.MakeStringObj(filepath.Dir(path.String()))
}

func filepathBasename(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	if ac := moduleutils.CheckArgs("basename", 1, args...); ac != nil {
		return ac
	}

	path, ok := args[0].(*object.String)
	if !ok {
		return object.NewException("basename expected a string, got %s", args[0].Type().String())
	}

	return object.MakeStringObj(filepath.Base(path.String()))
}

func filepathExt(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	if ac := moduleutils.CheckArgs("ext", 1, args...); ac != nil {
		return ac
	}

	path, ok := args[0].(*object.String)
	if !ok {
		return object.NewException("ext expected a string, got %s", args[0].Type().String())
	}

	return object.MakeStringObj(filepath.Ext(path.String()))
}

func filepathAbs(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	if ac := moduleutils.CheckArgs("abs", 1, args...); ac != nil {
		return ac
	}

	path, ok := args[0].(*object.String)
	if !ok {
		return object.NewException("abs expected a string, got %s", args[0].Type().String())
	}

	abs, _ := filepath.Abs(path.String())
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

	return object.MakeStringObj(filepath.Join(pathParts...))
}
