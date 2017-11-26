package vm

import (
	"fmt"
	"io"
	"math"
	"os"

	"github.com/nitrogen-lang/nitrogen/src/ast"
	"github.com/nitrogen-lang/nitrogen/src/compiler"
	"github.com/nitrogen-lang/nitrogen/src/object"
	"github.com/nitrogen-lang/nitrogen/src/vm/opcode"
)

type Settings struct {
	Debug bool
}

type VirtualMachine struct {
	callStack    *frameStack
	currentFrame *Frame
	returnValue  object.Object
	settings     *Settings
}

func NewVM(settings *Settings) *VirtualMachine {
	if settings == nil {
		settings = &Settings{}
	}
	return &VirtualMachine{
		callStack: newFrameStack(),
		settings:  settings,
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
	f := &Frame{
		code:       code,
		stack:      make([]object.Object, code.MaxStackSize),
		blockStack: make([]block, code.MaxBlockSize),
		Env:        object.NewEnvironment(),
	}
	return vm.runFrame(f)
}

func (vm *VirtualMachine) CurrentFrame() *Frame {
	return vm.currentFrame
}

func (vm *VirtualMachine) runFrame(f *Frame) object.Object {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
			fmt.Println("Stack Trace:")
			frame := vm.currentFrame
			for frame != nil {
				fmt.Printf("\t%s: %s\n", frame.code.Filename, frame.code.Name)
				frame = frame.lastFrame
			}
		}
	}()
	f.lastFrame = vm.currentFrame
	vm.callStack.Push(f)
	vm.currentFrame = f

	for {
		if vm.currentFrame.pc >= len(vm.currentFrame.code.Code) {
			panic(fmt.Sprintf("Program counter %d outside bounds of bytecode %d", vm.currentFrame.pc, len(vm.currentFrame.code.Code)))
		}
		code := vm.fetchByte()
		if vm.settings.Debug {
			fmt.Printf("Executing %s\n", opcode.Names[code])
		}

		switch code {
		case opcode.Noop:
			break
		case opcode.BinaryAdd:
			r := vm.currentFrame.popStack()
			l := vm.currentFrame.popStack()
			res := vm.evalBinaryExpression("+", l, r)
			if object.ObjectIs(res, object.ExceptionObj) {
				fmt.Printf("Exception: %s\n", res.Inspect())
			}
			vm.currentFrame.pushStack(res)
		case opcode.BinarySub:
			r := vm.currentFrame.popStack()
			l := vm.currentFrame.popStack()
			res := vm.evalBinaryExpression("-", l, r)
			if object.ObjectIs(res, object.ExceptionObj) {
				fmt.Printf("Exception: %s\n", res.Inspect())
			}
			vm.currentFrame.pushStack(res)
		case opcode.BinaryMul:
			r := vm.currentFrame.popStack()
			l := vm.currentFrame.popStack()
			res := vm.evalBinaryExpression("*", l, r)
			if object.ObjectIs(res, object.ExceptionObj) {
				fmt.Printf("Exception: %s\n", res.Inspect())
			}
			vm.currentFrame.pushStack(res)
		case opcode.BinaryDivide:
			r := vm.currentFrame.popStack()
			l := vm.currentFrame.popStack()
			res := vm.evalBinaryExpression("/", l, r)
			if object.ObjectIs(res, object.ExceptionObj) {
				fmt.Printf("Exception: %s\n", res.Inspect())
			}
			vm.currentFrame.pushStack(res)
		case opcode.BinaryMod:
			r := vm.currentFrame.popStack()
			l := vm.currentFrame.popStack()
			res := vm.evalBinaryExpression("%", l, r)
			if object.ObjectIs(res, object.ExceptionObj) {
				fmt.Printf("Exception: %s\n", res.Inspect())
			}
			vm.currentFrame.pushStack(res)
		case opcode.BinaryShiftL:
			r := vm.currentFrame.popStack()
			l := vm.currentFrame.popStack()
			res := vm.evalBinaryExpression("<<", l, r)
			if object.ObjectIs(res, object.ExceptionObj) {
				fmt.Printf("Exception: %s\n", res.Inspect())
			}
			vm.currentFrame.pushStack(res)
		case opcode.BinaryShiftR:
			r := vm.currentFrame.popStack()
			l := vm.currentFrame.popStack()
			res := vm.evalBinaryExpression(">>", l, r)
			if object.ObjectIs(res, object.ExceptionObj) {
				fmt.Printf("Exception: %s\n", res.Inspect())
			}
			vm.currentFrame.pushStack(res)
		case opcode.BinaryAnd:
			r := vm.currentFrame.popStack()
			l := vm.currentFrame.popStack()
			res := vm.evalBinaryExpression("&", l, r)
			if object.ObjectIs(res, object.ExceptionObj) {
				fmt.Printf("Exception: %s\n", res.Inspect())
			}
			vm.currentFrame.pushStack(res)
		case opcode.BinaryOr:
			r := vm.currentFrame.popStack()
			l := vm.currentFrame.popStack()
			res := vm.evalBinaryExpression("|", l, r)
			if object.ObjectIs(res, object.ExceptionObj) {
				fmt.Printf("Exception: %s\n", res.Inspect())
			}
			vm.currentFrame.pushStack(res)
		case opcode.BinaryNot:
			r := vm.currentFrame.popStack()
			l := vm.currentFrame.popStack()
			res := vm.evalBinaryExpression("^", l, r)
			if object.ObjectIs(res, object.ExceptionObj) {
				fmt.Printf("Exception: %s\n", res.Inspect())
			}
			vm.currentFrame.pushStack(res)
		case opcode.BinaryAndNot:
			r := vm.currentFrame.popStack()
			l := vm.currentFrame.popStack()
			res := vm.evalBinaryExpression("&^", l, r)
			if object.ObjectIs(res, object.ExceptionObj) {
				fmt.Printf("Exception: %s\n", res.Inspect())
			}
			vm.currentFrame.pushStack(res)
		case opcode.UnaryNeg:
			l := vm.currentFrame.popStack().(*object.Integer)
			vm.currentFrame.pushStack(&object.Integer{Value: -l.Value})
		case opcode.UnaryNot:
			l := vm.currentFrame.popStack().(*object.Boolean)
			if l.Value {
				vm.currentFrame.pushStack(object.FalseConst)
			} else {
				vm.currentFrame.pushStack(object.TrueConst)
			}
		case opcode.LoadConst:
			vm.currentFrame.pushStack(vm.currentFrame.code.Constants[vm.getUint16()])
		case opcode.StoreConst:
			// Ensure constant isn't redefined
			name := vm.currentFrame.code.Locals[vm.getUint16()]
			if vm.currentFrame.Env.IsConst(name) {
				fmt.Printf("Redefined constant %s\n", name)
				return vm.returnValue
			}
			if _, err := vm.currentFrame.Env.CreateConst(name, vm.currentFrame.popStack()); err != nil {
				fmt.Println(err)
			}
		case opcode.Return:
			vm.returnValue = vm.currentFrame.popStack()
			vm.currentFrame = vm.currentFrame.lastFrame
			vm.callStack.Pop()
			if vm.currentFrame == nil {
				return vm.returnValue
			}
			vm.currentFrame.pushStack(vm.returnValue)
		case opcode.Pop:
			vm.currentFrame.popStack()
		case opcode.LoadFast:
			name := vm.currentFrame.code.Locals[vm.getUint16()]
			if val, ok := vm.currentFrame.Env.Get(name); ok {
				vm.currentFrame.pushStack(val)
				break
			}

			fmt.Printf("Unknown variable/constant %s\n", name)
			return vm.returnValue
		case opcode.StoreFast:
			// Ensure constant isn't redefined
			name := vm.currentFrame.code.Locals[vm.getUint16()]
			if vm.currentFrame.Env.IsConst(name) {
				fmt.Printf("Redefined constant %s\n", name)
				return vm.returnValue
			}
			if _, exists := vm.currentFrame.Env.Get(name); !exists {
				fmt.Printf("Variable %s undefined\n", name)
				return vm.returnValue
			}
			vm.currentFrame.Env.Set(name, vm.currentFrame.popStack())
		case opcode.Define:
			// Ensure constant isn't redefined
			name := vm.currentFrame.code.Locals[vm.getUint16()]
			if vm.currentFrame.Env.IsConst(name) {
				fmt.Printf("Redefined constant %s\n", name)
				return vm.returnValue
			}
			if _, exists := vm.currentFrame.Env.GetLocal(name); exists {
				fmt.Printf("Variable %s already defined\n", name)
				return vm.returnValue
			}
			vm.currentFrame.Env.Create(name, vm.currentFrame.popStack())
		case opcode.LoadGlobal:
			name := vm.currentFrame.code.Names[vm.getUint16()]
			if val, ok := vm.currentFrame.Env.Get(name); ok {
				vm.currentFrame.pushStack(val)
				break
			}
			if fn := getBuiltin(name); fn != nil {
				vm.currentFrame.pushStack(fn)
				break
			}

			fmt.Printf("Global %s doesn't exist\n", name)
			return vm.returnValue
		case opcode.StoreGlobal:
			// Ensure constant isn't redefined
			name := vm.currentFrame.code.Names[vm.getUint16()]
			if vm.currentFrame.Env.IsConst(name) {
				fmt.Printf("Redefined constant %s\n", name)
				return vm.returnValue
			}
			if _, exists := vm.currentFrame.Env.Get(name); !exists {
				fmt.Printf("Global variable %s not defined\n", name)
				return vm.returnValue
			}
			vm.currentFrame.Env.Set(name, vm.currentFrame.popStack())
		case opcode.LoadIndex:
			left := vm.currentFrame.popStack()
			index := vm.currentFrame.popStack()
			res := vm.evalIndexExpression(left, index)
			if object.ObjectIs(res, object.ExceptionObj) {
				fmt.Printf("Exception %s\n", res.Inspect())
				return vm.returnValue
			}
			vm.currentFrame.pushStack(res)
		case opcode.StoreIndex:
			left := vm.currentFrame.popStack()
			index := vm.currentFrame.popStack()
			value := vm.currentFrame.popStack()
			res := vm.assignIndexedValue(left, index, value)
			if object.ObjectIs(res, object.ExceptionObj) {
				fmt.Printf("Exception %s\n", res.Inspect())
				return vm.returnValue
			}
		case opcode.Call:
			numargs := vm.getUint16()
			args := make([]object.Object, numargs)
			fn := vm.currentFrame.popStack()
			for i := uint16(0); i < numargs; i++ {
				args[i] = vm.currentFrame.popStack()
			}

			switch fn := fn.(type) {
			case *vmBuiltin:
				vm.currentFrame.pushStack(fn.fn(vm, args...))
			case *VMFunction:
				newFrame := &Frame{
					code:       fn.Body,
					stack:      make([]object.Object, fn.Body.MaxStackSize),
					blockStack: make([]block, fn.Body.MaxBlockSize),
					Env:        object.NewEnclosedEnv(fn.Env),
					lastFrame:  vm.currentFrame,
				}

				for i, arg := range args {
					newFrame.Env.SetForce(fn.Parameters[i], arg, false)
				}

				vm.currentFrame = newFrame
				vm.callStack.Push(newFrame)
			}
		case opcode.Compare:
			r := vm.currentFrame.popStack()
			l := vm.currentFrame.popStack()
			op := vm.fetchByte()
			if op >= opcode.MaxCmpCodes {
				panic(fmt.Sprintf("Invalid comparison operator %x", op))
			}
			vm.currentFrame.pushStack(vm.compareObjects(l, r, op))
		case opcode.MakeFunction:
			fnName := vm.currentFrame.popStack().(*object.String)
			params := vm.currentFrame.popStack().(*object.Array)
			codeBlock := vm.currentFrame.popStack().(*compiler.CodeBlock)

			fn := &VMFunction{
				Name:       fnName.Value,
				Parameters: make([]string, len(params.Elements)),
				Body:       codeBlock,
				Env:        object.NewEnclosedEnv(vm.currentFrame.Env),
			}

			for i, p := range params.Elements {
				fn.Parameters[i] = p.(*object.String).Value
			}
			vm.currentFrame.pushStack(fn)
		case opcode.MakeArray:
			l := vm.getUint16()
			array := &object.Array{
				Elements: make([]object.Object, l),
			}

			for i := l; i > 0; i-- {
				array.Elements[i-1] = vm.currentFrame.popStack()
			}
			vm.currentFrame.pushStack(array)
		case opcode.MakeMap:
			l := vm.getUint16()
			hash := &object.Hash{
				Pairs: make(map[object.HashKey]object.HashPair, l),
			}

			for i := l; i > 0; i-- {
				key := vm.currentFrame.popStack()
				val := vm.currentFrame.popStack()
				hashKey, ok := key.(object.Hashable)
				if !ok {
					panic("Map key not valid")
				}
				hash.Pairs[hashKey.HashKey()] = object.HashPair{
					Key:   key,
					Value: val,
				}
			}
			vm.currentFrame.pushStack(hash)
		case opcode.PopJumpIfFalse:
			target := vm.getUint16()
			tos := vm.currentFrame.popStack()
			if tos == object.FalseConst {
				vm.currentFrame.pc = int(target)
			}
		case opcode.JumpAbsolute:
			vm.currentFrame.pc = int(vm.getUint16())
		case opcode.PopJumpIfTrue:
			target := vm.getUint16()
			tos := vm.currentFrame.popStack()
			if tos == object.TrueConst {
				vm.currentFrame.pc = int(target)
			}
		case opcode.JumpForward:
			jump := vm.getUint16()
			vm.currentFrame.pc += int(jump)
		case opcode.JumpIfTrueOrPop:
			target := vm.getUint16()
			tos := vm.currentFrame.getFrontStack()
			if tos == object.TrueConst {
				vm.currentFrame.pc = int(target)
			} else {
				vm.currentFrame.popStack()
			}
		case opcode.JumpIfFalseOrPop:
			target := vm.getUint16()
			tos := vm.currentFrame.getFrontStack()
			if tos == object.FalseConst {
				vm.currentFrame.pc = int(target)
			} else {
				vm.currentFrame.popStack()
			}

		case opcode.PrepareBlock:
			vm.currentFrame.Env = object.NewEnclosedEnv(vm.currentFrame.Env)
		case opcode.EndBlock:
			vm.currentFrame.Env = vm.currentFrame.Env.Parent().Parent()
			vm.currentFrame.popBlock()
		case opcode.StartLoop:
			loopEnd := vm.getUint16()
			iter := vm.getUint16()
			vm.currentFrame.pushBlock(vm.currentFrame.pc, int(iter), int(loopEnd))
			vm.currentFrame.Env = object.NewEnclosedEnv(vm.currentFrame.Env)
		case opcode.Continue:
			_, vm.currentFrame.pc, _ = vm.currentFrame.getCurrentBlock()
			vm.currentFrame.Env = object.NewEnclosedEnv(vm.currentFrame.Env.Parent())
		case opcode.NextIter:
			vm.currentFrame.pc, _, _ = vm.currentFrame.getCurrentBlock()
		case opcode.Break:
			_, _, vm.currentFrame.pc = vm.currentFrame.getCurrentBlock()
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

func (vm *VirtualMachine) evalIndexExpression(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ArrayObj && index.Type() == object.IntergerObj:
		return vm.evalArrayIndexExpression(left, index)
	case left.Type() == object.HashObj:
		return vm.evalHashIndexExpression(left, index)
	case left.Type() == object.StringObj && index.Type() == object.IntergerObj:
		return vm.evalStringIndexExpression(left, index)
	}
	return object.NewException("Index operator not allowed: %s", left.Type())
}

func (vm *VirtualMachine) evalArrayIndexExpression(array, index object.Object) object.Object {
	arrObj := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrObj.Elements))

	if idx > max-1 { // Check upper bound
		return object.NullConst
	}

	if idx < 0 { // Check lower bound
		// Convert a negative index to positive
		idx = max + idx

		if idx < 0 { // Check lower bound again
			return object.NullConst
		}
	}

	return arrObj.Elements[idx]
}

func (vm *VirtualMachine) evalHashIndexExpression(hash, index object.Object) object.Object {
	hashObj := hash.(*object.Hash)

	key, ok := index.(object.Hashable)
	if !ok {
		return object.NewException("Invalid map key: %s", index.Type())
	}

	pair, ok := hashObj.Pairs[key.HashKey()]
	if !ok {
		return object.NullConst
	}

	return pair.Value
}

func (vm *VirtualMachine) evalStringIndexExpression(array, index object.Object) object.Object {
	strObj := array.(*object.String)
	idx := index.(*object.Integer).Value
	max := int64(len(strObj.Value))

	if idx > max-1 { // Check upper bound
		return object.NullConst
	}

	if idx < 0 { // Check lower bound
		// Convert a negative index to positive
		idx = max + idx

		if idx < 0 { // Check lower bound again
			return object.NullConst
		}
	}

	return &object.String{Value: string(strObj.Value[idx])}
}

func (vm *VirtualMachine) assignIndexedValue(
	indexed object.Object,
	index object.Object,
	val object.Object) object.Object {
	switch indexed.Type() {
	case object.ArrayObj:
		return vm.assignArrayIndex(indexed.(*object.Array), index, val)
	case object.HashObj:
		return vm.assignHashMapIndex(indexed.(*object.Hash), index, val)
	}
	return object.NullConst
}

func (vm *VirtualMachine) assignArrayIndex(
	array *object.Array,
	index object.Object,
	val object.Object) object.Object {

	in, ok := index.(*object.Integer)
	if !ok {
		return object.NewException("Invalid array index type %s", index.(object.Object).Type())
	}

	if in.Value < 0 || in.Value > int64(len(array.Elements)-1) {
		return object.NewException("Index out of bounds: %s", index.Inspect())
	}

	array.Elements[in.Value] = val
	return object.NullConst
}

func (vm *VirtualMachine) assignHashMapIndex(
	hashmap *object.Hash,
	index object.Object,
	val object.Object) object.Object {

	hashable, ok := index.(object.Hashable)
	if !ok {
		return object.NewException("Invalid index type %s", index.Type())
	}

	hashmap.Pairs[hashable.HashKey()] = object.HashPair{
		Key:   index,
		Value: val,
	}
	return object.NullConst
}

func (vm *VirtualMachine) evalBinaryExpression(op string, left, right object.Object) object.Object {
	switch {
	case left.Type() != right.Type():
		panic(object.NewException("type mismatch: %s %s %s", left.Type(), op, right.Type()))
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
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	case "%":
		return &object.Integer{Value: leftVal % rightVal}
	case "<<":
		if rightVal < 0 {
			return object.NewException("Shift value must be non-negative")
		}
		return &object.Integer{Value: leftVal << uint64(rightVal)}
	case ">>":
		if rightVal < 0 {
			return object.NewException("Shift value must be non-negative")
		}
		return &object.Integer{Value: leftVal >> uint64(rightVal)}
	case "&":
		return &object.Integer{Value: leftVal & rightVal}
	case "&^":
		return &object.Integer{Value: leftVal &^ rightVal}
	case "|":
		return &object.Integer{Value: leftVal | rightVal}
	case "^":
		return &object.Integer{Value: leftVal ^ rightVal}
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

	switch op {
	case "+":
		return &object.String{Value: leftVal + rightVal}
	}

	return object.NewException("unknown operator: %s %s %s", left.Type(), op, right.Type())
}

func (vm *VirtualMachine) evalArrayBinaryExpression(op string, left, right object.Object) object.Object {
	leftVal := left.(*object.Array)
	rightVal := right.(*object.Array)

	switch op {
	case "+":
		leftLen := len(leftVal.Elements)
		rightLen := len(rightVal.Elements)
		newElements := make([]object.Object, leftLen+rightLen, leftLen+rightLen)
		copy(newElements, leftVal.Elements)
		copy(newElements[leftLen:], rightVal.Elements)
		return &object.Array{Elements: newElements}
	}

	return object.NewException("unknown operator: %s %s %s", left.Type(), op, right.Type())
}
