package command

import "errors"

// Start starts to run the command.
func (c *command) Start() error {
	if c.engine == nil {
		return errors.New("engine not set")
	}

	if c.cfg.Command == "" {
		return errors.New("command is required")
	}

	return c.engine.Start()
}
