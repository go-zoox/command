package engine

import (
	"errors"

	"github.com/go-zoox/command/config"
	"github.com/go-zoox/core-utils/safe"
)

var container = safe.NewMap[string, func(cfg *config.Config) (Engine, error)]()

// ErrEngineNotFound is the error returned when an engine is not found.
var ErrEngineExists = errors.New("engine exists")

// ErrEngineNotFound is the error returned when an engine is not found.
var ErrEngineNotFound = errors.New("engine not found")

// Register registers an engine.
func Register(name string, e func(cfg *config.Config) (Engine, error)) error {
	return container.Set(name, e)
}

// Get gets an engine.
func Get(name string) (func(cfg *config.Config) (Engine, error), error) {
	e := container.Get(name)
	if e == nil {
		return nil, ErrEngineNotFound
	}

	return e, nil
}
