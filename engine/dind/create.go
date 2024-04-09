package dind

import (
	"os"

	"github.com/go-zoox/command/engine/docker"
)

// create creates a container.
func (d *dind) create() (err error) {
	if len(d.cfg.AllowedSystemEnvKeys) != 0 {
		for _, key := range d.cfg.AllowedSystemEnvKeys {
			if d.cfg.Environment[key] == "" {
				if value, ok := os.LookupEnv(key); ok {
					d.cfg.Environment[key] = value
				}
			}
		}
	}

	d.client, err = docker.New(&docker.Config{
		ID: d.cfg.ID,
		//
		Command:        d.cfg.Command,
		WorkDir:        d.cfg.WorkDir,
		Environment:    d.cfg.Environment,
		User:           d.cfg.User,
		Shell:          d.cfg.Shell,
		ReadOnly:       d.cfg.ReadOnly,
		Image:          d.cfg.Image,
		Memory:         d.cfg.Memory,
		CPU:            d.cfg.CPU,
		Platform:       d.cfg.Platform,
		Network:        d.cfg.Network,
		DisableNetwork: d.cfg.DisableNetwork,
		Privileged:     true,
	})

	return
}
