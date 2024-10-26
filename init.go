package command

import (
	"github.com/go-zoox/command/config"
	"github.com/go-zoox/command/engine"
	"github.com/go-zoox/command/engine/caas"
	"github.com/go-zoox/command/engine/dind"
	"github.com/go-zoox/command/engine/docker"
	"github.com/go-zoox/command/engine/host"

	// "github.com/go-zoox/command/engine/idp"
	"github.com/go-zoox/command/engine/ssh"
)

func init() {
	// This is the init function.
	// It is called when the package is initialized

	// Register the engines

	// Register the host engine
	engine.Register(host.Name, func(cfg *config.Config) (engine.Engine, error) {
		engine, err := host.New(&host.Config{
			ID: cfg.ID,
			//
			Command:     cfg.Command,
			WorkDir:     cfg.WorkDir,
			Environment: cfg.Environment,
			User:        cfg.User,
			Shell:       cfg.Shell,
			//
			ReadOnly: cfg.ReadOnly,
			//
			IsHistoryDisabled: cfg.IsHistoryDisabled,
			//
			IsInheritEnvironmentEnabled: cfg.IsInheritEnvironmentEnabled,
			//
			AllowedSystemEnvKeys: cfg.AllowedSystemEnvKeys,
		})
		if err != nil {
			return nil, err
		}

		return engine, nil
	})

	// Register the docker engine
	engine.Register(docker.Name, func(cfg *config.Config) (engine.Engine, error) {
		engine, err := docker.New(&docker.Config{
			ID: cfg.ID,
			//
			Command:     cfg.Command,
			WorkDir:     cfg.WorkDir,
			Environment: cfg.Environment,
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
			//
			AllowedSystemEnvKeys: cfg.AllowedSystemEnvKeys,
		})
		if err != nil {
			return nil, err
		}

		return engine, nil
	})

	// Register the caas engine
	engine.Register(caas.Name, func(cfg *config.Config) (engine.Engine, error) {
		engine, err := caas.New(&caas.Config{
			ID: cfg.ID,
			//
			Command:     cfg.Command,
			WorkDir:     cfg.WorkDir,
			Environment: cfg.Environment,
			User:        cfg.User,
			Shell:       cfg.Shell,
			//
			ReadOnly: cfg.ReadOnly,
			//
			Server:       cfg.Server,
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.ClientSecret,
			//
			AllowedSystemEnvKeys: cfg.AllowedSystemEnvKeys,
		})
		if err != nil {
			return nil, err
		}

		return engine, nil
	})

	// Register the dind engine
	engine.Register(dind.Name, func(cfg *config.Config) (engine.Engine, error) {
		engine, err := dind.New(&dind.Config{
			ID: cfg.ID,
			//
			Command:     cfg.Command,
			WorkDir:     cfg.WorkDir,
			Environment: cfg.Environment,
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
			//
			AllowedSystemEnvKeys: cfg.AllowedSystemEnvKeys,
		})
		if err != nil {
			return nil, err
		}

		return engine, nil
	})

	// Register the ssh engine
	engine.Register(ssh.Name, func(cfg *config.Config) (engine.Engine, error) {
		engine, err := ssh.New(&ssh.Config{
			ID: cfg.ID,
			//
			Command:     cfg.Command,
			WorkDir:     cfg.WorkDir,
			Environment: cfg.Environment,
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
			//
			AllowedSystemEnvKeys: cfg.AllowedSystemEnvKeys,
		})
		if err != nil {
			return nil, err
		}

		return engine, nil
	})

	// // Register the idp engine
	// engine.Register(idp.Name, func(cfg *config.Config) (engine.Engine, error) {
	// 	engine, err := idp.New(&idp.Config{
	// 		ID: cfg.ID,
	// 		//
	// 		Command:     cfg.Command,
	// 		WorkDir:     cfg.WorkDir,
	// 		Environment: cfg.Environment,
	// 		User:        cfg.User,
	// 		Shell:       cfg.Shell,
	// 		//
	// 		ReadOnly: cfg.ReadOnly,
	// 		//
	// 		Server:       cfg.Server,
	// 		ClientID:     cfg.ClientID,
	// 		ClientSecret: cfg.ClientSecret,
	// 		//
	// 		AllowedSystemEnvKeys: cfg.AllowedSystemEnvKeys,
	// 	})
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	return engine, nil
	// })
}
