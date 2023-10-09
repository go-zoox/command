package caas

type Config struct {
	Command     string
	Environment map[string]string
	WorkDir     string
	User        string
	Shell       string

	Server       string
	ClientID     string
	ClientSecret string

	// Custom Command Runner ID
	ID string
}
