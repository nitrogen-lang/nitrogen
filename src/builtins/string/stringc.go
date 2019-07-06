package string

import (
	"strings"

	"github.com/nitrogen-lang/nitrogen/src/moduleutils"
	"github.com/nitrogen-lang/nitrogen/src/object"
	"github.com/nitrogen-lang/nitrogen/src/vm"
)

var moduleName = "std/string"

func init() {
	vm.RegisterModule(moduleName, &object.Module{
		Name: moduleName,
		Methods: map[string]object.BuiltinFunction{
			"contains":  strContains,
			"count":     strCount,
			"dedup":     strDedup,
			"format":    strFormat,
			"hasPrefix": strHasPrefix,
			"hasSuffix": strHasSuffix,
			"replace":   strReplace,
			"split":     strSplit,
			"splitN":    strSplitN,
			"trimSpace": strTrim,
		},
		Vars: map[string]object.Object{
			"String": &vm.BuiltinClass{
				Fields: map[string]object.Object{
					"str": object.NullConst,
				},
				VMClass: &vm.VMClass{
					Name:   "String",
					Parent: nil,
					Methods: map[string]object.ClassMethod{
						"contains":  vm.MakeBuiltinMethod(vmStrContains),
						"count":     vm.MakeBuiltinMethod(vmStrCount),
						"dedup":     vm.MakeBuiltinMethod(vmStrDedup),
						"format":    vm.MakeBuiltinMethod(vmStrFormat),
						"hasPrefix": vm.MakeBuiltinMethod(vmStrHasPrefix),
						"hasSuffix": vm.MakeBuiltinMethod(vmStrHasSuffix),
						"init":      vm.MakeBuiltinMethod(vmStringInit),
						"replace":   vm.MakeBuiltinMethod(vmStrReplace),
						"split":     vm.MakeBuiltinMethod(vmStrSplit),
						"splitN":    vm.MakeBuiltinMethod(vmStrSplitN),
						"trimSpace": vm.MakeBuiltinMethod(vmStrTrim),
					},
				},
			},
		},
	})
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
	return nil
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

func vmStrSplit(interpreter *vm.VirtualMachine, self *vm.VMInstance, env *object.Environment, args ...object.Object) object.Object {
	if ac := moduleutils.CheckArgs("split", 1, args...); ac != nil {
		return ac
	}

	return vmStrSplitN(interpreter, self, env, args[0], object.MakeIntObj(-1))
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

func vmStrFormat(interpreter *vm.VirtualMachine, self *vm.VMInstance, env *object.Environment, args ...object.Object) object.Object {
	selfStr, _ := self.Fields.Get("str")
	target, ok := selfStr.(*object.String)
	if !ok {
		return object.NewException("str field is not a string")
	}

	t := string(target.Value)

	for _, arg := range args {
		if !strings.Contains(t, "{}") {
			break
		}

		s := objectToString(arg, interpreter)
		t = strings.Replace(t, "{}", s, 1)
	}

	return object.MakeStringObj(t)
}

func objectToString(obj object.Object, machine *vm.VirtualMachine) string {
	if instance, ok := obj.(*vm.VMInstance); ok {
		toString := instance.GetBoundMethod("toString")
		if toString != nil {
			machine.CallFunction(0, toString, true, nil, false)
			return objectToString(machine.PopStack(), machine)
		}
	}

	return obj.Inspect()
}

func vmStrContains(interpreter *vm.VirtualMachine, self *vm.VMInstance, env *object.Environment, args ...object.Object) object.Object {
	if ac := moduleutils.CheckArgs("contains", 1, args...); ac != nil {
		return ac
	}

	selfStr, _ := self.Fields.Get("str")
	target, ok := selfStr.(*object.String)
	if !ok {
		return object.NewException("str field is not a string")
	}

	sub, ok := args[0].(*object.String)
	if !ok {
		return object.NewException("contains expected a string, got %s", args[1].Type().String())
	}

	return object.NativeBoolToBooleanObj(strings.Contains(target.String(), sub.String()))
}

func vmStrCount(interpreter *vm.VirtualMachine, self *vm.VMInstance, env *object.Environment, args ...object.Object) object.Object {
	if ac := moduleutils.CheckArgs("count", 1, args...); ac != nil {
		return ac
	}

	selfStr, _ := self.Fields.Get("str")
	target, ok := selfStr.(*object.String)
	if !ok {
		return object.NewException("str field is not a string")
	}

	sub, ok := args[0].(*object.String)
	if !ok {
		return object.NewException("count expected a string, got %s", args[0].Type().String())
	}
	if len(sub.Value) == 0 {
		return object.NewException("count argument 2 can't be empty")
	}

	return object.MakeIntObj(int64(strings.Count(target.String(), sub.String())))
}

func vmStrHasPrefix(interpreter *vm.VirtualMachine, self *vm.VMInstance, env *object.Environment, args ...object.Object) object.Object {
	if ac := moduleutils.CheckArgs("hasPrefix", 1, args...); ac != nil {
		return ac
	}

	selfStr, _ := self.Fields.Get("str")
	target, ok := selfStr.(*object.String)
	if !ok {
		return object.NewException("str field is not a string")
	}

	prefix, ok := args[0].(*object.String)
	if !ok {
		return object.NewException("hasPrefix expected a string, got %s", args[0].Type().String())
	}

	return object.NativeBoolToBooleanObj(strings.HasPrefix(target.String(), prefix.String()))
}

func vmStrHasSuffix(interpreter *vm.VirtualMachine, self *vm.VMInstance, env *object.Environment, args ...object.Object) object.Object {
	if ac := moduleutils.CheckArgs("hasSuffix", 1, args...); ac != nil {
		return ac
	}

	selfStr, _ := self.Fields.Get("str")
	target, ok := selfStr.(*object.String)
	if !ok {
		return object.NewException("str field is not a string")
	}

	prefix, ok := args[0].(*object.String)
	if !ok {
		return object.NewException("hasSuffix expected a string, got %s", args[0].Type().String())
	}

	return object.NativeBoolToBooleanObj(strings.HasSuffix(target.String(), prefix.String()))
}

func vmStrReplace(interpreter *vm.VirtualMachine, self *vm.VMInstance, env *object.Environment, args ...object.Object) object.Object {
	if ac := moduleutils.CheckArgs("replace", 3, args...); ac != nil {
		return ac
	}

	selfStr, _ := self.Fields.Get("str")
	target, ok := selfStr.(*object.String)
	if !ok {
		return object.NewException("str field is not a string")
	}

	old, ok := args[0].(*object.String)
	if !ok {
		return object.NewException("replace expected a string, got %s", args[0].Type().String())
	}
	if len(old.Value) == 0 {
		return object.NewException("replace argument 2 can't be empty")
	}

	new, ok := args[1].(*object.String)
	if !ok {
		return object.NewException("replace expected a string, got %s", args[1].Type().String())
	}

	n, ok := args[2].(*object.Integer)
	if !ok {
		return object.NewException("replace expected an integer, got %s", args[2].Type().String())
	}

	return object.MakeStringObj(strings.Replace(target.String(), old.String(), new.String(), int(n.Value)))
}
