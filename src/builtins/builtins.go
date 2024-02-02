package builtins

import (
	// This file imports all the different builtin modules. Each module is its own package for simplicity
	// and separation of concerns.
	_ "github.com/nitrogen-lang/nitrogen/src/builtins/classes"
	_ "github.com/nitrogen-lang/nitrogen/src/builtins/collections"
	_ "github.com/nitrogen-lang/nitrogen/src/builtins/errors"
	_ "github.com/nitrogen-lang/nitrogen/src/builtins/file"
	_ "github.com/nitrogen-lang/nitrogen/src/builtins/filepath"
	_ "github.com/nitrogen-lang/nitrogen/src/builtins/http"
	_ "github.com/nitrogen-lang/nitrogen/src/builtins/imports"
	_ "github.com/nitrogen-lang/nitrogen/src/builtins/io"
	_ "github.com/nitrogen-lang/nitrogen/src/builtins/opbuf"
	_ "github.com/nitrogen-lang/nitrogen/src/builtins/os"
	_ "github.com/nitrogen-lang/nitrogen/src/builtins/runtime"
	_ "github.com/nitrogen-lang/nitrogen/src/builtins/string"
	_ "github.com/nitrogen-lang/nitrogen/src/builtins/time"
	_ "github.com/nitrogen-lang/nitrogen/src/builtins/typing"
)
