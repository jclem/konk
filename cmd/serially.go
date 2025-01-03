package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/jclem/konk/konk"
	"github.com/mattn/go-shellwords"
	"github.com/spf13/cobra"
)

var sCommand = cobra.Command{
	Use:     "serially <command...>",
	Aliases: []string{"s"},
	Short:   "Run commands serially (alias: s)",
	Example: `# Run two commands in serial

konk run serially "echo foo" "echo bar"

# Run a set of npm commands in serial

konk run serially -n build -n deploy`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		if debug {
			cmd.DebugFlags()
		}

		if workingDirectory != "" {
			if err := os.Chdir(workingDirectory); err != nil {
				return fmt.Errorf("changing working directory: %w", err)
			}
		}

		commandStrings, cmdParts, err := collectCommands(args)
		if err != nil {
			return err
		}

		if len(names) > 0 && len(names) != len(cmdParts) {
			return errors.New("number of names must match number of commands")
		}

		labels := collectLabels(commandStrings)

		var errCmd error

		commands := make([]*konk.Command, len(cmdParts))

		for i, cmd := range cmdParts {
			var c *konk.Command

			if noShell {
				parts, err := shellwords.Parse(cmd)

				if err != nil {
					return fmt.Errorf("parsing command: %w", err)
				}

				c = konk.NewCommand(konk.CommandConfig{
					Name:    parts[0],
					Args:    parts[1:],
					Label:   labels[i],
					NoColor: noColor,
					Env:     make([]string, 0),
					OmitEnv: false,
				})
			} else {
				c = konk.NewShellCommand(konk.ShellCommandConfig{
					Command: cmd,
					Label:   labels[i],
					NoColor: noColor,
					Env:     make([]string, 0),
					OmitEnv: false,
				})
			}

			commands[i] = c
		}

		debugCommands(ctx, commands)

		for _, c := range commands {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			if err := c.Run(ctx, cancel, konk.RunCommandConfig{
				AggregateOutput: false,
				StopOnCancel:    false,
			}); err != nil && continueOnError {
				errCmd = err
			} else if err != nil {
				return fmt.Errorf("running command: %w", err)
			}
		}

		return errCmd
	},
}

func init() {
	runCommand.AddCommand(&sCommand)
}
