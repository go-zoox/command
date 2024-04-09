package ssh

import (
	"io"
	"os"

	sshx "golang.org/x/crypto/ssh"
)

// Start starts the command without waiting for it to finish.
func (s *ssh) Start() error {
	if err := applyStdin(s.session, s.stdin); err != nil {
		return nil
	}

	if err := applyStdout(s.session, s.stdout); err != nil {
		return nil
	}

	if err := applyStderr(s.session, s.stderr); err != nil {
		return nil
	}

	if len(s.cfg.AllowedSystemEnvKeys) != 0 {
		for _, key := range s.cfg.AllowedSystemEnvKeys {
			if value, ok := os.LookupEnv(key); ok {
				s.session.Setenv(key, value)
			}
		}

	}

	for k, v := range s.cfg.Environment {
		s.session.Setenv(k, v)
	}

	return s.session.Start(s.cfg.Command)
}

func applyStdin(cmd *sshx.Session, stdin io.Reader) error {
	cmd.Stdin = stdin
	return nil
}

func applyStdout(cmd *sshx.Session, stdout io.Writer) error {
	cmd.Stdout = stdout
	if cmd.Stderr == nil {
		return applyStderr(cmd, stdout)
	}

	return nil
}

func applyStderr(cmd *sshx.Session, stderr io.Writer) error {
	cmd.Stderr = stderr
	return nil
}
