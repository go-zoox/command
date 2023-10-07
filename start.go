package command

import "errors"

func (c *command) Start() error {
	if c.engine == nil {
		return errors.New("engine not set")
	}

	return c.engine.Start()
}
