package community_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/mux"
	"github.com/retrovibed/retrovibed/community"
	"github.com/retrovibed/retrovibed/httpauthtest"
	"github.com/retrovibed/retrovibed/internal/fsx"
	"github.com/retrovibed/retrovibed/internal/httptestx"
	"github.com/retrovibed/retrovibed/internal/jwtx"
	"github.com/retrovibed/retrovibed/internal/sqltestx"
	"github.com/retrovibed/retrovibed/internal/testx"
	"github.com/retrovibed/retrovibed/library"
	"github.com/retrovibed/retrovibed/meta"
	"github.com/retrovibed/retrovibed/metaapi"
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
			EncryptionSeed: uuid.Must(uuid.NewV4()).String(),
		}
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, lmd2).Scan(&lmd2))

		routes := mux.NewRouter()
		community.NewHTTP(
			q,
			community.HTTPOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
			community.HTTPOptionMediaStorage(fsx.DirVirtual(mediaDir)),
			community.HTTPOptionTorrentStorage(fsx.DirVirtual(torrentDir)),
		).Bind(routes.PathPrefix("/c").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))
		token := "Bearer " + httpauthtest.UnsafeToken(claims, httpauthtest.UnsafeJWTSecretSource)

		// Publish two items to the same community.
		for _, lid := range []string{lmd1.ID, lmd2.ID} {
			body, err := json.Marshal(&meta.PublishContentRequest{
				PublishedContent: &meta.PublishedContent{LibraryId: lid},
			})
			require.NoError(t, err)

			resp, req, err := httptestx.BuildRequestContextBytes(
				ctx,
				http.MethodPost,
				"/c/"+communityID+"/publish",
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
			"/c/"+communityID+"/published",
			nil,
			httptestx.RequestOptionAuthorization(token),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusOK, resp.Code)

		var result meta.PublishedContentListResponse
		require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &result))
		require.Len(t, result.Items, 2)

		ids := []string{result.Items[0].LibraryId, result.Items[1].LibraryId}
		require.Contains(t, ids, lmd1.ID)
		require.Contains(t, ids, lmd2.ID)
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
		community.NewHTTP(
			q,
			community.HTTPOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
			community.HTTPOptionMediaStorage(fsx.DirVirtual(mediaDir)),
			community.HTTPOptionTorrentStorage(fsx.DirVirtual(torrentDir)),
		).Bind(routes.PathPrefix("/c").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))
		token := "Bearer " + httpauthtest.UnsafeToken(claims, httpauthtest.UnsafeJWTSecretSource)

		resp, req, err := httptestx.BuildRequestContextBytes(
			ctx,
			http.MethodGet,
			"/c/"+communityID+"/published",
			nil,
			httptestx.RequestOptionAuthorization(token),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusOK, resp.Code)

		var result meta.PublishedContentListResponse
		require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &result))
		require.Empty(t, result.Items)
	})
}
