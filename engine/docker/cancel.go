package docker

import (
	"github.com/docker/docker/api/types"
)

// Cancel cancels the command.
func (d *docker) Cancel() error {
	return d.client.ContainerRemove(d.ctx, d.container.ID, types.ContainerRemoveOptions{
		Force: true,
	})
}
