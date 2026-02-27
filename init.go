package command

import (
	"github.com/go-zoox/command/config"
	"github.com/go-zoox/command/engine"
	"github.com/go-zoox/command/engine/caas"
	"github.com/go-zoox/command/engine/dind"
	"github.com/go-zoox/command/engine/docker"
	"github.com/go-zoox/command/engine/host"
	"github.com/go-zoox/command/engine/k8s"
	"github.com/go-zoox/command/engine/podman"

	"github.com/go-zoox/command/engine/ssh"
	"github.com/go-zoox/command/engine/wsl"
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
			ImageRegistry:         cfg.ImageRegistry,
			ImageRegistryUsername: cfg.ImageRegistryUsername,
			ImageRegistryPassword: cfg.ImageRegistryPassword,
			Runtime:               cfg.DockerRuntime,
			//
			AllowedSystemEnvKeys: cfg.AllowedSystemEnvKeys,
			//
			DataDirOuter: cfg.DataDirOuter,
			DataDirInner: cfg.DataDirInner,
			//
			Sandbox: cfg.Sandbox,
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

	// Register the k8s engine
	engine.Register(k8s.Name, func(cfg *config.Config) (engine.Engine, error) {
		k8sImage := cfg.K8sImage
		if k8sImage == "" {
			k8sImage = cfg.Image
		}
		engine, err := k8s.New(&k8s.Config{
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
			Kubeconfig:         cfg.K8sKubeconfig,
			Namespace:         cfg.K8sNamespace,
			Image:             k8sImage,
			JobTimeoutSeconds: cfg.K8sPodTimeoutSeconds,
			//
			AllowedSystemEnvKeys: cfg.AllowedSystemEnvKeys,
		})
		if err != nil {
			return nil, err
		}

		return engine, nil
	})

	// Register the podman engine
	engine.Register(podman.Name, func(cfg *config.Config) (engine.Engine, error) {
		podmanImage := cfg.Image
		if podmanImage == "" {
			podmanImage = "docker.io/library/alpine:latest"
		}
		engine, err := podman.New(&podman.Config{
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
			Image:          podmanImage,
			Memory:         cfg.Memory,
			CPU:            cfg.CPU,
			Platform:       cfg.Platform,
			Network:        cfg.Network,
			DisableNetwork: cfg.DisableNetwork,
			Privileged:     cfg.Privileged,
			//
			PodmanHost: cfg.PodmanHost,
			//
			AllowedSystemEnvKeys: cfg.AllowedSystemEnvKeys,
		})
		if err != nil {
			return nil, err
		}

		return engine, nil
	})

	// Register the wsl engine (Windows only)
	engine.Register(wsl.Name, func(cfg *config.Config) (engine.Engine, error) {
		engine, err := wsl.New(&wsl.Config{
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
			WSLDistro: cfg.WSLDistro,
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
}
