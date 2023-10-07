package command

import "io"

func (c *command) SetStdin(stdin io.Reader) error {
	return c.engine.SetStdin(stdin)
}

func (c *command) SetStdout(stdout io.Writer) error {
	return c.engine.SetStdout(stdout)
}

func (c *command) SetStderr(stderr io.Writer) error {
	return c.engine.SetStderr(stderr)
}
