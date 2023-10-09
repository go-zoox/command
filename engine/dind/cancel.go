package dind

// Cancel cancels the command.
func (d *dind) Cancel() error {
	return d.client.Cancel()
}
