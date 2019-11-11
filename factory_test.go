package cannon

import (
	"testing"

	"go.uber.org/zap"
)

func TestFactory(t *testing.T) {
	factoryFunc := LogFactoryFunc(func() (*zap.Logger, error) { return NewProduction() })
	_, err := factoryFunc.New()
	if err != nil {
		t.Fatalf("no error expected: %v", err)
	}
	RegisterFactory(factoryFunc)
	if factory == nil {
		t.Fatalf("expected factory to be the factory function")
	}
	log, err := NewFactoryLogger()
	if log == nil || err != nil {
		t.Fatalf("expected factory to create a new logger, but got %v", err)
	}
	ClearFactory()
	_, err = NewFactoryLogger()
	if err == nil {
		t.Fatalf("expected error when no factory defined")
	}
}
