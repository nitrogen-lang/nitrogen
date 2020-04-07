package vm

import (
	"github.com/nitrogen-lang/nitrogen/src/object"
)

var arrayIterator = &BuiltinClass{
	Fields: map[string]object.Object{
		"arr": object.NullConst,
		"i":   object.MakeIntObj(0),
	},
	VMClass: &VMClass{
		Name:   "ArrayIterator",
		Parent: nil,
		Methods: map[string]object.ClassMethod{
			"_next": MakeBuiltinMethod(arrayIteratorNext, 0),
		},
	},
}

func arrayIteratorNext(interpreter *VirtualMachine, self *VMInstance, env *object.Environment, args ...object.Object) object.Object {
	selfArrayObj, _ := self.Fields.Get("array")
	selfIndexObj, _ := self.Fields.Get("i")

	selfArray := selfArrayObj.(*object.Array)
	selfIndex := selfIndexObj.(*object.Integer)

	if int(selfIndex.Value) == len(selfArray.Elements) {
		return object.NullConst
	}

	elem := selfArray.Elements[selfIndex.Value]

	self.Fields.Set("i", object.MakeIntObj(selfIndex.Value+1))

	return &object.Array{
		Elements: []object.Object{
			object.MakeIntObj(selfIndex.Value),
			elem,
		},
	}
}

func makeArrayIter(array *object.Array) *VMInstance {
	env := object.NewEnvironment()
	env.SetForce("array", array, true)
	env.SetForce("i", object.MakeIntObj(0), false)

	return &VMInstance{
		Class:  arrayIterator.VMClass,
		Fields: env,
	}
}

var mapIterator = &BuiltinClass{
	Fields: map[string]object.Object{
		"map":  object.NullConst,
		"keys": object.NullConst,
		"i":    object.MakeIntObj(0),
	},
	VMClass: &VMClass{
		Name:   "MapIterator",
		Parent: nil,
		Methods: map[string]object.ClassMethod{
			"_next": MakeBuiltinMethod(mapIteratorNext, 0),
		},
	},
}

func mapIteratorNext(interpreter *VirtualMachine, self *VMInstance, env *object.Environment, args ...object.Object) object.Object {
	selfMapObj, _ := self.Fields.Get("map")
	selfKeysObj, _ := self.Fields.Get("keys")
	selfIndexObj, _ := self.Fields.Get("i")

	selfMap := selfMapObj.(*object.Hash)
	selfKeys := selfKeysObj.(*object.Array)
	selfIndex := selfIndexObj.(*object.Integer)

	if int(selfIndex.Value) == len(selfKeys.Elements) {
		return object.NullConst
	}

	mapKey := selfKeys.Elements[selfIndex.Value]
	hashKey := mapKey.(object.Hashable)
	hashVal := selfMap.Pairs[hashKey.HashKey()]

	self.Fields.Set("i", object.MakeIntObj(selfIndex.Value+1))

	return &object.Array{
		Elements: []object.Object{
			mapKey,
			hashVal.Value,
		},
	}
}

func makeMapIter(hashmap *object.Hash) *VMInstance {
	env := object.NewEnvironment()
	env.SetForce("map", hashmap, true)
	env.SetForce("keys", hashmap.Keys(), true)
	env.SetForce("i", object.MakeIntObj(0), false)

	return &VMInstance{
		Class:  mapIterator.VMClass,
		Fields: env,
	}
}

var stringIterator = &BuiltinClass{
	Fields: map[string]object.Object{
		"str": object.NullConst,
		"i":   object.MakeIntObj(0),
	},
	VMClass: &VMClass{
		Name:   "StringIterator",
		Parent: nil,
		Methods: map[string]object.ClassMethod{
			"_next": MakeBuiltinMethod(stringIteratorNext, 0),
		},
	},
}

func stringIteratorNext(interpreter *VirtualMachine, self *VMInstance, env *object.Environment, args ...object.Object) object.Object {
	selfStrObj, _ := self.Fields.Get("str")
	selfIndexObj, _ := self.Fields.Get("i")

	selfStr := selfStrObj.(*object.String)
	selfIndex := selfIndexObj.(*object.Integer)

	if int(selfIndex.Value) == len(selfStr.Value) {
		return object.NullConst
	}

	elem := selfStr.Value[selfIndex.Value]

	self.Fields.Set("i", object.MakeIntObj(selfIndex.Value+1))

	return &object.Array{
		Elements: []object.Object{
			object.MakeIntObj(selfIndex.Value),
			object.MakeStringObjRunes([]rune{elem}),
		},
	}
}

func makeStringIter(str *object.String) *VMInstance {
	env := object.NewEnvironment()
	env.SetForce("str", str, true)
	env.SetForce("i", object.MakeIntObj(0), false)

	return &VMInstance{
		Class:  stringIterator.VMClass,
		Fields: env,
	}
}
