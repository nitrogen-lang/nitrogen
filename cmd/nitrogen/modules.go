// +build linux,cgo darwin,cgo

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"plugin"
)

func loadModules(modulePath string) error {
	return filepath.Walk(modulePath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() || filepath.Ext(path) != ".so" {
			return nil
		}
		if err != nil {
			return err
		}

		fmt.Printf("Loading module %s\n", filepath.Base(path))
		_, err = plugin.Open(path)
		return err
	})
}
