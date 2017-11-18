package eval

import (
	"github.com/nitrogen-lang/nitrogen/src/ast"
	"github.com/nitrogen-lang/nitrogen/src/object"
	"github.com/nitrogen-lang/nitrogen/src/token"
)

func (i *Interpreter) evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := i.Eval(ie.Condition, env)
	if isException(condition) {
		return condition
	}

	if isTruthy(condition) {
		return i.Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return i.Eval(ie.Alternative, env)
	}

	return object.NullConst
}

func (i *Interpreter) evalForLoop(loop *ast.ForLoopStatement, env *object.Environment) object.Object {
	outterScope := object.NewEnclosedEnv(env)

	if loop.Init != nil {
		init := i.Eval(loop.Init, outterScope)
		if isException(init) {
			return init
		}

		// If the iterator is not an assignment, generate one using the ident from the initializer
		if _, ok := loop.Iter.(*ast.AssignStatement); !ok {
			loop.Iter = &ast.AssignStatement{
				Left:  loop.Init.Name,
				Value: loop.Iter,
			}
		}
	}

	for {
		// Check loop condition
		if loop.Condition != nil {
			condition := i.Eval(loop.Condition, outterScope)
			if isException(condition) {
				return condition
			}
			if !isTruthy(condition) {
				break
			}
		}

		// Execute body
		body := i.Eval(loop.Body, object.NewEnclosedEnv(outterScope))
		if isException(body) {
			return body
		}

		// Return if necessary
		rt := body.Type()
		if rt == object.ReturnObj {
			return body
		}

		// Break if necessary, continue automatically
		if rt == object.LoopControlObj {
			if !body.(*object.LoopControl).Continue {
				break
			}
		}

		// Execute iterator
		if loop.Iter != nil {
			iter := i.Eval(loop.Iter, outterScope)
			if isException(iter) {
				return iter
			}
		}
	}
	return object.NullConst
}

func (i *Interpreter) evalCompareExpression(node *ast.CompareExpression, env *object.Environment) object.Object {
	left := i.Eval(node.Left, env)
	if isException(left) {
		return left
	}

	lBool, valid := convertToBoolean(left)
	if !valid {
		return object.NewException("Left side of conditional must be truthy or falsey")
	}

	// Short circuit if possible
	if node.Token.Type == token.LOr && lBool {
		return object.TrueConst
	}
	if node.Token.Type == token.LAnd && !lBool {
		return object.FalseConst
	}

	right := i.Eval(node.Right, env)
	if isException(right) {
		return right
	}

	rBool, valid := convertToBoolean(right)
	if !valid {
		return object.NewException("Right side of condition must be truthy or falsey")
	}

	return object.NativeBoolToBooleanObj(rBool)
}
