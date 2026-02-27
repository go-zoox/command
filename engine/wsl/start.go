package wsl

import (
	"fmt"
	"os"
	"os/exec"
)

// Start starts the wsl command.
func (w *wsl) Start() error {
	w.cmd = exec.Command("wsl", w.args...)

	if err := applyEnv(w.cmd, w.cfg.Environment, w.cfg.AllowedSystemEnvKeys); err != nil {
		return err
	}
	if w.cfg.WorkDir != "" {
		w.cmd.Dir = w.cfg.WorkDir
	}

	w.cmd.Stdin = w.stdin
	w.cmd.Stdout = w.stdout
	w.cmd.Stderr = w.stderr

	return w.cmd.Start()
}

func applyEnv(cmd *exec.Cmd, environment map[string]string, allowedSystemEnvKeys []string) error {
	cmd.Env = append([]string{}, "TERM=xterm")
	for _, key := range allowedSystemEnvKeys {
		if value, ok := os.LookupEnv(key); ok {
			cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, value))
		}
	}
	for k, v := range environment {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}
	return nil
}
