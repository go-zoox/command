package command

import "github.com/go-zoox/command/terminal"

func (c *command) Terminal() (terminal.Terminal, error) {
	return c.engine.Terminal()
}
