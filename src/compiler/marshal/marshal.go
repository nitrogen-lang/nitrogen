package marshal

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math"

	"github.com/nitrogen-lang/nitrogen/src/compiler"
	"github.com/nitrogen-lang/nitrogen/src/object"
)

func Marshal(o object.Object) ([]byte, error) {
	switch o := o.(type) {
	case *object.Integer:
		out := make([]byte, 9)
		out[0] = 'i'
		copy(out[1:], encodeUint64(uint64(o.Value)))
		return out, nil
	case *object.Float:
		out := make([]byte, 9)
		out[0] = 'f'
		copy(out[1:], encodeUint64(math.Float64bits(o.Value)))
		return out, nil
	case *object.String:
		val := string(o.Value)
		slen := len(val)
		out := make([]byte, slen+5)
		out[0] = 's'
		binary.BigEndian.PutUint32(out[1:5], uint32(slen))
		copy(out[5:], []byte(val))
		return out, nil
	case *object.Boolean:
		out := make([]byte, 2)
		out[0] = 'b'
		if o.Value {
			out[1] = 1
		} else {
			out[1] = 0
		}
		return out, nil
	case *object.Null:
		return []byte{'n'}, nil
	case *compiler.CodeBlock:
		buf := new(bytes.Buffer)
		tmpStr := object.MakeStringObj(o.Name)
		res, _ := Marshal(tmpStr)
		buf.Write(res)

		tmpStr.Value = []rune(o.Filename) // Reuse String object
		res, _ = Marshal(tmpStr)
		buf.Write(res)

		if o.Native {
			buf.WriteByte(1)
		} else {
			buf.WriteByte(0)
		}

		if o.ClassMethod {
			buf.WriteByte(1)
		} else {
			buf.WriteByte(0)
		}
		buf.Write(encodeUint16(uint16(o.LocalCount)))
		buf.Write(encodeUint16(uint16(o.MaxStackSize)))
		buf.Write(encodeUint16(uint16(o.MaxBlockSize)))
		buf.Write(encodeUint16(uint16(len(o.Constants))))
		for _, c := range o.Constants {
			res, err := Marshal(c)
			if err != nil {
				return nil, err
			}
			buf.Write(res)
		}
		buf.Write(encodeUint16(uint16(len(o.Locals))))
		for _, l := range o.Locals {
			tmpStr.Value = []rune(l)
			res, _ := Marshal(tmpStr) // No error check, strings are almost guaranteed to work
			buf.Write(res)
		}
		buf.Write(encodeUint16(uint16(len(o.Names))))
		for _, l := range o.Names {
			tmpStr.Value = []rune(l)
			res, _ := Marshal(tmpStr) // No error check, strings are almost guaranteed to work
			buf.Write(res)
		}
		buf.Write(encodeUint16(uint16(len(o.Code))))
		buf.Write(o.Code)

		clen := buf.Len()
		out := make([]byte, clen+9)
		out[0] = 'c'
		binary.BigEndian.PutUint64(out[1:9], uint64(clen))
		copy(out[9:], buf.Bytes())
		return out, nil
	}
	return nil, fmt.Errorf("Object type %T doesn't have a marshal implementation", o)
}

func Unmarshal(in []byte) (object.Object, []byte, error) {
	switch in[0] {
	case 'i':
		v := decodeUint64(in[1:9])
		return object.MakeIntObj(int64(v)), in[9:], nil
	case 'f':
		v := decodeUint64(in[1:9])
		return object.MakeFloatObj(math.Float64frombits(v)), in[9:], nil
	case 's':
		slen := int(binary.BigEndian.Uint32(in[1:5]))
		if len(in) < slen+5 {
			return nil, in, errors.New("Malformed string")
		}
		s := make([]byte, slen)
		copy(s, in[5:])
		return object.MakeStringObj(string(s)), in[slen+5:], nil
	case 'b':
		return object.NativeBoolToBooleanObj(in[1] == 1), in[2:], nil
	case 'n':
		return object.NullConst, in[1:], nil
	case 'c':
		inslice := in[9:] // Length is bytes [1-8]

		cb := &compiler.CodeBlock{}
		name, inslice, _ := Unmarshal(inslice)
		cb.Name = string(name.(*object.String).Value)

		filename, inslice, _ := Unmarshal(inslice)
		cb.Filename = string(filename.(*object.String).Value)

		cb.Native = inslice[0] == 1
		inslice = inslice[1:]
		cb.ClassMethod = inslice[0] == 1
		inslice = inslice[1:]
		cb.LocalCount = int(decodeUint16(inslice[:2]))
		inslice = inslice[2:]
		cb.MaxStackSize = int(decodeUint16(inslice[:2]))
		inslice = inslice[2:]
		cb.MaxBlockSize = int(decodeUint16(inslice[:2]))
		inslice = inslice[2:]

		constantsLen := int(decodeUint16(inslice[:2]))
		inslice = inslice[2:]
		cb.Constants = make([]object.Object, constantsLen)
		var err error
		for i := range cb.Constants {
			cb.Constants[i], inslice, err = Unmarshal(inslice)
			if err != nil {
				return nil, inslice, err
			}
		}

		localsLen := int(decodeUint16(inslice[:2]))
		inslice = inslice[2:]
		cb.Locals = make([]string, localsLen)
		for i := range cb.Locals {
			var tmpStr object.Object
			tmpStr, inslice, err = Unmarshal(inslice)
			if err != nil {
				return nil, inslice, err
			}
			cb.Locals[i] = string(tmpStr.(*object.String).Value)
		}

		namesLen := int(decodeUint16(inslice[:2]))
		inslice = inslice[2:]
		cb.Names = make([]string, namesLen)
		for i := range cb.Names {
			var tmpStr object.Object
			tmpStr, inslice, err = Unmarshal(inslice)
			if err != nil {
				return nil, inslice, err
			}
			cb.Names[i] = string(tmpStr.(*object.String).Value)
		}

		codeLen := int(decodeUint16(inslice[:2]))
		inslice = inslice[2:]
		cb.Code = make([]byte, codeLen)
		copy(cb.Code, inslice)

		return cb, inslice[codeLen:], nil
	}
	return nil, in, fmt.Errorf("Unknown unmarshal for type char %c", in[0])
}

func decodeUint64(in []byte) uint64 {
	return binary.BigEndian.Uint64(in)
}

func encodeUint64(in uint64) []byte {
	out := make([]byte, 8)
	binary.BigEndian.PutUint64(out, in)
	return out
}

func decodeUint16(in []byte) uint16 {
	return binary.BigEndian.Uint16(in)
}

func encodeUint16(in uint16) []byte {
	out := make([]byte, 2)
	binary.BigEndian.PutUint16(out, in)
	return out
}
