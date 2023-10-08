package terminal

import "io"

// Terminal is the interface that a terminal must implement.
type Terminal interface {
	io.ReadWriteCloser

	//
	Resize(rows, cols int) error

	//
	ExitCode() int

	//
	Wait() error
}
