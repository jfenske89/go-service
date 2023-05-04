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
	return &GenericService{BaseService: goservice.NewBaseService()}
}

func main() {
	//
	// Build a GenericService or create your own custom service for more advanced use cases
	app := NewGenericService()

	//
	// Register shutdown handlers
	app.RegisterShutdownHandler(func(ctx context.Context) error {
		// Write your own graceful shutdown logic in here
		fmt.Println("Shutting down...")
		time.Sleep(2 * time.Second)
		fmt.Println("OK")
		return nil
	})

	//
	// Run main service logic
	if err := app.Run(func(ctx context.Context) error {
		//
		// Write your logic here, for example some kind of server
		fmt.Println("service logic...")
		time.Sleep(3 * time.Second)

		//
		// Return at any time to initiate shutdown logic, any error will be returned by the Run func
		return nil
	}); err != nil {
		//
		// Handle errors here, for example with logging
		panic("error running service: " + err.Error())
	}
}
