package cmd

import (
	"fmt"

	"github.com/jclem/konk/konk"
	"github.com/spf13/cobra"
)

var snames []string

var SCommand = cobra.Command{
	Use:   "s <command...>",
	Short: "Run commands in serial",
	RunE: func(cmd *cobra.Command, args []string) error {
		commandStrings, commands, err := collectCommands(args)
		if err != nil {
			return err
		}

		if len(pnames) > 0 && len(pnames) != len(commands) {
			return fmt.Errorf("number of names must match number of commands")
		}

		labels := collectLabels(commandStrings)

		for i, cmd := range commands {
			c := konk.NewCommand(konk.CommandConfig{
				Name:  cmd[0],
				Args:  cmd[1:],
				Label: labels[i],
			})

			if err := c.Run(); err != nil {
				return err
			}
		}

		return nil
	},
}

func init() {
	SCommand.Flags().StringArrayVarP(&snames, "name", "n", []string{}, "name prefix for the command")
	rootCmd.AddCommand(&SCommand)
}
