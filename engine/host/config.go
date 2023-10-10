package host

// Config is the configuration for a host engine.
type Config struct {
	Command     string
	Environment map[string]string
	WorkDir     string
	User        string
	Shell       string
	// ReadOnly means none-interactive for terminal, which is used for show log, like top
	ReadOnly bool

	//
	IsHistoryDisabled bool

	// Custom Command Runner ID
	ID string
}
