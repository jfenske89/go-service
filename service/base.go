package service

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
	Configor         *configor.Configor
	Config           interface{}
	ShutdownHandlers []func(config interface{}, ctx context.Context) error
}

// Configure used to load or re-load configuration
func (s *BaseService) Configure(config interface{}, files ...string) error {
	environment := os.Getenv("ENVIRONMENT")
	configorConfig := &configor.Config{}
	if environment != "" {
		configorConfig.Environment = environment
	}

	configuration := configor.New(configorConfig)
	return s.ConfigureAdvanced(config, configuration, files...)
}

// ConfigureAdvanced used to load or re-load configuration with specified Configor
func (s *BaseService) ConfigureAdvanced(config interface{}, configuration *configor.Configor, files ...string) error {
	s.Configor = configuration
	s.Config = config
	return s.Configor.Load(s.Config, files...)
}

// GetConfig return the current Configor
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
	for i := 0; i < len(s.ShutdownHandlers); i++ {
		err := s.ShutdownHandlers[i](s.GetConfig(), ctx)
		if err != nil {
			return nil
		}
	}
	return nil
}

// RegisterShutdownHandler register a shutdown handler which attempts to run before shutdown
func (s *BaseService) RegisterShutdownHandler(logic func(config interface{}, ctx context.Context) error) {
	s.ShutdownHandlers = append(s.ShutdownHandlers, logic)
}
