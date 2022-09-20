package test

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunConcurrently(t *testing.T) {
	t.Parallel()

	out := new(strings.Builder)
	cmd := exec.Command("bin/konk", "run", "concurrently", "-g",
		"echo a", "echo b", "echo c")
	cmd.Stdout = out
	cmd.Stderr = out
	if err := cmd.Run(); err != nil {
		t.Error(err)
	}
	assert.Equal(t, `[0] a
[1] b
[2] c
`, out.String(), "output did not match expected output")
}

func TestRunConcurrentlyWithLabels(t *testing.T) {
	t.Parallel()

	out := new(strings.Builder)
	cmd := exec.Command(
		"bin/konk", "run", "concurrently", "-g",
		"-l", "a", "-l", "b", "-l", "c",
		"echo a", "echo b", "echo c")
	cmd.Stdout = out
	cmd.Stderr = out
	if err := cmd.Run(); err != nil {
		t.Error(err)
	}
	assert.Equal(t, `[a] a
[b] b
[c] c
`, out.String(), "output did not match expected output")
}

func TestRunConcurrentlyWithLabelsMismatch(t *testing.T) {
	t.Parallel()

	out := new(strings.Builder)
	cmd := exec.Command(
		"bin/konk", "run", "concurrently", "-g",
		"-l", "a", "-l", "b",
		"echo a", "echo b", "echo c")
	cmd.Stdout = out
	cmd.Stderr = out

	err := cmd.Run()
	if assert.Error(t, err) {
		assert.IsType(t, &exec.ExitError{}, err)
	}

	assert.Equal(t, "Error: number of names must match number of commands\n", out.String(), "error output did not match expectation")
}

func TestRunConcurrentlyWithCommandLabels(t *testing.T) {
	t.Parallel()

	out := new(strings.Builder)
	cmd := exec.Command(
		"bin/konk", "run", "concurrently", "-gL",
		"echo a", "echo b", "echo c")
	cmd.Stdout = out
	if err := cmd.Run(); err != nil {
		t.Error(err)
	}
	assert.Equal(t, `[echo a] a
[echo b] b
[echo c] c
`, out.String(), "output did not match expected output")
}

func TestRunConcurrentlyWithNpm(t *testing.T) {
	t.Parallel()

	out := new(strings.Builder)
	cmd := exec.Command(
		"bin/konk", "run", "concurrently", "-g",
		"-w", "fixtures/npm",
		"--npm", "echo-a",
		"--npm", "echo-b")
	cmd.Stdout = out
	if err := cmd.Run(); err != nil {
		t.Error(err)
	}
	assert.Equal(t, `[0] 
[0] > echo-a
[0] > echo a
[0] 
[0] a
[1] 
[1] > echo-b
[1] > echo b
[1] 
[1] b
`, out.String(), "output did not match expected output")
}

func TestRunConcurrentlyWithNpmGlob(t *testing.T) {
	t.Parallel()

	out := new(strings.Builder)
	cmd := exec.Command(
		"bin/konk", "run", "concurrently", "-g",
		"-w", "fixtures/npm",
		"--npm", "echo-*")
	cmd.Stdout = out
	if err := cmd.Run(); err != nil {
		t.Error(err)
	}
	assert.Equal(t, `[0] 
[0] > echo-a
[0] > echo a
[0] 
[0] a
[1] 
[1] > echo-b
[1] > echo b
[1] 
[1] b
`, out.String(), "output did not match expected output")
}
