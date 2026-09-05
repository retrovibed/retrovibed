package metaapi_test

import (
	"encoding/base64"
	"encoding/json"
	"strings"
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/shallows/meta"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
	"github.com/stretchr/testify/require"
)

// Token carries the registered claims under their RFC 7519 abbreviations, and
// nothing but the proto field names enforces that now protojson is gone --
// protoc-gen-go derives the encoding/json tag from the proto field name, so
// renaming Jti back to Id would silently emit "id" and break every consumer
// that validates the bearer.
func TestTokenClaimsWireNames(t *testing.T) {
	signed, err := jwtx.Signed(
		[]byte("f3fd7b1d0a4f4c6ea1c1c67f2a2ab61f"),
		metaapi.NewJWTClaim(metaapi.TokenFromRegisterClaims(
			jwtx.NewJWTClaims("0a484599-ff07-72bc-4114-13ef3ace4786", jwtx.ClaimsOptionIssuer("72bc4114")),
			metaapi.TokenOptionFromAuthz(meta.Authz{LibraryRead: true, LocalOnly: true}),
		)),
	)
	require.NoError(t, err)

	segments := strings.Split(signed, ".")
	require.Len(t, segments, 3)

	payload, err := base64.RawURLEncoding.DecodeString(segments[1])
	require.NoError(t, err)

	claims := make(map[string]any)
	require.NoError(t, json.Unmarshal(payload, &claims))

	require.Contains(t, claims, "jti")
	require.Contains(t, claims, "iss")
	require.Contains(t, claims, "sub")
	require.Contains(t, claims, "iat")
	require.Contains(t, claims, "exp")
	require.Contains(t, claims, "nbf")
	require.Equal(t, "0a484599-ff07-72bc-4114-13ef3ace4786", claims["sub"])
	require.Equal(t, "72bc4114", claims["iss"])

	for _, long := range []string{"id", "issuer", "profile_id", "session_id", "issued", "expires", "not_before"} {
		require.NotContains(t, claims, long)
	}
}
