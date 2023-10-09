package ssh

// Wait waits for the command to exit.
func (s *ssh) Wait() error {
	return s.session.Wait()
}
