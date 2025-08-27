package vm

import (
	"github.com/nitrogen-lang/nitrogen/src/elemental/object"
	"github.com/nitrogen-lang/nitrogen/src/elemental/vm/opcode"
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
			var m object.ClassMethod

			m, exists := nativeMethods[class.Name+"."+method.Name]
			if !exists {
				m, exists = class.Methods[method.Name]
				if !exists {
					return object.FalseConst
				}
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
		if op == opcode.CmpNotEq {
			return object.TrueConst
		}
		return object.FalseConst
	case object.ObjectsAre(object.IntergerObj, left, right):
		return vm.evalIntegerInfixExpression(op, left, right)
	case object.ObjectsAre(object.FloatObj, left, right):
		return vm.evalFloatInfixExpression(op, left, right)
	case object.ObjectsAre(object.StringObj, left, right):
		return vm.evalStringInfixExpression(op, left, right)
	case object.ObjectsAre(object.ByteStringObj, left, right):
		return vm.evalByteStringInfixExpression(op, left, right)
	case object.ObjectsAre(object.BooleanObj, left, right):
		return vm.evalBoolInfixExpression(op, left, right)
	case object.ObjectsAre(object.NullObj, left, right):
		return vm.evalNullInfixExpression(op)
	case object.ObjectsAre(object.ArrayObj, left, right):
		return vm.evalArrayInfixExpression(op, left, right)
	}

	return object.NewException("comparison is not implemented for type %s", left.Type())
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

	return object.NewException("unknown operator: %s %s %s", left.Type(), opcode.CmpOps[op], right.Type())
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

	return object.NewException("unknown operator: %s %s %s", left.Type(), opcode.CmpOps[op], right.Type())
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

	return object.NewException("unknown operator: %s %s %s", left.Type(), opcode.CmpOps[op], right.Type())
}

func (vm *VirtualMachine) evalByteStringInfixExpression(op byte, left, right object.Object) object.Object {
	leftVal := left.(*object.ByteString).String()
	rightVal := right.(*object.ByteString).String()

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

	return object.NewException("unknown operator: %s %s %s", left.Type(), opcode.CmpOps[op], right.Type())
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

	return object.NewException("unknown operator: %s %s %s", left.Type(), opcode.CmpOps[op], right.Type())
}

func (vm *VirtualMachine) evalNullInfixExpression(op byte) object.Object {
	switch op {
	case opcode.CmpEq:
		return object.TrueConst
	case opcode.CmpNotEq:
		return object.FalseConst
	}

	return object.NewException("unknown operator: nil %c nil", op)
}

func (vm *VirtualMachine) evalArrayInfixExpression(op byte, left, right object.Object) object.Object {
	leftVal := left.(*object.Array).Elements
	rightVal := right.(*object.Array).Elements

	switch op {
	case opcode.CmpEq:
		if len(leftVal) != len(rightVal) {
			return object.FalseConst
		}

		for i, l := range leftVal {
			r := rightVal[i]
			res := vm.compareObjects(l, r, opcode.CmpEq)
			if res != object.TrueConst {
				return res
			}
		}
		return object.TrueConst
	case opcode.CmpNotEq:
		if len(leftVal) != len(rightVal) {
			return object.TrueConst
		}

		for i, l := range leftVal {
			r := rightVal[i]
			res := vm.compareObjects(l, r, opcode.CmpNotEq)
			if res != object.FalseConst {
				return res
			}
		}
		return object.FalseConst
	}

	return object.NewException("unknown operator: array %c array", op)
}
