package runtime

import (
	"fmt"
	"runtime"

	"github.com/nitrogen-lang/nitrogen/src/moduleutils"
	"github.com/nitrogen-lang/nitrogen/src/object"
	"github.com/nitrogen-lang/nitrogen/src/vm"
)

func init() {
	vm.RegisterNative("std.runtime.debugVal", debugBuiltin)
	vm.RegisterNative("std.runtime.dis", disassemble)
	vm.RegisterNative("std.runtime.dis_method", disassemble_method)
	vm.RegisterNative("std.runtime.osName", osName)
	vm.RegisterNative("std.runtime.osArch", osArch)
}

func osName(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	return object.MakeStringObj(runtime.GOOS)
}

func osArch(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	return object.MakeStringObj(runtime.GOARCH)
}

func debugBuiltin(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewException("Incorrect number of arguments. Got %d, expected 1", len(args))
	}

	fmt.Println(args[0].Inspect())

	return args[0]
}

func disassemble(machine object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	if ac := moduleutils.CheckMinArgs("dis", 1, args...); ac != nil {
		return ac
	}

	var fnObj object.Object

	if bm, ok := args[0].(*vm.BoundMethod); ok {
		fnObj = bm.Method
	} else {
		fnObj = args[0]
	}

	if cl, ok := fnObj.(*vm.VMClass); ok {
		fmt.Printf("Field Count: %d\nMethod Count: %d\n", 0, len(cl.Methods))
		return object.NullConst
	}

	if cl, ok := fnObj.(*vm.VMInstance); ok {
		fmt.Printf("%#v\n", cl.Class.Methods)
		fmt.Printf("Field Count: %d\nMethod Count: %d\n", 0, len(cl.Class.Methods))
		return object.NullConst
	}

	if iface, ok := fnObj.(*object.Interface); ok {
		fmt.Println(iface.Inspect(), "{")
		for _, def := range iface.Methods {
			fmt.Println("", def.Inspect())
		}
		fmt.Println("}")
		return object.NullConst
	}

	fn, ok := fnObj.(*vm.VMFunction)
	if !ok {
		return object.NewException("dis expected a func, got %s", fnObj.Type().String())
	}

	cb := fn.Body

	fmt.Printf("Name: %s\nFilename: %s\nLocalCount: %d\nMaxStackSize: %d\nMaxBlockSize: %d\n",
		cb.Name, cb.Filename, cb.LocalCount, cb.MaxStackSize, cb.MaxBlockSize)
	cb.Print(" ")
	return object.NullConst
}

func disassemble_method(machine object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	if ac := moduleutils.CheckMinArgs("dis_method", 2, args...); ac != nil {
		return ac
	}

	classObj, ok := args[0].(*vm.VMClass)
	if !ok {
		return object.NewException("dis_method expected first arg to be a Class, got %s", args[0].Type().String())
	}

	methodName, ok := args[1].(*object.String)
	if !ok {
		return object.NewException("dis_method expected second arg to be a String, got %s", args[1].Type().String())
	}

	method := classObj.GetMethod(string(methodName.String()))

	fn, ok := method.(*vm.VMFunction)
	if !ok {
		return object.NewException("dis_method expected a func, got %s", method.Type().String())
	}

	cb := fn.Body

	fmt.Printf("Name: %s\nFilename: %s\nLocalCount: %d\nMaxStackSize: %d\nMaxBlockSize: %d\n",
		cb.Name, cb.Filename, cb.LocalCount, cb.MaxStackSize, cb.MaxBlockSize)
	cb.Print(" ")
	return object.NullConst
}
