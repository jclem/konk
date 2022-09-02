package cmd

import (
	"context"
	"fmt"

	"github.com/jclem/konk/konk"
	"github.com/spf13/cobra"
)

var sCommand = cobra.Command{
	Use:   "s <command...>",
	Short: "Run commands in serial",
	RunE: func(cmd *cobra.Command, args []string) error {
		commandStrings, commands, err := collectCommands(args)
		if err != nil {
			return err
		}

		if len(names) > 0 && len(names) != len(commands) {
			return fmt.Errorf("number of names must match number of commands")
		}

		labels := collectLabels(commandStrings)

		for i, cmd := range commands {
			c := konk.NewCommand(konk.CommandConfig{
				Name:  cmd[0],
				Args:  cmd[1:],
				Label: labels[i],
			})

			if err := c.Run(context.Background(), konk.RunCommandConfig{}); err != nil && !continueOnError {
				return err
			}
		}

		return nil
	},
}

func init() {
	runCommand.AddCommand(&sCommand)
}
