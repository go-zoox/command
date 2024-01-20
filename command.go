package command

import (
	"context"
	"fmt"
	"io"

	"github.com/go-zoox/command/agent/client"
	"github.com/go-zoox/command/config"
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

// Config is the command runner config
type Config = config.Config

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
			ReadOnly: cfg.ReadOnly,
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
			Command:     cfg.Command,
			WorkDir:     cfg.WorkDir,
			Environment: environment,
			User:        cfg.User,
			Shell:       cfg.Shell,
			//
			ReadOnly: cfg.ReadOnly,
			//
			Image:          cfg.Image,
			Memory:         cfg.Memory,
			CPU:            cfg.CPU,
			Platform:       cfg.Platform,
			Network:        cfg.Network,
			DisableNetwork: cfg.DisableNetwork,
			Privileged:     cfg.Privileged,
			//
			DockerHost: cfg.DockerHost,
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
			ReadOnly: cfg.ReadOnly,
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
			Command:     cfg.Command,
			WorkDir:     cfg.WorkDir,
			Environment: environment,
			User:        cfg.User,
			Shell:       cfg.Shell,
			//
			ReadOnly: cfg.ReadOnly,
			//
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
			ReadOnly: cfg.ReadOnly,
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
