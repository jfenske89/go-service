package main

import (
	"context"
	"fmt"
	"time"

	"github.com/jfenske89/go-service/goservice"
)

// GenericService embed the base service and override any methods as needed
type GenericService struct {
	goservice.Service
}

func NewGenericService() goservice.Service {
	return &GenericService{
		// Create a new service with a 30 second shutdown deadline.
		// Kubernetes will send a SIGTERM to the process when it's time to stop.
		// Generally you 30 seconds to stop before it sends a SIGKILL, but this
		// can be configured to support other types of environments.
		goservice.NewServiceWithShutdownDeadline(10 * time.Second),
	}
}

func main() {
	// Optionally create a parent context
	globalContext := context.Background()

	//
	// Build your custom service
	service := NewGenericService()

	//
	// Run main service logic
	if err := service.RunWithContext(globalContext, func(ctx context.Context) error {
		// Connect to databases, message queues, etc.
		// ...

		// Configure graceful shutdown. For example: wait for messages to be processed, close connections, etc.
		service.RegisterShutdownHandler(func(ctx context.Context) error {
			fmt.Println("disconnecting...")

			// an error will be returned to the caller if a shutdown handler takes longer than the deadline
			// time.Sleep(15 * time.Second)
			return nil
		})

		// Write your logic here, for example some kind of server or message processor
		fmt.Println("running service logic now...")
		time.Sleep(2 * time.Second)

		// Return at any time to initiate shutdown logic (errors are returned to the caller)
		return nil
	}); err != nil {
		// Handle any error from the service or graceful shutdown logic
		panic("error running service: " + err.Error())
	}
}
