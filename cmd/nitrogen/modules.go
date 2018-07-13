// +build linux,cgo darwin,cgo

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"plugin"
)

const modulesSupported = true

func loadModules(searchPaths, modules []string) error {
	for _, module := range modules {
		for _, path := range searchPaths {
			fullpath := filepath.Join(path, module)

			if fileExists(fullpath) {
				if fullDebug {
					fmt.Printf("Loading module %s\n", filepath.Base(fullpath))
				}
				_, err := plugin.Open(fullpath)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func fileExists(file string) bool {
	_, err := os.Stat(file)
	return !os.IsNotExist(err)
}
