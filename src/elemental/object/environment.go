package object

import (
	"errors"
	"fmt"
)

var (
	constError        = errors.New("constant can't not be changed")
	errAlreadyDefined = errors.New("symbol already defined")
	errNotDefined     = errors.New("symbol not defined")
)

func IsConstErr(e error) bool {
	return e == constError
}

type eco struct {
	name     string
	v        Object
	n        *eco
	readonly bool
	export   bool
}

type Environment struct {
	root      *eco
	parent    *Environment
	localOnly bool
}

func NewEnvironment() *Environment {
	return NewSizedEnvironment(false)
}

func NewLocalEnvironment() *Environment {
	return NewSizedEnvironment(true)
}

func NewSizedEnvironment(localOnly bool) *Environment {
	return &Environment{localOnly: localOnly}
}

func NewEnclosedEnv(outer *Environment) *Environment {
	return NewSizedEnclosedEnv(outer, false)
}

func NewEnclosedLocalEnv(outer *Environment) *Environment {
	return NewSizedEnclosedEnv(outer, true)
}

func NewSizedEnclosedEnv(outer *Environment, localOnly bool) *Environment {
	env := NewSizedEnvironment(localOnly)
	env.parent = outer
	return env
}

func (e *Environment) Clone() *Environment {
	return &Environment{
		root:   e.root,
		parent: e.parent,
	}
}

func (e *Environment) SetParent(env *Environment) {
	e.parent = env
}

func (e *Environment) Parent() *Environment {
	if e == nil {
		return nil
	}
	return e.parent
}

func (e *Environment) Print(indent string) {
	if e == nil {
		fmt.Println("{}")
		return
	}

	if e.root != nil {
		v := e.root
		for v != nil {
			fmt.Printf("%s%s = %s\n  %sConst: %t\n", indent, v.name, v.v.Inspect(), indent, v.readonly)
			v = v.n
		}
	}

	if e.parent != nil {
		if e.parent.parent == nil {
			fmt.Printf("\n%sGlobal:\n", indent)
		} else {
			fmt.Printf("\n%sParent:\n", indent)
		}
		e.parent.Print(indent + "  ")
	}
}

func (e *Environment) find(name string) *eco {
	if e == nil {
		return nil
	}

	v := e.root
	for v != nil {
		if v.name == name {
			return v
		}
		v = v.n
	}
	return nil
}

func (e *Environment) Get(name string) (Object, bool) {
	obj := e.find(name)
	if obj != nil {
		return obj.v, true
	}

	if e.parent != nil {
		return e.parent.Get(name)
	}
	return nil, false
}

func (e *Environment) GetLocal(name string) (Object, bool) {
	obj := e.find(name)
	if obj == nil {
		return nil, false
	}
	return obj.v, true
}

func (e *Environment) IsConst(name string) bool {
	obj := e.find(name)
	if obj != nil {
		return obj.readonly
	}

	if !e.localOnly && e.parent != nil {
		return e.parent.IsConst(name)
	}
	return false
}

func (e *Environment) IsConstLocal(name string) bool {
	obj := e.find(name)
	if obj != nil {
		return obj.readonly
	}
	return false
}

func (e *Environment) Create(name string, val Object) (Object, error) {
	obj := e.find(name)
	if obj != nil {
		return nil, errAlreadyDefined
	}

	e.root = &eco{
		name: name,
		n:    e.root,
		v:    val,
	}
	return val, nil
}

func (e *Environment) CreateConst(name string, val Object) (Object, error) {
	obj := e.find(name)
	if obj != nil {
		return nil, errAlreadyDefined
	}

	e.root = &eco{
		name:     name,
		n:        e.root,
		v:        val,
		readonly: true,
	}
	return val, nil
}

func (e *Environment) Set(name string, val Object) (Object, error) {
	obj := e.find(name)
	if obj != nil {
		obj.v = val
		return val, nil
	}

	if e.parent != nil {
		return e.parent.Set(name, val)
	}
	return nil, errNotDefined
}

func (e *Environment) SetLocal(name string, val Object) (Object, error) {
	obj := e.find(name)
	if obj != nil {
		obj.v = val
		return val, nil
	}
	return nil, errNotDefined
}

func (e *Environment) SetForce(name string, val Object, readonly bool) {
	obj := e.find(name)
	if obj != nil {
		obj.v = val
		obj.readonly = readonly
		return
	}

	e.root = &eco{
		name:     name,
		n:        e.root,
		v:        val,
		readonly: readonly,
	}
}

func (e *Environment) findParentNode(name string) (*eco, *eco) {
	if e == nil {
		return nil, nil // No environment
	}

	v := e.root
	if v == nil {
		return nil, nil // Environment has no items
	}

	if v.name == name {
		return nil, v // Element is the root, no parent
	}

	for v.n != nil {
		if v.n.name == name {
			return v, v.n // Element was found and has parent
		}
		v = v.n
	}
	return nil, nil // Element not found
}

func (e *Environment) UnsetLocal(name string) {
	p, el := e.findParentNode(name)
	if p != nil {
		p.n = p.n.n
		return
	}
	if el != nil {
		e.root = el.n
	}
}

func (e *Environment) Unset(name string) {
	p, el := e.findParentNode(name)
	if p != nil {
		p.n = p.n.n
		return
	}
	if el != nil {
		e.root = el.n
		return
	}

	if e.parent != nil {
		e.parent.Unset(name)
	}
}

func (e *Environment) Export(name string) {
	obj := e.find(name)
	if obj == nil {
		return
	}
	obj.export = true
}

func (e *Environment) GetExported() map[string]Object {
	if e == nil {
		return nil
	}

	exported := make(map[string]Object)

	v := e.root
	for v != nil {
		if v.export {
			exported[v.name] = v.v
		}
		v = v.n
	}

	return exported
}
