package metaapi_test

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/bytesx"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/httptestx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/jwtx"
	"github.com/retrovibed/retrovibed/shallows/internal/md5x"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/testx"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"github.com/stretchr/testify/require"
)

func TestHTTPWireguardCreateUnparsable(t *testing.T) {
	var (
		claims jwt.RegisteredClaims
		buf    io.ReadCloser
		d      = md5.New()
	)
	ctx, done := testx.Context(t)
	defer done()

	q := sqltestx.Metadatabase(t)

	tmpdir := fsx.DirVirtual(t.TempDir())

	routes := mux.NewRouter()

	metaapi.NewHTTPWireguard(
		tmpdir.Path(),
		q,
		metaapi.HTTPWireguardOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
	).Bind(routes.PathPrefix("/").Subrouter())

	mimetype, buf, err := httpx.Multipart(func(w *multipart.Writer) error {
		part, lerr := w.CreatePart(httpx.NewMultipartHeader("application/octet-stream", "content", "example.bin"))
		if lerr != nil {
			return errorsx.Wrap(lerr, "unable to create archive part")
		}

		if _, lerr = io.Copy(part, io.TeeReader(io.LimitReader(rand.Reader, 16*bytesx.KiB), d)); lerr != nil {
			return errorsx.Wrap(lerr, "unable to copy archive")
		}

		return nil
	})
	require.NoError(t, err)

	claims = jwtx.NewJWTClaims(uuid.Nil.String(), jwtx.ClaimsOptionAuthnExpiration())

	resp, req, err := httptestx.BuildRequestContextBytes(
		ctx,
		http.MethodPost,
		"/",
		testx.IOBytes(buf),
		httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(&claims, httpauthtest.UnsafeJWTSecretSource)),
		httptestx.RequestOptionHeader("Content-Type", mimetype),
	)
	require.NoError(t, err)

	routes.ServeHTTP(resp, req)

	require.Error(t, httpx.ErrorCode(resp.Result()))
	require.Equal(t, 400, resp.Code)
}

func TestHTTPWireguardCreate(t *testing.T) {
	var (
		result metaapi.WireguardUploadResponse
		claims jwt.RegisteredClaims
		buf    io.ReadCloser
		d      = md5.New()
	)
	ctx, done := testx.Context(t)
	defer done()

	q := sqltestx.Metadatabase(t)

	tmpdir := fsx.DirVirtual(t.TempDir())

	routes := mux.NewRouter()

	metaapi.NewHTTPWireguard(
		tmpdir.Path(),
		q,
		metaapi.HTTPWireguardOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
	).Bind(routes.PathPrefix("/").Subrouter())

	mimetype, buf, err := httpx.Multipart(func(w *multipart.Writer) error {
		part, lerr := w.CreatePart(httpx.NewMultipartHeader("application/octet-stream", "content", "example.bin"))
		if lerr != nil {
			return errorsx.Wrap(lerr, "unable to create archive part")
		}

		if _, lerr = io.Copy(part, io.TeeReader(testx.Read(testx.Fixture("wireguard", "example.1.conf")), d)); lerr != nil {
			return errorsx.Wrap(lerr, "unable to copy archive")
		}

		return nil
	})
	require.NoError(t, err)

	claims = jwtx.NewJWTClaims(uuid.Nil.String(), jwtx.ClaimsOptionAuthnExpiration())

	resp, req, err := httptestx.BuildRequestContextBytes(
		ctx,
		http.MethodPost,
		"/",
		testx.IOBytes(buf),
		httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(&claims, httpauthtest.UnsafeJWTSecretSource)),
		httptestx.RequestOptionHeader("Content-Type", mimetype),
	)
	require.NoError(t, err)

	require.Equal(t, 0, sqltestx.Count(t, q, "SELECT COUNT(*) FROM meta_wireguard"))

	routes.ServeHTTP(resp, req)

	require.Equal(t, 1, sqltestx.Count(t, q, "SELECT COUNT(*) FROM meta_wireguard"))

	require.NoError(t, httpx.ErrorCode(resp.Result()))
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))

	require.Equal(t, md5x.FormatUUID(d), result.Wireguard.Id)
	require.Equal(t, "example.bin", result.Wireguard.Description)
	require.FileExists(t, tmpdir.Path(md5x.FormatUUID(d)))
}
