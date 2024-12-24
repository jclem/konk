package cmd

import (
	"context"
	"os"

	"github.com/jclem/konk/konk/debugger"
	"github.com/spf13/cobra"
)

var debug bool
var continueOnError bool
var noShell bool
var workingDirectory string
var noColor bool
var noLabel bool

var Version = "dev"

var rootCmd = &cobra.Command{
	Use:               "konk",
	Short:             "Konk is a tool for running multiple processes",
	Version:           Version,
	DisableAutoGenTag: true,
	PersistentPreRun: func(cmd *cobra.Command, _ []string) {
		// Ensures that usage isn't printed for errors such as non-zero exits.
		// SEE: https://github.com/spf13/cobra/issues/340#issuecomment-378726225
		cmd.SilenceUsage = true
	},
}

func Execute() {
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "D", false, "debug mode")
	ctx := debugger.WithDebugger(context.Background(), &debug)

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}
}
