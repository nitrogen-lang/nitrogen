package vm

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"github.com/nitrogen-lang/nitrogen/src/compiler"
	"github.com/nitrogen-lang/nitrogen/src/object"
)

var (
	builtins      = map[string]*object.Builtin{}
	modules       = map[string]*object.Module{}
	nativeFn      = map[string]*object.Builtin{}
	nativeMethods = map[string]*BuiltinMethod{}
	identRegex    = regexp.MustCompile(`[a-zA-Z_][a-zA-Z0-9_]*`)
)

// RegisterBuiltin allows other packages to register functions for availability in user code
func RegisterBuiltin(name string, fn object.BuiltinFunction) {
	if !validBuiltinIdent(name) {
		panic("Invalid VM builtin function name " + name)
	}

	if _, defined := builtins[name]; defined {
		// Panic because this should NEVER happen when built
		panic("Builtin VM function " + name + " already defined")
	}

	builtins[name] = &object.Builtin{Fn: fn}
}

func RegisterModule(name string, m *object.Module) {
	for k := range m.Methods {
		if !validBuiltinIdent(k) {
			panic("Invalid VM module function name " + name)
		}
	}

	if _, defined := modules[name]; defined {
		// Panic because this should NEVER happen when built
		panic("VM module " + name + " already defined")
	}

	modules[name] = m
}

func RegisterNative(name string, fn object.BuiltinFunction) {
	if _, defined := nativeFn[name]; defined {
		// Panic because this should NEVER happen when built
		panic("VM native func " + name + " already defined")
	}

	nativeFn[name] = &object.Builtin{Fn: fn}
}

func RegisterNativeMethod(name string, fn BuiltinMethodFunction) {
	if _, defined := nativeMethods[name]; defined {
		// Panic because this should NEVER happen when built
		panic("VM native method " + name + " already defined")
	}

	nativeMethods[name] = &BuiltinMethod{
		Name: name[strings.LastIndex(name, ".")+1:],
		Fn:   fn,
	}
}

func validBuiltinIdent(ident string) bool {
	return identRegex.Match([]byte(ident))
}

func getBuiltin(name string) object.Object {
	if builtin, defined := builtins[name]; defined {
		return builtin
	}
	return nil
}

// GetModule returns a Module object is a module with the given name is registered, otherwise nil.
func GetModule(name string) *object.Module {
	if module, defined := modules[name]; defined {
		return module
	}
	return nil
}

type VMFunction struct {
	Name       string
	Parameters []string
	Native     bool
	Body       *compiler.CodeBlock
	Env        *object.Environment
	Class      *VMClass
}

func (f *VMFunction) Inspect() string {
	var out bytes.Buffer

	out.WriteString("func")
	out.WriteByte(' ')
	out.WriteString(f.Name)
	out.WriteByte('(')
	out.WriteString(strings.Join(f.Parameters, ", "))
	out.WriteString(") {...}")

	return out.String()
}
func (f *VMFunction) Type() object.ObjectType { return object.FunctionObj }
func (f *VMFunction) Dup() object.Object      { return f }
func (f *VMFunction) ClassMethod()            {}

type VMClass struct {
	Name    string
	Parent  *VMClass
	Fields  *compiler.CodeBlock
	Methods map[string]object.ClassMethod
}

func (c *VMClass) Inspect() string {
	return fmt.Sprintf("Class %s with %d methods", c.Name, len(c.Methods))
}
func (c *VMClass) Type() object.ObjectType { return object.ClassObj }
func (c *VMClass) Dup() object.Object      { return object.NullConst }
func (c *VMClass) GetMethod(name string) object.ClassMethod {
	m, ok := c.Methods[name]
	if ok || c.Parent == nil {
		return m
	}
	return c.Parent.GetMethod(name)
}

type BuiltinClass struct {
	*VMClass
	Fields map[string]object.Object
}

type VMInstance struct {
	Class  *VMClass
	Fields *object.Environment
}

func (i *VMInstance) Inspect() string                          { return fmt.Sprintf("instance of %s", i.Class.Name) }
func (i *VMInstance) Type() object.ObjectType                  { return object.InstanceObj }
func (i *VMInstance) Dup() object.Object                       { return object.NullConst }
func (i *VMInstance) GetMethod(name string) object.ClassMethod { return i.Class.GetMethod(name) }
func (i *VMInstance) GetBoundMethod(name string) *BoundMethod {
	method := i.Class.GetMethod(name)
	if method == nil {
		return nil
	}

	return &BoundMethod{
		Method:   method,
		Instance: i,
		Parent:   i.Class.Parent,
	}
}

type BoundMethod struct {
	Method   object.ClassMethod
	Instance *VMInstance
	Parent   *VMClass
}

func (b *BoundMethod) Inspect() string {
	return fmt.Sprintf("method bound to instance of %s", b.Instance.Class.Name)
}
func (b *BoundMethod) Type() object.ObjectType { return object.BoundMethodObj }
func (b *BoundMethod) Dup() object.Object      { return object.NullConst }

type BuiltinMethodFunction func(i *VirtualMachine, self *VMInstance, env *object.Environment, args ...object.Object) object.Object

type BuiltinMethod struct {
	Fn   BuiltinMethodFunction
	Name string
}

func (b *BuiltinMethod) Inspect() string         { return "builtin method" }
func (b *BuiltinMethod) Type() object.ObjectType { return object.BuiltinMethodObj }
func (b *BuiltinMethod) Dup() object.Object      { return b }
func (b *BuiltinMethod) ClassMethod()            {}

func MakeBuiltinMethod(fn BuiltinMethodFunction) *BuiltinMethod {
	return &BuiltinMethod{Fn: fn}
}

func InstanceOf(class string, i *VMInstance) bool {
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
