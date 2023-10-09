package ssh

// Cancel cancels the command.
func (s *ssh) Cancel() error {
	if s.session != nil {
		if err := s.session.Close(); err != nil {
			return err
		}
	}

	if s.client != nil {
		if err := s.client.Close(); err != nil {
			return err
		}
	}

	return nil
}
