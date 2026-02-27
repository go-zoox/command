package k8s

// Config is the configuration for the k8s engine.
type Config struct {
	Command     string
	Environment map[string]string
	WorkDir     string
	User        string
	Shell       string
	// ReadOnly means none-interactive for terminal, which is used for show log, like top
	ReadOnly bool

	// engine = k8s
	// Kubeconfig is the path to kubeconfig file (optional, uses in-cluster or default rules if empty)
	Kubeconfig string
	// Namespace is the Kubernetes namespace to run the Job in
	Namespace string
	// Image is the container image for the Job
	Image string
	// JobTimeoutSeconds is the optional timeout for the Job (0 = no timeout)
	JobTimeoutSeconds int64

	// Custom Command Runner ID (used as Job name prefix)
	ID string

	// AllowedSystemEnvKeys is the allowed system environment keys, which will be inherited to the command
	AllowedSystemEnvKeys []string
}
