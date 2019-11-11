package main

import (
	"os"

	"github.com/BTBurke/cannon"
	"go.uber.org/zap"
)

func main() {
	// create a LogFactoryFunc that checks to see if you are in a dev or prod environment
	factoryLogger := cannon.LogFactoryFunc(func() (*zap.Logger, error) {
		switch os.Getenv("ENVIRONMENT") {
		case "PRODUCTION":
			return cannon.NewProduction()
		default:
			return cannon.NewDevelopment()
		}
	})
	// register the global factory function
	cannon.RegisterFactory(factoryLogger)
	// create a new logger based on the factory function
	logger, _ := cannon.NewFactoryLogger()
	// log stuff and emit the cannnonical log line when you are done with the request
	logger.Info("processing request", zap.String("logger", "factoryLogger"))
	cannon.Emit(logger, zap.String("status", "success"))

	// Output:
	// 2019-11-11T11:20:28.477-0500    INFO    example/factory_example.go:25   processing request      {"logger": "factoryLogger"}
	// 2019-11-11T11:20:28.477-0500    INFO    cannonical_log_line     {"logger": "factoryLogger", "status": "success"}
}
