package media_test

import (
	"bytes"
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/mux"
	"github.com/retrovibed/retroapi/blockcache"
	"github.com/retrovibed/retroapi/deeppool"
	"github.com/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/asyncx"
	"github.com/retrovibed/retrovibed/shallows/internal/cryptox"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/httptestx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/md5x"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/testx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/retrovibed/retrovibed/shallows/media"
	"github.com/retrovibed/retrovibed/shallows/meta"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"github.com/stretchr/testify/require"
)

// deeppoolStore acts as both the archiver.Upload target and the HTTP download server.
// Archive writes ChaCha20-encrypted bytes into the store; serveDownload replays them
// so the download path can decrypt and return the original plaintext.
type deeppoolStore struct {
	mu    sync.Mutex
	store map[string][]byte // archiveID → encrypted bytes
}

func (m *deeppoolStore) Upload(_ context.Context, _ string, r io.Reader) (*deeppool.Media, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	id := uuid.Must(uuid.NewV4()).String()
	m.mu.Lock()
	m.store[id] = data
	m.mu.Unlock()
	return &deeppool.Media{
		Id:    id,
		Bytes: uint64(len(data)),
		Usage: uint64(len(data)),
	}, nil
}

func (m *deeppoolStore) serveDownload(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	m.mu.Lock()
	data := m.store[id]
	m.mu.Unlock()
	if data == nil {
		http.NotFound(w, r)
		return
	}
	http.ServeContent(w, r, "content.bin", time.Time{}, io.NewSectionReader(bytes.NewReader(data), 0, int64(len(data))))
}

func TestLibraryDownload(t *testing.T) {
	t.Run("serve from local storage", func(t *testing.T) {
		var (
			p  meta.Profile
			v  meta.Authz
			md library.Metadata
		)
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, testx.Fake(&v, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, v).Scan(&v))

		require.NoError(t, testx.Fake(&md, library.MetadataOptionTestDefaults, library.MetadataOptionTestRandomID))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, md).Scan(&md))

		storage := fsx.DirVirtual(t.TempDir())

		// Write plaintext directly into the local block cache for this media ID.
		expected := md5.New()
		dcache, err := blockcache.NewDirectoryCache(storage.Path(md.ID))
		require.NoError(t, err)
		_, err = io.Copy(io.MultiWriter(io.NewOffsetWriter(dcache, 0), expected), io.LimitReader(cryptox.NewChaCha8(t.Name()), int64(md.Bytes)))
		require.NoError(t, err)

		routes := mux.NewRouter()
		media.NewHTTPLibrary(
			q,
			asyncx.NewWakeup(t.Context()),
			storage,
			nil,
			media.HTTPLibraryOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodGet,
			fmt.Sprintf("/%s", md.ID),
			nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.Equal(t, http.StatusOK, resp.Result().StatusCode)
		require.Equal(t, md5x.FormatUUID(expected), testx.IOMD5(resp.Body))
	})

	// deeppool stores data encrypted: ChaCha20(plaintext). When the local cache misses,
	// the handler downloads the encrypted bytes and XORs them again with the same keystream,
	// yielding the plaintext on disk. The HTTP response then serves the plaintext.
	t.Run("deeppool - explicit encryption seed", func(t *testing.T) {
		var (
			p         meta.Profile
			v         meta.Authz
			md        library.Metadata
			archiveID = errorsx.Must(uuid.NewV4())
			seed      = errorsx.Must(uuid.NewV4())
		)
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, testx.Fake(&v, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, v).Scan(&v))

		require.NoError(t, testx.Fake(&md, library.MetadataOptionTestDefaults, library.MetadataOptionTestRandomID,
			library.MetadataOptionArchiveID(archiveID.String()),
			library.MetadataOptionEncryptionSeed(seed.String()),
		))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, md).Scan(&md))

		storage := fsx.DirVirtual(t.TempDir())

		// Build plaintext and capture its MD5.
		plaintextBuf := new(bytes.Buffer)
		expected := md5.New()
		_, err := io.Copy(io.MultiWriter(plaintextBuf, expected), io.LimitReader(cryptox.NewChaCha8(t.Name()), int64(md.Bytes)))
		require.NoError(t, err)

		// Deeppool stores the plaintext XOR'd with the ChaCha20 keystream.
		// The download path XORs again, cancelling out to produce plaintext on disk.
		encryptedReader, err := cryptox.NewReaderChaCha20(library.MetadataChaCha8(md), bytes.NewReader(plaintextBuf.Bytes()))
		require.NoError(t, err)
		encryptedBuf := new(bytes.Buffer)
		_, err = io.Copy(encryptedBuf, encryptedReader)
		require.NoError(t, err)

		deeppoolRoutes := mux.NewRouter()
		deeppoolRoutes.HandleFunc("/m/{id}/download", func(w http.ResponseWriter, r *http.Request) {
			http.ServeContent(w, r, "content.bin", time.Time{}, io.NewSectionReader(bytes.NewReader(encryptedBuf.Bytes()), 0, int64(md.Bytes)))
		})
		srv := httptest.NewServer(deeppoolRoutes)
		defer srv.Close()

		c := &http.Client{}
		c.Transport = httpx.RewriteHostTransport(errorsx.Must(url.ParseRequestURI(srv.URL)), c.Transport)

		routes := mux.NewRouter()
		media.NewHTTPLibrary(
			q,
			asyncx.NewWakeup(t.Context()),
			storage,
			c,
			media.HTTPLibraryOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodGet,
			fmt.Sprintf("/%s", md.ID),
			nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.Equal(t, http.StatusOK, resp.Result().StatusCode)
		require.Equal(t, md5x.FormatUUID(expected), testx.IOMD5(resp.Body))
	})

	// When EncryptionSeed is uuid.Nil, MetadataChaCha8 falls back to md.ID as the key.
	t.Run("deeppool - id used as encryption key when seed is nil uuid", func(t *testing.T) {
		var (
			p         meta.Profile
			v         meta.Authz
			md        library.Metadata
			archiveID = errorsx.Must(uuid.NewV4())
		)
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, testx.Fake(&v, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, v).Scan(&v))

		// EncryptionSeed = uuid.Nil causes MetadataChaCha8 to use md.ID as the PRNG seed.
		require.NoError(t, testx.Fake(&md, library.MetadataOptionTestDefaults, library.MetadataOptionTestRandomID,
			library.MetadataOptionArchiveID(archiveID.String()),
			library.MetadataOptionEncryptionSeed(uuid.Nil.String()),
		))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, md).Scan(&md))

		storage := fsx.DirVirtual(t.TempDir())

		plaintextBuf := new(bytes.Buffer)
		expected := md5.New()
		_, err := io.Copy(io.MultiWriter(plaintextBuf, expected), io.LimitReader(cryptox.NewChaCha8(t.Name()), int64(md.Bytes)))
		require.NoError(t, err)

		encryptedReader, err := cryptox.NewReaderChaCha20(library.MetadataChaCha8(md), bytes.NewReader(plaintextBuf.Bytes()))
		require.NoError(t, err)
		encryptedBuf := new(bytes.Buffer)
		_, err = io.Copy(encryptedBuf, encryptedReader)
		require.NoError(t, err)

		deeppoolRoutes := mux.NewRouter()
		deeppoolRoutes.HandleFunc("/m/{id}/download", func(w http.ResponseWriter, r *http.Request) {
			http.ServeContent(w, r, "content.bin", time.Time{}, io.NewSectionReader(bytes.NewReader(encryptedBuf.Bytes()), 0, int64(md.Bytes)))
		})
		srv := httptest.NewServer(deeppoolRoutes)
		defer srv.Close()

		c := &http.Client{}
		c.Transport = httpx.RewriteHostTransport(errorsx.Must(url.ParseRequestURI(srv.URL)), c.Transport)

		routes := mux.NewRouter()
		media.NewHTTPLibrary(
			q,
			asyncx.NewWakeup(t.Context()),
			storage,
			c,
			media.HTTPLibraryOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodGet,
			fmt.Sprintf("/%s", md.ID),
			nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.Equal(t, http.StatusOK, resp.Result().StatusCode)
		require.Equal(t, md5x.FormatUUID(expected), testx.IOMD5(resp.Body))
	})

	// End-to-end: Archive encrypts local plaintext and uploads to deeppool.
	// After evicting the local cache, the download path fetches the encrypted
	// bytes and XORs with the same ChaCha20 keystream, recovering the plaintext.
	t.Run("archive then serve from deeppool", func(t *testing.T) {
		var (
			p  meta.Profile
			v  meta.Authz
			md library.Metadata
		)
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, testx.Fake(&v, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, v).Scan(&v))

		require.NoError(t, testx.Fake(&md, library.MetadataOptionTestDefaults, library.MetadataOptionTestRandomID))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, md).Scan(&md))

		storage := fsx.DirVirtual(t.TempDir())

		// Write plaintext into local storage and record its MD5.
		expected := md5.New()
		dcache, err := blockcache.NewDirectoryCache(storage.Path(md.ID))
		require.NoError(t, err)
		_, err = io.Copy(io.MultiWriter(io.NewOffsetWriter(dcache, 0), expected), io.LimitReader(cryptox.NewChaCha8(t.Name()), int64(md.Bytes)))
		require.NoError(t, err)

		// Archive: reads plaintext, encrypts with ChaCha20, stores in mock deeppool.
		mock := &deeppoolStore{store: make(map[string][]byte)}
		require.NoError(t, library.Archive(ctx, q, &md, storage, mock))
		require.NotEmpty(t, md.ArchiveID)

		// Evict local cache to force a deeppool pull on the next read.
		require.NoError(t, os.RemoveAll(storage.Path(md.ID)))

		// Serve the captured encrypted bytes at the deeppool download endpoint.
		deeppoolRoutes := mux.NewRouter()
		deeppoolRoutes.HandleFunc("/m/{id}/download", mock.serveDownload)
		srv := httptest.NewServer(deeppoolRoutes)
		defer srv.Close()

		c := &http.Client{}
		c.Transport = httpx.RewriteHostTransport(errorsx.Must(url.ParseRequestURI(srv.URL)), c.Transport)

		routes := mux.NewRouter()
		media.NewHTTPLibrary(
			q,
			asyncx.NewWakeup(t.Context()),
			storage,
			c,
			media.HTTPLibraryOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodGet,
			fmt.Sprintf("/%s", md.ID),
			nil,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.Equal(t, http.StatusOK, resp.Result().StatusCode)
		require.Equal(t, md5x.FormatUUID(expected), testx.IOMD5(resp.Body))
	})
}
