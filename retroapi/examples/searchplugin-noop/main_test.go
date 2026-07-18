package main

import (
	"bytes"
	"encoding/json"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/ddiscapi"
	"github.com/stretchr/testify/require"
)

func TestNoopBuildBakesSourceViaLdflags(t *testing.T) {
	bin := filepath.Join(t.TempDir(), "noop")
	build := exec.Command("go", "build", "-ldflags", "-X main.source=demo-site", "-o", bin, ".")
	out, err := build.CombinedOutput()
	require.NoError(t, err, string(out))

	var stdout bytes.Buffer
	run := exec.Command(bin, "plugin", "--query", "ubuntu")
	run.Stdout = &stdout
	require.NoError(t, run.Run())

	var imp ddiscapi.Import
	require.NoError(t, json.Unmarshal(stdout.Bytes(), &imp))
	require.Equal(t, "demo-site", imp.Source)
}

func TestNoopBuildDefaultsSourceEmptyWithoutLdflags(t *testing.T) {
	bin := filepath.Join(t.TempDir(), "noop")
	build := exec.Command("go", "build", "-o", bin, ".")
	out, err := build.CombinedOutput()
	require.NoError(t, err, string(out))

	var stdout bytes.Buffer
	run := exec.Command(bin, "plugin", "--query", "ubuntu")
	run.Stdout = &stdout
	require.NoError(t, run.Run())

	var imp ddiscapi.Import
	require.NoError(t, json.Unmarshal(stdout.Bytes(), &imp))
	require.Empty(t, imp.Source)
}

func TestNoopRuntimeFlagOverridesBakedSource(t *testing.T) {
	bin := filepath.Join(t.TempDir(), "noop")
	build := exec.Command("go", "build", "-ldflags", "-X main.source=baked-site", "-o", bin, ".")
	out, err := build.CombinedOutput()
	require.NoError(t, err, string(out))

	var stdout bytes.Buffer
	run := exec.Command(bin, "plugin", "--query", "ubuntu", "--source", "runtime-site")
	run.Stdout = &stdout
	require.NoError(t, run.Run())

	var imp ddiscapi.Import
	require.NoError(t, json.Unmarshal(stdout.Bytes(), &imp))
	require.Equal(t, "runtime-site", imp.Source)
}
