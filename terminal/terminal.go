package terminal

import "io"

type Terminal interface {
	io.ReadWriteCloser

	//
	Resize(rows, cols int) error
	//
	Wait() error

	//
	ExitCode() int
}
