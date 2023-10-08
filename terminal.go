package command

import "github.com/go-zoox/command/terminal"

// Terminal returns a terminal for the command.
func (c *command) Terminal() (terminal.Terminal, error) {
	return c.engine.Terminal()
}
