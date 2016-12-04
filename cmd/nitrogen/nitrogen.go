package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/lfkeitel/nitrogen/src/eval"
	"github.com/lfkeitel/nitrogen/src/lexer"
	"github.com/lfkeitel/nitrogen/src/object"
	"github.com/lfkeitel/nitrogen/src/parser"
)

const PROMPT = ">> "

func main() {
	fmt.Print("Nitrogen Programming Language\n")
	fmt.Print("Type in commands at the prompt\n")
	startRepl(os.Stdin, os.Stdout)
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

		l := lexer.New(line)
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
