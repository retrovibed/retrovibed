package media_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/gorilla/mux"
	"github.com/retrovibed/retrovibed/httpauthtest"
	"github.com/retrovibed/retrovibed/internal/asyncx"
	"github.com/retrovibed/retrovibed/internal/httptestx"
	"github.com/retrovibed/retrovibed/internal/jwtx"
	"github.com/retrovibed/retrovibed/internal/langx"
	"github.com/retrovibed/retrovibed/internal/sqltestx"
	"github.com/retrovibed/retrovibed/internal/testx"
	"github.com/retrovibed/retrovibed/library"
	"github.com/retrovibed/retrovibed/media"
	"github.com/retrovibed/retrovibed/meta"
	"github.com/retrovibed/retrovibed/metaapi"
	"github.com/stretchr/testify/require"
)

func TestMediaUpdate(t *testing.T) {
	var (
		p      meta.Profile
		authz  meta.Authz
		v      library.Metadata
		result media.MediaUpdateResponse
	)
	ctx, done := testx.Context(t)
	defer done()

	q := sqltestx.Metadatabase(t)

	require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
	require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
	require.NoError(t, testx.Fake(&authz, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
	require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, authz).Scan(&authz))
	require.NoError(t, testx.Fake(&v, library.MetadataOptionTestDefaults))
	require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, v).Scan(&v))

	v.KnownMediaID = "2a945256-140e-430e-9128-5e17e29d9a28"

	routes := mux.NewRouter()
	media.NewHTTPLibrary(
		q,
		asyncx.NewWakeup(t.Context()),
		nil,
		nil,
		media.HTTPLibraryOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
	).Bind(routes.PathPrefix("/").Subrouter())

	claims := metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(authz)))

	encoded, err := json.Marshal(media.MediaUpdateRequest{
		Media: langx.Autoptr(langx.Clone(media.Media{}, media.MediaOptionFromLibraryMetadata(v))),
	})
	require.NoError(t, err)

	resp, req, err := httptestx.BuildRequestBytes(
		http.MethodPost,
		fmt.Sprintf("/%s", v.ID),
		encoded,
		httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(claims, httpauthtest.UnsafeJWTSecretSource)),
	)
	require.NoError(t, err)

	routes.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Result().StatusCode)
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
	require.Equal(t, v.ID, result.Media.Id)
	require.Equal(t, v.KnownMediaID, result.Media.KnownMediaId)
}
