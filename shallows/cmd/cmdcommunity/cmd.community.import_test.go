package cmdcommunity

import (
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/stretchr/testify/require"
)

func TestCommunityImport(t *testing.T) {
	t.Run("dry run outputs records", func(t *testing.T) {
		var (
			ctx, cancel = testx.Context(t)
			tmpdir      = t.TempDir()
			inputPath   = filepath.Join(tmpdir, "input.jsonl")
			outputPath  = filepath.Join(tmpdir, "output.jsonl")
		)
		defer cancel()

		input := `{"id":"550e8400-e29b-41d4-a716-446655440000","domain":"testcommunity"}
{"known_media_id":"550e8400-e29b-41d4-a716-446655440001","magnet_uri":"magnet:?xt=urn:btih:abc123"}
{"known_media_id":"550e8400-e29b-41d4-a716-446655440002","magnet_uri":"magnet:?xt=urn:btih:def456"}
`
		require.NoError(t, os.WriteFile(inputPath, []byte(input), 0644))

		cmd := cmdCommunityImport{
			DryRun: true,
			Input:  inputPath,
			Output: outputPath,
		}
		gctx := &cmdopts.Global{
			Context:  ctx,
			Shutdown: cancel,
			Cleanup:  &sync.WaitGroup{},
		}

		require.NoError(t, cmd.Run(gctx))

		outFile, err := os.Open(outputPath)
		require.NoError(t, err)
		defer outFile.Close()

		require.Equal(t, "eeb5c509-fa5f-5600-0241-aec899a97032", testx.IOMD5(outFile))
	})

	t.Run("fails on invalid json", func(t *testing.T) {
		var (
			ctx, cancel = testx.Context(t)
			tmpdir      = t.TempDir()
			inputPath   = filepath.Join(tmpdir, "input.jsonl")
			outputPath  = filepath.Join(tmpdir, "output.jsonl")
		)
		defer cancel()

		input := `not valid json`
		require.NoError(t, os.WriteFile(inputPath, []byte(input), 0644))

		cmd := cmdCommunityImport{
			DryRun: true,
			Input:  inputPath,
			Output: outputPath,
		}
		gctx := &cmdopts.Global{
			Context:  ctx,
			Shutdown: cancel,
			Cleanup:  &sync.WaitGroup{},
		}

		err := cmd.Run(gctx)
		require.Error(t, err)
	})
}
