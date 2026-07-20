//go:build windows

package commands

import "github.com/go-zoox/command/terminal"

func startResizeWatcher(t terminal.Terminal) {
	// SIGWINCH is not available on Windows
}
