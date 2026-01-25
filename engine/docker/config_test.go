package docker

import (
	"testing"
)

func TestConfig_SandboxField(t *testing.T) {
	cfg := &Config{
		Sandbox: true,
	}

	if !cfg.Sandbox {
		t.Error("expected Sandbox to be true")
	}

	cfg.Sandbox = false
	if cfg.Sandbox {
		t.Error("expected Sandbox to be false")
	}
}

func TestConfig_SandboxWithOtherSettings(t *testing.T) {
	cfg := &Config{
		Command:      "echo test",
		Sandbox:      true,
		Privileged:   false,
		DisableNetwork: true,
		Memory:       512,
		CPU:          1.0,
	}

	if !cfg.Sandbox {
		t.Error("expected Sandbox to be true")
	}

	if cfg.Privileged {
		t.Error("expected Privileged to be false in sandbox mode")
	}

	if !cfg.DisableNetwork {
		t.Error("expected DisableNetwork to be true in sandbox mode")
	}

	if cfg.Memory != 512 {
		t.Errorf("expected Memory to be 512, got %d", cfg.Memory)
	}

	if cfg.CPU != 1.0 {
		t.Errorf("expected CPU to be 1.0, got %f", cfg.CPU)
	}
}
