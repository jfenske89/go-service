package service

import (
	"context"
	"github.com/jinzhu/configor"
)

type IService interface {
	Configure(config interface{}, files ...string) error
	ConfigureAdvanced(config interface{}, configuration *configor.Configor, files ...string) error
	GetConfig() interface{}
	Run(logic func(config interface{}) error) error
	Shutdown(ctx context.Context) error
	RegisterShutdownHandler(logic func(config interface{}, ctx context.Context) error)
}
