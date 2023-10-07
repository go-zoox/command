package docker

import "io"

func (d *docker) SetStdin(stdin io.Reader) error {
	d.stdin = stdin
	return nil
}

func (d *docker) SetStdout(stdout io.Writer) error {
	d.stdout = stdout
	return nil
}

func (d *docker) SetStderr(stderr io.Writer) error {
	d.stderr = stderr
	return nil
}
