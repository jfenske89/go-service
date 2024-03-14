package goservice_test

import (
	"context"
	"errors"
	"strconv"
	"sync"
	"testing"
	"time"

	goservice "github.com/jfenske89/go-service"
	"github.com/stretchr/testify/assert"
)

func TestRunWithoutError(t *testing.T) {
	concResults := &sync.Map{}
	service := goservice.NewService()
	registerHandlers(service, concResults)

	assert.NoError(t, service.Run(func(_ context.Context) error {
		return nil
	}))

	verifyShutdownHandlers(t, concResults)
}

func TestRunWithError(t *testing.T) {
	concResults := &sync.Map{}
	service := goservice.NewService()
	registerHandlers(service, concResults)

	assert.Error(t, service.Run(func(_ context.Context) error {
		return errors.New("test")
	}))

	verifyShutdownHandlers(t, concResults)
}

func TestRunWithShutdownDeadline(t *testing.T) {
	service := goservice.NewServiceWithShutdownDeadline(1 * time.Millisecond)

	service.RegisterShutdownHandler(func(_ context.Context) error {
		<-time.After(2 * time.Millisecond)
		return nil
	})

	assert.Error(t, service.Run(func(_ context.Context) error {
		return nil
	}))

	assert.Error(t, service.Run(func(_ context.Context) error {
		return nil
	}))
}

func registerHandlers(service goservice.Service, concResults *sync.Map) {
	for i := 1; i <= 3; i++ {
		key := strconv.Itoa(i)
		service.RegisterShutdownHandler(func(_ context.Context) error {
			concResults.Store(key, true)
			return nil
		})
	}
}

func verifyShutdownHandlers(t *testing.T, concResults *sync.Map) {
	results := make(map[string]bool, 0)
	concResults.Range(func(key, value interface{}) bool {
		results[key.(string)] = value.(bool)
		return true
	})

	assert.Len(t, results, 3)
	assert.True(t, results["1"])
	assert.True(t, results["2"])
	assert.True(t, results["3"])
}
