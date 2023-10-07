package host

func (h *host) Wait() error {
	return h.cmd.Wait()
}
