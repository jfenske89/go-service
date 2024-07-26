# Archived project

This repository is archived and the code was moved here: https://github.com/jfenske89/gobasrv

# go-service

A base Go service implementation, with graceful shutdown handling.

## Running

Once your service has been initialized, execute your logic using the `Run` or `RunWithContext` function.

```go
// Run executes the main service logic
Run(logic func(context.Context) error) error

// RunWithContext executes the main service logic with a parent context
RunWithContext(context.Context, logic func(context.Context) error) error
```

## Graceful shutdown

Define graceful shutdown routines using the `RegisterShutdownHandler` function.

```go
// Shutdown executes shutdown functions and exits
RegisterShutdownHandler(logic func(context.Context) error)
```

These functions are executed in parallel before the application exits.

The context passed to these handlers has a 30 second deadline. So execution cannot exceed that time.

## Examples

See [./examples/basic/main.go](./examples/basic/main.go) for a basic example.
