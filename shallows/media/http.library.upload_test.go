package media_test

import (
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/mux"
	"github.com/retrovibed/retrovibed/retroapi/blockcache"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/asyncx"
	"github.com/retrovibed/retrovibed/shallows/internal/bytesx"
	"github.com/retrovibed/retrovibed/shallows/internal/cryptox"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/httptestx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/mimex"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/testx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/retrovibed/retrovibed/shallows/media"
	"github.com/retrovibed/retrovibed/shallows/meta"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"github.com/stretchr/testify/require"
)

func TestLibraryUploadFile(t *testing.T) {
	t.Run("basic upload", func(t *testing.T) {
		var (
			p      meta.Profile
			v      meta.Authz
			result media.MediaUploadResponse
		)
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, testx.Fake(&v, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, v).Scan(&v))

		vfs := fsx.DirVirtual(t.TempDir())
		routes := mux.NewRouter()

		media.NewHTTPLibrary(
			q,
			asyncx.NewWakeup(t.Context()),
			vfs,
			nil,
			media.HTTPLibraryOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		prng := cryptox.NewChaCha8(errorsx.Must(uuid.NewV4()).String())

		mimetype, buf, err := httpx.Multipart(func(w *multipart.Writer) error {
			part, lerr := w.CreatePart(httpx.NewMultipartHeader("application/octet-stream", "content", "example.bin"))
			if lerr != nil {
				return errorsx.Wrap(lerr, "unable to create archive part")
			}

			if _, lerr = io.Copy(part, io.LimitReader(prng, 16*bytesx.KiB)); lerr != nil {
				return errorsx.Wrap(lerr, "unable to copy archive")
			}

			return nil
		})
		require.NoError(t, err)

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodPost,
			"/",
			testx.IOBytes(buf),
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
			httptestx.RequestOptionHeader("Content-Type", mimetype),
		)
		require.NoError(t, err)

		require.Equal(t, 0, testx.Must(sqlx.Count(ctx, q, "SELECT COUNT(*) FROM library_metadata"))(t))

		routes.ServeHTTP(resp, req)

		require.Equal(t, http.StatusOK, resp.Result().StatusCode)
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))

		var md library.Metadata
		require.NoError(t, library.MetadataFindByID(ctx, q, result.Media.Id).Scan(&md))

		bcache, err := blockcache.NewDirectoryCache(vfs.Path(md.ID))
		require.NoError(t, err)

		require.Equal(t, result.Media.Id, testx.IOMD5(io.NewSectionReader(bcache, 0, int64(md.Bytes))))
	})

	t.Run("simple_filename", func(t *testing.T) {
		var (
			p      meta.Profile
			v      meta.Authz
			result media.MediaUploadResponse
		)
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, testx.Fake(&v, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, v).Scan(&v))

		vfs := fsx.DirVirtual(t.TempDir())
		routes := mux.NewRouter()

		media.NewHTTPLibrary(
			q,
			asyncx.NewWakeup(t.Context()),
			vfs,
			nil,
			media.HTTPLibraryOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		prng := cryptox.NewChaCha8(errorsx.Must(uuid.NewV4()).String())

		mimetype, buf, err := httpx.Multipart(func(w *multipart.Writer) error {
			part, lerr := w.CreatePart(httpx.NewMultipartHeader(mimex.RetrovibedMediaArchive, "content", "example.mp4"))
			if lerr != nil {
				return errorsx.Wrap(lerr, "unable to create archive part")
			}

			if _, lerr = io.Copy(part, io.LimitReader(prng, 16*bytesx.KiB)); lerr != nil {
				return errorsx.Wrap(lerr, "unable to copy archive")
			}

			return nil
		})
		require.NoError(t, err)

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodPost,
			"/",
			testx.IOBytes(buf),
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
			httptestx.RequestOptionHeader("Content-Type", mimetype),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.Equal(t, http.StatusOK, resp.Result().StatusCode)
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))

		var dbMD library.Metadata
		require.NoError(t, library.MetadataFindByID(ctx, q, result.Media.Id).Scan(&dbMD))

		require.Equal(t, "example.mp4", dbMD.Description)
		require.Equal(t, "example mp4", dbMD.AutoDescription)
	})

	t.Run("filename_with_dashes", func(t *testing.T) {
		var (
			p      meta.Profile
			v      meta.Authz
			result media.MediaUploadResponse
		)
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, testx.Fake(&v, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, v).Scan(&v))

		vfs := fsx.DirVirtual(t.TempDir())
		routes := mux.NewRouter()

		media.NewHTTPLibrary(
			q,
			asyncx.NewWakeup(t.Context()),
			vfs,
			nil,
			media.HTTPLibraryOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		prng := cryptox.NewChaCha8(errorsx.Must(uuid.NewV4()).String())

		mimetype, buf, err := httpx.Multipart(func(w *multipart.Writer) error {
			part, lerr := w.CreatePart(httpx.NewMultipartHeader(mimex.RetrovibedMediaArchive, "content", "my-example-video.mp4"))
			if lerr != nil {
				return errorsx.Wrap(lerr, "unable to create archive part")
			}

			if _, lerr = io.Copy(part, io.LimitReader(prng, 16*bytesx.KiB)); lerr != nil {
				return errorsx.Wrap(lerr, "unable to copy archive")
			}

			return nil
		})
		require.NoError(t, err)

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodPost,
			"/",
			testx.IOBytes(buf),
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
			httptestx.RequestOptionHeader("Content-Type", mimetype),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.Equal(t, http.StatusOK, resp.Result().StatusCode)
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))

		var dbMD library.Metadata
		require.NoError(t, library.MetadataFindByID(ctx, q, result.Media.Id).Scan(&dbMD))

		require.Equal(t, "my-example-video.mp4", dbMD.Description)
		require.Equal(t, "my example video mp4", dbMD.AutoDescription)
	})

	t.Run("filename_with_underscores", func(t *testing.T) {
		var (
			p      meta.Profile
			v      meta.Authz
			result media.MediaUploadResponse
		)
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, testx.Fake(&v, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, v).Scan(&v))

		vfs := fsx.DirVirtual(t.TempDir())
		routes := mux.NewRouter()

		media.NewHTTPLibrary(
			q,
			asyncx.NewWakeup(t.Context()),
			vfs,
			nil,
			media.HTTPLibraryOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		prng := cryptox.NewChaCha8(errorsx.Must(uuid.NewV4()).String())

		mimetype, buf, err := httpx.Multipart(func(w *multipart.Writer) error {
			part, lerr := w.CreatePart(httpx.NewMultipartHeader(mimex.RetrovibedMediaArchive, "content", "my_example_video.mp4"))
			if lerr != nil {
				return errorsx.Wrap(lerr, "unable to create archive part")
			}

			if _, lerr = io.Copy(part, io.LimitReader(prng, 16*bytesx.KiB)); lerr != nil {
				return errorsx.Wrap(lerr, "unable to copy archive")
			}

			return nil
		})
		require.NoError(t, err)

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodPost,
			"/",
			testx.IOBytes(buf),
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
			httptestx.RequestOptionHeader("Content-Type", mimetype),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.Equal(t, http.StatusOK, resp.Result().StatusCode)
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))

		var dbMD library.Metadata
		require.NoError(t, library.MetadataFindByID(ctx, q, result.Media.Id).Scan(&dbMD))

		require.Equal(t, "my_example_video.mp4", dbMD.Description)
		require.Equal(t, "my example video mp4", dbMD.AutoDescription)
	})

	t.Run("filename_with_dots", func(t *testing.T) {
		var (
			p      meta.Profile
			v      meta.Authz
			result media.MediaUploadResponse
		)
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, testx.Fake(&v, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, v).Scan(&v))

		vfs := fsx.DirVirtual(t.TempDir())
		routes := mux.NewRouter()

		media.NewHTTPLibrary(
			q,
			asyncx.NewWakeup(t.Context()),
			vfs,
			nil,
			media.HTTPLibraryOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		prng := cryptox.NewChaCha8(errorsx.Must(uuid.NewV4()).String())

		mimetype, buf, err := httpx.Multipart(func(w *multipart.Writer) error {
			part, lerr := w.CreatePart(httpx.NewMultipartHeader(mimex.RetrovibedMediaArchive, "content", "video.file.name.mp4"))
			if lerr != nil {
				return errorsx.Wrap(lerr, "unable to create archive part")
			}

			if _, lerr = io.Copy(part, io.LimitReader(prng, 16*bytesx.KiB)); lerr != nil {
				return errorsx.Wrap(lerr, "unable to copy archive")
			}

			return nil
		})
		require.NoError(t, err)

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodPost,
			"/",
			testx.IOBytes(buf),
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
			httptestx.RequestOptionHeader("Content-Type", mimetype),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.Equal(t, http.StatusOK, resp.Result().StatusCode)
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))

		var dbMD library.Metadata
		require.NoError(t, library.MetadataFindByID(ctx, q, result.Media.Id).Scan(&dbMD))

		require.Equal(t, "video.file.name.mp4", dbMD.Description)
		require.Equal(t, "video file name mp4", dbMD.AutoDescription)
	})

	t.Run("filename_with_spaces", func(t *testing.T) {
		var (
			p      meta.Profile
			v      meta.Authz
			result media.MediaUploadResponse
		)
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, testx.Fake(&v, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, v).Scan(&v))

		vfs := fsx.DirVirtual(t.TempDir())
		routes := mux.NewRouter()

		media.NewHTTPLibrary(
			q,
			asyncx.NewWakeup(t.Context()),
			vfs,
			nil,
			media.HTTPLibraryOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		prng := cryptox.NewChaCha8(errorsx.Must(uuid.NewV4()).String())

		mimetype, buf, err := httpx.Multipart(func(w *multipart.Writer) error {
			part, lerr := w.CreatePart(httpx.NewMultipartHeader(mimex.RetrovibedMediaArchive, "content", "my video file.mp4"))
			if lerr != nil {
				return errorsx.Wrap(lerr, "unable to create archive part")
			}

			if _, lerr = io.Copy(part, io.LimitReader(prng, 16*bytesx.KiB)); lerr != nil {
				return errorsx.Wrap(lerr, "unable to copy archive")
			}

			return nil
		})
		require.NoError(t, err)

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodPost,
			"/",
			testx.IOBytes(buf),
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
			httptestx.RequestOptionHeader("Content-Type", mimetype),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.Equal(t, http.StatusOK, resp.Result().StatusCode)
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))

		var dbMD library.Metadata
		require.NoError(t, library.MetadataFindByID(ctx, q, result.Media.Id).Scan(&dbMD))

		require.Equal(t, "my video file.mp4", dbMD.Description)
		require.Equal(t, "my video file mp4", dbMD.AutoDescription)
	})

	t.Run("filename_with_punctuation", func(t *testing.T) {
		var (
			p      meta.Profile
			v      meta.Authz
			result media.MediaUploadResponse
		)
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, testx.Fake(&v, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, v).Scan(&v))

		vfs := fsx.DirVirtual(t.TempDir())
		routes := mux.NewRouter()

		media.NewHTTPLibrary(
			q,
			asyncx.NewWakeup(t.Context()),
			vfs,
			nil,
			media.HTTPLibraryOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		prng := cryptox.NewChaCha8(errorsx.Must(uuid.NewV4()).String())

		mimetype, buf, err := httpx.Multipart(func(w *multipart.Writer) error {
			part, lerr := w.CreatePart(httpx.NewMultipartHeader(mimex.RetrovibedMediaArchive, "content", "video:part.two.mp4"))
			if lerr != nil {
				return errorsx.Wrap(lerr, "unable to create archive part")
			}

			if _, lerr = io.Copy(part, io.LimitReader(prng, 16*bytesx.KiB)); lerr != nil {
				return errorsx.Wrap(lerr, "unable to copy archive")
			}

			return nil
		})
		require.NoError(t, err)

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodPost,
			"/",
			testx.IOBytes(buf),
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
			httptestx.RequestOptionHeader("Content-Type", mimetype),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.Equal(t, http.StatusOK, resp.Result().StatusCode)
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))

		var dbMD library.Metadata
		require.NoError(t, library.MetadataFindByID(ctx, q, result.Media.Id).Scan(&dbMD))

		require.Equal(t, "video:part.two.mp4", dbMD.Description)
		require.Equal(t, "video part two mp4", dbMD.AutoDescription)
	})

	t.Run("filename_with_numbers", func(t *testing.T) {
		var (
			p      meta.Profile
			v      meta.Authz
			result media.MediaUploadResponse
		)
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, testx.Fake(&v, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, v).Scan(&v))

		vfs := fsx.DirVirtual(t.TempDir())
		routes := mux.NewRouter()

		media.NewHTTPLibrary(
			q,
			asyncx.NewWakeup(t.Context()),
			vfs,
			nil,
			media.HTTPLibraryOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		prng := cryptox.NewChaCha8(errorsx.Must(uuid.NewV4()).String())

		mimetype, buf, err := httpx.Multipart(func(w *multipart.Writer) error {
			part, lerr := w.CreatePart(httpx.NewMultipartHeader(mimex.RetrovibedMediaArchive, "content", "movie-2024.v1.mp4"))
			if lerr != nil {
				return errorsx.Wrap(lerr, "unable to create archive part")
			}

			if _, lerr = io.Copy(part, io.LimitReader(prng, 16*bytesx.KiB)); lerr != nil {
				return errorsx.Wrap(lerr, "unable to copy archive")
			}

			return nil
		})
		require.NoError(t, err)

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodPost,
			"/",
			testx.IOBytes(buf),
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
			httptestx.RequestOptionHeader("Content-Type", mimetype),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.Equal(t, http.StatusOK, resp.Result().StatusCode)
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))

		var dbMD library.Metadata
		require.NoError(t, library.MetadataFindByID(ctx, q, result.Media.Id).Scan(&dbMD))

		require.Equal(t, "movie-2024.v1.mp4", dbMD.Description)
		require.Equal(t, "movie 2024 v1 mp4", dbMD.AutoDescription)
	})

	t.Run("basic database assertions", func(t *testing.T) {
		var (
			p      meta.Profile
			v      meta.Authz
			result media.MediaUploadResponse
		)
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, testx.Fake(&v, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, v).Scan(&v))

		vfs := fsx.DirVirtual(t.TempDir())
		routes := mux.NewRouter()

		media.NewHTTPLibrary(
			q,
			asyncx.NewWakeup(t.Context()),
			vfs,
			nil,
			media.HTTPLibraryOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		prng := cryptox.NewChaCha8(errorsx.Must(uuid.NewV4()).String())

		mimetype, buf, err := httpx.Multipart(func(w *multipart.Writer) error {
			part, lerr := w.CreatePart(httpx.NewMultipartHeader(mimex.RetrovibedMediaArchive, "content", "test_upload.mp4"))
			if lerr != nil {
				return errorsx.Wrap(lerr, "unable to create archive part")
			}

			if _, lerr = io.Copy(part, io.LimitReader(prng, 32*bytesx.KiB)); lerr != nil {
				return errorsx.Wrap(lerr, "unable to copy archive")
			}

			return nil
		})
		require.NoError(t, err)

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodPost,
			"/",
			testx.IOBytes(buf),
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
			httptestx.RequestOptionHeader("Content-Type", mimetype),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.Equal(t, http.StatusOK, resp.Result().StatusCode)
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))

		var md library.Metadata
		require.NoError(t, library.MetadataFindByID(ctx, q, result.Media.Id).Scan(&md))

		require.Equal(t, result.Media.Id, md.ID)
		require.Equal(t, uint64(32*bytesx.KiB), md.Bytes)
		require.Equal(t, mimex.RetrovibedMediaArchive, md.Mimetype)
		require.Equal(t, "test_upload.mp4", md.Description)
		require.Equal(t, "test upload mp4", md.AutoDescription)
		require.Equal(t, uuid.Nil.String(), md.ArchiveID)
		require.Equal(t, uuid.Nil.String(), md.TorrentID)
		require.NotEqual(t, uuid.Nil.String(), md.EncryptionSeed)
		require.Equal(t, timex.Inf(), md.HiddenAt)
		require.Equal(t, timex.Inf(), md.TombstonedAt)
		require.False(t, md.CreatedAt.IsZero())
		require.False(t, md.UpdatedAt.IsZero())

		bcache, err := blockcache.NewDirectoryCache(vfs.Path(md.ID))
		require.NoError(t, err)

		reader := io.NewSectionReader(bcache, 0, int64(md.Bytes))
		size, err := reader.Seek(0, io.SeekEnd)
		require.NoError(t, err)
		require.Equal(t, int64(md.Bytes), size)
	})

	t.Run("duplicate file succeeds", func(t *testing.T) {
		var (
			p       meta.Profile
			v       meta.Authz
			result  media.MediaUploadResponse
			result2 media.MediaUploadResponse
		)
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, testx.Fake(&v, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, v).Scan(&v))

		vfs := fsx.DirVirtual(t.TempDir())
		routes := mux.NewRouter()

		media.NewHTTPLibrary(
			q,
			asyncx.NewWakeup(t.Context()),
			vfs,
			nil,
			media.HTTPLibraryOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		prng := cryptox.NewChaCha8(errorsx.Must(uuid.NewV4()).String())

		mimetype, buf, err := httpx.Multipart(func(w *multipart.Writer) error {
			part, lerr := w.CreatePart(httpx.NewMultipartHeader(mimex.RetrovibedMediaArchive, "content", "duplicate.mp4"))
			if lerr != nil {
				return errorsx.Wrap(lerr, "unable to create archive part")
			}

			if _, lerr = io.Copy(part, io.LimitReader(prng, 16*bytesx.KiB)); lerr != nil {
				return errorsx.Wrap(lerr, "unable to copy archive")
			}

			return nil
		})
		require.NoError(t, err)

		multipartData, err := io.ReadAll(buf)
		require.NoError(t, err)

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))

		resp1, req1, err := httptestx.BuildRequestBytes(
			http.MethodPost,
			"/",
			multipartData,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
			httptestx.RequestOptionHeader("Content-Type", mimetype),
		)
		require.NoError(t, err)
		routes.ServeHTTP(resp1, req1)

		require.Equal(t, http.StatusOK, resp1.Result().StatusCode)
		require.NoError(t, json.NewDecoder(resp1.Body).Decode(&result))

		resp2, req2, err := httptestx.BuildRequestBytes(
			http.MethodPost,
			"/",
			multipartData,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
			httptestx.RequestOptionHeader("Content-Type", mimetype),
		)
		require.NoError(t, err)
		routes.ServeHTTP(resp2, req2)

		require.Equal(t, http.StatusOK, resp2.Result().StatusCode)
		require.NoError(t, json.NewDecoder(resp2.Body).Decode(&result2))
		require.Equal(t, result.Media.Id, result2.Media.Id)

		count := testx.Must(sqlx.Count(ctx, q, "SELECT COUNT(*) FROM library_metadata WHERE id = ?", result.Media.Id))(t)
		require.Equal(t, 1, count)

		var md library.Metadata
		require.NoError(t, library.MetadataFindByID(ctx, q, result.Media.Id).Scan(&md))
		require.Equal(t, "duplicate.mp4", md.Description)

		bcache, err := blockcache.NewDirectoryCache(vfs.Path(md.ID))
		require.NoError(t, err)
		reader := io.NewSectionReader(bcache, 0, int64(md.Bytes))
		size, err := reader.Seek(0, io.SeekEnd)
		require.NoError(t, err)
		require.Equal(t, int64(md.Bytes), size)
	})

	t.Run("tombstoned file reupload", func(t *testing.T) {
		var (
			p      meta.Profile
			v      meta.Authz
			result media.MediaUploadResponse
		)
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, testx.Fake(&v, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, v).Scan(&v))

		vfs := fsx.DirVirtual(t.TempDir())
		routes := mux.NewRouter()

		media.NewHTTPLibrary(
			q,
			asyncx.NewWakeup(t.Context()),
			vfs,
			nil,
			media.HTTPLibraryOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		prng := cryptox.NewChaCha8(errorsx.Must(uuid.NewV4()).String())

		mimetype, buf, err := httpx.Multipart(func(w *multipart.Writer) error {
			part, lerr := w.CreatePart(httpx.NewMultipartHeader(mimex.RetrovibedMediaArchive, "content", "tombstoned.mp4"))
			if lerr != nil {
				return errorsx.Wrap(lerr, "unable to create archive part")
			}

			if _, lerr = io.Copy(part, io.LimitReader(prng, 16*bytesx.KiB)); lerr != nil {
				return errorsx.Wrap(lerr, "unable to copy archive")
			}

			return nil
		})
		require.NoError(t, err)

		multipartData, err := io.ReadAll(buf)
		require.NoError(t, err)

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))

		resp1, req1, err := httptestx.BuildRequestBytes(
			http.MethodPost,
			"/",
			multipartData,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
			httptestx.RequestOptionHeader("Content-Type", mimetype),
		)
		require.NoError(t, err)
		routes.ServeHTTP(resp1, req1)
		require.Equal(t, http.StatusOK, resp1.Result().StatusCode)
		require.NoError(t, json.NewDecoder(resp1.Body).Decode(&result))

		var md library.Metadata
		require.NoError(t, library.MetadataTombstoneByID(ctx, q, result.Media.Id).Scan(&md))
		require.NotEqual(t, timex.Inf(), md.TombstonedAt)

		resp2, req2, err := httptestx.BuildRequestBytes(
			http.MethodPost,
			"/",
			multipartData,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
			httptestx.RequestOptionHeader("Content-Type", mimetype),
		)
		require.NoError(t, err)
		routes.ServeHTTP(resp2, req2)

		require.Equal(t, http.StatusOK, resp2.Result().StatusCode)

		var restored library.Metadata
		require.NoError(t, library.MetadataFindByID(ctx, q, result.Media.Id).Scan(&restored))
		require.Equal(t, result.Media.Id, restored.ID)
		require.Equal(t, uint64(16*bytesx.KiB), restored.Bytes)
		require.Equal(t, mimex.RetrovibedMediaArchive, restored.Mimetype)
		require.Equal(t, "tombstoned.mp4", restored.Description)
		require.Equal(t, timex.Inf(), restored.HiddenAt)
		require.Equal(t, timex.Inf(), restored.TombstonedAt)

		count := testx.Must(sqlx.Count(ctx, q, "SELECT COUNT(*) FROM library_metadata WHERE id = ?", result.Media.Id))(t)
		require.Equal(t, 1, count)

		bcache, err := blockcache.NewDirectoryCache(vfs.Path(restored.ID))
		require.NoError(t, err)

		reader := io.NewSectionReader(bcache, 0, int64(restored.Bytes))
		size, err := reader.Seek(0, io.SeekEnd)
		require.NoError(t, err)
		require.Equal(t, int64(restored.Bytes), size)
	})

	t.Run("minimal content upload", func(t *testing.T) {
		var (
			p meta.Profile
			v meta.Authz
		)
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, testx.Fake(&v, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, v).Scan(&v))

		vfs := fsx.DirVirtual(t.TempDir())
		routes := mux.NewRouter()

		media.NewHTTPLibrary(
			q,
			asyncx.NewWakeup(t.Context()),
			vfs,
			nil,
			media.HTTPLibraryOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		mimetype, buf, err := httpx.Multipart(func(w *multipart.Writer) error {
			part, err := w.CreateFormFile("content", "empty.mp4")
			if err != nil {
				return err
			}
			_, err = io.WriteString(part, "test")
			return err
		})
		require.NoError(t, err)

		multipartData, err := io.ReadAll(buf)
		require.NoError(t, err)

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodPost,
			"/",
			multipartData,
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
			httptestx.RequestOptionHeader("Content-Type", mimetype),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.Equal(t, http.StatusOK, resp.Result().StatusCode)
	})

	t.Run("retrovibed_archive", func(t *testing.T) {
		var (
			p      meta.Profile
			v      meta.Authz
			result media.MediaUploadResponse
		)
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, testx.Fake(&v, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, v).Scan(&v))

		vfs := fsx.DirVirtual(t.TempDir())
		routes := mux.NewRouter()

		media.NewHTTPLibrary(
			q,
			asyncx.NewWakeup(t.Context()),
			vfs,
			nil,
			media.HTTPLibraryOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		prng := cryptox.NewChaCha8(errorsx.Must(uuid.NewV4()).String())

		mimetype, buf, err := httpx.Multipart(func(w *multipart.Writer) error {
			part, lerr := w.CreatePart(httpx.NewMultipartHeader(mimex.RetrovibedMediaArchive, "content", "archive.mp4"))
			if lerr != nil {
				return errorsx.Wrap(lerr, "unable to create archive part")
			}

			if _, lerr = io.Copy(part, io.LimitReader(prng, 16*bytesx.KiB)); lerr != nil {
				return errorsx.Wrap(lerr, "unable to copy archive")
			}

			return nil
		})
		require.NoError(t, err)

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodPost,
			"/",
			testx.IOBytes(buf),
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
			httptestx.RequestOptionHeader("Content-Type", mimetype),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.Equal(t, http.StatusOK, resp.Result().StatusCode)
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))

		var md library.Metadata
		require.NoError(t, library.MetadataFindByID(ctx, q, result.Media.Id).Scan(&md))
		require.Equal(t, mimex.RetrovibedMediaArchive, md.Mimetype)
		require.Equal(t, "archive.mp4", md.Description)
	})

	t.Run("binary", func(t *testing.T) {
		var (
			p      meta.Profile
			v      meta.Authz
			result media.MediaUploadResponse
		)
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, testx.Fake(&v, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, v).Scan(&v))

		vfs := fsx.DirVirtual(t.TempDir())
		routes := mux.NewRouter()

		media.NewHTTPLibrary(
			q,
			asyncx.NewWakeup(t.Context()),
			vfs,
			nil,
			media.HTTPLibraryOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		prng := cryptox.NewChaCha8(errorsx.Must(uuid.NewV4()).String())

		mimetype, buf, err := httpx.Multipart(func(w *multipart.Writer) error {
			part, lerr := w.CreatePart(httpx.NewMultipartHeader(mimex.Binary, "content", "binary.bin"))
			if lerr != nil {
				return errorsx.Wrap(lerr, "unable to create archive part")
			}

			if _, lerr = io.Copy(part, io.LimitReader(prng, 16*bytesx.KiB)); lerr != nil {
				return errorsx.Wrap(lerr, "unable to copy archive")
			}

			return nil
		})
		require.NoError(t, err)

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodPost,
			"/",
			testx.IOBytes(buf),
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
			httptestx.RequestOptionHeader("Content-Type", mimetype),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.Equal(t, http.StatusOK, resp.Result().StatusCode)
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))

		var md library.Metadata
		require.NoError(t, library.MetadataFindByID(ctx, q, result.Media.Id).Scan(&md))
		require.Equal(t, mimex.Binary, md.Mimetype)
		require.Equal(t, "binary.bin", md.Description)
	})

	t.Run("json", func(t *testing.T) {
		var (
			p      meta.Profile
			v      meta.Authz
			result media.MediaUploadResponse
		)
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, testx.Fake(&v, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, v).Scan(&v))

		vfs := fsx.DirVirtual(t.TempDir())
		routes := mux.NewRouter()

		media.NewHTTPLibrary(
			q,
			asyncx.NewWakeup(t.Context()),
			vfs,
			nil,
			media.HTTPLibraryOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		prng := cryptox.NewChaCha8(errorsx.Must(uuid.NewV4()).String())

		mimetype, buf, err := httpx.Multipart(func(w *multipart.Writer) error {
			part, lerr := w.CreatePart(httpx.NewMultipartHeader(mimex.JSON, "content", "data.json"))
			if lerr != nil {
				return errorsx.Wrap(lerr, "unable to create archive part")
			}

			if _, lerr = io.Copy(part, io.LimitReader(prng, 16*bytesx.KiB)); lerr != nil {
				return errorsx.Wrap(lerr, "unable to copy archive")
			}

			return nil
		})
		require.NoError(t, err)

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodPost,
			"/",
			testx.IOBytes(buf),
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
			httptestx.RequestOptionHeader("Content-Type", mimetype),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.Equal(t, http.StatusOK, resp.Result().StatusCode)
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))

		var md library.Metadata
		require.NoError(t, library.MetadataFindByID(ctx, q, result.Media.Id).Scan(&md))
		require.Equal(t, mimex.JSON, md.Mimetype)
		require.Equal(t, "data.json", md.Description)
	})

	t.Run("rss", func(t *testing.T) {
		var (
			p      meta.Profile
			v      meta.Authz
			result media.MediaUploadResponse
		)
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, testx.Fake(&v, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, v).Scan(&v))

		vfs := fsx.DirVirtual(t.TempDir())
		routes := mux.NewRouter()

		media.NewHTTPLibrary(
			q,
			asyncx.NewWakeup(t.Context()),
			vfs,
			nil,
			media.HTTPLibraryOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/").Subrouter())

		prng := cryptox.NewChaCha8(errorsx.Must(uuid.NewV4()).String())

		mimetype, buf, err := httpx.Multipart(func(w *multipart.Writer) error {
			part, lerr := w.CreatePart(httpx.NewMultipartHeader(mimex.RSS, "content", "feed.rss"))
			if lerr != nil {
				return errorsx.Wrap(lerr, "unable to create archive part")
			}

			if _, lerr = io.Copy(part, io.LimitReader(prng, 16*bytesx.KiB)); lerr != nil {
				return errorsx.Wrap(lerr, "unable to copy archive")
			}

			return nil
		})
		require.NoError(t, err)

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))

		resp, req, err := httptestx.BuildRequestBytes(
			http.MethodPost,
			"/",
			testx.IOBytes(buf),
			httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
			httptestx.RequestOptionHeader("Content-Type", mimetype),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)

		require.Equal(t, http.StatusOK, resp.Result().StatusCode)
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))

		var md library.Metadata
		require.NoError(t, library.MetadataFindByID(ctx, q, result.Media.Id).Scan(&md))
		require.Equal(t, mimex.RSS, md.Mimetype)
		require.Equal(t, "feed.rss", md.Description)
	})
}
