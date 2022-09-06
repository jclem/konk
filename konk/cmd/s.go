package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/jclem/konk/konk"
	"github.com/mattn/go-shellwords"
	"github.com/spf13/cobra"
)

var sCommand = cobra.Command{
	Use:   "s <command...>",
	Short: "Run commands in serial",
	RunE: func(cmd *cobra.Command, args []string) error {
		if workingDirectory != "" {
			if err := os.Chdir(workingDirectory); err != nil {
				return err
			}
		}

		commandStrings, commands, err := collectCommands(args)
		if err != nil {
			return err
		}

		if len(names) > 0 && len(names) != len(commands) {
			return fmt.Errorf("number of names must match number of commands")
		}

		labels := collectLabels(commandStrings)

		var cmdErr error

		for i, cmd := range commands {
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

			err := c.Run(context.Background(), konk.RunCommandConfig{})

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
