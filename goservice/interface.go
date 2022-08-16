package goservice

import (
	"context"
)

type IService interface {
	Run(logic func(ctx context.Context) error) error
	Shutdown(ctx context.Context) error
	RegisterShutdownHandler(logic func(ctx context.Context) error)
}
