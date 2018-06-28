// +build linux,cgo darwin,cgo

package vm

import (
	"plugin"

	"github.com/nitrogen-lang/nitrogen/src/object"
)

func importSharedModule(vm *VirtualMachine, scriptPath, name string) object.Object {
	p, err := plugin.Open(scriptPath)
	if err != nil {
		return object.NewException("%s", err)
	}

	// Check module name
	moduleNameSym, err := p.Lookup("ModuleName")
	if err != nil {
		// The module didn't declare a name
		return object.NewException("Invalid module %s, no name declared", name)
	}

	if module := GetModule(*(moduleNameSym.(*string))); module != nil {
		return module
	}
	return object.NullConst
}
