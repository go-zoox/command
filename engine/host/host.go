package host

import (
	"context"
	"io"
	"os"
	"os/exec"

	"github.com/go-zoox/command/engine"
)

// Name is the name of the engine.
const Name = "host"

type host struct {
	ctx context.Context
	//
	cfg *Config
	//
	cmd *exec.Cmd
	//

	//
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
}

// New creates a new host engine.
func New(ctx context.Context, cfg *Config) (engine.Engine, error) {
	if cfg.Shell == "" {
		cfg.Shell = "/bin/sh"
	}

	h := &host{
		ctx: ctx,
		cfg: cfg,
		//
		stdin:  os.Stdin,
		stdout: os.Stdout,
		stderr: os.Stderr,
	}

	if err := h.create(); err != nil {
		return nil, err
	}

	return h, nil
}
