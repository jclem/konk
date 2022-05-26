package cmd

import (
	"fmt"
	"strings"

	"github.com/jclem/konk/konk"
	"github.com/mattn/go-shellwords"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

var pnames []string

var PCommand = cobra.Command{
	Use:   "p <command...>",
	Short: "Run commands in parallel",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		commands := make([][]string, len(args))
		labels := make([]string, len(args))

		if len(pnames) > 0 && len(pnames) != len(args) {
			return fmt.Errorf("number of names must match number of commands")
		}

		for i, cmd := range args {
			args, err := shellwords.Parse(cmd)

			if err != nil {
				return err
			}

			commands[i] = args

			if len(pnames) > 0 {
				labels[i] = pnames[i]
			} else {
				labels[i] = fmt.Sprintf("%d", i)
			}
		}

		var maxLabelLen int

		for _, label := range labels {
			if len(label) > maxLabelLen {
				maxLabelLen = len(label)
			}
		}

		for i, label := range labels {
			labels[i] = fmt.Sprintf("%s%s", label, strings.Repeat(" ", maxLabelLen-len(label)))
		}

		var eg errgroup.Group

		for i, cmd := range commands {
			cmd := cmd
			i := i

			eg.Go(func() error {
				c := konk.NewCommand(konk.CommandConfig{
					Name:  cmd[0],
					Args:  cmd[1:],
					Label: labels[i],
				})

				return c.Run()
			})
		}

		return eg.Wait()
	},
}

func init() {
	PCommand.Flags().StringArrayVarP(&pnames, "name", "n", []string{}, "name prefix for the command")
	rootCmd.AddCommand(&PCommand)
}
