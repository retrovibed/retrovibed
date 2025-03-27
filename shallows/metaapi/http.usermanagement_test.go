package metaapi_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/retrovibed/retrovibed/cmd/cmdmeta"
	"github.com/retrovibed/retrovibed/httpauthtest"
	"github.com/retrovibed/retrovibed/internal/formx"
	"github.com/retrovibed/retrovibed/internal/httptestx"
	"github.com/retrovibed/retrovibed/internal/httpx"
	"github.com/retrovibed/retrovibed/internal/jwtx"
	"github.com/retrovibed/retrovibed/internal/sqltestx"
	"github.com/retrovibed/retrovibed/internal/testx"
	"github.com/retrovibed/retrovibed/meta"
	"github.com/retrovibed/retrovibed/metaapi"
	"github.com/stretchr/testify/require"
)

func TestHTTPUserManagementSearch(t *testing.T) {
	var (
		p      meta.Profile
		result metaapi.ProfileSearchResponse
		claims jwt.RegisteredClaims
	)

	ctx, done := testx.Context(t)
	defer done()

	q := sqltestx.Metadatabase(t)
	defer q.Close()

	require.NoError(t, cmdmeta.InitializeDatabase(ctx, q))

	require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
	require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
	routes := mux.NewRouter()

	metaapi.NewHTTPUsermanagement(
		q,
		metaapi.HTTPUsermanagementOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
	).Bind(routes.PathPrefix("/").Subrouter())

	b := testx.Must(formx.NewEncoder().Encode(&metaapi.ProfileSearchRequest{
		Offset: 0,
	}))(t)

	claims = jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration())

	resp, req, err := httptestx.BuildRequest(http.MethodGet, "/?"+b.Encode(), nil, httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(&claims, httpauthtest.UnsafeJWTSecretSource)))
	require.NoError(t, err)

	routes.ServeHTTP(resp, req)

	require.NoError(t, httpx.ErrorCode(resp.Result()))
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))

	encoded := testx.Must(metaapi.NewProfileFromMetaProfile(p))(t)
	require.Equal(t, result.Next.Offset, uint64(1))
	require.Contains(t, result.Items, encoded)
}

func TestHTTPUserManagementFind(t *testing.T) {
	var (
		p      meta.Profile
		result metaapi.ProfileLookupResponse
	)

	ctx, done := testx.Context(t)
	defer done()

	q := sqltestx.Metadatabase(t)
	defer q.Close()

	require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults, meta.ProfileOptionTimezoneUTC))
	require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))

	routes := mux.NewRouter()
	metaapi.NewHTTPUsermanagement(
		q,
		metaapi.HTTPUsermanagementOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
	).Bind(routes.PathPrefix("/").Subrouter())

	mut := p

	token := httpauthtest.UnsafeClaimsToken(jwtx.NewJWTClaims(p.ID, jwtx.ClaimsOptionAuthnExpiration()), httpauthtest.UnsafeJWTSecretSource)
	resp, req, err := httptestx.BuildRequest(http.MethodGet, fmt.Sprintf("/%s", mut.ID), nil, httptestx.RequestOptionAuthorization(token))
	require.NoError(t, err)
	routes.ServeHTTP(resp, req)

	require.NoError(t, httpx.ErrorCode(resp.Result()))
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))

	require.Equal(t, result.Profile.Id, p.ID)
}

// var _ = Describe("HTTPUsermanagement", func() {
// 	Describe("update", func() {
// 		It("should update the profile", func(ctx context.Context) {
// 			var (
// 				p      profiles.Profile
// 				result meta.UpdateRequest
// 				claims meta.Token
// 			)
// 			Expect(testx.Fake(&p, profiles.OptionTestDefaults, profiles.OptionTimezoneUTC)).Should(Succeed())
// 			Expect(profiles.ProfileInsertWithDefaults(ctx, testx.TX, p).Scan(&p)).Should(Succeed())

// 			routes := mux.NewRouter()
// 			meta.NewHTTPUsermanagement(
// 				testx.TX,
// 				meta.HTTPUsermanagementOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
// 			).Bind(routes.PathPrefix("/u12t").Subrouter())

// 			Expect(p.DisabledManuallyAt).To(Equal(timex.Inf()))
// 			mut := p
// 			mut.DisabledManuallyAt = time.Now()
// 			mut = langx.Clone(mut, profiles.OptionTimezoneUTC, profiles.OptionJSONSafeEncode)
// 			encoded := (&meta.Profile{}).FromProfile(&mut)
// 			b, err := json.Marshal(&meta.UpdateRequest{
// 				Profile: encoded,
// 			})
// 			Expect(err).Should(Succeed())

// 			Expect(testx.Fake(&claims, meta.TokenOptionTestAuthorized)).To(Succeed())

// 			resp, req, err := httptestx.BuildRequest(http.MethodPatch, fmt.Sprintf("/u12t/%s", mut.ID), b, httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(&claims)))
// 			Expect(err).Should(Succeed())

// 			routes.ServeHTTP(resp, req)
// 			Expect(httpx.ErrorCode(resp.Result())).Should(Succeed())
// 			Expect(json.NewDecoder(resp.Body).Decode(&result)).Should(Succeed())

// 			Expect(result.Profile.DisabledManuallyAt).ToNot(Equal(timex.Inf()))
// 			result.Profile.DisabledAt = encoded.DisabledAt
// 			result.Profile.DisabledManuallyAt = encoded.DisabledManuallyAt
// 			result.Profile.DisabledPendingApprovalAt = encoded.DisabledPendingApprovalAt
// 			Expect(result.Profile).To(Equal(encoded))
// 		})
// 	})

// 	Describe("disable", func() {
// 		It("should return the disabled profile", func(ctx context.Context) {
// 			var (
// 				p      profiles.Profile
// 				result meta.DisableResponse
// 				claims meta.Token
// 			)
// 			Expect(testx.Fake(&p, profiles.OptionTestDefaults, profiles.OptionTimezoneUTC)).Should(Succeed())
// 			Expect(profiles.ProfileInsertWithDefaults(ctx, testx.TX, p).Scan(&p)).Should(Succeed())

// 			routes := mux.NewRouter()
// 			meta.NewHTTPUsermanagement(
// 				testx.TX,
// 				meta.HTTPUsermanagementOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
// 			).Bind(routes.PathPrefix("/u12t").Subrouter())

// 			Expect(testx.Fake(&claims, meta.TokenOptionTestAuthorized)).To(Succeed())
// 			encoded := (&meta.Profile{}).FromProfile(&p)
// 			b, err := json.Marshal(&meta.UpdateRequest{
// 				Profile: encoded,
// 			})
// 			Expect(err).Should(Succeed())

// 			resp, req, err := httptestx.BuildRequest(http.MethodDelete, fmt.Sprintf("/u12t/%s", p.ID), b, httptestx.RequestOptionAuthorization(httpauthtest.UnsafeClaimsToken(&claims)))
// 			Expect(err).Should(Succeed())
// 			routes.ServeHTTP(resp, req)

// 			Expect(httpx.ErrorCode(resp.Result())).Should(Succeed())
// 			Expect(json.NewDecoder(resp.Body).Decode(&result)).Should(Succeed())
// 			p = langx.Clone(p, profiles.OptionJSONSafeEncode, profiles.OptionTimezoneUTC)
// 			Expect(result.Profile.Id).To(Equal(p.ID))
// 			Expect(result.Profile.AccountId).To(Equal(p.AccountID))
// 		})
// 	})
// })
