package caas

import (
	"io"
	"os"

	"github.com/go-zoox/command/engine"
	cs "github.com/go-zoox/commands-as-a-service/client"
)

// Name is the name of the engine.
const Name = "caas"

type caas struct {
	//
	cfg *Config
	//
	client cs.Client

	//
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
}

func New(cfg *Config) (engine.Engine, error) {
	c := &caas{
		cfg: cfg,
		//
		stdin:  os.Stdin,
		stdout: os.Stdout,
		stderr: os.Stderr,
	}

	if err := c.create(); err != nil {
		return nil, err
	}

	return c, nil
}
