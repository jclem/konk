package cmd

import (
	"fmt"

	"github.com/jclem/konk/konk"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

var aggregateOutput bool

var pCommand = cobra.Command{
	Use:   "p <command...>",
	Short: "Run commands in parallel",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmdStrings, cmdParts, err := collectCommands(args)
		if err != nil {
			return err
		}

		if len(names) > 0 && len(names) != len(cmdParts) {
			return fmt.Errorf("number of names must match number of commands")
		}

		labels := collectLabels(cmdStrings)

		var eg errgroup.Group

		commands := make([]*konk.Command, len(cmdParts))

		for i, cmd := range cmdParts {
			c := konk.NewCommand(konk.CommandConfig{
				Name:  cmd[0],
				Args:  cmd[1:],
				Label: labels[i],
			})

			commands[i] = c
		}

		for _, cmd := range commands {
			cmd := cmd

			eg.Go(func() error {
				return cmd.Run(konk.RunCommandConfig{
					AggregateOutput: aggregateOutput,
				})
			})
		}

		if err := eg.Wait(); err != nil {
			return err
		}

		if aggregateOutput {
			for _, cmd := range commands {
				fmt.Print(cmd.ReadOut())
			}
		}

		return nil
	},
}

func init() {
	pCommand.Flags().BoolVarP(&aggregateOutput, "aggregate-output", "g", false, "aggregate output")
	runCommand.AddCommand(&pCommand)
}
