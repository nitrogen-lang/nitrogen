# Runtime VM

This implementation of Nitrogen runs on a bytecode virtual machine core called
Elemental. In general it's faster and more efficient than the old interpreter.
With the VM, source code goes through a full compile stage before being executed
including when using the REPL. This implementation does not have a JIT compiler
and does not compile to machine assembly but instead to a higher level
assembly-like bytecode.

## Opcodes

These are all the opcodes used in this implementation.

### NOOP

### LOAD_CONST

### LOAD_FAST

### STORE_FAST

### DELETE_FAST

### DEFINE

### LOAD_GLOBAL

### STORE_GLOBAL

### LOAD_INDEX

### STORE_INDEX

### LOAD_ATTRIBUTE

### STORE_ATTRIBUTE

### BINARY_ADD

### BINARY_SUB

### BINARY_MUL

### BINARY_DIVIDE

### BINARY_MOD

### BINARY_SHIFTL

### BINARY_SHIFTR

### BINARY_AND

### BINARY_OR

### BINARY_NOT

### BINARY_ANDNOT

### IMPLEMENTS

### UNARY_NEG

### UNARY_NOT

### COMPARE

### CALL

### RETURN

### POP

### MAKE_ARRAY

### MAKE_MAP

### MAKE_FUNCTION

### POP_JUMP_IF_TRUE

### POP_JUMP_IF_FALSE

### JUMP_IF_TRUE_OR_POP

### JUMP_IF_FALSE_OR_POP

### JUMP_ABSOLUTE

### JUMP_FORWARD

### START_BLOCK

### END_BLOCK

### START_LOOP

### CONTINUE

### NEXT_ITER

### BREAK

### RECOVER

### BUILD_CLASS

### MAKE_INSTANCE

### IMPORT

### DUP

### GET_ITER

### BREAKPOINT
