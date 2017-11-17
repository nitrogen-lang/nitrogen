package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/nitrogen-lang/nitrogen/src/builtins"
	"github.com/nitrogen-lang/nitrogen/src/eval"
	"github.com/nitrogen-lang/nitrogen/src/object"
)

func init() {
	eval.RegisterBuiltin("openFile", openFile)
	eval.RegisterBuiltin("closeFile", closeFile)
	eval.RegisterBuiltin("writeFile", writeFile)
	eval.RegisterBuiltin("readFullFile", readFullFile)
	eval.RegisterBuiltin("deleteFile", deleteFile)
	eval.RegisterBuiltin("fileExists", fileExists)
	eval.RegisterBuiltin("renameFile", renameFile)
}

func main() {}

type fileResource struct {
	file *os.File
}

func (f *fileResource) Inspect() string         { return "File resource" }
func (f *fileResource) Type() object.ObjectType { return object.RESOURCE_OBJ }

var modes = map[string]int{
	"r":  os.O_RDONLY,
	"r+": os.O_RDWR,
	"w":  os.O_WRONLY | os.O_TRUNC | os.O_CREATE,
	"w+": os.O_RDWR | os.O_TRUNC | os.O_CREATE,
	"a":  os.O_APPEND | os.O_WRONLY | os.O_CREATE,
	"a+": os.O_APPEND | os.O_RDWR | os.O_CREATE,
}

func openFile(env *object.Environment, args ...object.Object) object.Object {
	if ac := builtins.CheckArgs("openFile", 2, args...); ac != nil {
		return ac
	}

	filepath, ok := args[0].(*object.String)
	if !ok {
		return object.NewException("openFile expected a string, got %s", args[0].Type().String())
	}

	mode, ok := args[1].(*object.String)
	if !ok {
		return object.NewException("openFile expected a string, got %s", args[0].Type().String())
	}

	fileMode, ok := modes[mode.Value]
	if !ok {
		return object.NewException("Invalid file mode %s", mode.Value)
	}

	file, err := os.OpenFile(filepath.Value, fileMode, 0644)
	if err != nil {
		return object.NewException("Error opening file %s", err.Error())
	}

	return &fileResource{file}
}

func closeFile(env *object.Environment, args ...object.Object) object.Object {
	if ac := builtins.CheckArgs("closeFile", 1, args...); ac != nil {
		return ac
	}

	file, ok := args[0].(*fileResource)
	if !ok {
		return object.NewException("closeFile expected a file resource, got %s", args[0].Type().String())
	}

	file.file.Close()

	return object.NULL
}

func writeFile(env *object.Environment, args ...object.Object) object.Object {
	if ac := builtins.CheckArgs("writeFile", 2, args...); ac != nil {
		return ac
	}

	file, ok := args[0].(*fileResource)
	if !ok {
		return object.NewException("writeFile expected a file resource, got %s", args[0].Type().String())
	}

	str, ok := args[1].(*object.String)
	if !ok {
		return object.NewException("writeFile expected a string, got %s", args[1].Type().String())
	}

	written, err := file.file.WriteString(str.Value)
	if err != nil {
		fmt.Println(err)
	}

	return &object.Integer{Value: int64(written)}
}

func readFullFile(env *object.Environment, args ...object.Object) object.Object {
	if ac := builtins.CheckArgs("readFullFile", 1, args...); ac != nil {
		return ac
	}

	filepath, ok := args[0].(*object.String)
	if !ok {
		return object.NewException("readFullFile expected a string, got %s", args[0].Type().String())
	}

	file, err := ioutil.ReadFile(filepath.Value)
	if err != nil {
		return object.NewException("Error reading file %s", err.Error())
	}

	return &object.String{Value: string(file)}
}

func deleteFile(env *object.Environment, args ...object.Object) object.Object {
	if ac := builtins.CheckArgs("deleteFile", 1, args...); ac != nil {
		return ac
	}

	filepath, ok := args[0].(*object.String)
	if !ok {
		return object.NewException("deleteFile expected a string, got %s", args[0].Type().String())
	}

	if !fileExistsCheck(filepath.Value) {
		return object.NULL
	}

	if err := os.Remove(filepath.Value); err != nil {
		return object.NewException("Error reading file %s", err.Error())
	}

	return object.NULL
}

func fileExists(env *object.Environment, args ...object.Object) object.Object {
	if ac := builtins.CheckArgs("fileExists", 1, args...); ac != nil {
		return ac
	}

	filepath, ok := args[0].(*object.String)
	if !ok {
		return object.NewException("fileExists expected a string, got %s", args[0].Type().String())
	}

	return &object.Boolean{Value: fileExistsCheck(filepath.Value)}
}

func fileExistsCheck(file string) bool {
	_, err := os.Stat(file)
	return !os.IsNotExist(err)
}

func renameFile(env *object.Environment, args ...object.Object) object.Object {
	if ac := builtins.CheckArgs("renameFile", 2, args...); ac != nil {
		return ac
	}

	oldPath, ok := args[0].(*object.String)
	if !ok {
		return object.NewException("renameFile expected a string, got %s", args[0].Type().String())
	}

	newPath, ok := args[1].(*object.String)
	if !ok {
		return object.NewException("renameFile expected a string, got %s", args[0].Type().String())
	}

	if err := os.Rename(oldPath.Value, newPath.Value); err != nil {
		return object.NewError("Error renaming file %s", err.Error())
	}

	return object.NULL
}
