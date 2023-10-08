package main

import (
	"github.com/go-zoox/cli"
	"github.com/go-zoox/command"
	"github.com/go-zoox/command/cmd/command-runner/commands"
)

func main() {
	app := cli.NewMultipleProgram(&cli.MultipleProgramConfig{
		Name:    "command-runner",
		Usage:   "Powerful command runner",
		Version: command.Version,
	})

	commands.Exec(app)

	app.Run()
}
