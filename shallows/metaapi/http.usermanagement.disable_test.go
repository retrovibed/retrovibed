package metaapi_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/httptestx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/testx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/meta"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"github.com/stretchr/testify/require"
)

func TestHTTPUserManagementDisable(t *testing.T) {
	t.Run("should disable the profile", func(t *testing.T) {
		var (
			p      meta.Profile
			v      meta.Authz
			result metaapi.ProfileUpdateRequest
		)

		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults, timex.UTCEncodeOption))
		require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
		require.EqualValues(t, p.DisabledAt, timex.Inf())
		require.NoError(t, testx.Fake(&v, meta.AuthzOptionProfileID(p.ID), meta.AuthzOptionAdmin))
		require.NoError(t, meta.AuthzInsertWithDefaults(ctx, q, v).Scan(&v))

		routes := mux.NewRouter()
		metaapi.NewHTTPUsermanagement(
			q,
			metaapi.HTTPUsermanagementOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/u12t").Subrouter())

		mut := p
		mut.DisabledManuallyAt = time.Now()
		mut = langx.Clone(mut, timex.JSONSafeEncodeOption, timex.UTCEncodeOption)
		encoded, err := metaapi.NewProfileFromMetaProfile(mut)
		require.NoError(t, err)
		b, err := json.Marshal(&metaapi.ProfileDisableRequest{
			Profile: encoded,
		})
		require.NoError(t, err)

		token := httpauthtest.UnsafeClaimsToken(metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), metaapi.TokenOptionFromAuthz(v))), httpauthtest.UnsafeJWTSecretSource)

		resp, req, err := httptestx.BuildRequestContext(ctx, http.MethodDelete, fmt.Sprintf("/u12t/%s", mut.ID), bytes.NewReader(b), httptestx.RequestOptionAuthorization(token))
		require.NoError(t, err)

		routes.ServeHTTP(resp, req)
		require.NoError(t, httpx.ErrorCode(resp.Result()))
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))

		require.NotEqualValues(t, timex.Inf(), result.Profile.DisabledManuallyAt)
		encoded.UpdatedAt = result.Profile.UpdatedAt
		result.Profile.DisabledAt = encoded.DisabledAt
		result.Profile.DisabledManuallyAt = encoded.DisabledManuallyAt
		result.Profile.DisabledPendingApprovalAt = encoded.DisabledPendingApprovalAt
		require.EqualValues(t, encoded, result.Profile)
	})
}
