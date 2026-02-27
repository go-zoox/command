package wsl

import "io"

// SetStdin sets the stdin for the wsl engine.
func (w *wsl) SetStdin(stdin io.Reader) error {
	w.stdin = stdin
	return nil
}

// SetStdout sets the stdout for the wsl engine.
func (w *wsl) SetStdout(stdout io.Writer) error {
	w.stdout = stdout
	return nil
}

// SetStderr sets the stderr for the wsl engine.
func (w *wsl) SetStderr(stderr io.Writer) error {
	w.stderr = stderr
	return nil
}
