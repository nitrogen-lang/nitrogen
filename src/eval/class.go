package eval

import (
	"github.com/nitrogen-lang/nitrogen/src/ast"
	"github.com/nitrogen-lang/nitrogen/src/object"
)

func (i *Interpreter) evalLookupAttribute(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.InstanceObj && index.Type() == object.StringObj:
		return i.evalLookupAttributeInstance(left, index)
	case left.Type() == object.ClassObj && index.Type() == object.StringObj:
		return i.evalLookupAttributeClass(left, index)
	}
	return object.NewException("Index operator not allowed: %s", left.Type())
}

func (i *Interpreter) evalLookupAttributeInstance(instance, index object.Object) object.Object {
	instanceObj := instance.(*object.Instance)
	key := index.(*object.String)

	method := instanceObj.GetMethod(key.Value)
	if method != nil {
		switch m := method.(type) {
		case *object.Function:
			fn := &object.Function{
				Name:       m.Name,
				Parameters: m.Parameters,
				Body:       m.Body,
				Env:        object.NewEnclosedEnv(instanceObj.Fields),
				Instance:   instanceObj,
			}
			fn.Env.CreateConst("this", instanceObj)
			if instanceObj.Class.Parent != nil {
				fn.Env.CreateConst("parent", instanceObj.Class.Parent)
			}
			return fn
		case *object.BuiltinMethod:
			return &object.BuiltinMethod{
				Fn:       m.Fn,
				Instance: instanceObj,
			}
		}
	}

	val, ok := instanceObj.Fields.Get(key.Value)
	if ok {
		return val
	}
	return object.NullConst
}

func (i *Interpreter) evalLookupAttributeClass(class, index object.Object) object.Object {
	classObj := class.(*object.Class)
	if !object.InstanceOf(classObj.Name, i.currentInstance) {
		return object.NullConst
	}
	key := index.(*object.String)

	method := classObj.GetMethod(key.Value)
	if method != nil {
		switch m := method.(type) {
		case *object.Function:
			fn := &object.Function{
				Name:       m.Name,
				Parameters: m.Parameters,
				Body:       m.Body,
				Env:        i.currentInstance.Fields,
				Instance:   i.currentInstance,
			}
			if classObj.Parent != nil {
				fn.Env.CreateConst("parent", classObj.Parent)
			}
			return fn
		case *object.BuiltinMethod:
			return &object.BuiltinMethod{
				Fn:       m.Fn,
				Instance: i.currentInstance,
			}
		}
	}

	return object.NullConst
}

func (i *Interpreter) evalAssignAttribute(
	e *ast.AttributeExpression,
	val ast.Expression,
	env *object.Environment) object.Object {
	indexed := i.Eval(e.Left, env)
	if isException(indexed) {
		return indexed
	}

	index := i.Eval(e.Index, env)
	if isException(indexed) {
		return indexed
	}
	instance := indexed.(*object.Instance)

	hashable, ok := index.(*object.String)
	if !ok {
		return object.NewException("Invalid index type %s", index.Type())
	}

	if _, ok := instance.Fields.Get(hashable.Value); !ok {
		return object.NewException("Instance has no field %s", hashable.Value)
	}

	if instance.Fields.IsConst(hashable.Value) {
		return object.NewException("Assignment to constant field %s", hashable.Value)
	}

	value := i.Eval(val, env)
	if isException(value) {
		return value
	}

	instance.Fields.SetForce(hashable.Value, value, false)
	return object.NullConst
}

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
