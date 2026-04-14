package community_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/mux"
	"github.com/retrovibed/retrovibed/shallows/community"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/httptestx"
	"github.com/retrovibed/retrovibed/shallows/internal/jwtx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/testx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/retrovibed/retrovibed/shallows/meta"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"github.com/retrovibed/retrovibed/shallows/tracking"
	"github.com/stretchr/testify/require"
)

func TestPublishEndpoint(t *testing.T) {
	t.Run("returns immediately without torrent for new content", func(t *testing.T) {
		var (
			ctx, done   = testx.Context(t)
			q           = sqltestx.Metadatabase(t)
			p           meta.Profile
			v           meta.Authz
			communityID = uuid.Must(uuid.NewV7()).String()
			libraryID   = uuid.Must(uuid.NewV7()).String()
			mediaDir    = t.TempDir()
			torrentDir  = t.TempDir()
		)
		defer done()

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
		require.Equal(t, uuid.Nil.String(), lmd.TorrentID)
		require.Equal(t, uuid.Nil.String(), lmd.ArchiveID)

		routes := mux.NewRouter()
		community.NewHTTP(
			q,
			community.HTTPOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
			community.HTTPOptionMediaStorage(fsx.DirVirtual(mediaDir)),
			community.HTTPOptionTorrentStorage(fsx.DirVirtual(torrentDir)),
		).Bind(routes.PathPrefix("/c").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))
		body, err := json.Marshal(&meta.PublishContentRequest{
			PublishedContent: &meta.PublishedContent{
				LibraryId: libraryID,
			},
		})
		require.NoError(t, err)

		resp, req, err := httptestx.BuildRequestContextBytes(
			ctx,
			http.MethodPost,
			"/c/"+communityID+"/publish",
			body,
			httptestx.RequestOptionAuthorization("Bearer "+httpauthtest.UnsafeToken(claims, httpauthtest.UnsafeJWTSecretSource)),
			httptestx.RequestOptionContent("application/json"),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusOK, resp.Code)

		var result meta.PublishContentResponse
		require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &result))
		require.NotEmpty(t, result.PublishedContent.Id)
		require.Equal(t, "", result.PublishedContent.MagnetUri)

		var pc community.PublishedContent
		require.NoError(t, community.PublishedContentFindByID(ctx, q, result.PublishedContent.Id).Scan(&pc))
		require.Equal(t, "", pc.MagnetURI)
		require.Equal(t, libraryID, pc.LibraryID)

		var updatedLmd library.Metadata
		require.NoError(t, library.MetadataFindByID(ctx, q, libraryID).Scan(&updatedLmd))
		require.Equal(t, uuid.Nil.String(), updatedLmd.ArchiveID)
	})

	t.Run("listed stores publish mode and skips archival marking", func(t *testing.T) {
		var (
			ctx, done   = testx.Context(t)
			q           = sqltestx.Metadatabase(t)
			p           meta.Profile
			v           meta.Authz
			communityID = uuid.Must(uuid.NewV7()).String()
			libraryID   = uuid.Must(uuid.NewV7()).String()
			mediaDir    = t.TempDir()
			torrentDir  = t.TempDir()
		)
		defer done()

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

		routes := mux.NewRouter()
		community.NewHTTP(
			q,
			community.HTTPOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
			community.HTTPOptionHTTPClient(&http.Client{}),
			community.HTTPOptionMediaStorage(fsx.DirVirtual(mediaDir)),
			community.HTTPOptionTorrentStorage(fsx.DirVirtual(torrentDir)),
		).Bind(routes.PathPrefix("/c").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))
		body, err := json.Marshal(&meta.PublishContentRequest{
			PublishedContent: &meta.PublishedContent{
				LibraryId: libraryID,
			},
			PublishMode: meta.PublishMode_LISTED,
		})
		require.NoError(t, err)

		resp, req, err := httptestx.BuildRequestContextBytes(
			ctx,
			http.MethodPost,
			"/c/"+communityID+"/publish",
			body,
			httptestx.RequestOptionAuthorization("Bearer "+httpauthtest.UnsafeToken(claims, httpauthtest.UnsafeJWTSecretSource)),
			httptestx.RequestOptionContent("application/json"),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusOK, resp.Code)

		var result meta.PublishContentResponse
		require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &result))
		require.NotEmpty(t, result.PublishedContent.Id)

		var pc community.PublishedContent
		require.NoError(t, community.PublishedContentFindByID(ctx, q, result.PublishedContent.Id).Scan(&pc))
		require.Equal(t, int32(meta.PublishMode_LISTED), pc.PublishMode)

		var updatedLmd library.Metadata
		require.NoError(t, library.MetadataFindByID(ctx, q, libraryID).Scan(&updatedLmd))
		require.Equal(t, uuid.Nil.String(), updatedLmd.ArchiveID)
	})

	t.Run("syndicated marks library as archivable", func(t *testing.T) {
		var (
			ctx, done   = testx.Context(t)
			q           = sqltestx.Metadatabase(t)
			p           meta.Profile
			v           meta.Authz
			communityID = uuid.Must(uuid.NewV7()).String()
			libraryID   = uuid.Must(uuid.NewV7()).String()
			mediaDir    = t.TempDir()
			torrentDir  = t.TempDir()
		)
		defer done()

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

		routes := mux.NewRouter()
		community.NewHTTP(
			q,
			community.HTTPOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
			community.HTTPOptionHTTPClient(&http.Client{}),
			community.HTTPOptionMediaStorage(fsx.DirVirtual(mediaDir)),
			community.HTTPOptionTorrentStorage(fsx.DirVirtual(torrentDir)),
		).Bind(routes.PathPrefix("/c").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))
		body, err := json.Marshal(&meta.PublishContentRequest{
			PublishedContent: &meta.PublishedContent{
				LibraryId: libraryID,
			},
			PublishMode: meta.PublishMode_SYNDICATED,
		})
		require.NoError(t, err)

		resp, req, err := httptestx.BuildRequestContextBytes(
			ctx,
			http.MethodPost,
			"/c/"+communityID+"/publish",
			body,
			httptestx.RequestOptionAuthorization("Bearer "+httpauthtest.UnsafeToken(claims, httpauthtest.UnsafeJWTSecretSource)),
			httptestx.RequestOptionContent("application/json"),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusOK, resp.Code)

		var result meta.PublishContentResponse
		require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &result))
		require.NotEmpty(t, result.PublishedContent.Id)

		var pc community.PublishedContent
		require.NoError(t, community.PublishedContentFindByID(ctx, q, result.PublishedContent.Id).Scan(&pc))
		require.Equal(t, int32(meta.PublishMode_SYNDICATED), pc.PublishMode)

		var updatedLmd library.Metadata
		require.NoError(t, library.MetadataFindByID(ctx, q, libraryID).Scan(&updatedLmd))
		require.Equal(t, uuid.Max.String(), updatedLmd.ArchiveID)
	})

	t.Run("returns existing magnet uri when torrent exists", func(t *testing.T) {
		var (
			ctx, done   = testx.Context(t)
			q           = sqltestx.Metadatabase(t)
			p           meta.Profile
			v           meta.Authz
			communityID = uuid.Must(uuid.NewV7()).String()
			libraryID   = uuid.Must(uuid.NewV7()).String()
			torrentID   = uuid.Must(uuid.NewV7()).String()
			mediaDir    = t.TempDir()
			torrentDir  = t.TempDir()
		)
		defer done()

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, testx.Fake(&v, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, v).Scan(&v))

		tmd := tracking.Metadata{
			ID:             torrentID,
			Description:    "test torrent",
			Infohash:       []byte{0x0b, 0xee, 0xc7, 0xb5, 0xea, 0x3f, 0x0f, 0xdb, 0xc9, 0x5d, 0x0d, 0xd4, 0x7f, 0x3c, 0x5b, 0xc2, 0x75, 0xda, 0x8a, 0x33},
			KnownMediaID:   uuid.Nil.String(),
			EncryptionSeed: uuid.Must(uuid.NewV4()).String(),
		}
		require.NoError(t, tracking.MetadataInsertWithDefaults(ctx, q, tmd).Scan(&tmd))

		lmd := library.Metadata{
			ID:             libraryID,
			Description:    "test media",
			Bytes:          1024,
			TorrentID:      torrentID,
			KnownMediaID:   uuid.Nil.String(),
			ArchiveID:      uuid.Nil.String(),
			EncryptionSeed: uuid.Must(uuid.NewV4()).String(),
		}
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, lmd).Scan(&lmd))

		routes := mux.NewRouter()
		community.NewHTTP(
			q,
			community.HTTPOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
			community.HTTPOptionMediaStorage(fsx.DirVirtual(mediaDir)),
			community.HTTPOptionTorrentStorage(fsx.DirVirtual(torrentDir)),
		).Bind(routes.PathPrefix("/c").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))
		body, err := json.Marshal(&meta.PublishContentRequest{
			PublishedContent: &meta.PublishedContent{
				LibraryId: libraryID,
			},
		})
		require.NoError(t, err)

		resp, req, err := httptestx.BuildRequestContextBytes(
			ctx,
			http.MethodPost,
			"/c/"+communityID+"/publish",
			body,
			httptestx.RequestOptionAuthorization("Bearer "+httpauthtest.UnsafeToken(claims, httpauthtest.UnsafeJWTSecretSource)),
			httptestx.RequestOptionContent("application/json"),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusOK, resp.Code)

		var result meta.PublishContentResponse
		require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &result))
		require.NotEmpty(t, result.PublishedContent.MagnetUri)
		require.Contains(t, result.PublishedContent.MagnetUri, "magnet:?")
	})
}
