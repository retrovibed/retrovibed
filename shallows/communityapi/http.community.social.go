package communityapi

import (
	"log"
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/go-playground/form/v4"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/retrovibed/retrovibed/retroapi/httpauth"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/shallows/community"
	"github.com/retrovibed/retrovibed/shallows/internal/duckdbx"
	"github.com/retrovibed/retrovibed/shallows/internal/env"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/formx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/lucenex"
	"github.com/retrovibed/retrovibed/shallows/internal/md5x"
	"github.com/retrovibed/retrovibed/shallows/internal/numericx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
)

type HTTPSocialOption func(*HTTPSocial)

func HTTPSocialOptionJWTSecret(j jwtx.SecretSource) HTTPSocialOption {
	return func(t *HTTPSocial) {
		t.jwtsecret = j
	}
}

func NewHTTPSocial(q sqlx.Queryer, options ...HTTPSocialOption) *HTTPSocial {
	svc := langx.Clone(HTTPSocial{
		q:         q,
		jwtsecret: env.JWTSecret,
		decoder:   formx.NewDecoder(),
		lucene:    duckdbx.NewLucene(),
	}, options...)

	return &svc
}

type HTTPSocial struct {
	q         sqlx.Queryer
	jwtsecret jwtx.SecretSource
	decoder   *form.Decoder
	lucene    lucenex.Driver
}

func (t *HTTPSocial) Bind(r *mux.Router) {
	r.StrictSlash(false)

	r.Path("/").Methods(http.MethodGet).Handler(alice.New(
		httpx.RouteInvoked,
		httpx.ContextBufferPool512(),
		httpx.ParseForm,
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.Timeout2s(),
	).ThenFunc(t.search))

	r.Path("/{communityId}/publishers/{publisherId}").Methods(http.MethodPost).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.Timeout2s(),
	).ThenFunc(t.enable))

	r.Path("/{communityId}/publishers/{publisherId}").Methods(http.MethodDelete).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.Timeout2s(),
	).ThenFunc(t.disable))
}

// search returns the authenticated account's communities, each with its
// currently-enabled publishers, plus the full publisher catalog so the
// console can render the available toggles.
func (t *HTTPSocial) search(w http.ResponseWriter, r *http.Request) {
	var resp SocialsSearchResponse
	resp.Next = &SocialsSearchRequest{Limit: 100}

	if err := t.decoder.Decode(resp.Next, r.Form); err != nil {
		log.Println(errorsx.Wrap(err, "unable to decode search request"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}
	resp.Next.Limit = numericx.Min(resp.Next.Limit, 100)

	_, pid, err := httpauth.IssuerSubjectID(r.Context(), t.jwtsecret, r)
	if err != nil {
		log.Println(errorsx.Wrap(err, "unable to retrieve token"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusBadRequest))
		return
	}

	catalog := community.PluginPublisherFindAll(r.Context(), t.q)
	ci := sqlx.Scan(catalog)
	for p := range ci.Iter() {
		resp.Catalog = append(resp.Catalog, NewPluginPublisher(PluginPublisherOptionFromDB(langx.Clone(p, timex.JSONSafeEncodeOption))))
	}
	if err := ci.Err(); err != nil {
		log.Println(errorsx.Wrap(err, "unable to list plugin publishers"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	communities := community.CommunitySearch(r.Context(), t.q, community.CommunitySearchBuilder().
		Distinct().
		Join("community_publisher ON community_publisher.community_id = community.id").
		Where(
			squirrel.And{
				squirrel.Eq{"community.account_id": pid},
				lucenex.Query(t.lucene, resp.Next.Query, lucenex.WithDefaultField("description")),
			},
		).Offset(resp.Next.Offset*resp.Next.Limit).Limit(resp.Next.Limit))
	qi := sqlx.Scan(communities)
	for c := range qi.Iter() {
		social := NewCommunitySocial(func(s *CommunitySocial) {
			s.Community = NewCommunity(CommunityOptionFromDB(langx.Clone(c, timex.JSONSafeEncodeOption)))
		})

		enabled := community.CommunityPublisherFindByCommunityID(r.Context(), t.q, c.ID)
		ei := sqlx.Scan(enabled)
		for cp := range ei.Iter() {
			social.Enabled = append(social.Enabled, NewCommunityPublisher(CommunityPublisherOptionFromDB(langx.Clone(cp, timex.JSONSafeEncodeOption))))
		}
		if err := ei.Err(); err != nil {
			log.Println(errorsx.Wrap(err, "unable to list enabled publishers"))
			errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
			return
		}

		resp.Items = append(resp.Items, social)
	}
	if err := qi.Err(); err != nil {
		log.Println(errorsx.Wrap(err, "unable to list communities"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err := httpx.WriteJSON(w, httpx.GetBuffer(r), &resp); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

func (t *HTTPSocial) enable(w http.ResponseWriter, r *http.Request) {
	var existing community.CommunityPublisher

	vars := mux.Vars(r)
	communityID, publisherID := vars["communityId"], vars["publisherId"]

	cp := community.CommunityPublisher{
		ID:          md5x.FormatUUID(md5x.Digest(communityID, publisherID)),
		CommunityID: communityID,
		PublisherID: publisherID,
	}

	if err := community.CommunityPublisherInsertWithDefaults(r.Context(), t.q, cp).Scan(&existing); err != nil {
		log.Println(errorsx.Wrap(err, "unable to enable publisher"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	if err := httpx.WriteJSON(w, httpx.GetBuffer(r), &CommunityPublisherEnableResponse{
		Enabled: NewCommunityPublisher(CommunityPublisherOptionFromDB(langx.Clone(existing, timex.JSONSafeEncodeOption))),
	}); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write response"))
		return
	}
}

func (t *HTTPSocial) disable(w http.ResponseWriter, r *http.Request) {
	// var existing community.CommunityPublisher

	// vars := mux.Vars(r)

	// err := community.CommunityPublisherDeleteByCommunityIDAndPublisherID(r.Context(), t.q, vars["communityId"], vars["publisherId"]).Scan(&existing)
	// if errors.Is(err, sql.ErrNoRows) {
	// 	errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusNotFound))
	// 	return
	// } else if err != nil {
	// 	log.Println(errorsx.Wrap(err, "unable to disable publisher"))
	// 	errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
	// 	return
	// }

	// if err := httpx.WriteJSON(w, httpx.GetBuffer(r), &CommunityPublisherDisableResponse{
	// 	Disabled: NewCommunityPublisher(CommunityPublisherOptionFromDB(langx.Clone(existing, timex.JSONSafeEncodeOption))),
	// }); err != nil {
	// 	log.Println(errorsx.Wrap(err, "unable to write response"))
	// 	return
	// }
}
