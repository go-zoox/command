package main

import (
	"github.com/go-zoox/cli"
	"github.com/go-zoox/command"
	"github.com/go-zoox/command/cmd/cmd/commands"
)

func main() {
	app := cli.NewMultipleProgram(&cli.MultipleProgramConfig{
		Name:    "cmd",
		Usage:   "Powerful command runner",
		Version: command.Version,
	})

	commands.Exec(app)

	app.Run()
}
