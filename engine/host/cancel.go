package host

import "errors"

func (h *host) Cancel() error {
	if h.cmd == nil {
		return errors.New("command: not started")
	}

	if err := h.cmd.Process.Kill(); err != nil {
		return err
	}

	return nil
}
