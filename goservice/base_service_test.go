package goservice_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/jfenske89/go-service/goservice"
	"github.com/stretchr/testify/assert"
)

func TestShutdownHandlers(t *testing.T) {
	shutdownVerification := make(map[string]bool)
	service := goservice.BaseService{}

	for i := 1; i <= 3; i++ {
		key := strconv.Itoa(i)
		service.RegisterShutdownHandler(func(_ context.Context) error {
			shutdownVerification[key] = true
			return nil
		})
	}

	err := service.Shutdown(context.Background())
	assert.NoError(t, err)

	assert.Len(t, shutdownVerification, 3)
	assert.True(t, shutdownVerification["1"])
	assert.True(t, shutdownVerification["2"])
	assert.True(t, shutdownVerification["3"])
}
