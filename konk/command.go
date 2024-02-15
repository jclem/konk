package konk

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"

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
	Env     []string
	OmitEnv bool
}

func NewShellCommand(conf ShellCommandConfig) *Command {
	c := exec.Command("/bin/sh", "-c", conf.Command) //nolint:gosec // Intentional user-defined sub-process.
	setEnv(c, conf.Env, conf.OmitEnv)
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
	Env     []string
	OmitEnv bool
}

func setEnv(c *exec.Cmd, env []string, omitEnv bool) {
	if !omitEnv {
		c.Env = os.Environ()
	}

	c.Env = append(c.Env, env...)
}

func NewCommand(conf CommandConfig) *Command {
	c := exec.Command(conf.Name, conf.Args...) //nolint:gosec // Intentional user-defined sub-process.
	setEnv(c, conf.Env, conf.OmitEnv)
	prefix := getPrefix(conf.Label, conf.NoColor)

	return &Command{
		c:      c,
		prefix: prefix,
	}
}

func (c *Command) Run(ctx context.Context, cancel context.CancelFunc, conf RunCommandConfig) error {
	stdout, err := c.c.StdoutPipe()
	if err != nil {
		return fmt.Errorf("getting stdout pipe: %w", err)
	}
	c.c.Stderr = c.c.Stdout

	out := make(chan string)
	scanner := bufio.NewScanner(stdout)
	scannerDone := make(chan bool)
	scannerErr := make(chan error)
	allDone := make(chan error)

	if err := c.c.Start(); err != nil {
		return fmt.Errorf("starting command: %w", err)
	}

	// Start a goroutine to read the command's output. Send that output to the
	// `out` channel and notify `scannerDone` when complete.
	go func() {
		for scanner.Scan() {
			out <- scanner.Text()
		}

		if err := scanner.Err(); err != nil {
			scannerErr <- err
		}

		scannerDone <- true
	}()

	// Read from the `out` channel and print or aggregate output. Send a signal to
	// `allDone` if our context is canceled, there is a scanner error, or if the
	// scanner is done.
	go func() {
		for {
			select {
			case t := <-out:
				line := fmt.Sprintf("%s %s\n", c.prefix, t)

				if conf.AggregateOutput {
					c.out.WriteString(line)
				} else {
					fmt.Fprint(os.Stdout, line)
				}
			case <-ctx.Done():
				if conf.KillOnCancel {
					_ = c.c.Process.Signal(syscall.SIGTERM)
					allDone <- nil
					return
				}
			case err := <-scannerErr:
				allDone <- err
			case <-scannerDone:
				allDone <- nil
			}
		}
	}()

	// We wait for `allDone` to ensure we have fully read the output (or
	// encountedred an error) *before* we call `Wait()` below.
	// SEE: https://pkg.go.dev/os/exec@go1.19.1#Cmd.StdoutPipe
	// "Wait will close the pipe after seeing the command exit, so most callers
	// need not close the pipe themselves. It is thus incorrect to call Wait
	// before all reads from the pipe have completed."
	if err := <-allDone; err != nil {
		return err
	}

	if conf.AggregateOutput {
		fmt.Fprint(os.Stdout, c.ReadOut())
	}

	if err := c.c.Wait(); err != nil {
		cancel()

		var xerr *exec.ExitError
		if errors.As(err, &xerr) {
			exitErr := newExitError(c.prefix, xerr)
			fmt.Fprintln(os.Stdout, exitErr)
			return exitErr
		}

		return fmt.Errorf("waiting for command: %w", err)
	}

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

func getPrefix(label string, noColor bool) string {
	var prefix string

	if noColor {
		prefix = fmt.Sprintf("[%s]", label)
	} else {
		prefixColor := rand.Intn(16) + 1 //nolint:gosec // No need for cryptographic randomness for process labels.
		prefixStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(strconv.Itoa(prefixColor)))
		prefix = prefixStyle.Render(fmt.Sprintf("[%s]", label))
	}

	return prefix
}
