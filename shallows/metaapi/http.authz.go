package metaapi

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-playground/form/v4"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/retrovibed/retrovibed/authn"
	"github.com/retrovibed/retrovibed/httpauth"
	"github.com/retrovibed/retrovibed/internal/errorsx"
	"github.com/retrovibed/retrovibed/internal/formx"
	"github.com/retrovibed/retrovibed/internal/httpx"
	"github.com/retrovibed/retrovibed/internal/jwtx"
	"github.com/retrovibed/retrovibed/internal/langx"
	"github.com/retrovibed/retrovibed/internal/sqlx"
	"github.com/retrovibed/retrovibed/meta"
)

type HTTPAuthzOption func(*HTTPAuthz)

func HTTPAuthzOptionJWTSecret(j jwtx.SecretSource) HTTPAuthzOption {
	return func(t *HTTPAuthz) {
		t.jwtsecret = j
	}
}

func NewHTTPAuthz(q sqlx.Queryer, options ...HTTPAuthzOption) *HTTPAuthz {
	svc := langx.Clone(HTTPAuthz{
		q:         q,
		jwtsecret: authn.JWTSecretFromEnv,
		decoder:   formx.NewDecoder(),
	}, options...)

	return &svc
}

type HTTPAuthz struct {
	q         sqlx.Queryer
	jwtsecret jwtx.SecretSource
	decoder   *form.Decoder
}

func (t *HTTPAuthz) Bind(r *mux.Router) {
	r.StrictSlash(false)

	r.Path("/").Methods(http.MethodGet).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpx.ParseForm,
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.Timeout2s(),
	).ThenFunc(t.authz))

	r.Path("/{id}").Methods(http.MethodGet).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpx.ParseForm,
		AuthzTokenHTTP(t.jwtsecret, AuthzPermUsermanagement),
		httpx.Timeout2s(),
	).ThenFunc(t.profile))

	r.Path("/{id}").Methods(http.MethodPost).Handler(alice.New(
		httpx.ContextBufferPool512(),
		AuthzTokenHTTP(t.jwtsecret, AuthzPermUsermanagement),
		httpx.ParseForm,
		httpx.Timeout2s(),
	).ThenFunc(t.grant))

	r.Path("/{id}").Methods(http.MethodDelete).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpx.ParseForm,
		AuthzTokenHTTP(t.jwtsecret, AuthzPermUsermanagement),
		httpx.Timeout2s(),
	).ThenFunc(t.revoke))
}

func (t *HTTPAuthz) authz(w http.ResponseWriter, r *http.Request) {
	var (
		err    error
		pid    string
		bearer string
		a      meta.Authz
	)

	if _, pid, err = httpauth.IssuerSubjectID(r.Context(), t.jwtsecret, r); err != nil {
		httpx.ErrorHeader(w, http.StatusBadRequest, errorsx.Wrap(err, "failed to retrieve token"))
		return
	}

	if err = meta.AuthzFindByProfileID(r.Context(), t.q, sqlx.NewNullString(pid)).Scan(&a); sqlx.IgnoreNoRows(err) != nil {
		httpx.ErrorHeader(w, http.StatusInternalServerError, errorsx.Wrap(err, "failed to retrieve authz"))
		return
	}

	registered := jwtx.NewJWTClaims(
		pid,
		jwtx.ClaimsOptionAuthzExpiration(),
	)

	claims := NewJWTClaim(TokenFromRegisterClaims(registered, TokenOptionFromAuthz(a)))

	if bearer, err = jwtx.Signed(t.jwtsecret(), claims); err != nil {
		httpx.ErrorHeader(w, http.StatusInternalServerError, errorsx.Wrap(err, "failed to generate signed token"))
		return
	}

	auth := &AuthzResponse{
		Bearer: bearer,
		Token:  claims.Token,
	}

	if err = httpx.WriteJSON(w, httpx.GetBuffer(r), &auth); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

func (t *HTTPAuthz) profile(w http.ResponseWriter, r *http.Request) {
	var (
		err    error
		a      meta.Authz
		pid, _ = mux.Vars(r)["id"]
	)

	if err = meta.AuthzFindByProfileID(r.Context(), t.q, sqlx.NewNullString(pid)).Scan(&a); sqlx.IgnoreNoRows(err) != nil {
		httpx.ErrorHeader(w, http.StatusInternalServerError, errorsx.Wrap(err, "failed to retrieve authz"))
		return
	}

	registered := jwtx.NewJWTClaims(
		pid,
		jwtx.ClaimsOptionAuthzExpiration(),
	)

	if err = httpx.WriteJSON(w, httpx.GetBuffer(r), &AuthzProfileResponse{
		Token: TokenFromRegisterClaims(registered, TokenOptionFromAuthz(a)),
	}); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

func (t *HTTPAuthz) grant(w http.ResponseWriter, r *http.Request) {
	var (
		err    error
		msg    AuthzGrantRequest
		a      meta.Authz
		pid, _ = mux.Vars(r)["id"]
	)

	if err = json.NewDecoder(r.Body).Decode(&msg); err != nil {
		log.Println("unable to decode request", err)
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	a = meta.Authz{
		ProfileID:      pid,
		Usermanagement: msg.Token.Usermanagement,
	}

	if err = meta.AuthzUpsertWithDefaults(r.Context(), t.q, a).Scan(&a); err != nil {
		log.Println("upsert failed", err)
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	registered := jwtx.NewJWTClaims(
		pid,
		jwtx.ClaimsOptionAuthzExpiration(),
	)

	g := &AuthzGrantResponse{
		Token: TokenFromRegisterClaims(registered, TokenOptionFromAuthz(a)),
	}
	if err = httpx.WriteJSON(w, httpx.GetBuffer(r), g); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

func (t *HTTPAuthz) revoke(w http.ResponseWriter, r *http.Request) {
	var (
		err    error
		a      meta.Authz
		pid, _ = mux.Vars(r)["id"]
	)

	if err = meta.AuthzDeleteByProfileID(r.Context(), t.q, pid).Scan(&a); err != nil {
		log.Println("unable to delete record", err)
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	registered := jwtx.NewJWTClaims(
		pid,
		jwtx.ClaimsOptionAuthzExpiration(),
	)

	g := &AuthzRevokeResponse{
		Token: TokenFromRegisterClaims(registered),
	}

	if err = httpx.WriteJSON(w, httpx.GetBuffer(r), g); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}
