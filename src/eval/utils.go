package eval

import (
	"github.com/nitrogen-lang/nitrogen/src/object"
)

func nativeBoolToBooleanObj(input bool) *object.Boolean {
	if input {
		return object.TRUE
	}
	return object.FALSE
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}
	return obj
}

func isTruthy(obj object.Object) bool {
	// Null or false is immediately not true
	if obj == object.NULL || obj == object.FALSE {
		return false
	}

	// True is immediately true
	if obj == object.TRUE {
		return true
	}

	// If the object is an INT, it's truthy if it doesn't equal 0
	if obj.Type() == object.INTEGER_OBJ {
		return (obj.(*object.Integer).Value != 0)
	}

	// Same as above if but with floats
	if obj.Type() == object.FLOAT_OBJ {
		return (obj.(*object.Float).Value != 0.0)
	}

	// Empty string is false, non-empty is true
	if obj.Type() == object.STRING_OBJ {
		return (obj.(*object.String).Value != "")
	}

	// Assume value is false
	return false
}

func isError(obj object.Object) bool {
	return (obj != nil && obj.Type() == object.ERROR_OBJ)
}

func typesEqualTo(t object.ObjectType, a, b object.Object) bool {
	return (a.Type() == t && b.Type() == t)
}

func objectIs(o object.Object, t ...object.ObjectType) bool {
	for _, ot := range t {
		if o.Type() == ot {
			return true
		}
	}
	return false
}
