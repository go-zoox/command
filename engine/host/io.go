package host

import "io"

// SetStdin sets the stdin for the command.
func (h *host) SetStdin(stdin io.Reader) error {
	h.stdin = stdin
	return nil
}

// SetStdout sets the stdout for the command.
func (h *host) SetStdout(stdout io.Writer) error {
	h.stdout = stdout
	return nil
}

// SetStderr sets the stderr for the command.
func (h *host) SetStderr(stderr io.Writer) error {
	h.stderr = stderr
	return nil
}
