package internal

import (
	"go.uber.org/zap/zapcore"
	"sync"
)

var _ zapcore.Core = &CannonicalLog{}

type CannonicalLog struct {
	WrappedCore zapcore.Core
	EmptyCore   zapcore.Core
	Fields      []zapcore.Field

	m sync.Mutex
}

func (c *CannonicalLog) Enabled(l zapcore.Level) bool {
	return c.WrappedCore.Enabled(l)
}

func (c *CannonicalLog) With(f []zapcore.Field) zapcore.Core {
	return &CannonicalLog{
		EmptyCore:   c,
		WrappedCore: c.WrappedCore.With(f),
		Fields:      append(c.Fields, f...),
	}
}

func (c *CannonicalLog) Check(e zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(e.Level) {
		return ce.AddCore(e, c)
	}
	return ce
}

func (c *CannonicalLog) Write(e zapcore.Entry, f []zapcore.Field) error {
	c.m.Lock()
	defer c.m.Unlock()
	c.Fields = append(c.Fields, f...)
	return c.WrappedCore.Write(e, f)
}

func (c *CannonicalLog) Sync() error {
	return c.WrappedCore.Sync()
}

// Reset will clear all set fields to allow the core to be reused
func (c *CannonicalLog) Reset() {
	c.Fields = []zapcore.Field{}
	return
}
