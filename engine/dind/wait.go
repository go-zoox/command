package dind

// Wait waits for the command to finish.
func (d *dind) Wait() error {
	return d.client.Wait()
}
