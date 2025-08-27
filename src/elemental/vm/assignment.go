package vm

import "github.com/nitrogen-lang/nitrogen/src/elemental/object"

func (vm *VirtualMachine) assignIndexedValue(indexed, index, val object.Object) object.Object {

	switch i := indexed.(type) {
	case *object.Array:
		return vm.assignArrayIndex(i, index, val)
	case *object.Hash:
		return vm.assignHashMapIndex(i, index, val)
	case *object.String:
		return vm.assignStringIndex(i, index, val)
	case *object.ByteString:
		return vm.assignByteStringIndex(i, index, val)
	}
	return object.NullConst
}

func (vm *VirtualMachine) assignStringIndex(
	str *object.String,
	index, val object.Object) object.Object {

	in, ok := index.(*object.Integer)
	if !ok {
		return object.NewException("Invalid string index type %s", index.Type())
	}

	if in.Value < 0 || in.Value > int64(len(str.Value)-1) {
		return object.NewException("Index out of bounds: %s", index.Inspect())
	}

	replace, ok := val.(*object.String)
	if !ok {
		return object.NewException("Invalid string index value type %s", val.Type())
	}

	str.Value[in.Value] = replace.Value[0]
	return object.NullConst
}

func (vm *VirtualMachine) assignByteStringIndex(
	str *object.ByteString,
	index, val object.Object) object.Object {

	in, ok := index.(*object.Integer)
	if !ok {
		return object.NewException("Invalid byte string index type %s", index.Type())
	}

	if in.Value < 0 || in.Value > int64(len(str.Value)-1) {
		return object.NewException("Index out of bounds: %s", index.Inspect())
	}

	replace, ok := val.(*object.ByteString)
	if !ok {
		return object.NewException("Invalid byte string index value type %s", val.Type())
	}

	str.Value[in.Value] = replace.Value[0]
	return object.NullConst
}

func (vm *VirtualMachine) assignArrayIndex(
	array *object.Array,
	index, val object.Object) object.Object {

	in, ok := index.(*object.Integer)
	if !ok {
		return object.NewException("Invalid array index type %s", index.Type())
	}

	if in.Value < 0 || in.Value > int64(len(array.Elements)-1) {
		return object.NewException("Index out of bounds: %s", index.Inspect())
	}

	array.Elements[in.Value] = val
	return object.NullConst
}

func (vm *VirtualMachine) assignHashMapIndex(
	hashmap *object.Hash,
	index, val object.Object) object.Object {

	hashable, ok := index.(object.Hashable)
	if !ok {
		return object.NewException("Invalid index type %s", index.Type())
	}

	hashmap.Pairs[hashable.HashKey()] = object.HashPair{
		Key:   index,
		Value: val,
	}
	return object.NullConst
}

func (vm *VirtualMachine) assignModuleAttr(module *object.Module, key string, val object.Object) object.Object {
	// Methods have priority over variables
	_, ok := module.Methods[key]
	if ok {
		return object.NewPanic("Cannot assign to module method %s", key)
	}

	module.Vars[key] = val
	return object.NullConst
}
