package vm

import (
	"fmt"
	"io"
	"os"

	"github.com/nitrogen-lang/nitrogen/src/ast"
	"github.com/nitrogen-lang/nitrogen/src/compiler"
	"github.com/nitrogen-lang/nitrogen/src/object"
	"github.com/nitrogen-lang/nitrogen/src/vm/opcode"
)

type VirtualMachine struct {
	callStack    *frameStack
	currentFrame *frame
	returnValue  object.Object
}

func NewVM() *VirtualMachine {
	return &VirtualMachine{
		callStack: newFrameStack(),
	}
}

func (vm *VirtualMachine) Eval(node ast.Node, env *object.Environment) object.Object {
	return object.NullConst
}
func (vm *VirtualMachine) GetCurrentScriptPath() string { return vm.currentFrame.code.Filename }
func (vm *VirtualMachine) GetStdout() io.Writer         { return os.Stdout }
func (vm *VirtualMachine) GetStderr() io.Writer         { return os.Stderr }
func (vm *VirtualMachine) GetStdin() io.Reader          { return os.Stdout }

func (vm *VirtualMachine) Execute(code *compiler.CodeBlock) object.Object {
	f := &frame{
		code:  code,
		stack: object.NewStack(),
		pc:    0,
	}
	return vm.runFrame(f)
}

func (vm *VirtualMachine) runFrame(f *frame) object.Object {
	f.lastFrame = vm.currentFrame
	vm.callStack.Push(f)
	vm.currentFrame = f
	if f.locals == nil {
		f.locals = make([]object.Object, vm.currentFrame.code.LocalCount)
	}

	for {
		code := vm.fetchByte()
		fmt.Printf("Executing %s\n", opcode.Names[code])

		switch code {
		case opcode.Noop:
			break
		case opcode.BinaryAdd:
			r := vm.currentFrame.stack.Pop().(*object.Integer)
			l := vm.currentFrame.stack.Pop().(*object.Integer)
			vm.currentFrame.stack.Push(&object.Integer{Value: l.Value + r.Value})
		case opcode.BinarySub:
			r := vm.currentFrame.stack.Pop().(*object.Integer)
			l := vm.currentFrame.stack.Pop().(*object.Integer)
			vm.currentFrame.stack.Push(&object.Integer{Value: l.Value - r.Value})
		case opcode.BinaryMul:
			r := vm.currentFrame.stack.Pop().(*object.Integer)
			l := vm.currentFrame.stack.Pop().(*object.Integer)
			vm.currentFrame.stack.Push(&object.Integer{Value: l.Value * r.Value})
		case opcode.BinaryDivide:
			r := vm.currentFrame.stack.Pop().(*object.Integer)
			l := vm.currentFrame.stack.Pop().(*object.Integer)
			vm.currentFrame.stack.Push(&object.Integer{Value: l.Value / r.Value})
		case opcode.BinaryMod:
			r := vm.currentFrame.stack.Pop().(*object.Integer)
			l := vm.currentFrame.stack.Pop().(*object.Integer)
			vm.currentFrame.stack.Push(&object.Integer{Value: l.Value % r.Value})
		case opcode.BinaryShiftL:
			r := vm.currentFrame.stack.Pop().(*object.Integer)
			l := vm.currentFrame.stack.Pop().(*object.Integer)
			vm.currentFrame.stack.Push(&object.Integer{Value: l.Value << uint64(r.Value)})
		case opcode.BinaryShiftR:
			r := vm.currentFrame.stack.Pop().(*object.Integer)
			l := vm.currentFrame.stack.Pop().(*object.Integer)
			vm.currentFrame.stack.Push(&object.Integer{Value: l.Value >> uint64(r.Value)})
		case opcode.BinaryAnd:
			r := vm.currentFrame.stack.Pop().(*object.Integer)
			l := vm.currentFrame.stack.Pop().(*object.Integer)
			vm.currentFrame.stack.Push(&object.Integer{Value: l.Value & r.Value})
		case opcode.BinaryOr:
			r := vm.currentFrame.stack.Pop().(*object.Integer)
			l := vm.currentFrame.stack.Pop().(*object.Integer)
			vm.currentFrame.stack.Push(&object.Integer{Value: l.Value | r.Value})
		case opcode.BinaryNot:
			r := vm.currentFrame.stack.Pop().(*object.Integer)
			l := vm.currentFrame.stack.Pop().(*object.Integer)
			vm.currentFrame.stack.Push(&object.Integer{Value: l.Value ^ r.Value})
		case opcode.BinaryAndNot:
			r := vm.currentFrame.stack.Pop().(*object.Integer)
			l := vm.currentFrame.stack.Pop().(*object.Integer)
			vm.currentFrame.stack.Push(&object.Integer{Value: l.Value &^ r.Value})
		case opcode.UnaryNeg:
			l := vm.currentFrame.stack.Pop().(*object.Integer)
			vm.currentFrame.stack.Push(&object.Integer{Value: -l.Value})
		case opcode.UnaryNot:
			l := vm.currentFrame.stack.Pop().(*object.Boolean)
			if l.Value {
				vm.currentFrame.stack.Push(object.FalseConst)
			} else {
				vm.currentFrame.stack.Push(object.TrueConst)
			}
		case opcode.LoadConst:
			vm.currentFrame.stack.Push(vm.currentFrame.code.Constants[vm.getUint16()])
		case opcode.Return:
			vm.returnValue = vm.currentFrame.stack.Pop()
			vm.currentFrame = vm.currentFrame.lastFrame
			vm.callStack.Pop()
			if vm.currentFrame == nil {
				return vm.returnValue
			}
		case opcode.Pop:
			vm.currentFrame.stack.Pop()
		case opcode.LoadFast:
			vm.currentFrame.stack.Push(vm.currentFrame.locals[vm.getUint16()])
		case opcode.StoreFast:
			vm.currentFrame.locals[vm.getUint16()] = vm.currentFrame.stack.Pop()
		case opcode.LoadGlobal:
			name := vm.currentFrame.code.Names[vm.getUint16()]
			fn := getBuiltin(name)
			if fn == nil {
				fmt.Printf("Global %s doesn't exist\n", name)
				return vm.returnValue
			}
			vm.currentFrame.stack.Push(fn)
		case opcode.Call:
			numargs := vm.getUint16()
			args := make([]object.Object, numargs)
			fn := vm.currentFrame.stack.Pop()
			for i := uint16(0); i < numargs; i++ {
				args[i] = vm.currentFrame.stack.Pop()
			}

			vm.currentFrame.stack.Push(fn.(*vmBuiltin).fn(vm, args...))
		case opcode.Compare:
			r := vm.currentFrame.stack.Pop()
			l := vm.currentFrame.stack.Pop()
			op := vm.fetchByte()
			if op >= opcode.MaxCmpCodes {
				panic(fmt.Sprintf("Invalid comparison operator %x", op))
			}
			vm.currentFrame.stack.Push(vm.compareObjects(l, r, op))
		case opcode.MakeArray:
			l := vm.getUint16()
			array := &object.Array{
				Elements: make([]object.Object, l),
			}

			for i := l; i > 0; i-- {
				array.Elements[i-1] = vm.currentFrame.stack.Pop()
			}
			vm.currentFrame.stack.Push(array)
		case opcode.MakeMap:
			l := vm.getUint16()
			hash := &object.Hash{
				Pairs: make(map[object.HashKey]object.HashPair, l),
			}

			for i := l; i > 0; i-- {
				key := vm.currentFrame.stack.Pop()
				val := vm.currentFrame.stack.Pop()
				hashKey, ok := key.(object.Hashable)
				if !ok {
					panic("Map key not valid")
				}
				hash.Pairs[hashKey.HashKey()] = object.HashPair{
					Key:   key,
					Value: val,
				}
			}
			vm.currentFrame.stack.Push(hash)
		case opcode.PopJumpIfFalse:
			target := vm.getUint16()
			tos := vm.currentFrame.stack.Pop()
			if tos == object.FalseConst {
				vm.currentFrame.pc = int(target)
			}
		case opcode.JumpAbsolute:
			vm.currentFrame.pc = int(vm.getUint16())
		case opcode.PopJumpIfTrue:
			target := vm.getUint16()
			tos := vm.currentFrame.stack.Pop()
			if tos == object.TrueConst {
				vm.currentFrame.pc = int(target)
			}
		case opcode.JumpForward:
			jump := vm.getUint16()
			vm.currentFrame.pc += int(jump)
		case opcode.JumpIfTrueOrPop:
			target := vm.getUint16()
			tos := vm.currentFrame.stack.GetFront()
			if tos == object.TrueConst {
				vm.currentFrame.pc = int(target)
			} else {
				vm.currentFrame.stack.Pop()
			}
		case opcode.JumpIfFalseOrPop:
			target := vm.getUint16()
			tos := vm.currentFrame.stack.GetFront()
			if tos == object.FalseConst {
				vm.currentFrame.pc = int(target)
			} else {
				vm.currentFrame.stack.Pop()
			}
		default:
			panic(fmt.Sprintf("Opcode %s is not supported", opcode.Names[code]))
		}
	}
}

func (vm *VirtualMachine) fetchByte() byte {
	b := vm.currentFrame.code.Code[vm.currentFrame.pc]
	vm.currentFrame.pc++
	return b
}

func (vm *VirtualMachine) getUint16() uint16 {
	return (uint16(vm.fetchByte()) << 8) + uint16(vm.fetchByte())
}

func (vm *VirtualMachine) compareObjects(left, right object.Object, op byte) object.Object {
	switch {
	case left.Type() != right.Type():
		return object.NewException("type mismatch: %s %s %s", left.Type(), opcode.CmpOps[op], right.Type())
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
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value

	switch op {
	case opcode.CmpEq:
		return object.NativeBoolToBooleanObj(leftVal == rightVal)
	case opcode.CmpNotEq:
		return object.NativeBoolToBooleanObj(leftVal != rightVal)
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
