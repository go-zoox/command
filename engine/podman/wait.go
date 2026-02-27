package podman

import (
	"context"
	"fmt"
	"io"

	"github.com/docker/docker/api/types/container"
	"github.com/go-zoox/command/errors"
)

// Wait waits for the command to finish.
func (p *podman) Wait() error {
	resultC, errC := p.client.ContainerWait(context.Background(), p.container.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errC:
		if err != nil && err != io.EOF {
			return fmt.Errorf("podman: container wait: %w", err)
		}
	case result := <-resultC:
		if result.StatusCode != 0 {
			return &errors.ExitError{
				Code:    int(result.StatusCode),
				Message: fmt.Sprintf("container exited with non-zero status: %d", result.StatusCode),
			}
		}
	}

	return nil
}
