package konk

import (
	"context"
	"fmt"

	"github.com/jclem/konk/konk/internal/env"
	"github.com/mattn/go-shellwords"
	"golang.org/x/sync/errgroup"
)

type RunConcurrentlyConfig struct {
	Commands        []string
	Labels          []string
	Env             []string
	OmitEnv         bool
	AggregateOutput bool
	ContinueOnError bool
	NoColor         bool
	NoShell         bool
}

func RunConcurrently(ctx context.Context, cfg RunConcurrentlyConfig) ([]*Command, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	eg, ctx := errgroup.WithContext(ctx)

	commands := make([]*Command, len(cfg.Commands))

	env, err := env.Parse(cfg.Env)
	if err != nil {
		return nil, fmt.Errorf("parsing env: %w", err)
	}

	for i, cmd := range cfg.Commands {
		var c *Command

		if cfg.NoShell {
			parts, err := shellwords.Parse(cmd)

			if err != nil {
				return nil, fmt.Errorf("parsing command: %w", err)
			}

			c = NewCommand(CommandConfig{
				Name:    parts[0],
				Args:    parts[1:],
				Label:   cfg.Labels[i],
				Env:     env,
				OmitEnv: cfg.OmitEnv,
				NoColor: cfg.NoColor,
			})
		} else {
			c = NewShellCommand(ShellCommandConfig{
				Command: cmd,
				Label:   cfg.Labels[i],
				Env:     env,
				OmitEnv: cfg.OmitEnv,
				NoColor: cfg.NoColor,
			})
		}

		commands[i] = c
	}

	for _, cmd := range commands {
		eg.Go(func() error {
			return cmd.Run(ctx, cancel, RunCommandConfig{
				AggregateOutput: cfg.AggregateOutput,
				StopOnCancel:    !cfg.ContinueOnError,
			})
		})
	}

	err = eg.Wait()
	if err != nil {
		err = fmt.Errorf("running commands: %w", err)
	}

	return commands, err
}
