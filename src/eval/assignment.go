package eval

import (
	"github.com/nitrogen-lang/nitrogen/src/ast"
	"github.com/nitrogen-lang/nitrogen/src/object"
)

func evalAssignment(stmt *ast.AssignStatement, env *object.Environment) object.Object {
	if left, ok := stmt.Left.(*ast.IndexExpression); ok {
		return assignIndexedValue(left, stmt.Value, env)
	}

	ident, ok := stmt.Left.(*ast.Identifier)
	if !ok {
		return object.NewException("Invalid variable name, expected identifier, got %s",
			stmt.Left.String())
	}

	return assignIdentValue(ident, stmt.Value, false, env)
}

func assignIdentValue(
	name *ast.Identifier,
	val ast.Expression,
	new bool,
	env *object.Environment) object.Object {
	// Protect builtin functions
	if builtin := getBuiltin(name.Value); builtin != nil {
		return object.NewException(
			"Attempted redeclaration of builtin function '%s'",
			name.Value,
		)
	}

	if !new { // Variables must be declared before use
		if _, exists := env.Get(name.Value); !exists {
			return object.NewException("Assignment to uninitialized variable %s", name.Value)
		}
	}

	if env.IsConst(name.Value) {
		return object.NewException("Assignment to declared constant %s", name.Value)
	}

	evaled := Eval(val, env)
	if isException(evaled) {
		return evaled
	}

	// Ignore error since we check for consant above
	if new {
		env.Create(name.Value, evaled)
	} else {
		env.Set(name.Value, evaled)
	}
	return object.NULL
}

func assignConstIdentValue(
	name *ast.Identifier,
	val ast.Expression,
	env *object.Environment) object.Object {
	// Protect builtin functions
	if builtin := getBuiltin(name.Value); builtin != nil {
		return object.NewException(
			"Attempted redeclaration of builtin function '%s'",
			name.Value,
		)
	}

	if _, exists := env.Get(name.Value); exists { // Constants can't redeclare an existing var
		return object.NewException("Can't assign constant to variable `%s`", name.Value)
	}

	evaled := Eval(val, env)
	if isException(evaled) {
		return evaled
	}

	if !object.ObjectIs(evaled, object.INTEGER_OBJ, object.FLOAT_OBJ, object.STRING_OBJ, object.NULL_OBJ, object.BOOLEAN_OBJ) {
		return object.NewException("Constants must be int, float, string, bool or null")
	}

	// Ignore error since we check above
	env.CreateConst(name.Value, evaled)
	return object.NULL
}

func assignIndexedValue(
	e *ast.IndexExpression,
	val ast.Expression,
	env *object.Environment) object.Object {
	indexed := Eval(e.Left, env)
	if isException(indexed) {
		return indexed
	}

	index := Eval(e.Index, env)
	if isException(indexed) {
		return indexed
	}

	switch indexed.Type() {
	case object.ARRAY_OBJ:
		return assignArrayIndex(indexed.(*object.Array), index, val, env)
	case object.HASH_OBJ:
		return assignHashMapIndex(indexed.(*object.Hash), index, val, env)
	}
	return object.NULL
}

func assignArrayIndex(
	array *object.Array,
	index object.Object,
	val ast.Expression,
	env *object.Environment) object.Object {

	in, ok := index.(*object.Integer)
	if !ok {
		return object.NewException("Invalid array index type %s", index.(object.Object).Type())
	}

	value := Eval(val, env)
	if isException(value) {
		return value
	}

	if in.Value < 0 || in.Value > int64(len(array.Elements)-1) {
		return object.NewException("Index out of bounds: %s", index.Inspect())
	}

	array.Elements[in.Value] = value
	return object.NULL
}

func assignHashMapIndex(
	hashmap *object.Hash,
	index object.Object,
	val ast.Expression,
	env *object.Environment) object.Object {

	hashable, ok := index.(object.Hashable)
	if !ok {
		return object.NewException("Invalid index type %s", index.Type())
	}

	value := Eval(val, env)
	if isException(value) {
		return value
	}

	hashmap.Pairs[hashable.HashKey()] = object.HashPair{
		Key:   index,
		Value: value,
	}
	return object.NULL
}
