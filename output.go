package command

import (
	"bytes"
)

// Output gets the command output.
func (c *command) Output() ([]byte, error) {
	var stdout bytes.Buffer
	// var stderr bytes.Buffer
	c.SetStdout(&stdout)
	c.SetStderr(&stdout)

	if err := c.Run(); err != nil {
		return nil, err
	}

	return stdout.Bytes(), nil
}
