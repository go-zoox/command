package host

import (
	"os/exec"

	"github.com/go-zoox/command/errors"
)

// Wait waits for the command to finish.
func (h *host) Wait() error {
	if err := h.cmd.Wait(); err != nil {
		v, ok := err.(*exec.ExitError)
		if !ok {
			return &errors.ExitError{
				Code:    1,
				Message: err.Error(),
			}
		}

		return &errors.ExitError{
			Code:    v.ExitCode(),
			Message: v.Error(),
		}
	}

	return nil
}
