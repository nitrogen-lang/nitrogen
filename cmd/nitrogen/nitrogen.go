package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"plugin"
	"strings"

	"github.com/nitrogen-lang/nitrogen/src/eval"
	"github.com/nitrogen-lang/nitrogen/src/lexer"
	"github.com/nitrogen-lang/nitrogen/src/object"
	"github.com/nitrogen-lang/nitrogen/src/parser"

	_ "github.com/nitrogen-lang/nitrogen/src/builtins"
)

const PROMPT = ">> "

var (
	interactive bool
	printAst    bool
	modulePath  string
)

func init() {
	flag.BoolVar(&interactive, "i", false, "Interactive mode")
	flag.BoolVar(&printAst, "ast", false, "Print AST and exit")
	flag.StringVar(&modulePath, "modules", "", "Module directory")
}

func main() {
	flag.Parse()

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

	l, err := lexer.NewFile(flag.Arg(0))
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	p := parser.New(l)
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		printParserErrors(os.Stdout, p.Errors())
		os.Exit(1)
	}

	if printAst {
		fmt.Println(program.String())
		os.Exit(1)
	}

	env := object.NewEnvironment()
	env.CreateConst("_ENV", getEnvironment())

	env.CreateConst("_ARGV", getScriptArgs(flag.Arg(0)))

	result := eval.Eval(program, env)
	if result != nil && result != object.NullConst {
		os.Stdout.WriteString(result.Inspect())
		os.Stdout.WriteString("\n")

		if _, ok := result.(*object.Exception); ok {
			os.Exit(1)
		}
	}
}

func getEnvironment() *object.Hash {
	env := os.Environ()
	m := &object.Hash{Pairs: make(map[object.HashKey]object.HashPair)}
	for _, v := range env {
		val := strings.SplitN(v, "=", 2)
		key := &object.String{Value: val[0]}
		m.Pairs[key.HashKey()] = object.HashPair{
			Key:   key,
			Value: &object.String{Value: val[1]},
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

	for {
		fmt.Fprint(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		if line == ".quit" {
			return
		}

		l := lexer.NewString(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		result := eval.Eval(program, env)
		if result != nil {
			io.WriteString(out, result.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		fmt.Fprintf(out, "ERROR: %s\n", msg)
	}
}

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
