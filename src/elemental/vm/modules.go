//go:build (linux && cgo) || (darwin && cgo)

package vm

import (
	"os"
	"path/filepath"
	"plugin"

	"github.com/nitrogen-lang/nitrogen/src/elemental/object"
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

const ModulesSupported = true

func PreloadModules(searchPaths, modules []string) error {
	for _, module := range modules {
		for _, path := range searchPaths {
			fullpath := filepath.Join(path, module)

			if fileExists(fullpath) {
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
