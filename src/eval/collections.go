package eval

import (
	"github.com/nitrogen-lang/nitrogen/src/ast"
	"github.com/nitrogen-lang/nitrogen/src/object"
)

func evalIndexExpression(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpression(left, index)
	case left.Type() == object.HASH_OBJ:
		return evalHashIndexExpression(left, index)
	}
	return object.NewException("Index operator not allowed: %s", left.Type())
}

func evalArrayIndexExpression(array, index object.Object) object.Object {
	arrObj := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrObj.Elements))

	if idx > max-1 { // Check upper bound
		return object.NULL
	}

	if idx < 0 { // Check lower bound
		// Convert a negative index to positive
		idx = max + idx

		if idx < 0 { // Check lower bound again
			return object.NULL
		}
	}

	return arrObj.Elements[idx]
}

func evalHashIndexExpression(hash, index object.Object) object.Object {
	hashObj := hash.(*object.Hash)

	key, ok := index.(object.Hashable)
	if !ok {
		return object.NewException("Invalid map key: %s", index.Type())
	}

	pair, ok := hashObj.Pairs[key.HashKey()]
	if !ok {
		return object.NULL
	}

	return pair.Value
}

func evalHashLiteral(node *ast.HashLiteral, env *object.Environment) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)

	for keyNode, valueNode := range node.Pairs {
		key := Eval(keyNode, env)
		if isException(key) {
			return key
		}

		hashKey, ok := key.(object.Hashable)
		if !ok {
			return object.NewException("Invalid map key: %s", key.Type())
		}

		value := Eval(valueNode, env)
		if isException(value) {
			return value
		}

		hashed := hashKey.HashKey()
		pairs[hashed] = object.HashPair{Key: key, Value: value}
	}

	return &object.Hash{Pairs: pairs}
}
