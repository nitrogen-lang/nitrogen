package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/nitrogen-lang/nitrogen/src/moduleutils"
	"github.com/nitrogen-lang/nitrogen/src/vm"

	"github.com/nitrogen-lang/nitrogen/src/object"
)

var (
	scgiSock          string
	scgiWorkers       int
	scgiWorkerTimeout int
)

func init() {
	flag.StringVar(&scgiSock, "scgi-sock", "tcp:0.0.0.0:9000", "Socket to listen on for SCGI")
	flag.IntVar(&scgiWorkers, "scgi-workers", 5, "Number of workers to service SCGI requests")
	flag.IntVar(&scgiWorkerTimeout, "scgi-worker-timeout", 10, "Number of seconds to wait for an available worker before giving up")
}

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

	fmt.Printf("Creating %d workers\n", scgiWorkers)
	workerPool := make(chan *worker, scgiWorkers)
	for i := 0; i < scgiWorkers; i++ {
		workerPool <- &worker{id: i, workerPool: workerPool}
	}

	fmt.Printf("SCGI listening on %s\n", scgiSock)

	workerTimeout := time.Duration(scgiWorkerTimeout) * time.Second
	workerTimeoutTimer := time.NewTimer(workerTimeout)
	for {
		conn, err := ln.Accept()
		if err != nil {
			os.Stderr.WriteString(err.Error())
			os.Exit(1)
		}

		if !workerTimeoutTimer.Stop() { // Consume any time already in channel. See time.Timer docs.
			<-workerTimeoutTimer.C
		}
		workerTimeoutTimer.Reset(workerTimeout)

		select {
		case w := <-workerPool:
			go w.run(conn)
		case <-workerTimeoutTimer.C:
			os.Stderr.WriteString("Not enough workers\n")
			conn.Close()
		}
	}
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

type worker struct {
	id         int
	workerPool chan *worker
}

func (w *worker) run(conn net.Conn) {
	defer func() {
		conn.Close()
		w.workerPool <- w
	}()

	// Get Headers
	headerBytes, err := getNetString(conn)
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Stderr.Write([]byte{'\n'})
		return
	}

	values := bytes.Split(headerBytes, []byte{0})
	headers := make(map[string]string, len(values)/2)

	// Convert headers into a map
	for i := 0; i < len(values)-1; i += 2 {
		headers[string(values[i])] = string(values[i+1])
	}

	// Read body
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

	// Check valid SCGI request
	if headers["SCGI"] != "1" {
		os.Stderr.WriteString("Invalid SCGI header\n")
	}

	// Get script file path
	scriptFilename := headers["SCRIPT_FILENAME"]
	if scriptFilename == "" {
		scriptName := headers["SCRIPT_NAME"]
		docRoot := headers["DOCUMENT_ROOT"]
		scriptFilename = filepath.Join(docRoot, scriptName)
	}

	// Ensure all expected headers are set
	for _, header := range cgiHeaderNames {
		if headers[header] == "" {
			headers[header] = ""
		}
	}

	// Execute script
	code, err := moduleutils.CodeBlockCache.GetBlock(scriptFilename)
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Stderr.Write([]byte{'\n'})
		return
	}

	env := makeEnv(code.Filename)
	env.CreateConst("_FILE", object.MakeStringObj(code.Filename))
	env.SetForce("_SERVER", object.StringMapToHash(headers), true)

	vmsettings := vm.NewSettings()
	vmsettings.Stdout = conn

	result := vm.NewVM(vmsettings).Execute(code, env)

	if result != nil && result != object.NullConst {
		if e, ok := result.(*object.Exception); ok {
			os.Stderr.WriteString(e.Message)
			os.Stderr.Write([]byte{'\n'})
		}
	}
}
