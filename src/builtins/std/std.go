package std

import (
	"github.com/nitrogen-lang/nitrogen/src/builtins/file"
	"github.com/nitrogen-lang/nitrogen/src/builtins/filepath"
	"github.com/nitrogen-lang/nitrogen/src/builtins/opbuf"
	"github.com/nitrogen-lang/nitrogen/src/builtins/os"
	"github.com/nitrogen-lang/nitrogen/src/builtins/runtime"
	nstring "github.com/nitrogen-lang/nitrogen/src/builtins/string"
	"github.com/nitrogen-lang/nitrogen/src/builtins/time"
	"github.com/nitrogen-lang/nitrogen/src/object"
	"github.com/nitrogen-lang/nitrogen/src/vm"
)

func init() {
	vm.RegisterModule("std", &object.Module{
		Name: "std",
		Vars: map[string]object.Object{
			"os":       os.Init(),
			"time":     time.Init(),
			"string":   nstring.Init(),
			"runtime":  runtime.Init(),
			"opbuf":    opbuf.Init(),
			"filepath": filepath.Init(),
			"file":     file.Init(),
		},
	})
}
