package docker

import "io"

// SetStdin sets the stdin for the Docker engine.
func (d *docker) SetStdin(stdin io.Reader) error {
	d.stdin = stdin
	return nil
}

// SetStdout sets the stdout for the Docker engine.
func (d *docker) SetStdout(stdout io.Writer) error {
	d.stdout = stdout
	return nil
}

// SetStderr sets the stderr for the Docker engine.
func (d *docker) SetStderr(stderr io.Writer) error {
	d.stderr = stderr
	return nil
}
