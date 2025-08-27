//go:build (linux && cgo) || (darwin && cgo)
// +build linux,cgo darwin,cgo

package imports

import (
	"github.com/nitrogen-lang/nitrogen/src/elemental/object"
)

func moduleSupport(i object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	return object.NativeBoolToBooleanObj(true)
}
