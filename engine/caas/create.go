package caas

import (
	"errors"

	cs "github.com/go-zoox/commands-as-a-service/client"
)

// create creates the engine.
func (c *caas) create() error {
	if c.client != nil {
		return errors.New("command: already created")
	}

	if c.cfg.Server == "" {
		return errors.New("command: server is required")
	}

	c.client = cs.New(&cs.Config{
		Server:       c.cfg.Server,
		ClientID:     c.cfg.ClientID,
		ClientSecret: c.cfg.ClientSecret,
		Stdout:       c.stdout,
		Stderr:       c.stderr,
	})

	return nil
}
