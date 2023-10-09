package dind

import (
	"io"
	"os"

	"github.com/go-zoox/command/engine"
)

// Name is the name of the engine.
const Name = "dind"

type dind struct {
	cfg *Config
	//
	client engine.Engine

	//
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
}

// New creates a new dind engine.
func New(cfg *Config) (engine.Engine, error) {
	cfg.Image = "whatwewant/dind:v24-1"

	d := &dind{
		cfg: cfg,
		//
		stdin:  os.Stdin,
		stdout: os.Stdout,
		stderr: os.Stderr,
	}

	if err := d.create(); err != nil {
		return nil, err
	}

	return d, nil
}
