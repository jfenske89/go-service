package goservice

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type BaseService struct {
	IService

	// A generic configuration object
	Config interface{}

	// ShutdownHandlers these run in descending order before shutdown
	ShutdownHandlers []func(ctx context.Context) error
}

func NewBaseService() *BaseService {
	return &BaseService{}
}

// SetConfig will store service configuration
func (s *BaseService) SetConfig(config interface{}) {
	s.Config = config
}

// GetConfig will return service configuration
func (s *BaseService) GetConfig() interface{} {
	return s.Config
}

// Run is the main business logic with graceful shutdown handling
func (s *BaseService) Run(logic func(config interface{}) error) error {
	c := make(chan os.Signal, syscall.SIGTERM)
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	var logicError error
	go func() {
		logicError = logic(s.GetConfig())
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

	err := s.Shutdown(ctxShutdown)
	if logicError != nil {
		return logicError
	}
	return err
}

// Shutdown will attempt to run registered shutdown handlers
func (s *BaseService) Shutdown(ctx context.Context) error {
	for i := len(s.ShutdownHandlers) - 1; i >= 0; i-- {
		err := s.ShutdownHandlers[i](ctx)
		if err != nil {
			return nil
		}
	}
	return nil
}

// RegisterShutdownHandler will register a shutdown handler which will all run before shutdown, in descending order
func (s *BaseService) RegisterShutdownHandler(logic func(ctx context.Context) error) {
	s.ShutdownHandlers = append(s.ShutdownHandlers, logic)
}
