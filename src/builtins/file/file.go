package file

import (
	"bufio"
	"io/ioutil"
	"os"

	"github.com/nitrogen-lang/nitrogen/src/elemental/object"
	"github.com/nitrogen-lang/nitrogen/src/elemental/vm"
	"github.com/nitrogen-lang/nitrogen/src/moduleutils"
)

const (
	fileResourceID = "std.file"
)

func init() {
	vm.RegisterNative("std.file.readFile", readFullFile)
	vm.RegisterNative("std.file.remove", deleteFile)
	vm.RegisterNative("std.file.exists", fileExists)
	vm.RegisterNative("std.file.rename", renameFile)
	vm.RegisterNative("std.file.dirlist", directoryList)
	vm.RegisterNative("std.file.isdir", isDirectory)

	vm.RegisterNativeMethod("std.file.File.init", vmFileOpenFile, 1)
	vm.RegisterNativeMethod("std.file.File.close", vmFileCloseFile, 0)
	vm.RegisterNativeMethod("std.file.File.write", vmFileWriteFile, 1)
	vm.RegisterNativeMethod("std.file.File.readAll", vmFileReadFullFile, 0)
	vm.RegisterNativeMethod("std.file.File.readLine", vmFileReadLine, 0)
	vm.RegisterNativeMethod("std.file.File.readChar", vmFileReadChar, 0)
	vm.RegisterNativeMethod("std.file.File.remove", vmFileDeleteFile, 0)
	vm.RegisterNativeMethod("std.file.File.rename", vmFileRenameFile, 1)
}

type fileResource struct {
	file   *os.File
	mode   int
	reader *bufio.Reader
}

func (f *fileResource) Inspect() string         { return "File resource" }
func (f *fileResource) Type() object.ObjectType { return object.ResourceObj }
func (f *fileResource) Dup() object.Object      { return object.NullConst } // Duplicating a file resource isn't allowed
func (f *fileResource) ResourceID() string      { return fileResourceID }

var modes = map[string]int{
	"r":  os.O_RDONLY,
	"r+": os.O_RDWR,
	"w":  os.O_WRONLY | os.O_TRUNC | os.O_CREATE,
	"w+": os.O_RDWR | os.O_TRUNC | os.O_CREATE,
	"a":  os.O_APPEND | os.O_WRONLY | os.O_CREATE,
	"a+": os.O_APPEND | os.O_RDWR | os.O_CREATE,
}

func readFullFile(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	if ac := moduleutils.CheckArgs("readFullFile", 1, args...); ac != nil {
		return ac
	}

	filepath, ok := args[0].(*object.String)
	if !ok {
		return object.NewException("readFullFile expected a string, got %s", args[0].Type().String())
	}

	file, err := ioutil.ReadFile(filepath.String())
	if err != nil {
		return object.NewException("Error reading file %s", err.Error())
	}

	return object.MakeStringObj(string(file))
}

func deleteFile(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	if ac := moduleutils.CheckArgs("deleteFile", 1, args...); ac != nil {
		return ac
	}

	filepath, ok := args[0].(*object.String)
	if !ok {
		return object.NewException("deleteFile expected a string, got %s", args[0].Type().String())
	}

	if !fileExistsCheck(filepath.String()) {
		return object.NullConst
	}

	if err := os.Remove(filepath.String()); err != nil {
		return object.NewException("Error reading file %s", err.Error())
	}

	return object.NullConst
}

func fileExists(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	if ac := moduleutils.CheckArgs("fileExists", 1, args...); ac != nil {
		return ac
	}

	filepath, ok := args[0].(*object.String)
	if !ok {
		return object.NewException("fileExists expected a string, got %s", args[0].Type().String())
	}

	return &object.Boolean{Value: fileExistsCheck(filepath.String())}
}

func fileExistsCheck(file string) bool {
	_, err := os.Stat(file)
	return !os.IsNotExist(err)
}

func renameFile(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	if ac := moduleutils.CheckArgs("renameFile", 2, args...); ac != nil {
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

	if err := os.Rename(oldPath.String(), newPath.String()); err != nil {
		return object.NewError("Error renaming file %s", err.Error())
	}

	return object.NullConst
}

func directoryList(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	if ac := moduleutils.CheckArgs("dirList", 1, args...); ac != nil {
		return ac
	}

	filepath, ok := args[0].(*object.String)
	if !ok {
		return object.NewException("dirList expected a string, got %s", args[0].Type().String())
	}

	file, err := os.Open(filepath.String())
	if err != nil {
		return object.NewException("Error opening directory %s", err.Error())
	}
	defer file.Close()

	dirlist, err := file.Readdirnames(0)
	if err != nil {
		return object.NewException("Error reading directory list %s %s", string(filepath.Value), err.Error())
	}
	return object.MakeStringArray(dirlist)
}

func isDirectory(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	if ac := moduleutils.CheckArgs("dirList", 1, args...); ac != nil {
		return ac
	}

	filepath, ok := args[0].(*object.String)
	if !ok {
		return object.NewException("dirList expected a string, got %s", args[0].Type().String())
	}

	file, err := os.Stat(filepath.String())
	if err != nil {
		return object.NewException("Error opening directory %s", err.Error())
	}
	return object.NativeBoolToBooleanObj(file.IsDir())
}
