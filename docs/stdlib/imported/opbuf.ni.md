# opbuf.ni

Manage the output buffer.

To use: `import 'std/opbuf'`

## start(): nil

Start output buffering. Buffering can't be nested so `start` will throw if buffering is
already started.

## stop(): nil

Stop output buffering. `stop` will throw if buffering is already stopped.

## isStarted(): bool

Returns if output buffering is running.

## clear(): nil

Clear the current buffer.

## flush(): nil

Flush current buffer to stdout.

## get(): string

Get the contents of the buffer as a string.

## stopAndGet(): string

Stop the current buffer and return the contents of the buffer as a string.
`stopAndGet` will throw if buffering is already stopped
