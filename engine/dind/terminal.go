package dind

import (
	"github.com/go-zoox/command/terminal"
)

// Terminal returns a terminal.
func (d *dind) Terminal() (terminal.Terminal, error) {
	return d.client.Terminal()
}
