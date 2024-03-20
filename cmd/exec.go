package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/jclem/konk/konk/debugger"
	"github.com/jclem/konk/konk/konkfile"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var konkfilePath string

var execCommand = cobra.Command{
	Use:     "exec <command>",
	Aliases: []string{"e"},
	Short:   "Execute a command from a konkfile (alias: e)",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dbg := debugger.Get(cmd.Context())
		dbg.Flags(cmd)

		if workingDirectory != "" {
			if err := os.Chdir(workingDirectory); err != nil {
				return fmt.Errorf("changing working directory: %w", err)
			}
		}

		kfsearch := []string{"konkfile", "konkfile.json", "konkfile.toml", "konkfile.yaml", "konkfile.yml"}
		if konkfilePath != "" {
			kfsearch = []string{konkfilePath}
		}

		var kf []byte
		var kfpath string

		for _, kfp := range kfsearch {
			b, err := os.ReadFile(kfp)
			if err != nil {
				if os.IsNotExist(err) {
					continue
				}

				return fmt.Errorf("reading konkfile: %w", err)
			}

			kf = b
			kfpath = kfp
		}

		ext := filepath.Ext(kfpath)
		var file konkfile.File

		if ext == "" {
			if err := json.Unmarshal(kf, &file); err != nil {
				if err := yaml.Unmarshal(kf, &file); err != nil {
					if err := toml.Unmarshal(kf, &file); err != nil {
						return fmt.Errorf("unmarshalling konkfile: %w", err)
					}
				}
			}
		} else if ext == ".yaml" || ext == ".yml" {
			if err := yaml.Unmarshal(kf, &file); err != nil {
				return fmt.Errorf("unmarshalling konkfile: %w", err)
			}
		} else if ext == ".toml" {
			if err := toml.Unmarshal(kf, &file); err != nil {
				return fmt.Errorf("unmarshalling konkfile: %w", err)
			}
		} else {
			if err := json.Unmarshal(kf, &file); err != nil {
				return fmt.Errorf("unmarshalling konkfile: %w", err)
			}
		}

		if err := konkfile.Execute(cmd.Context(), file, args[0], konkfile.ExecuteConfig{
			AggregateOutput: aggregateOutput,
			ContinueOnError: continueOnError,
			NoColor:         noColor,
			NoShell:         noShell,
		}); err != nil {
			return fmt.Errorf("executing command: %w", err)
		}

		return nil
	},
}

func init() {
	execCommand.Flags().StringVarP(&workingDirectory, "working-directory", "w", "", "set the working directory for all commands")
	execCommand.Flags().BoolVarP(&aggregateOutput, "aggregate-output", "g", false, "aggregate command output")
	execCommand.Flags().BoolVarP(&continueOnError, "continue-on-error", "c", false, "continue running commands after a failure")
	execCommand.Flags().StringVarP(&konkfilePath, "konkfile", "k", "", "path to konkfile")
	rootCmd.AddCommand(&execCommand)
}
