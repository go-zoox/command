package main

import (
	"github.com/go-zoox/cli"
	"github.com/go-zoox/command"
)

func main() {
	app := cli.NewMultipleProgram(&cli.MultipleProgramConfig{
		Name:    "command-runner",
		Usage:   "Powerful command runner",
		Version: command.Version,
	})

	registerServerCommand(app)
	registerClientCommand(app)

	app.Run()
}
