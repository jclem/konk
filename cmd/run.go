package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

var cmdAsLabel bool
var npmCmds []string
var names []string
var continueOnError bool
var noShell bool
var workingDirectory string

var runCommand = cobra.Command{
	Use:   "run <subcommand>",
	Short: "Run commands serially or concurrently",
	RunE: func(cmd *cobra.Command, args []string) error {
		_ = cmd.Help()
		os.Exit(1)
		return nil
	},
}

func init() {
	runCommand.PersistentFlags().StringVarP(&workingDirectory, "working-directory", "w", "", "set the working directory for all commands")
	runCommand.PersistentFlags().BoolVarP(&continueOnError, "continue-on-error", "c", false, "continue running commands after a failure")
	runCommand.PersistentFlags().BoolVarP(&cmdAsLabel, "command-as-label", "L", false, "use each command as its own label")
	runCommand.PersistentFlags().BoolVarP(&noShell, "no-subshell", "S", false, "do not run commands in a subshell")
	runCommand.PersistentFlags().StringArrayVarP(&npmCmds, "npm", "n", []string{}, "npm command")
	runCommand.PersistentFlags().StringArrayVarP(&names, "label", "l", []string{}, "label prefix for the command")
	rootCmd.AddCommand(&runCommand)
}

func collectCommands(args []string) ([]string, []string, error) {
	commandStrings := []string{}
	commands := []string{}

	for _, cmd := range args {
		commandStrings = append(commandStrings, cmd)
		commands = append(commands, cmd)
	}

	for _, cmd := range npmCmds {
		if strings.HasSuffix(cmd, "*") {
			prefix := strings.TrimSuffix(cmd, "*")
			pkgFile, err := os.ReadFile("package.json")
			if err != nil {
				return nil, nil, err
			}
			var pkg map[string]interface{}
			if err := json.Unmarshal(pkgFile, &pkg); err != nil {
				return nil, nil, err
			}

			// See if any "scripts" match our prefix
			matchingScripts := []string{}
			for script := range pkg["scripts"].(map[string]interface{}) {
				if strings.HasPrefix(script, prefix) {
					matchingScripts = append(matchingScripts, script)
				}
			}

			sort.Strings(matchingScripts)

			for _, script := range matchingScripts {
				commandStrings = append(commandStrings, script)
				commands = append(commands, fmt.Sprintf("npm run %s", script))
			}

			continue
		}

		script := fmt.Sprintf("npm run %s", cmd)
		commandStrings = append(commandStrings, cmd)
		commands = append(commands, script)
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
