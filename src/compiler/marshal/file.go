package marshal

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io/fs"
	"io/ioutil"
	"os"
	"time"

	"github.com/nitrogen-lang/nitrogen/src/compiler"
)

var (
	ByteFileHeader = []byte{31, 'N', 'I', 'B'}
	VersionNumber  = []byte{0, 0, 0, 9}

	ErrVersion = errors.New("File does not match current version")
)

const execHeader = "#!/usr/bin/nitrogenrun\n"

func IsErrVersion(err error) bool {
	return err == ErrVersion
}

func WriteFile(name string, cb *compiler.CodeBlock, ts time.Time, executable bool) error {
	marshaled, err := Marshal(cb)
	if err != nil {
		return err
	}

	fileMode := 0644
	if executable {
		fileMode = 0755
	}

	file, err := os.OpenFile(name, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, fs.FileMode(fileMode))
	if err != nil {
		return err
	}
	defer file.Close()

	if executable {
		file.WriteString(execHeader)
	}

	if ts.IsZero() {
		ts = time.Now()
	}
	ts = ts.Round(time.Second)

	file.Write(ByteFileHeader)
	file.Write(VersionNumber)
	t := make([]byte, 8)
	binary.BigEndian.PutUint64(t, uint64(ts.Unix()))
	file.Write(t)
	file.Write(marshaled)
	return nil
}

type FileInfo struct {
	Filename string
	Version  []byte
	ModTime  time.Time
}

func ReadFile(name string) (*compiler.CodeBlock, *FileInfo, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	fi := &FileInfo{Filename: name}

	firstChar := make([]byte, 1)
	file.Read(firstChar)
	if firstChar[0] == '#' {
		file.Seek(int64(len(execHeader)), 0)
	} else {
		file.Seek(0, 0)
	}

	fileHeader := make([]byte, 4)
	file.Read(fileHeader)
	if !bytes.Equal(ByteFileHeader, fileHeader) {
		return nil, nil, errors.New("File is not Nitrogen bytecode")
	}

	fileVersion := make([]byte, 4)
	file.Read(fileVersion)
	if !bytes.Equal(VersionNumber, fileVersion) {
		return nil, nil, ErrVersion
	}
	fi.Version = fileVersion

	fileTime := make([]byte, 8)
	n, _ := file.Read(fileTime)
	if n != 8 {
		return nil, nil, errors.New("Invalid timestamp")
	}
	// fileTime is checked by caller if they care about it
	fi.ModTime = time.Unix(int64(binary.BigEndian.Uint64(fileTime)), 0)

	rest, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, nil, err
	}

	cb, _, err := Unmarshal(rest)
	if err != nil {
		return nil, nil, err
	}
	return cb.(*compiler.CodeBlock), fi, nil
}
