package object

import (
	"fmt"

	"github.com/nitrogen-lang/nitrogen/src/ast"
)

type Class struct {
	Name    string
	Parent  string
	Fields  []*ast.DefStatement
	Methods map[string]*Function
}

func (c *Class) Inspect() string                 { return "class " + c.Name }
func (c *Class) Type() ObjectType                { return ClassObj }
func (c *Class) Dup() Object                     { return NullConst }
func (c *Class) GetMethod(name string) *Function { return c.Methods[name] }

type Instance struct {
	Class  *Class
	Fields *Environment
}

func (i *Instance) Inspect() string                 { return fmt.Sprintf("instance of %s", i.Class.Name) }
func (i *Instance) Type() ObjectType                { return InstanceObj }
func (i *Instance) Dup() Object                     { return NullConst }
func (i *Instance) GetMethod(name string) *Function { return i.Class.GetMethod(name) }

func InstanceOf(class string, instance Object) *Instance {
	i, ok := instance.(*Instance)
	if !ok || i.Class.Name != class {
		return nil
	}
	return i
}

func InstanceOfAny(instance Object, classes ...string) bool {
	i, ok := instance.(*Instance)
	if !ok {
		return false
	}

	for _, class := range classes {
		if i.Class.Name == class {
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
