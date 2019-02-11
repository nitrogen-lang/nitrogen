# Nitrogen Environment Globals

The interpreter exposes a few global variables that can be used during execution.
Nitrogen reserves all variables starting with a single underscore for interpreter
provided values. Use of variables starting with a single underscore in user code
is discouraged.

## _FILE

`_FILE` is the absolute path to the currently executing script.

## _SEARCH_PATHS

`_SEARCH_PATHS` is an array containing the paths used for import search.

## _SERVER

`_SERVER` is a map containing values given by a web server when using CGI or SCGI.
If a script is run in another way, `_SERVER` will be nil.
