# Runtime VM

This implementation of Nitrogen runs on a bytecode virtual machine core called
Elemental. In general it's faster and more efficient than the old interpreter.
With the VM, source code goes through a full compile stage before being executed
including when using the REPL. This implementation does not have a JIT compiler
and does not compile to machine assembly but instead to a higher level
assembly-like bytecode.

## Opcodes

These are all the opcodes used in this implementation.

### BINARY_ADD

### BINARY_AND

### BINARY_ANDNOT

### BINARY_DIVIDE

### BINARY_MOD

### BINARY_MUL

### BINARY_NOT

### BINARY_OR

### BINARY_SHIFTL

### BINARY_SHIFTR

### BINARY_SUB

### BREAK

### BUILD_CLASS

### CALL

### COMPARE

### CONTINUE

### DEFINE

### DELETE_FAST

### DUP

### END_BLOCK

### GET_ITER

### IMPLEMENTS

### IMPORT

### JUMP_ABSOLUTE

### JUMP_FORWARD

### JUMP_IF_FALSE_OR_POP

### JUMP_IF_TRUE_OR_POP

### LOAD_ATTRIBUTE

### LOAD_CONST

### LOAD_FAST

### LOAD_GLOBAL

### LOAD_INDEX

### MAKE_ARRAY

### MAKE_FUNCTION

### MAKE_INSTANCE

### MAKE_MAP

### NEXT_ITER

### NOOP

### POP

### POP_JUMP_IF_FALSE

### POP_JUMP_IF_TRUE

### RECOVER

### RETURN

### START_BLOCK

### START_LOOP

### STORE_ATTRIBUTE

### STORE_CONST

### STORE_FAST

### STORE_GLOBAL

### STORE_INDEX

### UNARY_NEG

### UNARY_NOT
