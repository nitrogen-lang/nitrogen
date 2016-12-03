package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/lfkeitel/nitrogen/src/lexer"
	"github.com/lfkeitel/nitrogen/src/token"
)

const PROMPT = ">> "

func main() {
	fmt.Print("Nitrogen Programming Language\n")
	fmt.Print("Type in commands at the prompt\n")
	startRepl(os.Stdin, os.Stdout)
}

func startRepl(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

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

		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Fprintf(out, "%+v\n", tok)
		}
	}
}
