package integration_test

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunSerially(t *testing.T) {
	t.Parallel()

	out, err := newSerialRunner().
		withFlags("echo a", "echo b", "echo c").
		run(t)
	require.NoError(t, err)

	assert.Equal(t, `[0] a
[1] b
[2] c
`, out, "output did not match expected output")
}

func TestRunSeriallyWithLabels(t *testing.T) {
	t.Parallel()

	out, err := newSerialRunner().
		withFlags("-l", "a", "-l", "b", "-l", "c",
			"echo a", "echo b", "echo c").
		run(t)
	require.NoError(t, err)

	assert.Equal(t, `[a] a
[b] b
[c] c
`, out, "output did not match expected output")
}

func TestRunSeriallyWithLabelsMismatch(t *testing.T) {
	t.Parallel()

	out, err := newSerialRunner().
		withFlags("-l", "a", "-l", "b",
			"echo a", "echo b", "echo c").
		run(t)

	if assert.Error(t, err) {
		assert.IsType(t, &exec.ExitError{}, err)
	}

	assert.Equal(t, "Error: number of names must match number of commands\n", out, "error output did not match expectation")
}

func TestRunSeriallyWithCommandLabels(t *testing.T) {
	t.Parallel()

	out, err := newSerialRunner().
		withFlags("-L",
			"echo a", "echo b", "echo c").
		run(t)
	require.NoError(t, err)

	assert.Equal(t, `[echo a] a
[echo b] b
[echo c] c
`, out, "output did not match expected output")
}

func TestRunSeriallyWithNpm(t *testing.T) {
	t.Parallel()

	out, err := newSerialRunner().
		withFlags(
			"-w", "fixtures/npm",
			"--npm", "echo-a",
			"--npm", "echo-b",
		).
		run(t)
	require.NoError(t, err)

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
`, out, "output did not match expected output")
}

func TestRunSeriallyWithNpmGlob(t *testing.T) {
	t.Parallel()

	out, err := newSerialRunner().
		withFlags(
			"-w", "fixtures/npm",
			"--npm", "echo-*",
		).
		run(t)
	require.NoError(t, err)

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
`, out, "output did not match expected output")
}

func TestRunSeriallyNoLabel(t *testing.T) {
	t.Parallel()

	out, err := newSerialRunner().withFlags("-B", "echo a", "echo b", "echo c").run(t)
	require.NoError(t, err)

	assert.Equal(t, `a
b
c
`, out, "output did not match expected output")
}

func newSerialRunner() runner {
	return newRunner("run").withFlags("serially")
}
