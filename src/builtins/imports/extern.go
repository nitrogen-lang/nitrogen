package imports

import (
	"path/filepath"
	"plugin"

	"github.com/nitrogen-lang/nitrogen/src/ast"
	"github.com/nitrogen-lang/nitrogen/src/eval"
	"github.com/nitrogen-lang/nitrogen/src/lexer"
	"github.com/nitrogen-lang/nitrogen/src/moduleutils"
	"github.com/nitrogen-lang/nitrogen/src/object"
	"github.com/nitrogen-lang/nitrogen/src/parser"
)

var included map[string]*ast.Program

func init() {
	eval.RegisterBuiltin("include", includeScript)
	eval.RegisterBuiltin("require", requireScript)
	eval.RegisterBuiltin("evalScript", evalScript)
	eval.RegisterBuiltin("module", importModule)

	included = make(map[string]*ast.Program)
}

func includeScript(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	return commonInclude(false, true, interpreter, env, args...)
}

func requireScript(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	return commonInclude(true, true, interpreter, env, args...)
}

func importModule(i object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	if ac := moduleutils.CheckMinArgs("module", 1, args...); ac != nil {
		return ac
	}

	filepathArg, ok := args[0].(*object.String)
	if !ok {
		return object.NewException("module expected a string, got %s", args[0].Type().String())
	}

	required := false
	if len(args) > 1 {
		requiredArg, ok := args[1].(*object.Boolean)
		if !ok {
			return object.NewException("module expected a boolean for second argument, got %s", args[1].Type().String())
		}
		required = requiredArg.Value
	}

	includedPath := filepath.Clean(filepath.Join(filepath.Dir(i.GetCurrentScriptPath()), filepathArg.Value))
	//return &object.String{Value: includedPath}
	_, err := plugin.Open(includedPath)
	if err != nil {
		if required {
			return &object.Exception{Message: err.Error()}
		}
		return &object.Error{Message: err.Error()}
	}
	return object.NullConst
}

func evalScript(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	cleanEnv := object.NewEnvironment()

	envvar, _ := env.Get("_ARGV")
	cleanEnv.CreateConst("_ARGV", envvar.Dup())

	envvar, _ = env.Get("_ENV")
	cleanEnv.CreateConst("_ENV", envvar.Dup())

	return commonInclude(false, false, interpreter, cleanEnv, args...)
}

func commonInclude(require bool, save bool, i object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	funcName := "include"
	if require {
		funcName = "require"
	}

	if ac := moduleutils.CheckMinArgs(funcName, 1, args...); ac != nil {
		return ac
	}

	filepathArg, ok := args[0].(*object.String)
	if !ok {
		return object.NewException("%s expected a string, got %s", funcName, args[0].Type().String())
	}

	once := false
	if len(args) > 1 {
		includeOnce, ok := args[1].(*object.Boolean)
		if !ok {
			return object.NewException("%s expected a boolean for second argument, got %s", funcName, args[1].Type().String())
		}
		once = includeOnce.Value
	}

	includedFile := filepath.Clean(filepath.Join(filepath.Dir(i.GetCurrentScriptPath()), filepathArg.Value))

	program, exists := included[includedFile]
	if exists {
		if once || program == nil {
			return object.NullConst
		}
		return i.Eval(program, object.NewEnclosedEnv(env))
	}

	l, err := lexer.NewFile(includedFile)
	if err != nil {
		if require {
			return object.NewException("including %s failed %s", includedFile, err.Error())
		}
		return object.NewError("including %s failed %s", includedFile, err.Error())
	}

	p := parser.New(l)
	program = p.ParseProgram()
	if len(p.Errors()) != 0 {
		if require {
			return object.NewException("including %s failed %s", includedFile, p.Errors()[0])
		}
		return object.NewError("including %s failed %s", includedFile, p.Errors()[0])
	}

	if save {
		if once {
			// Create the key, but don't save the parsed script since we don't need it anymore.
			included[includedFile] = nil
		} else {
			included[includedFile] = program
		}
	}
	return i.Eval(program, object.NewEnclosedEnv(env))
}
