package metaapi

import (
	"log"
	"net/http"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-playground/form/v4"
	"github.com/gofrs/uuid/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/retrovibed/retrovibed/authn"
	"github.com/retrovibed/retrovibed/httpauth"
	"github.com/retrovibed/retrovibed/internal/duckdbx"
	"github.com/retrovibed/retrovibed/internal/errorsx"
	"github.com/retrovibed/retrovibed/internal/formx"
	"github.com/retrovibed/retrovibed/internal/httpx"
	"github.com/retrovibed/retrovibed/internal/jwtx"
	"github.com/retrovibed/retrovibed/internal/langx"
	"github.com/retrovibed/retrovibed/internal/sqlx"
	"github.com/retrovibed/retrovibed/meta"
	"github.com/retrovibed/retrovibed/meta/identityssh"
	"golang.org/x/time/rate"
)

// HTTPOption allows setting options for the http router.
type HTTPOption func(*HTTPSSO)

// HTTPOptionJWTSecret set the jwt secret
func HTTPOptionJWTSecret(v jwtx.SecretSource) HTTPOption {
	return func(dst *HTTPSSO) {
		dst.jwtsecret = v
	}
}

// NewHTTP creates the routes for handling sso requests.
func NewHTTP(q sqlx.Queryer, options ...HTTPOption) (svc *HTTPSSO) {
	svc = &HTTPSSO{
		q:         q,
		jwtsecret: authn.JWTSecretFromEnv,
		decoder:   formx.NewDecoder(),
	}

	for _, opt := range options {
		opt(svc)
	}

	return svc
}

// HTTPSSO consumes a request and returns a set of oauth2 auth urls
// to the client for authenticating the client.
type HTTPSSO struct {
	decoder   *form.Decoder
	q         sqlx.Queryer
	jwtsecret jwtx.SecretSource
}

func (t *HTTPSSO) Bind(r *mux.Router) {
	r = r.StrictSlash(true)
	r.Use(httpx.RouteInvoked)
	// r.Use(httpx.DebugRequest)

	r.Methods(http.MethodGet).Path("/").Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpx.RouteRateLimited(rate.NewLimiter(rate.Every(time.Second/20), 20)),
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.Timeout4s(),
	).ThenFunc(t.current))

	r.Methods(http.MethodPost).Path("/register").Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpx.Timeout4s(),
		httpx.RouteRateLimited(rate.NewLimiter(rate.Every(20*time.Millisecond), 200)),
		httpx.ParseForm,
	).ThenFunc(t.register))
}

func (t *HTTPSSO) current(w http.ResponseWriter, r *http.Request) {
	var (
		err    error
		pid    string
		bearer string
		a      meta.Authz
		p      meta.Profile
		mp     *Profile
	)

	if _, pid, err = httpauth.IssuerSubjectID(r.Context(), t.jwtsecret, r); err != nil {
		httpx.Unauthorized(w, errorsx.Wrap(err, "unable to decode authentication token"))
		return
	}

	if err = meta.ProfileFindByID(r.Context(), t.q, pid).Scan(&p); err != nil {
		log.Println(errorsx.Wrap(err, "unable to locate profile"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusNotFound))
		return
	}

	if err = meta.AuthzFindByProfileID(r.Context(), t.q, sqlx.NewNullString(pid)).Scan(&a); sqlx.IgnoreNoRows(err) != nil {
		log.Println(errorsx.Wrap(err, "unable to locate authz"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	registered := jwtx.NewJWTClaims(pid, jwtx.ClaimsOptionAuthzExpiration())
	claims := NewJWTClaim(TokenFromRegisterClaims(registered, TokenOptionFromAuthz(a)))
	if bearer, err = jwtx.Signed(t.jwtsecret(), claims); err != nil {
		log.Println(errorsx.Wrap(err, "unable to sign token"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if mp, err = NewProfileFromMetaProfile(p); err != nil {
		log.Println(errorsx.Wrap(err, "unable to encode profile"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err = httpx.WriteJSON(w, httpx.GetBuffer(r), &Authn{
		Token:   bearer,
		Profile: mp,
	}); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
	}
}

// register is used to create a new profile within a particular account.
func (t *HTTPSSO) register(w http.ResponseWriter, r *http.Request) {
	type tokenJSON struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int32  `json:"expires_in"`
	}

	var (
		err      error
		claims   jwt.RegisteredClaims
		consumed meta.ConsumedToken
		reg      Identity
		iden     identityssh.Identity
		p        meta.Profile
		bearer   string
	)

	if bearer, err = authn.AuthorizationToken(r.Context(), authn.JWTRegistrationSecretFromEnv, r, &claims); err != nil {
		httpx.Forbidden(w, errorsx.Wrap(err, "unable to decode authentication token"))
		return
	}

	if err = t.decoder.Decode(&reg, r.Form); err != nil {
		httpx.Forbidden(w, errorsx.Wrap(err, "unable to decode registration"))
		return
	}

	if err = meta.ConsumedTokenInsertWithDefaults(r.Context(), t.q, meta.ConsumedTokenFromJWTClaims(bearer, claims)).Scan(&consumed); err != nil {
		httpx.Forbidden(w, errorsx.Wrap(err, "unable to consume the token"))
		return
	}

	if err = identityssh.IdentityFindByID(r.Context(), t.q, sqlx.NewNullString(claims.Issuer)).Scan(&iden); err != nil {
		httpx.Forbidden(w, errorsx.Wrapf(err, "unable to locate user information %s", spew.Sdump(claims)))
		return
	}

	code := http.StatusConflict
	if uuid.FromStringOrNil(iden.ProfileID).String() == uuid.Nil.String() {
		if err = meta.ProfileInsertWithDefaults(r.Context(), t.q, langx.Clone(meta.Profile{}, meta.ProfileOptionDisplay(reg.Display))).Scan(&p); err != nil {
			httpx.Unauthorized(w, errorsx.Wrap(err, "profile creation failed"))
			return
		}

		iden.ProfileID = p.ID

		if err = identityssh.IdentityInsertWithDefaults(r.Context(), t.q, iden).Scan(&iden); duckdbx.ErrUniqueConstraintViolation(err) != nil {
			// ignore unique violations.
		} else if err != nil {
			httpx.Unauthorized(w, errorsx.Wrap(err, "unable to associate profile"))
			return
		}

		code = http.StatusOK
	}

	aclaims := jwtx.NewJWTClaims(
		p.ID,
		jwtx.ClaimsOptionAuthzExpiration(),
	)
	ast, err := jwtx.Signed(authn.JWTRegistrationSecretFromEnv(), aclaims)

	if err != nil {
		log.Println("unable to generate access token", err)
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	errorsx.Log(httpx.WriteJSONCode(w, code, httpx.GetBuffer(r), tokenJSON{
		AccessToken: ast,
		TokenType:   "BEARER",
		ExpiresIn:   int32(time.Until(aclaims.ExpiresAt.Time) / time.Second),
	}))
}
