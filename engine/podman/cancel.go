package podman

import (
	"context"

	"github.com/docker/docker/api/types/container"
)

// Cancel cancels the command.
func (p *podman) Cancel() error {
	return p.client.ContainerRemove(context.Background(), p.container.ID, container.RemoveOptions{
		Force: true,
	})
}
