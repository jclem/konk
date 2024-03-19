package konkfile

import (
	"context"
	"fmt"

	"github.com/jclem/konk/konk"
	"github.com/jclem/konk/konk/debugger"
	"github.com/mattn/go-shellwords"
)

type ExecuteConfig struct {
	NoColor bool
	NoShell bool
}

func Execute(ctx context.Context, file File, command string, cfg ExecuteConfig) error {
	ex := newExecutor(file, cfg)
	return ex.execute(ctx, command)
}

type executor struct {
	file File
	cfg  ExecuteConfig
}

func (e *executor) execute(ctx context.Context, cmdName string) error {
	dbg := debugger.Get(ctx)

	cmd, ok := e.file.Commands[cmdName]
	if !ok {
		return fmt.Errorf("command not found: %s", cmdName)
	}

	for _, dep := range cmd.Dependencies {
		if err := e.execute(ctx, dep); err != nil {
			return fmt.Errorf("running dependency %s: %w", dep, err)
		}
	}

	if cmd.Run == "" {
		return nil
	}

	var c *konk.Command

	if e.cfg.NoShell {
		parts, err := shellwords.Parse(cmd.Run)
		if err != nil {
			return fmt.Errorf("parsing command: %w", err)
		}

		c = konk.NewCommand(konk.CommandConfig{
			Name:    parts[0],
			Args:    parts[1:],
			Label:   cmdName,
			NoColor: e.cfg.NoColor,
		})
	} else {
		c = konk.NewShellCommand(konk.ShellCommandConfig{
			Command: cmd.Run,
			Label:   cmdName,
			NoColor: e.cfg.NoColor,
		})
	}

	dbg.Prettyln(c)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if err := c.Run(ctx, cancel, konk.RunCommandConfig{}); err != nil {
		return fmt.Errorf("running command: %w", err)
	}

	return nil
}

func newExecutor(file File, cfg ExecuteConfig) *executor {
	return &executor{file, cfg}
}
