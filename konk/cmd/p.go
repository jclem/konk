package cmd

import (
	"context"
	"fmt"

	"github.com/jclem/konk/konk"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

var aggregateOutput bool
var continueOnError bool

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

		eg, ctx := errgroup.WithContext(context.Background())

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
				return cmd.Run(ctx, konk.RunCommandConfig{
					AggregateOutput: aggregateOutput,
					KillOnCancel:    !continueOnError,
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
	pCommand.Flags().BoolVarP(&continueOnError, "continue-on-error", "c", false, "continue running commands after a failure")
	pCommand.Flags().BoolVarP(&aggregateOutput, "aggregate-output", "g", false, "aggregate output")
	runCommand.AddCommand(&pCommand)
}
