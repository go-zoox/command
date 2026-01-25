package command

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/go-zoox/command/agent/client"
	"github.com/go-zoox/command/config"
	"github.com/go-zoox/command/engine"
	"github.com/go-zoox/command/engine/host"
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
	Output() ([]byte, error)
	//
	SetStdin(stdin io.Reader) error
	SetStdout(stdout io.Writer) error
	SetStderr(stderr io.Writer) error
	//
	Terminal() (terminal.Terminal, error)
}

// Config is the command runner config
type Config = config.Config

// New creates a new command runner.
func New(cfg *Config) (cmd Command, err error) {
	if cfg.Context == nil {
		cfg.Context = context.Background()
	}

	// If sandbox mode is enabled, force docker engine and apply security settings
	if cfg.Sandbox {
		if cfg.Engine != "" && cfg.Engine != "docker" {
			return nil, fmt.Errorf("sandbox mode requires docker engine, but got: %s", cfg.Engine)
		}
		cfg.Engine = "docker"

		// Apply default sandbox security settings if not explicitly set
		if !cfg.DisableNetwork && cfg.Network == "" {
			// Default to no network in sandbox mode for security
			cfg.DisableNetwork = true
		}
		if cfg.Privileged {
			// Force non-privileged in sandbox mode
			cfg.Privileged = false
		}
		// Set default resource limits if not set
		if cfg.Memory == 0 {
			cfg.Memory = 512 // Default 512MB memory limit
		}
		if cfg.CPU == 0 {
			cfg.CPU = 1.0 // Default 1 CPU core limit
		}
	}

	// Set default engine if not set and sandbox is not enabled
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
		//
		"HOME": os.Getenv("HOME"),
		"USER": os.Getenv("USER"),
		"PATH": os.Getenv("PATH"),
	}
	for k, v := range cfg.Environment {
		environment[k] = v
	}
	cfg.Environment = environment

	// support agent
	if cfg.Agent != "" {
		agent, err := client.New(func(opt *client.Option) {
			opt.Server = cfg.Agent
		})
		if err != nil {
			return nil, err
		}

		if err := agent.Connect(); err != nil {
			return nil, err
		}

		err = agent.New(&config.Config{
			// Context:                          cfg.Context,
			Timeout:                          cfg.Timeout,
			Engine:                           cfg.Engine,
			Sandbox:                          cfg.Sandbox,
			Command:                          cfg.Command,
			WorkDir:                          cfg.WorkDir,
			Environment:                      environment,
			User:                             cfg.User,
			Shell:                            cfg.Shell,
			ReadOnly:                         cfg.ReadOnly,
			IsHistoryDisabled:                cfg.IsHistoryDisabled,
			Image:                            cfg.Image,
			Memory:                           cfg.Memory,
			CPU:                              cfg.CPU,
			Platform:                         cfg.Platform,
			Network:                          cfg.Network,
			DisableNetwork:                   cfg.DisableNetwork,
			Privileged:                       cfg.Privileged,
			DockerHost:                       cfg.DockerHost,
			Server:                           cfg.Server,
			ClientID:                         cfg.ClientID,
			ClientSecret:                     cfg.ClientSecret,
			SSHHost:                          cfg.SSHHost,
			SSHPort:                          cfg.SSHPort,
			SSHUser:                          cfg.SSHUser,
			SSHPass:                          cfg.SSHPass,
			SSHPrivateKey:                    cfg.SSHPrivateKey,
			SSHPrivateKeySecret:              cfg.SSHPrivateKeySecret,
			SSHIsIgnoreStrictHostKeyChecking: cfg.SSHIsIgnoreStrictHostKeyChecking,
			SSHKnowHostsFilePath:             cfg.SSHKnowHostsFilePath,
			ID:                               cfg.ID,
		})
		if err != nil {
			return nil, err
		}

		return agent, nil
	}

	var eg engine.Engine
	if createEngine, err := engine.Get(cfg.Engine); err != nil {
		return nil, fmt.Errorf("unsupported command engine: %s", cfg.Engine)
	} else {
		eg, err = createEngine(cfg)
		if err != nil {
			return nil, err
		}
	}

	go func() {
		<-cfg.Context.Done()
		eg.Cancel()
	}()

	return &command{
		cfg:    cfg,
		engine: eg,
	}, nil
}

type command struct {
	cfg *Config
	//
	engine engine.Engine
}
