package cmdtestx

import (
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
)

func NewTLSServer(t *testing.T, q sqlx.Queryer, routes *mux.Router) *httptest.Server {
	t.Helper()

	metamux := routes.PathPrefix("/meta").Subrouter()
	metaapi.NewHTTPAuthz(
		q,
		metaapi.HTTPAuthzOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
	).Bind(metamux.PathPrefix("/authz").Subrouter())

	oauth2mux := routes.PathPrefix("/oauth2").Subrouter()
	metaapi.NewSSHOauth2(
		q,
		metaapi.SSHOauth2OptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
	).Bind(oauth2mux.PathPrefix("/ssh").Subrouter())

	srv := httptest.NewTLSServer(routes)
	t.Cleanup(srv.Close)
	return srv
}
