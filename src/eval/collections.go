package eval

import (
	"github.com/nitrogen-lang/nitrogen/src/ast"
	"github.com/nitrogen-lang/nitrogen/src/object"
)

func (i *Interpreter) evalIndexExpression(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ArrayObj && index.Type() == object.IntergerObj:
		return i.evalArrayIndexExpression(left, index)
	case left.Type() == object.HashObj:
		return i.evalHashIndexExpression(left, index)
	case left.Type() == object.StringObj && index.Type() == object.IntergerObj:
		return i.evalStringIndexExpression(left, index)
	case left.Type() == object.ModuleObj && index.Type() == object.StringObj:
		return i.evalModuleLookupExpression(left, index)
	case left.Type() == object.InstanceObj && index.Type() == object.StringObj:
		return i.evalInstanceLookupExpression(left, index)
	}
	return object.NewException("Index operator not allowed: %s", left.Type())
}

func (i *Interpreter) evalArrayIndexExpression(array, index object.Object) object.Object {
	arrObj := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrObj.Elements))

	if idx > max-1 { // Check upper bound
		return object.NullConst
	}

	if idx < 0 { // Check lower bound
		// Convert a negative index to positive
		idx = max + idx

		if idx < 0 { // Check lower bound again
			return object.NullConst
		}
	}

	return arrObj.Elements[idx]
}

func (i *Interpreter) evalHashIndexExpression(hash, index object.Object) object.Object {
	hashObj := hash.(*object.Hash)

	key, ok := index.(object.Hashable)
	if !ok {
		return object.NewException("Invalid map key: %s", index.Type())
	}

	pair, ok := hashObj.Pairs[key.HashKey()]
	if !ok {
		return object.NullConst
	}

	return pair.Value
}

func (i *Interpreter) evalModuleLookupExpression(module, index object.Object) object.Object {
	moduleObj := module.(*object.Module)
	key := index.(*object.String)

	// Methods have priority over variables
	method, ok := moduleObj.Methods[key.Value]
	if ok {
		return &object.Builtin{Fn: method}
	}

	variable, ok := moduleObj.Vars[key.Value]
	if ok {
		return variable
	}
	return object.NullConst
}

func (i *Interpreter) evalHashLiteral(node *ast.HashLiteral, env *object.Environment) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)

	for keyNode, valueNode := range node.Pairs {
		key := i.Eval(keyNode, env)
		if isException(key) {
			return key
		}

		hashKey, ok := key.(object.Hashable)
		if !ok {
			return object.NewException("Invalid map key: %s", key.Type())
		}

		value := i.Eval(valueNode, env)
		if isException(value) {
			return value
		}

		hashed := hashKey.HashKey()
		pairs[hashed] = object.HashPair{Key: key, Value: value}
	}

	return &object.Hash{Pairs: pairs}
}

func (i *Interpreter) evalStringIndexExpression(array, index object.Object) object.Object {
	strObj := array.(*object.String)
	idx := index.(*object.Integer).Value
	max := int64(len(strObj.Value))

	if idx > max-1 { // Check upper bound
		return object.NullConst
	}

	if idx < 0 { // Check lower bound
		// Convert a negative index to positive
		idx = max + idx

		if idx < 0 { // Check lower bound again
			return object.NullConst
		}
	}

	return &object.String{Value: string(strObj.Value[idx])}
}

func (i *Interpreter) evalInstanceLookupExpression(instance, index object.Object) object.Object {
	instanceObj := instance.(*object.Instance)
	key := index.(*object.String)

	method := instanceObj.GetMethod(key.Value)
	if method != nil {
		return &object.Function{
			Name:       method.Name,
			Parameters: method.Parameters,
			Body:       method.Body,
			Env:        object.NewEnclosedEnv(instanceObj.Fields),
		}
	}

	val, ok := instanceObj.Fields.Get(key.Value)
	if ok {
		return val
	}
	return object.NullConst
}
