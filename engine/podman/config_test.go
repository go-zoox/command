package podman

import (
	"testing"
)

func TestConfig_PodmanHostDefault(t *testing.T) {
	cfg := &Config{
		Command: "echo test",
	}

	if cfg.PodmanHost != "" {
		t.Errorf("expected empty PodmanHost by default, got %q", cfg.PodmanHost)
	}
}

func TestConfig_PodmanHostCustom(t *testing.T) {
	cfg := &Config{
		PodmanHost: "unix:///run/user/1000/podman/podman.sock",
		Image:      "alpine:latest",
	}

	if cfg.PodmanHost != "unix:///run/user/1000/podman/podman.sock" {
		t.Errorf("expected custom PodmanHost, got %q", cfg.PodmanHost)
	}
}

func TestConfig_ImageAndResources(t *testing.T) {
	cfg := &Config{
		Image:  "docker.io/library/alpine:latest",
		Memory: 512,
		CPU:    1.0,
	}

	if cfg.Image != "docker.io/library/alpine:latest" {
		t.Errorf("expected Image docker.io/library/alpine:latest, got %q", cfg.Image)
	}
	if cfg.Memory != 512 {
		t.Errorf("expected Memory 512, got %d", cfg.Memory)
	}
	if cfg.CPU != 1.0 {
		t.Errorf("expected CPU 1.0, got %f", cfg.CPU)
	}
}
