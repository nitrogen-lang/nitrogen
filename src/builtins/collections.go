package builtins

import (
	"sort"

	"github.com/nitrogen-lang/nitrogen/src/eval"
	"github.com/nitrogen-lang/nitrogen/src/moduleutils"
	"github.com/nitrogen-lang/nitrogen/src/object"
)

func init() {
	eval.RegisterBuiltin("len", lenBuiltin)
	eval.RegisterBuiltin("first", firstBuiltin)
	eval.RegisterBuiltin("last", lastBuiltin)
	eval.RegisterBuiltin("rest", restBuiltin)
	eval.RegisterBuiltin("push", pushBuiltin)
	eval.RegisterBuiltin("sort", sortArrayBuiltin)
	eval.RegisterBuiltin("hashMerge", hashMergeBuiltin)
	eval.RegisterBuiltin("hashKeys", hashKeysBuiltin)
}

func lenBuiltin(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewException("Incorrect number of arguments. Got %d, expected 1", len(args))
	}

	switch arg := args[0].(type) {
	case *object.String:
		return &object.Integer{Value: int64(len(arg.Value))}
	case *object.Array:
		return &object.Integer{Value: int64(len(arg.Elements))}
	case *object.Null:
		return &object.Integer{Value: 0}
	}

	return object.NewException("Unsupported type %s", args[0].Type())
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
	return s.a.Elements[i].(*object.String).Value < s.a.Elements[j].(*object.String).Value
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
