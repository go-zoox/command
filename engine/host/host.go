package host

import (
	"io"
	"os"
	"os/exec"

	"github.com/go-zoox/command/engine"
)

const Engine = "host"

type host struct {
	cfg *Config
	//
	cmd *exec.Cmd
	//

	//
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
}

func New(cfg *Config) (engine.Engine, error) {
	h := &host{
		cfg: cfg,
		//
		stdin:  os.Stdin,
		stdout: os.Stdout,
		stderr: os.Stderr,
	}

	return h, nil
}
