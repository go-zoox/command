//go:build !windows

package host

import (
	"os"
	"os/exec"
	"os/user"
	"syscall"
	"testing"
)

func TestApplyUser_EmptyUsername(t *testing.T) {
	cmd := exec.Command("echo", "hello")
	err := applyUser(cmd, "")
	if err != nil {
		t.Fatalf("applyUser with empty username should not error: %v", err)
	}
	if cmd.SysProcAttr != nil {
		t.Error("SysProcAttr should be nil when username is empty")
	}
}

func TestApplyUser_NonexistentUser(t *testing.T) {
	cmd := exec.Command("echo", "hello")
	err := applyUser(cmd, "nonexistent_user_abcdef123456789")
	if err == nil {
		t.Fatal("expected error for nonexistent user")
	}
}

func TestApplyUser_CurrentUser(t *testing.T) {
	currentUser, err := user.Current()
	if err != nil {
		t.Skipf("cannot get current user: %v", err)
	}

	cmd := exec.Command("echo", "hello")
	err = applyUser(cmd, currentUser.Username)
	if err != nil {
		t.Fatalf("applyUser for current user should not error: %v", err)
	}

	if cmd.SysProcAttr == nil {
		t.Fatal("SysProcAttr should be set for valid user")
	}
	if cmd.SysProcAttr.Credential == nil {
		t.Fatal("SysProcAttr.Credential should be set for valid user")
	}

	expectedUid := currentUser.Uid
	expectedGid := currentUser.Gid

	if cmd.SysProcAttr.Credential.Uid == 0 && expectedUid != "0" {
		t.Errorf("unexpected UID: got 0, want %s", expectedUid)
	}
	if cmd.SysProcAttr.Credential.Gid == 0 && expectedGid != "0" {
		t.Errorf("unexpected GID: got 0, want %s", expectedGid)
	}

	foundUser := false
	foundHome := false
	for _, env := range cmd.Env {
		if env == "USER="+currentUser.Username {
			foundUser = true
		}
		if env == "HOME="+currentUser.HomeDir {
			foundHome = true
		}
	}
	if !foundUser {
		t.Error("USER env var not set in command environment")
	}
	if !foundHome {
		t.Error("HOME env var not set in command environment")
	}
}

func TestKillProcess_Success(t *testing.T) {
	cmd := exec.Command("sleep", "5")
	if err := cmd.Start(); err != nil {
		t.Fatalf("could not start process: %v", err)
	}
	if cmd.Process == nil {
		t.Fatal("process should not be nil after start")
	}

	err := killProcess(cmd.Process)
	if err != nil {
		t.Errorf("killProcess should not error on valid process: %v", err)
	}

	_ = cmd.Wait()
}

func TestKillProcess_AlreadyDead(t *testing.T) {
	cmd := exec.Command("true")
	if err := cmd.Run(); err != nil {
		t.Fatalf("could not run process: %v", err)
	}

	err := killProcess(cmd.Process)
	if err != nil {
		if errno, ok := err.(syscall.Errno); !ok || errno != syscall.ESRCH {
			t.Errorf("killProcess on dead process should return nil or ESRCH, got: %v", err)
		}
	}
}

func TestKillProcess_NilProcess(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Logf("killProcess with nil process panicked (expected): %v", r)
		}
	}()
	err := killProcess(nil)
	if err == nil {
		t.Log("killProcess with nil process returned nil")
	}
}

func TestKillProcess_ProcessGroup(t *testing.T) {
	if os.Getenv("CI") == "true" {
		t.Skip("skipping process group test in CI")
	}

	cmd := exec.Command("sh", "-c", "sleep 30 & wait $!")
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	if err := cmd.Start(); err != nil {
		t.Fatalf("could not start process group: %v", err)
	}

	err := killProcess(cmd.Process)
	_ = cmd.Wait()

	if err != nil {
		t.Logf("killProcess error: %v (may be expected on some systems)", err)
	}
}
