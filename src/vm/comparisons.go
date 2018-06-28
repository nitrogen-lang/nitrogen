package vm

import (
	"github.com/nitrogen-lang/nitrogen/src/object"
	"github.com/nitrogen-lang/nitrogen/src/vm/opcode"
)

func (vm *VirtualMachine) compareObjects(left, right object.Object, op byte) object.Object {
	switch {
	case left.Type() != right.Type():
		panic(object.NewException("type mismatch: %s %s %s", left.Type(), opcode.CmpOps[op], right.Type()))
	case object.ObjectsAre(object.IntergerObj, left, right):
		return vm.evalIntegerInfixExpression(op, left, right)
	case object.ObjectsAre(object.FloatObj, left, right):
		return vm.evalFloatInfixExpression(op, left, right)
	case object.ObjectsAre(object.StringObj, left, right):
		return vm.evalStringInfixExpression(op, left, right)
	case object.ObjectsAre(object.BooleanObj, left, right):
		return vm.evalBoolInfixExpression(op, left, right)
	}

	return object.NullConst
}

func (vm *VirtualMachine) evalIntegerInfixExpression(op byte, left, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch op {
	case opcode.CmpLT:
		return object.NativeBoolToBooleanObj(leftVal < rightVal)
	case opcode.CmpGT:
		return object.NativeBoolToBooleanObj(leftVal > rightVal)
	case opcode.CmpEq:
		return object.NativeBoolToBooleanObj(leftVal == rightVal)
	case opcode.CmpNotEq:
		return object.NativeBoolToBooleanObj(leftVal != rightVal)
	case opcode.CmpLTEq:
		return object.NativeBoolToBooleanObj(leftVal <= rightVal)
	case opcode.CmpGTEq:
		return object.NativeBoolToBooleanObj(leftVal >= rightVal)
	}

	return object.NewException("unknown operator: %s %s %s", left.Type(), op, right.Type())
}

func (vm *VirtualMachine) evalFloatInfixExpression(op byte, left, right object.Object) object.Object {
	leftVal := left.(*object.Float).Value
	rightVal := right.(*object.Float).Value

	switch op {
	case opcode.CmpLT:
		return object.NativeBoolToBooleanObj(leftVal < rightVal)
	case opcode.CmpGT:
		return object.NativeBoolToBooleanObj(leftVal > rightVal)
	case opcode.CmpEq:
		return object.NativeBoolToBooleanObj(leftVal == rightVal)
	case opcode.CmpNotEq:
		return object.NativeBoolToBooleanObj(leftVal != rightVal)
	case opcode.CmpLTEq:
		return object.NativeBoolToBooleanObj(leftVal <= rightVal)
	case opcode.CmpGTEq:
		return object.NativeBoolToBooleanObj(leftVal >= rightVal)
	}

	return object.NewException("unknown operator: %s %s %s", left.Type(), op, right.Type())
}

func (vm *VirtualMachine) evalStringInfixExpression(op byte, left, right object.Object) object.Object {
	leftVal := left.(*object.String).String()
	rightVal := right.(*object.String).String()

	switch op {
	case opcode.CmpEq:
		return object.NativeBoolToBooleanObj(leftVal == rightVal)
	case opcode.CmpNotEq:
		return object.NativeBoolToBooleanObj(leftVal != rightVal)
	case opcode.CmpGT:
		return object.NativeBoolToBooleanObj(leftVal > rightVal)
	case opcode.CmpGTEq:
		return object.NativeBoolToBooleanObj(leftVal >= rightVal)
	case opcode.CmpLT:
		return object.NativeBoolToBooleanObj(leftVal < rightVal)
	case opcode.CmpLTEq:
		return object.NativeBoolToBooleanObj(leftVal <= rightVal)
	}

	return object.NewException("unknown operator: %s %s %s", left.Type(), op, right.Type())
}

func (vm *VirtualMachine) evalBoolInfixExpression(op byte, left, right object.Object) object.Object {
	leftVal := left.(*object.Boolean).Value
	rightVal := right.(*object.Boolean).Value

	switch op {
	case opcode.CmpEq:
		return object.NativeBoolToBooleanObj(leftVal == rightVal)
	case opcode.CmpNotEq:
		return object.NativeBoolToBooleanObj(leftVal != rightVal)
	}

	return object.NewException("unknown operator: %s %s %s", left.Type(), op, right.Type())
}
