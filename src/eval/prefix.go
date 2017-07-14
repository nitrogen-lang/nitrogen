package eval

import (
	"github.com/nitrogen-lang/nitrogen/src/object"
)

func evalPrefixExpression(op string, right object.Object) object.Object {
	switch op {
	case "!":
		return evalBangOpExpression(right)
	case "-":
		return evalMinusPreOpExpression(right)
	}

	return object.NewError("unknown operator: %s%s", op, right.Type())
}

func evalBangOpExpression(right object.Object) object.Object {
	if right == object.FALSE || right == object.NULL {
		return object.TRUE
	}

	if right.Type() == object.INTEGER_OBJ && right.(*object.Integer).Value == 0 {
		return object.TRUE
	}

	return object.FALSE
}

func evalMinusPreOpExpression(right object.Object) object.Object {
	switch right.Type() {
	case object.INTEGER_OBJ:
		value := right.(*object.Integer).Value
		return &object.Integer{Value: -value}
	case object.FLOAT_OBJ:
		value := right.(*object.Float).Value
		return &object.Float{Value: -value}
	}

	return object.NewError("unknown operator: -%s", right.Type())
}
