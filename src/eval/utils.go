package eval

import (
	"github.com/nitrogen-lang/nitrogen/src/object"
)

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

// The first value is obj expressed as boolean
// The second is if obj is a valid bool-like object
func convertToBoolean(obj object.Object) (bool, bool) {
	isValid := object.ObjectIs(
		obj,
		object.BOOLEAN_OBJ,
		object.INTEGER_OBJ,
		object.FLOAT_OBJ,
		object.STRING_OBJ,
		object.NULL_OBJ,
	)

	return isTruthy(obj), isValid
}

func isException(obj object.Object) bool {
	return (obj != nil && obj.Type() == object.EXCEPTION_OBJ)
}

type stringStack struct {
	head   *stackElement
	length int
}

type stackElement struct {
	val  string
	prev *stackElement
}

func newStringStack() *stringStack {
	return &stringStack{}
}

func (s *stringStack) push(val string) {
	s.head = &stackElement{
		val:  val,
		prev: s.head,
	}
	s.length++
}

func (s *stringStack) getFront() string {
	if s.head == nil {
		return ""
	}
	return s.head.val
}

func (s *stringStack) pop() string {
	if s.head == nil {
		return ""
	}
	r := s.head.val
	s.head = s.head.prev
	s.length--
	return r
}

func (s *stringStack) len() int {
	return s.length
}
