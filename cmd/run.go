package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var cmdAsLabel bool
var npmCmds []string
var names []string

var runCommand = cobra.Command{
	Use:     "run <subcommand>",
	Aliases: []string{"r"},
	Short:   "Run commands serially or concurrently (alias: r)",
	RunE: func(cmd *cobra.Command, _ []string) error {
		_ = cmd.Help()
		os.Exit(1)
		return nil
	},
}

func init() {
	runCommand.PersistentFlags().StringVarP(&workingDirectory, "working-directory", "w", "", "set the working directory for all commands")
	runCommand.PersistentFlags().BoolVarP(&continueOnError, "continue-on-error", "c", false, "continue running commands after a failure")
	runCommand.PersistentFlags().BoolVarP(&noShell, "no-subshell", "S", false, "do not run commands in a subshell")
	runCommand.PersistentFlags().BoolVarP(&noColor, "no-color", "C", false, "do not colorize label output")

	runCommand.PersistentFlags().BoolVarP(&cmdAsLabel, "command-as-label", "L", false, "use each command as its own label")
	runCommand.PersistentFlags().StringArrayVarP(&npmCmds, "npm", "n", []string{}, "npm command")
	runCommand.PersistentFlags().StringArrayVarP(&names, "label", "l", []string{}, "label prefix for the command")
	runCommand.PersistentFlags().BoolVarP(&noLabel, "no-label", "B", false, "do not attach label/prefix to output")
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
				return nil, nil, fmt.Errorf("reading package.json: %w", err)
			}
			var pkg map[string]any
			if err := json.Unmarshal(pkgFile, &pkg); err != nil {
				return nil, nil, fmt.Errorf("unmarshalling package.json: %w", err)
			}

			// See if any "scripts" match our prefix
			matchingScripts := []string{}

			scripts, ok := pkg["scripts"].(map[string]any)
			if !ok {
				return nil, nil, errors.New("invalid scripts in package.json")
			}

			for script := range scripts {
				if strings.HasPrefix(script, prefix) {
					matchingScripts = append(matchingScripts, script)
				}
			}

			sort.Strings(matchingScripts)

			for _, script := range matchingScripts {
				commandStrings = append(commandStrings, script)
				commands = append(commands, "npm run "+script)
			}

			continue
		}

		script := "npm run " + cmd
		commandStrings = append(commandStrings, cmd)
		commands = append(commands, script)
	}

	return commandStrings, commands, nil
}

func collectLabels(commandStrings []string) []string {
	if noLabel {
		return make([]string, len(commandStrings))
	}

	labels := make([]string, len(commandStrings))

	for i, cmdStr := range commandStrings {
		switch {
		case cmdAsLabel:
			labels[i] = cmdStr
		case len(names) > 0:
			labels[i] = names[i]
		default:
			labels[i] = strconv.Itoa(i)
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
