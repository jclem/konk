package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/mattn/go-shellwords"
	"github.com/spf13/cobra"
)

var cmdAsLabel bool
var npmCmds []string
var names []string
var continueOnError bool
var workingDirectory string

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
	runCommand.PersistentFlags().StringVarP(&workingDirectory, "working-directory", "w", "", "working directory")
	runCommand.PersistentFlags().BoolVarP(&continueOnError, "continue-on-error", "c", false, "continue running commands after a failure")
	runCommand.PersistentFlags().BoolVarP(&cmdAsLabel, "command-as-label", "L", false, "use command as label")
	runCommand.PersistentFlags().StringArrayVar(&npmCmds, "npm", []string{}, "npm command")
	runCommand.PersistentFlags().StringArrayVarP(&names, "name", "n", []string{}, "name prefix for the command")
	rootCmd.AddCommand(&runCommand)
}

func collectCommands(args []string) ([]string, [][]string, error) {
	commandStrings := []string{}
	commands := [][]string{}

	for _, cmd := range args {
		parts, err := shellwords.Parse(cmd)

		if err != nil {
			return nil, nil, err
		}

		commands = append(commands, parts)
		commandStrings = append(commandStrings, cmd)
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
				parts, err := shellwords.Parse(fmt.Sprintf("npm run %s", script))
				if err != nil {
					return nil, nil, err
				}
				commands = append(commands, parts)
				commandStrings = append(commandStrings, script)
			}

			continue
		}

		parts, err := shellwords.Parse(fmt.Sprintf("npm run %s", cmd))

		if err != nil {
			return nil, nil, err
		}

		commands = append(commands, parts)
		commandStrings = append(commandStrings, cmd)
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
