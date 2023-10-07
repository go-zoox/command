package engine

import "io"

type Engine interface {
	Start() error
	Wait() error
	Cancel() error
	//
	SetStdin(stdin io.Reader) error
	SetStdout(stdout io.Writer) error
	SetStderr(stderr io.Writer) error
}

type Config struct {
	Shell       string
	Command     string
	Environment map[string]string
	WorkDir     string
}
