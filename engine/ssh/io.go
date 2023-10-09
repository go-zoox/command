package ssh

import "io"

// SetStdin sets the stdin for the command.
func (s *ssh) SetStdin(stdin io.Reader) error {
	s.stdin = stdin
	return nil
}

// SetStdout sets the stdout for the command.
func (s *ssh) SetStdout(stdout io.Writer) error {
	s.stdout = stdout
	return nil
}

// SetStderr sets the stderr for the command.
func (s *ssh) SetStderr(stderr io.Writer) error {
	s.stderr = stderr
	return nil
}
