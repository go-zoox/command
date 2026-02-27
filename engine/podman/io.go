package podman

import "io"

// SetStdin sets the stdin for the podman engine.
func (p *podman) SetStdin(stdin io.Reader) error {
	p.stdin = stdin
	return nil
}

// SetStdout sets the stdout for the podman engine.
func (p *podman) SetStdout(stdout io.Writer) error {
	p.stdout = stdout
	return nil
}

// SetStderr sets the stderr for the podman engine.
func (p *podman) SetStderr(stderr io.Writer) error {
	p.stderr = stderr
	return nil
}
