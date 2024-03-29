package engine

import (
	"io"

	"github.com/go-zoox/command/terminal"
)

// Engine is the interface that an command engine must implement.
type Engine interface {
	Start() error
	Wait() error
	Cancel() error
	//
	SetStdin(stdin io.Reader) error
	SetStdout(stdout io.Writer) error
	SetStderr(stderr io.Writer) error
	//
	Terminal() (terminal.Terminal, error)
}
