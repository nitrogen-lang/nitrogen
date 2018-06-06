package dis

import (
	"fmt"

	"github.com/nitrogen-lang/nitrogen/src/object"
	"github.com/nitrogen-lang/nitrogen/src/vm"
)

func init() {
	vm.RegisterBuiltin("dis", disassemble)
}

func disassemble(machine object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	cb := args[0].(*vm.VMFunction).Body

	fmt.Printf("Name: %s\nFilename: %s\nLocalCount: %d\nMaxStackSize: %d\nMaxBlockSize: %d\n",
		cb.Name, cb.Filename, cb.LocalCount, cb.MaxStackSize, cb.MaxBlockSize)
	cb.Print(" ")
	return object.NullConst
}
