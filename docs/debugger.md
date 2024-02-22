# VM Debugger

Run `nitrogen` with the `-debug` flag. In this mode, `breakpoint` statements in
a file will pause VM execution and start an interactive debugger prompt.

## Debugger Commands

- `quit`: Exit the program.
- `continue`: Continue execution until the end or the next breakpoint.
- `step`: Execute the next instruction then break again.
- `frame`: Print current frame information from VM.
- `contwithinstrs`: Continue execution until the end of the script and print all
  instruction debug information.
- `frames`: Print the current call frame stack.
- `env`: Print all defined variables in the current scope.
