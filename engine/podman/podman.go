package podman

import (
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/go-zoox/command/engine"
	"github.com/go-zoox/uuid"
)

// Name is the name of the engine.
const Name = "podman"

type podman struct {
	cfg *Config
	//
	args []string
	env  []string
	//
	client *client.Client
	//
	container container.CreateResponse
	//
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
}

// New creates a new podman engine.
func New(cfg *Config) (engine.Engine, error) {
	if cfg.Image == "" {
		cfg.Image = "docker.io/library/alpine:latest"
	}
	if cfg.Shell == "" {
		cfg.Shell = "/bin/sh"
	}
	if cfg.ID == "" {
		cfg.ID = fmt.Sprintf("go-zoox_command_%s", uuid.V4())
	}

	p := &podman{
		cfg:    cfg,
		stdin:  os.Stdin,
		stdout: os.Stdout,
		stderr: os.Stderr,
	}

	if err := p.create(); err != nil {
		return nil, err
	}

	return p, nil
}
