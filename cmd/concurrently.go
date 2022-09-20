package cmd

import (
	"fmt"
	"os"

	"github.com/jclem/konk/konk"
	"github.com/jclem/konk/konk/debugger"
	"github.com/spf13/cobra"
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

		commands, err := konk.RunConcurrently(cmd.Context(), konk.RunConcurrentlyConfig{
			Commands:        cmdParts,
			Labels:          labels,
			AggregateOutput: aggregateOutput,
			ContinueOnError: continueOnError,
			NoColor:         noColor,
			NoShell:         noShell,
		})

		if commands != nil {
			dbg.Prettyln(commands)
		}

		if commands != nil && aggregateOutput {
			for _, cmd := range commands {
				fmt.Print(cmd.ReadOut())
			}
		}

		return err
	},
}

func init() {
	cCommand.Flags().BoolVarP(&aggregateOutput, "aggregate-output", "g", false, "aggregate command output")
	runCommand.AddCommand(&cCommand)
}
