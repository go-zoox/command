package host

func (h *host) Start() error {
	if err := applyStdin(h.cmd, h.stdin); err != nil {
		return nil
	}

	if err := applyStdout(h.cmd, h.stdout); err != nil {
		return nil
	}

	if err := applyStderr(h.cmd, h.stderr); err != nil {
		return nil
	}

	return h.cmd.Start()
}
