package caas

func (c *caas) Cancel() error {
	return c.client.Close()
}
