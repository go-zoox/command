package host

type Config struct {
	Command     string
	Environment map[string]string
	WorkDir     string
	User        string
	Shell       string

	//
	IsHistoryDisabled bool
}
