package file

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/nitrogen-lang/nitrogen/src/moduleutils"
	"github.com/nitrogen-lang/nitrogen/src/object"
	"github.com/nitrogen-lang/nitrogen/src/vm"
)

const (
	moduleName     = "std/file"
	fileResourceID = "std.file"
)

func init() {
	vm.RegisterModule(moduleName, &object.Module{
		Name: moduleName,
		Methods: map[string]object.BuiltinFunction{
			"open":     openFile,
			"close":    closeFile,
			"write":    writeFile,
			"readAll":  readFullFile,
			"readLine": readLine,
			"readChar": readChar,
			"remove":   deleteFile,
			"exists":   fileExists,
			"rename":   renameFile,
			"dirlist":  directoryList,
			"isdir":    isDirectory,
		},
		Vars: map[string]object.Object{
			"name":           object.MakeStringObj(moduleName),
			"fileResourceID": object.MakeStringObj(fileResourceID),
		},
	})
}

type fileResource struct {
	file   *os.File
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

func openFile(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	if ac := moduleutils.CheckArgs("openFile", 2, args...); ac != nil {
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

	fileMode, ok := modes[mode.String()]
	if !ok {
		return object.NewException("Invalid file mode %s", mode.String())
	}

	file, err := os.OpenFile(filepath.String(), fileMode, 0644)
	if err != nil {
		return object.NewException("Error opening file %s", err.Error())
	}

	return &fileResource{
		file:   file,
		reader: bufio.NewReader(file),
	}
}

func closeFile(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	if ac := moduleutils.CheckArgs("closeFile", 1, args...); ac != nil {
		return ac
	}

	file, ok := args[0].(*fileResource)
	if !ok {
		return object.NewException("closeFile expected a file resource, got %s", args[0].Type().String())
	}

	file.file.Close()

	return object.NullConst
}

func writeFile(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	if ac := moduleutils.CheckArgs("writeFile", 2, args...); ac != nil {
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

	written, err := file.file.WriteString(str.String())
	if err != nil {
		fmt.Fprintln(interpreter.GetStderr(), err)
	}

	return object.MakeIntObj(int64(written))
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

func readLine(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	if ac := moduleutils.CheckArgs("readLine", 1, args...); ac != nil {
		return ac
	}

	file, ok := args[0].(*fileResource)
	if !ok {
		return object.NewException("readLine expected a file resource, got %s", args[0].Type().String())
	}

	line, err := file.reader.ReadString('\n')
	if err != nil {
		if err == io.EOF {
			if line != "" {
				object.MakeStringObj(line)
			}
			return object.NullConst
		}
		return object.NewException(err.Error())
	}

	return object.MakeStringObj(line[:len(line)-1])
}

func readChar(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	if ac := moduleutils.CheckArgs("readChar", 1, args...); ac != nil {
		return ac
	}

	file, ok := args[0].(*fileResource)
	if !ok {
		return object.NewException("readChar expected a file resource, got %s", args[0].Type().String())
	}

	r, _, err := file.reader.ReadRune()
	if err != nil {
		if err == io.EOF {
			if r != 0 {
				return object.MakeStringObjRunes([]rune{r})
			}
			return object.NullConst
		}
		return object.NewException(err.Error())
	}

	return object.MakeStringObjRunes([]rune{r})
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
		return object.NewException("Error reading directory list %s %s", filepath.Value, err.Error())
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
