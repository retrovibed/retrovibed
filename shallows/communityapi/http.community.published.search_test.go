package communityapi_test

import (
	"encoding/json"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/mux"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/community"
	"github.com/retrovibed/retrovibed/shallows/communityapi"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/httptestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/retrovibed/retrovibed/shallows/meta"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"github.com/stretchr/testify/require"
)

func TestPublishedListEndpoint(t *testing.T) {
	t.Run("returns published content for a community", func(t *testing.T) {
		var (
			ctx, done  = testx.Context(t)
			q          = sqltestx.Metadatabase(t)
			p          meta.Profile
			v          meta.Authz
			mediaDir   = t.TempDir()
			torrentDir = t.TempDir()
		)
		defer done()

		communityID := uuid.Must(uuid.NewV7()).String()

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, testx.Fake(&v, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, v).Scan(&v))

		lmd1 := library.Metadata{
			ID:             uuid.Must(uuid.NewV7()).String(),
			Description:    "first media",
			Bytes:          1024,
			TorrentID:      uuid.Nil.String(),
			KnownMediaID:   uuid.Nil.String(),
			ArchiveID:      uuid.Nil.String(),
			DirectoryID:    uuid.Nil.String(),
			EncryptionSeed: uuid.Must(uuid.NewV4()).String(),
		}
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, lmd1).Scan(&lmd1))

		lmd2 := library.Metadata{
			ID:             uuid.Must(uuid.NewV7()).String(),
			Description:    "second media",
			Bytes:          2048,
			TorrentID:      uuid.Nil.String(),
			KnownMediaID:   uuid.Nil.String(),
			ArchiveID:      uuid.Nil.String(),
			DirectoryID:    uuid.Nil.String(),
			EncryptionSeed: uuid.Must(uuid.NewV4()).String(),
		}
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, lmd2).Scan(&lmd2))

		routes := mux.NewRouter()
		communityapi.NewHTTPPublished(
			q,
			communityapi.HTTPPublishedOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
			communityapi.HTTPPublishedOptionHTTPClient(&http.Client{}),
			communityapi.HTTPPublishedOptionMediaStorage(fsx.DirVirtual(mediaDir)),
			communityapi.HTTPPublishedOptionTorrentStorage(fsx.DirVirtual(torrentDir)),
		).Bind(routes.PathPrefix("/c").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))
		token := "Bearer " + httpauthtest.UnsafeToken(claims, httpauthtest.UnsafeJWTSecretSource)

		// Publish two items to the same community.
		for _, lid := range []string{lmd1.ID, lmd2.ID} {
			body, err := json.Marshal(&communityapi.PublishContentRequest{
				PublishedContent: &communityapi.PublishedContent{LibraryId: lid},
			})
			require.NoError(t, err)

			resp, req, err := httptestx.BuildRequestContextBytes(
				ctx,
				http.MethodPost,
				"/c/"+communityID,
				body,
				httptestx.RequestOptionAuthorization(token),
				httptestx.RequestOptionContent("application/json"),
			)
			require.NoError(t, err)
			routes.ServeHTTP(resp, req)
			require.Equal(t, http.StatusOK, resp.Code)
		}

		// List published content.
		resp, req, err := httptestx.BuildRequestContextBytes(
			ctx,
			http.MethodGet,
			"/c/"+communityID,
			nil,
			httptestx.RequestOptionAuthorization(token),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusOK, resp.Code)

		var result communityapi.PublishedContentSearchResponse
		require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &result))
		require.Len(t, result.Items, 2)

		ids := []string{result.Items[0].LibraryId, result.Items[1].LibraryId}
		require.Contains(t, ids, lmd1.ID)
		require.Contains(t, ids, lmd2.ID)
	})

	t.Run("excludes tombstoned content from results", func(t *testing.T) {
		var (
			ctx, done  = testx.Context(t)
			q          = sqltestx.Metadatabase(t)
			p          meta.Profile
			v          meta.Authz
			mediaDir   = t.TempDir()
			torrentDir = t.TempDir()
		)
		defer done()

		communityID := uuid.Must(uuid.NewV7()).String()

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, testx.Fake(&v, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, v).Scan(&v))

		lmd1 := library.Metadata{
			ID:             uuid.Must(uuid.NewV7()).String(),
			Description:    "active media",
			Bytes:          1024,
			TorrentID:      uuid.Nil.String(),
			KnownMediaID:   uuid.Nil.String(),
			ArchiveID:      uuid.Nil.String(),
			DirectoryID:    uuid.Nil.String(),
			EncryptionSeed: uuid.Must(uuid.NewV4()).String(),
		}
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, lmd1).Scan(&lmd1))

		lmd2 := library.Metadata{
			ID:             uuid.Must(uuid.NewV7()).String(),
			Description:    "tombstoned media",
			Bytes:          2048,
			TorrentID:      uuid.Nil.String(),
			KnownMediaID:   uuid.Nil.String(),
			ArchiveID:      uuid.Nil.String(),
			DirectoryID:    uuid.Nil.String(),
			EncryptionSeed: uuid.Must(uuid.NewV4()).String(),
		}
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, lmd2).Scan(&lmd2))

		var pc1 community.PublishedContent
		require.NoError(t, testx.Fake(&pc1, community.PublishedContentOptionTestDefaults, func(p *community.PublishedContent) {
			p.CommunityID = communityID
			p.LibraryID = lmd1.ID
			p.Bytes = lmd1.Bytes
		}))
		require.NoError(t, community.PublishedContentInsertWithDefaults(ctx, q, pc1).Scan(&pc1))

		var pc2 community.PublishedContent
		require.NoError(t, testx.Fake(&pc2, community.PublishedContentOptionTestDefaults, func(p *community.PublishedContent) {
			p.CommunityID = communityID
			p.LibraryID = lmd2.ID
			p.Bytes = lmd2.Bytes
		}))
		require.NoError(t, community.PublishedContentInsertWithDefaults(ctx, q, pc2).Scan(&pc2))
		require.NoError(t, community.PublishedContentTombstone(ctx, q, pc2.ID).Scan(&pc2))

		routes := mux.NewRouter()
		communityapi.NewHTTPPublished(
			q,
			communityapi.HTTPPublishedOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
			communityapi.HTTPPublishedOptionHTTPClient(&http.Client{}),
			communityapi.HTTPPublishedOptionMediaStorage(fsx.DirVirtual(mediaDir)),
			communityapi.HTTPPublishedOptionTorrentStorage(fsx.DirVirtual(torrentDir)),
		).Bind(routes.PathPrefix("/c").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))
		token := "Bearer " + httpauthtest.UnsafeToken(claims, httpauthtest.UnsafeJWTSecretSource)

		resp, req, err := httptestx.BuildRequestContextBytes(
			ctx,
			http.MethodGet,
			"/c/"+communityID,
			nil,
			httptestx.RequestOptionAuthorization(token),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusOK, resp.Code)

		var result communityapi.PublishedContentSearchResponse
		require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &result))
		require.Len(t, result.Items, 1)
		require.Equal(t, lmd1.ID, result.Items[0].LibraryId)
	})

	t.Run("returns content sorted by published_at descending", func(t *testing.T) {
		var (
			ctx, done  = testx.Context(t)
			q          = sqltestx.Metadatabase(t)
			p          meta.Profile
			v          meta.Authz
			mediaDir   = t.TempDir()
			torrentDir = t.TempDir()
		)
		defer done()

		communityID := uuid.Must(uuid.NewV7()).String()

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, testx.Fake(&v, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, v).Scan(&v))

		lmd1 := library.Metadata{
			ID:             uuid.Must(uuid.NewV7()).String(),
			Bytes:          1024,
			TorrentID:      uuid.Nil.String(),
			KnownMediaID:   uuid.Nil.String(),
			ArchiveID:      uuid.Nil.String(),
			DirectoryID:    uuid.Nil.String(),
			EncryptionSeed: uuid.Must(uuid.NewV4()).String(),
		}
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, lmd1).Scan(&lmd1))

		lmd2 := library.Metadata{
			ID:             uuid.Must(uuid.NewV7()).String(),
			Bytes:          2048,
			TorrentID:      uuid.Nil.String(),
			KnownMediaID:   uuid.Nil.String(),
			ArchiveID:      uuid.Nil.String(),
			DirectoryID:    uuid.Nil.String(),
			EncryptionSeed: uuid.Must(uuid.NewV4()).String(),
		}
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, lmd2).Scan(&lmd2))

		now := time.Now()

		var pc1 community.PublishedContent
		require.NoError(t, testx.Fake(&pc1, community.PublishedContentOptionTestDefaults, func(p *community.PublishedContent) {
			p.CommunityID = communityID
			p.LibraryID = lmd1.ID
			p.Bytes = lmd1.Bytes
		}))
		require.NoError(t, community.PublishedContentInsertWithDefaults(ctx, q, pc1).Scan(&pc1))
		require.NoError(t, community.PublishedContentUpdatePublishedAt(ctx, q, pc1.ID, now.Add(-48*time.Hour)).Scan(&pc1))

		var pc2 community.PublishedContent
		require.NoError(t, testx.Fake(&pc2, community.PublishedContentOptionTestDefaults, func(p *community.PublishedContent) {
			p.CommunityID = communityID
			p.LibraryID = lmd2.ID
			p.Bytes = lmd2.Bytes
		}))
		require.NoError(t, community.PublishedContentInsertWithDefaults(ctx, q, pc2).Scan(&pc2))
		require.NoError(t, community.PublishedContentUpdatePublishedAt(ctx, q, pc2.ID, now.Add(-24*time.Hour)).Scan(&pc2))

		routes := mux.NewRouter()
		communityapi.NewHTTPPublished(
			q,
			communityapi.HTTPPublishedOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
			communityapi.HTTPPublishedOptionHTTPClient(&http.Client{}),
			communityapi.HTTPPublishedOptionMediaStorage(fsx.DirVirtual(mediaDir)),
			communityapi.HTTPPublishedOptionTorrentStorage(fsx.DirVirtual(torrentDir)),
		).Bind(routes.PathPrefix("/c").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))
		token := "Bearer " + httpauthtest.UnsafeToken(claims, httpauthtest.UnsafeJWTSecretSource)

		resp, req, err := httptestx.BuildRequestContextBytes(
			ctx,
			http.MethodGet,
			"/c/"+communityID,
			nil,
			httptestx.RequestOptionAuthorization(token),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusOK, resp.Code)

		var result communityapi.PublishedContentSearchResponse
		require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &result))
		require.Len(t, result.Items, 2)
		require.Equal(t, lmd2.ID, result.Items[0].LibraryId) // newer published_at first
		require.Equal(t, lmd1.ID, result.Items[1].LibraryId)
	})

	t.Run("filters content by query matching title", func(t *testing.T) {
		var (
			ctx, done  = testx.Context(t)
			q          = sqltestx.Metadatabase(t)
			p          meta.Profile
			v          meta.Authz
			mediaDir   = t.TempDir()
			torrentDir = t.TempDir()
		)
		defer done()

		communityID := uuid.Must(uuid.NewV7()).String()

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, testx.Fake(&v, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, v).Scan(&v))

		lmd1 := library.Metadata{
			ID:             uuid.Must(uuid.NewV7()).String(),
			Bytes:          1024,
			TorrentID:      uuid.Nil.String(),
			KnownMediaID:   uuid.Nil.String(),
			ArchiveID:      uuid.Nil.String(),
			DirectoryID:    uuid.Nil.String(),
			EncryptionSeed: uuid.Must(uuid.NewV4()).String(),
		}
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, lmd1).Scan(&lmd1))

		lmd2 := library.Metadata{
			ID:             uuid.Must(uuid.NewV7()).String(),
			Bytes:          2048,
			TorrentID:      uuid.Nil.String(),
			KnownMediaID:   uuid.Nil.String(),
			ArchiveID:      uuid.Nil.String(),
			DirectoryID:    uuid.Nil.String(),
			EncryptionSeed: uuid.Must(uuid.NewV4()).String(),
		}
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, lmd2).Scan(&lmd2))

		routes := mux.NewRouter()
		communityapi.NewHTTPPublished(
			q,
			communityapi.HTTPPublishedOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
			communityapi.HTTPPublishedOptionHTTPClient(&http.Client{}),
			communityapi.HTTPPublishedOptionMediaStorage(fsx.DirVirtual(mediaDir)),
			communityapi.HTTPPublishedOptionTorrentStorage(fsx.DirVirtual(torrentDir)),
		).Bind(routes.PathPrefix("/c").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))
		token := "Bearer " + httpauthtest.UnsafeToken(claims, httpauthtest.UnsafeJWTSecretSource)

		for _, item := range []struct {
			lid   string
			title string
		}{
			{lmd1.ID, "alpha video"},
			{lmd2.ID, "beta video"},
		} {
			body, err := json.Marshal(&communityapi.PublishContentRequest{
				PublishedContent: &communityapi.PublishedContent{LibraryId: item.lid, Title: item.title},
			})
			require.NoError(t, err)

			resp, req, err := httptestx.BuildRequestContextBytes(
				ctx,
				http.MethodPost,
				"/c/"+communityID,
				body,
				httptestx.RequestOptionAuthorization(token),
				httptestx.RequestOptionContent("application/json"),
			)
			require.NoError(t, err)
			routes.ServeHTTP(resp, req)
			require.Equal(t, http.StatusOK, resp.Code)
		}

		resp, req, err := httptestx.BuildRequestContextBytes(
			ctx,
			http.MethodGet,
			"/c/"+communityID+"?query=alpha",
			nil,
			httptestx.RequestOptionAuthorization(token),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusOK, resp.Code)

		var result communityapi.PublishedContentSearchResponse
		require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &result))
		require.Len(t, result.Items, 1)
		require.Equal(t, lmd1.ID, result.Items[0].LibraryId)
	})

	t.Run("returns empty list when query matches no titles", func(t *testing.T) {
		var (
			ctx, done  = testx.Context(t)
			q          = sqltestx.Metadatabase(t)
			p          meta.Profile
			v          meta.Authz
			mediaDir   = t.TempDir()
			torrentDir = t.TempDir()
		)
		defer done()

		communityID := uuid.Must(uuid.NewV7()).String()

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, testx.Fake(&v, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, v).Scan(&v))

		lmd1 := library.Metadata{
			ID:             uuid.Must(uuid.NewV7()).String(),
			Bytes:          1024,
			TorrentID:      uuid.Nil.String(),
			KnownMediaID:   uuid.Nil.String(),
			ArchiveID:      uuid.Nil.String(),
			DirectoryID:    uuid.Nil.String(),
			EncryptionSeed: uuid.Must(uuid.NewV4()).String(),
		}
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, lmd1).Scan(&lmd1))

		routes := mux.NewRouter()
		communityapi.NewHTTPPublished(
			q,
			communityapi.HTTPPublishedOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
			communityapi.HTTPPublishedOptionHTTPClient(&http.Client{}),
			communityapi.HTTPPublishedOptionMediaStorage(fsx.DirVirtual(mediaDir)),
			communityapi.HTTPPublishedOptionTorrentStorage(fsx.DirVirtual(torrentDir)),
		).Bind(routes.PathPrefix("/c").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))
		token := "Bearer " + httpauthtest.UnsafeToken(claims, httpauthtest.UnsafeJWTSecretSource)

		body, err := json.Marshal(&communityapi.PublishContentRequest{
			PublishedContent: &communityapi.PublishedContent{LibraryId: lmd1.ID, Title: "specific title"},
		})
		require.NoError(t, err)

		resp, req, err := httptestx.BuildRequestContextBytes(
			ctx,
			http.MethodPost,
			"/c/"+communityID,
			body,
			httptestx.RequestOptionAuthorization(token),
			httptestx.RequestOptionContent("application/json"),
		)
		require.NoError(t, err)
		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusOK, resp.Code)

		resp, req, err = httptestx.BuildRequestContextBytes(
			ctx,
			http.MethodGet,
			"/c/"+communityID+"?query=nomatch",
			nil,
			httptestx.RequestOptionAuthorization(token),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusOK, resp.Code)

		var result communityapi.PublishedContentSearchResponse
		require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &result))
		require.Empty(t, result.Items)
	})

	t.Run("returns empty list for community with no content", func(t *testing.T) {
		var (
			ctx, done   = testx.Context(t)
			q           = sqltestx.Metadatabase(t)
			p           meta.Profile
			v           meta.Authz
			communityID = uuid.Must(uuid.NewV7()).String()
			mediaDir    = t.TempDir()
			torrentDir  = t.TempDir()
		)
		defer done()

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, testx.Fake(&v, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, v).Scan(&v))

		routes := mux.NewRouter()
		communityapi.NewHTTPPublished(
			q,
			communityapi.HTTPPublishedOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
			communityapi.HTTPPublishedOptionMediaStorage(fsx.DirVirtual(mediaDir)),
			communityapi.HTTPPublishedOptionTorrentStorage(fsx.DirVirtual(torrentDir)),
		).Bind(routes.PathPrefix("/c").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))
		token := "Bearer " + httpauthtest.UnsafeToken(claims, httpauthtest.UnsafeJWTSecretSource)

		resp, req, err := httptestx.BuildRequestContextBytes(
			ctx,
			http.MethodGet,
			"/c/"+communityID,
			nil,
			httptestx.RequestOptionAuthorization(token),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusOK, resp.Code)

		var result communityapi.PublishedContentSearchResponse
		require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &result))
		require.Empty(t, result.Items)
	})

	t.Run("paginates results using limit and offset", func(t *testing.T) {
		var (
			ctx, done  = testx.Context(t)
			q          = sqltestx.Metadatabase(t)
			p          meta.Profile
			v          meta.Authz
			mediaDir   = t.TempDir()
			torrentDir = t.TempDir()
		)
		defer done()

		communityID := uuid.Must(uuid.NewV7()).String()

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, testx.Fake(&v, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, v).Scan(&v))

		now := time.Now()
		lmds := make([]library.Metadata, 3)
		pcs := make([]community.PublishedContent, 3)
		for i := range lmds {
			lmds[i] = library.Metadata{
				ID:             uuid.Must(uuid.NewV7()).String(),
				Bytes:          1024,
				TorrentID:      uuid.Nil.String(),
				KnownMediaID:   uuid.Nil.String(),
				ArchiveID:      uuid.Nil.String(),
				DirectoryID:    uuid.Nil.String(),
				EncryptionSeed: uuid.Must(uuid.NewV4()).String(),
			}
			require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, lmds[i]).Scan(&lmds[i]))

			require.NoError(t, testx.Fake(&pcs[i], community.PublishedContentOptionTestDefaults, func(p *community.PublishedContent) {
				p.CommunityID = communityID
				p.LibraryID = lmds[i].ID
				p.Bytes = lmds[i].Bytes
			}))
			require.NoError(t, community.PublishedContentInsertWithDefaults(ctx, q, pcs[i]).Scan(&pcs[i]))
			// oldest -> newest: lmds[0], lmds[1], lmds[2]; published_at DESC returns lmds[2] first.
			require.NoError(t, community.PublishedContentUpdatePublishedAt(ctx, q, pcs[i].ID, now.Add(-time.Duration(len(lmds)-i)*24*time.Hour)).Scan(&pcs[i]))
		}

		routes := mux.NewRouter()
		communityapi.NewHTTPPublished(
			q,
			communityapi.HTTPPublishedOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
			communityapi.HTTPPublishedOptionHTTPClient(&http.Client{}),
			communityapi.HTTPPublishedOptionMediaStorage(fsx.DirVirtual(mediaDir)),
			communityapi.HTTPPublishedOptionTorrentStorage(fsx.DirVirtual(torrentDir)),
		).Bind(routes.PathPrefix("/c").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))
		token := "Bearer " + httpauthtest.UnsafeToken(claims, httpauthtest.UnsafeJWTSecretSource)

		expectedOrder := []string{lmds[2].ID, lmds[1].ID, lmds[0].ID}
		var seen []string
		for page, offset := range []int{0, 1, 2} {
			resp, req, err := httptestx.BuildRequestContextBytes(
				ctx,
				http.MethodGet,
				"/c/"+communityID+"?limit=1&offset="+strconv.Itoa(offset),
				nil,
				httptestx.RequestOptionAuthorization(token),
			)
			require.NoError(t, err)

			routes.ServeHTTP(resp, req)
			require.Equal(t, http.StatusOK, resp.Code)

			var result communityapi.PublishedContentSearchResponse
			require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &result))
			require.Len(t, result.Items, 1, "page %d", page)
			require.Equal(t, expectedOrder[page], result.Items[0].LibraryId, "page %d", page)
			seen = append(seen, result.Items[0].LibraryId)
		}

		require.ElementsMatch(t, expectedOrder, seen, "pagination must cover every item exactly once")

		// requesting beyond the last page returns no items.
		resp, req, err := httptestx.BuildRequestContextBytes(
			ctx,
			http.MethodGet,
			"/c/"+communityID+"?limit=1&offset=3",
			nil,
			httptestx.RequestOptionAuthorization(token),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusOK, resp.Code)

		var result communityapi.PublishedContentSearchResponse
		require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &result))
		require.Empty(t, result.Items)
	})
}
