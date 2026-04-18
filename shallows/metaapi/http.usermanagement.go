package metaapi

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/go-playground/form/v4"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/retrovibed/retrovibed/retroapi/httpauth"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/shallows/internal/duckdbx"
	"github.com/retrovibed/retrovibed/shallows/internal/env"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/formx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/lucenex"
	"github.com/retrovibed/retrovibed/shallows/internal/numericx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/sshx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/meta"
	"github.com/retrovibed/retrovibed/shallows/meta/identityssh"
	"golang.org/x/crypto/ssh"
)

type HTTPUsermanagementOption func(*HTTPUsermanagement)

func HTTPUsermanagementOptionJWTSecret(j jwtx.SecretSource) HTTPUsermanagementOption {
	return func(t *HTTPUsermanagement) {
		t.jwtsecret = j
	}
}

func NewHTTPUsermanagement(q sqlx.Queryer, options ...HTTPUsermanagementOption) *HTTPUsermanagement {
	svc := langx.Clone(HTTPUsermanagement{
		q:         q,
		jwtsecret: env.JWTSecret,
		decoder:   formx.NewDecoder(),
	}, options...)

	return &svc
}

type HTTPUsermanagement struct {
	q         sqlx.Queryer
	jwtsecret jwtx.SecretSource
	decoder   *form.Decoder
}

func (t *HTTPUsermanagement) Bind(r *mux.Router) {
	r.StrictSlash(false)

	r.Path("/").Methods(http.MethodGet).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpx.ParseForm,
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.Timeout2s(),
	).ThenFunc(t.search))

	r.Path("/").Methods(http.MethodPost).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpx.ParseForm,
		AuthzTokenHTTP(t.jwtsecret, AuthzPermUsermanagement),
		httpx.Timeout2s(),
	).ThenFunc(t.create))

	r.Path("/{id}").Methods(http.MethodPatch).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpx.ParseForm,
		AuthzTokenHTTP(t.jwtsecret, AuthzPermUsermanagement),
		httpx.Timeout2s(),
	).ThenFunc(t.update))

	r.Path("/{id}").Methods(http.MethodDelete).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpx.ParseForm,
		AuthzTokenHTTP(t.jwtsecret, AuthzPermUsermanagement),
		httpx.Timeout2s(),
	).ThenFunc(t.disable))

	r.Path("/{id}").Methods(http.MethodGet).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpx.ParseForm,
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.Timeout2s(),
	).ThenFunc(t.find))
}

func (t *HTTPUsermanagement) search(w http.ResponseWriter, r *http.Request) {
	var (
		err  error
		resp = ProfileSearchResponse{
			Next: &ProfileSearchRequest{
				Offset: 0,
				Limit:  100,
			},
		}
	)

	if err = t.decoder.Decode(resp.Next, r.Form); err != nil {
		log.Println(errorsx.Wrap(err, "unable to decode request"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}
	resp.Next.Limit = numericx.Min(resp.Next.Limit, 100)

	query := meta.ProfileSearchBuilder().
		Where(
			squirrel.And{
				meta.QueryIsEnabled(resp.Next.Status),
				lucenex.Query(duckdbx.NewLucene(), resp.Next.Query, lucenex.WithDefaultField("display")),
			},
		).
		Offset(resp.Next.Offset * resp.Next.Limit).Limit(resp.Next.Limit)

	q := sqlx.Scan(meta.ProfileSearch(r.Context(), t.q, query))
	for p := range q.Iter() {
		var (
			encoded *Profile
		)

		if encoded, err = NewProfileFromMetaProfile(p); err != nil {
			log.Println(errorsx.Wrap(err, "profile generation failed"))
			errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
			return
		}

		resp.Items = append(resp.Items, encoded)
	}

	if err = q.Err(); err != nil {
		log.Println(errorsx.Wrap(err, "profile generation failed"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err = httpx.WriteJSON(w, httpx.GetBuffer(r), &resp); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

func (t *HTTPUsermanagement) create(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		msg ProfileCreateRequest
		p   meta.Profile
	)

	if err = json.NewDecoder(r.Body).Decode(&msg); err != nil {
		log.Println("unable to decode request", err)
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	pubkey, comment, _, _, err := ssh.ParseAuthorizedKey([]byte(msg.PublicKey))
	if err != nil {
		log.Println("invalid public key", err)
		httpx.ErrorHeader(w, http.StatusBadRequest, errorsx.Wrap(err, "invalid ssh public key"))
		return
	}

	if p, err = NewMetaProfileFromProfile(msg.Profile, meta.ProfileOptionAutoDisplay(comment)); err != nil {
		log.Println("unable to convert profile", err)
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	p.ID = sshx.FingerprintMD5(pubkey)

	if err = meta.ProfileInsertWithID(r.Context(), t.q, p).Scan(&p); err != nil {
		log.Println("unable to create profile", err)
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err = meta.ProfileEnable(r.Context(), t.q, p.ID).Scan(&p); err != nil {
		log.Println("unable to enable profile", err)
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	iden := identityssh.Identity{
		ID:        sshx.FingerprintMD5(pubkey),
		PublicKey: sshx.EncodeBase64PublicKey(pubkey),
		ProfileID: p.ID,
	}

	if err = identityssh.IdentityInsertWithDefaults(r.Context(), t.q, iden).Scan(&iden); err != nil {
		log.Println("unable to create identity", err)
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	profile, err := NewProfileFromMetaProfile(p)
	if err != nil {
		log.Println("unable to convert profile", err)
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err = httpx.WriteJSON(w, httpx.GetBuffer(r), &ProfileCreateResponse{Profile: profile}); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

func (t *HTTPUsermanagement) find(w http.ResponseWriter, r *http.Request) {
	var (
		err  error
		resp = ProfileLookupResponse{
			Profile: &Profile{},
		}
		vars = mux.Vars(r)
		p    meta.Profile
	)

	if err = meta.ProfileFindByID(r.Context(), t.q, vars["id"]).Scan(&p); err != nil {
		log.Println(errorsx.Wrap(err, "unable to find profile"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusNotFound))
		return
	}

	if resp.Profile, err = NewProfileFromMetaProfile(p); err != nil {
		log.Println(errorsx.Wrap(err, "failed to encode profile"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err = httpx.WriteJSON(w, httpx.GetBuffer(r), &resp); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

func (t *HTTPUsermanagement) update(w http.ResponseWriter, r *http.Request) {
	var (
		err  error
		msg  ProfileUpdateRequest
		p    meta.Profile
		vars = mux.Vars(r)
	)

	if err = json.NewDecoder(r.Body).Decode(&msg); err != nil {
		log.Println(errorsx.Wrap(err, "unable to decode request"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	if p, err = NewMetaProfileFromProfile(msg.Profile); err != nil {
		log.Println(errorsx.Wrap(err, "converting profile failed"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	p.ID = vars["id"]

	if err = meta.ProfileUpdate(r.Context(), t.q, p).Scan(&p); err != nil {
		log.Println(errorsx.Wrap(err, "unable to insert profile"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	pp, err := NewProfileFromMetaProfile(langx.Clone(p, timex.JSONSafeEncodeOption))
	if err != nil {
		log.Println(errorsx.Wrap(err, "unable to encode profile"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	if err = httpx.WriteJSON(w, httpx.GetBuffer(r), &ProfileUpdateResponse{
		Profile: pp,
	}); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

func (t *HTTPUsermanagement) disable(w http.ResponseWriter, r *http.Request) {
	var (
		err  error
		p    meta.Profile
		vars = mux.Vars(r)
	)

	if err = meta.ProfileDisableByID(r.Context(), t.q, vars["id"]).Scan(&p); err != nil {
		log.Println(errorsx.Wrap(err, "unable to disable profile"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	pp, err := NewProfileFromMetaProfile(langx.Clone(p, timex.JSONSafeEncodeOption))
	if err != nil {
		log.Println(errorsx.Wrap(err, "unable to encode profile"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	if err = httpx.WriteJSON(w, httpx.GetBuffer(r), &ProfileDisableRequest{
		Profile: pp,
	}); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}
