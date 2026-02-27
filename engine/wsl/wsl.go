package wsl

import (
	"errors"
	"io"
	"os"
	"os/exec"
	"runtime"

	"github.com/go-zoox/command/engine"
)

// Name is the name of the engine.
const Name = "wsl"

// ErrNotWindows is returned when the wsl engine is used on a non-Windows system.
var ErrNotWindows = errors.New("wsl engine is only available on Windows")

type wsl struct {
	cfg *Config
	//
	cmd  *exec.Cmd
	args []string
	//
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
}

// New creates a new wsl engine. Returns an error if not on Windows.
func New(cfg *Config) (engine.Engine, error) {
	if runtime.GOOS != "windows" {
		return nil, ErrNotWindows
	}
	if cfg.Shell == "" {
		cfg.Shell = "/bin/sh"
	}

	w := &wsl{
		cfg:    cfg,
		stdin:  os.Stdin,
		stdout: os.Stdout,
		stderr: os.Stderr,
	}

	if err := w.create(); err != nil {
		return nil, err
	}

	return w, nil
}
