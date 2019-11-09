package cannon

import (
	"errors"
	"time"

	"github.com/BTBurke/cannon/internal"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger() zap.Option {
	return zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		core := &internal.CannonicalLog{
			EmptyCore:   c,
			WrappedCore: c,
		}
		return core
	})
}

func Emit(log *zap.Logger) error {
	c, ok := log.Core().(*internal.CannonicalLog)
	if !ok {
		return errors.New("unknown logger type")
	}
	if err := c.EmptyCore.Write(zapcore.Entry{
		Time:    time.Now(),
		Message: "cannonical_log_line",
	}, c.Fields); err != nil {
		return err
	}
	if err := c.EmptyCore.Sync(); err != nil {
		return err
	}
	return nil
}
