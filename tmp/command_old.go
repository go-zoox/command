package command

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
)

// command is better os/exec command
type command struct {
	Script      string            `json:"content"`
	Context     string            `json:"context"`
	Environment map[string]string `json:"environment"`
	Shell       string            `json:"shell"`
	//
	Stdout io.Writer
	Stderr io.Writer
	//
	cmd *exec.Cmd
}

// Run runs the command
func (c *command) Run() error {
	environment := os.Environ()

	for k, v := range c.Environment {
		environment = append(environment, fmt.Sprintf("%s=%s", k, v))
	}

	shell := c.Shell
	if shell == "" {
		shell = os.Getenv("SHELL")
		if shell == "" {
			shell = "sh"
		}
	}

	cmd := exec.Command(shell, "-c", c.Script)
	cmd.Dir = c.Context
	cmd.Env = environment

	cmd.Stdout = c.Stdout
	if cmd.Stdout == nil {
		cmd.Stdout = os.Stdout
	}

	cmd.Stderr = c.Stderr
	if cmd.Stderr == nil {
		cmd.Stderr = os.Stderr
	}

	c.cmd = cmd

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

// Config gets the config of command
func (c *command) Config() (string, error) {
	cfg, err := json.MarshalIndent(c, "", " ")
	if err != nil {
		return "", err
	}

	return string(cfg), nil
}

// MustConfig gets the config of command
func (c *command) MustConfig() string {
	cfg, err := c.Config()
	if err != nil {
		return ""
	}

	return cfg
}

// ExitCode gets the exit code of process
func (c *command) ExitCode() int {
	return c.cmd.ProcessState.ExitCode()
}
