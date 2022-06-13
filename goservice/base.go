package goservice

import (
	"context"
	"github.com/jinzhu/configor"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type BaseService struct {
	IService

	// Configor An optional instance of Configor for handling configuration loading
	Configor *configor.Configor

	// A generic configuration object
	Config interface{}

	// ShutdownHandlers these run in descending order before shutdown
	ShutdownHandlers []func(ctx context.Context) error
}

func NewBaseService() *BaseService {
	return &BaseService{}
}

// Configure used to load or re-load configuration using Configor
func (s *BaseService) Configure(config interface{}, files ...string) error {
	environment := os.Getenv("ENVIRONMENT")
	configorConfig := &configor.Config{}
	if environment != "" {
		configorConfig.Environment = environment
	}

	configuration := configor.New(configorConfig)
	s.Configor = configuration
	s.Config = config

	return s.Configor.Load(s.Config, files...)
}

// SetConfig set the configuration (if you don't want to use Configor)
func (s *BaseService) SetConfig(config interface{}) {
	s.Config = config
}

// GetConfig return the current configuration
func (s *BaseService) GetConfig() interface{} {
	return s.Config
}

// Run main business logic with graceful shutdown handling
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

// Shutdown attempt to run registered shutdown handlers
func (s *BaseService) Shutdown(ctx context.Context) error {
	for i := len(s.ShutdownHandlers) - 1; i >= 0; i-- {
		err := s.ShutdownHandlers[i](ctx)
		if err != nil {
			return nil
		}
	}
	return nil
}

// RegisterShutdownHandler register a shutdown handler which will all run before shutdown, in descending order
func (s *BaseService) RegisterShutdownHandler(logic func(ctx context.Context) error) {
	s.ShutdownHandlers = append(s.ShutdownHandlers, logic)
}
