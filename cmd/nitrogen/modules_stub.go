// +build !linux,!darwin !cgo

package main

import "fmt"

func loadModules(searchPaths []string, modules []string) error {
	if fullDebug {
		fmt.Println("This version of Nitrogen was built without shared module support.")
	}
	return nil
}
