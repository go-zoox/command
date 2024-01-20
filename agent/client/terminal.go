package client

import (
	"fmt"

	"github.com/go-zoox/command/terminal"
)

func (c *client) Terminal() (terminal.Terminal, error) {
	return nil, fmt.Errorf("not implemented in agent")
}
