package eval

import (
	"github.com/nitrogen-lang/nitrogen/src/object"
)

func (i *Interpreter) evalPrefixExpression(op string, right object.Object) object.Object {
	switch op {
	case "!":
		return i.evalBangOpExpression(right)
	case "-":
		return i.evalMinusPreOpExpression(right)
	}

	return object.NewException("unknown operator: %s%s", op, right.Type())
}

func (i *Interpreter) evalBangOpExpression(right object.Object) object.Object {
	if right == object.FalseConst || right == object.NullConst {
		return object.TrueConst
	}

	if right.Type() == object.IntergerObj && right.(*object.Integer).Value == 0 {
		return object.TrueConst
	}

	return object.FalseConst
}

func (i *Interpreter) evalMinusPreOpExpression(right object.Object) object.Object {
	switch right.Type() {
	case object.IntergerObj:
		value := right.(*object.Integer).Value
		return &object.Integer{Value: -value}
	case object.FloatObj:
		value := right.(*object.Float).Value
		return &object.Float{Value: -value}
	}

	return object.NewException("unknown operator: -%s", right.Type())
}
