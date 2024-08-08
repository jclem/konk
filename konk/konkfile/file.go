package konkfile

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/jclem/konk/konk"
	"github.com/jclem/konk/konk/konkfile/internal/dag"
	"golang.org/x/exp/maps"
	"golang.org/x/sync/errgroup"
)

type File struct {
	Commands map[string]Command `json:"commands" toml:"commands" yaml:"commands"`
}

type Command struct {
	Run       string   `json:"run"       toml:"run"       yaml:"run"`
	Needs     []string `json:"needs"     toml:"needs"     yaml:"needs"`
	Exclusive bool     `json:"exclusive" toml:"exclusive" yaml:"exclusive"`
	Watch     []string `json:"watch"     toml:"watch"     yaml:"watch"`
}

type ExecuteConfig struct {
	AggregateOutput bool
	ContinueOnError bool
}

func (f File) Execute(originalCtx context.Context, command string, cfg ExecuteConfig) error {
	g := dag.New[string]()
	g.AddNodes(maps.Keys(f.Commands)...)

	for name, cmd := range f.Commands {
		if err := g.AddEdges(name, cmd.Needs...); err != nil {
			return fmt.Errorf("adding edge: %w", err)
		}
	}

	s := newScheduler(len(g.Nodes()))
	s.add(g.Nodes()...)

	mut := new(sync.Mutex)
	wg := new(sync.WaitGroup)

	onNode := func(n string) error {
		mut.Lock()

		cmd, ok := f.Commands[n]
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

		ctx, cancel := context.WithCancel(originalCtx)
		defer cancel()

		var runCommand func() error
		runCommand = func() error {
			ctx, cancel = context.WithCancel(originalCtx)
			defer cancel()

			c := konk.NewShellCommand(konk.ShellCommandConfig{
				Command: cmd.Run,
				Label:   n,
			})

			if err := c.Run(ctx, cancel, konk.RunCommandConfig{
				AggregateOutput: cfg.AggregateOutput,
				StopOnCancel:    cmd.Watch != nil || !cfg.ContinueOnError,
			}); err != nil {
				var xerr *konk.ExitError
				if errors.As(err, &xerr) {
					return runCommand()
				}

				return fmt.Errorf("running command: %w", err)
			}

			return nil
		}

		if cmd.Watch != nil {
			watcher, err := fsnotify.NewWatcher()
			if err != nil {
				return fmt.Errorf("creating watcher: %w", err)
			}
			defer watcher.Close()

			go func() {
				for {
					select {
					case event, ok := <-watcher.Events:
						if !ok {
							return
						}

						fmt.Printf("event: %s\n", event)

						cancel()
					case err, ok := <-watcher.Errors:
						if !ok {
							return
						}

						fmt.Printf("error: %s\n", err)
					}
				}
			}()

			for _, path := range cmd.Watch {
				if err := watcher.Add(path); err != nil {
					return fmt.Errorf("watching path: %w", err)
				}
			}

			if err := runCommand(); err != nil {
				return err
			}
		} else {
			if err := runCommand(); err != nil {
				return err
			}
		}

		return nil
	}

	path, err := g.VisitBreadthFirst(command)
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

func newScheduler(size int) *scheduler {
	return &scheduler{wgs: make(map[string]*sync.WaitGroup, size)}
}

func (s *scheduler) add(ns ...string) {
	for _, n := range ns {
		s.wgs[n] = new(sync.WaitGroup)
		s.wgs[n].Add(1)
	}
}

func (s *scheduler) run(n string, deps []string, onNode func(string) error) error {
	nodewg, ok := s.wgs[n]
	if !ok {
		return fmt.Errorf("node not found: %s", n)
	}

	defer nodewg.Done()

	for _, dep := range deps {
		depwg, ok := s.wgs[dep]
		if !ok {
			return fmt.Errorf("dependency not found: %s", dep)
		}

		depwg.Wait()
	}

	if err := onNode(n); err != nil {
		return fmt.Errorf("running node: %w", err)
	}

	return nil
}
