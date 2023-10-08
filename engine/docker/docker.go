package docker

import (
	"context"
	"io"
	"os"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/go-zoox/command/engine"
)

// Name is the name of the engine.
const Name = "docker"

type docker struct {
	cfg *Config
	//
	args []string
	env  []string
	//
	ctx context.Context
	//
	client *client.Client
	//
	container container.CreateResponse

	//
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
}

// New creates a new docker engine.
func New(ctx context.Context, cfg *Config) (engine.Engine, error) {
	if cfg.Image == "" {
		cfg.Image = "whatwewant/zmicro:v1"
	}

	if cfg.Shell == "" {
		cfg.Shell = "/bin/sh"
	}

	d := &docker{
		ctx: ctx,
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
