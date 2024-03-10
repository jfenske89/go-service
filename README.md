# go-service

A base Go service implementation with graceful shutdown handling.

## Running

Simply execute your logic using the `Run` or `RunWithContext` function.

```go
// Run executes the main service logic
Run(logic func(context.Context) error) error

// RunWithContext executes the main service logic with a parent context
RunWithContext(context.Context, logic func(context.Context) error) error
```

## Graceful shutdown

Define graceful shutdown logic with `RegisterShutdownHandler`.

For example: flush logs, close connections, wait for active work to finish, etc...

```go
// Shutdown executes shutdown functions and exits
RegisterShutdownHandler(logic func(context.Context) error)
```

These functions are executed in parallel before the application exits.

Shutdown has a 30 second deadline by default. This can be customized.

## Examples

See [./examples/basic/main.go](./examples/basic/main.go) for a basic example.