package metaapi_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/retrovibed/retrovibed/authn"
	"github.com/retrovibed/retrovibed/httpauthtest"
	"github.com/retrovibed/retrovibed/internal/httpx"
	"github.com/retrovibed/retrovibed/internal/jwtx"
	"github.com/retrovibed/retrovibed/internal/sqltestx"
	"github.com/retrovibed/retrovibed/internal/sshx"
	"github.com/retrovibed/retrovibed/internal/testx"
	"github.com/retrovibed/retrovibed/meta"
	"github.com/retrovibed/retrovibed/meta/identityssh"
	"github.com/retrovibed/retrovibed/metaapi"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
)

func TestHTTPSSHOauth2WithValidIdentiy(t *testing.T) {
	var (
		p      meta.Profile
		iden   identityssh.Identity
		claims jwt.RegisteredClaims
	)

	ctx, done := testx.Context(t)
	defer done()

	q := sqltestx.Metadatabase(t)
	defer q.Close()

	signer, err := sshx.SignerFromGenerator(sshx.UnsafeNewKeyGen())
	require.NoError(t, err)

	require.NoError(t, testx.Fake(&p, meta.ProfileOptionTestDefaults))
	require.NoError(t, meta.ProfileInsertWithDefaults(ctx, q, p).Scan(&p))
	require.NoError(t, identityssh.IdentityInsertWithDefaults(ctx, q, identityssh.Identity{
		ID:        sshx.FingerprintMD5(signer.PublicKey()),
		ProfileID: p.ID,
		PublicKey: sshx.EncodeBase64PublicKey(signer.PublicKey()),
	}).Scan(&iden))

	routes := mux.NewRouter()
	routes.Use(httpx.RouteInvoked)

	metaapi.NewSSHOauth2(
		q,
		metaapi.SSHOauth2OptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
	).Bind(routes.PathPrefix("/oauth2/ssh").Subrouter())

	s := httptest.NewServer(routes)
	defer s.Close()

	state, err := authn.AutoTokenState(signer)
	require.NoError(t, err)

	endpoint := authn.EndpointSSHAuth(s.URL)

	cfg := authn.OAuth2SSHConfig(signer, "", endpoint)

	authzuri := cfg.AuthCodeURL(
		state,
		oauth2.AccessTypeOffline,
	)

	c := &http.Client{}
	r, err := authn.RetrieveAuthCode(ctx, c, authzuri)
	require.NoError(t, err)
	require.Equal(t, state, r.State)

	token, err := cfg.Exchange(ctx, r.Code, oauth2.AccessTypeOffline)
	require.NoError(t, err)

	require.NotEqual(t, "", token.AccessToken)          // should be an access token
	require.NotEqual(t, "", token.RefreshToken)         // should be an refresh token
	require.Less(t, int64(0), time.Until(token.Expiry)) // should have an expiration

	require.NoError(t, jwtx.Validate(httpauthtest.UnsafeJWTSecretSource, token.AccessToken, &claims))

	require.Equal(t, p.ID, claims.Subject)
	require.Equal(t, iden.ID, claims.Issuer)
	require.True(t, claims.ExpiresAt.After(time.Now()))
	require.True(t, claims.IssuedAt.Before(time.Now()))
}

func TestHTTPSSHOauth2WithUnknownPublicKey(t *testing.T) {
	ctx, done := testx.Context(t)
	defer done()

	q := sqltestx.Metadatabase(t)
	defer q.Close()

	signer, err := sshx.SignerFromGenerator(sshx.UnsafeNewKeyGen())
	require.NoError(t, err)

	routes := mux.NewRouter()
	routes.Use(httpx.RouteInvoked)

	metaapi.NewSSHOauth2(
		q,
		metaapi.SSHOauth2OptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
	).Bind(routes.PathPrefix("/oauth2/ssh").Subrouter())

	s := httptest.NewServer(routes)
	defer s.Close()

	state, err := authn.AutoTokenState(signer)
	require.NoError(t, err)

	endpoint := authn.EndpointSSHAuth(s.URL)

	cfg := authn.OAuth2SSHConfig(signer, "", endpoint)

	authzuri := cfg.AuthCodeURL(
		state,
		oauth2.AccessTypeOffline,
	)

	c := &http.Client{}
	r, err := authn.RetrieveAuthCode(ctx, c, authzuri)
	require.NoError(t, err)
	require.Equal(t, state, r.State)

	_, err = cfg.Exchange(ctx, r.Code, oauth2.AccessTypeOffline)
	require.Error(t, err)
}
