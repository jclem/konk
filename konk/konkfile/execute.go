package konkfile

import (
	"context"
	"fmt"
	"sync"

	"github.com/jclem/konk/konk"
	"golang.org/x/sync/errgroup"
)

type ExecuteConfig struct {
	NoColor bool
	NoShell bool
}

func Execute(ctx context.Context, file File, command string, cfg ExecuteConfig) error {
	dag := newDAG[string]()

	for name := range file.Commands {
		dag.addNode(name)
	}

	for name, cmd := range file.Commands {
		for _, dep := range cmd.Dependencies {
			if err := dag.addEdge(edge[string]{name, dep}); err != nil {
				return fmt.Errorf("adding edge: %w", err)
			}
		}
	}

	s := &scheduler{wgs: make(map[string]*sync.WaitGroup, 0)}

	for _, node := range dag.nodes {
		s.wgs[node] = new(sync.WaitGroup)
		s.wgs[node].Add(1)
	}

	mut := new(sync.Mutex)

	wg := new(sync.WaitGroup)

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

		c := konk.NewShellCommand(konk.ShellCommandConfig{
			Command: cmd.Run,
			Label:   n,
			NoColor: false,
		})

		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		if err := c.Run(ctx, cancel, konk.RunCommandConfig{}); err != nil {
			return fmt.Errorf("running command: %w", err)
		}

		return nil
	}

	var eg errgroup.Group
	for _, n := range dag.nodes {
		n := n
		eg.Go(func() error {
			return s.run(n, dag.from(n), onNode)
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
