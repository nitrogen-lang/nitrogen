package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/nitrogen-lang/nitrogen/src/compiler"
	"github.com/nitrogen-lang/nitrogen/src/compiler/marshal"
	"github.com/nitrogen-lang/nitrogen/src/elemental/vm"
	"github.com/nitrogen-lang/nitrogen/src/moduleutils"
)

var (
	printAst     bool
	printVersion bool
	fullDebug    bool
	outputFile   string

	version   = "Unknown"
	buildTime = ""
	builder   = ""
)

func init() {
	flag.BoolVar(&printAst, "ast", false, "Print AST and exit")
	flag.BoolVar(&printVersion, "version", false, "Print version information")
	flag.BoolVar(&fullDebug, "debug", false, "Enable debug mode")
	flag.StringVar(&outputFile, "o", "", "Output file of compiled bytecode")
}

func main() {
	flag.Parse()

	if printVersion {
		versionInfo()
		return
	}

	if flag.NArg() == 0 {
		fmt.Println("No script given")
		os.Exit(1)
	}

	sourceFile := flag.Arg(0)

	if filepath.Ext(sourceFile) == ".nib" {
		fmt.Print("The file is already Nitrogen bytecode")
		return
	}

	program, err := moduleutils.ASTCache.GetTree(sourceFile)
	if err != nil {
		fmt.Print("There were errors compiling the program:\n\n")
		fmt.Println(err.Error())
		os.Exit(1)
		return
	}

	if printAst {
		fmt.Println(program.String())
		return
	}

	code := compiler.Compile(program, "__main")

	if outputFile == "" {
		sourceFileDir := filepath.Dir(sourceFile)
		sourceFileBase := filepath.Base(sourceFile)
		sourceFilename := sourceFileBase[:strings.LastIndexByte(sourceFileBase, '.')]
		outputFile = filepath.Join(sourceFileDir, sourceFilename+".nib")
	}
	fmt.Printf("Bytecode writen to %s\n", outputFile)
	marshal.WriteFile(outputFile, code, moduleutils.FileModTime(sourceFile), true)
}

func versionInfo() {
	fmt.Printf(`Nitrogen - (C) 2018 Lee Keitel
Version:           %s
Built:             %s
Compiled by:       %s
Go version:        %s %s/%s
Modules Supported: %t
`, version, buildTime, builder, runtime.Version(), runtime.GOOS, runtime.GOARCH, vm.ModulesSupported)
}
