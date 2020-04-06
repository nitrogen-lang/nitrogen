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

func vmFileOpenFile(interpreter *vm.VirtualMachine, self *vm.VMInstance, env *object.Environment, args ...object.Object) object.Object {
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

	res := &fileResource{
		file:   file,
		reader: bufio.NewReader(file),
		mode:   fileMode,
	}

	self.Fields.SetForce("res", res, true)
	return nil
}

func vmFileCloseFile(interpreter *vm.VirtualMachine, self *vm.VMInstance, env *object.Environment, args ...object.Object) object.Object {
	if ac := moduleutils.CheckArgs("closeFile", 0, args...); ac != nil {
		return ac
	}

	res, exists := self.Fields.Get("res")
	if !exists {
		return object.NewException("File object doesn't contain a resource")
	}

	file, ok := res.(*fileResource)
	if !ok {
		return object.NewException("closeFile expected a file resource, got %s", res.Type().String())
	}

	file.file.Close()

	self.Fields.SetForce("res", object.NullConst, true)
	return object.NullConst
}

func vmFileWriteFile(interpreter *vm.VirtualMachine, self *vm.VMInstance, env *object.Environment, args ...object.Object) object.Object {
	if ac := moduleutils.CheckArgs("writeFile", 1, args...); ac != nil {
		return ac
	}

	res, exists := self.Fields.Get("res")
	if !exists {
		return object.NewException("File object doesn't contain a resource")
	}

	file, ok := res.(*fileResource)
	if !ok {
		return object.NewException("writeFile expected a file resource, got %s", res.Type().String())
	}

	str, ok := args[0].(*object.String)
	if !ok {
		return object.NewException("writeFile expected a string, got %s", args[0].Type().String())
	}

	written, err := file.file.WriteString(str.String())
	if err != nil {
		fmt.Fprintln(interpreter.GetStderr(), err)
	}

	return object.MakeIntObj(int64(written))
}

func vmFileReadFullFile(interpreter *vm.VirtualMachine, self *vm.VMInstance, env *object.Environment, args ...object.Object) object.Object {
	if ac := moduleutils.CheckArgs("readFullFile", 0, args...); ac != nil {
		return ac
	}

	res, exists := self.Fields.Get("res")
	if !exists {
		return object.NewException("File object doesn't contain a resource")
	}

	file, ok := res.(*fileResource)
	if !ok {
		return object.NewException("writeFile expected a file resource, got %s", res.Type().String())
	}

	bytes, err := ioutil.ReadAll(file.file)
	if err != nil {
		return object.NewException("Error reading file %s", err.Error())
	}

	return object.MakeStringObj(string(bytes))
}

func vmFileReadLine(interpreter *vm.VirtualMachine, self *vm.VMInstance, env *object.Environment, args ...object.Object) object.Object {
	if ac := moduleutils.CheckArgs("readLine", 0, args...); ac != nil {
		return ac
	}

	res, exists := self.Fields.Get("res")
	if !exists {
		return object.NewException("File object doesn't contain a resource")
	}

	file, ok := res.(*fileResource)
	if !ok {
		return object.NewException("writeFile expected a file resource, got %s", res.Type().String())
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

func vmFileReadChar(interpreter *vm.VirtualMachine, self *vm.VMInstance, env *object.Environment, args ...object.Object) object.Object {
	if ac := moduleutils.CheckArgs("readChar", 0, args...); ac != nil {
		return ac
	}

	res, exists := self.Fields.Get("res")
	if !exists {
		return object.NewException("File object doesn't contain a resource")
	}

	file, ok := res.(*fileResource)
	if !ok {
		return object.NewException("writeFile expected a file resource, got %s", res.Type().String())
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

func vmFileDeleteFile(interpreter *vm.VirtualMachine, self *vm.VMInstance, env *object.Environment, args ...object.Object) object.Object {
	if ac := moduleutils.CheckArgs("deleteFile", 0, args...); ac != nil {
		return ac
	}

	res, exists := self.Fields.Get("res")
	if !exists {
		return object.NewException("File object doesn't contain a resource")
	}

	file, ok := res.(*fileResource)
	if !ok {
		return object.NewException("writeFile expected a file resource, got %s", res.Type().String())
	}

	filepath := file.file.Name()
	file.file.Close()

	if !fileExistsCheck(filepath) {
		return object.NullConst
	}

	if err := os.Remove(filepath); err != nil {
		return object.NewException("Error reading file %s", err.Error())
	}

	self.Fields.SetForce("res", object.NullConst, true)
	return object.NullConst
}

func vmFileRenameFile(interpreter *vm.VirtualMachine, self *vm.VMInstance, env *object.Environment, args ...object.Object) object.Object {
	if ac := moduleutils.CheckArgs("renameFile", 1, args...); ac != nil {
		return ac
	}

	// Get resource
	res, exists := self.Fields.Get("res")
	if !exists {
		return object.NewException("File object doesn't contain a resource")
	}

	file, ok := res.(*fileResource)
	if !ok {
		return object.NewException("writeFile expected a file resource, got %s", res.Type().String())
	}

	// Get old and new names
	oldPath := file.file.Name()

	newPath, ok := args[0].(*object.String)
	if !ok {
		return object.NewException("renameFile expected a string, got %s", args[0].Type().String())
	}

	// Close file and rename
	file.file.Close()

	if err := os.Rename(oldPath, newPath.String()); err != nil {
		return object.NewError("Error renaming file %s", err.Error())
	}

	self.Fields.SetForce("res", object.NullConst, true)
	return object.NullConst
}
