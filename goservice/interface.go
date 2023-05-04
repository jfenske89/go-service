package goservice

import (
	"context"
)

type BaseService interface {
	// Run executes the main service logic
	Run(logic func(ctx context.Context) error) error

	// RunWithContext executes the main service logic with a parent context
	RunWithContext(parentContext context.Context, logic func(ctx context.Context) error) error

	// Shutdown executes shutdown functions and exits
	Shutdown(ctx context.Context) error

	// RegisterShutdownHandler adds a function to run during graceful shutdown
	RegisterShutdownHandler(logic func(ctx context.Context) error)
}
