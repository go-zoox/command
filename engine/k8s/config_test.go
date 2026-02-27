package k8s

import (
	"testing"
)

func TestConfig_Defaults(t *testing.T) {
	cfg := &Config{
		Namespace: "myns",
		Image:     "busybox:1.36",
	}

	if cfg.Namespace != "myns" {
		t.Errorf("expected Namespace myns, got %q", cfg.Namespace)
	}
	if cfg.Image != "busybox:1.36" {
		t.Errorf("expected Image busybox:1.36, got %q", cfg.Image)
	}
}

func TestConfig_KubeconfigAndNamespace(t *testing.T) {
	cfg := &Config{
		Kubeconfig: "/path/to/kubeconfig",
		Namespace:  "default",
		Command:    "echo hello",
	}

	if cfg.Kubeconfig != "/path/to/kubeconfig" {
		t.Errorf("expected Kubeconfig /path/to/kubeconfig, got %q", cfg.Kubeconfig)
	}
	if cfg.Namespace != "default" {
		t.Errorf("expected Namespace default, got %q", cfg.Namespace)
	}
}

func TestConfig_JobTimeoutSeconds(t *testing.T) {
	cfg := &Config{
		JobTimeoutSeconds: 120,
	}

	if cfg.JobTimeoutSeconds != 120 {
		t.Errorf("expected JobTimeoutSeconds 120, got %d", cfg.JobTimeoutSeconds)
	}
}
