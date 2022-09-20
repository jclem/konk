package test

import (
	"os/exec"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProc(t *testing.T) {
	out := new(strings.Builder)
	cmd := exec.Command(
		"bin/konk", "proc",
		"-w", "fixtures/proc")
	cmd.Stdout = out
	cmd.Stderr = out
	if err := cmd.Run(); err != nil {
		t.Error(err)
	}

	lines := strings.Split(out.String(), "\n")
	sort.Strings(lines)
	sortedOut := strings.Join(lines, "\n")

	assert.Equal(t, `
[echo-a] a
[echo-b] b`, sortedOut, "output did not match expected output")
}
