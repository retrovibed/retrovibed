package cmdtorrent_test

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"testing"
	"time"

	"github.com/alecthomas/kong"
	"github.com/james-lawrence/torrent/metainfo"
	"github.com/james-lawrence/torrent/torrenttest"
	"github.com/retrovibed/retrovibed/retroapi/bytesx"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdtestx"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdtorrent"
	"github.com/retrovibed/retrovibed/shallows/internal/env"
	"github.com/stretchr/testify/require"
)

func TestInspect(t *testing.T) {
	t.Run("prints details for a single torrent file", func(t *testing.T) {
		info, _, err := torrenttest.Random(t.TempDir(), 16*bytesx.KiB)
		require.NoError(t, err)

		md := metainfo.MetaInfo{
			InfoBytes:    testx.Must(metainfo.Encode(info))(t),
			Comment:      "inspect test torrent",
			CreatedBy:    "cmdtorrent_test",
			CreationDate: time.Now().Unix(),
		}
		path := filepath.Join(t.TempDir(), "example.torrent")
		require.NoError(t, os.WriteFile(path, testx.Must(metainfo.Encode(md))(t), 0600))

		var buf bytes.Buffer
		genparser := cmdtestx.Genparser(cmdtorrent.Commands{}, kong.Writers(&buf, nil), kong.Vars{
			"env_torrent_private":    env.TorrentPrivateNetwork,
			"vars_cores":             strconv.Itoa(runtime.GOMAXPROCS(0)),
			"vars_timestamp_started": time.Now().UTC().Format(time.RFC3339),
		})
		require.NoError(t, cmdtestx.Execute(t, genparser(t), "command", "inspect", path))

		require.Contains(t, buf.String(), md.HashInfoBytes().String())
		require.Contains(t, buf.String(), md.Comment)
		require.Contains(t, buf.String(), md.CreatedBy)
	})

	t.Run("prints details for each of multiple torrent files", func(t *testing.T) {
		info1, _, err := torrenttest.Random(t.TempDir(), 16*bytesx.KiB)
		require.NoError(t, err)
		md1 := metainfo.MetaInfo{InfoBytes: testx.Must(metainfo.Encode(info1))(t), CreationDate: time.Now().Unix()}
		path1 := filepath.Join(t.TempDir(), "first.torrent")
		require.NoError(t, os.WriteFile(path1, testx.Must(metainfo.Encode(md1))(t), 0600))

		info2, _, err := torrenttest.Random(t.TempDir(), 16*bytesx.KiB)
		require.NoError(t, err)
		md2 := metainfo.MetaInfo{InfoBytes: testx.Must(metainfo.Encode(info2))(t), CreationDate: time.Now().Unix()}
		path2 := filepath.Join(t.TempDir(), "second.torrent")
		require.NoError(t, os.WriteFile(path2, testx.Must(metainfo.Encode(md2))(t), 0600))

		var buf bytes.Buffer
		genparser := cmdtestx.Genparser(cmdtorrent.Commands{}, kong.Writers(&buf, nil), kong.Vars{
			"env_torrent_private":    env.TorrentPrivateNetwork,
			"vars_cores":             strconv.Itoa(runtime.GOMAXPROCS(0)),
			"vars_timestamp_started": time.Now().UTC().Format(time.RFC3339),
		})
		require.NoError(t, cmdtestx.Execute(t, genparser(t), "command", "inspect", path1, path2))

		require.Contains(t, buf.String(), md1.HashInfoBytes().String())
		require.Contains(t, buf.String(), md2.HashInfoBytes().String())
	})

	t.Run("wraps the error for a missing file", func(t *testing.T) {
		missing := filepath.Join(t.TempDir(), "does.not.exist.torrent")

		var buf bytes.Buffer
		genparser := cmdtestx.Genparser(cmdtorrent.Commands{}, kong.Writers(&buf, nil), kong.Vars{
			"env_torrent_private":    env.TorrentPrivateNetwork,
			"vars_cores":             strconv.Itoa(runtime.GOMAXPROCS(0)),
			"vars_timestamp_started": time.Now().UTC().Format(time.RFC3339),
		})
		err := cmdtestx.Execute(t, genparser(t), "command", "inspect", missing)
		require.Error(t, err)
		require.ErrorContains(t, err, "failed to read file")
		require.ErrorContains(t, err, missing)
	})

	t.Run("wraps the error for a torrent file with unparsable info", func(t *testing.T) {
		// a syntactically valid bencode integer is not a valid Info dict, so
		// LoadFromFile succeeds but UnmarshalInfo fails during printing.
		md := metainfo.MetaInfo{InfoBytes: []byte("i5e")}
		path := filepath.Join(t.TempDir(), "corrupt.torrent")
		require.NoError(t, os.WriteFile(path, testx.Must(metainfo.Encode(md))(t), 0600))

		var buf bytes.Buffer
		genparser := cmdtestx.Genparser(cmdtorrent.Commands{}, kong.Writers(&buf, nil), kong.Vars{
			"env_torrent_private":    env.TorrentPrivateNetwork,
			"vars_cores":             strconv.Itoa(runtime.GOMAXPROCS(0)),
			"vars_timestamp_started": time.Now().UTC().Format(time.RFC3339),
		})
		err := cmdtestx.Execute(t, genparser(t), "command", "inspect", path)
		require.Error(t, err)
		require.ErrorContains(t, err, "failed to format torrent")
		require.ErrorContains(t, err, path)
	})
}
