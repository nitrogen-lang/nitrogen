//go:build (!linux && !darwin) || !cgo

package vm

import (
	"fmt"

	"github.com/nitrogen-lang/nitrogen/src/object"
)

func importSharedModule(vm *VirtualMachine, scriptPath, name string) object.Object {
	return object.NewException("Shared object modules are not supported in this build")
}

const ModulesSupported = false

func PreloadModules(searchPaths []string, modules []string) error {
	fmt.Println("This version of Nitrogen was built without shared module support.")
	return nil
}
