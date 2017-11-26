package vm

import "github.com/nitrogen-lang/nitrogen/src/object"

func (vm *VirtualMachine) evalIndexExpression(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ArrayObj && index.Type() == object.IntergerObj:
		return vm.evalArrayIndexExpression(left, index)
	case left.Type() == object.HashObj:
		return vm.evalHashIndexExpression(left, index)
	case left.Type() == object.StringObj && index.Type() == object.IntergerObj:
		return vm.evalStringIndexExpression(left, index)
	case left.Type() == object.ModuleObj && index.Type() == object.StringObj:
		return vm.evalModuleLookupExpression(left, index)
	}
	return object.NewException("Index operator not allowed: %s", left.Type())
}

func (vm *VirtualMachine) evalArrayIndexExpression(array, index object.Object) object.Object {
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

func (vm *VirtualMachine) evalHashIndexExpression(hash, index object.Object) object.Object {
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

func (vm *VirtualMachine) evalStringIndexExpression(array, index object.Object) object.Object {
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

func (vm *VirtualMachine) evalModuleLookupExpression(module, index object.Object) object.Object {
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
