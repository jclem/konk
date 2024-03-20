package konkfile

import (
	"context"
	"fmt"
	"sync"

	"github.com/jclem/konk/konk"
	"github.com/jclem/konk/konk/konkfile/internal/dag"
	"github.com/mattn/go-shellwords"
	"golang.org/x/sync/errgroup"
)

type ExecuteConfig struct {
	AggregateOutput bool
	ContinueOnError bool
	NoColor         bool
	NoShell         bool
}

func Execute(ctx context.Context, file File, command string, cfg ExecuteConfig) error {
	g := dag.New[string]()

	for name := range file.Commands {
		g.AddNode(name)
	}

	for name, cmd := range file.Commands {
		for _, dep := range cmd.Dependencies {
			if err := g.AddEdge(name, dep); err != nil {
				return fmt.Errorf("adding edge: %w", err)
			}
		}
	}

	s := &scheduler{wgs: make(map[string]*sync.WaitGroup, 0)}

	for _, n := range g.Nodes() {
		s.wgs[n] = new(sync.WaitGroup)
		s.wgs[n].Add(1)
	}

	mut := new(sync.Mutex)
	wg := new(sync.WaitGroup)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	onNode := func(n string) error {
		mut.Lock()

		cmd, ok := file.Commands[n]
		if !ok {
			return fmt.Errorf("command not found: %s", n)
		}

		if cmd.Exclusive {
			defer mut.Unlock()
			wg.Wait()
		} else {
			wg.Add(1)
			defer wg.Done()
			mut.Unlock()
		}

		var c *konk.Command
		if cfg.NoShell {
			parts, err := shellwords.Parse(cmd.Run)
			if err != nil {
				return fmt.Errorf("parsing command: %w", err)
			}

			c = konk.NewCommand(konk.CommandConfig{
				Name:    parts[0],
				Args:    parts[1:],
				Label:   n,
				NoColor: cfg.NoColor,
			})
		} else {
			c = konk.NewShellCommand(konk.ShellCommandConfig{
				Command: cmd.Run,
				Label:   n,
				NoColor: false,
			})
		}

		if err := c.Run(ctx, cancel, konk.RunCommandConfig{
			AggregateOutput: cfg.AggregateOutput,
			KillOnCancel:    !cfg.ContinueOnError,
		}); err != nil {
			return fmt.Errorf("running command: %w", err)
		}

		return nil
	}

	path, err := g.Visit(command)
	if err != nil {
		return fmt.Errorf("visiting node: %w", err)
	}

	var eg errgroup.Group
	for _, n := range path {
		n := n
		eg.Go(func() error {
			from := g.From(n)
			return s.run(n, from, onNode)
		})
	}

	if err := eg.Wait(); err != nil {
		return fmt.Errorf("running commands: %w", err)
	}

	return nil
}

type scheduler struct {
	wgs map[string]*sync.WaitGroup
}

func (s *scheduler) run(n string, deps []string, onNode func(string) error) error {
	defer s.wgs[n].Done()

	for _, dep := range deps {
		s.wgs[dep].Wait()
	}

	if err := onNode(n); err != nil {
		return fmt.Errorf("running node: %w", err)
	}

	return nil
}
