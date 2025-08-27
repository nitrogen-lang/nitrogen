package moduleutils_test

import (
	"testing"

	"github.com/nitrogen-lang/nitrogen/src/compiler"
	"github.com/nitrogen-lang/nitrogen/src/elemental/object"
	"github.com/nitrogen-lang/nitrogen/src/elemental/vm"
	"github.com/nitrogen-lang/nitrogen/src/lexer"
	"github.com/nitrogen-lang/nitrogen/src/parser"
)

func TestEval(input string) object.Object {
	p := parser.New(lexer.NewString(input), nil)
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		return object.MakeStringArray(p.Errors())
	}

	code := compiler.Compile(program, "__main")
	env := object.NewEnvironment()
	vmSettings := vm.NewSettings()
	// Force the virtual machine to not panic on uncaught exceptions
	vmSettings.ReturnExceptions = true
	machine := vm.NewVM(vmSettings)

	ret, _ := machine.Execute(code, env, "__main")
	return ret
}

// Verification functions
func TestIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. expected=%d, got=%T (%+v)",
			expected,
			obj,
			ShowError(obj),
		)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d",
			result.Value, expected)
		return false
	}
	return true
}

func TestFloatObject(t *testing.T, obj object.Object, expected float64) bool {
	result, ok := obj.(*object.Float)
	if !ok {
		t.Errorf("object is not Float. expected=%g, got=%T (%+v)",
			expected,
			obj,
			ShowError(obj),
		)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. got=%g, want=%g",
			result.Value, expected)
		return false
	}
	return true
}

func TestBoolObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Bool. expected=%t, got=%T (%+v)",
			expected,
			obj,
			ShowError(obj),
		)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t",
			result.Value, expected)
		return false
	}
	return true
}

func ShowError(obj object.Object) string {
	if obj == nil {
		return "nil"
	}

	switch o := obj.(type) {
	case *object.Exception:
		return o.Inspect()
	case *object.Array:
		return o.Inspect()
	default:
		return obj.Type().String()
	}
}

func TestLiteralErrorObjects(t *testing.T, got object.Object, expected interface{}) {
	switch expected := expected.(type) {
	case int:
		TestIntegerObject(t, got, int64(expected))
	case float64:
		TestFloatObject(t, got, expected)
	case bool:
		TestBoolObject(t, got, expected)
	case string:
		errObj, ok := got.(*object.Exception)
		if !ok {
			t.Errorf("object is not Error. got=%T (%+v)", got, ShowError(got))
			return
		}

		if errObj.Message != expected {
			t.Errorf("wrong error message. expected=%q, got=%q", expected, errObj.Message)
		}
	}
}
