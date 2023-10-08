package docker

import (
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/go-zoox/core-utils/cast"
	"github.com/go-zoox/core-utils/strings"
	"github.com/go-zoox/uuid"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

// create creates a container.
func (d *docker) create() (err error) {
	if d.cfg.Command != "" {
		d.args = append(d.args, "-c", d.cfg.Command)
	}

	for k, v := range d.cfg.Environment {
		d.env = append(d.env, fmt.Sprintf("%s=%s", k, v))
	}

	d.client, err = client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}

	cfg := &container.Config{
		Hostname:     "go-zoox",
		Image:        d.cfg.Image,
		Cmd:          append([]string{d.cfg.Shell}, d.args...),
		User:         d.cfg.User,
		WorkingDir:   d.cfg.WorkDir,
		Env:          d.env,
		Tty:          true,
		OpenStdin:    true,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		StdinOnce:    true,
	}

	hostCfg := &container.HostConfig{
		Resources: container.Resources{
			// Memory:    d.cfg.Memory,
			// CPUPeriod: 100000,
			// CPUQuota:  cast.ToInt64(100000 * d.cfg.CPU),
		},
		Mounts: []mount.Mount{
			// {
			// 	Type:     mount.TypeBind,
			// 	Source:   d.cfg.WorkDir,
			// 	Target:   d.cfg.WorkDir,
			// 	ReadOnly: false,
			// },
		},
	}
	if d.cfg.Memory != 0 {
		hostCfg.Resources.Memory = d.cfg.Memory * 1024 * 1024
	}

	if d.cfg.CPU != 0 {
		hostCfg.Resources.CPUPeriod = 100000
		hostCfg.Resources.CPUQuota = cast.ToInt64(float64(hostCfg.Resources.CPUPeriod) * d.cfg.CPU)
	}
	if d.cfg.WorkDir != "" {
		hostCfg.Mounts = append(hostCfg.Mounts, mount.Mount{
			Type:     mount.TypeBind,
			Source:   d.cfg.WorkDir,
			Target:   d.cfg.WorkDir,
			ReadOnly: false,
		})
	}

	platformCfg := &ocispec.Platform{
		OS:           "linux",
		Architecture: "amd64",
	}
	if d.cfg.Platform != "" {
		switch d.cfg.Platform {
		case "linux/amd64", "linux/arm64":
		default:
			return fmt.Errorf("invalid platform: %s, available: linux/amd64, linux/arm64", d.cfg.Platform)
		}

		osArch := strings.Split(d.cfg.Platform, "/")
		platformCfg.OS = osArch[0]
		platformCfg.Architecture = osArch[1]
	}

	d.container, err = d.client.ContainerCreate(d.ctx, cfg, hostCfg, nil, platformCfg, uuid.V4())
	if err != nil {
		return err
	}

	return nil
}
