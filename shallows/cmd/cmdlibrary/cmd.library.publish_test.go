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
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/shallows/community"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/jsonl"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/testx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/retrovibed/retrovibed/shallows/meta"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"github.com/stretchr/testify/require"
)

func communityLibraryPublishServer(t *testing.T, q *sql.DB) *mux.Router {
	t.Helper()

	routes := mux.NewRouter()
	community.NewHTTP(
		q,
		community.HTTPOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		community.HTTPOptionHTTPClient(&http.Client{}),
		community.HTTPOptionMediaStorage(fsx.DirVirtual(t.TempDir())),
		community.HTTPOptionTorrentStorage(fsx.DirVirtual(t.TempDir())),
	).Bind(routes.PathPrefix("/c").Subrouter())

	return routes
}

func encodeLibraryItem(t *testing.T, enc *jsonl.Encoder, lmd library.Metadata) {
	t.Helper()
	require.NoError(t, enc.Encode(langx.Clone(lmd, timex.JSONSafeEncodeOption)))
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
			EncryptionSeed: uuid.Must(uuid.NewV4()).String(),
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

		com := &meta.Community{Id: communityID}
		var input bytes.Buffer
		encodeLibraryItem(t, jsonl.NewEncoder(&input), lmd)

		var output bytes.Buffer
		cmd := cmdPublish{Endpoint: "localhost:9998", DryRun: true}
		require.NoError(t, cmd.run(ctx, jsonl.NewEncoder(&output), &input, c, com))

		require.False(t, called, "endpoint must not be called during dry run")
		var decoded library.Metadata
		require.NoError(t, json.NewDecoder(&output).Decode(&decoded))
		require.Equal(t, libraryID, decoded.ID)
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
			EncryptionSeed: uuid.Must(uuid.NewV4()).String(),
		}
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, lmd).Scan(&lmd))

		srv := httptest.NewServer(communityLibraryPublishServer(t, q))
		defer srv.Close()

		com := &meta.Community{Id: communityID}
		var input bytes.Buffer
		encodeLibraryItem(t, jsonl.NewEncoder(&input), lmd)

		var output bytes.Buffer
		cmd := cmdPublish{Endpoint: srv.URL, DryRun: false}
		token := httpauthtest.UnsafeToken(metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v))), httpauthtest.UnsafeJWTSecretSource)
		headers := http.Header{"Authorization": []string{"Bearer " + token}}
		c := &http.Client{
			Transport: httpx.NewHeadersTransport(headers),
		}
		require.NoError(t, cmd.run(ctx, jsonl.NewEncoder(&output), &input, c, com))

		var result meta.PublishedContent
		require.NoError(t, json.NewDecoder(&output).Decode(&result))
		require.NotEmpty(t, result.Id)
		require.Equal(t, libraryID, result.LibraryId)
		require.Equal(t, communityID, result.CommunityId)

		var pc community.PublishedContent
		require.NoError(t, community.PublishedContentFindByID(ctx, q, result.Id).Scan(&pc))
		require.Equal(t, libraryID, pc.LibraryID)
		require.Equal(t, communityID, pc.CommunityID)
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

		var input bytes.Buffer
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
				EncryptionSeed: uuid.Must(uuid.NewV4()).String(),
			}
			require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, lmd).Scan(&lmd))
			encodeLibraryItem(t, enc, lmd)
		}

		srv := httptest.NewServer(communityLibraryPublishServer(t, q))
		defer srv.Close()

		com := &meta.Community{Id: communityID}
		var output bytes.Buffer
		cmd := cmdPublish{Endpoint: srv.URL, DryRun: false}
		token := httpauthtest.UnsafeToken(metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v))), httpauthtest.UnsafeJWTSecretSource)
		headers := http.Header{"Authorization": []string{"Bearer " + token}}
		c := &http.Client{
			Transport: httpx.NewHeadersTransport(headers),
		}

		require.NoError(t, cmd.run(ctx, jsonl.NewEncoder(&output), &input, c, com))

		dec := json.NewDecoder(&output)
		var results []*meta.PublishedContent
		for dec.More() {
			var pc meta.PublishedContent
			require.NoError(t, dec.Decode(&pc))
			results = append(results, &pc)
		}
		require.Len(t, results, 3)
		for i, result := range results {
			require.Equal(t, libraryIDs[i], result.LibraryId)
			require.Equal(t, communityID, result.CommunityId)
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
			EncryptionSeed: uuid.Must(uuid.NewV4()).String(),
		}
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, lmd).Scan(&lmd))

		srv := httptest.NewServer(communityLibraryPublishServer(t, q))
		defer srv.Close()

		com := &meta.Community{Id: communityID, DefaultPublishMode: meta.PublishMode_LISTED}
		var input bytes.Buffer
		encodeLibraryItem(t, jsonl.NewEncoder(&input), lmd)

		var output bytes.Buffer
		cmd := cmdPublish{Endpoint: srv.URL, DryRun: false}
		token := httpauthtest.UnsafeToken(metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v))), httpauthtest.UnsafeJWTSecretSource)
		headers := http.Header{"Authorization": []string{"Bearer " + token}}
		c := &http.Client{
			Transport: httpx.NewHeadersTransport(headers),
		}

		require.NoError(t, cmd.run(ctx, jsonl.NewEncoder(&output), &input, c, com))

		var result meta.PublishedContent
		require.NoError(t, json.NewDecoder(&output).Decode(&result))

		var pc community.PublishedContent
		require.NoError(t, community.PublishedContentFindByID(ctx, q, result.Id).Scan(&pc))
		require.Equal(t, int32(meta.PublishMode_LISTED), pc.PublishMode)
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
		lmd := library.Metadata{ID: libraryID}
		var input bytes.Buffer
		require.NoError(t, jsonl.NewEncoder(&input).Encode(lmd))

		com := &meta.Community{Id: communityID}
		cmd := cmdPublish{Endpoint: srv.URL, DryRun: false}
		token := httpauthtest.UnsafeToken(metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v))), httpauthtest.UnsafeJWTSecretSource)
		headers := http.Header{"Authorization": []string{"Bearer " + token}}
		c := &http.Client{
			Transport: httpx.NewHeadersTransport(headers),
		}

		require.Error(t, cmd.run(ctx, jsonl.NewEncoder(&bytes.Buffer{}), &input, c, com))
	})
}
