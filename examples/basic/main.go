package main

import (
	"context"
	"fmt"
	"time"

	"github.com/jfenske89/go-service/goservice"
)

// GenericService example service that inherits the base version (handle your configs, etc...)
type GenericService struct {
	goservice.BaseService
}

func NewGenericService() *GenericService {
	return &GenericService{}
}

func main() {
	//
	// Build your custom service
	service := NewGenericService()

	//
	// Register global shutdown handlers as needed
	service.RegisterShutdownHandler(func(ctx context.Context) error {
		// Write your own graceful shutdown logic in here
		fmt.Println("Shutting down...")
		time.Sleep(2 * time.Second)
		fmt.Println("OK")
		return nil
	})

	//
	// Run main service logic
	if err := service.Run(func(ctx context.Context) error {
		//
		// You could reference the service within this logic, for example to register conditional shutdown handlers
		service.RegisterShutdownHandler(func(ctx context.Context) error {
			fmt.Println("other shutdown handler...")
			return nil
		})

		//
		// Write your logic here, for example some kind of server
		fmt.Println("service logic...")
		time.Sleep(3 * time.Second)

		if false {
			//
			// Any error returned by this logic, will also be returned by the Run function
			return fmt.Errorf("oh no")
		}

		//
		// Return at any time to initiate shutdown logic
		return nil
	}); err != nil {
		//
		// Handle errors here (for example the "oh no" error above)
		panic("error running service: " + err.Error())
	}
}
