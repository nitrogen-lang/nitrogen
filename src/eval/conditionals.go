package eval

import (
	"github.com/nitrogen-lang/nitrogen/src/ast"
	"github.com/nitrogen-lang/nitrogen/src/object"
	"github.com/nitrogen-lang/nitrogen/src/token"
)

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	}

	return object.NULL
}

func evalForLoop(loop *ast.ForLoopStatement, env *object.Environment) object.Object {
	scope := object.NewEnclosedEnv(env)

	if loop.Init != nil {
		init := Eval(loop.Init, scope)
		if isError(init) {
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
			condition := Eval(loop.Condition, scope)
			if isError(condition) {
				return condition
			}
			if !isTruthy(condition) {
				break
			}
		}

		// Execute body
		body := Eval(loop.Body, scope)
		if isError(body) {
			return body
		}

		// Return if necessary
		rt := body.Type()
		if rt == object.RETURN_OBJ {
			return body
		}

		// Break if necessary, continue automatically
		if rt == object.LOOP_CONTROL_OBJ {
			if !body.(*object.LoopControl).Continue {
				break
			}
		}

		// Execute iterator
		if loop.Iter != nil {
			iter := Eval(loop.Iter, scope)
			if isError(iter) {
				return iter
			}
		}
	}
	return object.NULL
}

func evalCompareExpression(node *ast.CompareExpression, env *object.Environment) object.Object {
	left := Eval(node.Left, env)
	if isError(left) {
		return left
	}

	lBool, valid := convertToBoolean(left)
	if !valid {
		return object.NewError("Left side of conditional must be truthy or falsey")
	}

	// Short circuit if possible
	if node.Token.Type == token.LOr && lBool {
		return object.TRUE
	}
	if node.Token.Type == token.LAnd && !lBool {
		return object.FALSE
	}

	right := Eval(node.Right, env)
	if isError(right) {
		return right
	}

	rBool, valid := convertToBoolean(right)
	if !valid {
		return object.NewError("Right side of condition must be truthy or falsey")
	}

	return object.NativeBoolToBooleanObj(rBool)
}
