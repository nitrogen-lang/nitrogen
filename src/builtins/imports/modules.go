// +build linux,cgo darwin,cgo

package imports

import (
	"path/filepath"
	"plugin"

	"github.com/nitrogen-lang/nitrogen/src/config"
	"github.com/nitrogen-lang/nitrogen/src/moduleutils"
	"github.com/nitrogen-lang/nitrogen/src/object"
	"github.com/nitrogen-lang/nitrogen/src/vm"
)

func init() {
	vm.RegisterBuiltin("module", importModuleVM)
	vm.RegisterBuiltin("modulesSupported", moduleSupport)
}

func moduleSupport(i object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	return object.NativeBoolToBooleanObj(true)
}

func importModuleVM(i object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	if ac := moduleutils.CheckMinArgs("module", 1, args...); ac != nil {
		return ac
	}

	filepathArg, ok := args[0].(*object.String)
	if !ok {
		return object.NewException("module expected a string, got %s", args[0].Type().String())
	}

	if filepathArg.Value == "" {
		return object.NewException("module name expected")
	}

	required := false
	if len(args) > 1 {
		requiredArg, ok := args[1].(*object.Boolean)
		if !ok {
			return object.NewException("module expected a boolean for second argument, got %s", args[1].Type().String())
		}
		required = requiredArg.Value
	}

	// Return already registered, named module
	if module := vm.GetModule(filepathArg.Value); module != nil {
		return module
	}

	modulepath := ""
	if filepathArg.Value[0] == '/' { // Absolute path
		if fileExists(filepathArg.Value) {
			modulepath = filepathArg.Value
		}
	} else if filepathArg.Value[0] == '.' { // Relative path to script file
		fullpath := filepath.Clean(filepath.Join(filepath.Dir(i.GetCurrentScriptPath()), filepathArg.Value))
		if fileExists(fullpath) {
			modulepath = fullpath
		}
	} else { // Search for module
		// TODO: Use the _SEARCH_PATHS variable when loading instead
		for _, path := range config.ModulePaths {
			fullpath := filepath.Join(path, filepathArg.Value)
			if fileExists(fullpath) {
				modulepath = fullpath
			}
		}
	}

	if modulepath == "" {
		if required {
			return object.NewException("Module %s not found", filepathArg.Value)
		}
		return object.NewError("Module %s not found", filepathArg.Value)
	}

	p, err := plugin.Open(modulepath)
	if err != nil {
		if required {
			return object.NewException("%s", err)
		}
		return object.NewError("%s", err)
	}

	// Check module name
	moduleNameSym, err := p.Lookup("ModuleName")
	if err != nil {
		// The module didn't declare a name
		return object.NewException("Invalid module %s", filepathArg.Value)
	}

	if module := vm.GetModule(*(moduleNameSym.(*string))); module != nil {
		return module
	}
	return object.NullConst
}
