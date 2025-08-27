package object

import (
	"bytes"
	"fmt"

	"github.com/nitrogen-lang/nitrogen/src/ast"
)

type ClassMethod interface {
	Object
	ClassMethod()
}

type Class struct {
	Name    string
	Parent  *Class
	Fields  []*ast.DefStatement
	Methods map[string]ClassMethod
}

func (c *Class) Inspect() string  { return "class " + c.Name }
func (c *Class) Type() ObjectType { return ClassObj }
func (c *Class) Dup() Object      { return NullConst }
func (c *Class) GetMethod(name string) ClassMethod {
	m, ok := c.Methods[name]
	if ok || c.Parent == nil {
		return m
	}
	return c.Parent.GetMethod(name)
}

type Instance struct {
	Class  *Class
	Fields *Environment
}

func (i *Instance) Inspect() string                   { return fmt.Sprintf("instance of %s", i.Class.Name) }
func (i *Instance) Type() ObjectType                  { return InstanceObj }
func (i *Instance) Dup() Object                       { return NullConst }
func (i *Instance) GetMethod(name string) ClassMethod { return i.Class.GetMethod(name) }

func InstanceOf(class string, i *Instance) bool {
	if i == nil {
		return false
	}

	c := i.Class
	for {
		if c.Name == class {
			return true
		}
		if c.Parent == nil {
			return false
		}
		c = c.Parent
	}
}

func InstanceOfAny(instance Object, classes ...string) bool {
	i, ok := instance.(*Instance)
	if !ok {
		return false
	}

	for _, class := range classes {
		if InstanceOf(class, i) {
			return true
		}
	}
	return false
}

func GetSelf(env *Environment) *Instance {
	self, ok := env.Get("self")
	if !ok {
		return nil
	}

	i, ok := self.(*Instance)
	if !ok {
		return nil
	}
	return i
}

type IfaceMethodDef struct {
	Name       string
	Parameters []string
}

func (i *IfaceMethodDef) Inspect() string {
	var out bytes.Buffer

	out.WriteString("fn ")
	out.WriteString(i.Name)
	out.WriteByte('(')
	for j := 0; j < len(i.Parameters); j++ {
		out.WriteString(i.Parameters[j])
		if j < len(i.Parameters)-1 {
			out.WriteString(", ")
		}
	}
	out.WriteByte(')')

	return out.String()
}

type Interface struct {
	Name    string
	Methods map[string]*IfaceMethodDef
}

func (i *Interface) Inspect() string  { return "interface " + i.Name }
func (i *Interface) Type() ObjectType { return InterfaceObj }
func (i *Interface) Dup() Object      { return NullConst }
func (i *Interface) GetMethod(name string) *IfaceMethodDef {
	return i.Methods[name]
}
