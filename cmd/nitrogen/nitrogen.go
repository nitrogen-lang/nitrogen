package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"plugin"

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

	if modulePath != "" {
		if err := loadModules(modulePath); err != nil {
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

	if len(flag.Args()) == 0 {
		fmt.Print("No file given")
		os.Exit(1)
	}

	l, err := lexer.NewFileList(flag.Args())
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

	result := eval.Eval(program, object.NewEnvironment())
	if result != nil && result != object.NULL {
		os.Stdout.WriteString(result.Inspect())
		os.Stdout.WriteString("\n")

		if _, ok := result.(*object.Exception); ok {
			os.Exit(1)
		}
	}
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
