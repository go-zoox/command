package wsl

// Cancel kills the wsl process.
func (w *wsl) Cancel() error {
	if w.cmd == nil || w.cmd.Process == nil {
		return nil
	}
	return w.cmd.Process.Kill()
}
