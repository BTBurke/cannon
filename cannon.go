package cannon

import (
	"errors"
	"time"

	"github.com/BTBurke/cannon/internal"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewCore creates a new zapcore.Core that exends zap's logging core to enable cannonical logging when the Emit method
// is called.  This is a low level primitive that allows you to pass additional options to the zap logging contructor to
// get exactly the functionality you want.  For a higher level log constructor, you can use `cannon.NewProduction()`,
// `cannon.NewDevelopment()`, or pass your own log factory to `RegisterFactory` and then call `cannon.NewLogger()`
func NewCore() zap.Option {
	return zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		core := &internal.CannonicalLog{
			EmptyCore:   c,
			WrappedCore: c,
		}
		return core
	})
}

// Emit is called at the end of the request to emit the cannonical log line that contains all fields set throughout
// the lifetime of the call
func Emit(log *zap.Logger, fields ...zap.Field) error {
	c, ok := log.Core().(*internal.CannonicalLog)
	if !ok {
		return errors.New("unknown logger type")
	}
	if err := c.EmptyCore.Write(zapcore.Entry{
		Time:    time.Now(),
		Message: "cannonical_log_line",
	}, append(c.Fields, fields...)); err != nil {
		return err
	}
	if err := c.EmptyCore.Sync(); err != nil {
		return err
	}
	return nil
}

// NewDevelopment gives you a `zap.NewDevelopment` configuration with the ability to emit a cannonical logline
// at the end of the request
func NewDevelopment(options ...zap.Option) (*zap.Logger, error) {
	return zap.NewDevelopment(append(options, NewCore())...)
}

// NewProduction gives you a `zap.NewProduction` configuration with the ability to emit a cannonical logline
// at the end of the request
func NewProduction(options ...zap.Option) (*zap.Logger, error) {
	return zap.NewProduction(append(options, NewCore())...)
}
