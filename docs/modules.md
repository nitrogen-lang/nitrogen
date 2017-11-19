# Language Modules

The Nitrogen interpreter supports dynamically linking to shared Go plugins to extend the functionality of the interpreter
by providing additional userland functions. **Modules are only supported on Linux.**

## Using Modules

Modules can be imported in two ways. The first is when the interpreter starts but before a script is executed. In this case,
modules should be in a single folder and the nitrogen binary needs to be given the `-modules` flag with the path to the modules
directory. Any file with the extension `.so` is loaded as a module. The modules folder can also be given with the environment
variable `NITROGEN_MODULES`. These modules can be retrieved in user code by just using the module name in `module()`. For example,
if `os.so` is loaded this way, in code the module can be used with `let os = module('os')`.

The second way is at execution time. The `module()` function can take a filepath and will attempt to load a shared module
from that path. The path is relative to the executing file, not the interpreter working directory. Modules loaded in this
way can later be retrieved using just the module name. For example if a module named `database` is at `/modules/database.so`
on disk, the first import for the module is done by call `module('/modules/database.so')` (in this case an absolute path is used).
Later, the same module can be "imported" again by simply calling `module('database')`.

A module can register global functions, create a Module object to encapsulate functionality, or even both. If a module registers
a Module object, that object will be returned with the `module()` function call. Registered global functions are available
immediately after import.
