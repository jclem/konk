package main

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunSerially(t *testing.T) {
	out := new(strings.Builder)
	cmd := exec.Command("go", "run", "..", "run", "serially", "echo a", "echo b", "echo c")
	cmd.Stdout = out
	if err := cmd.Run(); err != nil {
		t.Error(err)
	}
	assert.Equal(t, `[0] a
[1] b
[2] c
`, out.String(), "output did not match expected output")
}

func TestRunSeriallyWithLabels(t *testing.T) {
	out := new(strings.Builder)
	cmd := exec.Command(
		"go", "run", "..", "run", "serially",
		"-l", "a", "-l", "b", "-l", "c",
		"echo a", "echo b", "echo c")
	cmd.Stdout = out
	if err := cmd.Run(); err != nil {
		t.Error(err)
	}
	assert.Equal(t, `[a] a
[b] b
[c] c
`, out.String(), "output did not match expected output")
}
