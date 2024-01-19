package caas

import "github.com/go-zoox/commands-as-a-service/entities"

// Wait waits for the command to finish.
func (c *caas) Wait() error {
	return c.client.Exec(&entities.Command{
		ID:          c.cfg.ID,
		Script:      c.cfg.Command,
		Environment: c.cfg.Environment,
		// WorkDir:     c.cfg.WorkDir,
		User: c.cfg.User,
		// Shell:       c.cfg.Shell,
	})
}
