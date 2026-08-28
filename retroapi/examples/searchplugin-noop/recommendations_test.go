package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/ddiscapi"
	"github.com/stretchr/testify/require"
)

func TestNoopRecommendationsDefaultLimit(t *testing.T) {
	bin := filepath.Join(t.TempDir(), "noop")
	build := exec.Command("go", "build", "-o", bin, ".")
	out, err := build.CombinedOutput()
	require.NoError(t, err, string(out))

	var stdout bytes.Buffer
	run := exec.Command(bin, "recommendations")
	run.Stdout = &stdout
	require.NoError(t, run.Run())

	var count int
	scanner := bufio.NewScanner(bytes.NewReader(stdout.Bytes()))
	for scanner.Scan() {
		var imp ddiscapi.Import
		require.NoError(t, json.Unmarshal(scanner.Bytes(), &imp))
		count++
	}
	require.NoError(t, scanner.Err())
	require.Equal(t, 5, count)
}

func TestNoopRecommendationsRespectsLimit(t *testing.T) {
	bin := filepath.Join(t.TempDir(), "noop")
	build := exec.Command("go", "build", "-o", bin, ".")
	out, err := build.CombinedOutput()
	require.NoError(t, err, string(out))

	var stdout bytes.Buffer
	run := exec.Command(bin, "recommendations", "--limit", "3")
	run.Stdout = &stdout
	require.NoError(t, run.Run())

	var count int
	scanner := bufio.NewScanner(bytes.NewReader(stdout.Bytes()))
	for scanner.Scan() {
		var imp ddiscapi.Import
		require.NoError(t, json.Unmarshal(scanner.Bytes(), &imp))
		count++
	}
	require.NoError(t, scanner.Err())
	require.Equal(t, 3, count)
}

func TestNoopRecommendationsTagsMimetype(t *testing.T) {
	bin := filepath.Join(t.TempDir(), "noop")
	build := exec.Command("go", "build", "-o", bin, ".")
	out, err := build.CombinedOutput()
	require.NoError(t, err, string(out))

	var stdout bytes.Buffer
	run := exec.Command(bin, "recommendations", "--limit", "2", "--mimetype", "video/mp4")
	run.Stdout = &stdout
	require.NoError(t, run.Run())

	var count int
	scanner := bufio.NewScanner(bytes.NewReader(stdout.Bytes()))
	for scanner.Scan() {
		var imp ddiscapi.Import
		require.NoError(t, json.Unmarshal(scanner.Bytes(), &imp))
		require.Equal(t, "video/mp4", imp.Mimetype)
		count++
	}
	require.NoError(t, scanner.Err())
	require.Equal(t, 2, count)
}

func TestNoopRecommendationsBakesSourceViaLdflags(t *testing.T) {
	bin := filepath.Join(t.TempDir(), "noop")
	build := exec.Command("go", "build", "-ldflags", "-X main.source=demo-site", "-o", bin, ".")
	out, err := build.CombinedOutput()
	require.NoError(t, err, string(out))

	var stdout bytes.Buffer
	run := exec.Command(bin, "recommendations", "--limit", "1")
	run.Stdout = &stdout
	require.NoError(t, run.Run())

	var imp ddiscapi.Import
	require.NoError(t, json.Unmarshal(stdout.Bytes(), &imp))
	require.Equal(t, "demo-site", imp.Source)
}
