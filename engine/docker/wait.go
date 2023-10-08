package docker

import (
	"context"
	"fmt"
	"io"

	"github.com/docker/docker/api/types/container"
	"github.com/go-zoox/command/errors"
)

// Wait waits for the command to finish.
func (d *docker) Wait() error {
	result, err := d.client.ContainerWait(context.Background(), d.container.ID, container.WaitConditionNotRunning)
	select {
	case err := <-err:
		if err != nil && err != io.EOF {
			return fmt.Errorf("container exit error: %#v", err)
		}
	case result := <-result:
		if result.StatusCode != 0 {
			// return fmt.Errorf("container exited with non-zero status: %d", result.StatusCode)
			return &errors.ExitError{
				Code:    int(result.StatusCode),
				Message: fmt.Sprintf("container exited with non-zero status: %d", result.StatusCode),
			}
		}
	}

	return nil
}
