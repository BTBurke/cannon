package cannon

import (
	"context"
	"errors"

	"go.uber.org/zap"
)

type ctxLogger string

// CtxLogger will pass the logger in the context to subsequent request handlers
func CtxLogger(parent context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(parent, ctxLogger("cannon"), logger)
}

// LoggerFromCtx will extract a logger from the context
func LoggerFromCtx(ctx context.Context) (*zap.Logger, error) {
	log := ctx.Value(ctxLogger("cannon"))
	logger, ok := log.(*zap.Logger)
	if !ok || logger == nil {
		return nil, errors.New("unknown logger type")
	}
	return logger, nil
}
