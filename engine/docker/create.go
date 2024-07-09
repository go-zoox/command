package docker

import (
	"context"
	"fmt"
	"os"

	"github.com/docker/cli/cli/streams"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/go-zoox/core-utils/cast"
	"github.com/go-zoox/core-utils/strings"
	"github.com/go-zoox/datetime"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

// create creates a container.
func (d *docker) create() (err error) {
	if d.cfg.Command != "" {
		d.args = append(d.args, "-c", d.cfg.Command)
	}

	if len(d.cfg.AllowedSystemEnvKeys) != 0 {
		for _, key := range d.cfg.AllowedSystemEnvKeys {
			if value, ok := os.LookupEnv(key); ok {
				d.env = append(d.env, fmt.Sprintf("%s=%s", key, value))
			}
		}
	}

	for k, v := range d.cfg.Environment {
		d.env = append(d.env, fmt.Sprintf("%s=%s", k, v))
	}

	d.stderr.Write([]byte(fmt.Sprintf("[%s][docker] start to prepare docker environment ...\n", datetime.Now().Format())))

	d.client, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation(), func(c *client.Client) error {
		if d.cfg.DockerHost != "" {
			if err := client.WithHost(d.cfg.DockerHost)(c); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		d.stderr.Write([]byte(fmt.Sprintf("[%s][docker] failed to connect docker engine .\n", datetime.Now().Format())))
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
		// NetworkMode: "none",
		Privileged: d.cfg.Privileged,
	}
	if d.cfg.Memory != 0 {
		hostCfg.Resources.Memory = d.cfg.Memory * 1024 * 1024
	}
	if d.cfg.DisableNetwork {
		hostCfg.NetworkMode = "none"
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

	networkCfg := &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{},
	}
	if d.cfg.Network != "" {
		d.stderr.Write([]byte(fmt.Sprintf("[%s][docker] inspect network %s ...\n", datetime.Now().Format(), d.cfg.Network)))
		networkIns, err := d.client.NetworkInspect(context.Background(), d.cfg.Network, network.InspectOptions{})
		if err != nil {
			d.stderr.Write([]byte(fmt.Sprintf("[%s][docker] failed to inspect network: %s\n", datetime.Now().Format(), err)))
			return err
		}

		networkCfg.EndpointsConfig[d.cfg.Network] = &network.EndpointSettings{
			NetworkID: networkIns.ID,
		}
	}

	platformCfg := &ocispec.Platform{
		// OS:           "linux",
		// Architecture: "amd64",
	}
	if d.cfg.Platform != "" {
		d.stderr.Write([]byte(fmt.Sprintf("[%s][docker] platform: %s ...\n", datetime.Now().Format(), d.cfg.Platform)))
		switch d.cfg.Platform {
		case "linux/amd64", "linux/arm64":
		default:
			return fmt.Errorf("invalid platform: %s, available: linux/amd64, linux/arm64", d.cfg.Platform)
		}

		osArch := strings.Split(d.cfg.Platform, "/")
		platformCfg.OS = osArch[0]
		platformCfg.Architecture = osArch[1]
	}

	_, _, err = d.client.ImageInspectWithRaw(context.Background(), d.cfg.Image)
	if err != nil {
		d.stderr.Write([]byte(fmt.Sprintf("[%s][docker] pull image %s ...\n", datetime.Now().Format(), d.cfg.Image)))
		imagePullReader, err := d.client.ImagePull(context.Background(), d.cfg.Image, image.PullOptions{
			Platform: d.cfg.Platform,
		})
		if err != nil {
			d.stderr.Write([]byte(fmt.Sprintf("[%s][docker] failed to pull image %s ...\n", datetime.Now().Format(), err)))
			return err
		}
		defer imagePullReader.Close()

		if err := jsonmessage.DisplayJSONMessagesToStream(imagePullReader, streams.NewOut(d.stderr), nil); err != nil {
			return err
		}
	}

	d.container, err = d.client.ContainerCreate(context.Background(), cfg, hostCfg, networkCfg, platformCfg, d.cfg.ID)
	if err != nil {
		d.stderr.Write([]byte(fmt.Sprintf("[%s][docker] failed to create container: %s\n", datetime.Now().Format(), err)))
		return err
	}

	d.stderr.Write([]byte(fmt.Sprintf("[%s][docker] succeed to prepare docker environment.\n", datetime.Now().Format())))

	return nil
}
