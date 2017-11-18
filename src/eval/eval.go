package eval

import (
	"io"
	"os"

	"github.com/nitrogen-lang/nitrogen/src/ast"
	"github.com/nitrogen-lang/nitrogen/src/object"
)

type Interpreter struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer

	scriptNameStack *stringStack
}

func NewInterpreter() *Interpreter {
	return &Interpreter{
		Stdin:           os.Stdin,
		Stdout:          os.Stdout,
		Stderr:          os.Stderr,
		scriptNameStack: newStringStack(),
	}
}

func (i *Interpreter) Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	// Statements
	case *ast.Program:
		return i.evalProgram(node, env)
	case *ast.ExpressionStatement:
		return i.Eval(node.Expression, env)
	case *ast.BlockStatement:
		return i.evalBlockStatements(node, env)
	case *ast.ReturnStatement:
		val := i.Eval(node.Value, env)
		if isException(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	case *ast.DefStatement:
		if node.Const {
			return i.assignConstIdentValue(node.Name, node.Value, env)
		}
		return i.assignIdentValue(node.Name, node.Value, true, env)
	case *ast.AssignStatement:
		return i.evalAssignment(node, env)

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
		elements := i.evalExpressions(node.Elements, env)
		if len(elements) == 1 && isException(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}
	case *ast.Boolean:
		return object.NativeBoolToBooleanObj(node.Value)
	case *ast.HashLiteral:
		return i.evalHashLiteral(node, env)

	// Expressions
	case *ast.Identifier:
		return i.evalIdent(node, env)
	case *ast.PrefixExpression:
		right := i.Eval(node.Right, env)
		if isException(right) {
			return right
		}
		return i.evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		right := i.Eval(node.Right, env)
		if isException(right) {
			return right
		}

		left := i.Eval(node.Left, env)
		if isException(left) {
			return left
		}

		return i.evalInfixExpression(node.Operator, left, right)
	case *ast.IndexExpression:
		left := i.Eval(node.Left, env)
		if isException(left) {
			return left
		}

		index := i.Eval(node.Index, env)
		if isException(index) {
			return index
		}
		return i.evalIndexExpression(left, index)

	// Conditionals
	case *ast.IfExpression:
		return i.evalIfExpression(node, env)
	case *ast.CompareExpression:
		return i.evalCompareExpression(node, env)
	case *ast.ForLoopStatement:
		return i.evalForLoop(node, env)
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
		function := i.Eval(node.Function, env)
		if isException(function) {
			if ident, ok := node.Function.(*ast.Identifier); ok {
				return object.NewException("function not found: %s", ident.Value)
			}
			return function
		}

		args := i.evalExpressions(node.Arguments, env)
		if len(args) == 1 && isException(args[0]) {
			return args[0]
		}

		return i.applyFunction(function, args, env)
	}

	return nil
}

// GetCurrentScriptPath returns the filepath of the current executing script
func (i *Interpreter) GetCurrentScriptPath() string {
	return i.scriptNameStack.getFront()
}

func (i *Interpreter) GetStdout() io.Writer {
	return i.Stdout
}
func (i *Interpreter) GetStderr() io.Writer {
	return i.Stderr
}
func (i *Interpreter) GetStdin() io.Reader {
	return i.Stdin
}

func (i *Interpreter) evalProgram(p *ast.Program, env *object.Environment) object.Object {
	var result object.Object
	i.scriptNameStack.push(p.Filename)
	defer i.scriptNameStack.pop()
	env.CreateConst("_FILE", &object.String{Value: p.Filename})

	for _, statement := range p.Statements {
		result = i.Eval(statement, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Exception:
			return result
		}
	}

	return result
}

func (i *Interpreter) evalBlockStatements(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = i.Eval(statement, env)

		if result != nil {
			rt := result.Type()
			if rt == object.ReturnObj || rt == object.ExceptionObj || rt == object.LoopControlObj {
				return result
			}
		}
	}

	return result
}

func (i *Interpreter) evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, e := range exps {
		evaled := i.Eval(e, env)
		if isException(evaled) {
			return []object.Object{evaled}
		}
		result = append(result, evaled)
	}

	return result
}

func (i *Interpreter) evalIdent(node *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}
	if builtin := getBuiltin(node.Value); builtin != nil {
		return builtin
	}
	return object.NewException("identifier not found: %s", node.Value)
}
