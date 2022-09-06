package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/jclem/konk/konk"
	"github.com/jclem/konk/konk/debugger"
	"github.com/mattn/go-shellwords"
	"github.com/spf13/cobra"
)

var sCommand = cobra.Command{
	Use:     "serially <command...>",
	Aliases: []string{"s"},
	Short:   "Run commands serially",
	RunE: func(cmd *cobra.Command, args []string) error {
		dbg := debugger.Get(cmd.Context())
		dbg.Flags(cmd)

		if workingDirectory != "" {
			if err := os.Chdir(workingDirectory); err != nil {
				return err
			}
		}

		commandStrings, cmdParts, err := collectCommands(args)
		if err != nil {
			return err
		}

		if len(names) > 0 && len(names) != len(cmdParts) {
			return fmt.Errorf("number of names must match number of commands")
		}

		labels := collectLabels(commandStrings)

		var cmdErr error

		commands := make([]*konk.Command, len(cmdParts))

		for i, cmd := range cmdParts {
			var c *konk.Command

			if noShell {
				parts, err := shellwords.Parse(cmd)

				if err != nil {
					return err
				}

				c = konk.NewCommand(konk.CommandConfig{
					Name:  parts[0],
					Args:  parts[1:],
					Label: labels[i],
				})
			} else {
				c = konk.NewShellCommand(konk.ShellCommandConfig{
					Command: cmd,
					Label:   labels[i],
				})
			}

			commands[i] = c
		}

		dbg.Prettyln(commands)

		for _, c := range commands {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			err := c.Run(ctx, cancel, konk.RunCommandConfig{})

			if err != nil && !continueOnError {
				return err
			}

			if err != nil {
				cmdErr = err
			}
		}

		return cmdErr
	},
}

func init() {
	runCommand.AddCommand(&sCommand)
}
