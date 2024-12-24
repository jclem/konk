package cmd

import (
	"context"
	"log/slog"
	"os"
	"strconv"

	"github.com/golang-cz/devslog"
	"github.com/jclem/konk/konk"
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

		level := slog.LevelInfo
		if debug {
			level = slog.LevelDebug
		}

		slog.SetDefault(slog.New(devslog.NewHandler(os.Stdout, &devslog.Options{ //nolint:exhaustruct // Fields not needed.
			HandlerOptions: &slog.HandlerOptions{Level: level}, //nolint:exhaustruct // Fields not needed.
		})))
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "D", false, "debug mode")
}

func debugCommands(ctx context.Context, commands []*konk.Command) {
	if commands != nil {
		var attrs []any
		for i, c := range commands {
			attrs = append(attrs, slog.Any(strconv.Itoa(i), c))
		}
		slog.DebugContext(ctx, "commands", attrs...)
	}
}
