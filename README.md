# Cannon - Cannonical Log Lines for Go

Cannon builds on top of [Uber's zap logging library](https://github.com/uber-go/zap) to facilitate canonical logging in Go.  Canonical log lines (i.e. a single summary line emitted per request per service) are a great way to assist with performance monitoring.  [Stripe's engineering blog](https://stripe.com/blog/canonical-log-lines) has a good overview of the benefits of canonical log lines.  This library plugs into zap to let you emit a canonical log line at the end of processing each request by aggregating your structured key-value log entries throughout the duration of the request.

## Canonical Logging Example

Many applications emit structured logs throughout the course of a request, such as the following:

```
[2019-03-18 22:48:32.990] Request started http_method=POST http_path=/v1/charges request_id=req_123

[2019-03-18 22:48:32.991] User authenticated auth_type=api_key key_id=mk_123 user_id=usr_123

[2019-03-18 22:48:32.992] Rate limiting ran rate_allowed=true rate_quota=100 rate_remaining=99

[2019-03-18 22:48:32.998] Charge created charge_id=ch_123 permissions_used=account_write team=acquiring

[2019-03-18 22:48:32.999] Request finished alloc_count=9123 database_queries=34 duration=0.009 http_status=200
```

Canonical logging summarizes the structured values above at the end of the request in a wide log line:

```
[2019-03-18 22:48:32.999] canonical-log-line alloc_count=9123 auth_type=api_key database_queries=34 duration=0.009 http_method=POST http_path=/v1/charges http_status=200 key_id=mk_123 permissions_used=account_write rate_allowed=true rate_quota=100 rate_remaining=99 request_id=req_123 team=acquiring user_id=usr_123
```

Using these per-request summaries allow you easily calculate important metrics about your application's performance in your log aggregator of choice.  For example, you could generate stats about rate limited requests by querying your canonical log lines for those that were rate limited:

```
canonical-log-line rate_allowed=false | stats count by user_id
```

Many more examples are listed on Stripe's blog post linked above.

# Using Cannon

Cannon wraps Zap's logging core to aggregate your structured log entries.  It's not another logging library, but rather a plugin for Zap that adds a bit of state to your normal logging calls.  You get the same configurability and performance of Zap (with a small bit of overhead) since Cannon delegates the log processing to Zap.

To use cannon, you log as you normally would with Zap, then call `cannon.Emit` when your request is finished:

```go
    // create a new zap logger with the default zap.NewDevelopment configuration that keeps track of all
	// logging calls throughout the request and logs a wide cannonical log line at the end of the request
	log, _ := cannon.NewDevelopment()
	start := time.Now()

	// you can log as normal using any zap methods, such as adding common logging fields
	logger := log.With(
		zap.String("request_id", "001"),
	)

	// you can add additional fields at each logging call
	logger.Info("auth success", zap.String("auth_role", "user_rw"))

	// the logger can be passed along in a context to handlers and other services
	ctx := cannon.CtxLogger(context.Background(), logger)
	requestHandler(ctx, req, resp)

	// when finished with this request, call cannon.Emit (with optional additional fields) to log
	// a single wide log line with every field added throughout the entire request

	// cannonical log lines make it easy to gather all of the relevant context for each request in one place
	// and allow you to aggregate statistics across requests for a better view of how your application is performing
	cannon.Emit(logger, zap.Duration("request_duration", time.Now().Sub(start)))

	// Output:
	// 2019-11-11T11:11:09.567-0500    INFO    example/basic.go:23     auth success    {"request_id": "001", "auth_role": "user_rw"}
	// 2019-11-11T11:11:09.567-0500    INFO    example/basic.go:45     started request handler {"request_id": "001", "service": "defaultHandler"}
	// 2019-11-11T11:11:09.567-0500    INFO    example/basic.go:46     updated user password   {"request_id": "001", "db_upsert": "300ms"}
	// 2019-11-11T11:11:09.567-0500    INFO    cannonical_log_line     {"request_id": "001", "auth_role": "user_rw", "service": "defaultHandler", "db_upsert": "300ms", "request_duration": "86.803µs"}
```

To create a new cannon logger, you can use convenience methods for Zap's development or production loggers, or create your own with any of Zap's options:

| Method | Result |
| ------ | ------ |
| `cannon.NewDevelopment(...zap.Option)` | Wraps Zap's development logger designed for human-readability |
| `cannon.NewProduction(...zap.Option)` | Wraps Zap's production logger with JSON output |

# Passing the logger through context

Cannon provides convenience methods for passing the logger through context to other services or request handlers.  Create a new logger for each request then log as you normally would through each handler and service, then call `cannon.Emit` when the request is finished.

| `CtxLogger(parent context.Context, logger *zap.Logger) context.Context` | Passes the logger in the context |
| `LoggerFromCtx(ctx context.Context) (*zap.Logger, error)` | Get logger from context

For example, a simple middleware might be:

```go
func CannonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		logger, _ := cannon.NewProduction()
		req := r.WithContext(cannon.CtxLogger(r.Context(), logger))

		next.ServeHTTP(w, req)

		cannon.Emit(logger,
			zap.Duration("request_duration", time.Now().Sub(start)),
		)
	})
}
```

# Factory logging construction

Cannon also provides a facility for a global logger factory if you want more control over how the logger is created.  

```go
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
```

# License

MIT


