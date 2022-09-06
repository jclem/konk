package konk

import (
	"bufio"
	"context"
	"fmt"
	"math/rand"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/charmbracelet/lipgloss"
)

type Command struct {
	c      *exec.Cmd
	out    strings.Builder
	prefix string
}

type RunCommandConfig struct {
	AggregateOutput bool
	KillOnCancel    bool
}

type ShellCommandConfig struct {
	Command string
	Label   string
	NoColor bool
}

func NewShellCommand(conf ShellCommandConfig) *Command {
	c := exec.Command("/bin/sh", "-c", conf.Command)
	prefix := getPrefix(conf.Label, conf.NoColor)

	return &Command{
		c:      c,
		prefix: prefix,
	}
}

type CommandConfig struct {
	Name    string
	Args    []string
	Label   string
	NoColor bool
}

func NewCommand(conf CommandConfig) *Command {
	c := exec.Command(conf.Name, conf.Args...)
	prefix := getPrefix(conf.Label, conf.NoColor)

	return &Command{
		c:      c,
		prefix: prefix,
	}
}

func (c *Command) Run(ctx context.Context, cancel context.CancelFunc, conf RunCommandConfig) error {
	stdout, err := c.c.StdoutPipe()
	if err != nil {
		return err
	}
	c.c.Stderr = c.c.Stdout

	out := make(chan string)
	done := make(chan bool)
	scanner := bufio.NewScanner(stdout)

	if err := c.c.Start(); err != nil {
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
			select {
			case <-ctx.Done():
				if conf.KillOnCancel {
					_ = c.c.Process.Signal(syscall.SIGTERM)
					return
				}
			case t := <-out:
				line := fmt.Sprintf("%s %s\n", c.prefix, t)

				if conf.AggregateOutput {
					c.out.WriteString(line)
				} else {
					fmt.Print(line)
				}
			}
		}
	}()

	if err := c.c.Wait(); err != nil {
		cancel()

		if execExitErr, ok := err.(*exec.ExitError); ok {
			exitErr := newExitError(c.prefix, execExitErr)
			fmt.Println(exitErr)
			return exitErr
		}

		return err
	}

	// Flush remainder of scanner
	<-done

	return nil
}

func (c *Command) ReadOut() string {
	return c.out.String()
}

type ExitError struct {
	label string
	err   error
}

func (e *ExitError) Error() string {
	return fmt.Sprintf("%s exited with error: %s", e.label, e.err)
}

func newExitError(label string, err error) error {
	return &ExitError{
		label: label,
		err:   err,
	}
}

func init() {
	// Seed random for random prefix colors.
	rand.Seed(time.Now().UnixNano())
}

func getPrefix(label string, noColor bool) string {

	var prefix string

	if noColor {
		prefix = fmt.Sprintf("[%s]", label)
	} else {
		prefixColor := rand.Intn(16) + 1
		prefixStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(fmt.Sprint(prefixColor)))
		prefix = prefixStyle.Render(fmt.Sprintf("[%s]", label))
	}

	return prefix
}
