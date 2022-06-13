package goservice

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestShutdownHandlersRunInDescendingOrder(t *testing.T) {
	var shutdownVerification []string
	service := NewBaseService()

	service.RegisterShutdownHandler(func(ctx context.Context) error {
		shutdownVerification = append(shutdownVerification, "last")
		return nil
	})
	service.RegisterShutdownHandler(func(ctx context.Context) error {
		shutdownVerification = append(shutdownVerification, "second")
		return nil
	})
	service.RegisterShutdownHandler(func(ctx context.Context) error {
		shutdownVerification = append(shutdownVerification, "first")
		return nil
	})

	err := service.Shutdown(context.Background())
	assert.NoError(t, err)

	assert.Len(t, shutdownVerification, 3)
	assert.Equal(t, shutdownVerification[0], "first")
	assert.Equal(t, shutdownVerification[1], "second")
	assert.Equal(t, shutdownVerification[2], "last")
}
