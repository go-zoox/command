package ssh

import (
	"io"
	"os"

	"github.com/go-zoox/command/engine"
	sshx "golang.org/x/crypto/ssh"
)

// Name is the name of the engine.
const Name = "ssh"

// Config is the config for the ssh engine.
type Config struct {
	Command     string
	Environment map[string]string
	WorkDir     string
	Shell       string
	// ReadOnly means none-interactive for terminal, which is used for show log, like top
	ReadOnly bool
	//
	Host             string
	Port             int
	User             string
	Pass             string
	PrivateKey       string
	PrivateKeySecret string
	//
	IsIgnoreStrictHostKeyChecking bool
	//
	KnowHostsFilePath string

	//
	ID string
}

type ssh struct {
	cfg *Config
	//
	client *sshx.Client
	//
	session *sshx.Session
	//
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
}

// New creates a new ssh engine.
func New(cfg *Config) (engine.Engine, error) {
	if cfg.Shell == "" {
		cfg.Shell = "/bin/sh"
	}

	s := &ssh{
		cfg: cfg,
		//
		stdin:  os.Stdin,
		stdout: os.Stdout,
		stderr: os.Stderr,
	}

	if err := s.create(); err != nil {
		return nil, err
	}

	return s, nil
}
