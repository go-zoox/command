package command

import (
	"io"

	cio "github.com/go-zoox/core-utils/io"
)

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

// SetStdinWrapFunc sets the stdin wrap function for the command.
func (c *command) SetStdinWrapFunc(stdinFunc func(b []byte) (n int, err error)) error {
	return c.SetStdin(cio.ReadWrapFunc(stdinFunc))
}

// SetStdoutWrapFunc sets the stdout wrap function for the command.
func (c *command) SetStdoutWrapFunc(stdoutFunc func(b []byte) (n int, err error)) error {
	return c.SetStdout(cio.WriterWrapFunc(stdoutFunc))
}

// SetStderrWrapFunc sets the stderr wrap function for the command.
func (c *command) SetStderrWrapFunc(stderrFunc func(b []byte) (n int, err error)) error {
	return c.SetStderr(cio.WriterWrapFunc(stderrFunc))
}
