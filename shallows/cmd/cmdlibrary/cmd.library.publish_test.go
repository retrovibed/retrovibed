package cmdlibrary

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/mux"
	"github.com/retrovibed/retrovibed/retroapi/jsonx"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/community"
	"github.com/retrovibed/retrovibed/shallows/communityapi"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/jsonl"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/retrovibed/retrovibed/shallows/meta"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"github.com/stretchr/testify/require"
)

func communityLibraryPublishServer(t *testing.T, q *sql.DB) *mux.Router {
	t.Helper()

	routes := mux.NewRouter()
	communityapi.NewHTTPPublished(
		q,
		communityapi.HTTPPublishedOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		communityapi.HTTPPublishedOptionHTTPClient(&http.Client{}),
		communityapi.HTTPPublishedOptionMediaStorage(fsx.DirVirtual(t.TempDir())),
		communityapi.HTTPPublishedOptionTorrentStorage(fsx.DirVirtual(t.TempDir())),
	).Bind(routes.PathPrefix("/c/p").Subrouter())

	return routes
}

func TestCommunityLibraryPublish(t *testing.T) {
	t.Run("dry run outputs library items without calling endpoint", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)
		communityID := uuid.Must(uuid.NewV7()).String()
		libraryID := uuid.Must(uuid.NewV7()).String()

		lmd := library.Metadata{
			ID:             libraryID,
			Description:    "test media",
			Bytes:          1024,
			TorrentID:      uuid.Nil.String(),
			KnownMediaID:   uuid.Nil.String(),
			ArchiveID:      uuid.Nil.String(),
			DirectoryID:    uuid.Nil.String(),
			EncryptionSeed: uuid.Must(uuid.NewV4()).String(),
			Mimetype:       mimex.RetrovibedMediaArchive,
		}
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, lmd).Scan(&lmd))

		called := false
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
		}))
		defer srv.Close()

		c := &http.Client{
			Transport: httpx.RewriteHostTransport(testx.Must(url.ParseRequestURI(srv.URL))(t), nil),
		}

		com := &communityapi.Community{Id: communityID}
		var input bytes.Buffer
		require.NoError(t, json.NewEncoder(&input).Encode(com))
		require.NoError(t, jsonl.NewEncoder(&input).Encode(langx.Clone(lmd, timex.JSONSafeEncodeOption)))

		var output bytes.Buffer
		cmd := cmdPublish{DryRun: true}
		require.NoError(t, cmd.run(ctx, "localhost:9998", jsonl.NewEncoder(&output), &input, c))

		require.False(t, called, "endpoint must not be called during dry run")
		var decoded library.Metadata
		require.NoError(t, jsonx.UnmarshalRead(&output, &decoded))
		require.Equal(t, libraryID, decoded.ID)
		require.Equal(t, mimex.RetrovibedMediaArchive, decoded.Mimetype)
	})

	t.Run("publishes library item to community", func(t *testing.T) {
		var (
			p meta.Profile
			v meta.Authz
		)
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)
		communityID := uuid.Must(uuid.NewV7()).String()
		libraryID := uuid.Must(uuid.NewV7()).String()

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, testx.Fake(&v, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, v).Scan(&v))

		lmd := library.Metadata{
			ID:             libraryID,
			Description:    "test media",
			Bytes:          1024,
			TorrentID:      uuid.Nil.String(),
			KnownMediaID:   uuid.Nil.String(),
			ArchiveID:      uuid.Nil.String(),
			DirectoryID:    uuid.Nil.String(),
			EncryptionSeed: uuid.Must(uuid.NewV4()).String(),
			Mimetype:       mimex.RetrovibedMediaArchive,
		}
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, lmd).Scan(&lmd))

		srv := httptest.NewServer(communityLibraryPublishServer(t, q))
		defer srv.Close()

		com := &communityapi.Community{Id: communityID}
		var input bytes.Buffer
		require.NoError(t, json.NewEncoder(&input).Encode(com))
		require.NoError(t, jsonl.NewEncoder(&input).Encode(langx.Clone(lmd, timex.JSONSafeEncodeOption)))

		var output bytes.Buffer
		cmd := cmdPublish{DryRun: false}
		token := httpauthtest.UnsafeToken(metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v))), httpauthtest.UnsafeJWTSecretSource)
		headers := http.Header{"Authorization": []string{"Bearer " + token}}
		c := &http.Client{
			Transport: httpx.NewHeadersTransport(headers),
		}
		require.NoError(t, cmd.run(ctx, srv.URL, jsonl.NewEncoder(&output), &input, c))

		var result communityapi.PublishedContent
		require.NoError(t, jsonx.UnmarshalRead(&output, &result))
		require.NotEmpty(t, result.Id)
		require.Equal(t, libraryID, result.LibraryId)
		require.Equal(t, communityID, result.CommunityId)
		require.Equal(t, mimex.RetrovibedMediaArchive, result.Mimetype)

		var pc community.PublishedContent
		require.NoError(t, community.PublishedContentFindByID(ctx, q, result.Id).Scan(&pc))
		require.Equal(t, libraryID, pc.LibraryID)
		require.Equal(t, communityID, pc.CommunityID)
		require.Equal(t, mimex.RetrovibedMediaArchive, pc.Mimetype)
	})

	t.Run("publishes multiple library items", func(t *testing.T) {
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

		communityID := uuid.Must(uuid.NewV7()).String()

		com := &communityapi.Community{Id: communityID}
		var input bytes.Buffer
		require.NoError(t, json.NewEncoder(&input).Encode(com))
		enc := jsonl.NewEncoder(&input)
		var libraryIDs []string
		for range 3 {
			libraryID := uuid.Must(uuid.NewV7()).String()
			libraryIDs = append(libraryIDs, libraryID)
			lmd := library.Metadata{
				ID:             libraryID,
				Description:    "test media",
				Bytes:          1024,
				TorrentID:      uuid.Nil.String(),
				KnownMediaID:   uuid.Nil.String(),
				ArchiveID:      uuid.Nil.String(),
				DirectoryID:    uuid.Nil.String(),
				EncryptionSeed: uuid.Must(uuid.NewV4()).String(),
				Mimetype:       mimex.RetrovibedMediaArchive,
			}
			require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, lmd).Scan(&lmd))
			require.NoError(t, enc.Encode(langx.Clone(lmd, timex.JSONSafeEncodeOption)))
		}

		srv := httptest.NewServer(communityLibraryPublishServer(t, q))
		defer srv.Close()

		var output bytes.Buffer
		cmd := cmdPublish{DryRun: false}
		token := httpauthtest.UnsafeToken(metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v))), httpauthtest.UnsafeJWTSecretSource)
		headers := http.Header{"Authorization": []string{"Bearer " + token}}
		c := &http.Client{
			Transport: httpx.NewHeadersTransport(headers),
		}

		require.NoError(t, cmd.run(ctx, srv.URL, jsonl.NewEncoder(&output), &input, c))

		dec := json.NewDecoder(&output)
		var results []*communityapi.PublishedContent
		for dec.More() {
			var pc communityapi.PublishedContent
			require.NoError(t, dec.Decode(&pc))
			results = append(results, &pc)
		}
		require.Len(t, results, 3)
		for i, result := range results {
			require.Equal(t, libraryIDs[i], result.LibraryId)
			require.Equal(t, communityID, result.CommunityId)
			require.Equal(t, mimex.RetrovibedMediaArchive, result.Mimetype)
		}
	})

	t.Run("uses community default publish mode", func(t *testing.T) {
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

		communityID := uuid.Must(uuid.NewV7()).String()
		libraryID := uuid.Must(uuid.NewV7()).String()

		lmd := library.Metadata{
			ID:             libraryID,
			Description:    "test media",
			Bytes:          1024,
			TorrentID:      uuid.Nil.String(),
			KnownMediaID:   uuid.Nil.String(),
			ArchiveID:      uuid.Nil.String(),
			DirectoryID:    uuid.Nil.String(),
			EncryptionSeed: uuid.Must(uuid.NewV4()).String(),
		}
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, lmd).Scan(&lmd))

		srv := httptest.NewServer(communityLibraryPublishServer(t, q))
		defer srv.Close()

		com := &communityapi.Community{Id: communityID, DefaultPublishMode: communityapi.PublishMode_LISTED}
		var input bytes.Buffer
		require.NoError(t, json.NewEncoder(&input).Encode(com))
		require.NoError(t, jsonl.NewEncoder(&input).Encode(langx.Clone(lmd, timex.JSONSafeEncodeOption)))

		var output bytes.Buffer
		cmd := cmdPublish{DryRun: false}
		token := httpauthtest.UnsafeToken(metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v))), httpauthtest.UnsafeJWTSecretSource)
		headers := http.Header{"Authorization": []string{"Bearer " + token}}
		c := &http.Client{
			Transport: httpx.NewHeadersTransport(headers),
		}

		require.NoError(t, cmd.run(ctx, srv.URL, jsonl.NewEncoder(&output), &input, c))

		var result communityapi.PublishedContent
		require.NoError(t, jsonx.UnmarshalRead(&output, &result))

		var pc community.PublishedContent
		require.NoError(t, community.PublishedContentFindByID(ctx, q, result.Id).Scan(&pc))
		require.Equal(t, int32(communityapi.PublishMode_LISTED), pc.PublishMode)
	})

	t.Run("returns error when library item does not exist", func(t *testing.T) {
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

		communityID := uuid.Must(uuid.NewV7()).String()
		libraryID := uuid.Must(uuid.NewV7()).String()

		srv := httptest.NewServer(communityLibraryPublishServer(t, q))
		defer srv.Close()

		// library item not inserted — endpoint returns an error status
		com := &communityapi.Community{Id: communityID}
		lmd := library.Metadata{ID: libraryID}
		var input bytes.Buffer
		require.NoError(t, json.NewEncoder(&input).Encode(com))
		require.NoError(t, jsonl.NewEncoder(&input).Encode(lmd))

		cmd := cmdPublish{DryRun: false}
		token := httpauthtest.UnsafeToken(metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v))), httpauthtest.UnsafeJWTSecretSource)
		headers := http.Header{"Authorization": []string{"Bearer " + token}}
		c := &http.Client{
			Transport: httpx.NewHeadersTransport(headers),
		}

		require.Error(t, cmd.run(ctx, srv.URL, jsonl.NewEncoder(&bytes.Buffer{}), &input, c))
	})
}
