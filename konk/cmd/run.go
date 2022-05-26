package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/mattn/go-shellwords"
	"github.com/spf13/cobra"
)

var cmdAsLabel bool
var npmCmd []string
var names []string

var runCommand = cobra.Command{
	Use:   "run <subcommand>",
	Short: "Run commands in serial or parallel",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Help()
		os.Exit(1)
		return nil
	},
}

func init() {
	runCommand.PersistentFlags().BoolVar(&cmdAsLabel, "command-as-label", false, "use command as label")
	runCommand.PersistentFlags().StringArrayVar(&npmCmd, "npm", []string{}, "npm command")
	runCommand.PersistentFlags().StringArrayVarP(&names, "name", "n", []string{}, "name prefix for the command")
	rootCmd.AddCommand(&runCommand)
}

func collectCommands(args []string) ([]string, [][]string, error) {
	commandStrings := make([]string, len(args)+len(npmCmd))
	commands := make([][]string, len(args)+len(npmCmd))

	for i, cmd := range args {
		parts, err := shellwords.Parse(cmd)

		if err != nil {
			return nil, nil, err
		}

		commandStrings[i] = cmd
		commands[i] = parts
	}

	for i, cmd := range npmCmd {
		parts, err := shellwords.Parse(fmt.Sprintf("npm run %s", cmd))

		if err != nil {
			return nil, nil, err
		}

		commandStrings[i+len(args)] = cmd
		commands[i+len(args)] = parts
	}

	return commandStrings, commands, nil
}

func collectLabels(commandStrings []string) []string {
	labels := make([]string, len(commandStrings))

	for i, cmdStr := range commandStrings {
		if cmdAsLabel {
			labels[i] = cmdStr
		} else if len(names) > 0 {
			labels[i] = names[i]
		} else {
			labels[i] = fmt.Sprintf("%d", i)
		}
	}

	var maxLabelLen int

	for _, label := range labels {
		if len(label) > maxLabelLen {
			maxLabelLen = len(label)
		}
	}

	for i, label := range labels {
		labels[i] = fmt.Sprintf("%s%s", label, strings.Repeat(" ", maxLabelLen-len(label)))
	}

	return labels
}
