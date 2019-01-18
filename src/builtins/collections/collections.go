package collections

import (
	"sort"

	"github.com/nitrogen-lang/nitrogen/src/moduleutils"
	"github.com/nitrogen-lang/nitrogen/src/object"
	"github.com/nitrogen-lang/nitrogen/src/vm"
)

func init() {
	// Register with virtual machine
	vm.RegisterBuiltin("len", lenBuiltin)
	vm.RegisterBuiltin("first", firstBuiltin)
	vm.RegisterBuiltin("last", lastBuiltin)
	vm.RegisterBuiltin("rest", restBuiltin)
	vm.RegisterBuiltin("pop", popBuiltin)
	vm.RegisterBuiltin("push", pushBuiltin)
	vm.RegisterBuiltin("prepend", prependBuiltin)
	vm.RegisterBuiltin("splice", spliceBuiltin)
	vm.RegisterBuiltin("slice", sliceBuiltin)
	vm.RegisterBuiltin("sort", sortArrayBuiltin)
	vm.RegisterBuiltin("hashMerge", hashMergeBuiltin)
	vm.RegisterBuiltin("hashKeys", hashKeysBuiltin)
	vm.RegisterBuiltin("hasKey", hasKeyBuiltin)
}

func lenBuiltin(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewException("Incorrect number of arguments. Got %d, expected 1", len(args))
	}

	switch arg := args[0].(type) {
	case *object.String:
		return object.MakeIntObj(int64(len(arg.Value)))
	case *object.Array:
		return object.MakeIntObj(int64(len(arg.Elements)))
	case *object.Hash:
		return object.MakeIntObj(int64(len(arg.Pairs)))
	case *object.Null:
		return object.MakeIntObj(0)
	}

	return object.NewException("len(): Unsupported type %s", args[0].Type())
}

func firstBuiltin(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewException("Incorrect number of arguments. Got %d, expected 1", len(args))
	}
	if args[0].Type() != object.ArrayObj {
		return object.NewException("Argument to `first` must be ARRAY, got %s", args[0].Type())
	}

	arr := args[0].(*object.Array)
	if len(arr.Elements) > 0 {
		return arr.Elements[0]
	}

	return object.NullConst
}

func lastBuiltin(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewException("Incorrect number of arguments. Got %d, expected 1", len(args))
	}
	if args[0].Type() != object.ArrayObj {
		return object.NewException("Argument to `last` must be ARRAY, got %s", args[0].Type())
	}

	arr := args[0].(*object.Array)
	length := len(arr.Elements)
	if length > 0 {
		return arr.Elements[length-1]
	}

	return object.NullConst
}

func restBuiltin(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewException("Incorrect number of arguments. Got %d, expected 1", len(args))
	}
	if args[0].Type() != object.ArrayObj {
		return object.NewException("Argument to `rest` must be ARRAY, got %s", args[0].Type())
	}

	arr := args[0].(*object.Array)
	length := len(arr.Elements)
	if length > 0 {
		newElements := make([]object.Object, length-1, length-1)
		copy(newElements, arr.Elements[1:length])
		return &object.Array{Elements: newElements}
	}

	return object.NullConst
}

func popBuiltin(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewException("Incorrect number of arguments. Got %d, expected 1", len(args))
	}
	if args[0].Type() != object.ArrayObj {
		return object.NewException("Argument to `pop` must be ARRAY, got %s", args[0].Type())
	}

	arr := args[0].(*object.Array)
	length := len(arr.Elements)
	if length > 0 {
		newElements := make([]object.Object, length-1, length-1)
		copy(newElements, arr.Elements[:length-1])
		return &object.Array{Elements: newElements}
	}

	return object.NullConst
}

func pushBuiltin(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	if len(args) != 2 {
		return object.NewException("Incorrect number of arguments. Got %d, expected 2", len(args))
	}
	if args[0].Type() != object.ArrayObj {
		return object.NewException("Argument to `push` must be ARRAY, got %s", args[0].Type())
	}

	arr := args[0].(*object.Array)
	length := len(arr.Elements)
	newElements := make([]object.Object, length+1, length+1)
	copy(newElements, arr.Elements)
	newElements[length] = args[1]

	return &object.Array{Elements: newElements}
}

func prependBuiltin(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	if len(args) != 2 {
		return object.NewException("Incorrect number of arguments. Got %d, expected 2", len(args))
	}
	if args[0].Type() != object.ArrayObj {
		return object.NewException("Argument to `prepend` must be ARRAY, got %s", args[0].Type())
	}

	arr := args[0].(*object.Array)
	length := len(arr.Elements)
	newElements := make([]object.Object, length+1, length+1)
	copy(newElements[1:], arr.Elements)
	newElements[0] = args[1]

	return &object.Array{Elements: newElements}
}

func spliceBuiltin(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	if ac := moduleutils.CheckMinArgs("splice", 2, args...); ac != nil {
		return ac
	}

	if args[0].Type() != object.ArrayObj {
		return object.NewException("Argument 1 to `splice` must be ARRAY, got %s", args[0].Type())
	}
	arr := args[0].(*object.Array)

	offsetObj, ok := args[1].(*object.Integer)
	if !ok {
		return object.NewException("Argument 2 to `splice` must be INTEGER, got %s", args[1].Type())
	}
	offset := int(offsetObj.Value)
	if offset == 0 {
		return &object.Array{Elements: []object.Object{}}
	} else if offset < 0 {
		return object.NewException("Argument 2 to `splice` must be positive, got %d", offset)
	}

	orgLen := len(arr.Elements)
	length := orgLen - offset
	if len(args) > 2 {
		lenObj, ok := args[2].(*object.Integer)
		if !ok {
			return object.NewException("Argument 3 to `splice` must be INTEGER, got %s", args[2].Type())
		}
		length = int(lenObj.Value)
	}
	if length == 0 {
		return arr
	} else if length < 0 {
		return object.NewException("Argument 3 to `splice` must be positive, got %d", length)
	}

	newElements := make([]object.Object, orgLen-length, orgLen-length)
	copy(newElements, arr.Elements[:offset])
	copy(newElements[offset:], arr.Elements[offset+length:])

	return &object.Array{Elements: newElements}
}

func sliceBuiltin(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	if ac := moduleutils.CheckMinArgs("slice", 2, args...); ac != nil {
		return ac
	}

	if args[0].Type() != object.ArrayObj {
		return object.NewException("Argument 1 to `slice` must be ARRAY, got %s", args[0].Type())
	}
	arr := args[0].(*object.Array)

	offsetObj, ok := args[1].(*object.Integer)
	if !ok {
		return object.NewException("Argument 2 to `slice` must be INTEGER, got %s", args[1].Type())
	}
	offset := int(offsetObj.Value)
	if offset == 0 {
		return arr.Dup()
	} else if offset < 0 {
		return object.NewException("Argument 2 to `slice` must be positive, got %d", offset)
	}

	orgLen := len(arr.Elements)
	length := orgLen - offset
	if len(args) > 2 {
		lenObj, ok := args[2].(*object.Integer)
		if !ok {
			return object.NewException("Argument 3 to `slice` must be INTEGER, got %s", args[2].Type())
		}
		length = int(lenObj.Value)
	}
	if length == 0 {
		return arr
	} else if length < 0 {
		return object.NewException("Argument 3 to `slice` must be positive, got %d", length)
	}

	newElements := make([]object.Object, length, length)
	copy(newElements, arr.Elements[offset:length+1])

	return &object.Array{Elements: newElements}
}

func sortArrayBuiltin(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	if ac := moduleutils.CheckMinArgs("sort", 1, args...); ac != nil {
		return ac
	}

	arr, ok := args[0].(*object.Array)
	if !ok {
		return object.NewException("Argument to `sort` must be ARRAY, got %s", args[0].Type())
	}

	sorter := &arraySorter{arr.Dup().(*object.Array)}
	sort.Sort(sorter)

	return sorter.a
}

type arraySorter struct {
	a *object.Array
}

func (s *arraySorter) Len() int { return len(s.a.Elements) }
func (s *arraySorter) Less(i, j int) bool {
	return s.a.Elements[i].(*object.String).String() < s.a.Elements[j].(*object.String).String()
}
func (s *arraySorter) Swap(i, j int) {
	s.a.Elements[i], s.a.Elements[j] = s.a.Elements[j], s.a.Elements[i]
}

func hashMergeBuiltin(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	if len(args) < 2 {
		return object.NewException("hashMerge requires at least 2 arguments. Got %d", len(args))
	}

	if !object.ObjectsAre(object.HashObj, args[:2]...) {
		return object.NewException("First two arguments must be maps")
	}

	overwrite := true
	if len(args) > 2 {
		if args[2].Type() == object.BooleanObj {
			overwrite = args[2].(*object.Boolean).Value
		}
	}

	newMap := &object.Hash{
		Pairs: make(map[object.HashKey]object.HashPair),
	}

	for k, v := range args[0].(*object.Hash).Pairs {
		newMap.Pairs[k] = v
	}

	for k, v := range args[1].(*object.Hash).Pairs {
		if _, exists := newMap.Pairs[k]; !exists || (exists && overwrite) {
			newMap.Pairs[k] = v
		}
	}

	return newMap
}

func hashKeysBuiltin(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	if ac := moduleutils.CheckArgs("hashKeys", 1, args...); ac != nil {
		return ac
	}

	hash, ok := args[0].(*object.Hash)
	if !ok {
		return object.NewException("hashKeys expects a hash map")
	}

	arr := &object.Array{
		Elements: make([]object.Object, 0, len(hash.Pairs)),
	}

	for _, pair := range hash.Pairs {
		arr.Elements = append(arr.Elements, pair.Key)
	}

	return arr
}

func hasKeyBuiltin(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	if ac := moduleutils.CheckArgs("hasKey", 2, args...); ac != nil {
		return ac
	}

	hash, ok := args[0].(*object.Hash)
	if !ok {
		return object.NewException("hasKey arg 1 expects a hash map")
	}

	key, ok := args[1].(object.Hashable)
	if !ok {
		return object.NewException("hasKey arg 2 expects a valid hash key")
	}

	_, has := hash.Pairs[key.HashKey()]
	return object.NativeBoolToBooleanObj(has)
}
