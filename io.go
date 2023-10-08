package command

import "io"

// SetStdin sets the stdin for the command.
func (c *command) SetStdin(stdin io.Reader) error {
	return c.engine.SetStdin(stdin)
}

// SetStdout sets the stdout for the command.
func (c *command) SetStdout(stdout io.Writer) error {
	return c.engine.SetStdout(stdout)
}

// SetStderr sets the stderr for the command.
func (c *command) SetStderr(stderr io.Writer) error {
	return c.engine.SetStderr(stderr)
}
