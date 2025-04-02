package vm

import "github.com/nitrogen-lang/nitrogen/src/object"

func (vm *VirtualMachine) evalIndexExpression(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ArrayObj && index.Type() == object.IntergerObj:
		return vm.evalArrayIndexExpression(left.(*object.Array), index)
	case left.Type() == object.HashObj:
		return vm.lookupHashIndex(left.(*object.Hash), index)
	case left.Type() == object.StringObj && index.Type() == object.IntergerObj:
		return vm.evalStringIndexExpression(left.(*object.String), index)
	case left.Type() == object.ByteStringObj && index.Type() == object.IntergerObj:
		return vm.evalByteStringIndexExpression(left.(*object.ByteString), index)
	}
	return object.NewException("Index operator not allowed on type %s", left.Type())
}

func (vm *VirtualMachine) evalArrayIndexExpression(array *object.Array, index object.Object) object.Object {
	idx := index.(*object.Integer).Value
	max := int64(len(array.Elements))

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

	return array.Elements[idx]
}

func (vm *VirtualMachine) lookupHashIndex(hash *object.Hash, index object.Object) object.Object {
	key, ok := index.(object.Hashable)
	if !ok {
		return object.NewException("Invalid map key: %s", index.Type())
	}

	pair, ok := hash.Pairs[key.HashKey()]
	if !ok {
		return object.NullConst
	}

	return pair.Value
}

func (vm *VirtualMachine) evalStringIndexExpression(str *object.String, index object.Object) object.Object {
	idx := index.(*object.Integer).Value
	max := int64(len(str.Value))

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

	return &object.String{Value: []rune{str.Value[idx]}}
}

func (vm *VirtualMachine) evalByteStringIndexExpression(str *object.ByteString, index object.Object) object.Object {
	idx := index.(*object.Integer).Value
	max := int64(len(str.Value))

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

	return &object.ByteString{Value: []byte{str.Value[idx]}}
}

func (vm *VirtualMachine) lookupModuleAttr(module *object.Module, key string) object.Object {
	// Methods have priority over variables
	method, ok := module.Methods[key]
	if ok {
		return &object.Builtin{Fn: method}
	}

	variable, ok := module.Vars[key]
	if ok {
		return variable
	}

	return object.NullConst
}
