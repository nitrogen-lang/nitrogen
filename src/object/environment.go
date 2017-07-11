package object

import "errors"

var (
	constError = errors.New("constant can't not be changed")
)

func IsConstErr(e error) bool {
	return e == constError
}

type eco struct {
	v        Object
	readonly bool
}

type Environment struct {
	store  map[string]*eco
	parent *Environment
}

func NewEnvironment() *Environment {
	return &Environment{
		store: make(map[string]*eco),
	}
}

func NewEnclosedEnv(outer *Environment) *Environment {
	env := NewEnvironment()
	env.parent = outer
	return env
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if ok {
		return obj.v, ok
	}

	if e.parent != nil {
		return e.parent.Get(name)
	}
	return nil, false
}

func (e *Environment) IsConst(name string) bool {
	obj, ok := e.store[name]
	if ok {
		return obj.readonly
	}

	if e.parent != nil {
		return e.parent.IsConst(name)
	}
	return false
}

func (e *Environment) isConst(name string) bool {
	obj, ok := e.store[name]
	if ok {
		return obj.readonly
	}
	return false
}

func (e *Environment) Set(name string, val Object) (Object, error) {
	if e.isConst(name) {
		return nil, constError
	}

	e.store[name] = &eco{v: val}
	return val, nil
}

func (e *Environment) SetConst(name string, val Object) (Object, error) {
	if e.isConst(name) {
		return nil, constError
	}

	e.store[name] = &eco{
		v:        val,
		readonly: true,
	}
	return val, nil
}
