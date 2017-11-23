package eval

import (
	"github.com/nitrogen-lang/nitrogen/src/ast"
	"github.com/nitrogen-lang/nitrogen/src/object"
)

func (i *Interpreter) evalMakeInstance(ms *ast.MakeInstance, env *object.Environment) object.Object {
	class := i.Eval(ms.Class, env)
	if class == object.NullConst {
		return object.NewException("Class %s not defined", ms.Class)
	}

	classObj, ok := class.(*object.Class)
	if !ok {
		return object.NewException("%s is not a class", ms.Class)
	}

	cClass := classObj
	classChain := make([]*object.Class, 0, 3)
	classChain = append(classChain, classObj)
	for cClass.Parent != nil {
		classChain = append(classChain, cClass.Parent)
		cClass = cClass.Parent
	}

	iFields := object.NewEnvironment()
	iFields.SetParent(env)

	for c := len(classChain) - 1; c >= 0; c-- {
		for _, def := range classChain[c].Fields {
			i.Eval(def, iFields)
		}
		iFields = object.NewEnclosedEnv(iFields)
	}

	instance := &object.Instance{
		Class:  classObj,
		Fields: iFields.Parent(),
	}

	init := classObj.GetMethod("init")
	if init == nil {
		return instance
	}

	initEnv := object.NewEnclosedEnv(instance.Fields)
	initEnv.CreateConst("this", instance)
	if len(classChain) > 1 {
		initEnv.CreateConst("parent", classChain[1])
	}
	oldInstance := i.currentInstance
	i.currentInstance = instance
	args := i.evalExpressions(ms.Arguments, env)
	i.currentInstance = oldInstance
	if len(args) == 1 && isException(args[0]) {
		return args[0]
	}

	if ret := i.applyFunctionDirect(init, args, initEnv); isException(ret) || isPanic(ret) {
		return ret
	}
	return instance
}
