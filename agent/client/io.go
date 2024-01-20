package client

import "io"

func (c *client) SetStdin(stdin io.Reader) error {
	c.stdin = stdin
	return nil
}

func (c *client) SetStdout(stdout io.Writer) error {
	c.stdout = stdout
	return nil
}

func (c *client) SetStderr(stderr io.Writer) error {
	c.stderr = stderr
	return nil
}
