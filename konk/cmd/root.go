package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cmdAsLabel bool

var rootCmd = cobra.Command{
	Use:   "konk",
	Short: "Konk is a tool for running multiple processes",
}

func Execute() {
	rootCmd.PersistentFlags().BoolVar(&cmdAsLabel, "command-as-label", false, "use command as label")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
