package vm

import (
	"path/filepath"

	"github.com/nitrogen-lang/nitrogen/src/ast"
	"github.com/nitrogen-lang/nitrogen/src/moduleutils"
	"github.com/nitrogen-lang/nitrogen/src/object"
)

var included = make(map[string]*ast.Program)

func (vm *VirtualMachine) importPackage(name, path string) {
	searchPaths, ok := vm.currentFrame.Env.Get("_SEARCH_PATHS")
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
		module = importSharedModule(vm, includedFile)
	} else {
		module = importScriptFile(vm, includedFile)
	}

	if object.ObjectIs(module, object.ExceptionObj) {
		vm.currentFrame.pushStack(module)
		vm.throw()
		return
	}

	vm.currentFrame.Env.SetForce(name, module, true)
}

func importScriptFile(vm *VirtualMachine, scriptPath string) object.Object {
	code, err := moduleutils.CodeBlockCache.GetBlock(scriptPath)
	if err != nil {
		return object.NewException("importing %s failed %s", scriptPath, err.Error())
	}

	env := object.NewEnclosedEnv(vm.currentFrame.Env)
	env.CreateConst("_FILE", object.MakeStringObj(code.Filename))
	return vm.RunFrame(vm.MakeFrame(code, env), true)
}

func findModule(name, scriptPath string, searchPaths []string) string {
	if name[0] == '/' { // Absolute path
		if moduleutils.FileExists(name) {
			return name
		}
	} else if name[0] == '.' { // Relative path to script file
		fullpath := filepath.Clean(filepath.Join(filepath.Dir(scriptPath), name))
		if moduleutils.FileExists(fullpath) {
			return fullpath
		}
	} else { // Search for module
		// TODO: Use the _SEARCH_PATHS variable when loading instead
		for _, path := range searchPaths {
			fullpath := filepath.Join(path, name)
			if moduleutils.FileExists(fullpath) {
				return fullpath
			}
		}
	}
	return ""
}
