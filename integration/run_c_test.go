package test

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunConcurrently(t *testing.T) {
	t.Parallel()

	out, err := newRunner("run").
		withFlags("concurrently", "-g", "echo a", "echo b", "echo c").
		run(t)
	assert.NoError(t, err)

	assert.Equal(t, `[0] a
[1] b
[2] c
`, sortOut(t, out), "output did not match expected output")
}

func TestRunConcurrentlyWithLabels(t *testing.T) {
	t.Parallel()

	out, err := newRunner("run").
		withFlags("concurrently", "-g",
			"-l", "a", "-l", "b", "-l", "c",
			"echo a", "echo b", "echo c").
		run(t)
	assert.NoError(t, err)

	assert.Equal(t, `[a] a
[b] b
[c] c
`, sortOut(t, out), "output did not match expected output")
}

func TestRunConcurrentlyWithLabelsMismatch(t *testing.T) {
	t.Parallel()

	out, err := newRunner("run").
		withFlags("concurrently", "-g",
			"-l", "a", "-l", "b",
			"echo a", "echo b", "echo c").
		run(t)

	if assert.Error(t, err) {
		assert.IsType(t, &exec.ExitError{}, err)
	}

	assert.Equal(t, "Error: number of names must match number of commands\n", out, "error output did not match expectation")
}

func TestRunConcurrentlyWithCommandLabels(t *testing.T) {
	t.Parallel()

	out, err := newRunner("run").
		withFlags("concurrently", "-gL",
			"echo a", "echo b", "echo c").
		run(t)
	assert.NoError(t, err)

	assert.Equal(t, `[echo a] a
[echo b] b
[echo c] c
`, sortOut(t, out), "output did not match expected output")
}

func TestRunConcurrentlyWithNpm(t *testing.T) {
	t.Parallel()

	out, err := newRunner("run").
		withFlags("concurrently", "-g",
			"-w", "fixtures/npm",
			"--npm", "echo-a",
			"--npm", "echo-b").
		run(t)
	assert.NoError(t, err)

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
`, sortOut(t, out), "output did not match expected output")
}

func TestRunConcurrentlyWithNpmGlob(t *testing.T) {
	t.Parallel()

	out, err := newRunner("run").
		withFlags("concurrently", "-g",
			"-w", "fixtures/npm",
			"--npm", "echo-*").
		run(t)
	assert.NoError(t, err)

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
`, sortOut(t, out), "output did not match expected output")
}
