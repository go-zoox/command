package config

import (
	"context"
	"testing"
)

func TestConfig_ZeroValue(t *testing.T) {
	var cfg Config
	if cfg.Command != "" {
		t.Errorf("zero Config should have empty Command, got %q", cfg.Command)
	}
	if cfg.Engine != "" {
		t.Errorf("zero Config should have empty Engine, got %q", cfg.Engine)
	}
}

func TestConfig_WithContext(t *testing.T) {
	ctx := context.Background()
	cfg := Config{
		Context: ctx,
		Command: "echo ok",
	}
	if cfg.Context != ctx {
		t.Error("Context not set correctly")
	}
	if cfg.Command != "echo ok" {
		t.Errorf("Command = %q, want echo ok", cfg.Command)
	}
}

func TestConfig_EngineFields(t *testing.T) {
	cfg := Config{
		Engine:  "docker",
		Image:   "alpine:latest",
		Memory:  512,
		CPU:     1.0,
		Sandbox: true,
	}
	if cfg.Engine != "docker" || cfg.Image != "alpine:latest" || cfg.Memory != 512 || cfg.CPU != 1.0 || !cfg.Sandbox {
		t.Errorf("engine fields not set: Engine=%q Image=%q Memory=%d CPU=%f Sandbox=%t",
			cfg.Engine, cfg.Image, cfg.Memory, cfg.CPU, cfg.Sandbox)
	}
}

func TestConfig_K8sFields(t *testing.T) {
	cfg := Config{
		Engine:               "k8s",
		K8sKubeconfig:        "/path/to/kubeconfig",
		K8sNamespace:         "default",
		K8sImage:             "busybox",
		K8sPodTimeoutSeconds: 300,
	}
	if cfg.K8sKubeconfig != "/path/to/kubeconfig" || cfg.K8sNamespace != "default" ||
		cfg.K8sImage != "busybox" || cfg.K8sPodTimeoutSeconds != 300 {
		t.Errorf("k8s fields not set correctly")
	}
}

func TestConfig_WSLAndDockerRuntime(t *testing.T) {
	cfg := Config{
		Engine:         "wsl",
		WSLDistro:      "Ubuntu",
		DockerRuntime:  "runsc",
	}
	if cfg.WSLDistro != "Ubuntu" || cfg.DockerRuntime != "runsc" {
		t.Errorf("WSLDistro=%q DockerRuntime=%q", cfg.WSLDistro, cfg.DockerRuntime)
	}
}
