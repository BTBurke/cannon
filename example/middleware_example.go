package main

import (
	"net/http"
	"time"

	"github.com/BTBurke/cannon"
	"go.uber.org/zap"
)

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
