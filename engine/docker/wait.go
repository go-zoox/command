package docker

import (
	"fmt"
	"io"

	"github.com/docker/docker/api/types/container"
)

func (d *docker) Wait() error {
	result, err := d.client.ContainerWait(d.ctx, d.container.ID, container.WaitConditionNotRunning)
	select {
	case err := <-err:
		if err != nil && err != io.EOF {
			return fmt.Errorf("container exit error: %#v", err)
		}
	case result := <-result:
		if result.StatusCode != 0 {
			return fmt.Errorf("container exited with non-zero status: %d", result.StatusCode)
		}
	}

	return nil
}
