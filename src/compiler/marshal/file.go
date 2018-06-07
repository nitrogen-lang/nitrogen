package marshal

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io/ioutil"
	"os"
	"time"

	"github.com/nitrogen-lang/nitrogen/src/compiler"
)

var (
	ByteFileHeader = []byte{31, 'N', 'I', 'B'}
	VersionNumber  = []byte{0, 0, 0, 1}
)

func WriteFile(name string, cb *compiler.CodeBlock) error {
	marshaled, err := Marshal(cb)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(name, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	file.Write(ByteFileHeader)
	file.Write(VersionNumber)
	t := make([]byte, 8)
	binary.BigEndian.PutUint64(t, uint64(time.Now().Unix()))
	file.Write(t)
	file.Write(marshaled)
	return nil
}

func ReadFile(name string) (*compiler.CodeBlock, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileHeader := make([]byte, 4)
	file.Read(fileHeader)
	if !bytes.Equal(ByteFileHeader, fileHeader) {
		return nil, errors.New("File is not Nitrogen bytecode")
	}

	fileVersion := make([]byte, 4)
	file.Read(fileVersion)
	if !bytes.Equal(VersionNumber, fileVersion) {
		return nil, errors.New("File does not match current version")
	}

	fileTime := make([]byte, 8)
	if n, _ := file.Read(fileTime); n != 8 {
		return nil, errors.New("Invalid timestamp")
	}
	// Eventually check fileTime against the main source file

	rest, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	cb, _, err := Unmarshal(rest)
	return cb.(*compiler.CodeBlock), err
}
