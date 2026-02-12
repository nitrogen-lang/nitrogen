package vm

import (
	"path/filepath"
	"strings"

	"github.com/nitrogen-lang/nitrogen/src/elemental/object"
	"github.com/nitrogen-lang/nitrogen/src/moduleutils"
)

func pathToName(path string) string {
	path = strings.ReplaceAll(path, "/", ".")
	path = strings.ReplaceAll(path, "\\", ".")
	path = strings.ReplaceAll(path, "..", ".")
	if path[0] == '.' {
		path = path[1:]
	}
	return path
}

func (vm *VirtualMachine) importPackage(path string) {
	name := pathToName(path)
	if mod, loaded := vm.loadedModules[name]; loaded {
		vm.currentFrame.pushStack(mod)
		return
	}

	mod := GetNativeModule(name)
	if mod != nil {
		vm.loadedModules[name] = mod
		vm.currentFrame.pushStack(mod)
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

	switch module.Type() {
	case object.ExceptionObj:
		vm.currentFrame.pushStack(module)
		vm.throw()
		return
	case object.ModuleObj:
		vm.loadedModules[name] = module.(*object.Module)
	}
	vm.currentFrame.pushStack(module)
}

func importScriptFile(vm *VirtualMachine, scriptPath, name string) object.Object {
	code, err := moduleutils.CodeBlockCache.GetBlock(scriptPath, name)
	if err != nil {
		return object.NewException("importing %s failed:\n%s", name, err.Error())
	}

	if filepath.Ext(scriptPath) == ".nib" {
		scriptPath = scriptPath[:len(scriptPath)-1]
	}

	env := object.NewEnclosedEnv(vm.globalEnv)
	env.CreateConst("_FILE", object.MakeStringObj(scriptPath))

	ret := vm.RunFrame(vm.MakeFrame(code, env, name), true)
	if object.ObjectIs(ret, object.ExceptionObj) {
		return ret
	}

	return &object.Module{
		Name: name,
		Vars: env.GetExported(),
	}
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
