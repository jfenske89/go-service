package goservice

import (
	"context"
)

// ServiceLogic a function for executing logic
type ServiceLogicFunc func(ctx context.Context) error

// Service an interface to a base service
type Service interface {
	// Run executes the main service logic
	Run(logic ServiceLogicFunc) error

	// RunWithContext executes the main service logic with a parent context
	RunWithContext(parentContext context.Context, logic ServiceLogicFunc) error

	// Shutdown executes shutdown functions and exits
	Shutdown(ctx context.Context) error

	// RegisterShutdownHandler adds a function to run during graceful shutdown
	RegisterShutdownHandler(logic ServiceLogicFunc)
}
