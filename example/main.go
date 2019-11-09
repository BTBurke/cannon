package main

import (
	"context"
	"time"

	"github.com/BTBurke/cannon"
	"go.uber.org/zap"
)

func main() {
	// create a new zap logger with the cannon.NewLogger() option that keeps track of all
	// logging calls throughout the request and logs a wide cannonical log line at the end of the request
	log, _ := zap.NewDevelopment(cannon.NewLogger())
	start := time.Now()

	// you can log as normal using any zap methods, such as adding common logging fields
	logger := log.With(
		zap.String("request_id", "001"),
	)

	// you can add additional fields at each logging call
	logger.Info("auth success", zap.String("auth_role", "user_rw"))

	// the logger can be passed along in a context to handlers and other services
	ctx := cannon.ContextLogger(context.Background(), logger)
	requestHandler(ctx)

	// when finished with this request, call cannon.Emit (with optional additional fields) to log
	// a single wide log line with every field added throughout the entire request
	cannon.Emit(logger, zap.Duration("request_duration", time.Now().Sub(start)))

	// cannonical log lines make it easy to gather all of the relevant context for each request in one place
	// and allow you to aggregate statistics across requests for a better view of how your application is performing

	// Output:
	// 2019-11-08T21:46:15.377-0500    INFO    example/main.go:18      auth success    {"request_id": "001", "auth_role": "user_rw"}
	// 2019-11-08T21:46:15.377-0500    INFO    example/main.go:26      started request handler {"request_id": "001", "service": "defaultHandler"}
	// 2019-11-08T21:46:15.377-0500    INFO    example/main.go:27      updated user password   {"request_id": "001", "db_upsert": "300ms"}
	// 2019-11-08T21:46:15.377-0500    INFO    cannonical_log_line     {"request_id": "001", "auth_role": "user_rw", "service": "defaultHandler", "db_upsert": "300ms", "request_duration": "104.36µs"}
}

func requestHandler(ctx context.Context) {
	logger, _ := cannon.LoggerFromCtx(ctx)
	logger.Info("started request handler", zap.String("service", "defaultHandler"))
	logger.Info("updated user password", zap.Duration("db_upsert", 300*time.Millisecond))

}
