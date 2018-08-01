package dis

import (
	"fmt"
	"runtime"

	"github.com/nitrogen-lang/nitrogen/src/moduleutils"
	"github.com/nitrogen-lang/nitrogen/src/object"
	"github.com/nitrogen-lang/nitrogen/src/vm"
)

var moduleName = "stdlib/runtime"

func init() {
	vm.RegisterModule(moduleName, &object.Module{
		Name: moduleName,
		Methods: map[string]object.BuiltinFunction{
			"dis": disassemble,
		},
		Vars: map[string]object.Object{
			"osName": object.MakeStringObj(runtime.GOOS),
			"osArch": object.MakeStringObj(runtime.GOARCH),
		},
	})
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
