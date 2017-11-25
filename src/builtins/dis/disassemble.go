package dis

import (
	"github.com/nitrogen-lang/nitrogen/src/object"
	"github.com/nitrogen-lang/nitrogen/src/vm"
)

func init() {
	vm.RegisterBuiltin("dis", disassemble)
}

func disassemble(machine object.Interpreter, args ...object.Object) object.Object {
	args[0].(*vm.VMFunction).Body.Print()
	return object.NullConst
}
