package goservice

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/sourcegraph/conc/pool"
)

// ServiceLogic a function for executing logic
type ServiceLogic func(context.Context) error

// Service an interface to a base service
type Service interface {
	// Run executes the main service logic
	Run(ServiceLogic) error

	// RunWithContext executes the main service logic with a parent context
	RunWithContext(context.Context, ServiceLogic) error

	// Shutdown cancels the service context which triggers graceful shutdown
	Shutdown(context.Context) error

	// RegisterShutdownHandler queues a function to run during graceful shutdown
	RegisterShutdownHandler(ServiceLogic)
}

type serviceImpl struct {
	// shutdownDeadline is the maximum time allowed for graceful shutdown
	shutdownDeadline time.Duration

	// shutdownHandlers are functions to run during graceful shutdown
	shutdownHandlers []ServiceLogic

	// cancelFunc can be used to cancel the service context and trigger graceful shutdown
	cancelFunc *context.CancelFunc

	// mutex is used to for concurrent writes to shutdownHandlers
	mutex sync.Mutex
}

func NewService() Service {
	return &serviceImpl{
		// default shutdown deadline is 30 seconds
		shutdownDeadline: 30 * time.Second,
	}
}

func NewServiceWithShutdownDeadline(deadline time.Duration) Service {
	if deadline <= 0 {
		// use the default service if the deadline is invalid
		return NewService()
	}

	return &serviceImpl{
		shutdownDeadline: deadline,
	}
}

// Run executes the main service logic
func (s *serviceImpl) Run(logic ServiceLogic) error {
	return s.RunWithContext(context.Background(), logic)
}

// RunWithContext executes the main service logic with a parent context
func (s *serviceImpl) RunWithContext(parentContext context.Context, logic ServiceLogic) error {
	c := make(chan os.Signal, syscall.SIGTERM)
	ctx, cancel := signal.NotifyContext(
		parentContext,
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	s.cancelFunc = &cancel

	defer func() {
		cancel()
		close(c)
	}()

	// run the service logic...
	var logicError error
	go func() {
		defer cancel()

		logicError = logic(ctx)
	}()

	// wait until the service logic finishes (or a shutdown is initiated)
	<-ctx.Done()

	// always initiate graceful shutdown
	if gracefulShutdownError := s.gracefulShutdown(); logicError != nil {
		// prefer the logic error when present
		return logicError
	} else if gracefulShutdownError != nil {
		// return the shutdown error if that failed
		return gracefulShutdownError
	}

	return nil
}

func (s *serviceImpl) gracefulShutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownDeadline)

	defer cancel()

	var err error
	var finished bool

	go func() {
		defer cancel()

		p := pool.New().WithContext(ctx)

		for i := len(s.shutdownHandlers) - 1; i >= 0; i-- {
			handler := s.shutdownHandlers[i]

			p.Go(func(pctx context.Context) error {
				return handler(pctx)
			})
		}

		err = p.Wait()
		finished = true
	}()

	<-ctx.Done()

	if !finished {
		err = fmt.Errorf("graceful shutdown deadline exceeded")
	}

	return err
}

// Shutdown executes shutdown functions and exits
func (s *serviceImpl) Shutdown(ctx context.Context) error {
	if s.cancelFunc == nil {
		return fmt.Errorf("service not running")
	}

	(*s.cancelFunc)()

	return nil
}

// RegisterShutdownHandler adds a function to run during graceful shutdown
func (s *serviceImpl) RegisterShutdownHandler(logic ServiceLogic) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.shutdownHandlers = append(s.shutdownHandlers, logic)
}
