//go:build windows

package host

import "os"

func killProcess(process *os.Process) error {
	return process.Kill()
}
