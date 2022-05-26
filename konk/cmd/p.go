package cmd

import (
	"fmt"

	"github.com/jclem/konk/konk"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

var pnames []string

var PCommand = cobra.Command{
	Use:   "p <command...>",
	Short: "Run commands in parallel",
	RunE: func(cmd *cobra.Command, args []string) error {
		commandStrings, commands, err := collectCommands(args)
		if err != nil {
			return err
		}

		if len(pnames) > 0 && len(pnames) != len(commands) {
			return fmt.Errorf("number of names must match number of commands")
		}

		labels := collectLabels(commandStrings)

		var eg errgroup.Group

		for i, cmd := range commands {
			cmd := cmd
			i := i

			eg.Go(func() error {
				c := konk.NewCommand(konk.CommandConfig{
					Name:  cmd[0],
					Args:  cmd[1:],
					Label: labels[i],
				})

				return c.Run()
			})
		}

		return eg.Wait()
	},
}

func init() {
	PCommand.Flags().StringArrayVarP(&pnames, "name", "n", []string{}, "name prefix for the command")
	rootCmd.AddCommand(&PCommand)
}
