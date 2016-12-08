package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/nitrogen-lang/nitrogen/src/eval"
	"github.com/nitrogen-lang/nitrogen/src/lexer"
	"github.com/nitrogen-lang/nitrogen/src/object"
	"github.com/nitrogen-lang/nitrogen/src/parser"
)

const PROMPT = ">> "

var (
	interactive bool
	scriptFile  string
)

func init() {
	flag.StringVar(&scriptFile, "f", "", "Filename to execute")
	flag.BoolVar(&interactive, "i", false, "Interactive mode")
}

func main() {
	flag.Parse()

	if interactive {
		fmt.Print("Nitrogen Programming Language\n")
		fmt.Print("Type in commands at the prompt\n")
		startRepl(os.Stdin, os.Stdout)
		return
	}

	if scriptFile == "" {
		fmt.Print("No file given")
		os.Exit(1)
	}

	file, err := os.Open(scriptFile)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	l := lexer.New(file)
	p := parser.New(l)
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		printParserErrors(os.Stdout, p.Errors())
		os.Exit(1)
	}

	result := eval.Eval(program, object.NewEnvironment())
	if result != nil && result != eval.NULL {
		io.WriteString(os.Stdout, result.Inspect())
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
		fmt.Fprintf(out, "\t%s\n", msg)
	}
}
