package command

import (
	"context"
	"fmt"
	"io"

	"github.com/go-zoox/command/engine"
	"github.com/go-zoox/command/engine/caas"
	"github.com/go-zoox/command/engine/dind"
	"github.com/go-zoox/command/engine/docker"
	"github.com/go-zoox/command/engine/host"
	"github.com/go-zoox/command/engine/ssh"
	"github.com/go-zoox/command/terminal"
	"github.com/go-zoox/uuid"
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
	Context context.Context

	// Engine is the command engine, available: host, docker
	Engine string

	// engine common
	Command     string
	WorkDir     string
	Environment map[string]string
	User        string
	Shell       string

	// engine = host
	IsHistoryDisabled bool

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
	// Privileged enables privileged mode
	Privileged bool

	// engine = caas
	// Server is the command server address
	Server string
	// ClientID is the client ID for server auth
	ClientID string
	// ClientSecret is the client secret for server auth
	ClientSecret string

	// engine = ssh
	SSHHost                          string
	SSHPort                          int
	SSHUser                          string
	SSHPass                          string
	SSHPrivateKey                    string
	SSHPrivateKeySecret              string
	SSHIsIgnoreStrictHostKeyChecking bool
	SSHKnowHostsFilePath             string

	// Custom Command Runner ID
	ID string
}

// New creates a new command runner.
func New(cfg *Config) (cmd Command, err error) {
	if cfg.Context == nil {
		cfg.Context = context.Background()
	}

	if cfg.Engine == "" {
		cfg.Engine = host.Name
	}

	if cfg.Shell == "" {
		cfg.Shell = "/bin/sh"
	}

	if cfg.ID == "" {
		cfg.ID = fmt.Sprintf("go-zoox_command_%s", uuid.V4())
	}

	environment := map[string]string{
		"GO_ZOOX_COMMAND_ENGINE":          cfg.Engine,
		"GO_ZOOX_COMMAND_ID":              cfg.ID,
		"GO_ZOOX_COMMAND_SHELL":           cfg.Shell,
		"GO_ZOOX_COMMAND_USER":            cfg.User,
		"GO_ZOOX_COMMAND_WORKDIR":         cfg.WorkDir,
		"GO_ZOOX_COMMAND_COMMAND":         cfg.Command,
		"GO_ZOOX_COMMAND_IMAGE":           cfg.Image,
		"GO_ZOOX_COMMAND_MEMORY":          fmt.Sprintf("%d", cfg.Memory),
		"GO_ZOOX_COMMAND_CPU":             fmt.Sprintf("%f", cfg.CPU),
		"GO_ZOOX_COMMAND_PLATFORM":        cfg.Platform,
		"GO_ZOOX_COMMAND_NETWORK":         cfg.Network,
		"GO_ZOOX_COMMAND_DISABLE_NETWORK": fmt.Sprintf("%t", cfg.DisableNetwork),
	}
	for k, v := range cfg.Environment {
		environment[k] = v
	}

	var engine engine.Engine
	switch cfg.Engine {
	case host.Name:
		engine, err = host.New(&host.Config{
			ID: cfg.ID,
			//
			Command:     cfg.Command,
			WorkDir:     cfg.WorkDir,
			Environment: environment,
			User:        cfg.User,
			Shell:       cfg.Shell,
			//
			IsHistoryDisabled: cfg.IsHistoryDisabled,
		})
		if err != nil {
			return nil, err
		}
	case docker.Name:
		engine, err = docker.New(&docker.Config{
			ID: cfg.ID,
			//
			Command:        cfg.Command,
			WorkDir:        cfg.WorkDir,
			Environment:    environment,
			User:           cfg.User,
			Shell:          cfg.Shell,
			Image:          cfg.Image,
			Memory:         cfg.Memory,
			CPU:            cfg.CPU,
			Platform:       cfg.Platform,
			Network:        cfg.Network,
			DisableNetwork: cfg.DisableNetwork,
			Privileged:     cfg.Privileged,
		})
		if err != nil {
			return nil, err
		}
	case caas.Name:
		engine, err = caas.New(&caas.Config{
			ID: cfg.ID,
			//
			Command:     cfg.Command,
			WorkDir:     cfg.WorkDir,
			Environment: environment,
			User:        cfg.User,
			Shell:       cfg.Shell,
			//
			Server:       cfg.Server,
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.ClientSecret,
		})
		if err != nil {
			return nil, err
		}
	case dind.Name:
		engine, err = dind.New(&dind.Config{
			ID: cfg.ID,
			//
			Command:        cfg.Command,
			WorkDir:        cfg.WorkDir,
			Environment:    environment,
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
	case ssh.Name:
		engine, err = ssh.New(&ssh.Config{
			ID: cfg.ID,
			//
			Command:     cfg.Command,
			WorkDir:     cfg.WorkDir,
			Environment: environment,
			// User:        cfg.User,
			Shell: cfg.Shell,
			//
			Host:             cfg.SSHHost,
			Port:             cfg.SSHPort,
			User:             cfg.SSHUser,
			Pass:             cfg.SSHPass,
			PrivateKey:       cfg.SSHPrivateKey,
			PrivateKeySecret: cfg.SSHPrivateKeySecret,
			//
			IsIgnoreStrictHostKeyChecking: cfg.SSHIsIgnoreStrictHostKeyChecking,
			KnowHostsFilePath:             cfg.SSHKnowHostsFilePath,
		})
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsupported command engine: %s", cfg.Engine)
	}

	go func() {
		<-cfg.Context.Done()
		engine.Cancel()
	}()

	return &command{
		cfg:    cfg,
		engine: engine,
	}, nil
}

type command struct {
	cfg *Config
	//
	engine engine.Engine
}
