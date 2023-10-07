package commands

import (
	"github.com/go-zoox/cli"
	"github.com/go-zoox/command"
)

func Run(app *cli.MultipleProgram) {
	app.Register("run", &cli.Command{
		Name:  "run",
		Usage: "command run",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "engine",
				Usage:   "command engine",
				Aliases: []string{"e"},
				EnvVars: []string{"ENGINE"},
				Value:   "host",
			},
			&cli.StringFlag{
				Name:     "command",
				Usage:    "the command",
				Aliases:  []string{"c"},
				EnvVars:  []string{"COMMAND"},
				Required: true,
			},
			&cli.StringFlag{
				Name:    "shell",
				Usage:   "the command shell",
				Aliases: []string{"s"},
				EnvVars: []string{"SHELL"},
			},
		},
		Action: func(ctx *cli.Context) (err error) {
			cmd, err := command.New(&command.Config{
				Engine:  ctx.String("engine"),
				Command: ctx.String("command"),
				Shell:   ctx.String("shell"),
			})
			if err != nil {
				return err
			}

			return cmd.Run()
		},
	})
}
