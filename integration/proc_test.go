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
[echo-b] b
[echo-c] `, sortedOut, "output did not match expected output")
}

func TestProcWithExternalEnvNoEnv(t *testing.T) {
	out := new(strings.Builder)
	cmd := exec.Command(
		"bin/konk", "proc", "-E",
		"-w", "fixtures/proc")
	cmd.Env = append(cmd.Env, "A=new-a")
	cmd.Env = append(cmd.Env, "B=new-b")
	cmd.Env = append(cmd.Env, "C=new-c")
	cmd.Stdout = out
	cmd.Stderr = out
	if err := cmd.Run(); err != nil {
		t.Error(err)
	}

	lines := strings.Split(out.String(), "\n")
	sort.Strings(lines)
	sortedOut := strings.Join(lines, "\n")

	assert.Equal(t, `
[echo-a] new-a
[echo-b] new-b
[echo-c] new-c`, sortedOut, "output did not match expected output")
}

func TestProcWithExternalEnvAndEnv(t *testing.T) {
	out := new(strings.Builder)
	cmd := exec.Command(
		"bin/konk", "proc",
		"-w", "fixtures/proc")
	cmd.Env = append(cmd.Env, "C=c")
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
[echo-b] b
[echo-c] c`, sortedOut, "output did not match expected output")
}
