package konkfile

import (
	"context"
	"fmt"
	"sync"

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
	wg   *sync.WaitGroup
	mut  *sync.Mutex
}

func (e *executor) execute(ctx context.Context, cmdName string) error {
	dbg := debugger.Get(ctx)

	cmd, ok := e.file.Commands[cmdName]
	if !ok {
		return fmt.Errorf("command not found: %s", cmdName)
	}

	wg := new(sync.WaitGroup)
	for _, dep := range cmd.Dependencies {
		wg.Add(1)

		go func(dep string) {
			defer wg.Done()
			if err := e.execute(ctx, dep); err != nil {
				panic(fmt.Errorf("running dependency %s: %w", dep, err))
			}
		}(dep)
	}

	wg.Wait()

	if cmd.Run == "" {
		return nil
	}

	// Concurrency control:
	// - The mutex ensures that no other commands run while an exclusive command
	// 	 is running.
	// - The wait group ensures an exclusive command waits for all other commands
	//   to complete.

	e.mut.Lock()

	if cmd.Exclusive {
		defer e.mut.Unlock()
		e.wg.Wait()
	} else {
		e.wg.Add(1)
		defer e.wg.Done()
		e.mut.Unlock()
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
	return &executor{
		file: file,
		cfg:  cfg,
		mut:  new(sync.Mutex),
		wg:   new(sync.WaitGroup),
	}
}
