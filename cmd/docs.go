package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var docsFormat string

var docsCmd = &cobra.Command{
	Use:   "docs",
	Short: "Print documentation",
	RunE: func(cmd *cobra.Command, _ []string) error {
		switch docsFormat {
		case "markdown":
			md, err := genMarkdown(rootCmd)
			if err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "%s", md)
		}

		return nil
	},
}

func genMarkdown(cmd *cobra.Command) (string, error) {
	b := strings.Builder{}

	linkHandler := func(s string) string {
		s = strings.ReplaceAll(s, "_", "-")
		s = strings.TrimSuffix(s, ".md")
		s = "#" + s
		return s
	}

	var gen func(cmd *cobra.Command) error
	gen = func(cmd *cobra.Command) error {
		if err := doc.GenMarkdownCustom(cmd, &b, linkHandler); err != nil {
			return fmt.Errorf("generate markdown: %w", err)
		}

		for _, c := range cmd.Commands() {
			if !c.IsAvailableCommand() || c.IsAdditionalHelpTopicCommand() {
				continue
			}

			if err := gen(c); err != nil {
				return err
			}
		}

		return nil
	}

	if err := gen(cmd); err != nil {
		return "", err
	}

	return b.String(), nil
}

func init() {
	docsCmd.Flags().StringVarP(&docsFormat, "format", "f", "markdown", "output format")
	rootCmd.AddCommand(docsCmd)
}
