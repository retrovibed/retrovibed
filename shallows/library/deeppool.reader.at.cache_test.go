package library_test

import (
	"bytes"
	"context"
	"crypto/md5"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/retrovibed/retrovibed/retroapi/blockcache"
	"github.com/retrovibed/retrovibed/shallows/internal/bytesx"
	"github.com/retrovibed/retrovibed/shallows/internal/cryptox"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/md5x"
	"github.com/retrovibed/retrovibed/shallows/internal/testx"
	"github.com/retrovibed/retrovibed/shallows/internal/uuidx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/stretchr/testify/require"
)

func TestDeeppoolReaderAtDownload(t *testing.T) {
	download := func(src io.Reader, n int64) http.HandlerFunc {
		var (
			buf bytes.Buffer
		)

		_n, err := io.Copy(&buf, io.LimitReader(src, n))
		require.NoError(t, err)
		require.EqualValues(t, n, _n)

		return func(w http.ResponseWriter, r *http.Request) {
			http.ServeContent(w, r, "test.bin", time.Date(2025, time.July, 01, 0, 0, 0, 0, time.UTC), io.NewSectionReader(bytes.NewReader(buf.Bytes()), 0, n))
		}
	}

	t.Run("full download - n % DefaultBlockSize == 0", func(t *testing.T) {
		const nlen = 128 * bytesx.MiB

		var (
			digest = md5.New()
			seed   = errorsx.Must(uuid.NewV4())
		)

		md := library.NewMetadata(uuidx.WithSuffix(1), library.MetadataOptionEncryptionSeed(seed.String()), library.MetadataOptionBytes(nlen))
		src, err := cryptox.NewReaderChaCha20(library.MetadataChaCha8(md), io.TeeReader(cryptox.NewChaCha8(seed.String()), digest))
		require.NoError(t, err)

		routes := mux.NewRouter()
		routes.Handle("/m/{id}/download", alice.New().ThenFunc(download(src, nlen)))

		srv := httptest.NewServer(routes)
		defer srv.Close()

		c := &http.Client{}
		c.Transport = httpx.RewriteHostTransport(testx.Must(url.ParseRequestURI(srv.URL))(t), c.Transport)

		reader := library.NewDeeppoolReaderAt(c, md, testx.Must(blockcache.NewDirectoryCache(t.TempDir()))(t))

		downloaded := md5.New()
		n, err := io.CopyN(downloaded, io.NewSectionReader(reader, 0, nlen), nlen)
		require.NoError(t, err)
		require.EqualValues(t, nlen, n)

		require.Equal(t, md5x.FormatUUID(digest), md5x.FormatUUID(downloaded))
	})

	t.Run("full download - n % DefaultBlockSize == 1", func(t *testing.T) {
		const nlen = 128*bytesx.MiB + 1

		var (
			digest = md5.New()
			seed   = errorsx.Must(uuid.NewV4())
		)

		md := library.NewMetadata(uuidx.WithSuffix(1), library.MetadataOptionEncryptionSeed(seed.String()), library.MetadataOptionBytes(nlen))
		src, err := cryptox.NewReaderChaCha20(library.MetadataChaCha8(md), io.TeeReader(cryptox.NewChaCha8(seed.String()), digest))
		require.NoError(t, err)

		routes := mux.NewRouter()
		routes.Handle("/m/{id}/download", alice.New().ThenFunc(download(src, nlen)))

		srv := httptest.NewServer(routes)
		defer srv.Close()

		c := &http.Client{}
		c.Transport = httpx.RewriteHostTransport(testx.Must(url.ParseRequestURI(srv.URL))(t), c.Transport)

		reader := library.NewDeeppoolReaderAt(c, md, testx.Must(blockcache.NewDirectoryCache(t.TempDir()))(t))

		downloaded := md5.New()
		n, err := io.CopyN(downloaded, io.NewSectionReader(reader, 0, nlen), nlen)
		require.NoError(t, err)
		require.EqualValues(t, nlen, n)
		require.Equal(t, md5x.FormatUUID(digest), md5x.FormatUUID(downloaded))
	})

	t.Run("full download - n % DefaultBlockSize == DefaultBlockSize-1", func(t *testing.T) {
		const nlen = 128*bytesx.MiB - 1

		var (
			digest = md5.New()
			seed   = errorsx.Must(uuid.NewV4())
		)

		md := library.NewMetadata(uuidx.WithSuffix(1), library.MetadataOptionEncryptionSeed(seed.String()), library.MetadataOptionBytes(nlen))
		src, err := cryptox.NewReaderChaCha20(library.MetadataChaCha8(md), io.TeeReader(cryptox.NewChaCha8(seed.String()), digest))
		require.NoError(t, err)

		routes := mux.NewRouter()
		routes.Handle("/m/{id}/download", alice.New().ThenFunc(download(src, nlen)))

		srv := httptest.NewServer(routes)
		defer srv.Close()

		c := &http.Client{}
		c.Transport = httpx.RewriteHostTransport(testx.Must(url.ParseRequestURI(srv.URL))(t), c.Transport)

		reader := library.NewDeeppoolReaderAt(c, md, testx.Must(blockcache.NewDirectoryCache(t.TempDir()))(t))

		downloaded := md5.New()
		n, err := io.CopyN(downloaded, io.NewSectionReader(reader, 0, nlen), nlen)
		require.NoError(t, err)
		require.EqualValues(t, nlen, n)
		require.Equal(t, md5x.FormatUUID(digest), md5x.FormatUUID(downloaded))
	})

	t.Run("range read the 2nd 16 KiB middle block of the data", func(t *testing.T) {
		const (
			nlen = 128 * bytesx.MiB
			clen = 33*bytesx.MiB + 16*bytesx.KiB
		)

		var ()

		md := library.NewMetadata(uuidx.WithSuffix(1), library.MetadataOptionEncryptionSeed(md5x.String(t.Name())))
		src, err := cryptox.NewReaderChaCha20(library.MetadataChaCha8(md), cryptox.NewChaCha8(t.Name()))
		require.NoError(t, err)

		routes := mux.NewRouter()
		routes.Handle("/m/{id}/download", alice.New().ThenFunc(download(src, nlen)))

		srv := httptest.NewServer(routes)
		defer srv.Close()

		c := &http.Client{}
		c.Transport = httpx.RewriteHostTransport(testx.Must(url.ParseRequestURI(srv.URL))(t), c.Transport)

		var (
			buf    bytes.Buffer
			digest = md5.New()
		)

		_n, err := io.Copy(&buf, io.LimitReader(cryptox.NewChaCha8(t.Name()), nlen))
		require.NoError(t, err)
		require.EqualValues(t, nlen, _n)
		require.EqualValues(t, nlen, len(buf.Bytes()))

		_n, err = io.Copy(digest, io.NewSectionReader(bytes.NewReader(buf.Bytes()), clen, 16*bytesx.KiB))
		require.NoError(t, err)
		require.EqualValues(t, 16*bytesx.KiB, _n)

		reader := library.NewDeeppoolReaderAt(c, md, testx.Must(blockcache.NewDirectoryCache(t.TempDir()))(t))

		downloaded := md5.New()
		n, err := io.Copy(downloaded, io.NewSectionReader(reader, clen, 16*bytesx.KiB))
		require.NoError(t, err)
		require.EqualValues(t, 16*bytesx.KiB, n)
		require.Equal(t, md5x.FormatUUID(digest), md5x.FormatUUID(downloaded))
	})

	t.Run("full download with non-zero DiskOffset via blockcache.File", func(t *testing.T) {
		// Regression test for double-offset bug: blockcache.File.ReadAt adds DiskOffset
		// before calling DeeppoolReaderAtCache.ReadAt, which then adds it again when
		// computing the archive download offset. With DiskOffset=0 this is silent;
		// with any non-zero value the wrong archive range is fetched and/or the
		// plaintext is read back from the wrong position in local storage.
		const (
			nlen       = 4 * bytesx.KiB
			diskOffset = 512
		)

		var (
			digest = md5.New()
			seed   = errorsx.Must(uuid.NewV4())
		)

		md := library.NewMetadata(
			uuidx.WithSuffix(1),
			library.MetadataOptionEncryptionSeed(seed.String()),
			library.MetadataOptionBytes(nlen),
			library.MetadataOptionOffset(diskOffset),
		)

		src, err := cryptox.NewReaderChaCha20(library.MetadataChaCha8(md), io.TeeReader(cryptox.NewChaCha8(seed.String()), digest))
		require.NoError(t, err)

		routes := mux.NewRouter()
		routes.Handle("/m/{id}/download", alice.New().ThenFunc(download(src, nlen)))

		srv := httptest.NewServer(routes)
		defer srv.Close()

		c := &http.Client{}
		c.Transport = httpx.RewriteHostTransport(testx.Must(url.ParseRequestURI(srv.URL))(t), c.Transport)

		vfs := library.New(c, fsx.DirVirtual(t.TempDir()), func(ctx context.Context, s string) (*library.Metadata, error) {
			return &md, nil
		})

		file, err := vfs.Open(md.ID)
		require.NoError(t, err)
		defer file.Close()

		downloaded := md5.New()
		n, err := io.CopyN(downloaded, file, nlen)
		require.NoError(t, err)
		require.EqualValues(t, nlen, n)
		require.Equal(t, md5x.FormatUUID(digest), md5x.FormatUUID(downloaded))
	})

	t.Run("range read the 2nd 16 KiB block of the data", func(t *testing.T) {
		const (
			nlen = 128 * bytesx.KiB
			clen = 16 * bytesx.KiB
		)

		var ()

		md := library.NewMetadata(uuidx.WithSuffix(1), library.MetadataOptionEncryptionSeed(md5x.String(t.Name())))
		src, err := cryptox.NewReaderChaCha20(library.MetadataChaCha8(md), cryptox.NewChaCha8(t.Name()))
		require.NoError(t, err)

		routes := mux.NewRouter()
		routes.Handle("/m/{id}/download", alice.New().ThenFunc(download(src, nlen)))

		srv := httptest.NewServer(routes)
		defer srv.Close()

		c := &http.Client{}
		c.Transport = httpx.RewriteHostTransport(testx.Must(url.ParseRequestURI(srv.URL))(t), c.Transport)

		var (
			buf    bytes.Buffer
			digest = md5.New()
		)

		_n, err := io.Copy(&buf, io.LimitReader(cryptox.NewChaCha8(t.Name()), nlen))
		require.NoError(t, err)
		require.EqualValues(t, nlen, _n)
		require.EqualValues(t, nlen, len(buf.Bytes()))

		_n, err = io.Copy(digest, io.NewSectionReader(bytes.NewReader(buf.Bytes()), clen, clen))
		require.NoError(t, err)
		require.EqualValues(t, clen, _n)

		reader := library.NewDeeppoolReaderAt(c, md, testx.Must(blockcache.NewDirectoryCache(t.TempDir()))(t))

		downloaded := md5.New()
		n, err := io.Copy(downloaded, io.NewSectionReader(reader, clen, clen))
		require.NoError(t, err)
		require.EqualValues(t, clen, n)
		require.Equal(t, md5x.FormatUUID(digest), md5x.FormatUUID(downloaded))
	})

	t.Run("within-block cache hit - second read needs no network", func(t *testing.T) {
		// After a block is downloaded, subsequent reads within that block are served
		// from local storage without touching the network.

		const nlen = 128 * bytesx.KiB // fits in one 32 MiB block

		md := library.NewMetadata(uuidx.WithSuffix(1), library.MetadataOptionEncryptionSeed(md5x.String(t.Name())), library.MetadataOptionBytes(nlen))
		src, err := cryptox.NewReaderChaCha20(library.MetadataChaCha8(md), cryptox.NewChaCha8(t.Name()))
		require.NoError(t, err)

		routes := mux.NewRouter()
		routes.Handle("/m/{id}/download", alice.New().ThenFunc(download(src, nlen)))

		srv := httptest.NewServer(routes)
		defer srv.Close()

		c := &http.Client{}
		c.Transport = httpx.RewriteHostTransport(testx.Must(url.ParseRequestURI(srv.URL))(t), c.Transport)

		dcache := testx.Must(blockcache.NewDirectoryCache(t.TempDir()))(t)
		reader := library.NewDeeppoolReaderAt(c, md, dcache)

		// first read: triggers the download and populates the cache
		first := md5.New()
		n, err := io.CopyN(first, io.NewSectionReader(reader, 0, nlen/2), nlen/2)
		require.NoError(t, err)
		require.EqualValues(t, nlen/2, n)

		// stop the server — any further network access would cause an error
		srv.Close()

		// second read: different range within the same cached block, no network needed
		second := md5.New()
		n, err = io.CopyN(second, io.NewSectionReader(reader, nlen/2, nlen/2), nlen/2)
		require.NoError(t, err)
		require.EqualValues(t, nlen/2, n)

		// both halves together should equal the full plaintext
		full := md5.New()
		_, err = io.Copy(full, io.LimitReader(cryptox.NewChaCha8(t.Name()), nlen))
		require.NoError(t, err)

		combined := md5.New()
		plaintext := make([]byte, nlen)
		_, err = io.ReadFull(cryptox.NewChaCha8(t.Name()), plaintext)
		require.NoError(t, err)
		_, err = combined.Write(plaintext[:nlen/2])
		require.NoError(t, err)
		_, err = combined.Write(plaintext[nlen/2:])
		require.NoError(t, err)

		require.Equal(t, md5x.FormatUUID(combined), md5x.FormatUUID(full))
	})

	t.Run("download failure propagates error", func(t *testing.T) {
		// A read against an unreachable server (simulating download failure) must
		// propagate the error rather than returning zeroed or stale data.
		const nlen = 128 * bytesx.KiB

		md := library.NewMetadata(uuidx.WithSuffix(1), library.MetadataOptionEncryptionSeed(md5x.String(t.Name())), library.MetadataOptionBytes(nlen))

		routes := mux.NewRouter()
		routes.Handle("/m/{id}/download", alice.New().ThenFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}))

		srv := httptest.NewServer(routes)
		defer srv.Close()

		c := &http.Client{}
		c.Transport = httpx.RewriteHostTransport(testx.Must(url.ParseRequestURI(srv.URL))(t), c.Transport)

		reader := library.NewDeeppoolReaderAt(c, md, testx.Must(blockcache.NewDirectoryCache(t.TempDir()))(t))

		buf := make([]byte, 1024)
		_, err := reader.ReadAt(buf, 0)
		require.Error(t, err)
	})

	t.Run("cache hit on first attempt skips download", func(t *testing.T) {
		// A cached read (local storage hit on first attempt) must not reach the network.

		const nlen = 128 * bytesx.KiB

		md := library.NewMetadata(uuidx.WithSuffix(1), library.MetadataOptionEncryptionSeed(md5x.String(t.Name())), library.MetadataOptionBytes(nlen))
		src, err := cryptox.NewReaderChaCha20(library.MetadataChaCha8(md), cryptox.NewChaCha8(t.Name()))
		require.NoError(t, err)

		routes := mux.NewRouter()
		routes.Handle("/m/{id}/download", alice.New().ThenFunc(download(src, nlen)))

		srv := httptest.NewServer(routes)
		defer srv.Close()

		c := &http.Client{}
		c.Transport = httpx.RewriteHostTransport(testx.Must(url.ParseRequestURI(srv.URL))(t), c.Transport)

		dcache := testx.Must(blockcache.NewDirectoryCache(t.TempDir()))(t)
		reader := library.NewDeeppoolReaderAt(c, md, dcache)

		// warm the cache
		digest := md5.New()
		n, err := io.CopyN(digest, io.NewSectionReader(reader, 0, nlen), nlen)
		require.NoError(t, err)
		require.EqualValues(t, nlen, n)

		// stop server: proves subsequent read is cache-only
		srv.Close()

		// re-read the exact same range
		again := md5.New()
		n, err = io.CopyN(again, io.NewSectionReader(reader, 0, nlen), nlen)
		require.NoError(t, err)
		require.EqualValues(t, nlen, n)
		require.Equal(t, md5x.FormatUUID(digest), md5x.FormatUUID(again))
	})
}
