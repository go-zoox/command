package podman

import (
	"context"
	"fmt"
	"os"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/go-zoox/core-utils/cast"
)

// defaultPodmanHost is the default Podman socket (Docker-compatible API).
const defaultPodmanHost = "unix:///run/podman/podman.sock"

// create creates a container via Podman's Docker-compatible API.
func (p *podman) create() (err error) {
	if p.cfg.Command != "" {
		p.args = append(p.args, "-c", p.cfg.Command)
	}

	if len(p.cfg.AllowedSystemEnvKeys) != 0 {
		for _, key := range p.cfg.AllowedSystemEnvKeys {
			if value, ok := os.LookupEnv(key); ok {
				p.env = append(p.env, fmt.Sprintf("%s=%s", key, value))
			}
		}
	}

	for k, v := range p.cfg.Environment {
		p.env = append(p.env, fmt.Sprintf("%s=%s", k, v))
	}

	host := p.cfg.PodmanHost
	if host == "" {
		host = defaultPodmanHost
	}
	if v := os.Getenv("PODMAN_HOST"); v != "" && p.cfg.PodmanHost == "" {
		host = v
	}

	p.client, err = client.NewClientWithOpts(client.WithHost(host), client.WithAPIVersionNegotiation())
	if err != nil {
		return fmt.Errorf("podman: connect: %w", err)
	}

	cfg := &container.Config{
		Hostname:     "go-zoox",
		Image:        p.cfg.Image,
		Cmd:          append([]string{p.cfg.Shell}, p.args...),
		User:         p.cfg.User,
		WorkingDir:   p.cfg.WorkDir,
		Env:          p.env,
		Tty:          true,
		OpenStdin:    true,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		StdinOnce:    true,
	}

	hostCfg := &container.HostConfig{
		AutoRemove: true,
		Resources:  container.Resources{},
		Privileged: p.cfg.Privileged,
	}

	if p.cfg.Memory != 0 {
		hostCfg.Resources.Memory = p.cfg.Memory * 1024 * 1024
	}
	if p.cfg.CPU != 0 {
		hostCfg.Resources.CPUPeriod = 100000
		hostCfg.Resources.CPUQuota = cast.ToInt64(float64(hostCfg.Resources.CPUPeriod) * p.cfg.CPU)
	}
	if p.cfg.DisableNetwork {
		hostCfg.NetworkMode = "none"
	}

	p.container, err = p.client.ContainerCreate(context.Background(), cfg, hostCfg, nil, nil, p.cfg.ID)
	if err != nil {
		return fmt.Errorf("podman: create container: %w", err)
	}

	return nil
}
