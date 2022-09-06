package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var debug bool

var rootCmd = cobra.Command{
	Use:   "konk",
	Short: "Konk is a tool for running multiple processes",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Ensures that usage isn't printed for errors such as non-zero exits.
		// SEE: https://github.com/spf13/cobra/issues/340#issuecomment-378726225
		cmd.SilenceUsage = true
	},
}

func Execute() {
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "D", false, "debug mode")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
