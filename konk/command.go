package konk

import (
	"bufio"
	"fmt"
	"math/rand"
	"os/exec"
	"time"

	"github.com/charmbracelet/lipgloss"
)

type Command struct {
	c      *exec.Cmd
	prefix string
}

type CommandConfig struct {
	Name  string
	Args  []string
	Label string
}

func NewCommand(conf CommandConfig) *Command {
	c := exec.Command(conf.Name, conf.Args...)
	c.Stderr = c.Stdout
	prefixColor := rand.Intn(16) + 1
	prefixStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(fmt.Sprint(prefixColor)))
	prefix := prefixStyle.Render(fmt.Sprintf("[%s]", conf.Label))

	return &Command{
		c:      c,
		prefix: prefix,
	}
}

func (c *Command) Run() error {
	stdout, err := c.c.StdoutPipe()
	if err != nil {
		return err
	}

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
			t := <-out
			fmt.Println(fmt.Sprintf("%s %s", c.prefix, t))
		}
	}()

	if err := c.c.Wait(); err != nil {
		return err
	}

	return nil
}

func init() {
	// Seed random for random prefix colors.
	rand.Seed(time.Now().UnixNano())
}
