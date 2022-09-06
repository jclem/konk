package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var debug bool

var rootCmd = cobra.Command{
	Use:   "konk",
	Short: "Konk is a tool for running multiple processes",
}

func Execute() {
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "x", false, "debug mode")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
