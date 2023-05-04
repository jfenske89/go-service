package goservice

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

type BaseService struct {
	ShutdownHandlers []func(ctx context.Context) error
}

// Run executes the main service logic
func (s *BaseService) Run(logic func(ctx context.Context) error) error {
	return s.RunWithContext(context.Background(), logic)
}

// RunWithContext executes the main service logic with a parent context
func (s *BaseService) RunWithContext(parentContext context.Context, logic func(ctx context.Context) error) error {
	c := make(chan os.Signal, syscall.SIGTERM)
	ctx, cancel := signal.NotifyContext(parentContext, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	// Run the service logic...
	var logicError error
	go func() {
		logicError = logic(ctx)
		cancel()
	}()

	defer func() {
		cancel()
		close(c)
	}()

	// Wait until the service logic finishes (or a shutdown is initiated)
	<-ctx.Done()

	// Execute graceful shutdown routines
	ctxShutdown, ctxShutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() {
		ctxShutdownCancel()
	}()

	var err error
	go func() {
		err = s.Shutdown(ctxShutdown)
		ctxShutdownCancel()
	}()

	// Wait for the graceful shutdown logic to complete (or timeout)
	<-ctxShutdown.Done()
	if logicError != nil {
		return logicError
	} else if err != nil {
		return err
	}

	if ctxErr := ctxShutdown.Err(); ctxErr != nil && ctxErr != context.Canceled {
		return ctxErr
	}
	return nil
}

// Shutdown executes shutdown functions and exits
func (s *BaseService) Shutdown(ctx context.Context) error {
	g, gctx := errgroup.WithContext(ctx)
	for i := len(s.ShutdownHandlers) - 1; i >= 0; i-- {
		handler := s.ShutdownHandlers[i]
		g.Go(func() error {
			return handler(gctx)
		})
	}
	return g.Wait()
}

// RegisterShutdownHandler adds a function to run during graceful shutdown
func (s *BaseService) RegisterShutdownHandler(logic func(ctx context.Context) error) {
	s.ShutdownHandlers = append(s.ShutdownHandlers, logic)
}
