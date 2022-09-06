package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/jclem/konk/konk"
	"github.com/mattn/go-shellwords"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

var aggregateOutput bool

var pCommand = cobra.Command{
	Use:   "p <command...>",
	Short: "Run commands in parallel",
	RunE: func(cmd *cobra.Command, args []string) error {
		if workingDirectory != "" {
			if err := os.Chdir(workingDirectory); err != nil {
				return err
			}
		}

		cmdStrings, cmdParts, err := collectCommands(args)
		if err != nil {
			return err
		}

		if len(names) > 0 && len(names) != len(cmdParts) {
			return fmt.Errorf("number of names must match number of commands")
		}

		labels := collectLabels(cmdStrings)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		eg, ctx := errgroup.WithContext(ctx)

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

		for _, cmd := range commands {
			cmd := cmd

			eg.Go(func() error {
				return cmd.Run(ctx, cancel, konk.RunCommandConfig{
					AggregateOutput: aggregateOutput,
					KillOnCancel:    !continueOnError,
				})
			})
		}

		waitErr := eg.Wait()

		if aggregateOutput {
			for _, cmd := range commands {
				fmt.Print(cmd.ReadOut())
			}
		}

		if waitErr != nil {
			return waitErr
		}

		return nil
	},
}

func init() {
	pCommand.Flags().BoolVarP(&aggregateOutput, "aggregate-output", "g", false, "aggregate command output")
	runCommand.AddCommand(&pCommand)
}
