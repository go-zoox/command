package ssh

import (
	"io"

	"github.com/go-zoox/command/terminal"
	sshx "golang.org/x/crypto/ssh"
)

// Terminal returns a terminal.
func (s *ssh) Terminal() (terminal.Terminal, error) {
	s.session.Stdout = s.stdout
	s.session.Stderr = s.stderr

	sessionStdin, err := s.session.StdinPipe()
	if err != nil {
		return nil, err
	}

	modes := sshx.TerminalModes{
		sshx.ECHO:          0,     // disable echoing
		sshx.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		sshx.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	if err := s.session.RequestPty("xterm", 100, 300, modes); err != nil {
		return nil, err
	}

	if err := s.session.Shell(); err != nil {
		return nil, err
	}

	return &Terminal{
		Session:      s.session,
		SessionStdin: sessionStdin,
	}, nil
}

// Terminal is a terminal.
type Terminal struct {
	terminal.Terminal

	Session *sshx.Session

	// @TODO
	SessionStdin io.WriteCloser
}

// Resize resizes the terminal.
func (t *Terminal) Resize(rows, cols int) error {
	return t.Session.WindowChange(rows, cols)
}

// Close closes the terminal.
func (t *Terminal) Close() error {
	return t.Session.Close()
}

// Read reads from the terminal.
func (t *Terminal) Read(p []byte) (n int, err error) {
	// @TODO
	sessionStdout, err := t.Session.StdoutPipe()
	if err != nil {
		return 0, err
	}

	return sessionStdout.Read(p)
}

// Write writes to the terminal.
func (t *Terminal) Write(p []byte) (n int, err error) {
	return t.SessionStdin.Write(p)
}

// ExitCode returns the exit code.
func (t *Terminal) ExitCode() int {
	return 1
}
