package eval

import (
	"io"
	"os"

	"github.com/nitrogen-lang/nitrogen/src/ast"
	"github.com/nitrogen-lang/nitrogen/src/object"
)

var (
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer

	scriptNameStack = newStringStack()
)

func init() {
	Stdin = os.Stdin
	Stdout = os.Stdout
	Stderr = os.Stderr
}

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	// Statements
	case *ast.Program:
		return evalProgram(node, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.BlockStatement:
		return evalBlockStatements(node, env)
	case *ast.ReturnStatement:
		val := Eval(node.Value, env)
		if isException(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	case *ast.DefStatement:
		if node.Const {
			return assignConstIdentValue(node.Name, node.Value, env)
		}
		return assignIdentValue(node.Name, node.Value, true, env)
	case *ast.AssignStatement:
		return evalAssignment(node, env)

	// Literals
	case *ast.NullLiteral:
		return object.NullConst
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.FloatLiteral:
		return &object.Float{Value: node.Value}
	case *ast.Array:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && isException(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}
	case *ast.Boolean:
		return object.NativeBoolToBooleanObj(node.Value)
	case *ast.HashLiteral:
		return evalHashLiteral(node, env)

	// Expressions
	case *ast.Identifier:
		return evalIdent(node, env)
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isException(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		right := Eval(node.Right, env)
		if isException(right) {
			return right
		}

		left := Eval(node.Left, env)
		if isException(left) {
			return left
		}

		return evalInfixExpression(node.Operator, left, right)
	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if isException(left) {
			return left
		}

		index := Eval(node.Index, env)
		if isException(index) {
			return index
		}
		return evalIndexExpression(left, index)

	// Conditionals
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.CompareExpression:
		return evalCompareExpression(node, env)
	case *ast.ForLoopStatement:
		return evalForLoop(node, env)
	case *ast.ContinueStatement:
		return &object.LoopControl{Continue: true}
	case *ast.BreakStatement:
		return &object.LoopControl{Continue: false}

	// Functions
	case *ast.FunctionLiteral:
		return &object.Function{
			Name:       node.Name,
			Parameters: node.Parameters,
			Body:       node.Body,
			Env:        env,
		}
	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isException(function) {
			if ident, ok := node.Function.(*ast.Identifier); ok {
				return object.NewException("function not found: %s", ident.Value)
			}
			return function
		}

		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isException(args[0]) {
			return args[0]
		}

		return applyFunction(function, args, env)
	}

	return nil
}

// GetCurrentScriptPath returns the filepath of the current executing script
func GetCurrentScriptPath() string {
	return scriptNameStack.getFront()
}

func evalProgram(p *ast.Program, env *object.Environment) object.Object {
	var result object.Object
	scriptNameStack.push(p.Filename)
	defer scriptNameStack.pop()
	env.CreateConst("_FILE", &object.String{Value: p.Filename})

	for _, statement := range p.Statements {
		result = Eval(statement, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Exception:
			return result
		}
	}

	return result
}

func evalBlockStatements(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement, env)

		if result != nil {
			rt := result.Type()
			if rt == object.ReturnObj || rt == object.ExceptionObj || rt == object.LoopControlObj {
				return result
			}
		}
	}

	return result
}

func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, e := range exps {
		evaled := Eval(e, env)
		if isException(evaled) {
			return []object.Object{evaled}
		}
		result = append(result, evaled)
	}

	return result
}

func evalIdent(node *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}
	if builtin := getBuiltin(node.Value); builtin != nil {
		return builtin
	}
	return object.NewException("identifier not found: %s", node.Value)
}
