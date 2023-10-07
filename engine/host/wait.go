package host

import "errors"

func (h *host) Wait() error {
	if h.cmd == nil {
		return errors.New("command: not started")
	}

	return h.cmd.Wait()
}
