package cannon

import (
	"context"
	"errors"

	"go.uber.org/zap"
)

type ctxLogger string

func ContextLogger(parent context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(parent, ctxLogger("cannon"), logger)
}

func LoggerFromCtx(ctx context.Context) (*zap.Logger, error) {
	log := ctx.Value(ctxLogger("cannon"))
	if log == nil {
		return nil, errors.New("no logger found")
	}
	logger, ok := log.(*zap.Logger)
	if !ok {
		return nil, errors.New("unknown logger type")
	}
	return logger, nil
}
