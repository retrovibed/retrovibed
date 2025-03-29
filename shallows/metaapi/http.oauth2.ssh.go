package metaapi

import (
	"encoding/base64"
	"log"
	"net/http"
	"time"

	"github.com/go-playground/form/v4"
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/retrovibed/retrovibed/authn"
	"github.com/retrovibed/retrovibed/internal/errorsx"
	"github.com/retrovibed/retrovibed/internal/formx"
	"github.com/retrovibed/retrovibed/internal/httpx"
	"github.com/retrovibed/retrovibed/internal/jwtx"
	"github.com/retrovibed/retrovibed/internal/md5x"
	"github.com/retrovibed/retrovibed/internal/slicesx"
	"github.com/retrovibed/retrovibed/internal/sqlx"
	"github.com/retrovibed/retrovibed/internal/sshx"
	"github.com/retrovibed/retrovibed/meta"
	"github.com/retrovibed/retrovibed/meta/identityssh"
	"golang.org/x/crypto/ssh"
	"golang.org/x/time/rate"
)

// allows setting options for the http router.
type SSHOauth2Option func(*HTTPSSHOauth2)

// HTTPOptionQueryer database to insert into.
func SSHOauth2OptionJWTSecret(j jwtx.SecretSource) SSHOauth2Option {
	return func(dst *HTTPSSHOauth2) {
		dst.jwtsecret = j
	}
}

// HTTPAutoConfig automatically generates the sso http service from the environment.
func SSHOauth2AutoConfig(q sqlx.Queryer) (_ *HTTPSSHOauth2, err error) {
	return NewSSHOauth2(
		q,
	), nil
}

// NewHTTP creates the routes for handling sso requests.
func NewSSHOauth2(q sqlx.Queryer, options ...SSHOauth2Option) (svc *HTTPSSHOauth2) {
	svc = &HTTPSSHOauth2{
		q:         q,
		jwtsecret: authn.JWTSecretFromEnv,
		decoder:   formx.NewDecoder(),
	}

	for _, opt := range options {
		opt(svc)
	}

	return svc
}

// HTTP consumes a request and returns a set of oauth2 auth urls
// to the client for authenticating the client.
type HTTPSSHOauth2 struct {
	jwtsecret jwtx.SecretSource
	decoder   *form.Decoder
	q         sqlx.Queryer
}

func (t *HTTPSSHOauth2) Bind(r *mux.Router) {
	r = r.StrictSlash(true)

	r.Methods(http.MethodGet).Path("/auth").Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpx.RouteRateLimited(rate.NewLimiter(rate.Every(20*time.Millisecond), 200)),
		httpx.ParseForm,
		httpx.Timeout1s(),
	).ThenFunc(t.auth))

	r.Methods(http.MethodPost).Path("/token").Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpx.RouteRateLimited(rate.NewLimiter(rate.Every(20*time.Millisecond), 200)),
		httpx.ParseForm,
		httpx.Timeout2s(),
	).ThenFunc(t.token))
}

func (t *HTTPSSHOauth2) auth(w http.ResponseWriter, req *http.Request) {
	// The ssh auth endpoint is a trust of first use. we insert the identity with the
	// public key the first time its encountered. we sign the generated code using that
	// public key.
	type reqmsg struct {
		AccessType   string `json:"access_type"`
		ClientID     string `json:"client_id"`
		ResponseType string `json:"response_type"`
		State        string `json:"state"`
	}

	type reqstate struct {
		PublicKey []byte `json:"pkey"`
		Email     string `json:"email"`
		Display   string `json:"display"`
	}

	type authJSON struct {
		Code  string `json:"code"`
		State string `json:"state"`
	}

	var (
		err      error
		msg      reqmsg
		state    reqstate
		iden     identityssh.Identity
		consumed meta.ConsumedToken
	)

	if err = t.decoder.Decode(&msg, req.Form); err != nil {
		log.Println("unable to decode message", err)
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	if err = jwtx.DecodeJSON(msg.State, &state); err != nil {
		log.Println("unable to decode state", err)
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	pub, err := ssh.ParsePublicKey(state.PublicKey)
	if err != nil {
		log.Println("unable to parse public key", err)
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	iden = identityssh.Identity{
		ID:        sshx.FingerprintMD5(pub),
		PublicKey: sshx.EncodeBase64PublicKey(pub),
		ProfileID: uuid.Nil.String(),
	}

	err = identityssh.IdentityInsertWithDefaults(req.Context(), t.q, iden).Scan(&iden)
	if err != nil {
		log.Println("unable to register public key", err)
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	claims := jwtx.NewJWTClaims(
		iden.ID,
		jwtx.ClaimsOptionExpiration(time.Minute),
		jwtx.ClaimsOptionIssuer(iden.ID),
	)

	if err = meta.ConsumedTokenInsertWithDefaults(req.Context(), t.q, meta.ConsumedTokenString(msg.State, claims.ExpiresAt.Time)).Scan(&consumed); err != nil {
		log.Println("unable to consume otp", err)
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	// Here we sign the code using the identity public key + the jwtsecret.
	// ensures the code is secure against manipulation.
	code, err := jwtx.Signed([]byte(md5x.String(iden.PublicKey+md5x.FormatUUID(md5x.Digest(t.jwtsecret())))), claims)
	if err != nil {
		log.Println("unable to sign code", err)
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	errorsx.Log(httpx.WriteJSON(w, httpx.GetBuffer(req), authJSON{
		Code:  code,
		State: msg.State,
	}))
}

func (t *HTTPSSHOauth2) token(w http.ResponseWriter, req *http.Request) {
	type reqmsg struct {
		GrantType string `json:"grant_type"`
		Username  string `json:"username"`
		Password  string `json:"password"`
		Code      string `json:"code"`
	}

	type tokenJSON struct {
		AccessToken  string `json:"access_token"`
		TokenType    string `json:"token_type"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int32  `json:"expires_in"`
	}

	var (
		err          error
		ok           bool
		msg          reqmsg
		rawsignature []byte
		iden         identityssh.Identity
		p            meta.Profile
		sig          ssh.Signature
	)

	if err = t.decoder.Decode(&msg, req.Form); err != nil {
		log.Println("unable to decode message", err)
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	switch msg.GrantType {
	case "password":
		var (
			password string
		)

		if _, password, ok = req.BasicAuth(); !ok {
			log.Println("unable to password grant failed", err)
			errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
			return
		}

		if err = identityssh.IdentityFindByID(req.Context(), t.q, sqlx.NewNullString(md5x.String(msg.Username))).Scan(&iden); sqlx.ErrNoRows(err) != nil {
			errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
			return
		} else if err != nil {
			log.Println("unable to locate identity", err)
			errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
			return
		}

		pkey, err := sshx.DecodeBase64PublicKey(iden.PublicKey)
		if err != nil {
			log.Println("unable to parse public key", err)
			errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
			return
		}

		if rawsignature, err = base64.RawURLEncoding.DecodeString(msg.Password); err != nil {
			log.Println("unable to decode password", err)
			errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
			return
		}

		if err = ssh.Unmarshal(rawsignature, &sig); err != nil {
			log.Println("unable to decode signature", err)
			errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
			return
		}

		if err = pkey.Verify(uuid.FromStringOrNil(password).Bytes(), &sig); err != nil {
			log.Println("unable to verifcation failure", err)
			errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
			return
		}
	case "authorization_code":
		var claims jwt.RegisteredClaims
		// we intentionally ignore the validation error here.
		// because we want the subject id so we can properly validate
		// using the public key.
		_ = jwtx.Validate(t.jwtsecret, msg.Code, &claims)

		if err = identityssh.IdentityFindByID(req.Context(), t.q, sqlx.NewNullString(claims.Subject)).Scan(&iden); sqlx.ErrNoRows(err) != nil {
			log.Println("unable to locate identity, no rows", claims.Subject)
			errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
			return
		} else if err != nil {
			log.Println("unable to locate identity", err)
			errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
			return
		}

		if err = jwtx.Validate(func() []byte {
			return []byte(md5x.String(iden.PublicKey + md5x.FormatUUID(md5x.Digest(t.jwtsecret()))))
		}, msg.Code, &claims); err != nil {
			errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
			return
		}
	case "refresh_token":
		var (
			claims jwt.RegisteredClaims
		)

		if err = jwtx.Validate(t.jwtsecret, req.Form.Get("refresh_token"), &claims); err != nil {
			log.Println("unable to validate bearer", req.Form.Get("refresh_token"), err)
			errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
			return
		}

		if err = identityssh.IdentityFindByID(req.Context(), t.q, sqlx.NewNullString(slicesx.LastOrZero(claims.Audience...))).Scan(&iden); sqlx.ErrNoRows(err) != nil {
			log.Println("unable to find identity", slicesx.LastOrZero(claims.Audience...), err)
			errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
			return
		} else if err != nil {
			log.Println("unable to locate identity", slicesx.LastOrZero(claims.Audience...), err)
			errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
			return
		}
	default:
		log.Println("unknown grant type", msg.GrantType)
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	if uuid.FromStringOrNil(iden.ProfileID).String() == uuid.Nil.String() {
		log.Println("no profile associated with public key", iden.ID)
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	if err = meta.ProfileFindByID(req.Context(), t.q, iden.ProfileID).Scan(&p); err != nil {
		log.Println("unable to locate profile", err)
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	// generate refresh token.
	rst := ""
	if req.Form.Get("access_type") == "offline" {
		rst, err = jwtx.Signed(t.jwtsecret(), jwtx.NewJWTClaims(
			p.ID,
			jwtx.ClaimsOptionAuthnRefreshExpiration(),
			jwtx.ClaimsOptionIssuer(iden.ID),
			// jwtx.ClaimsOptionAudience(),
		))

		if err != nil {
			log.Println("unable to generate refresh token", err)
			errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
			return
		}
	}

	aclaims := jwtx.NewJWTClaims(
		p.ID,
		jwtx.ClaimsOptionAuthnExpiration(),
		jwtx.ClaimsOptionIssuer(iden.ID),
		// jwtx.ClaimsOptionAudience(),
	)
	ast, err := jwtx.Signed(t.jwtsecret(), aclaims)

	if err != nil {
		log.Println("unable to generate access token", err)
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	errorsx.Log(httpx.WriteJSON(w, httpx.GetBuffer(req), tokenJSON{
		AccessToken:  ast,
		RefreshToken: rst,
		TokenType:    "BEARER",
		ExpiresIn:    int32(time.Until(aclaims.ExpiresAt.Time) / time.Second),
	}))
}
