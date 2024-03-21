package docker

import (
	"context"

	"github.com/docker/docker/api/types/container"
)

// Cancel cancels the command.
func (d *docker) Cancel() error {
	return d.client.ContainerRemove(context.Background(), d.container.ID, container.RemoveOptions{
		Force: true,
	})
}
