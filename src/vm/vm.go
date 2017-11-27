package vm

import (
	"fmt"
	"io"
	"os"
	"runtime/debug"

	"github.com/nitrogen-lang/nitrogen/src/ast"
	"github.com/nitrogen-lang/nitrogen/src/compiler"
	"github.com/nitrogen-lang/nitrogen/src/object"
	"github.com/nitrogen-lang/nitrogen/src/vm/opcode"
)

type Settings struct {
	Debug bool

	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

func NewSettings() *Settings {
	return &Settings{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
}

type VirtualMachine struct {
	callStack    *frameStack
	currentFrame *Frame
	returnValue  object.Object
	settings     *Settings
}

func NewVM(settings *Settings) *VirtualMachine {
	if settings == nil {
		settings = &Settings{
			Stdin:  os.Stdin,
			Stdout: os.Stdout,
			Stderr: os.Stderr,
		}
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
func (vm *VirtualMachine) GetStdout() io.Writer         { return vm.settings.Stdout }
func (vm *VirtualMachine) GetStderr() io.Writer         { return vm.settings.Stderr }
func (vm *VirtualMachine) GetStdin() io.Reader          { return vm.settings.Stdin }

func (vm *VirtualMachine) Execute(code *compiler.CodeBlock, env *object.Environment) object.Object {
	if env == nil {
		env = object.NewEnvironment()
	}
	return vm.RunFrame(vm.MakeFrame(code, env), false)
}

func (vm *VirtualMachine) CurrentFrame() *Frame {
	return vm.currentFrame
}

func (vm *VirtualMachine) MakeFrame(code *compiler.CodeBlock, env *object.Environment) *Frame {
	return &Frame{
		code:       code,
		stack:      make([]object.Object, code.MaxStackSize+1), // +1 to make room for a runtime exception if thrown
		blockStack: make([]block, code.MaxBlockSize),
		Env:        env,
	}
}

func (vm *VirtualMachine) RunFrame(f *Frame, immediateReturn bool) object.Object {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(object.Object); ok {
				fmt.Fprintln(vm.GetStderr(), r)
				fmt.Fprintln(vm.GetStderr(), "Stack Trace:")
				frame := vm.currentFrame
				for frame != nil {
					fmt.Fprintf(vm.GetStderr(), "\t%s: %s\n", frame.code.Filename, frame.code.Name)
					frame = frame.lastFrame
				}
			} else {
				fmt.Fprintln(vm.GetStderr(), r)
				fmt.Fprintln(vm.GetStderr(), string(debug.Stack()))
			}
		}
	}()
	f.lastFrame = vm.currentFrame
	vm.callStack.Push(f)
	vm.currentFrame = f

mainLoop:
	for {
		if vm.currentFrame.pc >= len(vm.currentFrame.code.Code) {
			panic(fmt.Sprintf("Program counter %d outside bounds of bytecode %d", vm.currentFrame.pc, len(vm.currentFrame.code.Code)))
		}
		code := vm.fetchByte()
		if vm.settings.Debug {
			fmt.Fprintf(vm.GetStdout(), "Executing %s\n", opcode.Names[code])
		}

		switch code {
		case opcode.Noop:
			continue mainLoop
		case opcode.BinaryAdd:
			r := vm.currentFrame.popStack()
			l := vm.currentFrame.popStack()
			res := vm.evalBinaryExpression("+", l, r)
			vm.currentFrame.pushStack(res)
			if object.ObjectIs(res, object.ExceptionObj) {
				vm.throw()
			}
		case opcode.BinarySub:
			r := vm.currentFrame.popStack()
			l := vm.currentFrame.popStack()
			res := vm.evalBinaryExpression("-", l, r)
			vm.currentFrame.pushStack(res)
			if object.ObjectIs(res, object.ExceptionObj) {
				vm.throw()
			}
		case opcode.BinaryMul:
			r := vm.currentFrame.popStack()
			l := vm.currentFrame.popStack()
			res := vm.evalBinaryExpression("*", l, r)
			vm.currentFrame.pushStack(res)
			if object.ObjectIs(res, object.ExceptionObj) {
				vm.throw()
			}
		case opcode.BinaryDivide:
			r := vm.currentFrame.popStack()
			l := vm.currentFrame.popStack()
			res := vm.evalBinaryExpression("/", l, r)
			vm.currentFrame.pushStack(res)
			if object.ObjectIs(res, object.ExceptionObj) {
				vm.throw()
			}
		case opcode.BinaryMod:
			r := vm.currentFrame.popStack()
			l := vm.currentFrame.popStack()
			res := vm.evalBinaryExpression("%", l, r)
			vm.currentFrame.pushStack(res)
			if object.ObjectIs(res, object.ExceptionObj) {
				vm.throw()
			}
		case opcode.BinaryShiftL:
			r := vm.currentFrame.popStack()
			l := vm.currentFrame.popStack()
			res := vm.evalBinaryExpression("<<", l, r)
			vm.currentFrame.pushStack(res)
			if object.ObjectIs(res, object.ExceptionObj) {
				vm.throw()
			}
		case opcode.BinaryShiftR:
			r := vm.currentFrame.popStack()
			l := vm.currentFrame.popStack()
			res := vm.evalBinaryExpression(">>", l, r)
			vm.currentFrame.pushStack(res)
			if object.ObjectIs(res, object.ExceptionObj) {
				vm.throw()
			}
		case opcode.BinaryAnd:
			r := vm.currentFrame.popStack()
			l := vm.currentFrame.popStack()
			res := vm.evalBinaryExpression("&", l, r)
			vm.currentFrame.pushStack(res)
			if object.ObjectIs(res, object.ExceptionObj) {
				vm.throw()
			}
		case opcode.BinaryOr:
			r := vm.currentFrame.popStack()
			l := vm.currentFrame.popStack()
			res := vm.evalBinaryExpression("|", l, r)
			vm.currentFrame.pushStack(res)
			if object.ObjectIs(res, object.ExceptionObj) {
				vm.throw()
			}
		case opcode.BinaryNot:
			r := vm.currentFrame.popStack()
			l := vm.currentFrame.popStack()
			res := vm.evalBinaryExpression("^", l, r)
			vm.currentFrame.pushStack(res)
			if object.ObjectIs(res, object.ExceptionObj) {
				vm.throw()
			}
		case opcode.BinaryAndNot:
			r := vm.currentFrame.popStack()
			l := vm.currentFrame.popStack()
			res := vm.evalBinaryExpression("&^", l, r)
			vm.currentFrame.pushStack(res)
			if object.ObjectIs(res, object.ExceptionObj) {
				vm.throw()
			}
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
				vm.currentFrame.pushStack(object.NewException("Redefined constant %s\n", name))
				vm.throw()
			}
			if _, err := vm.currentFrame.Env.CreateConst(name, vm.currentFrame.popStack()); err != nil {
				fmt.Println(err)
			}
		case opcode.Return:
			vm.returnValue = vm.currentFrame.popStack()
			vm.currentFrame = vm.currentFrame.lastFrame
			vm.callStack.Pop()
			if vm.currentFrame == nil || immediateReturn {
				return vm.returnValue
			}
			vm.currentFrame.pushStack(vm.returnValue)
		case opcode.Pop:
			vm.currentFrame.popStack()
		case opcode.LoadFast:
			name := vm.currentFrame.code.Locals[vm.getUint16()]
			if val, ok := vm.currentFrame.Env.GetLocal(name); ok {
				vm.currentFrame.pushStack(val)
				break
			}

			vm.currentFrame.pushStack(object.NewException("Unknown variable/constant %s\n", name))
			vm.throw()
		case opcode.StoreFast:
			// Ensure constant isn't redefined
			name := vm.currentFrame.code.Locals[vm.getUint16()]
			if vm.currentFrame.Env.IsConst(name) {
				vm.currentFrame.pushStack(object.NewException("Redefined constant %s\n", name))
				vm.throw()
			}
			if _, exists := vm.currentFrame.Env.GetLocal(name); !exists {
				vm.currentFrame.pushStack(object.NewException("Variable %s undefined\n", name))
				vm.throw()
			}
			vm.currentFrame.Env.SetLocal(name, vm.currentFrame.popStack())
		case opcode.Define:
			// Ensure constant isn't redefined
			name := vm.currentFrame.code.Locals[vm.getUint16()]
			if vm.currentFrame.Env.IsConst(name) {
				vm.currentFrame.pushStack(object.NewException("Redefined constant %s\n", name))
				vm.throw()
			}
			if _, exists := vm.currentFrame.Env.GetLocal(name); exists {
				vm.currentFrame.pushStack(object.NewException("Variable %s already defined\n", name))
				vm.throw()
			}
			vm.currentFrame.Env.Create(name, vm.currentFrame.popStack())
		case opcode.LoadGlobal:
			name := vm.currentFrame.code.Names[vm.getUint16()]
			p := vm.currentFrame.Env.Parent()
			if p == nil {
				vm.currentFrame.pushStack(object.NewException("Global variable %s not defined\n", name))
				vm.throw()
			}
			if val, ok := p.Get(name); ok {
				vm.currentFrame.pushStack(val)
				break
			}
			if fn := getBuiltin(name); fn != nil {
				vm.currentFrame.pushStack(fn)
				break
			}

			vm.currentFrame.pushStack(object.NewException("Global %s doesn't exist\n", name))
			vm.throw()
		case opcode.StoreGlobal:
			// Ensure constant isn't redefined
			name := vm.currentFrame.code.Names[vm.getUint16()]
			if vm.currentFrame.Env.IsConst(name) {
				vm.currentFrame.pushStack(object.NewException("Redefined constant %s\n", name))
				vm.throw()
			}
			p := vm.currentFrame.Env.Parent()
			if p == nil {
				vm.currentFrame.pushStack(object.NewException("Global variable %s not defined\n", name))
				vm.throw()
			}
			if _, exists := p.Get(name); !exists {
				vm.currentFrame.pushStack(object.NewException("Global variable %s not defined\n", name))
				vm.throw()
			}
			vm.currentFrame.Env.Set(name, vm.currentFrame.popStack())
		case opcode.LoadIndex:
			left := vm.currentFrame.popStack()
			index := vm.currentFrame.popStack()
			res := vm.evalIndexExpression(left, index)
			vm.currentFrame.pushStack(res)
			if object.ObjectIs(res, object.ExceptionObj) {
				vm.throw()
			}
		case opcode.StoreIndex:
			left := vm.currentFrame.popStack()
			index := vm.currentFrame.popStack()
			value := vm.currentFrame.popStack()
			res := vm.assignIndexedValue(left, index, value)
			if object.ObjectIs(res, object.ExceptionObj) {
				vm.currentFrame.pushStack(res)
				vm.throw()
			}
		case opcode.Call:
			numargs := vm.getUint16()
			fn := vm.currentFrame.popStack()

			switch fn := fn.(type) {
			case *object.Builtin:
				args := make([]object.Object, numargs)
				for i := uint16(0); i < numargs; i++ {
					args[i] = vm.currentFrame.popStack()
				}

				result := fn.Fn(vm, vm.currentFrame.Env, args...)
				if result == nil {
					result = object.NullConst
				}

				vm.returnValue = result
				vm.currentFrame.pushStack(result)

				if object.ObjectIs(result, object.ExceptionObj) {
					vm.throw()
				}
			case *VMFunction:
				newFrame := vm.MakeFrame(fn.Body, object.NewSizedEnclosedEnv(fn.Env, fn.Body.LocalCount))
				newFrame.lastFrame = vm.currentFrame

				for i := 0; i < int(numargs); i++ {
					newFrame.Env.SetForce(fn.Parameters[i], vm.currentFrame.popStack(), false)
				}

				vm.currentFrame = newFrame
				vm.callStack.Push(newFrame)
			default:
				for i := 0; i < int(numargs); i++ {
					vm.currentFrame.popStack()
				}
				vm.currentFrame.pushStack(object.NewPanic("TOS is not a function for CALL opcode"))
				vm.throw()
			}
		case opcode.Compare:
			r := vm.currentFrame.popStack()
			l := vm.currentFrame.popStack()
			op := vm.fetchByte()
			if op >= opcode.MaxCmpCodes {
				ex := object.NewPanic("Invalid comparison operator %x", op)
				vm.currentFrame.pushStack(ex)
				vm.throw()
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
					ex := object.NewPanic("Map key %s not valid", key.Inspect())
					vm.currentFrame.pushStack(ex)
					vm.throw()
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
			lb := &forLoopBlock{
				start: vm.currentFrame.pc,
				iter:  int(iter),
				end:   int(loopEnd),
			}
			vm.currentFrame.pushBlock(lb)
			vm.currentFrame.Env = object.NewEnclosedEnv(vm.currentFrame.Env)
		case opcode.Continue:
			vm.currentFrame.pc = vm.currentFrame.popBlockUntil(loopBlockT).(*forLoopBlock).iter
		case opcode.NextIter:
			vm.currentFrame.pc = vm.currentFrame.popBlockUntil(loopBlockT).(*forLoopBlock).start
			vm.currentFrame.Env = object.NewEnclosedEnv(vm.currentFrame.Env.Parent())
		case opcode.Break:
			vm.currentFrame.pc = vm.currentFrame.popBlockUntil(loopBlockT).(*forLoopBlock).end
		case opcode.StartTry:
			catch := vm.getUint16()
			tcb := &tryBlock{
				catch: int(catch),
				sp:    vm.currentFrame.sp,
			}
			vm.currentFrame.pushBlock(tcb)
			vm.currentFrame.Env = object.NewEnclosedEnv(vm.currentFrame.Env)
		case opcode.Throw:
			vm.throw()
		default:
			codename := opcode.Names[code]
			if codename == "" {
				codename = fmt.Sprintf("%X", code)
			}
			ex := object.NewPanic("Opcode %s is not supported", codename)
			vm.currentFrame.pushStack(ex)
			vm.throw()
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

// throw takes the top of stack item as an exception object
// it will then progressivly unwind the block stack and call stack until
// a try block is found. If none is found, it will panic with an uncaught
// exception and the VM will print a stack trace.
func (vm *VirtualMachine) throw() {
	exception := vm.currentFrame.popStack()
	if exception.Type() != object.ExceptionObj {
		exception = object.NewException(exception.Inspect())
	}

	if ex := exception.(*object.Exception); !ex.Catchable {
		panic(object.NewException("Runtime Exception: %s", exception.Inspect()))
	}

	vm.currentFrame.Env = vm.currentFrame.Env.Parent().Parent() // Unwind block scoping
	cframe := vm.currentFrame
	for {
		// Unwind block stack until there's a try block
		catchBlock := vm.currentFrame.popBlockUntil(tryBlockT)
		if catchBlock != nil { // Try block found
			tryBlockS := catchBlock.(*tryBlock)
			vm.currentFrame.sp = tryBlockS.sp    // Unwind data stack
			vm.currentFrame.pc = tryBlockS.catch // Set program counter to catch block
			break
		}
		vm.currentFrame = vm.currentFrame.lastFrame // This frame doesn't have a try block, unwind call stack
		if vm.currentFrame == nil {                 // Call stack exhausted
			vm.currentFrame = cframe // Reset frame for stack trace
			panic(object.NewException("Uncaught Exception: %s", exception.Inspect()))
		}
	}

	// Enclose once for new block (like a PREPARE_BLOCK) and another for block scope
	// END_BLOCK removes two layers of environments
	vm.currentFrame.Env = object.NewEnclosedEnv(object.NewEnclosedEnv(vm.currentFrame.Env))
	vm.currentFrame.pushStack(exception)
}
