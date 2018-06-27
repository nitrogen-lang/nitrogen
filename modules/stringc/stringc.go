package main

import (
	"strings"

	"github.com/nitrogen-lang/nitrogen/src/moduleutils"
	"github.com/nitrogen-lang/nitrogen/src/object"
	"github.com/nitrogen-lang/nitrogen/src/vm"
)

func init() {
	vm.RegisterModule(ModuleName, &object.Module{
		Name:    ModuleName,
		Methods: map[string]object.BuiltinFunction{},
		Vars: map[string]object.Object{
			"string": &vm.BuiltinClass{
				Fields: map[string]object.Object{
					"str": object.NullConst,
				},
				VMClass: &vm.VMClass{
					Name:   "string",
					Parent: nil,
					Methods: map[string]object.ClassMethod{
						"init":      vm.MakeBuiltinMethod(vmStringInit),
						"splitN":    vm.MakeBuiltinMethod(vmStrSplitN),
						"trimSpace": vm.MakeBuiltinMethod(vmStrTrim),
						"dedup":     vm.MakeBuiltinMethod(vmStrDedup),
					},
				},
			},
		},
	})
}

func main() {}

var ModuleName = "stringc"

func stringInit(interpreter object.Interpreter, self *object.Instance, env *object.Environment, args ...object.Object) object.Object {
	_, ok := args[0].(*object.String)
	if !ok {
		return object.NewException("string expected a string, got %s", args[1].Type().String())
	}

	env.Set("str", args[0])
	return object.NullConst
}

func strSplitN(interpreter object.Interpreter, self *object.Instance, env *object.Environment, args ...object.Object) object.Object {
	if ac := moduleutils.CheckArgs("strSplitN", 2, args...); ac != nil {
		return ac
	}

	selfStr, _ := self.Fields.Get("str")
	target, ok := selfStr.(*object.String)
	if !ok {
		return object.NewException("str field is not a string")
	}

	sep, ok := args[0].(*object.String)
	if !ok {
		return object.NewException("splitN expected a string, got %s", args[1].Type().String())
	}

	count, ok := args[1].(*object.Integer)
	if !ok {
		return object.NewException("splitN expected an int, got %s", args[1].Type().String())
	}

	return object.MakeStringArray(strings.SplitN(target.String(), sep.String(), int(count.Value)))
}

func strTrim(interpreter object.Interpreter, self *object.Instance, env *object.Environment, args ...object.Object) object.Object {
	selfStr, _ := self.Fields.Get("str")
	target, ok := selfStr.(*object.String)
	if !ok {
		return object.NewException("str field is not a string")
	}

	return object.MakeStringObj(strings.TrimSpace(target.String()))
}

func strDedup(interpreter object.Interpreter, self *object.Instance, env *object.Environment, args ...object.Object) object.Object {
	if ac := moduleutils.CheckArgs("strDedup", 1, args...); ac != nil {
		return ac
	}

	selfStr, _ := self.Fields.Get("str")
	target, ok := selfStr.(*object.String)
	if !ok {
		return object.NewException("str field is not a string")
	}

	dedup, ok := args[0].(*object.String)
	if !ok {
		return object.NewException("strDedup expected a string, got %s", args[0].Type().String())
	}

	if len(dedup.Value) != 1 {
		return object.NewException("Dedup string must be one byte")
	}

	return object.MakeStringObj(dedupString(target.Value, dedup.Value[0]))
}

func dedupString(str []rune, c rune) string {
	newstr := make([]rune, 0, int(float32(len(str))*0.75))

	var lastc rune
	for _, char := range str {
		if char == c && char == lastc {
			continue
		}
		newstr = append(newstr, char)
		lastc = char
	}

	return string(newstr)
}

func vmStringInit(interpreter *vm.VirtualMachine, self *vm.VMInstance, env *object.Environment, args ...object.Object) object.Object {
	_, ok := args[0].(*object.String)
	if !ok {
		return object.NewException("string expected a string, got %s", args[1].Type().String())
	}

	env.Set("str", args[0])
	return object.NullConst
}

func vmStrSplitN(interpreter *vm.VirtualMachine, self *vm.VMInstance, env *object.Environment, args ...object.Object) object.Object {
	if ac := moduleutils.CheckArgs("strSplitN", 2, args...); ac != nil {
		return ac
	}

	selfStr, _ := self.Fields.Get("str")
	target, ok := selfStr.(*object.String)
	if !ok {
		return object.NewException("str field is not a string")
	}

	sep, ok := args[0].(*object.String)
	if !ok {
		return object.NewException("splitN expected a string, got %s", args[1].Type().String())
	}

	count, ok := args[1].(*object.Integer)
	if !ok {
		return object.NewException("splitN expected an int, got %s", args[1].Type().String())
	}

	return object.MakeStringArray(strings.SplitN(target.String(), sep.String(), int(count.Value)))
}

func vmStrTrim(interpreter *vm.VirtualMachine, self *vm.VMInstance, env *object.Environment, args ...object.Object) object.Object {
	selfStr, _ := self.Fields.Get("str")
	target, ok := selfStr.(*object.String)
	if !ok {
		return object.NewException("str field is not a string")
	}

	return object.MakeStringObj(strings.TrimSpace(target.String()))
}

func vmStrDedup(interpreter *vm.VirtualMachine, self *vm.VMInstance, env *object.Environment, args ...object.Object) object.Object {
	if ac := moduleutils.CheckArgs("strDedup", 1, args...); ac != nil {
		return ac
	}

	selfStr, _ := self.Fields.Get("str")
	target, ok := selfStr.(*object.String)
	if !ok {
		return object.NewException("str field is not a string")
	}

	dedup, ok := args[0].(*object.String)
	if !ok {
		return object.NewException("strDedup expected a string, got %s", args[0].Type().String())
	}

	if len(dedup.Value) != 1 {
		return object.NewException("Dedup string must be one byte")
	}

	return object.MakeStringObj(dedupString(target.Value, dedup.Value[0]))
}
