package metaapi_test

import (
	"encoding/json"
	"net/url"
	"testing"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/retrovibed/retrovibed/authn"
	"github.com/retrovibed/retrovibed/httpauthtest"
	"github.com/retrovibed/retrovibed/internal/formx"
	"github.com/retrovibed/retrovibed/internal/httptestx"
	"github.com/retrovibed/retrovibed/internal/httpx"
	"github.com/retrovibed/retrovibed/internal/jwtx"
	"github.com/retrovibed/retrovibed/internal/sqltestx"
	"github.com/retrovibed/retrovibed/internal/sqlx"
	"github.com/retrovibed/retrovibed/internal/testx"
	"github.com/retrovibed/retrovibed/meta"
	"github.com/retrovibed/retrovibed/meta/identityssh"
	"github.com/retrovibed/retrovibed/metaapi"
	"github.com/stretchr/testify/require"
)

func TestHTTPSSO(t *testing.T) {
	type tokenJSON struct {
		AccessToken  string `json:"access_token"`
		TokenType    string `json:"token_type"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int32  `json:"expires_in"`
	}

	t.Run("register", func(t *testing.T) {
		t.Run("create a profile in the pending state", func(t *testing.T) {
			var (
				iden   identityssh.Identity
				result tokenJSON
				claims jwt.RegisteredClaims
				p      meta.Profile
			)

			q := sqltestx.Metadatabase(t)

			routes := mux.NewRouter()
			metaapi.NewHTTP(
				q,
				metaapi.HTTPOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
			).Bind(routes)

			require.NoError(
				t,
				testx.Fake(
					&iden,
					identityssh.OptionTestDefaults,
				),
			)
			require.NoError(
				t,
				identityssh.IdentityInsertZeroWithDefaults(t.Context(), q, iden).Scan(&iden),
			)

			token := httpauthtest.UnsafeClaimsToken(
				jwtx.NewJWTClaims(
					uuid.Nil.String(),
					jwtx.ClaimsOptionAuthnExpiration(),
					jwtx.ClaimsOptionID(uuid.Nil.String()),
					jwtx.ClaimsOptionIssuer(iden.ID),
				),
				authn.JWTRegistrationSecretFromEnv,
			)

			v, err := formx.NewEncoder().Encode(metaapi.Identity{Display: "derp"})
			require.NoError(t, err)

			w, r, err := httptestx.BuildPostRequest(
				nil,
				httptestx.RequestOptionURL(&url.URL{Path: "/register", RawQuery: v.Encode()}),
				httptestx.RequestOptionAuthorization(
					token,
				),
			)
			require.NoError(t, err)

			routes.ServeHTTP(w, r)

			require.NoError(t, httpx.ErrorCode(w.Result()))
			require.NoError(t, json.NewDecoder(w.Body).Decode(&result))

			_, err = jwt.ParseWithClaims(result.AccessToken, &claims, func(t *jwt.Token) (interface{}, error) {
				return authn.JWTRegistrationSecretFromEnv(), nil
			})
			require.NoError(t, err)

			require.WithinRange(t, claims.NotBefore.Time, time.Now().Add(-1*time.Minute), time.Now().Add(time.Hour))
			require.WithinRange(t, claims.ExpiresAt.Time, time.Now().Add(-1*time.Minute), time.Now().Add(7*24*time.Hour))
			require.NotEqual(t, uuid.Nil.String(), claims.Subject)
			require.NoError(t, meta.ProfileFindByID(t.Context(), q, claims.Subject).Scan(&p))
			require.Equal(t, p.Display, "derp")
			require.NoError(t, identityssh.IdentityFindByID(t.Context(), q, sqlx.NewNullString(iden.ID)).Scan(&iden))
			require.Equal(t, p.ID, iden.ProfileID)
		})
	})
}
