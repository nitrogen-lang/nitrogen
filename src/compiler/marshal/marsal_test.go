package marshal

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/nitrogen-lang/nitrogen/src/compiler"
	"github.com/nitrogen-lang/nitrogen/src/moduleutils"

	"github.com/nitrogen-lang/nitrogen/src/object"
)

func TestIntegerMarshal(t *testing.T) {
	i1 := object.MakeIntObj(-2)
	bytes1, _ := Marshal(i1)
	bytes1Exp := []byte{0x69, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xfe}
	if !bytes.Equal(bytes1, bytes1Exp) {
		t.Errorf("Incorrect marshaled bytes. Expected %#v, got %#v", bytes1Exp, bytes1)
	}
	i1newo, _, _ := Unmarshal(bytes1)
	i1new := i1newo.(*object.Integer)
	if i1new.Value != -2 {
		t.Errorf("Incorrect unmarshaled value. Expected -2, got %d", i1new.Value)
	}

	i2 := object.MakeIntObj(2)
	bytes2, _ := Marshal(i2)
	bytes2Exp := []byte{0x69, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2}
	if !bytes.Equal(bytes2, bytes2Exp) {
		t.Errorf("Incorrect marshaled bytes. Expected %#v, got %#v", bytes2Exp, bytes2)
	}
	i2newo, _, _ := Unmarshal(bytes2)
	i2new := i2newo.(*object.Integer)
	if i2new.Value != 2 {
		t.Errorf("Incorrect unmarshaled value. Expected 2, got %d", i2new.Value)
	}
}

func TestFloatMarshal(t *testing.T) {
	f1 := object.MakeFloatObj(-2.1)
	bytes1, _ := Marshal(f1)
	bytes1Exp := []byte{0x66, 0xc0, 0x0, 0xcc, 0xcc, 0xcc, 0xcc, 0xcc, 0xcd}
	if !bytes.Equal(bytes1, bytes1Exp) {
		t.Errorf("Incorrect marshaled bytes. Expected %#v, got %#v", bytes1Exp, bytes1)
	}
	f1newo, _, _ := Unmarshal(bytes1)
	f1new := f1newo.(*object.Float)
	if f1new.Value != -2.1 {
		t.Errorf("Incorrect unmarshaled value. Expected -2.0, got %f", f1new.Value)
	}

	f2 := object.MakeFloatObj(2.2)
	bytes2, _ := Marshal(f2)
	bytes2Exp := []byte{0x66, 0x40, 0x1, 0x99, 0x99, 0x99, 0x99, 0x99, 0x9a}
	if !bytes.Equal(bytes2, bytes2Exp) {
		t.Errorf("Incorrect marshaled bytes. Expected %#v, got %#v", bytes2Exp, bytes2)
	}
	f2newo, _, _ := Unmarshal(bytes2)
	f2new := f2newo.(*object.Float)
	if f2new.Value != 2.2 {
		t.Errorf("Incorrect unmarshaled value. Expected 2.0, got %f", f2new.Value)
	}
}

func TestBooleanMarshal(t *testing.T) {
	b1 := object.NativeBoolToBooleanObj(true)
	bytes1, _ := Marshal(b1)
	bytes1Exp := []byte{'b', 1}
	if !bytes.Equal(bytes1, bytes1Exp) {
		t.Errorf("Incorrect marshaled bytes. Expected %#v, got %#v", bytes1Exp, bytes1)
	}
	b1new, _, _ := Unmarshal(bytes1)
	if b1new != object.TrueConst {
		t.Errorf("Incorrect unmarshaled value. Expected true, got false")
	}

	b2 := object.NativeBoolToBooleanObj(false)
	bytes2, _ := Marshal(b2)
	bytes2Exp := []byte{'b', 0}
	if !bytes.Equal(bytes2, bytes2Exp) {
		t.Errorf("Incorrect marshaled bytes. Expected %#v, got %#v", bytes2Exp, bytes2)
	}
	b2new, _, _ := Unmarshal(bytes2)
	if b2new != object.FalseConst {
		t.Errorf("Incorrect unmarshaled value. Expected false, got true")
	}
}

func TestNullMarshal(t *testing.T) {
	n1 := object.NullConst
	bytes1, _ := Marshal(n1)
	bytes1Exp := []byte{'n'}
	if !bytes.Equal(bytes1, bytes1Exp) {
		t.Errorf("Incorrect marshaled bytes. Expected %#v, got %#v", bytes1Exp, bytes1)
	}
	n1new, _, _ := Unmarshal(bytes1)
	if n1new != object.NullConst {
		t.Errorf("Incorrect unmarshaled value. Expected null constant")
	}
}

func TestStringMarshal(t *testing.T) {
	tests := []struct {
		s string
		b []byte
	}{
		{
			s: "Hello, World!",
			b: []byte{'s', 0x0, 0x0, 0x0, 0xd, 0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x2c, 0x20, 0x57, 0x6f, 0x72, 0x6c, 0x64, 0x21},
		},
		{
			s: "",
			b: []byte{'s', 0, 0, 0, 0},
		},
	}

	for i, test := range tests {
		s1 := object.MakeStringObj(test.s)
		bytes1, _ := Marshal(s1)
		if !bytes.Equal(bytes1, test.b) {
			t.Errorf("Test %d: Incorrect marshaled bytes. Expected %#v, got %#v", i+1, test.b, bytes1)
		}
		s1newo, _, _ := Unmarshal(bytes1)
		s1new := s1newo.(*object.String)
		if s1new.Value != test.s {
			t.Errorf("Test %d: Incorrect unmarshaled value. Expected -2, got %s", i+1, s1new.Value)
		}
	}
}

func TestCodeBlockMarshal(t *testing.T) {
	program, err := moduleutils.ASTCache.GetTree("./testdata/simple.ni")
	if err != nil {
		t.Fatal(err)
	}
	code := compiler.Compile(program, "__main")
	bytes, _ := Marshal(code)

	newcode, _, _ := Unmarshal(bytes)
	newcodeObj := newcode.(*compiler.CodeBlock)
	if !reflect.DeepEqual(code, newcodeObj) {
		t.Fatal("Code objects are not the same")
	}
}
