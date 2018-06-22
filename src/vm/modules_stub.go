// +build !linux,!darwin !cgo

package vm

import (
	"github.com/nitrogen-lang/nitrogen/src/object"
)

func importSharedModule(vm *VirtualMachine, scriptPath string) object.Object {
	return object.NewException("Shared object modules are not supported in this build")
}
