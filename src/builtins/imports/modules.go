// +build linux,cgo darwin,cgo

package imports

import (
	"plugin"

	"github.com/nitrogen-lang/nitrogen/src/object"
	"github.com/nitrogen-lang/nitrogen/src/vm"
)

func moduleSupport(i object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	return object.NativeBoolToBooleanObj(true)
}

func importSharedModule(scriptPath string, required bool, interpreter object.Interpreter, env *object.Environment) object.Object {
	p, err := plugin.Open(scriptPath)
	if err != nil {
		if required {
			return object.NewException("%s", err)
		}
		return object.NewError("%s", err)
	}

	// Check module name
	moduleNameSym, err := p.Lookup("ModuleName")
	if err != nil {
		// The module didn't declare a name
		return object.NewException("Invalid module %s", scriptPath)
	}

	if module := vm.GetModule(*(moduleNameSym.(*string))); module != nil {
		return module
	}
	return object.NullConst
}
