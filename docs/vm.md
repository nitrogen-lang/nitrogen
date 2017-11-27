# Bytecode VM

A bytecode virtual machine is software that compiles a language to a set of assembly-like instructions that run in a software
defined virtual machine. Bytecode VMs are generally much faster and more efficient than tree-walking interpreters.

This implementation of Nitrogen has both a tree-walking interpreter and a VM implementation. By default the interpreter is
used. The VM can used instead by running Nitrogen with the `-compile` flag. The flag will also make the SCGI server use
the VM when serving web pages.

## Differences

For 98% of the language, the VM will execute code exactly the same as the interpreter. But there are a few differences. Some
just need to be implemented and others will be reflected back into the original interpreter.

- For loop iterators are not longer auto-converted to an assignment. Previously `i + 1` was converted to `i += 1` at runtime.
That is no longer the case. The VM will render the expression as is and if it doesn't increment the loop counter, that's a bug
in the source, not the runtime.
- Try/catch blocks have a separate enclosed scope. Previously, variables declared in a try/catch were placed in the surrounding
scope. Now they will stay in the block scope of the try or catch. Try/catch still has an implicit return.
- Classes are currently not implemented. They will be eventually.

## Why use the VM?

Performance. The virtual machine is much faster than the tree-walking interpreter even with an extra compile step. At the moment,
the VM doesn't have full feature parity with the interpreter. But if the changes above don't affect you, just use the VM. All
standard library functions and all but one external modules are available.

## Opcodes

These are all the opcodes used in this implementation.

### NOOP

### LOAD\_CONST

### STORE\_CONST

### LOAD\_FAST

### STORE\_FAST

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

### PREPARE\_BLOCK

### END\_BLOCK

### START\_LOOP

### CONTINUE

### NEXT\_ITER

### BREAK

### START\_TRY
### THROW
