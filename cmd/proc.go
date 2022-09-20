package cmd

import (
	"os"
	"strings"

	"github.com/jclem/konk/konk"
	"github.com/jclem/konk/konk/debugger"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var envFile string
var noEnvFile bool
var procfile string

var procCommand = cobra.Command{
	Use:     "proc",
	Aliases: []string{"p"},
	Short:   "Run commands defined in a Procfile (alias: p)",
	RunE: func(cmd *cobra.Command, args []string) error {
		dbg := debugger.Get(cmd.Context())
		dbg.Flags(cmd)

		if workingDirectory != "" {
			if err := os.Chdir(workingDirectory); err != nil {
				return err
			}
		}

		procfile, err := os.ReadFile(procfile)
		if err != nil {
			return err
		}

		var procfileMap map[string]string
		if err := yaml.Unmarshal(procfile, &procfileMap); err != nil {
			return err
		}

		envLines := []string{}
		if !noEnvFile {
			envFile, err := os.ReadFile(envFile)
			if err != nil {
				return err
			}
			envLines = strings.Split(string(envFile), "\n")
		}

		commandStrings := make([]string, 0, len(procfileMap))
		commandLabels := make([]string, 0, len(procfileMap))

		for label, command := range procfileMap {
			commandStrings = append(commandStrings, command)
			commandLabels = append(commandLabels, label)
		}

		commands, err := konk.RunConcurrently(cmd.Context(), konk.RunConcurrentlyConfig{
			Commands:        commandStrings,
			Labels:          commandLabels,
			ContinueOnError: continueOnError,
			Env:             envLines,
			NoColor:         noColor,
			NoShell:         noShell,
		})

		if commands != nil {
			dbg.Prettyln(commands)
		}

		return err
	},
}

func init() {
	procCommand.Flags().StringVarP(&workingDirectory, "working-directory", "w", "", "set the working directory for all commands")
	procCommand.Flags().BoolVarP(&continueOnError, "continue-on-error", "c", false, "continue running commands after a failure")
	procCommand.Flags().BoolVarP(&noShell, "no-subshell", "S", false, "do not run commands in a subshell")
	procCommand.Flags().BoolVarP(&noColor, "no-color", "C", false, "do not colorize label output")

	procCommand.Flags().StringVarP(&procfile, "procfile", "p", "Procfile", "Path to the Procfile")
	procCommand.Flags().StringVarP(&envFile, "env-file", "e", ".env", "Path to the env file")
	procCommand.Flags().BoolVarP(&noEnvFile, "no-env-file", "E", false, "Don't load the env file")
	rootCmd.AddCommand(&procCommand)
}
