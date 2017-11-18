// +build !linux,!darwin !cgo

package main

import "errors"

func loadModules(modulePath string) error {
	return errors.New("This version of Nitrogen was built without shared module support.")
}
