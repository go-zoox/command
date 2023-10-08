package docker

import (
	"context"

	"github.com/docker/docker/api/types"
)

// Cancel cancels the command.
func (d *docker) Cancel() error {
	return d.client.ContainerRemove(context.Background(), d.container.ID, types.ContainerRemoveOptions{
		Force: true,
	})
}
