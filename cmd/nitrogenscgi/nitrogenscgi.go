package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/nitrogen-lang/nitrogen/src/elemental/object"
	"github.com/nitrogen-lang/nitrogen/src/elemental/vm"
	"github.com/nitrogen-lang/nitrogen/src/moduleutils"
	"github.com/nitrogen-lang/nitrogen/src/scgi"

	_ "github.com/nitrogen-lang/nitrogen/src/builtins"
)

type strSliceFlag []string

func (s *strSliceFlag) String() string {
	return strings.Join(*s, ":")
}

func (s *strSliceFlag) Set(st string) error {
	*s = append(*s, st)
	return nil
}

var (
	printVersion bool
	disableNibs  bool

	extraModulePaths strSliceFlag
	autoloadModules  strSliceFlag
	modulePaths      []string

	version         = "Unknown"
	buildTime       = ""
	builder         = ""
	builtinModPaths = ""
)

func init() {
	flag.BoolVar(&disableNibs, "nonibs", false, "Disable creation of .nib files")
	flag.BoolVar(&printVersion, "version", false, "Print version information")

	flag.Var(&extraModulePaths, "M", "Module search paths")
	flag.Var(&autoloadModules, "al", "Autoload modules")
}

func main() {
	flag.Parse()

	modulePaths = make([]string, 0, len(extraModulePaths)+6)

	// Package paths from command line flag
	modulePaths = append(modulePaths, extraModulePaths...)

	// Package paths from environment variable
	envModPath := os.Getenv("NITROGEN_MODULES")
	if envModPath != "" {
		modulePaths = append(modulePaths, strings.Split(envModPath, ":")...)
	}

	// Add working directory to path
	pwd, _ := os.Getwd()
	modulePaths = append(modulePaths, pwd)

	// Add compile time paths
	if builtinModPaths != "" {
		modulePaths = append(modulePaths, strings.Split(builtinModPaths, ":")...)
	}

	// Add Noble package manager path
	homeDir, _ := os.UserHomeDir()
	if homeDir != "" {
		modulePaths = append(modulePaths, filepath.Join(homeDir, ".noble", "pkgs"))
	}

	if printVersion {
		versionInfo()
		return
	}

	if disableNibs {
		moduleutils.WriteCompiledScripts = false
	}

	if len(autoloadModules) > 0 {
		if err := vm.PreloadModules(modulePaths, autoloadModules); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	scgi.StartSCGIServer(getScriptArgs("nitrogen"), object.MakeStringArray(modulePaths), getServerEnv())
}

func getServerEnv() *object.Hash {
	if os.Getenv("GATEWAY_INTERFACE") != "CGI/1.1" {
		return object.MakeEmptyHash()
	}

	headers := make(map[string]string, len(scgi.CGIHeaderNames))
	for _, header := range scgi.CGIHeaderNames {
		headers[header] = os.Getenv(header)
	}

	return object.StringMapToHash(headers)
}

func getScriptArgs(filepath string) *object.Array {
	var s []string
	if flag.NArg() > 1 {
		s = flag.Args()[1:]
	}
	length := len(s) + 1
	newElements := make([]object.Object, length)
	newElements[0] = object.MakeStringObj(filepath)
	for i, v := range s {
		newElements[i+1] = object.MakeStringObj(v)
	}
	return &object.Array{Elements: newElements}
}

func versionInfo() {
	fmt.Printf(`Nitrogen
Version:           %s
Built:             %s
Compiled by:       %s
Go version:        %s %s/%s
Builtin Mod Path:  %s
Pkg Path:          %s
Native Modules Supported: %t
`, version, buildTime, builder, runtime.Version(), runtime.GOOS, runtime.GOARCH, builtinModPaths, modulePaths, vm.ModulesSupported)
}
