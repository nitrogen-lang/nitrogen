# Nitrogen Standard Library Documentation

The Nitrogen Standard Library (NSL) is split up into global and imported functions.

Global functions are ones implemented directly in the interpreter and are available
to all code at all times. They can not be reassigned directly but can be overshadowed
by another variable.

Imported library modules are implemented in Nitrogen and require being imported before
they can be used. Each module exports its public api for use. Imports follow the same
rules as user-defined modules. The NSL directory must be in the module import paths.
This does allow alternate NSL implementations to be used instead of the default. So
long as the public APIs match, no code would need to be updated. This can also allow
for stripping the library of extraneous modules when not needed.

- [Global](global): Globally defined NSL functions.
- [Imported](imported): NSL modules require import before use.
