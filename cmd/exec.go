package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/jclem/konk/konk/debugger"
	"github.com/jclem/konk/konk/konkfile"
	"github.com/spf13/cobra"
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

		kf, err := os.ReadFile(konkfilePath)
		if err != nil {
			return fmt.Errorf("reading konkfile: %w", err)
		}

		var file konkfile.File
		if err := json.Unmarshal(kf, &file); err != nil {
			return fmt.Errorf("unmarshalling konkfile: %w", err)
		}

		if err := konkfile.Execute(cmd.Context(), file, args[0], konkfile.ExecuteConfig{
			NoColor: noColor,
			NoShell: noShell,
		}); err != nil {
			return fmt.Errorf("executing command: %w", err)
		}

		return nil
	},
}

func init() {
	execCommand.Flags().StringVarP(&konkfilePath, "konkfile", "k", "konkfile.json", "path to konkfile")
	rootCmd.AddCommand(&execCommand)
}
