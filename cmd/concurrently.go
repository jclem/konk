package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/jclem/konk/konk"
	"github.com/jclem/konk/konk/debugger"
	"github.com/mattn/go-shellwords"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

var aggregateOutput bool

var cCommand = cobra.Command{
	Use:     "concurrently <command...>",
	Aliases: []string{"c"},
	Short:   "Run commands concurrently (alias: c)",
	RunE: func(cmd *cobra.Command, args []string) error {
		dbg := debugger.Get(cmd.Context())
		dbg.Flags(cmd)

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

		ctx, cancel := context.WithCancel(cmd.Context())
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
					Name:    parts[0],
					Args:    parts[1:],
					Label:   labels[i],
					NoColor: noColor,
				})
			} else {
				c = konk.NewShellCommand(konk.ShellCommandConfig{
					Command: cmd,
					Label:   labels[i],
					NoColor: noColor,
				})
			}

			commands[i] = c
		}

		dbg.Prettyln(commands)

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
	cCommand.Flags().BoolVarP(&aggregateOutput, "aggregate-output", "g", false, "aggregate command output")
	runCommand.AddCommand(&cCommand)
}