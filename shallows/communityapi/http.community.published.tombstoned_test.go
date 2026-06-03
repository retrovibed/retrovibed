package communityapi_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/mux"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/shallows/community"
	"github.com/retrovibed/retrovibed/shallows/communityapi"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/httptestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/testx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/retrovibed/retrovibed/shallows/meta"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"github.com/stretchr/testify/require"
)

func TestTombstonedEndpoint(t *testing.T) {
	t.Run("returns 404 when published content does not exist", func(t *testing.T) {
		var (
			ctx, done     = testx.Context(t)
			q             = sqltestx.Metadatabase(t)
			p             meta.Profile
			v             meta.Authz
			nonExistentID = uuid.Must(uuid.NewV7()).String()
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
			communityapi.HTTPPublishedOptionMediaStorage(fsx.DirVirtual(t.TempDir())),
			communityapi.HTTPPublishedOptionTorrentStorage(fsx.DirVirtual(t.TempDir())),
		).Bind(routes.PathPrefix("/c").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))
		resp, req, err := httptestx.BuildRequestContextBytes(
			ctx,
			http.MethodDelete,
			"/c/"+nonExistentID,
			nil,
			httptestx.RequestOptionAuthorization("Bearer "+httpauthtest.UnsafeToken(claims, httpauthtest.UnsafeJWTSecretSource)),
			httptestx.RequestOptionContent("application/json"),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusNotFound, resp.Code)
	})

	t.Run("deletes published content and returns it in response", func(t *testing.T) {
		var (
			ctx, done   = testx.Context(t)
			q           = sqltestx.Metadatabase(t)
			p           meta.Profile
			v           meta.Authz
			communityID = uuid.Must(uuid.NewV7()).String()
		)
		defer done()

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, testx.Fake(&v, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, v).Scan(&v))

		lmd := library.Metadata{
			ID:             uuid.Must(uuid.NewV7()).String(),
			Description:    "test media",
			Bytes:          1024,
			TorrentID:      uuid.Nil.String(),
			KnownMediaID:   uuid.Nil.String(),
			ArchiveID:      uuid.Nil.String(),
			EncryptionSeed: uuid.Must(uuid.NewV4()).String(),
		}
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, lmd).Scan(&lmd))

		pc := community.NewPublishedContent(community.PublishedContent{
			CommunityID: communityID,
			LibraryID:   lmd.ID,
			Bytes:       lmd.Bytes,
		})
		require.NoError(t, community.PublishedContentInsertWithDefaults(ctx, q, pc).Scan(&pc))

		routes := mux.NewRouter()
		communityapi.NewHTTPPublished(
			q,
			communityapi.HTTPPublishedOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
			communityapi.HTTPPublishedOptionMediaStorage(fsx.DirVirtual(t.TempDir())),
			communityapi.HTTPPublishedOptionTorrentStorage(fsx.DirVirtual(t.TempDir())),
		).Bind(routes.PathPrefix("/c").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))
		resp, req, err := httptestx.BuildRequestContextBytes(
			ctx,
			http.MethodDelete,
			"/c/"+pc.ID,
			nil,
			httptestx.RequestOptionAuthorization("Bearer "+httpauthtest.UnsafeToken(claims, httpauthtest.UnsafeJWTSecretSource)),
			httptestx.RequestOptionContent("application/json"),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusOK, resp.Code)

		var result communityapi.PublishContentDeleteResponse
		require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &result))
		require.NotNil(t, result.PublishedContent)
		require.Equal(t, pc.ID, result.PublishedContent.Id)
		require.Equal(t, lmd.ID, result.PublishedContent.LibraryId)

		var deleted community.PublishedContent
		require.Error(t, community.PublishedContentFindByID(ctx, q, pc.ID).Scan(&deleted))
	})

	t.Run("deletes one of multiple published content items leaving others intact", func(t *testing.T) {
		var (
			ctx, done   = testx.Context(t)
			q           = sqltestx.Metadatabase(t)
			p           meta.Profile
			v           meta.Authz
			communityID = uuid.Must(uuid.NewV7()).String()
		)
		defer done()

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.NoError(t, testx.Fake(&v, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, v).Scan(&v))

		insertPC := func(libraryID string) community.PublishedContent {
			lmd := library.Metadata{
				ID:             libraryID,
				Description:    "test media " + libraryID,
				Bytes:          512,
				TorrentID:      uuid.Nil.String(),
				KnownMediaID:   uuid.Nil.String(),
				ArchiveID:      uuid.Nil.String(),
				EncryptionSeed: uuid.Must(uuid.NewV4()).String(),
			}
			require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, lmd).Scan(&lmd))
			pc := community.NewPublishedContent(community.PublishedContent{
				CommunityID: communityID,
				LibraryID:   lmd.ID,
				Bytes:       lmd.Bytes,
			})
			require.NoError(t, community.PublishedContentInsertWithDefaults(ctx, q, pc).Scan(&pc))
			return pc
		}

		pc1 := insertPC(uuid.Must(uuid.NewV7()).String())
		pc2 := insertPC(uuid.Must(uuid.NewV7()).String())
		require.NotEmpty(t, pc1.ID)
		require.NotEmpty(t, pc2.ID)

		routes := mux.NewRouter()
		communityapi.NewHTTPPublished(
			q,
			communityapi.HTTPPublishedOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
			communityapi.HTTPPublishedOptionMediaStorage(fsx.DirVirtual(t.TempDir())),
			communityapi.HTTPPublishedOptionTorrentStorage(fsx.DirVirtual(t.TempDir())),
		).Bind(routes.PathPrefix("/c").Subrouter())

		claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v)))
		resp, req, err := httptestx.BuildRequestContextBytes(
			ctx,
			http.MethodDelete,
			"/c/"+pc1.ID,
			nil,
			httptestx.RequestOptionAuthorization("Bearer "+httpauthtest.UnsafeToken(claims, httpauthtest.UnsafeJWTSecretSource)),
			httptestx.RequestOptionContent("application/json"),
		)
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)
		require.Equal(t, http.StatusOK, resp.Code)

		var deleted community.PublishedContent
		require.Error(t, community.PublishedContentFindByID(ctx, q, pc1.ID).Scan(&deleted))

		var remaining community.PublishedContent
		require.NoError(t, community.PublishedContentFindByID(ctx, q, pc2.ID).Scan(&remaining))
		require.Equal(t, pc2.ID, remaining.ID)
	})
}
