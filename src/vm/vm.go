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
	Debug            bool
	ReturnExceptions bool

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

	unwind bool
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
		env:        env,
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
				vm.unwind = true
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
		if vm.unwind || vm.currentFrame == nil {
			if vm.returnValue == nil {
				vm.returnValue = object.NullConst
			}
			return vm.returnValue
		}

		if vm.currentFrame.pc >= len(vm.currentFrame.code.Code) {
			panic(fmt.Sprintf("Program counter %d outside bounds of bytecode %d", vm.currentFrame.pc, len(vm.currentFrame.code.Code)-1))
		}
		code := vm.fetchOpcode()
		if vm.settings.Debug {
			fmt.Fprintf(vm.GetStdout(), "Executing %d -> %s\n", vm.currentFrame.pc-1, opcode.Names[code])
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
				break
			}

		case opcode.BinarySub:
			r := vm.currentFrame.popStack()
			l := vm.currentFrame.popStack()
			res := vm.evalBinaryExpression("-", l, r)
			vm.currentFrame.pushStack(res)
			if object.ObjectIs(res, object.ExceptionObj) {
				vm.throw()
				break
			}

		case opcode.BinaryMul:
			r := vm.currentFrame.popStack()
			l := vm.currentFrame.popStack()
			res := vm.evalBinaryExpression("*", l, r)
			vm.currentFrame.pushStack(res)
			if object.ObjectIs(res, object.ExceptionObj) {
				vm.throw()
				break
			}

		case opcode.BinaryDivide:
			r := vm.currentFrame.popStack()
			l := vm.currentFrame.popStack()
			res := vm.evalBinaryExpression("/", l, r)
			vm.currentFrame.pushStack(res)
			if object.ObjectIs(res, object.ExceptionObj) {
				vm.throw()
				break
			}

		case opcode.BinaryMod:
			r := vm.currentFrame.popStack()
			l := vm.currentFrame.popStack()
			res := vm.evalBinaryExpression("%", l, r)
			vm.currentFrame.pushStack(res)
			if object.ObjectIs(res, object.ExceptionObj) {
				vm.throw()
				break
			}

		case opcode.BinaryShiftL:
			r := vm.currentFrame.popStack()
			l := vm.currentFrame.popStack()
			res := vm.evalBinaryExpression("<<", l, r)
			vm.currentFrame.pushStack(res)
			if object.ObjectIs(res, object.ExceptionObj) {
				vm.throw()
				break
			}

		case opcode.BinaryShiftR:
			r := vm.currentFrame.popStack()
			l := vm.currentFrame.popStack()
			res := vm.evalBinaryExpression(">>", l, r)
			vm.currentFrame.pushStack(res)
			if object.ObjectIs(res, object.ExceptionObj) {
				vm.throw()
				break
			}

		case opcode.BinaryAnd:
			r := vm.currentFrame.popStack()
			l := vm.currentFrame.popStack()
			res := vm.evalBinaryExpression("&", l, r)
			vm.currentFrame.pushStack(res)
			if object.ObjectIs(res, object.ExceptionObj) {
				vm.throw()
				break
			}

		case opcode.BinaryOr:
			r := vm.currentFrame.popStack()
			l := vm.currentFrame.popStack()
			res := vm.evalBinaryExpression("|", l, r)
			vm.currentFrame.pushStack(res)
			if object.ObjectIs(res, object.ExceptionObj) {
				vm.throw()
				break
			}

		case opcode.BinaryNot:
			r := vm.currentFrame.popStack()
			l := vm.currentFrame.popStack()
			res := vm.evalBinaryExpression("^", l, r)
			vm.currentFrame.pushStack(res)
			if object.ObjectIs(res, object.ExceptionObj) {
				vm.throw()
				break
			}

		case opcode.BinaryAndNot:
			r := vm.currentFrame.popStack()
			l := vm.currentFrame.popStack()
			res := vm.evalBinaryExpression("&^", l, r)
			vm.currentFrame.pushStack(res)
			if object.ObjectIs(res, object.ExceptionObj) {
				vm.throw()
				break
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
			if vm.currentFrame.env.IsConst(name) {
				vm.currentFrame.pushStack(object.NewException("Redefined constant %s\n", name))
				vm.throw()
				break
			}
			if _, err := vm.currentFrame.env.CreateConst(name, vm.currentFrame.popStack()); err != nil {
				fmt.Println(err)
			}

		case opcode.Return:
			if vm.currentFrame.sp == 0 {
				vm.returnValue = object.NullConst
			} else {
				vm.returnValue = vm.currentFrame.popStack()
			}

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
			if val, ok := vm.currentFrame.env.GetLocal(name); ok {
				vm.currentFrame.pushStack(val)
				break
			}

			vm.currentFrame.pushStack(object.NewException("Unknown variable/constant %s\n", name))
			vm.throw()
			break

		case opcode.StoreFast:
			// Ensure constant isn't redefined
			name := vm.currentFrame.code.Locals[vm.getUint16()]
			if vm.currentFrame.env.IsConstLocal(name) {
				vm.currentFrame.pushStack(object.NewException("Redefined constant %s\n", name))
				vm.throw()
				break
			}
			if _, exists := vm.currentFrame.env.GetLocal(name); !exists {
				vm.currentFrame.pushStack(object.NewException("Variable %s undefined\n", name))
				vm.throw()
				break
			}
			vm.currentFrame.env.SetLocal(name, vm.currentFrame.popStack())

		case opcode.DeleteFast:
			name := vm.currentFrame.code.Locals[vm.getUint16()]
			if vm.currentFrame.env.IsConstLocal(name) {
				vm.currentFrame.pushStack(object.NewException("Cannot delete constant %s\n", name))
				vm.throw()
				break
			}
			vm.currentFrame.env.UnsetLocal(name)

		case opcode.Define:
			// Ensure constant isn't redefined
			name := vm.currentFrame.code.Locals[vm.getUint16()]
			if vm.currentFrame.env.IsConst(name) {
				vm.currentFrame.pushStack(object.NewException("Redefined constant %s\n", name))
				vm.throw()
				break
			}
			if _, exists := vm.currentFrame.env.GetLocal(name); exists {
				vm.currentFrame.pushStack(object.NewException("Variable %s already defined\n", name))
				vm.throw()
				break
			}
			vm.currentFrame.env.Create(name, vm.currentFrame.popStack())

		case opcode.LoadGlobal:
			name := vm.currentFrame.code.Names[vm.getUint16()]
			p := vm.currentFrame.env.Parent()
			if p == nil {
				p = vm.currentFrame.env
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
			break

		case opcode.StoreGlobal:
			// Ensure constant isn't redefined
			name := vm.currentFrame.code.Names[vm.getUint16()]
			if vm.currentFrame.env.IsConst(name) {
				vm.currentFrame.pushStack(object.NewException("Redefined constant %s\n", name))
				vm.throw()
				break
			}
			p := vm.currentFrame.env.Parent()
			if p == nil {
				p = vm.currentFrame.env
			}
			if _, exists := p.Get(name); !exists {
				vm.currentFrame.pushStack(object.NewException("Global variable %s not defined\n", name))
				vm.throw()
				break
			}
			vm.currentFrame.env.Set(name, vm.currentFrame.popStack())

		case opcode.LoadIndex:
			left := vm.currentFrame.popStack()
			index := vm.currentFrame.popStack()
			res := vm.evalIndexExpression(left, index)
			vm.currentFrame.pushStack(res)
			if object.ObjectIs(res, object.ExceptionObj) {
				vm.throw()
				break
			}

		case opcode.StoreIndex:
			left := vm.currentFrame.popStack()
			index := vm.currentFrame.popStack()
			value := vm.currentFrame.popStack()
			res := vm.assignIndexedValue(left, index, value)
			if object.ObjectIs(res, object.ExceptionObj) {
				vm.currentFrame.pushStack(res)
				vm.throw()
				break
			}

		case opcode.Call:
			numargs := vm.getUint16()
			fn := vm.currentFrame.popStack()
			vm.CallFunction(numargs, fn, false)

		case opcode.Compare:
			r := vm.currentFrame.popStack()
			l := vm.currentFrame.popStack()
			op := vm.fetchByte()
			if op >= opcode.MaxCmpCodes {
				ex := object.NewPanic("Invalid comparison operator %x", op)
				vm.currentFrame.pushStack(ex)
				vm.throw()
				break
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
				Env:        object.NewEnclosedEnv(vm.currentFrame.env),
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
					break
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

		case opcode.OpenScope:
			vm.currentFrame.env = object.NewEnclosedEnv(vm.currentFrame.env)

		case opcode.CloseScope:
			vm.currentFrame.env = vm.currentFrame.env.Parent()

		case opcode.EndBlock:
			vm.currentFrame.popBlock()
			if vm.currentFrame.sp == 0 {
				vm.currentFrame.pushStack(object.NullConst)
			}

		case opcode.StartLoop:
			loopEnd := vm.getUint16()
			iter := vm.getUint16()
			lb := &forLoopBlock{
				start: vm.currentFrame.pc,
				iter:  int(iter),
				end:   int(loopEnd),
			}
			vm.currentFrame.pushBlock(lb)
			vm.currentFrame.env = object.NewEnclosedEnv(vm.currentFrame.env)

		case opcode.Continue:
			vm.currentFrame.pc = vm.currentFrame.popBlockUntil(loopBlockT).(*forLoopBlock).iter

		case opcode.NextIter:
			vm.currentFrame.pc = vm.currentFrame.popBlockUntil(loopBlockT).(*forLoopBlock).start
			vm.currentFrame.env = object.NewEnclosedEnv(vm.currentFrame.env.Parent())

		case opcode.Break:
			vm.currentFrame.pc = vm.currentFrame.popBlockUntil(loopBlockT).(*forLoopBlock).end

		case opcode.Import:
			name := vm.currentFrame.code.Locals[vm.getUint16()]
			path := vm.currentFrame.popStack().(*object.String)
			vm.importPackage(name, path.Value)

		case opcode.StartTry:
			catch := vm.getUint16()
			tcb := &tryBlock{
				catch: int(catch),
				sp:    vm.currentFrame.sp,
			}
			vm.currentFrame.pushBlock(tcb)
			vm.currentFrame.env = object.NewEnclosedEnv(vm.currentFrame.env)

		case opcode.Throw:
			exc := vm.throw()
			if exc != nil {
				return exc
			}

		case opcode.BuildClass:
			methodNum := vm.getUint16()
			class := &VMClass{}
			class.Name = vm.currentFrame.popStack().(*object.String).Value
			parent := vm.currentFrame.popStack()
			if parent != object.NullConst {
				class.Parent = parent.(*VMClass)
			}
			class.Fields = vm.currentFrame.popStack().(*compiler.CodeBlock)
			class.Methods = make(map[string]object.ClassMethod, methodNum)
			for i := methodNum; i > 0; i-- {
				method := vm.currentFrame.popStack().(*VMFunction)
				class.Methods[method.Name] = method
				method.Class = class
			}
			vm.currentFrame.pushStack(class)

		case opcode.MakeInstance:
			argLen := vm.getUint16()
			class := vm.currentFrame.popStack()
			vm.makeInstance(argLen, class)

		case opcode.LoadAttribute:
			name := vm.currentFrame.code.Names[vm.getUint16()]
			instance := vm.currentFrame.popStack()

			switch instance := instance.(type) {
			case *VMInstance:
				if method := instance.GetBoundMethod(name); method != nil {
					vm.currentFrame.pushStack(method)
				} else {
					val, ok := instance.Fields.Get(name)
					if ok {
						vm.currentFrame.pushStack(val)
					} else {
						vm.currentFrame.pushStack(object.NullConst)
					}
				}
			case *VMClass:
				this := vm.currentFrame.frontInstance()
				if this == nil {
					vm.currentFrame.pushStack(object.NewException("Method call outside instance"))
					vm.throw()
					break
				}

				method := instance.GetMethod(name)
				if method != nil {
					vm.currentFrame.pushStack(&BoundMethod{
						Method:   method,
						Instance: vm.currentFrame.frontInstance(),
						Parent:   vm.currentFrame.frontInstance().Class.Parent,
					})
				} else {
					vm.currentFrame.pushStack(object.NullConst)
				}
			case *object.Module:
				vm.currentFrame.pushStack(vm.lookupModuleAttr(instance, name))
			case *object.Hash:
				vm.currentFrame.pushStack(vm.lookupHashIndex(instance, object.MakeStringObj(name)))
			default:
				vm.currentFrame.pushStack(object.NewPanic("Attribute lookup on non-object"))
				vm.throw()
			}

		case opcode.StoreAttribute:
			name := vm.currentFrame.code.Names[vm.getUint16()]
			instance := vm.currentFrame.popStack()
			val := vm.currentFrame.popStack()

			switch instance := instance.(type) {
			case *VMInstance:
				if _, ok := instance.Fields.Get(name); !ok {
					vm.currentFrame.pushStack(object.NewException("Instance has no field %s", name))
					vm.throw()
					break
				}

				if instance.Fields.IsConst(name) {
					vm.currentFrame.pushStack(object.NewException("Assignment to constant field %s", name))
					vm.throw()
					break
				}
				instance.Fields.SetForce(name, val, false)
			case *object.Module:
				ret := vm.assignModuleAttr(instance, name, val)
				if ret != object.NullConst {
					vm.currentFrame.pushStack(ret)
					vm.throw()
				}
			case *object.Hash:
				vm.assignHashMapIndex(instance, object.MakeStringObj(name), val)
			default:
				vm.currentFrame.pushStack(object.NewPanic("Attribute lookup on non-object"))
				vm.throw()
			}

		default:
			codename := opcode.Names[code]
			if codename == "" {
				codename = fmt.Sprintf("%X", code)
			}
			ex := object.NewPanic("Opcode %s is not supported", codename)
			vm.currentFrame.pushStack(ex)
			vm.throw()
			break
		}
	}
}

func (vm *VirtualMachine) fetchOpcode() opcode.Opcode {
	return opcode.Opcode(vm.fetchByte())
}

func (vm *VirtualMachine) fetchByte() byte {
	b := vm.currentFrame.code.Code[vm.currentFrame.pc]
	vm.currentFrame.pc++
	return b
}

func (vm *VirtualMachine) getUint16() uint16 {
	return (uint16(vm.fetchByte()) << 8) + uint16(vm.fetchByte())
}

func (vm *VirtualMachine) PopStack() object.Object {
	return vm.currentFrame.popStack()
}

// throw takes the top of stack item as an exception object
// it will then progressivly unwind the block stack and call stack until
// a try block is found. If none is found, it will panic with an uncaught
// exception and the VM will print a stack trace.
func (vm *VirtualMachine) throw() object.Object {
	exception := vm.currentFrame.popStack()
	if exception.Type() != object.ExceptionObj {
		exception = object.NewException(exception.Inspect())
	}

	if ex := exception.(*object.Exception); !ex.Catchable {
		panic(object.NewException("Runtime Exception: %s", exception.Inspect()))
	}

	cframe := vm.currentFrame
	for {
		// Unwind block stack until there's a try block
		catchBlock := vm.currentFrame.popBlockUntil(tryBlockT)
		if catchBlock != nil { // Try block found
			tryBlockS := catchBlock.(*tryBlock)
			if !tryBlockS.caught {
				tryBlockS.caught = true
				vm.currentFrame.sp = tryBlockS.sp    // Unwind data stack
				vm.currentFrame.pc = tryBlockS.catch // Set program counter to catch block
				break
			}
		}
		vm.currentFrame = vm.currentFrame.lastFrame // This frame doesn't have a try block, unwind call stack
		if vm.currentFrame == nil {                 // Call stack exhausted
			exc := object.NewException("Uncaught Exception: %s", exception.Inspect())

			if vm.settings.ReturnExceptions {
				return exc
			}

			vm.currentFrame = cframe // Reset frame for stack trace
			panic(object.NewException("Uncaught Exception: %s", exception.Inspect()))
		}
	}

	vm.currentFrame.env = vm.currentFrame.env.Parent().Parent() // Unwind block scoping
	// Enclose once for new block (like a PREPARE_BLOCK) and another for block scope
	// END_BLOCK removes two layers of environments
	vm.currentFrame.env = object.NewEnclosedEnv(object.NewEnclosedEnv(vm.currentFrame.env))
	vm.currentFrame.pushStack(exception)
	return nil
}

func (vm *VirtualMachine) makeInstance(argLen uint16, class object.Object) {
	var instance *VMInstance

	if class, ok := class.(*VMClass); ok {
		cClass := class
		classChain := make([]*VMClass, 1, 3)
		classChain[0] = class
		for cClass.Parent != nil {
			classChain = append(classChain, cClass.Parent)
			cClass = cClass.Parent
		}

		iFields := object.NewEnvironment()
		iFields.SetParent(vm.currentFrame.env)

		for _, c := range classChain {
			frame := vm.MakeFrame(c.Fields, iFields)
			vm.RunFrame(frame, true)
		}

		iFields.SetParent(nil)

		instance = &VMInstance{
			Class:  class,
			Fields: iFields,
		}
	}

	if class, ok := class.(*BuiltinClass); ok {
		iFields := object.NewEnvironment()
		for k, v := range class.Fields {
			iFields.SetForce(k, v, false)
		}

		instance = &VMInstance{
			Class:  class.VMClass,
			Fields: iFields,
		}
	}

	init := instance.GetBoundMethod("init")
	if init == nil {
		for i := argLen; i > 0; i-- {
			vm.currentFrame.popStack()
		}
		vm.currentFrame.pushStack(instance)
		return
	}

	vm.CallFunction(argLen, init, true)
	vm.currentFrame.pushStack(instance)
	return
}

func (vm *VirtualMachine) CallFunction(argc uint16, fn object.Object, now bool) {
	switch fn := fn.(type) {
	case *object.Builtin:
		args := make([]object.Object, argc)
		for i := uint16(0); i < argc; i++ {
			args[i] = vm.currentFrame.popStack()
		}

		env := vm.currentFrame.env
		this := vm.currentFrame.frontInstance()
		if this != nil {
			env = object.NewEnclosedEnv(env)
			env.SetForce("this", this, true)
		}

		result := fn.Fn(vm, env, args...)
		if result == nil {
			result = object.NullConst
		}

		vm.returnValue = result
		vm.currentFrame.pushStack(result)

		if object.ObjectIs(result, object.ExceptionObj) {
			vm.throw()
		}
	case *BuiltinMethod:
		args := make([]object.Object, argc)
		for i := uint16(0); i < argc; i++ {
			args[i] = vm.currentFrame.popStack()
		}

		this := vm.currentFrame.frontInstance()
		result := fn.Fn(vm, this, this.Fields, args...)
		if result == nil {
			result = object.NullConst
		}

		vm.returnValue = result
		vm.currentFrame.pushStack(result)

		if object.ObjectIs(result, object.ExceptionObj) {
			vm.throw()
		}
	case *VMFunction:
		var env *object.Environment
		this := vm.currentFrame.frontInstance()
		env = object.NewEnclosedEnv(fn.Env)
		if this != nil {
			env.SetForce("this", this, true)
			if fn.Class != nil && fn.Class.Parent != nil {
				env.SetForce("parent", fn.Class.Parent, true)
			}
		}

		paramLen := len(fn.Parameters)

		if int(argc) < paramLen {
			vm.currentFrame.pushStack(object.NewException("Func expected %d args but was given %d", paramLen, argc))
			vm.throw()
			return
		}

		newFrame := vm.MakeFrame(fn.Body, env)
		newFrame.lastFrame = vm.currentFrame
		if this != nil {
			newFrame.pushInstance(this)
		}

		for i := 0; i < paramLen; i++ {
			newFrame.env.SetForce(fn.Parameters[i], vm.currentFrame.popStack(), false)
		}

		if int(argc) > paramLen {
			remaining := int(argc) - paramLen
			rest := make([]object.Object, remaining)
			for i := 0; i < remaining; i++ {
				rest[i] = vm.currentFrame.popStack()
			}
			newFrame.env.SetForce("arguments", &object.Array{Elements: rest}, false)
		} else {
			newFrame.env.SetForce("arguments", &object.Array{Elements: []object.Object{}}, false)
		}

		if now {
			vm.currentFrame.pushStack(vm.RunFrame(newFrame, true))
		} else {
			vm.currentFrame = newFrame
			vm.callStack.Push(newFrame)
		}
	case *BoundMethod:
		vm.currentFrame.pushInstance(fn.Instance)
		vm.CallFunction(argc, fn.Method, true)
		vm.currentFrame.popInstance()
	case *VMClass:
		this := vm.currentFrame.frontInstance()
		if this == nil {
			vm.currentFrame.pushStack(object.NewPanic("Can't call class method outside of object"))
			vm.throw()
			return
		}

		if !InstanceOf(fn.Name, this) {
			vm.currentFrame.pushStack(object.NewPanic("Object not instance of %s", fn.Name))
			vm.throw()
			return
		}

		init := fn.GetMethod("init")
		if init == nil {
			for i := argc; i > 0; i-- {
				vm.currentFrame.popStack()
			}
			return
		}
		vm.CallFunction(argc, init, true)
	default:
		for i := 0; i < int(argc); i++ {
			vm.currentFrame.popStack()
		}
		vm.currentFrame.pushStack(object.NewPanic("%s is not a function", fn.Type()))
		vm.throw()
	}
}
