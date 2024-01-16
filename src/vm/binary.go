package vm

import (
	"math"

	"github.com/nitrogen-lang/nitrogen/src/object"
)

func (vm *VirtualMachine) evalBinaryExpression(op string, left, right object.Object) object.Object {
	switch {
	case left.Type() != right.Type():
		return object.NewException("type mismatch: %s %s %s", left.Type(), op, right.Type())
	case object.ObjectsAre(object.IntergerObj, left, right):
		return vm.evalIntegerBinaryExpression(op, left, right)
	case object.ObjectsAre(object.FloatObj, left, right):
		return vm.evalFloatBinaryExpression(op, left, right)
	case object.ObjectsAre(object.StringObj, left, right):
		return vm.evalStringBinaryExpression(op, left, right)
	case object.ObjectsAre(object.ArrayObj, left, right):
		return vm.evalArrayBinaryExpression(op, left, right)
	}

	return object.NewException("unknown operator: %s %s %s", left.Type(), op, right.Type())
}

func (vm *VirtualMachine) evalIntegerBinaryExpression(op string, left, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch op {
	case "+":
		return object.MakeIntObj(leftVal + rightVal)
	case "-":
		return object.MakeIntObj(leftVal - rightVal)
	case "*":
		return object.MakeIntObj(leftVal * rightVal)
	case "/":
		return object.MakeIntObj(leftVal / rightVal)
	case "%":
		return object.MakeIntObj(leftVal % rightVal)
	case "<<":
		if rightVal < 0 {
			return object.NewException("Shift value must be non-negative")
		}
		return object.MakeIntObj(leftVal << uint64(rightVal))
	case ">>":
		if rightVal < 0 {
			return object.NewException("Shift value must be non-negative")
		}
		return object.MakeIntObj(leftVal >> uint64(rightVal))
	case "&":
		return object.MakeIntObj(leftVal & rightVal)
	case "&^":
		return object.MakeIntObj(leftVal &^ rightVal)
	case "|":
		return object.MakeIntObj(leftVal | rightVal)
	case "^":
		return object.MakeIntObj(leftVal ^ rightVal)
	}

	return object.NewException("unknown operator: %s %s %s", left.Type(), op, right.Type())
}

func (vm *VirtualMachine) evalFloatBinaryExpression(op string, left, right object.Object) object.Object {
	leftVal := left.(*object.Float).Value
	rightVal := right.(*object.Float).Value

	switch op {
	case "+":
		return &object.Float{Value: leftVal + rightVal}
	case "-":
		return &object.Float{Value: leftVal - rightVal}
	case "*":
		return &object.Float{Value: leftVal * rightVal}
	case "/":
		return &object.Float{Value: leftVal / rightVal}
	case "%":
		return &object.Float{Value: math.Mod(leftVal, rightVal)}
	}

	return object.NewException("unknown operator: %s %s %s", left.Type(), op, right.Type())
}

func (vm *VirtualMachine) evalStringBinaryExpression(op string, left, right object.Object) object.Object {
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value

	if op == "+" {
		return &object.String{Value: append(leftVal, rightVal...)}
	}

	return object.NewException("unknown operator: %s %s %s", left.Type(), op, right.Type())
}

func (vm *VirtualMachine) evalArrayBinaryExpression(op string, left, right object.Object) object.Object {
	leftVal := left.(*object.Array)
	rightVal := right.(*object.Array)

	if op == "+" {
		leftLen := len(leftVal.Elements)
		rightLen := len(rightVal.Elements)
		newElements := make([]object.Object, leftLen+rightLen, leftLen+rightLen)
		copy(newElements, leftVal.Elements)
		copy(newElements[leftLen:], rightVal.Elements)
		return &object.Array{Elements: newElements}
	}

	return object.NewException("unknown operator: %s %s %s", left.Type(), op, right.Type())
}
