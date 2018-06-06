package vm

import "github.com/nitrogen-lang/nitrogen/src/object"

func (vm *VirtualMachine) assignIndexedValue(
	indexed object.Object,
	index object.Object,
	val object.Object) object.Object {
	switch indexed.Type() {
	case object.ArrayObj:
		return vm.assignArrayIndex(indexed.(*object.Array), index, val)
	case object.HashObj:
		return vm.assignHashMapIndex(indexed.(*object.Hash), index, val)
	}
	return object.NullConst
}

func (vm *VirtualMachine) assignArrayIndex(
	array *object.Array,
	index object.Object,
	val object.Object) object.Object {

	in, ok := index.(*object.Integer)
	if !ok {
		return object.NewException("Invalid array index type %s", index.(object.Object).Type())
	}

	if in.Value < 0 || in.Value > int64(len(array.Elements)-1) {
		return object.NewException("Index out of bounds: %s", index.Inspect())
	}

	array.Elements[in.Value] = val
	return object.NullConst
}

func (vm *VirtualMachine) assignHashMapIndex(
	hashmap *object.Hash,
	index object.Object,
	val object.Object) object.Object {

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
