package wsl

// Config is the configuration for the wsl engine.
type Config struct {
	Command     string
	Environment map[string]string
	WorkDir     string
	User        string
	Shell       string
	// ReadOnly means none-interactive for terminal, which is used for show log, like top
	ReadOnly bool

	// engine = wsl (Windows only)
	// WSLDistro is the WSL distribution name (optional, e.g. "Ubuntu")
	WSLDistro string

	// Custom Command Runner ID
	ID string

	// AllowedSystemEnvKeys is the allowed system environment keys, which will be inherited to the command
	AllowedSystemEnvKeys []string
}
