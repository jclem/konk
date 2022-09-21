package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProc(t *testing.T) {
	t.Parallel()

	out, err := newProcRunner().run(t)
	assert.NoError(t, err)

	assert.Equal(t, `[echo-a] a
[echo-b] b
[echo-c] 
`, sortOut(t, out), "output did not match expected output")
}

func TestProcEnvSpaces(t *testing.T) {
	t.Parallel()

	out, err := newProcRunner().withFlags(
		"-e", ".env-spaces",
		"-p", "Procfile-spaces").
		run(t)
	assert.NoError(t, err)

	assert.Equal(t, `[echo-abc] a b c
[echo-def] d "e" f
`, sortOut(t, out), "output did not match expected output")
}

func TestProcWithExternalEnvNoEnv(t *testing.T) {
	t.Parallel()

	out, err := newProcRunner().
		withFlags("-E").
		withEnv("A=new-a", "B=new-b", "C=new-c").
		run(t)
	assert.NoError(t, err)

	assert.Equal(t, `[echo-a] new-a
[echo-b] new-b
[echo-c] new-c
`, sortOut(t, out), "output did not match expected output")
}

func TestProcWithExternalEnvAndEnv(t *testing.T) {
	t.Parallel()

	out, err := newProcRunner().
		withEnv("C=c").
		run(t)
	assert.NoError(t, err)

	assert.Equal(t, `[echo-a] a
[echo-b] b
[echo-c] c
`, sortOut(t, out), "output did not match expected output")
}

func newProcRunner() runner {
	return newRunner("proc").withFlags("-w", "fixtures/proc")
}
