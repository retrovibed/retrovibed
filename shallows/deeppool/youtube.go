package deeppool

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/retrovibed/retrovibed/internal/httpx"
	"golang.org/x/oauth2"
)

type GoogleTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	Scope        string `json:"scope"`
}

func NewYouTube(c *http.Client) YouTube {
	return YouTube{
		c:        c,
		endpoint: Deeppool(),
	}
}

type YouTube struct {
	c        *http.Client
	endpoint string
}

// Exchange trades an authorization code for tokens via deeppool.
func (t YouTube) Exchange(ctx context.Context, code, redirectURI string) (*GoogleTokenResponse, error) {
	var (
		err  error
		req  *http.Request
		resp *http.Response
		msg  GoogleTokenResponse
	)

	body := url.Values{
		"grant_type":   {"authorization_code"},
		"code":         {code},
		"redirect_uri": {redirectURI},
	}

	uri := fmt.Sprintf("https://%s/oauth2/proxy/google/token", t.endpoint)
	if req, err = http.NewRequestWithContext(ctx, http.MethodPost, uri, strings.NewReader(body.Encode())); err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if resp, err = httpx.AsError(t.c.Do(req)); err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&msg); err != nil {
		return nil, err
	}

	return &msg, nil
}

// TokenSource returns an oauth2.TokenSource that auto-refreshes via the deeppool proxy.
func (t YouTube) TokenSource(token *oauth2.Token) oauth2.TokenSource {
	return oauth2.ReuseTokenSource(token, &youtubeTokenSource{yt: t, refreshToken: token.RefreshToken})
}

type youtubeTokenSource struct {
	yt           YouTube
	refreshToken string
}

func (s *youtubeTokenSource) Token() (*oauth2.Token, error) {
	refreshed, err := s.yt.Refresh(context.Background(), s.refreshToken)
	if err != nil {
		return nil, err
	}

	tok := &oauth2.Token{
		AccessToken: refreshed.AccessToken,
		TokenType:   refreshed.TokenType,
		Expiry:      time.Now().Add(time.Duration(refreshed.ExpiresIn) * time.Second),
	}

	if refreshed.RefreshToken != "" {
		tok.RefreshToken = refreshed.RefreshToken
		s.refreshToken = refreshed.RefreshToken
	}

	return tok, nil
}

// Refresh obtains new tokens using a refresh token via deeppool.
func (t YouTube) Refresh(ctx context.Context, refreshToken string) (*GoogleTokenResponse, error) {
	var (
		err  error
		req  *http.Request
		resp *http.Response
		msg  GoogleTokenResponse
	)

	body := url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {refreshToken},
	}

	uri := fmt.Sprintf("https://%s/oauth2/proxy/google/token", t.endpoint)
	if req, err = http.NewRequestWithContext(ctx, http.MethodPost, uri, strings.NewReader(body.Encode())); err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if resp, err = httpx.AsError(t.c.Do(req)); err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&msg); err != nil {
		return nil, err
	}

	return &msg, nil
}
