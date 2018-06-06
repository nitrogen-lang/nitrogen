package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"strings"
	"time"

	"github.com/nitrogen-lang/nitrogen/src/ast"
	"github.com/nitrogen-lang/nitrogen/src/compiler"
	"github.com/nitrogen-lang/nitrogen/src/compiler/marshal"
	"github.com/nitrogen-lang/nitrogen/src/lexer"
	"github.com/nitrogen-lang/nitrogen/src/moduleutils"
	"github.com/nitrogen-lang/nitrogen/src/object"
	"github.com/nitrogen-lang/nitrogen/src/parser"
	"github.com/nitrogen-lang/nitrogen/src/vm"

	_ "github.com/nitrogen-lang/nitrogen/src/builtins"
)

const (
	interactivePrompt = ">> "
	version           = "0.1.0"
)

var (
	interactive       bool
	printAst          bool
	startSCGI         bool
	scgiSock          string
	scgiWorkers       int
	scgiWorkerTimeout int
	modulePath        string
	printVersion      bool
	fullDebug         bool
	cpuprofile        string
	memprofile        string
	outputFile        string
)

func init() {
	flag.BoolVar(&interactive, "i", false, "Interactive mode")
	flag.BoolVar(&printAst, "ast", false, "Print AST and exit")
	flag.BoolVar(&startSCGI, "scgi", false, "Start as an SCGI server")
	flag.StringVar(&scgiSock, "scgi-sock", "tcp:0.0.0.0:9000", "Socket to listen on for SCGI")
	flag.IntVar(&scgiWorkers, "scgi-workers", 5, "Number of workers to service SCGI requests")
	flag.IntVar(&scgiWorkerTimeout, "scgi-worker-timeout", 10, "Number of seconds to wait for an available worker before giving up")
	flag.StringVar(&modulePath, "modules", "", "Module directory")
	flag.BoolVar(&printVersion, "version", false, "Print version information")
	flag.BoolVar(&fullDebug, "debug", false, "Enable debug mode")
	flag.StringVar(&cpuprofile, "cpuprofile", "", "File to write CPU profile data")
	flag.StringVar(&memprofile, "memprofile", "", "File to write memory profile data")
	flag.StringVar(&outputFile, "o", "", "Output file of compiled bytecode")
}

func main() {
	flag.Parse()

	if printVersion {
		fmt.Printf("Nitrogen Interpreter Version %s\n", version)
		return
	}

	modulesPath := os.Getenv("NITROGEN_MODULES")
	if modulePath != "" {
		modulesPath = modulePath
	}
	if modulesPath != "" {
		if err := loadModules(modulesPath); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	if startSCGI {
		startSCGIServer()
		return
	}

	moduleutils.ParserSettings.Debug = fullDebug
	if interactive {
		fmt.Print("Nitrogen Programming Language\n")
		fmt.Print("Type in commands at the prompt\n")
		startRepl(os.Stdin, os.Stdout)
		return
	}

	if flag.NArg() == 0 {
		fmt.Print("No script given")
		os.Exit(1)
	}

	if cpuprofile != "" {
		f, err := os.Create(cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	env := object.NewEnvironment()
	env.CreateConst("_ENV", getEnvironment())
	env.CreateConst("_ARGV", getScriptArgs(flag.Arg(0)))

	var code *compiler.CodeBlock
	var program *ast.Program
	var err error
	if filepath.Ext(flag.Arg(0)) == ".nib" {
		code, err = marshal.ReadFile(flag.Arg(0))
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	} else {
		program, err = moduleutils.ASTCache.GetTree(flag.Arg(0))
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		if printAst {
			fmt.Println(program.String())
			return
		}
	}

	var result object.Object
	var start time.Time

	if code == nil {
		code = compiler.Compile(program, "__main")
	}

	if outputFile != "" {
		marshal.WriteFile(outputFile, code)
		return
	}

	start = time.Now()
	result = runCompiledCode(code, env)

	if fullDebug {
		fmt.Printf("Execution took %s\n", time.Now().Sub(start))
	}

	if result != nil && result != object.NullConst {
		if e, ok := result.(*object.Exception); ok {
			os.Stdout.WriteString("Uncaught Exception: ")
			os.Stdout.WriteString(e.Message)
			os.Stdout.Write([]byte{'\n'})
			os.Exit(1)
		}
		os.Stdout.WriteString(result.Inspect())
		os.Stdout.Write([]byte{'\n'})
	}

	if memprofile != "" {
		f, err := os.Create(memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
		f.Close()
	}
}

func runCompiledCode(code *compiler.CodeBlock, env *object.Environment) object.Object {
	if fullDebug {
		code.Print("")
	}

	env.CreateConst("_FILE", object.MakeStringObj(code.Filename))

	vmsettings := vm.NewSettings()
	vmsettings.Debug = fullDebug
	machine := vm.NewVM(vmsettings)

	return machine.Execute(code, env)
}

func getEnvironment() *object.Hash {
	return makeEnvironment(getEnvironmentMap())
}

func getEnvironmentMap() map[string]string {
	env := os.Environ()
	m := make(map[string]string, len(env))
	for _, v := range env {
		val := strings.SplitN(v, "=", 2)
		m[val[0]] = val[1]
	}
	return m
}

func makeEnvironment(env map[string]string) *object.Hash {
	m := &object.Hash{Pairs: make(map[object.HashKey]object.HashPair)}
	for k, v := range env {
		key := &object.String{Value: k}
		m.Pairs[key.HashKey()] = object.HashPair{
			Key:   key,
			Value: &object.String{Value: v},
		}
	}
	return m
}

func getScriptArgs(filepath string) *object.Array {
	s := flag.Args()[1:]
	length := len(s) + 1
	newElements := make([]object.Object, length, length)
	newElements[0] = &object.String{Value: filepath}
	for i, v := range s {
		newElements[i+1] = &object.String{Value: v}
	}
	return &object.Array{Elements: newElements}
}

func startRepl(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()
	env.CreateConst("_ENV", getEnvironment())

	var code *compiler.CodeBlock
	for {
		fmt.Fprint(out, interactivePrompt)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if line == ".quit" || line == ".exit" {
			return
		}

		l := lexer.NewString(line)
		p := parser.New(l, moduleutils.ParserSettings)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		code = compiler.Compile(program, "__main")

		result := runCompiledCode(code, env)
		if result != nil && result != object.NullConst {
			io.WriteString(out, result.Inspect())
		}
		io.WriteString(out, "\n")
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		fmt.Fprintf(out, "ERROR: %s\n", msg)
	}
}
