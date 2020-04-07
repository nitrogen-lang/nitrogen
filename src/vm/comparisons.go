package vm

import (
	"github.com/nitrogen-lang/nitrogen/src/object"
	"github.com/nitrogen-lang/nitrogen/src/vm/opcode"
)

func (vm *VirtualMachine) evalImplementsExpression(left, right object.Object) object.Object {
	var iface *object.Interface

	if inter, ok := right.(*object.Interface); ok {
		iface = inter
	} else {
		return object.NewException("implements must have an interface on the right side")
	}

	switch node := left.(type) {
	case *object.Interface:
		for _, method := range iface.Methods {
			m, exists := node.Methods[method.Name]
			if !exists || len(m.Parameters) != len(method.Parameters) {
				return object.FalseConst
			}
		}
	case *object.Class:
		for _, method := range iface.Methods {
			m, exists := node.Methods[method.Name]
			if !exists {
				return object.FalseConst
			}

			fn, ok := m.(*object.Function)
			if !ok || len(fn.Parameters) != len(method.Parameters) {
				return object.FalseConst
			}
		}
	case *object.Instance:
		class := node.Class
		for _, method := range iface.Methods {
			m, exists := class.Methods[method.Name]
			if !exists {
				return object.FalseConst
			}

			fn, ok := m.(*object.Function)
			if !ok || len(fn.Parameters) != len(method.Parameters) {
				return object.FalseConst
			}
		}
	case *VMClass:
		for _, method := range iface.Methods {
			m, exists := node.Methods[method.Name]
			if !exists {
				return object.FalseConst
			}

			switch m := m.(type) {
			case *object.Function:
				if len(m.Parameters) != len(method.Parameters) {
					return object.FalseConst
				}
			case *BuiltinMethod:
				if m.NumOfParams != len(method.Parameters) {
					return object.FalseConst
				}
			case *VMFunction:
				if len(m.Parameters) != len(method.Parameters) {
					return object.FalseConst
				}
			default:
				return object.FalseConst
			}
		}
	case *VMInstance:
		class := node.Class
		for _, method := range iface.Methods {
			m, exists := class.Methods[method.Name]
			if !exists {
				return object.FalseConst
			}

			switch m := m.(type) {
			case *object.Function:
				if len(m.Parameters) != len(method.Parameters) {
					return object.FalseConst
				}
			case *BuiltinMethod:
				if m.NumOfParams != len(method.Parameters) {
					return object.FalseConst
				}
			case *VMFunction:
				if len(m.Parameters) != len(method.Parameters) {
					return object.FalseConst
				}
			default:
				return object.FalseConst
			}
		}
	case *BuiltinClass:
		for _, method := range iface.Methods {
			m, exists := node.Methods[method.Name]
			if !exists {
				return object.FalseConst
			}

			switch m := m.(type) {
			case *object.Function:
				if len(m.Parameters) != len(method.Parameters) {
					return object.FalseConst
				}
			case *BuiltinMethod:
				if m.NumOfParams != len(method.Parameters) {
					return object.FalseConst
				}
			case *VMFunction:
				if len(m.Parameters) != len(method.Parameters) {
					return object.FalseConst
				}
			default:
				return object.FalseConst
			}
		}
	}

	return object.TrueConst
}

func (vm *VirtualMachine) compareObjects(left, right object.Object, op byte) object.Object {
	switch {
	case left.Type() != right.Type():
		return object.FalseConst
	case object.ObjectsAre(object.IntergerObj, left, right):
		return vm.evalIntegerInfixExpression(op, left, right)
	case object.ObjectsAre(object.FloatObj, left, right):
		return vm.evalFloatInfixExpression(op, left, right)
	case object.ObjectsAre(object.StringObj, left, right):
		return vm.evalStringInfixExpression(op, left, right)
	case object.ObjectsAre(object.BooleanObj, left, right):
		return vm.evalBoolInfixExpression(op, left, right)
	case object.ObjectsAre(object.NullObj, left, right):
		return vm.evalNullInfixExpression(op)
	}

	return object.NewException("comparison not implemented for type %s", left.Type())
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

func (vm *VirtualMachine) evalNullInfixExpression(op byte) object.Object {
	switch op {
	case opcode.CmpEq:
		return object.TrueConst
	case opcode.CmpNotEq:
		return object.FalseConst
	}

	return object.NewException("unknown operator: nil %s nil", op)
}
