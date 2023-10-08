package command

import (
	"context"
	"fmt"
	"io"

	"github.com/go-zoox/command/engine"
	"github.com/go-zoox/command/engine/docker"
	"github.com/go-zoox/command/engine/host"
	"github.com/go-zoox/command/terminal"
)

// Command is the command runner interface
type Command interface {
	Start() error
	Wait() error
	Cancel() error
	//
	Run() error
	//
	SetStdin(stdin io.Reader) error
	SetStdout(stdout io.Writer) error
	SetStderr(stderr io.Writer) error
	//
	Terminal() (terminal.Terminal, error)
}

// Config is the command config
type Config struct {
	Engine      string
	Command     string
	WorkDir     string
	Environment map[string]string
	User        string
	Shell       string

	// engine = docker
	Image string
	// Memory is the memory limit, unit: MB
	Memory int64
	// CPU is the CPU limit, unit: core
	CPU float64
	// Platform is the command platform, available: linux/amd64, linux/arm64
	Platform string
	// Network is the network name
	Network string
	// DisableNetwork disables network
	DisableNetwork bool
}

// New creates a new command runner.
func New(ctx context.Context, cfg *Config) (cmd Command, err error) {
	if cfg.Engine == "" {
		cfg.Engine = host.Name
	}

	if cfg.Shell == "" {
		cfg.Shell = "/bin/sh"
	}

	var engine engine.Engine
	switch cfg.Engine {
	case host.Name:
		engine, err = host.New(ctx, &host.Config{
			Command:     cfg.Command,
			WorkDir:     cfg.WorkDir,
			Environment: cfg.Environment,
			User:        cfg.User,
			Shell:       cfg.Shell,
		})
		if err != nil {
			return nil, err
		}
	case docker.Name:
		engine, err = docker.New(ctx, &docker.Config{
			Command:        cfg.Command,
			WorkDir:        cfg.WorkDir,
			Environment:    cfg.Environment,
			User:           cfg.User,
			Shell:          cfg.Shell,
			Image:          cfg.Image,
			Memory:         cfg.Memory,
			CPU:            cfg.CPU,
			Platform:       cfg.Platform,
			Network:        cfg.Network,
			DisableNetwork: cfg.DisableNetwork,
		})
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsupported command engine: %s", cfg.Engine)
	}

	return &command{
		engine: engine,
	}, nil
}

type command struct {
	engine engine.Engine
}
