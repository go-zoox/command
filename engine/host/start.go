package host

import (
	"io"
	"os/exec"
)

// Start starts the command.
func (h *host) Start() error {
	if err := applyStdin(h.cmd, h.stdin); err != nil {
		return nil
	}

	if err := applyStdout(h.cmd, h.stdout); err != nil {
		return nil
	}

	if err := applyStderr(h.cmd, h.stderr); err != nil {
		return nil
	}

	return h.cmd.Start()
}

func applyStdin(cmd *exec.Cmd, stdin io.Reader) error {
	cmd.Stdin = stdin
	return nil
}

func applyStdout(cmd *exec.Cmd, stdout io.Writer) error {
	cmd.Stdout = stdout
	if cmd.Stderr == nil {
		return applyStderr(cmd, stdout)
	}

	return nil
}

func applyStderr(cmd *exec.Cmd, stderr io.Writer) error {
	cmd.Stderr = stderr
	return nil
}
