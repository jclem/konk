package konk

import (
	"context"
	"fmt"

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

	env, err := parseEnv(cfg.Env)
	if err != nil {
		return nil, err
	}

	for i, cmd := range cfg.Commands {
		var c *Command

		if cfg.NoShell {
			parts, err := shellwords.Parse(cmd)

			if err != nil {
				return nil, err
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
		cmd := cmd

		eg.Go(func() error {
			return cmd.Run(ctx, cancel, RunCommandConfig{
				AggregateOutput: cfg.AggregateOutput,
				KillOnCancel:    !cfg.ContinueOnError,
			})
		})
	}

	err = eg.Wait()
	return commands, err
}

func parseEnv(env []string) ([]string, error) {
	var parsedEnv []string

	// Unquote any quoted .env vars.
	for _, line := range env {
		parsed, err := shellwords.Parse(line)
		if err != nil {
			return nil, err
		}

		if len(parsed) == 0 {
			continue
		}

		if len(parsed) != 1 {
			return nil, fmt.Errorf("invalid .env line: %s", line)
		}

		parsedEnv = append(parsedEnv, parsed[0])
	}

	return parsedEnv, nil
}
