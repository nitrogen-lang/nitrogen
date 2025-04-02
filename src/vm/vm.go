package vm

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"runtime/debug"
	"strconv"

	"github.com/nitrogen-lang/nitrogen/src/ast"
	"github.com/nitrogen-lang/nitrogen/src/compiler"
	"github.com/nitrogen-lang/nitrogen/src/object"
	"github.com/nitrogen-lang/nitrogen/src/vm/opcode"
)

type ErrExitCode struct {
	Code int
}

func (e ErrExitCode) Error() string {
	return strconv.Itoa(e.Code)
}

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
	returnErr    error
	Settings     *Settings
	globalEnv    *object.Environment
	instanceVars map[string]object.Object

	breakpoint bool
	unwind     bool
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
		callStack:    newFrameStack(),
		Settings:     settings,
		instanceVars: make(map[string]object.Object),
	}
}

func (vm *VirtualMachine) SetGlobalEnv(env *object.Environment) {
	vm.globalEnv = env
}

func (vm *VirtualMachine) Eval(node ast.Node, env *object.Environment) object.Object {
	return object.NullConst
}
func (vm *VirtualMachine) GetCurrentScriptPath() string { return vm.currentFrame.code.Filename }
func (vm *VirtualMachine) GetStdout() io.Writer         { return vm.Settings.Stdout }
func (vm *VirtualMachine) GetStderr() io.Writer         { return vm.Settings.Stderr }
func (vm *VirtualMachine) GetStdin() io.Reader          { return vm.Settings.Stdin }

func (vm *VirtualMachine) SetInstanceVar(key string, val object.Object) {
	vm.instanceVars[key] = val
}

func (vm *VirtualMachine) GetInstanceVar(key string) object.Object {
	return vm.instanceVars[key]
}

func (vm *VirtualMachine) GetOkInstanceVar(key string) (object.Object, bool) {
	val, ok := vm.instanceVars[key]
	return val, ok
}

func (vm *VirtualMachine) HasInstanceVar(key string) bool {
	_, ok := vm.instanceVars[key]
	return ok
}

func (vm *VirtualMachine) RemoveInstanceVar(key string) {
	delete(vm.instanceVars, key)
}

func (vm *VirtualMachine) Execute(code *compiler.CodeBlock, env *object.Environment, modulename string) (object.Object, error) {
	if env == nil {
		env = object.NewEnvironment()
	}
	env.SetParent(vm.globalEnv)
	return vm.RunFrame(vm.MakeFrame(code, env, modulename), false), vm.returnErr
}

func (vm *VirtualMachine) CurrentFrame() *Frame {
	return vm.currentFrame
}

func (vm *VirtualMachine) Exit(code int) {
	vm.returnValue = object.MakeIntObj(int64(code))
	vm.returnErr = ErrExitCode{code}
	vm.unwind = true
}

func (vm *VirtualMachine) MakeFrame(code *compiler.CodeBlock, env *object.Environment, module string) *Frame {
	return &Frame{
		code:       code,
		stack:      make([]object.Object, code.MaxStackSize+1), // +1 to make room for a runtime exception if thrown
		blockStack: make([]block, code.MaxBlockSize),
		env:        env,
		unwind:     true,
		module:     module,
	}
}

func (vm *VirtualMachine) emptyFrame(env *object.Environment, module string) *Frame {
	code := &compiler.CodeBlock{
		Name:     module,
		Filename: module,
	}
	return vm.MakeFrame(code, env, module)
}

func (vm *VirtualMachine) RunFrame(f *Frame, immediateReturn bool) (ret object.Object) {
	defer func() {
		if r := recover(); r != nil {
			if retObj, ok := r.(object.Object); ok {
				if exc, ok := retObj.(*object.Exception); ok && exc.HasStackTrace {
					ret = retObj
					return
				}

				stackBuf := bytes.Buffer{}
				fmt.Fprintln(&stackBuf, retObj)
				fmt.Fprintln(&stackBuf, "Stack Trace:")
				frame := vm.currentFrame
				for frame != nil {
					fmt.Fprintf(&stackBuf, "\t%s: %s:%d\n", frame.code.Filename, frame.code.Name, frame.lineno())
					frame = frame.lastFrame
				}
				vm.unwind = vm.currentFrame.unwind
				exc := object.NewException(stackBuf.String())
				exc.HasStackTrace = true
				ret = exc
				vm.currentFrame = vm.currentFrame.lastFrame
			} else {
				fmt.Fprintln(vm.GetStderr(), r)
				fmt.Fprintln(vm.GetStderr(), string(debug.Stack()))

				fmt.Fprintln(vm.GetStderr(), "VM Stack Trace:")
				frame := vm.currentFrame
				for frame != nil {
					fmt.Fprintf(vm.GetStderr(), "\t%s: %s:%d\n", frame.code.Filename, frame.code.Name, frame.lineno())
					frame = frame.lastFrame
				}
				vm.unwind = true
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

		if vm.currentFrame.sp > 0 &&
			object.ObjectIs(vm.currentFrame.getFrontStack(), object.ExceptionObj) &&
			!(vm.currentFrame.getFrontStack().(*object.Exception)).Caught {
			if vm.Settings.ReturnExceptions {
				return vm.currentFrame.popStack()
			}
			vm.throw()
		}

		if vm.currentFrame.pc >= len(vm.currentFrame.code.Code) {
			panic(fmt.Sprintf("Program counter %d outside bounds of bytecode %d", vm.currentFrame.pc, len(vm.currentFrame.code.Code)-1))
		}
		code := vm.fetchOpcode()
		if vm.Settings.Debug && vm.breakpoint {
			debugPrompt(vm)
		}

		switch code {
		case opcode.Noop:
			continue mainLoop

		case opcode.Breakpoint:
			vm.breakpoint = true

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
			vm.currentFrame.pushStack(object.MakeIntObj(-l.Value))

		case opcode.UnaryNot:
			l := vm.currentFrame.popStack().(*object.Boolean)
			if l.Value {
				vm.currentFrame.pushStack(object.FalseConst)
			} else {
				vm.currentFrame.pushStack(object.TrueConst)
			}

		case opcode.Implements:
			r := vm.currentFrame.popStack()
			l := vm.currentFrame.popStack()
			res := vm.evalImplementsExpression(l, r)
			vm.currentFrame.pushStack(res)
			if object.ObjectIs(res, object.ExceptionObj) {
				vm.throw()
			}

		case opcode.LoadConst:
			vm.currentFrame.pushStack(vm.currentFrame.code.Constants[vm.getUint16()])

		case opcode.StoreConst:
			// Ensure constant isn't redefined
			name := vm.currentFrame.code.Locals[vm.getUint16()]
			if vm.currentFrame.env.IsConst(name) {
				vm.currentFrame.pushStack(object.NewException("Redefined constant %s", name))
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
			if object.ObjectIs(vm.returnValue, object.ExceptionObj) {
				vm.throw()
			}

		case opcode.Pop:
			vm.currentFrame.popStack()

		case opcode.LoadFast:
			name := vm.currentFrame.code.Locals[vm.getUint16()]
			if val, ok := vm.currentFrame.env.GetLocal(name); ok {
				vm.currentFrame.pushStack(val)
				break
			}

			vm.currentFrame.pushStack(object.NewException("Unknown variable/constant %s", name))
			vm.throw()

		case opcode.StoreFast:
			// Ensure constant isn't redefined
			name := vm.currentFrame.code.Locals[vm.getUint16()]
			if vm.currentFrame.env.IsConstLocal(name) {
				vm.currentFrame.pushStack(object.NewException("Redefined constant %s", name))
				vm.throw()
				break
			}
			if _, exists := vm.currentFrame.env.GetLocal(name); !exists {
				vm.currentFrame.pushStack(object.NewException("Variable %s undefined", name))
				vm.throw()
				break
			}
			vm.currentFrame.env.SetLocal(name, vm.currentFrame.popStack())

		case opcode.DeleteFast:
			name := vm.currentFrame.code.Locals[vm.getUint16()]
			if vm.currentFrame.env.IsConstLocal(name) {
				vm.currentFrame.pushStack(object.NewException("Cannot delete constant %s", name))
				vm.throw()
				break
			}
			vm.currentFrame.env.UnsetLocal(name)

		case opcode.Define:
			// Ensure constant isn't redefined
			name := vm.currentFrame.code.Locals[vm.getUint16()]
			if vm.currentFrame.env.IsConst(name) {
				vm.currentFrame.pushStack(object.NewException("Redefined constant %s", name))
				vm.throw()
				break
			}
			if _, exists := vm.currentFrame.env.GetLocal(name); exists {
				vm.currentFrame.pushStack(object.NewException("Variable %s already defined", name))
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

			vm.currentFrame.pushStack(object.NewException("Global %s doesn't exist", name))
			vm.throw()

		case opcode.StoreGlobal:
			// Ensure constant isn't redefined
			name := vm.currentFrame.code.Names[vm.getUint16()]
			if vm.currentFrame.env.IsConst(name) {
				vm.currentFrame.pushStack(object.NewException("Redefined constant %s", name))
				vm.throw()
				break
			}
			p := vm.currentFrame.env.Parent()
			if p == nil {
				p = vm.currentFrame.env
			}
			if _, exists := p.Get(name); !exists {
				vm.currentFrame.pushStack(object.NewException("Global variable %s not defined", name))
				vm.throw()
				break
			}
			vm.currentFrame.env.Set(name, vm.currentFrame.popStack())

		case opcode.LoadIndex:
			index := vm.currentFrame.popStack()
			left := vm.currentFrame.popStack()
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
			this, exists := vm.currentFrame.env.GetLocal("this")
			if !exists {
				vm.CallFunction(numargs, fn, false, nil, !immediateReturn)
				break
			}
			instance, ok := this.(*VMInstance)
			if !ok {
				vm.CallFunction(numargs, fn, false, nil, !immediateReturn)
				break
			}
			vm.CallFunction(numargs, fn, false, instance, !immediateReturn)

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

			res := vm.compareObjects(l, r, op)
			vm.currentFrame.pushStack(res)
			if object.ObjectIs(res, object.ExceptionObj) {
				vm.throw()
			}

		case opcode.MakeFunction:
			fnName := vm.currentFrame.popStack().(*object.String)
			params := vm.currentFrame.popStack().(*object.Array)
			codeBlock := vm.currentFrame.popStack().(*compiler.CodeBlock)

			if codeBlock.Native {
				var (
					fn     object.Object
					exists bool
				)

				if codeBlock.ClassMethod {
					fn, exists = nativeMethods[codeBlock.Name]
					if !exists {
						ex := object.NewPanic("Native method not implemented %s", codeBlock.Name)
						vm.currentFrame.pushStack(ex)
						vm.throw()
						break
					}
				} else {
					fn, exists = nativeFn[codeBlock.Name]
					if !exists {
						ex := object.NewPanic("Native function not implemented %s", codeBlock.Name)
						vm.currentFrame.pushStack(ex)
						vm.throw()
						break
					}
				}
				vm.currentFrame.pushStack(fn)
				break
			}

			fn := &VMFunction{
				Name:       fnName.String(),
				Parameters: make([]string, len(params.Elements)),
				Body:       codeBlock,
				Env:        object.NewEnclosedEnv(vm.currentFrame.env),
			}

			for i, p := range params.Elements {
				fn.Parameters[i] = p.(*object.String).String()
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
					ex := object.NewException("Map key %s not valid", key.Inspect())
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

		case opcode.StartBlock:
			vm.currentFrame.pushBlock(&doBlock{})
			vm.currentFrame.env = object.NewEnclosedEnv(vm.currentFrame.env)

		case opcode.Recover:
			catch := vm.getUint16()
			tcb := &recoverBlock{
				pc: int(catch),
				sp: vm.currentFrame.sp,
			}
			vm.currentFrame.pushBlock(tcb)
			vm.currentFrame.env = object.NewEnclosedEnv(vm.currentFrame.env)

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

		case opcode.EndBlock:
			vm.currentFrame.popBlock()
			vm.currentFrame.env = vm.currentFrame.env.Parent()
			if vm.currentFrame.sp == 0 {
				vm.currentFrame.pushStack(object.NullConst)
			}

		case opcode.Continue:
			vm.currentFrame.pc = vm.currentFrame.popBlockUntil(loopBlockT).(*forLoopBlock).iter

		case opcode.NextIter:
			vm.currentFrame.pc = vm.currentFrame.popBlockUntil(loopBlockT).(*forLoopBlock).start
			vm.currentFrame.env = object.NewEnclosedEnv(vm.currentFrame.env.Parent())

		case opcode.Break:
			vm.currentFrame.pc = vm.currentFrame.popBlockUntil(loopBlockT).(*forLoopBlock).end

		case opcode.Import:
			path := vm.currentFrame.code.Constants[vm.getUint16()].(*object.String)
			vm.importPackage(path.String())

		case opcode.BuildClass:
			methodNum := vm.getUint16()
			class := &VMClass{}
			class.Name = vm.currentFrame.module + "." + vm.currentFrame.popStack().(*object.String).String()
			parent := vm.currentFrame.popStack()
			if parent != object.NullConst {
				class.Parent = parent.(*VMClass)
			}
			class.Fields = vm.currentFrame.popStack().(*compiler.CodeBlock)
			class.Methods = make(map[string]object.ClassMethod, methodNum)
			for i := methodNum; i > 0; i-- {
				method := vm.currentFrame.popStack()
				switch method := method.(type) {
				case *VMFunction:
					class.Methods[method.Name] = method
					method.Class = class
				case *BuiltinMethod:
					class.Methods[method.Name] = method
				}
			}
			vm.currentFrame.pushStack(class)

		case opcode.MakeInstance:
			argLen := vm.getUint16()
			class := vm.currentFrame.popStack()
			if !object.ObjectIs(class, object.ClassObj) {
				vm.currentFrame.pushStack(object.NewException("Cannot make instance from non-class object " + class.Type().String()))
				vm.throw()
			}
			vm.makeInstance(argLen, class)

		case opcode.LoadAttribute:
			name := vm.currentFrame.code.Names[vm.getUint16()]
			obj := vm.currentFrame.popStack()

			switch obj := obj.(type) {
			case *VMInstance:
				if method := obj.GetBoundMethod(name); method != nil {
					vm.currentFrame.pushStack(method)
				} else {
					val, ok := obj.Fields.Get(name)
					if ok {
						vm.currentFrame.pushStack(val)
					} else {
						vm.currentFrame.pushStack(object.NewException("Attribute " + name + " not found on object " + obj.Class.Name))
						vm.throw()
						break
					}
				}
			case *VMClass:
				this, exists := vm.currentFrame.env.GetLocal("this")
				if !exists {
					vm.currentFrame.pushStack(object.NewException("Method call outside instance"))
					vm.throw()
					break
				}

				instance, ok := this.(*VMInstance)
				if !ok {
					vm.currentFrame.pushStack(object.NewException("Method call outside instance"))
					vm.throw()
					break
				}

				method := obj.GetMethod(name)
				if method != nil {
					vm.currentFrame.pushStack(&BoundMethod{
						Method:   method,
						Instance: instance,
						Parent:   instance.Class.Parent,
					})
				} else {
					vm.currentFrame.pushStack(object.NewException("Attribute " + name + " not found on object " + obj.Name))
					vm.throw()
					break
				}
			case *object.Module:
				vm.currentFrame.pushStack(vm.lookupModuleAttr(obj, name))
			case *object.Hash:
				vm.currentFrame.pushStack(vm.lookupHashIndex(obj, object.MakeStringObj(name)))
			default:
				vm.currentFrame.pushStack(object.NewException("Attribute lookup on non-object type %s", obj.Type()))
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
				vm.currentFrame.pushStack(object.NewException("Attribute lookup on non-object type %s", instance.Type()))
				vm.throw()
			}

		case opcode.Dup:
			vm.currentFrame.pushStack(vm.currentFrame.getFrontStack())

		case opcode.GetIter:
			obj := vm.currentFrame.popStack()

			switch obj := obj.(type) {
			case *VMInstance:
				method := obj.GetBoundMethod("_iter")
				if method == nil {
					vm.currentFrame.pushStack(object.NewException("Instance does not implement _iter() %s", obj.Class.Name))
					vm.throw()
					break
				}
				vm.CallFunction(0, method, true, obj, !immediateReturn)
			case *object.Hash:
				vm.currentFrame.pushStack(makeMapIter(obj))
			case *object.Array:
				vm.currentFrame.pushStack(makeArrayIter(obj))
			case *object.String:
				vm.currentFrame.pushStack(makeStringIter(obj))
			case *object.ByteString:
				vm.currentFrame.pushStack(makeByteStringIter(obj))
			default:
				vm.currentFrame.pushStack(object.NewException("Attribute lookup on non-object type %s", obj.Type()))
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
		}
	}
}

func (vm *VirtualMachine) currentOpcode() opcode.Opcode {
	b := vm.currentFrame.code.Code[vm.currentFrame.pc-1]
	return opcode.Opcode(b)
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
			tryBlockS := catchBlock.(*recoverBlock)
			if !tryBlockS.caught {
				tryBlockS.caught = true
				vm.currentFrame.sp = tryBlockS.sp // Unwind data stack
				vm.currentFrame.pc = tryBlockS.pc // Set program counter to catch block
				(exception.(*object.Exception)).Caught = true
				break
			}
		}
		if !vm.currentFrame.unwind {
			exc := object.NewException(exception.Inspect())
			exc.HasStackTrace = exception.(*object.Exception).HasStackTrace
			panic(exc)
		}
		vm.currentFrame = vm.currentFrame.lastFrame // This frame doesn't have a try block, unwind call stack
		if vm.currentFrame == nil {                 // Call stack exhausted
			// exc := object.NewException("Uncaught Exception: %s", exception.Inspect())
			exc := object.NewException(exception.Inspect())
			exc.HasStackTrace = exception.(*object.Exception).HasStackTrace

			if vm.Settings.ReturnExceptions {
				return exc
			}

			vm.currentFrame = cframe // Reset frame for stack trace
			panic(exc)
		}
	}

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

		iFields := object.NewLocalEnvironment()
		iFields.SetParent(vm.currentFrame.env)

		for _, c := range classChain {
			frame := vm.MakeFrame(c.Fields, iFields, vm.currentFrame.module)
			ret := vm.RunFrame(frame, true)
			if ret.Type() == object.ExceptionObj {
				vm.currentFrame.pushStack(ret)
				vm.throw()
				return
			}
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

	if instance == nil {
		panic("Cannot make an instance from " + class.Type().String())
	}

	init := instance.GetBoundMethod("init")
	if init == nil {
		for i := argLen; i > 0; i-- {
			vm.currentFrame.popStack()
		}
		vm.currentFrame.pushStack(instance)
		return
	}

	vm.CallFunction(argLen, init, true, nil, false)
	ret := vm.currentFrame.popStack() // Pop return value of init function
	if ret.Type() == object.ExceptionObj {
		vm.currentFrame.pushStack(ret)
		vm.throw()
		return
	}

	if ret.Type() == object.ErrorObj {
		vm.currentFrame.pushStack(ret)
		return
	}

	vm.currentFrame.pushStack(instance)
}

func (vm *VirtualMachine) CallFunction(argc uint16, fn object.Object, now bool, this *VMInstance, unwind bool) {
	switch fn := fn.(type) {
	case *object.Builtin:
		if vm.Settings.Debug {
			fmt.Fprintf(vm.GetStdout(), "Calling builtin\n")
		}

		args := make([]object.Object, argc)
		for i := uint16(0); i < argc; i++ {
			args[i] = vm.currentFrame.popStack()
		}

		env := vm.currentFrame.env
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
		if vm.Settings.Debug {
			fmt.Fprintf(vm.GetStdout(), "Calling builtin method %s\n", fn.Name)
		}

		args := make([]object.Object, argc)
		for i := uint16(0); i < argc; i++ {
			args[i] = vm.currentFrame.popStack()
		}

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
		if vm.Settings.Debug {
			fmt.Fprintf(vm.GetStdout(), "Calling function %s\n", fn.Name)
		}

		var env *object.Environment
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

		newFrame := vm.MakeFrame(fn.Body, env, vm.currentFrame.module)
		newFrame.unwind = unwind
		newFrame.lastFrame = vm.currentFrame

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
			val := vm.RunFrame(newFrame, true)
			vm.currentFrame.pushStack(val)
		} else {
			vm.currentFrame = newFrame
			vm.callStack.Push(newFrame)
		}
	case *BoundMethod:
		vm.CallFunction(argc, fn.Method, true, fn.Instance, unwind)
	case *VMClass:
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
		vm.CallFunction(argc, init, true, this, unwind)
	default:
		for i := 0; i < int(argc); i++ {
			vm.currentFrame.popStack()
		}
		vm.currentFrame.pushStack(object.NewPanic("%s is not a function", fn.Type()))
		vm.throw()
	}
}

func (vm *VirtualMachine) ImportPreamble(name string) error {
	if name == "" {
		name = "std/preamble/main"
	}

	vm.currentFrame = vm.emptyFrame(vm.globalEnv, "__preamble__")

	vm.importPackage(name)
	module, ok := vm.PopStack().(*object.Hash)
	if !ok {
		return errors.New("preamble module did not return a hash")
	}

	for _, pair := range module.Pairs {
		vm.globalEnv.SetForce(pair.Key.Inspect(), pair.Value, true)
	}

	vm.currentFrame = nil
	return nil
}
