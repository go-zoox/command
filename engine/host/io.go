package host

import "io"

func (h *host) SetStdin(stdin io.Reader) error {
	h.stdin = stdin
	return nil
}

func (h *host) SetStdout(stdout io.Writer) error {
	h.stdout = stdout
	return nil
}

func (h *host) SetStderr(stderr io.Writer) error {
	h.stderr = stderr
	return nil
}
