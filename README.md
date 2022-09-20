# Konk

Konk runs a series of commands serially or concurrently. It is especially
well-suited to running multiple npm scripts.

## Why?

There are two npm packages I frequently already use for running npm scripts
serially or concurrently: `npm-run-all` and `concurrently`. I built konk because
I wanted something that could run serially and concurrently and did not need to
be installed as an npm package (note, however, that konk *can* be installed from
npm). In addition, I wanted to be able to use the same command line interface to
run processes defined in a Procfile. Finally, I have always been curious how to
build such a command line interface, so this is also a learning exercise for me.

There are currently feature gaps between `npm-run-all` and `concurrently`, but I
am working to fill them when I have time.

## Installation

Install the plain Go module:

```shell
$ go install github.com/jclem/konk@latest
```

Or, use or install directly from npm:

```shell
$ npx konk      # Run from npm
$ npm i -g konk # Install from npm
```

## Usage

### Run Commands: `konk run`

Konk can run an arbitrary list of commands passed to it, either serially or
concurrently.

Run some commands serially, one after the other. They will be run in the order
they're presented in the command line interface.

```shell
$ konk run s 'echo hello' 'echo world'
[0] hello
[1] world
```

Run some commands concurrently (notice the interleaved output).

```shell
$ konk run c ls ls
[1] README.md
[1] bin
[1] go.mod
[1] go.sum
[0] README.md
[1] konk
[1] main.go
[0] bin
[0] go.mod
[0] go.sum
[0] konk
[0] main.go
```

Use the `-n/--npm` flag to easily run scripts in your `package.json`. For
example, given these scripts:

```json
{
  "scripts": {
    "check:build": "tsc --noEmit",
    "check:format": "prettier --loglevel warn --check .",
    "check:lint": "eslint ."
  }
}
```

You can run all three concurrently with with `-n/--npm` flags:

```shell
$ konk run c --npm 'check:format' --npm 'check:lint' --npm 'check:types'
[2]
[2] > check:lint
[2] > eslint .
[2]
[1]
[1] > check:format
[0]
[0] > check:build
[0] > tsc --noEmit
[0]
[1] > prettier --loglevel warn --check .
[1]
```

Or, you can use a glob-like pattern for brevity:

```shell
$ konk run c -n 'check:*'
[2]
[2] > check:lint
[2] > eslint .
[2]
[1]
[1] > check:format
[0]
[0] > check:build
[0] > tsc --noEmit
[0]
[1] > prettier --loglevel warn --check .
[1]
```

If you want to run commands concurrently but want to ensure output is not
interleaved, use the `-g`/`--aggregate-output` flag:

```shell
$ konk run c -g -n 'check:*'
[0]
[0] > check:build
[0] > tsc --noEmit
[0]
[1]
[1] > check:format
[1] > prettier --loglevel warn --check .
[1]
[2]
[2] > check:lint
[2] > eslint .
[2]
```

You can also use the `-L`/`--command-as-label` flag to use the command itself as
the process label:

```shell
$ konk run c -gL -n 'check:*'
[check:build ]
[check:build ] > check:build
[check:build ] > tsc --noEmit
[check:build ]
[check:format]
[check:format] > check:format
[check:format] > prettier --loglevel warn --check .
[check:format]
[check:lint  ]
[check:lint  ] > check:lint
[check:lint  ] > eslint .
[check:lint  ]
```

There is also a `-c/--continue-on-error` flag that will ensure other commands
continue to run even if one fails. The default behavior is that all commands
halt when any other command exits with a non-zero exit code.

```shell
$ konk run c -cgL -n 'check:*'
[check:build ]
[check:build ] > check:build
[check:build ] > tsc --noEmit
[check:build ]
[check:format]
[check:format] > check:format
[check:format] > prettier --loglevel warn --check .
[check:format]
[check:lint  ]
[check:lint  ] > check:lint
[check:lint  ] > eslint .
[check:lint  ]
```

### Run Procfile: `konk proc`

Konk can also run a set of commands defined in a Procfile, using `konk proc`.

Given a Procfile like this:

```procfile
foo: echo foo
bar: echo bar
```

```shell
$ konk proc
[foo] foo
[bar] bar
```

Konk will run the "foo" and "bar" commands concurrently (as if they were run
with `konk run c`). In addition, Konk will pass any environment variables
defined in a `.env` file to the running commands.