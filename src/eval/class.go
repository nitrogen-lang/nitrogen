package eval

import (
	"github.com/nitrogen-lang/nitrogen/src/ast"
	"github.com/nitrogen-lang/nitrogen/src/object"
)

func (i *Interpreter) evalMakeInstance(ms *ast.MakeInstance, env *object.Environment) object.Object {
	class, exists := env.Get(ms.Class)
	if !exists {
		return object.NewException("Class %s not defined", ms.Class)
	}

	classObj, ok := class.(*object.Class)
	if !ok {
		return object.NewException("%s is not a class", ms.Class)
	}

	instance := &object.Instance{
		Class:  classObj,
		Fields: object.NewEnvironment(),
	}
	instance.Fields.SetParent(env)
	instance.Fields.CreateConst("this", instance)

	for _, def := range classObj.Fields {
		i.Eval(def, instance.Fields)
	}

	init := classObj.GetMethod("init")
	if init == nil {
		return instance
	}

	args := i.evalExpressions(ms.Arguments, env)
	if len(args) == 1 && isException(args[0]) {
		return args[0]
	}

	if ret := i.applyFunctionDirect(init, args, instance.Fields); isException(ret) || isPanic(ret) {
		return ret
	}
	return instance
}
