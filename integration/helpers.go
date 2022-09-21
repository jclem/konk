package test

import (
	"fmt"
	"os/exec"
	"sort"
	"strings"
	"testing"
)

type runner struct {
	cmd   string
	flags []string
	env   []string
}

func newRunner(cmd string) runner {
	return runner{cmd: cmd}
}

func (r runner) withFlags(flags ...string) runner {
	r.flags = append(r.flags, flags...)
	return r
}

func (r runner) withEnv(env ...string) runner {
	r.env = append(r.env, env...)
	return r
}

func (r runner) run(t *testing.T) (string, error) {
	t.Helper()

	out := new(strings.Builder)
	fullCmd := append([]string{r.cmd}, r.flags...)
	cmd := exec.Command("bin/konk", fullCmd...)
	cmd.Stdout = out
	cmd.Stderr = out
	cmd.Env = append(cmd.Env, r.env...)

	err := cmd.Run()
	return out.String(), err
}

func sortOut(t *testing.T, out string) string {
	t.Helper()

	lines := strings.Split(out, "\n")
	mapByPrefix := make(map[string][]string)

	for _, line := range lines {
		prefix := strings.SplitN(line, "]", 2)[0]

		if prefix == "" {
			continue
		}

		mapByPrefix[prefix] = append(mapByPrefix[prefix], line)
	}

	keys := make([]string, 0, len(mapByPrefix))
	for k := range mapByPrefix {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sortedLines []string

	for _, k := range keys {
		sortedLines = append(sortedLines, mapByPrefix[k]...)
	}

	// Our output always ends in a newline.
	return fmt.Sprintf("%s\n", strings.Join(sortedLines, "\n"))
}
