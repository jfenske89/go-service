package goservice

import (
	"context"
)

type IService interface {
	Configure(config interface{}, files ...string) error
	SetConfig(config interface{})
	GetConfig() interface{}
	Run(logic func(config interface{}) error) error
	Shutdown(ctx context.Context) error
	RegisterShutdownHandler(logic func(ctx context.Context) error)
}
