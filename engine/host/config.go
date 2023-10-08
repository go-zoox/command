package host

// Config is the configuration for a host engine.
type Config struct {
	Command     string
	Environment map[string]string
	WorkDir     string
	User        string
	Shell       string

	//
	IsHistoryDisabled bool
}
