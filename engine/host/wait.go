package host

// Wait waits for the command to finish.
func (h *host) Wait() error {
	return h.cmd.Wait()
}
