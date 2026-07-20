//go:build !windows

package host

import (
	"os"
	"syscall"
)

func killProcess(process *os.Process) error {
	pid := process.Pid
	if err := syscall.Kill(-pid, syscall.SIGKILL); err != nil {
		errno, _ := err.(syscall.Errno)
		if errno != syscall.ESRCH && errno != syscall.EINVAL {
			return process.Kill()
		}
	}
	return nil
}
