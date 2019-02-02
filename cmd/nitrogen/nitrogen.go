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
	"strconv"
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
	interactive  bool
	printAst     bool
	startSCGI    bool
	printVersion bool
	fullDebug    bool
	compileOnly  bool
	cpuprofile   string
	memprofile   string
	outputFile   string

	infoCmd bool

	modulePaths     strSliceFlag
	autoloadModules strSliceFlag

	version         = "Unknown"
	buildTime       = ""
	builder         = ""
	builtinModPaths = ""
)

func init() {
	pwd, _ := os.Getwd()
	modulePaths = append(modulePaths, pwd)

	envModPath := os.Getenv("NITROGEN_MODULES")
	if envModPath != "" {
		modulePaths = append(modulePaths, strings.Split(envModPath, ":")...)
	}

	flag.BoolVar(&interactive, "i", false, "Interactive mode")
	flag.BoolVar(&compileOnly, "c", false, "Compile code, print any errors, and exit")
	flag.BoolVar(&printAst, "ast", false, "Print AST and exit")
	flag.BoolVar(&startSCGI, "scgi", false, "Start as an SCGI server")
	flag.BoolVar(&printVersion, "version", false, "Print version information")
	flag.BoolVar(&fullDebug, "debug", false, "Enable debug mode")
	flag.StringVar(&cpuprofile, "cpuprofile", "", "File to write CPU profile data")
	flag.StringVar(&memprofile, "memprofile", "", "File to write memory profile data")
	flag.StringVar(&outputFile, "o", "", "Output file of compiled bytecode")

	flag.Var(&modulePaths, "M", "Module search paths")
	flag.Var(&autoloadModules, "al", "Autoload modules")

	flag.BoolVar(&infoCmd, "info", false, "Print information about a .nib file")
}

func main() {
	flag.Parse()

	if builtinModPaths != "" {
		modulePaths = append(modulePaths, strings.Split(builtinModPaths, ":")...)
	}

	if printVersion {
		versionInfo()
		return
	}

	if infoCmd {
		runInfoCmd()
		return
	}

	if len(autoloadModules) > 0 {
		if err := loadModules(modulePaths, autoloadModules); err != nil {
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
		fmt.Println("Nitrogen Programming Language")
		fmt.Println("Type in commands at the prompt")
		startRepl(os.Stdin, os.Stdout)
		return
	}

	if flag.NArg() == 0 {
		fmt.Println("No script given")
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

	sourceFile := flag.Arg(0)

	env := makeEnv(sourceFile)

	var code *compiler.CodeBlock
	var program *ast.Program
	var err error
	if filepath.Ext(sourceFile) == ".nib" {
		code, _, err = marshal.ReadFile(sourceFile)
		if err != nil {
			fmt.Print("There were errors reading compiled program:\n\n")
			fmt.Println(err.Error())
			return
		}
	} else {
		program, err = moduleutils.ASTCache.GetTree(sourceFile)
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
	}

	var result object.Object
	var start time.Time

	if code == nil {
		code = compiler.Compile(program, "__main")
	}

	if compileOnly {
		return
	}

	if outputFile != "" {
		marshal.WriteFile(outputFile, code, moduleutils.FileModTime(sourceFile))
		return
	}

	start = time.Now()
	result = runCompiledCode(code, env)

	if fullDebug {
		fmt.Printf("Execution took %s\n", time.Now().Sub(start))
	}

	if result != nil && result != object.NullConst {
		if e, ok := result.(*object.Exception); ok {
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
	machine.SetGlobalEnv(env)

	ret, err := machine.Execute(code, nil)
	if err != nil {
		if ex, ok := err.(vm.ErrExitCode); ok {
			os.Exit(ex.Code)
		}
	}
	return ret
}

func makeEnv(filepath string) *object.Environment {
	env := object.NewEnvironment()
	env.CreateConst("_ENV", getExternalEnv())
	env.CreateConst("_ARGV", getScriptArgs(filepath))
	env.CreateConst("_SERVER", getServerEnv())
	env.Create("_SEARCH_PATHS", object.MakeStringArray(modulePaths))
	return env
}

var cgiHeaderNames = []string{
	"AUTH_TYPE",
	"DOCUMENT_ROOT",
	"DOCUMENT_URI",
	"GATEWAY_INTERFACE",
	"HTTP_ACCEPT_CHARSET",
	"HTTP_ACCEPT_ENCODING",
	"HTTP_ACCEPT_LANGUAGE",
	"HTTP_ACCEPT",
	"HTTP_CONNECTION",
	"HTTP_HOST",
	"HTTP_REFERER",
	"HTTP_USER_AGENT",
	"HTTPS",
	"QUERY_STRING",
	"REDIRECT_REMOTE_USER",
	"REMOTE_ADDR",
	"REMOTE_HOST",
	"REMOTE_PORT",
	"REMOTE_USER",
	"REQUEST_METHOD",
	"REQUEST_TIME",
	"REQUEST_URI",
	"SCRIPT_FILENAME",
	"SCRIPT_NAME",
	"SERVER_ADDR",
	"SERVER_ADMIN",
	"SERVER_NAME",
	"SERVER_PORT",
	"SERVER_PROTOCOL",
	"SERVER_SIGNATURE",
	"SERVER_SOFTWARE",
}

func getServerEnv() object.Object {
	if os.Getenv("GATEWAY_INTERFACE") != "CGI/1.1" {
		return object.MakeEmptyHash()
	}

	headers := make(map[string]string, len(cgiHeaderNames))
	for _, header := range cgiHeaderNames {
		headers[header] = os.Getenv(header)
	}

	return object.StringMapToHash(headers)
}

func getExternalEnv() *object.Hash {
	return object.StringMapToHash(getExtEnvMap())
}

func getExtEnvMap() map[string]string {
	env := os.Environ()
	m := make(map[string]string, len(env))
	for _, v := range env {
		val := strings.SplitN(v, "=", 2)
		m[val[0]] = val[1]
	}
	return m
}

func getScriptArgs(filepath string) *object.Array {
	var s []string
	if flag.NArg() > 1 {
		s = flag.Args()[1:]
	}
	length := len(s) + 1
	newElements := make([]object.Object, length, length)
	newElements[0] = object.MakeStringObj(filepath)
	for i, v := range s {
		newElements[i+1] = object.MakeStringObj(v)
	}
	return &object.Array{Elements: newElements}
}

func startRepl(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := makeEnv("")

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

func versionInfo() {
	fmt.Printf(`Nitrogen - (C) 2018 Lee Keitel
Version:           %s
Built:             %s
Compiled by:       %s
Go version:        %s %s/%s
Modules Supported: %t
Builtin Mod Path:  %s
`, version, buildTime, builder, runtime.Version(), runtime.GOOS, runtime.GOARCH, modulesSupported, builtinModPaths)
}

func runInfoCmd() {
	sourceFile := flag.Arg(0)
	if sourceFile == "" {
		return
	}

	code, fileinfo, err := marshal.ReadFile(sourceFile)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("Filename: %s\n", fileinfo.Filename)
	fmt.Printf("Version:  %s\n", bytesToVersionNumber(fileinfo.Version))
	fmt.Printf("ModTime:  %s\n", fileinfo.ModTime)
	code.Print("")
}

func bytesToVersionNumber(b []byte) string {
	ver := ""
	for _, v := range b {
		ver += strconv.Itoa(int(v)) + "."
	}
	return strings.TrimRight(ver, ".")
}
