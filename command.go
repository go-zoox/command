package command

import (
	"context"
	"fmt"

	"github.com/go-zoox/command/engine"
	"github.com/go-zoox/command/engine/docker"
	"github.com/go-zoox/command/engine/host"
	"github.com/go-zoox/command/terminal"
)

type Command interface {
	Start() error
	Wait() error
	Cancel() error
	//
	Run() error
	//
	Terminal() (terminal.Terminal, error)
}

type Config struct {
	Engine      string
	Command     string
	WorkDir     string
	Environment map[string]string
	User        string
	Shell       string

	// engine = docker
	Image string
}

func New(ctx context.Context, cfg *Config) (cmd Command, err error) {
	if cfg.Engine == "" {
		cfg.Engine = host.Engine
	}

	if cfg.Shell == "" {
		cfg.Shell = "/bin/sh"
	}

	var engine engine.Engine
	switch cfg.Engine {
	case host.Engine:
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
	case docker.Engine:
		engine, err = docker.New(ctx, &docker.Config{
			Command:     cfg.Command,
			WorkDir:     cfg.WorkDir,
			Environment: cfg.Environment,
			User:        cfg.User,
			Shell:       cfg.Shell,
			Image:       cfg.Image,
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
