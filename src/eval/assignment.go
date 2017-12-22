package eval

import (
	"github.com/nitrogen-lang/nitrogen/src/ast"
	"github.com/nitrogen-lang/nitrogen/src/object"
)

func (i *Interpreter) evalAssignment(stmt *ast.AssignStatement, env *object.Environment) object.Object {
	if left, ok := stmt.Left.(*ast.IndexExpression); ok {
		return i.assignIndexedValue(left, stmt.Value, env)
	}
	if left, ok := stmt.Left.(*ast.AttributeExpression); ok {
		return i.evalAssignAttribute(left, stmt.Value, env)
	}

	ident, ok := stmt.Left.(*ast.Identifier)
	if !ok {
		return object.NewException("Invalid variable name, expected identifier, got %s",
			stmt.Left.String())
	}

	return i.assignIdentValue(ident, stmt.Value, false, env)
}

func (i *Interpreter) assignIdentValue(
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

	var evaled object.Object = object.NullConst
	if val != nil {
		evaled = i.Eval(val, env)
		if isException(evaled) {
			return evaled
		}
	}

	// Ignore error since we check for consant above
	if new {
		env.Create(name.Value, evaled)
	} else {
		env.Set(name.Value, evaled)
	}
	return object.NullConst
}

func (i *Interpreter) assignConstIdentValue(
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

	evaled := i.Eval(val, env)
	if isException(evaled) {
		return evaled
	}

	if !object.ObjectIs(evaled, object.IntergerObj, object.FloatObj, object.StringObj, object.NullObj, object.BooleanObj, object.ModuleObj) {
		return object.NewException("Constants must be int, float, string, bool or null")
	}

	// Ignore error since we check above
	env.CreateConst(name.Value, evaled)
	return object.NullConst
}

func (i *Interpreter) assignIndexedValue(
	e *ast.IndexExpression,
	val ast.Expression,
	env *object.Environment) object.Object {
	indexed := i.Eval(e.Left, env)
	if isException(indexed) {
		return indexed
	}

	index := i.Eval(e.Index, env)
	if isException(indexed) {
		return indexed
	}

	switch indexed.Type() {
	case object.ArrayObj:
		return i.assignArrayIndex(indexed.(*object.Array), index, val, env)
	case object.HashObj:
		return i.assignHashMapIndex(indexed.(*object.Hash), index, val, env)
	case object.ModuleObj:
		return i.assignModuleVariable(indexed.(*object.Module), index, val, env)
	}
	return object.NullConst
}

func (i *Interpreter) assignArrayIndex(
	array *object.Array,
	index object.Object,
	val ast.Expression,
	env *object.Environment) object.Object {

	in, ok := index.(*object.Integer)
	if !ok {
		return object.NewException("Invalid array index type %s", index.(object.Object).Type())
	}

	value := i.Eval(val, env)
	if isException(value) {
		return value
	}

	if in.Value < 0 || in.Value > int64(len(array.Elements)-1) {
		return object.NewException("Index out of bounds: %s", index.Inspect())
	}

	array.Elements[in.Value] = value
	return object.NullConst
}

func (i *Interpreter) assignHashMapIndex(
	hashmap *object.Hash,
	index object.Object,
	val ast.Expression,
	env *object.Environment) object.Object {

	hashable, ok := index.(object.Hashable)
	if !ok {
		return object.NewException("Invalid index type %s", index.Type())
	}

	value := i.Eval(val, env)
	if isException(value) {
		return value
	}

	hashmap.Pairs[hashable.HashKey()] = object.HashPair{
		Key:   index,
		Value: value,
	}
	return object.NullConst
}

func (i *Interpreter) assignModuleVariable(
	module *object.Module,
	index object.Object,
	val ast.Expression,
	env *object.Environment) object.Object {

	hashable, ok := index.(*object.String)
	if !ok {
		return object.NewException("Invalid index type %s", index.Type())
	}

	if _, exists := module.Vars[hashable.Value]; !exists {
		return object.NewException("Module %s has no assignable variable %s", module.Name, hashable.Value)
	}

	value := i.Eval(val, env)
	if isException(value) {
		return value
	}

	module.Vars[hashable.Value] = value
	return object.NullConst
}
