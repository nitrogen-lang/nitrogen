# Language Modules

The Nitrogen interpreter supports dynamically linking to shared Go plugins to extend the functionality of the interpreter
by providing additional userland functions. **Modules are only supported on Linux.**

## Using Modules

Modules can be imported in two ways. The first is when the interpreter starts but before a script is executed. In this case,
modules should be in a single folder and the nitrogen binary needs to be given the `-modules` flag with the path to the modules
directory. Any file with the extension `.so` is loaded as a module. The modules folder can also be given with the environment
variable `NITROGEN_MODULES`.

The second way is in a script. The `module(filepath)` function will import a shared library module which is then available for any script.
The module is not bound to a variable. Any functions created from the module are in the global scope so using `module` inside
a block scope has no effect.
