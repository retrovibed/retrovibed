package community

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-playground/form/v4"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/retrovibed/retrovibed/retroapi/deeppool"
	"github.com/retrovibed/retrovibed/retroapi/httpauth"
	"github.com/retrovibed/retrovibed/retroapi/jwtx"
	"github.com/retrovibed/retrovibed/shallows/internal/env"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/formx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
)

type HTTPYouTubeOption func(*HTTPYouTube)

func HTTPYouTubeOptionJWTSecret(j jwtx.SecretSource) HTTPYouTubeOption {
	return func(t *HTTPYouTube) {
		t.jwtsecret = j
	}
}

func NewHTTPYouTube(q sqlx.Queryer, httpc *http.Client, options ...HTTPYouTubeOption) *HTTPYouTube {
	svc := langx.Clone(HTTPYouTube{
		q:         q,
		jwtsecret: env.JWTSecret,
		decoder:   formx.NewDecoder(),
		httpc:     httpc,
	}, options...)

	return &svc
}

type HTTPYouTube struct {
	q         sqlx.Queryer
	jwtsecret jwtx.SecretSource
	decoder   *form.Decoder
	httpc     *http.Client
}

func (t *HTTPYouTube) Bind(r *mux.Router) {
	r.StrictSlash(false)

	r.Path("/callback").Methods(http.MethodGet).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpx.ParseForm,
		httpx.Timeout4s(),
	).ThenFunc(t.callback))

	r.Path("/status").Methods(http.MethodGet).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.Timeout2s(),
	).ThenFunc(t.status))

	r.Path("/token").Methods(http.MethodDelete).Handler(alice.New(
		httpx.ContextBufferPool512(),
		httpauth.AuthenticateWithToken(t.jwtsecret),
		httpx.Timeout2s(),
	).ThenFunc(t.revoke))
}

func callbackHTML(w http.ResponseWriter, status int, heading string, body string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	_, _ = fmt.Fprintf(w, `<!DOCTYPE html><html><body><h2>%s</h2><p>%s</p></body></html>`, heading, body)
}

func (t *HTTPYouTube) callback(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Code string `json:"code"`
	}

	if err := t.decoder.Decode(&req, r.Form); err != nil {
		log.Println(errorsx.Wrap(err, "unable to decode callback request"))
		callbackHTML(w, http.StatusBadRequest, "YouTube linking failed.", "Invalid request.")
		return
	}

	if req.Code == "" {
		log.Println("youtube callback: missing code parameter")
		callbackHTML(w, http.StatusBadRequest, "YouTube linking failed.", "Missing authorization code.")
		return
	}

	redirectURI := fmt.Sprintf("https://%s%s", r.Host, r.URL.Path)
	token, err := deeppool.NewYouTube(t.httpc).Exchange(r.Context(), req.Code, redirectURI)
	if err != nil {
		log.Println(errorsx.Wrap(err, "youtube token exchange failed"))
		callbackHTML(w, http.StatusBadGateway, "YouTube linking failed.", "Unable to complete authorization with YouTube.")
		return
	}

	row := OAuth2Google{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		TokenType:    token.TokenType,
		Expiry:       time.Now().Add(time.Duration(token.ExpiresIn) * time.Second),
		Scopes:       "https://www.googleapis.com/auth/youtube.upload",
	}

	if err = OAuth2GoogleInsertWithDefaults(r.Context(), t.q, row).Scan(&row); err != nil {
		log.Println(errorsx.Wrap(err, "unable to store youtube token"))
		callbackHTML(w, http.StatusInternalServerError, "YouTube linking failed.", "Unable to save authorization.")
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprint(w, `<!DOCTYPE html><html><body><script>window.close()</script><p>YouTube linked. You can close this window.</p></body></html>`)
}

func (t *HTTPYouTube) status(w http.ResponseWriter, r *http.Request) {
	var row OAuth2Google
	err := OAuth2GoogleFindFirst(r.Context(), t.q).Scan(&row)

	linked := err == nil && row.AccessToken != ""

	if err = httpx.WriteJSON(w, httpx.GetBuffer(r), struct {
		Linked bool   `json:"linked"`
		ID     string `json:"id,omitempty"`
	}{
		Linked: linked,
		ID:     row.ID,
	}); err != nil {
		log.Println(errorsx.Wrap(err, "unable to write youtube status"))
	}
}

func (t *HTTPYouTube) revoke(w http.ResponseWriter, r *http.Request) {
	var row OAuth2Google
	if err := OAuth2GoogleDeleteAll(r.Context(), t.q).Scan(&row); sqlx.IgnoreNoRows(err) != nil {
		log.Println(errorsx.Wrap(err, "unable to delete youtube token"))
		errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusInternalServerError))
		return
	}

	errorsx.Log(httpx.WriteEmptyJSON(w, http.StatusOK))
}
