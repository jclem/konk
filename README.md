# Konk

Konk runs a series of commands serially or concurrently. It is especially
well-suited to running multiple npm scripts.

## Why?

There are two npm packages I frequently already use for running npm scripts
serially or concurrently: `npm-run-all` and `concurrently`. I built konk because
I wanted something that could run serially and concurrently and did not need to
be installed as an npm package (note, however, that konk _can_ be installed from
npm). In addition, I wanted to be able to use the same command line interface to
run processes defined in a Procfile. Finally, I have always been curious how to
build such a command line interface, so this is also a learning exercise for me.

There are currently feature gaps between `npm-run-all` and `concurrently`, but I
am working to fill them when I have time.

## Installation

Install via Homebrew:

```shell
$ brew install jclem/tap/konk
```

Or, use or install directly from npm:

```shell
$ npx konk      # Run from npm
$ npm i -g konk # Install from npm
```

## konk

Konk is a tool for running multiple processes

### Options

```
  -D, --debug   debug mode
  -h, --help    help for konk
```

### SEE ALSO

- [konk completion](#konk-completion) - Generate the autocompletion script for the specified shell
- [konk docs](#konk-docs) - Print documentation
- [konk proc](#konk-proc) - Run commands defined in a Procfile (alias: p)
- [konk run](#konk-run) - Run commands serially or concurrently (alias: r)

## konk completion

Generate the autocompletion script for the specified shell

### Synopsis

Generate the autocompletion script for konk for the specified shell.
See each sub-command's help for details on how to use the generated script.

### Options

```
  -h, --help   help for completion
```

### Options inherited from parent commands

```
  -D, --debug   debug mode
```

### SEE ALSO

- [konk](#konk) - Konk is a tool for running multiple processes
- [konk completion bash](#konk-completion-bash) - Generate the autocompletion script for bash
- [konk completion fish](#konk-completion-fish) - Generate the autocompletion script for fish
- [konk completion powershell](#konk-completion-powershell) - Generate the autocompletion script for powershell
- [konk completion zsh](#konk-completion-zsh) - Generate the autocompletion script for zsh

## konk completion bash

Generate the autocompletion script for bash

### Synopsis

Generate the autocompletion script for the bash shell.

This script depends on the 'bash-completion' package.
If it is not installed already, you can install it via your OS's package manager.

To load completions in your current shell session:

    source <(konk completion bash)

To load completions for every new session, execute once:

#### Linux:

    konk completion bash > /etc/bash_completion.d/konk

#### macOS:

    konk completion bash > $(brew --prefix)/etc/bash_completion.d/konk

You will need to start a new shell for this setup to take effect.

```
konk completion bash
```

### Options

```
  -h, --help              help for bash
      --no-descriptions   disable completion descriptions
```

### Options inherited from parent commands

```
  -D, --debug   debug mode
```

### SEE ALSO

- [konk completion](#konk-completion) - Generate the autocompletion script for the specified shell

## konk completion fish

Generate the autocompletion script for fish

### Synopsis

Generate the autocompletion script for the fish shell.

To load completions in your current shell session:

    konk completion fish | source

To load completions for every new session, execute once:

    konk completion fish > ~/.config/fish/completions/konk.fish

You will need to start a new shell for this setup to take effect.

```
konk completion fish [flags]
```

### Options

```
  -h, --help              help for fish
      --no-descriptions   disable completion descriptions
```

### Options inherited from parent commands

```
  -D, --debug   debug mode
```

### SEE ALSO

- [konk completion](#konk-completion) - Generate the autocompletion script for the specified shell

## konk completion powershell

Generate the autocompletion script for powershell

### Synopsis

Generate the autocompletion script for powershell.

To load completions in your current shell session:

    konk completion powershell | Out-String | Invoke-Expression

To load completions for every new session, add the output of the above command
to your powershell profile.

```
konk completion powershell [flags]
```

### Options

```
  -h, --help              help for powershell
      --no-descriptions   disable completion descriptions
```

### Options inherited from parent commands

```
  -D, --debug   debug mode
```

### SEE ALSO

- [konk completion](#konk-completion) - Generate the autocompletion script for the specified shell

## konk completion zsh

Generate the autocompletion script for zsh

### Synopsis

Generate the autocompletion script for the zsh shell.

If shell completion is not already enabled in your environment you will need
to enable it. You can execute the following once:

    echo "autoload -U compinit; compinit" >> ~/.zshrc

To load completions in your current shell session:

    source <(konk completion zsh)

To load completions for every new session, execute once:

#### Linux:

    konk completion zsh > "${fpath[1]}/_konk"

#### macOS:

    konk completion zsh > $(brew --prefix)/share/zsh/site-functions/_konk

You will need to start a new shell for this setup to take effect.

```
konk completion zsh [flags]
```

### Options

```
  -h, --help              help for zsh
      --no-descriptions   disable completion descriptions
```

### Options inherited from parent commands

```
  -D, --debug   debug mode
```

### SEE ALSO

- [konk completion](#konk-completion) - Generate the autocompletion script for the specified shell

## konk docs

Print documentation

```
konk docs [flags]
```

### Options

```
  -f, --format string   output format (default "markdown")
  -h, --help            help for docs
```

### Options inherited from parent commands

```
  -D, --debug   debug mode
```

### SEE ALSO

- [konk](#konk) - Konk is a tool for running multiple processes

## konk proc

Run commands defined in a Procfile (alias: p)

```
konk proc [flags]
```

### Options

```
  -c, --continue-on-error          continue running commands after a failure
  -e, --env-file string            Path to the env file (default ".env")
  -h, --help                       help for proc
  -C, --no-color                   do not colorize label output
  -E, --no-env-file                Don't load the env file
  -B, --no-label                   do not attach label/prefix to output
  -S, --no-subshell                do not run commands in a subshell
      --omit-env                   Omit any existing runtime environment variables
  -p, --procfile string            Path to the Procfile (default "Procfile")
  -w, --working-directory string   set the working directory for all commands
```

### Options inherited from parent commands

```
  -D, --debug   debug mode
```

### SEE ALSO

- [konk](#konk) - Konk is a tool for running multiple processes

## konk run

Run commands serially or concurrently (alias: r)

```
konk run <subcommand> [flags]
```

### Options

```
  -b, --bun                        Run npm commands with Bun
  -L, --command-as-label           use each command as its own label
  -c, --continue-on-error          continue running commands after a failure
  -h, --help                       help for run
  -l, --label stringArray          label prefix for the command
  -C, --no-color                   do not colorize label output
  -B, --no-label                   do not attach label/prefix to output
  -S, --no-subshell                do not run commands in a subshell
  -n, --npm stringArray            npm command
  -w, --working-directory string   set the working directory for all commands
```

### Options inherited from parent commands

```
  -D, --debug   debug mode
```

### SEE ALSO

- [konk](#konk) - Konk is a tool for running multiple processes
- [konk run concurrently](#konk-run-concurrently) - Run commands concurrently (alias: c)
- [konk run serially](#konk-run-serially) - Run commands serially (alias: s)

## konk run concurrently

Run commands concurrently (alias: c)

```
konk run concurrently <command...> [flags]
```

### Examples

```
# Run two commands concurrently

konk run concurrently "script/api-server" "script/frontend-server"

# Run a set of npm commands concurrently

konk run concurrently -n lint -n test

# Run a set of npm commands concurrently, but aggregate their output

konk run concurrently -g -n lint -n test

# Run all npm commands prefixed with "check:" concurrently using Bun, ignore
# errors, aggregate output, and use the script name as the label

konk run concurrently -bgcL -n "check:*"
```

### Options

```
  -g, --aggregate-output   aggregate command output
  -h, --help               help for concurrently
```

### Options inherited from parent commands

```
  -b, --bun                        Run npm commands with Bun
  -L, --command-as-label           use each command as its own label
  -c, --continue-on-error          continue running commands after a failure
  -D, --debug                      debug mode
  -l, --label stringArray          label prefix for the command
  -C, --no-color                   do not colorize label output
  -B, --no-label                   do not attach label/prefix to output
  -S, --no-subshell                do not run commands in a subshell
  -n, --npm stringArray            npm command
  -w, --working-directory string   set the working directory for all commands
```

### SEE ALSO

- [konk run](#konk-run) - Run commands serially or concurrently (alias: r)

## konk run serially

Run commands serially (alias: s)

```
konk run serially <command...> [flags]
```

### Examples

```
# Run two commands in serial

konk run serially "echo foo" "echo bar"

# Run a set of npm commands in serial

konk run serially -n build -n deploy
```

### Options

```
  -h, --help   help for serially
```

### Options inherited from parent commands

```
  -b, --bun                        Run npm commands with Bun
  -L, --command-as-label           use each command as its own label
  -c, --continue-on-error          continue running commands after a failure
  -D, --debug                      debug mode
  -l, --label stringArray          label prefix for the command
  -C, --no-color                   do not colorize label output
  -B, --no-label                   do not attach label/prefix to output
  -S, --no-subshell                do not run commands in a subshell
  -n, --npm stringArray            npm command
  -w, --working-directory string   set the working directory for all commands
```

### SEE ALSO

- [konk run](#konk-run) - Run commands serially or concurrently (alias: r)
