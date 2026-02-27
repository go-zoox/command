package errors

import (
	"testing"
)

func TestExitError_Error(t *testing.T) {
	e := &ExitError{
		Code:    1,
		Message: "exit status 1",
	}
	if e.Error() != "exit status 1" {
		t.Errorf("Error() = %q, want %q", e.Error(), "exit status 1")
	}
}

func TestExitError_ExitCode(t *testing.T) {
	e := &ExitError{
		Code:    42,
		Message: "failed",
	}
	if e.ExitCode() != 42 {
		t.Errorf("ExitCode() = %d, want 42", e.ExitCode())
	}
}
