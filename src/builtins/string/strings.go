package string

import (
	"strings"

	"github.com/nitrogen-lang/nitrogen/src/moduleutils"
	"github.com/nitrogen-lang/nitrogen/src/object"
	"github.com/nitrogen-lang/nitrogen/src/vm"
)

func strSplitN(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	if ac := moduleutils.CheckArgs("strSplitN", 3, args...); ac != nil {
		return ac
	}

	target, ok := args[0].(*object.String)
	if !ok {
		return object.NewException("splitN expected a string, got %s", args[0].Type().String())
	}

	sep, ok := args[1].(*object.String)
	if !ok {
		return object.NewException("splitN expected a string, got %s", args[1].Type().String())
	}

	count, ok := args[2].(*object.Integer)
	if !ok {
		return object.NewException("splitN expected an int, got %s", args[2].Type().String())
	}

	return object.MakeStringArray(strings.SplitN(target.String(), sep.String(), int(count.Value)))
}

func strSplit(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	if ac := moduleutils.CheckArgs("split", 2, args...); ac != nil {
		return ac
	}

	return strSplitN(interpreter, env, args[0], args[1], object.MakeIntObj(-1))
}

func strTrim(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	target, ok := args[0].(*object.String)
	if !ok {
		return object.NewException("trim expected a string, got %s", args[0].Type().String())
	}

	return object.MakeStringObj(strings.TrimSpace(target.String()))
}

func strDedup(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	if ac := moduleutils.CheckArgs("strDedup", 2, args...); ac != nil {
		return ac
	}

	target, ok := args[0].(*object.String)
	if !ok {
		return object.NewException("strDedup expected a string, got %s", args[0].Type().String())
	}

	dedup, ok := args[1].(*object.String)
	if !ok {
		return object.NewException("strDedup expected a string, got %s", args[1].Type().String())
	}

	if len(dedup.Value) != 1 {
		return object.NewException("Dedup string must be one byte")
	}

	return object.MakeStringObj(dedupString(target.Value, dedup.Value[0]))
}

func strFormat(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	target, ok := args[0].(*object.String)
	if !ok {
		return object.NewException("format expected a string, got %s", args[0].Type().String())
	}

	t := string(target.Value)

	for _, arg := range args[1:] {
		if !strings.Contains(t, "{}") {
			break
		}

		s := objectToString(arg, interpreter.(*vm.VirtualMachine))
		t = strings.Replace(t, "{}", s, 1)
	}

	return object.MakeStringObj(t)
}

func strContains(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	if ac := moduleutils.CheckArgs("contains", 2, args...); ac != nil {
		return ac
	}

	target, ok := args[0].(*object.String)
	if !ok {
		return object.NewException("contains expected a string, got %s", args[0].Type().String())
	}

	sub, ok := args[1].(*object.String)
	if !ok {
		return object.NewException("contains expected a string, got %s", args[1].Type().String())
	}

	return object.NativeBoolToBooleanObj(strings.Contains(target.String(), sub.String()))
}

func strCount(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	if ac := moduleutils.CheckArgs("count", 2, args...); ac != nil {
		return ac
	}

	target, ok := args[0].(*object.String)
	if !ok {
		return object.NewException("count expected a string, got %s", args[0].Type().String())
	}

	sub, ok := args[1].(*object.String)
	if !ok {
		return object.NewException("count expected a string, got %s", args[1].Type().String())
	}
	if len(sub.Value) == 0 {
		return object.NewException("count argument 2 can't be empty")
	}

	return object.MakeIntObj(int64(strings.Count(target.String(), sub.String())))
}

func strHasPrefix(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	if ac := moduleutils.CheckArgs("hasPrefix", 2, args...); ac != nil {
		return ac
	}

	target, ok := args[0].(*object.String)
	if !ok {
		return object.NewException("hasPrefix expected a string, got %s", args[0].Type().String())
	}

	prefix, ok := args[1].(*object.String)
	if !ok {
		return object.NewException("hasPrefix expected a string, got %s", args[1].Type().String())
	}

	return object.NativeBoolToBooleanObj(strings.HasPrefix(target.String(), prefix.String()))
}

func strHasSuffix(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	if ac := moduleutils.CheckArgs("hasSuffix", 2, args...); ac != nil {
		return ac
	}

	target, ok := args[0].(*object.String)
	if !ok {
		return object.NewException("hasSuffix expected a string, got %s", args[0].Type().String())
	}

	prefix, ok := args[1].(*object.String)
	if !ok {
		return object.NewException("hasSuffix expected a string, got %s", args[1].Type().String())
	}

	return object.NativeBoolToBooleanObj(strings.HasSuffix(target.String(), prefix.String()))
}

func strReplace(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	if ac := moduleutils.CheckArgs("replace", 4, args...); ac != nil {
		return ac
	}

	target, ok := args[0].(*object.String)
	if !ok {
		return object.NewException("replace expected a string, got %s", args[0].Type().String())
	}

	old, ok := args[1].(*object.String)
	if !ok {
		return object.NewException("replace expected a string, got %s", args[1].Type().String())
	}
	if len(old.Value) == 0 {
		return object.NewException("replace argument 2 can't be empty")
	}

	new, ok := args[2].(*object.String)
	if !ok {
		return object.NewException("replace expected a string, got %s", args[2].Type().String())
	}

	n, ok := args[3].(*object.Integer)
	if !ok {
		return object.NewException("replace expected an integer, got %s", args[3].Type().String())
	}

	return object.MakeStringObj(strings.Replace(target.String(), old.String(), new.String(), int(n.Value)))
}
