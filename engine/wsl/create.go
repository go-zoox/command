package wsl

// create builds the wsl command arguments (exec.Cmd is created in Start).
func (w *wsl) create() error {
	// wsl [-d Distro] -e shell -c "command"
	if w.cfg.WSLDistro != "" {
		w.args = append(w.args, "-d", w.cfg.WSLDistro)
	}
	w.args = append(w.args, "-e", w.cfg.Shell)
	if w.cfg.Command != "" {
		w.args = append(w.args, "-c", w.cfg.Command)
	} else {
		w.args = append(w.args, "-c", "sleep 0")
	}
	return nil
}
