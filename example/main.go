package main

import (
	"time"

	"github.com/BTBurke/cannon"
	"go.uber.org/zap"
)

func main() {
	log, _ := zap.NewDevelopment(cannon.NewLogger())

	logger := log.With(
		zap.String("request_id", "001"),
	)
	logger.Info("auth success", zap.String("auth_role", "user_rw"))
	logger.Info("updated user password", zap.Duration("db_upsert", 300*time.Millisecond))
	cannon.Emit(logger)
}
