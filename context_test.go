package cannon

import (
	"context"
	"testing"

	"go.uber.org/zap"
)

func TestContext(t *testing.T) {
	newLogger := func() *zap.Logger {
		logger, err := zap.NewProduction(Core())
		if err != nil {
			t.Fatalf("failed to create logger: %s", err)
		}
		return logger
	}
	nilLogger := func() *zap.Logger { return nil }

	tt := []struct {
		Name      string
		Logger    func() *zap.Logger
		ShouldErr bool
	}{
		{"embed in context", newLogger, false},
		{"no logger defined", nilLogger, true},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			factory = nil
			logger := tc.Logger()
			ctx := CtxLogger(context.Background(), logger)
			receivedLogger, err := LoggerFromCtx(ctx)
			switch {
			case tc.ShouldErr && err == nil:
				t.Fatalf("logger should have errored but did not, instead got logger: %v", receivedLogger)
				return
			case receivedLogger != logger:
				t.Fatalf("expected logger %v but got %v", logger, receivedLogger)
				return
			default:
			}
		})
	}
}
