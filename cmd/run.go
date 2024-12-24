package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var cmdAsLabel bool
var npmCmds []string
var runWithBun bool
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
	runCommand.PersistentFlags().StringVarP(&workingDirectory,
		"working-directory", "w", "", "set the working directory for all commands")
	runCommand.PersistentFlags().BoolVarP(&continueOnError,
		"continue-on-error", "c", false, "continue running commands after a failure")
	runCommand.PersistentFlags().BoolVarP(&noShell, "no-subshell", "S", false, "do not run commands in a subshell")
	runCommand.PersistentFlags().BoolVarP(&noColor, "no-color", "C", false, "do not colorize label output")

	runCommand.PersistentFlags().BoolVarP(&cmdAsLabel, "command-as-label", "L", false, "use each command as its own label")
	runCommand.PersistentFlags().StringArrayVarP(&npmCmds, "npm", "n", []string{}, "npm command")
	runCommand.PersistentFlags().BoolVarP(&runWithBun, "bun", "b", false, "Run npm commands with Bun")
	runCommand.PersistentFlags().StringArrayVarP(&names, "label", "l", []string{}, "label prefix for the command")
	runCommand.PersistentFlags().BoolVarP(&noLabel, "no-label", "B", false, "do not attach label/prefix to output")
	rootCmd.AddCommand(&runCommand)
}

func collectCommands(args []string) ([]string, []string, error) {
	// The commands as provided by the user
	providedCommands := []string{}

	// "Resolved" commands (e.g. prefixed with "bun run", etc.)
	runnableCommands := []string{}

	for _, cmd := range args {
		providedCommands = append(providedCommands, cmd)
		runnableCommands = append(runnableCommands, cmd)
	}

	var scripts []string

	if len(npmCmds) > 0 {
		var err error
		scripts, err = getPackageJSONScripts()
		if err != nil {
			return nil, nil, err
		}
	}

	for _, cmd := range npmCmds {
		if strings.HasSuffix(cmd, "*") {
			prefix := strings.TrimSuffix(cmd, "*")

			// See if any "scripts" match our prefix
			matchingScripts := []string{}

			for _, script := range scripts {
				if strings.HasPrefix(script, prefix) {
					matchingScripts = append(matchingScripts, script)
				}
			}

			for _, script := range matchingScripts {
				providedCommands = append(providedCommands, script)

				if runWithBun {
					runnableCommands = append(runnableCommands, "bun run "+script)
				} else {
					runnableCommands = append(runnableCommands, "npm run "+script)
				}
			}

			continue
		}

		providedCommands = append(providedCommands, cmd)

		if runWithBun {
			runnableCommands = append(runnableCommands, "bun run "+cmd)
		} else {
			runnableCommands = append(runnableCommands, "npm run "+cmd)
		}
	}

	return providedCommands, runnableCommands, nil
}

func getPackageJSONScripts() ([]string, error) {
	pkgFile, err := os.ReadFile("package.json")
	if err != nil {
		return nil, fmt.Errorf("reading package.json: %w", err)
	}

	var pkgJSON struct {
		Scripts map[string]string `json:"scripts"`
	}

	if err := json.Unmarshal(pkgFile, &pkgJSON); err != nil {
		return nil, fmt.Errorf("unmarshalling package.json: %w", err)
	}

	scripts := make([]string, 0, len(pkgJSON.Scripts))
	for script := range pkgJSON.Scripts {
		scripts = append(scripts, script)
	}

	sort.Strings(scripts)

	return scripts, nil
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
