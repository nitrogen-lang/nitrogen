# Importing Packages

Nitrogen application are split into separate packages for separation of functionality and
better maintainability. For in-depth documentation on imports and packages, please read the
[packages language docs](../../language/packages.md).

## import "module/path"[ as var]

Import a module or package into the application. Optionally assigning it to `var` otherwise it's assigned
to the last segment of the module name.

## use module.attribute[ as var]

The use statement is a convenience statement to assign an attribute to a variable. The same can be achieved
with `let` or `const`, but `use` conveys better meaning as to the intent. `use` makes it so the full module
name doesn't need to be used when accessing a module's properties.

### Example

```
import "string"

use string.String // Assigned to the variable 'String'

/*
 * The above use statement would be equivalent to:
 * const String = string.String
 */

// Now the module name isn't needed
println((new String("Hello, {}")).format("World!"))
```

## modulesSupported(): bool

Returns if the platform and build supports dynamic binary modules.
