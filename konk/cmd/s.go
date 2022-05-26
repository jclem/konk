package cmd

import (
	"bufio"
	"fmt"
	"os/exec"

	"github.com/mattn/go-shellwords"
	"github.com/spf13/cobra"
)

var names []string

var SCommand = cobra.Command{
	Use:   "s <command...>",
	Short: "Run commands in serial",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		commands := make([][]string, len(args))
		labels := make([]string, len(args))

		if len(names) > 0 && len(names) != len(args) {
			return fmt.Errorf("number of names must match number of commands")
		}

		for i, cmd := range args {
			args, err := shellwords.Parse(cmd)

			if err != nil {
				return err
			}

			commands[i] = args

			if len(names) > 0 {
				labels[i] = names[i]
			} else {
				labels[i] = fmt.Sprintf("%d", i)
			}
		}

		for i, cmd := range commands {
			c := exec.Command(cmd[0], cmd[1:]...)
			out := make(chan string)

			stdout, err := c.StdoutPipe()
			if err != nil {
				return err
			}
			c.Stderr = c.Stdout

			scanner := bufio.NewScanner(stdout)
			done := make(chan bool)

			if err := c.Start(); err != nil {
				return err
			}

			go func() {
				for scanner.Scan() {
					out <- scanner.Text()
				}
				done <- true
			}()

			go func() {
				for {
					t := <-out
					fmt.Println(fmt.Sprintf("[%s] %s", labels[i], t))
				}
			}()

			if err := c.Wait(); err != nil {
				return err
			}
		}
		return nil
	},
}

func init() {
	SCommand.Flags().StringArrayVarP(&names, "name", "n", []string{}, "name prefix for the command")
	rootCmd.AddCommand(&SCommand)
}
