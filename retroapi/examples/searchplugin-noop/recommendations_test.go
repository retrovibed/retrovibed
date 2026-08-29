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

func TestNoopRecommendations(t *testing.T) {
	t.Run("DefaultLimit", func(t *testing.T) {
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
	})

	t.Run("RespectsLimit", func(t *testing.T) {
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
	})

	t.Run("TagsMimetype", func(t *testing.T) {
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
	})

	t.Run("BakesSourceViaLdflags", func(t *testing.T) {
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
	})

	t.Run("CyclesLicenseStatus", func(t *testing.T) {
		bin := filepath.Join(t.TempDir(), "noop")
		build := exec.Command("go", "build", "-o", bin, ".")
		out, err := build.CombinedOutput()
		require.NoError(t, err, string(out))

		var stdout bytes.Buffer
		run := exec.Command(bin, "recommendations", "--limit", "4")
		run.Stdout = &stdout
		require.NoError(t, run.Run())

		var licensed []ddiscapi.Import_LicenseStatus
		scanner := bufio.NewScanner(bytes.NewReader(stdout.Bytes()))
		for scanner.Scan() {
			var imp ddiscapi.Import
			require.NoError(t, json.Unmarshal(scanner.Bytes(), &imp))
			licensed = append(licensed, imp.Licensed)
		}
		require.NoError(t, scanner.Err())
		require.Equal(t, []ddiscapi.Import_LicenseStatus{
			ddiscapi.Import_Unknown,
			ddiscapi.Import_Unlicensed,
			ddiscapi.Import_Licensed,
			ddiscapi.Import_Unknown,
		}, licensed)
	})

	t.Run("PublicYieldsZeroResults", func(t *testing.T) {
		bin := filepath.Join(t.TempDir(), "noop")
		build := exec.Command("go", "build", "-o", bin, ".")
		out, err := build.CombinedOutput()
		require.NoError(t, err, string(out))

		var stdout bytes.Buffer
		run := exec.Command(bin, "recommendations", "--public")
		run.Stdout = &stdout
		require.NoError(t, run.Run())

		require.Empty(t, stdout.Bytes())
	})
}
