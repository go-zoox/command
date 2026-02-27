package k8s

import (
	"io"
	"os"

	"github.com/go-zoox/command/engine"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// Name is the name of the engine.
const Name = "k8s"

type k8s struct {
	cfg *Config
	//
	clientset    kubernetes.Interface
	restConfig   *rest.Config
	jobNamespace string
	jobName      string
	//
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
}

// New creates a new k8s engine.
func New(cfg *Config) (engine.Engine, error) {
	if cfg.Shell == "" {
		cfg.Shell = "/bin/sh"
	}
	if cfg.Image == "" {
		cfg.Image = "alpine:latest"
	}
	if cfg.Namespace == "" {
		cfg.Namespace = "default"
	}

	k := &k8s{
		cfg:    cfg,
		stdin:  os.Stdin,
		stdout: os.Stdout,
		stderr: os.Stderr,
	}

	if err := k.create(); err != nil {
		return nil, err
	}

	return k, nil
}
