package docker

import (
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/go-zoox/uuid"
)

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

	hostCfg := &container.HostConfig{}

	d.container, err = d.client.ContainerCreate(d.ctx, cfg, hostCfg, nil, nil, uuid.V4())
	if err != nil {
		return err
	}

	return nil
}
