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
	IService

	// ShutdownHandlers these run in descending order before shutdown
	ShutdownHandlers []func(ctx context.Context) error
}

func NewBaseService() *BaseService {
	return &BaseService{}
}

// Run is the main business logic with graceful shutdown handling
func (s *BaseService) Run(logic func(ctx context.Context) error) error {
	c := make(chan os.Signal, syscall.SIGTERM)
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	var logicError error
	go func() {
		logicError = logic(ctx)
		cancel()
	}()

	defer func() {
		cancel()
		close(c)
	}()

	<-ctx.Done()

	ctxShutdown, ctxShutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() {
		ctxShutdownCancel()
	}()

	var err error
	go func() {
		err = s.Shutdown(ctxShutdown)
		ctxShutdownCancel()
	}()

	for {
		select {
		case <-ctxShutdown.Done():
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
	}
}

// Shutdown will attempt to run registered shutdown handlers
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

// RegisterShutdownHandler registers shutdown logic that will be invoked asynchronously during shutdown
func (s *BaseService) RegisterShutdownHandler(logic func(ctx context.Context) error) {
	s.ShutdownHandlers = append(s.ShutdownHandlers, logic)
}
