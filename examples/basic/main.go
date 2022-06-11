package main

import (
	"context"
	"fmt"
	"github.com/jfenske89/go-service/service"
	"time"
)

type Config struct {
	Example string `default:"example" env:"EXAMPLE"`
}

func main() {
	// Build a GenericService or create your own custom service for more advanced use cases
	app := service.NewGenericService()
	err := app.Configure(&Config{})
	if err != nil {
		// Handle errors here, for example with logging
		panic("configuration error: " + err.Error())
	}

	//
	// Register shutdown handlers
	app.RegisterShutdownHandler(func(config interface{}, ctx context.Context) error {
		// Write your own graceful shutdown logic in here
		fmt.Println("Shutting down...")
		time.Sleep(2 * time.Second)
		return nil
	})

	//
	// Run main service logic
	err = app.Run(func(configIface interface{}) error {
		// Convert your custom configuration object back into the struct to use more easily
		config := *configIface.(*Config)

		// Write your logic here, for example some kind of server
		fmt.Printf("%s service...\n", config.Example)
		time.Sleep(3 * time.Second)
		return nil
	})

	if err != nil {
		// Handle errors here, for example with logging
		panic("error running service: " + err.Error())
	}
}
