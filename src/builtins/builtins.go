package builtins

import (
	// This file imports all the different builtin modules. Each module is its own package for simplicity
	// and separation of concerns.
	_ "github.com/nitrogen-lang/nitrogen/src/builtins/classes"
	_ "github.com/nitrogen-lang/nitrogen/src/builtins/collections"
	_ "github.com/nitrogen-lang/nitrogen/src/builtins/dis"
	_ "github.com/nitrogen-lang/nitrogen/src/builtins/imports"
	_ "github.com/nitrogen-lang/nitrogen/src/builtins/io"
	_ "github.com/nitrogen-lang/nitrogen/src/builtins/string"
	_ "github.com/nitrogen-lang/nitrogen/src/builtins/typing"
)
