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

// CtxSLogger will pass a sugared logger in the context
func CtxSLogger(parent context.Context, logger *zap.SugaredLogger) context.Context {
	return CtxLogger(parent, logger.Desugar())
}

// SSLoggerFromCtx will extract a sugared logger from the context
func SLoggerFromCtx(ctx context.Context) (*zap.SugaredLogger, error) {
	l, err := LoggerFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	return l.Sugar(), nil
}
