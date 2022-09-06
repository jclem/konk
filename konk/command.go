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

type CommandConfig struct {
	Command string
	Label   string
}

type RunCommandConfig struct {
	AggregateOutput bool
	KillOnCancel    bool
}

func NewCommand(conf CommandConfig) *Command {
	c := exec.Command("/bin/sh", "-c", conf.Command)
	prefixColor := rand.Intn(16) + 1
	prefixStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(fmt.Sprint(prefixColor)))
	prefix := prefixStyle.Render(fmt.Sprintf("[%s]", conf.Label))

	return &Command{
		c:      c,
		prefix: prefix,
	}
}

func (c *Command) Run(ctx context.Context, conf RunCommandConfig) error {
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
					c.c.Process.Signal(syscall.SIGINT)
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
		return err
	}

	return nil
}

func (c *Command) ReadOut() string {
	return c.out.String()
}

func init() {
	// Seed random for random prefix colors.
	rand.Seed(time.Now().UnixNano())
}
