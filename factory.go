package cannon

import (
	"errors"

	"go.uber.org/zap"
)

var factory LogFactory

// LogFactoryFunc turns a function into a LogFactory
type LogFactoryFunc func() (*zap.Logger, error)

// LogFactory describes how to create a new logger
type LogFactory interface {
	New() (*zap.Logger, error)
}

// New will return a new logger from the log factory function
func (f LogFactoryFunc) New() (*zap.Logger, error) {
	return f()
}

// NewLogger will use the registered log factory to create a new logger
func NewLogger() (*zap.Logger, error) {
	switch {
	case factory == nil:
		return nil, errors.New("no factory function defined, cannot create logger")
	default:
		return factory.New()
	}
}

// RegisterFactory will register a global log factory that will be used for logger
// creation when `NewLogger()` is called.  This is useful when you want customize how
// the logger is created, such as creating different loggers for development and prod.
func RegisterFactory(f LogFactory) {
	factory = f
}

// ClearFactory will remove any registered log factory function
func ClearFactory() {
	factory = nil
}
