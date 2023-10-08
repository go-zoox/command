package host

// Cancel cancels the command.
func (h *host) Cancel() error {
	if err := h.cmd.Process.Kill(); err != nil {
		return err
	}

	return nil
}
