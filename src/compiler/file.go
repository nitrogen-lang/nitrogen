package compiler

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io/ioutil"
	"os"
	"time"
)

var ByteFileHeader = []byte{31, 'N', 'I', 'B'}

func WriteFile(name string, cb *CodeBlock) error {
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
	t := make([]byte, 8)
	binary.BigEndian.PutUint64(t, uint64(time.Now().Unix()))
	file.Write(t)
	file.Write(marshaled)
	return nil
}

func ReadFile(name string) (*CodeBlock, error) {
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
	return cb.(*CodeBlock), err
}
