package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/nitrogen-lang/nitrogen/src/ast"
	"github.com/nitrogen-lang/nitrogen/src/eval"
	"github.com/nitrogen-lang/nitrogen/src/lexer"
	"github.com/nitrogen-lang/nitrogen/src/object"
	"github.com/nitrogen-lang/nitrogen/src/parser"
)

func startSCGIServer() {
	addrSplit := strings.SplitN(scgiSock, ":", 2)
	if len(addrSplit) != 2 {
		os.Stderr.WriteString("Invalid listening socket address\n")
		os.Exit(1)
	}

	switch addrSplit[0] {
	case "tcp", "unix":
	default:
		os.Stderr.WriteString("Listening socket must be tcp or unix\n")
		os.Exit(1)
	}

	ln, err := net.Listen(addrSplit[0], addrSplit[1])
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}

	fmt.Printf("SCGI listening on %s\n", scgiSock)

	for {
		conn, err := ln.Accept()
		if err != nil {
			os.Stderr.WriteString(err.Error())
			os.Exit(1)
		}
		handleSCGIRequest(conn)
	}
}

func handleSCGIRequest(conn net.Conn) {
	defer conn.Close()

	headerBytes, err := getNetString(conn)
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Stderr.Write([]byte{'\n'})
		return
	}

	values := bytes.Split(headerBytes, []byte{0})
	headers := make(map[string]string, len(values)/2)

	for i := 0; i < len(values)-1; i = i + 2 {
		headers[string(values[i])] = string(values[i+1])
	}

	contentLength, _ := strconv.Atoi(headers["CONTENT_LENGTH"])
	body := make([]byte, contentLength)
	n, err := conn.Read(body)
	if n != contentLength {
		os.Stderr.WriteString("Body length and Content-Length don't match\n")
		return
	}
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Stderr.Write([]byte{'\n'})
		return
	}

	if headers["SCGI"] != "1" {
		os.Stderr.WriteString("Invalid SCGI header\n")
	}

	appEnv := getEnvironmentMap()
	for k, v := range headers {
		appEnv[k] = v
	}

	env := object.NewEnvironment()
	env.CreateConst("_ENV", makeEnvironment(appEnv))

	scriptFilename := headers["SCRIPT_FILENAME"]
	if scriptFilename == "" {
		scriptName := headers["SCRIPT_NAME"]
		docRoot := headers["DOCUMENT_ROOT"]
		scriptFilename = filepath.Join(docRoot, scriptName)
	}

	program, parseErrors := parseFile(scriptFilename)
	if len(parseErrors) != 0 {
		printParserErrors(os.Stderr, parseErrors[:1])
		return
	}

	interpreter := eval.NewInterpreter()
	interpreter.Stdout = conn
	result := interpreter.Eval(program, env)
	if result != nil && result != object.NullConst {
		if e, ok := result.(*object.Exception); ok {
			os.Stderr.WriteString(e.Message)
			os.Stderr.Write([]byte{'\n'})
		}
	}
}

func parseFile(pathname string) (*ast.Program, []string) {
	l, err := lexer.NewFile(pathname)
	if err != nil {
		return nil, []string{err.Error()}
	}
	p := parser.New(l)
	program := p.ParseProgram()
	return program, p.Errors()
}

func getNetString(r io.Reader) ([]byte, error) {
	reader := bufio.NewReader(r)
	length, err := reader.ReadString(':')
	if err != nil {
		return nil, errors.New("Invalid SCGI header netstring")
	}

	l, err := strconv.Atoi(length[:len(length)-1])
	if err != nil {
		return nil, errors.New("Invalid SCGI header length")
	}

	netstring := make([]byte, l)
	n, err := reader.Read(netstring)
	if n != l || err != nil {
		return nil, errors.New("Invalid SCGI header")
	}

	if b, _ := reader.ReadByte(); b != ',' {
		return nil, errors.New("Invalid SCGI header netstring")
	}

	return netstring, nil
}
