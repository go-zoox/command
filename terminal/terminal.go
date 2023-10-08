package terminal

import "io"

type Terminal interface {
	io.ReadWriteCloser

	//
	Resize(rows, cols int) error

	//
	ExitCode() int

	//
	Wait() error
}
