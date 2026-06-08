package communityapi_test

import (
	"encoding/json"
	"net/http"
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
	"github.com/retrovibed/retrovibed/shallows/internal/grpcx"
	"github.com/retrovibed/retrovibed/shallows/internal/httptestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
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
			Mimetype:       "audio/mpeg",
			TorrentID:      uuid.Nil.String(),
			KnownMediaID:   uuid.Nil.String(),
			ArchiveID:      uuid.Nil.String(),
			EncryptionSeed: uuid.Must(uuid.NewV4()).String(),
		}
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, lmd).Scan(&lmd))
		require.Equal(t, uuid.Nil.String(), lmd.TorrentID)
		require.Equal(t, uuid.Nil.String(), lmd.ArchiveID)

		routes := mux.NewRouter()
		communityapi.NewHTTPPublished(
			q,
			communityapi.HTTPPublishedOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
			communityapi.HTTPPublishedOptionHTTPClient(&http.Client{}),
			communityapi.HTTPPublishedOptionMediaStorage(fsx.DirVirtual(mediaDir)),
			communityapi.HTTPPublishedOptionTorrentStorage(fsx.DirVirtual(torrentDir)),
		).Bind(routes.PathPrefix("/c").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))
		body, err := json.Marshal(&communityapi.PublishContentRequest{
			PublishedContent: &communityapi.PublishedContent{
				LibraryId: libraryID,
			},
		})
		require.NoError(t, err)

		resp, req, err := httptestx.BuildRequestContextBytes(
			ctx,
			http.MethodPost,
			"/c/"+communityID,
			body,
			httptestx.RequestOptionAuthorization("Bearer "+httpauthtest.UnsafeToken(claims, httpauthtest.UnsafeJWTSecretSource)),
			httptestx.RequestOptionContent("application/json"),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusOK, resp.Code)

		var result communityapi.PublishContentResponse
		require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &result))
		require.NotEmpty(t, result.PublishedContent.Id)
		require.Equal(t, "", result.PublishedContent.MagnetUri)

		var pc community.PublishedContent
		require.NoError(t, community.PublishedContentFindByID(ctx, q, result.PublishedContent.Id).Scan(&pc))
		require.Equal(t, "", pc.MagnetURI)
		require.Equal(t, libraryID, pc.LibraryID)
		require.Equal(t, lmd.Mimetype, pc.Mimetype)
		require.Equal(t, lmd.Bytes, pc.Bytes)

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
			Mimetype:       "audio/mpeg",
			TorrentID:      uuid.Nil.String(),
			KnownMediaID:   uuid.Nil.String(),
			ArchiveID:      uuid.Nil.String(),
			EncryptionSeed: uuid.Must(uuid.NewV4()).String(),
		}
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, lmd).Scan(&lmd))

		routes := mux.NewRouter()
		communityapi.NewHTTPPublished(
			q,
			communityapi.HTTPPublishedOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
			communityapi.HTTPPublishedOptionHTTPClient(&http.Client{}),
			communityapi.HTTPPublishedOptionMediaStorage(fsx.DirVirtual(mediaDir)),
			communityapi.HTTPPublishedOptionTorrentStorage(fsx.DirVirtual(torrentDir)),
		).Bind(routes.PathPrefix("/c").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))
		body, err := json.Marshal(&communityapi.PublishContentRequest{
			PublishedContent: &communityapi.PublishedContent{
				LibraryId: libraryID,
			},
			PublishMode: communityapi.PublishMode_LISTED,
		})
		require.NoError(t, err)

		resp, req, err := httptestx.BuildRequestContextBytes(
			ctx,
			http.MethodPost,
			"/c/"+communityID,
			body,
			httptestx.RequestOptionAuthorization("Bearer "+httpauthtest.UnsafeToken(claims, httpauthtest.UnsafeJWTSecretSource)),
			httptestx.RequestOptionContent("application/json"),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusOK, resp.Code)

		var result communityapi.PublishContentResponse
		require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &result))
		require.NotEmpty(t, result.PublishedContent.Id)

		var pc community.PublishedContent
		require.NoError(t, community.PublishedContentFindByID(ctx, q, result.PublishedContent.Id).Scan(&pc))
		require.Equal(t, int32(communityapi.PublishMode_LISTED), pc.PublishMode)
		require.Equal(t, lmd.Mimetype, pc.Mimetype)
		require.Equal(t, lmd.Bytes, pc.Bytes)

		var updatedLmd library.Metadata
		require.NoError(t, library.MetadataFindByID(ctx, q, libraryID).Scan(&updatedLmd))
		require.Equal(t, uuid.Nil.String(), updatedLmd.ArchiveID)
	})

	t.Run("syndicated stores publish mode and skips archival marking", func(t *testing.T) {
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
			Mimetype:       "audio/mpeg",
			TorrentID:      uuid.Nil.String(),
			KnownMediaID:   uuid.Nil.String(),
			ArchiveID:      uuid.Nil.String(),
			EncryptionSeed: uuid.Must(uuid.NewV4()).String(),
		}
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, lmd).Scan(&lmd))

		routes := mux.NewRouter()
		communityapi.NewHTTPPublished(
			q,
			communityapi.HTTPPublishedOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
			communityapi.HTTPPublishedOptionHTTPClient(&http.Client{}),
			communityapi.HTTPPublishedOptionMediaStorage(fsx.DirVirtual(mediaDir)),
			communityapi.HTTPPublishedOptionTorrentStorage(fsx.DirVirtual(torrentDir)),
		).Bind(routes.PathPrefix("/c").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))
		body, err := json.Marshal(&communityapi.PublishContentRequest{
			PublishedContent: &communityapi.PublishedContent{
				LibraryId: libraryID,
			},
			PublishMode: communityapi.PublishMode_SYNDICATED,
		})
		require.NoError(t, err)

		resp, req, err := httptestx.BuildRequestContextBytes(
			ctx,
			http.MethodPost,
			"/c/"+communityID,
			body,
			httptestx.RequestOptionAuthorization("Bearer "+httpauthtest.UnsafeToken(claims, httpauthtest.UnsafeJWTSecretSource)),
			httptestx.RequestOptionContent("application/json"),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusOK, resp.Code)

		var result communityapi.PublishContentResponse
		require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &result))
		require.NotEmpty(t, result.PublishedContent.Id)

		var pc community.PublishedContent
		require.NoError(t, community.PublishedContentFindByID(ctx, q, result.PublishedContent.Id).Scan(&pc))
		require.Equal(t, int32(communityapi.PublishMode_SYNDICATED), pc.PublishMode)
		require.Equal(t, lmd.Mimetype, pc.Mimetype)
		require.Equal(t, lmd.Bytes, pc.Bytes)

		var updatedLmd library.Metadata
		require.NoError(t, library.MetadataFindByID(ctx, q, libraryID).Scan(&updatedLmd))
		require.Equal(t, uuid.Nil.String(), updatedLmd.ArchiveID)
	})

	t.Run("timestamps are set by db on insert when not provided in request", func(t *testing.T) {
		var (
			ctx, done   = testx.Context(t)
			q           = sqltestx.Metadatabase(t)
			p           meta.Profile
			v           meta.Authz
			communityID = uuid.Must(uuid.NewV7()).String()
			libraryID   = uuid.Must(uuid.NewV7()).String()
			mediaDir    = t.TempDir()
			torrentDir  = t.TempDir()
			before      = time.Now()
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
			Mimetype:       "audio/mpeg",
			TorrentID:      uuid.Nil.String(),
			KnownMediaID:   uuid.Nil.String(),
			ArchiveID:      uuid.Nil.String(),
			EncryptionSeed: uuid.Must(uuid.NewV4()).String(),
		}
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, lmd).Scan(&lmd))

		routes := mux.NewRouter()
		communityapi.NewHTTPPublished(
			q,
			communityapi.HTTPPublishedOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
			communityapi.HTTPPublishedOptionHTTPClient(&http.Client{}),
			communityapi.HTTPPublishedOptionMediaStorage(fsx.DirVirtual(mediaDir)),
			communityapi.HTTPPublishedOptionTorrentStorage(fsx.DirVirtual(torrentDir)),
		).Bind(routes.PathPrefix("/c").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))
		body, err := json.Marshal(&communityapi.PublishContentRequest{
			PublishedContent: &communityapi.PublishedContent{
				LibraryId: libraryID,
			},
		})
		require.NoError(t, err)

		resp, req, err := httptestx.BuildRequestContextBytes(
			ctx,
			http.MethodPost,
			"/c/"+communityID,
			body,
			httptestx.RequestOptionAuthorization("Bearer "+httpauthtest.UnsafeToken(claims, httpauthtest.UnsafeJWTSecretSource)),
			httptestx.RequestOptionContent("application/json"),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusOK, resp.Code)

		var result communityapi.PublishContentResponse
		require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &result))
		require.NotEmpty(t, result.PublishedContent.Id)

		var pc community.PublishedContent
		require.NoError(t, community.PublishedContentFindByID(ctx, q, result.PublishedContent.Id).Scan(&pc))
		require.False(t, pc.CreatedAt.IsZero())
		require.False(t, pc.UpdatedAt.IsZero())
		require.False(t, pc.PublishedAt.IsZero())
		require.WithinDuration(t, before, pc.CreatedAt, time.Minute)
		require.WithinDuration(t, before, pc.UpdatedAt, time.Minute)
		require.WithinDuration(t, timex.Inf(), pc.PublishedAt, time.Minute)
	})

	t.Run("response timestamps match stored timestamps", func(t *testing.T) {
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
			Mimetype:       "audio/mpeg",
			TorrentID:      uuid.Nil.String(),
			KnownMediaID:   uuid.Nil.String(),
			ArchiveID:      uuid.Nil.String(),
			EncryptionSeed: uuid.Must(uuid.NewV4()).String(),
		}
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, lmd).Scan(&lmd))

		routes := mux.NewRouter()
		communityapi.NewHTTPPublished(
			q,
			communityapi.HTTPPublishedOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
			communityapi.HTTPPublishedOptionHTTPClient(&http.Client{}),
			communityapi.HTTPPublishedOptionMediaStorage(fsx.DirVirtual(mediaDir)),
			communityapi.HTTPPublishedOptionTorrentStorage(fsx.DirVirtual(torrentDir)),
		).Bind(routes.PathPrefix("/c").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))
		body, err := json.Marshal(&communityapi.PublishContentRequest{
			PublishedContent: &communityapi.PublishedContent{
				LibraryId:   libraryID,
				CreatedAt:   grpcx.EncodeTime(time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)),
				UpdatedAt:   grpcx.EncodeTime(time.Date(2024, 2, 20, 12, 0, 0, 0, time.UTC)),
				PublishedAt: grpcx.EncodeTime(time.Date(2024, 3, 25, 14, 0, 0, 0, time.UTC)),
			},
		})
		require.NoError(t, err)

		resp, req, err := httptestx.BuildRequestContextBytes(
			ctx,
			http.MethodPost,
			"/c/"+communityID,
			body,
			httptestx.RequestOptionAuthorization("Bearer "+httpauthtest.UnsafeToken(claims, httpauthtest.UnsafeJWTSecretSource)),
			httptestx.RequestOptionContent("application/json"),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusOK, resp.Code)

		var result communityapi.PublishContentResponse
		require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &result))
		require.NotEmpty(t, result.PublishedContent.Id)

		var pc community.PublishedContent
		require.NoError(t, community.PublishedContentFindByID(ctx, q, result.PublishedContent.Id).Scan(&pc))

		// verify response timestamps round-trip correctly with what was stored
		storedCreatedAt, err := grpcx.DecodeTime(result.PublishedContent.CreatedAt)
		require.NoError(t, err)
		require.WithinDuration(t, storedCreatedAt, pc.CreatedAt, time.Second)

		storedUpdatedAt, err := grpcx.DecodeTime(result.PublishedContent.UpdatedAt)
		require.NoError(t, err)
		require.WithinDuration(t, storedUpdatedAt, pc.UpdatedAt, time.Second)

		storedPublishedAt, err := grpcx.DecodeTime(result.PublishedContent.PublishedAt)
		require.NoError(t, err)
		require.WithinDuration(t, storedPublishedAt, pc.PublishedAt, time.Second)
	})

	t.Run("returns immediately without magnet uri even when torrent exists", func(t *testing.T) {
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
		communityapi.NewHTTPPublished(
			q,
			communityapi.HTTPPublishedOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
			communityapi.HTTPPublishedOptionHTTPClient(&http.Client{}),
			communityapi.HTTPPublishedOptionMediaStorage(fsx.DirVirtual(mediaDir)),
			communityapi.HTTPPublishedOptionTorrentStorage(fsx.DirVirtual(torrentDir)),
		).Bind(routes.PathPrefix("/c").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))
		body, err := json.Marshal(&communityapi.PublishContentRequest{
			PublishedContent: &communityapi.PublishedContent{
				LibraryId: libraryID,
			},
		})
		require.NoError(t, err)

		resp, req, err := httptestx.BuildRequestContextBytes(
			ctx,
			http.MethodPost,
			"/c/"+communityID,
			body,
			httptestx.RequestOptionAuthorization("Bearer "+httpauthtest.UnsafeToken(claims, httpauthtest.UnsafeJWTSecretSource)),
			httptestx.RequestOptionContent("application/json"),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusOK, resp.Code)

		var result communityapi.PublishContentResponse
		require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &result))
		require.NotEmpty(t, result.PublishedContent.Id)
		require.Equal(t, "", result.PublishedContent.MagnetUri)
	})
}
