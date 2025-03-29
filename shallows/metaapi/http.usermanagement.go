package metaapi

import (
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
	"github.com/retrovibed/retrovibed/internal/numericx"
	"github.com/retrovibed/retrovibed/internal/sqlx"
	"github.com/retrovibed/retrovibed/internal/sqlxx"
	"github.com/retrovibed/retrovibed/meta"
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
		jwtsecret: authn.JWTSecretFromEnv,
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
	r.Use(httpx.RouteInvoked)

	r.Path("/").Methods(http.MethodGet).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpx.ParseForm,
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.Timeout2s(),
	).ThenFunc(t.search))

	// r.Path("/").Methods(http.MethodPost).Handler(alice.New(
	// 	httpx.ContextBufferPool512(),
	// 	httpx.ParseForm,
	// 	httpauth.AuthenticateWithToken(t.jwtsecret),
	// 	// AuthzTokenHTTP(t.jwtsecret, AuthzPermUsermanagement),
	// 	httpx.Timeout2s(),
	// ).ThenFunc(t.create))

	// r.Path("/{id}").Methods(http.MethodPatch).Handler(alice.New(
	// 	httpx.ContextBufferPool512(),
	// 	httpx.ParseForm,
	// 	// AuthzTokenHTTP(t.jwtsecret, AuthzPermUsermanagement),
	// 	httpx.Timeout2s(),
	// ).ThenFunc(t.update))

	// r.Path("/{id}").Methods(http.MethodDelete).Handler(alice.New(
	// 	httpx.ContextBufferPool512(),
	// 	httpx.ParseForm,
	// 	// AuthzTokenHTTP(t.jwtsecret, AuthzPermUsermanagement),
	// 	httpx.Timeout2s(),
	// ).ThenFunc(t.disable))

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
		// Where(
		// 	squirrel.And{
		// 		profiles.QuerySearchVector(resp.Next.Query),
		// 		profiles.QueryIsEnabled(resp.Next.Status),
		// 	},
		// ).
		Offset(resp.Next.Offset * resp.Next.Limit).Limit(resp.Next.Limit)

	err = sqlxx.ScanEach(meta.ProfileSearch(r.Context(), t.q, query), func(p *meta.Profile) (failed error) {
		var (
			encoded *Profile
		)

		if encoded, failed = NewProfileFromMetaProfile(*p); failed != nil {
			return failed
		}

		resp.Items = append(resp.Items, encoded)
		return nil
	})

	if err != nil {
		log.Println(errorsx.Wrap(err, "profile generation failed"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	resp.Next.Offset += 1

	if err = httpx.WriteJSON(w, httpx.GetBuffer(r), &resp); err != nil {
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

// func (t *HTTPUsermanagement) update(w http.ResponseWriter, r *http.Request) {
// 	var (
// 		err  error
// 		aid  string
// 		msg  ProfileUpdateRequest
// 		p    meta.Profile
// 		vars = mux.Vars(r)
// 	)

// 	if err = json.NewDecoder(r.Body).Decode(&msg); err != nil {
// 		log.Println(errorsx.Wrap(err, "unable to decode request"))
// 		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
// 		return
// 	}
// 	p.ID = vars["id"]

// 	if err = msg.Profile.ToProfile(&p, profiles.OptionAccountID(aid), profiles.OptionJSONSafeDecode); err != nil {
// 		log.Println(errorsx.Wrap(err, "converting profile failed"))
// 		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
// 		return
// 	}

// 	if err = profiles.ProfileUpdate(r.Context(), t.q, &p).Scan(&p); err != nil {
// 		log.Println(errorsx.Wrap(err, "unable to insert profile"))
// 		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
// 		return
// 	}

// 	if err = httpx.WriteJSON(w, httpx.GetBuffer(r), &LookupResponse{
// 		Profile: (&Profile{}).FromProfile(langx.Autoptr(langx.Clone(p, profiles.OptionJSONSafeEncode, profiles.OptionTimezoneUTC))),
// 	}); err != nil {
// 		log.Println(errorsx.Wrap(err, "unable to write response"))
// 		return
// 	}
// }

// func (t *HTTPUsermanagement) disable(w http.ResponseWriter, r *http.Request) {
// 	var (
// 		err error
// 		aid string
// 		msg DisableRequest
// 		p   profiles.Profile
// 	)

// 	if aid, err = httpauth.AccountID(r.Context(), t.jwtsecret, r); err != nil {
// 		log.Println(errorsx.Wrap(err, "unable to decode json request"))
// 		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
// 		return
// 	}

// 	if err = json.NewDecoder(r.Body).Decode(&msg); err != nil {
// 		log.Println(errorsx.Wrap(err, "unable to decode request"))
// 		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
// 		return
// 	}

// 	if err = profiles.ProfileDisableByID(r.Context(), t.q, aid, msg.Profile.Id).Scan(&p); err != nil {
// 		log.Println(errorsx.Wrap(err, "unable to disable profile"))
// 		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
// 		return
// 	}

// 	if err = httpx.WriteJSON(w, httpx.GetBuffer(r), &DisableResponse{
// 		Profile: (&Profile{}).FromProfile(langx.Autoptr(langx.Clone(p, profiles.OptionJSONSafeEncode))),
// 	}); err != nil {
// 		log.Println(errorsx.Wrap(err, "unable to write response"))
// 		return
// 	}
// }
