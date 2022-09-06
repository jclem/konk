# Konk

Konk runs a series of commands in parallel or in serial. It is especially
well-suited to running multiple npm scripts.

## Why?

There are two npm packages I frequently already use for running npm scripts
serially or concurrently: `npm-run-all` and `concurrently`. I built konk because
I wanted something that could run in serial and in parallel and did not need to
be installed as an npm package. In addition, I have always been curious how to
build such a command line interface, so this is also a learning exercise for me.

There are currently feature gaps between `npm-run-all` and `concurrently`, but I
am working to fill them when I have time. Generally, I recommend using one of
those tools, instead, in an npm project, as they can be declared as a dependency
in your package.json file.

## Installation

```shell
$ go install github.com/jclem/konk@v0.1.0
```

## Usage

Run some commands in serial, one after the other. They will be run in the order
they're presented in the command line interface.

```shell
$ konk run s 'echo hello' 'echo world'
[0] hello
[1] world
```

Run some commands concurrently (notice the interleaved output).

```shell
$ konk run p ls ls
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

Use the `--npm` flag to easily run scripts in your `package.json`. For example,
given these scripts:

```json
{
  "scripts": {
    "check:build": "tsc --noEmit",
    "check:format": "prettier --loglevel warn --check .",
    "check:lint": "eslint ."
  }
}
```

You can run all three concurrently with with `--npm` flags:

```shell
$ konk run p --npm 'check:format' --npm 'check:lint' --npm 'check:types'
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
$ konk run p --npm 'check:*'
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
$ konk run p -g --npm 'check:*'
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
$ konk run p -gL --npm 'check:*'
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