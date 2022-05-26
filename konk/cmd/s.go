package cmd

import (
	"fmt"

	"github.com/jclem/konk/konk"
	"github.com/mattn/go-shellwords"
	"github.com/spf13/cobra"
)

var snames []string

var SCommand = cobra.Command{
	Use:   "s <command...>",
	Short: "Run commands in serial",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		commands := make([][]string, len(args))
		labels := make([]string, len(args))

		if len(snames) > 0 && len(snames) != len(args) {
			return fmt.Errorf("number of names must match number of commands")
		}

		for i, cmd := range args {
			args, err := shellwords.Parse(cmd)

			if err != nil {
				return err
			}

			commands[i] = args

			if cmdAsLabel {
				labels[i] = cmd
			} else if len(snames) > 0 {
				labels[i] = snames[i]
			} else {
				labels[i] = fmt.Sprintf("%d", i)
			}
		}

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
