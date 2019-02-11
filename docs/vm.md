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

### LOAD\_CONST

### STORE\_CONST

### LOAD\_FAST

### STORE\_FAST

### DELETE\_FAST

### DEFINE

### LOAD\_GLOBAL

### STORE\_GLOBAL

### LOAD\_INDEX

### STORE\_INDEX

### BINARY\_ADD

### BINARY\_SUB

### BINARY\_MUL

### BINARY\_DIVIDE

### BINARY\_MOD

### BINARY\_SHIFTL

### BINARY\_SHIFTR

### BINARY\_AND

### BINARY\_OR

### BINARY\_NOT

### BINARY\_ANDNOT

### UNARY\_NEG

### UNARY\_NOT

### COMPARE

### CALL

### RETURN

### POP

### MAKE\_ARRAY

### MAKE\_MAP

### MAKE\_FUNCTION

### POP\_JUMP\_IF\_TRUE

### POP\_JUMP\_IF\_FALSE

### JUMP\_IF\_TRUE\_OR\_POP

### JUMP\_IF\_FALSE\_OR\_POP

### JUMP\_ABSOLUTE

### JUMP\_FORWARD

### OPEN\_SCOPE

### CLOSE\_SCOPE

### END\_BLOCK

### START\_LOOP

### CONTINUE

### NEXT\_ITER

### BREAK

### START\_TRY

### THROW
