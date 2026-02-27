package engine

import (
	"testing"
)

// TestGet_UnknownEngine verifies that an unknown engine returns ErrEngineNotFound.
// Known engines (host, docker, k8s, etc.) are registered in the main command
// package's init(), so they are only available when that package is imported.
func TestGet_UnknownEngine(t *testing.T) {
	_, err := Get("unknown-engine-name-xyz")
	if err == nil {
		t.Fatal("expected error for unknown engine")
	}
	if err != ErrEngineNotFound {
		t.Errorf("expected ErrEngineNotFound, got %v", err)
	}
}

func TestGet_EmptyName(t *testing.T) {
	_, err := Get("")
	if err == nil {
		t.Fatal("expected error for empty engine name")
	}
	if err != ErrEngineNotFound {
		t.Errorf("expected ErrEngineNotFound, got %v", err)
	}
}
