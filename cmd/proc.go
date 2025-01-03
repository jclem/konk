package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/jclem/konk/konk"
	"github.com/spf13/cobra"
)

var envFile string
var noEnvFile bool
var procfile string
var omitEnv bool

var procCommand = cobra.Command{
	Use:     "proc",
	Aliases: []string{"p"},
	Short:   "Run commands defined in a Procfile (alias: p)",
	RunE: func(cmd *cobra.Command, _ []string) error {
		ctx := cmd.Context()

		if debug {
			cmd.DebugFlags()
		}

		if workingDirectory != "" {
			if err := os.Chdir(workingDirectory); err != nil {
				return fmt.Errorf("changing working directory: %w", err)
			}
		}

		procfile, err := os.Open(procfile)
		if err != nil {
			return fmt.Errorf("opening procfile: %w", err)
		}
		defer procfile.Close()

		procfileMap := map[string]string{}
		scanner := bufio.NewScanner(procfile)

		for scanner.Scan() {
			procfileLine := strings.TrimSpace(scanner.Text())

			if procfileLine == "" {
				continue
			}

			line := strings.SplitN(procfileLine, ":", 2)
			procfileMap[strings.TrimSpace(line[0])] = strings.TrimSpace(line[1])
		}

		envLines := []string{}
		if !noEnvFile {
			envFile, err := os.ReadFile(envFile)
			if err != nil {
				return fmt.Errorf("reading env file: %w", err)
			}
			envLines = strings.Split(string(envFile), "\n")
		}

		commandStrings := make([]string, 0, len(procfileMap))
		commandLabels := make([]string, 0, len(procfileMap))

		for label, command := range procfileMap {
			commandStrings = append(commandStrings, command)
			if noLabel {
				commandLabels = append(commandLabels, "")
			} else {
				commandLabels = append(commandLabels, label)
			}
		}

		if !noLabel {
			var maxLabelLen int

			for _, label := range commandLabels {
				if len(label) > maxLabelLen {
					maxLabelLen = len(label)
				}
			}

			for i, label := range commandLabels {
				commandLabels[i] = fmt.Sprintf("%s%s", label, strings.Repeat(" ", maxLabelLen-len(label)))
			}
		}

		commands, err := konk.RunConcurrently(ctx, konk.RunConcurrentlyConfig{
			Commands:        commandStrings,
			Labels:          commandLabels,
			Env:             envLines,
			OmitEnv:         omitEnv,
			AggregateOutput: false,
			ContinueOnError: continueOnError,
			NoColor:         noColor,
			NoShell:         noShell,
		})

		debugCommands(ctx, commands)

		if err != nil {
			return fmt.Errorf("running commands: %w", err)
		}

		return nil
	},
}

func init() {
	procCommand.Flags().StringVarP(&workingDirectory,
		"working-directory", "w", "", "set the working directory for all commands")
	procCommand.Flags().BoolVarP(&continueOnError,
		"continue-on-error", "c", false, "continue running commands after a failure")
	procCommand.Flags().BoolVarP(&noShell, "no-subshell", "S", false, "do not run commands in a subshell")
	procCommand.Flags().BoolVarP(&noColor, "no-color", "C", false, "do not colorize label output")

	procCommand.Flags().StringVarP(&procfile, "procfile", "p", "Procfile", "Path to the Procfile")
	procCommand.Flags().StringVarP(&envFile, "env-file", "e", ".env", "Path to the env file")
	procCommand.Flags().BoolVar(&omitEnv, "omit-env", false, "Omit any existing runtime environment variables")
	procCommand.Flags().BoolVarP(&noEnvFile, "no-env-file", "E", false, "Don't load the env file")
	procCommand.Flags().BoolVarP(&noLabel, "no-label", "B", false, "do not attach label/prefix to output")
	rootCmd.AddCommand(&procCommand)
}
