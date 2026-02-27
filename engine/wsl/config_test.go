package wsl

import (
	"runtime"
	"testing"
)

func TestNew_NonWindows(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skip on Windows")
	}
	_, err := New(&Config{
		Command: "echo ok",
		Shell:   "/bin/sh",
	})
	if err != ErrNotWindows {
		t.Errorf("expected ErrNotWindows on non-Windows, got %v", err)
	}
}

func TestConfig_Defaults(t *testing.T) {
	cfg := &Config{
		Command: "echo test",
		Shell:   "/bin/sh",
	}
	if cfg.Shell != "/bin/sh" {
		t.Errorf("expected Shell /bin/sh, got %q", cfg.Shell)
	}
}

func TestConfig_WSLDistro(t *testing.T) {
	cfg := &Config{
		WSLDistro: "Ubuntu",
	}
	if cfg.WSLDistro != "Ubuntu" {
		t.Errorf("expected WSLDistro Ubuntu, got %q", cfg.WSLDistro)
	}
}
