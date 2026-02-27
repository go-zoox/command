package wsl

import (
	"os/exec"

	"github.com/go-zoox/command/errors"
)

// Wait waits for the command to finish.
func (w *wsl) Wait() error {
	if err := w.cmd.Wait(); err != nil {
		if v, ok := err.(*exec.ExitError); ok {
			return &errors.ExitError{
				Code:    v.ExitCode(),
				Message: v.Error(),
			}
		}
		return &errors.ExitError{
			Code:    1,
			Message: err.Error(),
		}
	}
	return nil
}
