package imports

import (
	"path/filepath"

	"github.com/nitrogen-lang/nitrogen/src/ast"
	"github.com/nitrogen-lang/nitrogen/src/config"
	"github.com/nitrogen-lang/nitrogen/src/moduleutils"
	"github.com/nitrogen-lang/nitrogen/src/object"
	"github.com/nitrogen-lang/nitrogen/src/vm"
)

var included map[string]*ast.Program

func init() {
	vm.RegisterBuiltin("import", includeScript) // import(filename[, throw])
	vm.RegisterBuiltin("modulesSupported", moduleSupport)
	// vm.RegisterBuiltin("evalScript", evalScript)

	included = make(map[string]*ast.Program)
}

// func evalScript(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
// 	cleanEnv := object.NewEnvironment()

// 	envvar, _ := env.Get("_ARGV")
// 	cleanEnv.CreateConst("_ARGV", envvar.Dup())

// 	envvar, _ = env.Get("_ENV")
// 	cleanEnv.CreateConst("_ENV", envvar.Dup())

// 	return commonInclude(false, false, interpreter, cleanEnv, args...)
// }

func includeScript(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	funcName := "import"

	if ac := moduleutils.CheckMinArgs(funcName, 1, args...); ac != nil {
		return ac
	}

	filepathArg, ok := args[0].(*object.String)
	if !ok {
		return object.NewException("%s expected a string, got %s", funcName, args[0].Type().String())
	}

	required := true
	if len(args) == 2 {
		reqArg, ok := args[1].(*object.Boolean)
		if !ok {
			return object.NewException("import expected a boolean as 2 argument, got %s", args[1].Type().String())
		}
		required = reqArg.Value
	}

	includedFile := findModule(filepathArg.Value, interpreter.GetCurrentScriptPath(), config.ModulePaths)
	if includedFile == "" {
		if required {
			return object.NewException("import failed, module not found %s", filepathArg.Value)
		}
		return object.NewError("import failed, module not found %s", filepathArg.Value)
	}

	if filepath.Ext(includedFile) == ".so" {
		return importSharedModule(includedFile, required, interpreter, env)
	}
	return importScriptFile(includedFile, required, interpreter, env)
}

func importScriptFile(scriptPath string, required bool, interpreter object.Interpreter, env *object.Environment) object.Object {
	code, err := moduleutils.CodeBlockCache.GetBlock(scriptPath)
	if err != nil {
		if required {
			return object.NewException("importing %s failed %s", scriptPath, err.Error())
		}
		return object.NewError("importing %s failed %s", scriptPath, err.Error())
	}

	i := interpreter.(*vm.VirtualMachine)
	env = object.NewEnclosedEnv(env)
	env.CreateConst("_FILE", object.MakeStringObj(code.Filename))
	return i.RunFrame(i.MakeFrame(code, env), true)
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
