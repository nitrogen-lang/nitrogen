package vm

import (
	"path/filepath"

	"github.com/nitrogen-lang/nitrogen/src/moduleutils"
	"github.com/nitrogen-lang/nitrogen/src/object"
)

var included = make(map[string]object.Object)

func (vm *VirtualMachine) importPackage(name, path string) {
	mod := GetModule(name)
	if mod != nil {
		_, err := vm.currentFrame.env.Set(name, mod)
		if err != nil {
			vm.currentFrame.env.SetForce(name, mod, false)
		}
		return
	}

	searchPaths, ok := vm.currentFrame.env.Get("_SEARCH_PATHS")
	if !ok {
		vm.currentFrame.pushStack(object.NewException("_SEARCH_PATHS variable not found, required for module lookup"))
		vm.throw()
		return
	}
	if !object.ObjectIs(searchPaths, object.ArrayObj) {
		vm.currentFrame.pushStack(object.NewException("_SEARCH_PATHS must be an array, required for module lookup"))
		vm.throw()
		return
	}

	includedFile := findModule(path, vm.GetCurrentScriptPath(), object.ArrayToStringSlice(searchPaths.(*object.Array)))
	if includedFile == "" {
		vm.currentFrame.pushStack(object.NewException("import failed, module not found %s", path))
		vm.throw()
		return
	}

	var module object.Object
	if filepath.Ext(includedFile) == ".so" {
		module = importSharedModule(vm, includedFile, name)
	} else {
		module = importScriptFile(vm, includedFile, name)
	}

	if object.ObjectIs(module, object.ExceptionObj) {
		vm.currentFrame.pushStack(module)
		vm.throw()
		return
	}

	_, err := vm.currentFrame.env.Set(name, module)
	if err != nil {
		vm.currentFrame.env.SetForce(name, module, false)
	}
}

func importScriptFile(vm *VirtualMachine, scriptPath, name string) object.Object {
	res, imported := included[scriptPath]
	if imported {
		return res
	}

	code, err := moduleutils.CodeBlockCache.GetBlock(scriptPath, name)
	if err != nil {
		return object.NewException("importing %s failed %s", name, err.Error())
	}

	env := object.NewEnclosedEnv(vm.currentFrame.env)
	env.CreateConst("_FILE", object.MakeStringObj(code.Filename))

	res = vm.RunFrame(vm.MakeFrame(code, env), true)
	included[scriptPath] = res
	return res
}

var extensions = []string{"", ".nib", ".ni", ".so"}

func findModule(name, scriptPath string, searchPaths []string) string {
	if name[0] == '/' { // Absolute path
		return testModulePath(name)
	} else if name[0] == '.' { // Relative path to script file
		fullpath := filepath.Clean(filepath.Join(filepath.Dir(scriptPath), name))
		return testModulePath(fullpath)
	}

	// Search for module
	for _, path := range searchPaths {
		mp := testModulePath(filepath.Join(path, name))
		if mp != "" {
			return mp
		}
	}
	return ""
}

func testModulePath(path string) string {
	for _, ext := range extensions {
		fullname := path + ext
		if moduleutils.IsDir(fullname) {
			mp := testModulePath(filepath.Join(path, "mod.nib"))
			if mp != "" {
				return mp
			}

			mp = testModulePath(filepath.Join(path, "mod.ni"))
			if mp != "" {
				return mp
			}
		}
		if moduleutils.FileExists(fullname) {
			return fullname
		}
	}
	return ""
}
