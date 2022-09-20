package test

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunSerially(t *testing.T) {
	t.Parallel()

	out := new(strings.Builder)
	cmd := exec.Command("bin/konk", "run", "serially", "echo a", "echo b", "echo c")
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

func TestRunSeriallyWithLabels(t *testing.T) {
	t.Parallel()

	out := new(strings.Builder)
	cmd := exec.Command(
		"bin/konk", "run", "serially",
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

func TestRunSeriallyWithLabelsMismatch(t *testing.T) {
	t.Parallel()

	out := new(strings.Builder)
	cmd := exec.Command(
		"bin/konk", "run", "serially",
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

func TestRunSeriallyWithCommandLabels(t *testing.T) {
	t.Parallel()

	out := new(strings.Builder)
	cmd := exec.Command(
		"bin/konk", "run", "serially", "-L",
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

func TestRunSeriallyWithNpm(t *testing.T) {
	t.Parallel()

	out := new(strings.Builder)
	cmd := exec.Command(
		"bin/konk", "run", "serially",
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

func TestRunSeriallyWithNpmGlob(t *testing.T) {
	t.Parallel()

	out := new(strings.Builder)
	cmd := exec.Command(
		"bin/konk", "run", "serially",
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
