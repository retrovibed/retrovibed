package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/retrovibed/retrovibed/retroapi/errorsx"
)

const (
	// ErrUnauthorized is what every request helper reports for a 401,
	// letting a caller distinguish "this session expired, log in again"
	// from a genuine failure. Lemmy's jwts outlive a single publish but
	// not forever, and the cached session on disk is exactly the thing
	// that goes stale.
	ErrUnauthorized = errorsx.String("lemmy rejected the credentials")
)

// session is the cached login state written beside the plugin's other
// scratch state in CACHE_DIRECTORY. It holds nothing that isn't
// recoverable by logging in again, which is the point: losing it costs one
// extra round trip, not a publish.
type session struct {
	JWT string `json:"jwt"`
}

// Post is the subset of Lemmy's post representation this plugin reports
// back to the registry.
type Post struct {
	ID   int64  `json:"id"`
	ApID string `json:"ap_id"`
	Name string `json:"name"`
}

// CreatePost mirrors Lemmy's /api/v3/post request body. Only Name and
// CommunityID are required; everything else is omitted when empty so Lemmy
// applies its own defaults rather than receiving explicit zero values.
//
// URL accepts a magnet link: Lemmy's post url validation allows exactly
// the http, https and magnet schemes, so the content's magnet URI goes in
// here directly rather than being buried in the body.
type CreatePost struct {
	Name            string `json:"name"`
	CommunityID     int64  `json:"community_id"`
	URL             string `json:"url,omitempty"`
	Body            string `json:"body,omitempty"`
	AltText         string `json:"alt_text,omitempty"`
	NSFW            bool   `json:"nsfw,omitempty"`
	LanguageID      int64  `json:"language_id,omitempty"`
	CustomThumbnail string `json:"custom_thumbnail,omitempty"`
}

type loginRequest struct {
	UsernameOrEmail string `json:"username_or_email"`
	Password        string `json:"password"`
	TOTP            string `json:"totp_2fa_token,omitempty"`
}

type loginResponse struct {
	JWT string `json:"jwt"`
}

type communityResponse struct {
	CommunityView struct {
		Community struct {
			ID int64 `json:"id"`
		} `json:"community"`
	} `json:"community_view"`
}

type postResponse struct {
	PostView struct {
		Post Post `json:"post"`
	} `json:"post_view"`
}

// pictrsResponse is what Lemmy's image endpoint returns. It is not part of
// the /api/v3 surface - image hosting is pict-rs sitting behind the same
// origin - so it neither shares the api prefix nor the error envelope.
type pictrsResponse struct {
	Message string `json:"msg"`
	Files   []struct {
		File        string `json:"file"`
		DeleteToken string `json:"delete_token"`
	} `json:"files"`
}

// Client talks to one Lemmy instance's v3 HTTP API. Everything that is
// fixed for the lifetime of the plugin invocation - which instance, which
// http client, where to cache the session - is configured here; what
// varies per call is a method argument.
type Client struct {
	instance *url.URL
	http     *http.Client
	token    string
	cachedir string
}

// Option customizes a Client built by NewClient.
type Option func(*Client)

// OptionHTTPClient substitutes the http client used for every request
// (a test server, a client with different timeouts). Defaults to
// http.DefaultClient, which under wasip1 is the wasinet-hijacked
// transport.
func OptionHTTPClient(c *http.Client) Option {
	return func(t *Client) { t.http = c }
}

// OptionToken presets the bearer token, skipping login entirely - for an
// operator who would rather paste a long lived token into the plugin's
// configuration than store a password.
func OptionToken(token string) Option {
	return func(t *Client) { t.token = token }
}

// OptionCacheDir points the session cache at dir (default: no caching, so
// every invocation logs in). The registry mounts a per-plugin cache
// directory and names it in CACHE_DIRECTORY, which is what the publish
// command passes here.
func OptionCacheDir(dir string) Option {
	return func(t *Client) { t.cachedir = dir }
}

// NewClient builds a client for the lemmy instance at uri (e.g.
// https://lemmy.ml).
func NewClient(uri string, options ...Option) (*Client, error) {
	instance, err := url.Parse(uri)
	if err != nil {
		return nil, errorsx.Wrapf(err, "unable to parse lemmy instance url: %s", uri)
	}

	if instance.Scheme == "" || instance.Host == "" {
		return nil, errorsx.String("lemmy instance url must be absolute, e.g. https://lemmy.ml")
	}

	c := &Client{instance: instance, http: http.DefaultClient}
	for _, opt := range options {
		opt(c)
	}

	return c, nil
}

// Authenticated reports whether the client currently holds a token, from
// OptionToken, a restored session, or a Login.
func (t *Client) Authenticated() bool {
	return t.token != ""
}

// sessionPath is where the cached jwt lives, or "" when no cache directory
// was configured.
func (t *Client) sessionPath() string {
	if t.cachedir == "" {
		return ""
	}

	return filepath.Join(t.cachedir, "session.json")
}

// RestoreSession adopts a previously cached jwt, if one is there. A missing
// or unreadable cache is not an error - it just means logging in - so the
// only thing a caller does with the result is decide whether to call Login.
func (t *Client) RestoreSession() bool {
	path := t.sessionPath()
	if path == "" {
		return false
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		return false
	}

	var s session
	if err := json.Unmarshal(raw, &s); err != nil || s.JWT == "" {
		return false
	}

	t.token = s.JWT

	return true
}

// storeSession caches the current token for the next invocation. Failing to
// cache is never fatal to a publish that already succeeded, so the error is
// returned for logging rather than propagated.
func (t *Client) storeSession() error {
	path := t.sessionPath()
	if path == "" {
		return nil
	}

	encoded, err := json.Marshal(session{JWT: t.token})
	if err != nil {
		return errorsx.Wrap(err, "unable to encode lemmy session")
	}

	return errorsx.Wrapf(os.WriteFile(path, encoded, 0600), "unable to cache lemmy session: %s", path)
}

// Login exchanges credentials for a jwt and caches it. totp may be empty
// unless the account has 2fa enabled.
func (t *Client) Login(ctx context.Context, username, password, totp string) error {
	var resp loginResponse

	req := loginRequest{UsernameOrEmail: username, Password: password, TOTP: totp}
	if err := t.do(ctx, http.MethodPost, "/api/v3/user/login", req, &resp); err != nil {
		return errorsx.Wrap(err, "unable to log into lemmy")
	}

	if resp.JWT == "" {
		return errorsx.String("lemmy accepted the login but returned no token; the account likely requires email verification or admin approval")
	}

	t.token = resp.JWT

	return t.storeSession()
}

// ResolveCommunity maps a community name - either "movies" for one local to
// this instance, or "movies@lemmy.ml" for a remote one - onto the numeric
// id CreatePost needs.
func (t *Client) ResolveCommunity(ctx context.Context, name string) (int64, error) {
	var resp communityResponse

	if err := t.do(ctx, http.MethodGet, "/api/v3/community?name="+url.QueryEscape(name), nil, &resp); err != nil {
		return 0, errorsx.Wrapf(err, "unable to resolve lemmy community: %s", name)
	}

	if resp.CommunityView.Community.ID == 0 {
		return 0, errorsx.Wrapf(errorsx.String("lemmy returned no community"), "community: %s", name)
	}

	return resp.CommunityView.Community.ID, nil
}

// CreatePost submits a post and returns what lemmy recorded.
func (t *Client) CreatePost(ctx context.Context, post CreatePost) (*Post, error) {
	var resp postResponse

	if err := t.do(ctx, http.MethodPost, "/api/v3/post", post, &resp); err != nil {
		return nil, errorsx.Wrap(err, "unable to create lemmy post")
	}

	return &resp.PostView.Post, nil
}

// UploadImage posts path to the instance's pict-rs endpoint and returns the
// publicly reachable url of the stored image, suitable for
// CreatePost.CustomThumbnail. Lemmy has no notion of a post attachment -
// an image is hosted separately and referenced by url - so this is only
// worth calling for content small enough, and of a type, that a thumbnail
// makes sense.
func (t *Client) UploadImage(ctx context.Context, path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", errorsx.Wrapf(err, "unable to open image: %s", path)
	}
	defer f.Close()

	body := bytes.NewBuffer(nil)
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("images[]", filepath.Base(path))
	if err != nil {
		return "", errorsx.Wrap(err, "unable to build image upload")
	}

	if _, err = io.Copy(part, f); err != nil {
		return "", errorsx.Wrapf(err, "unable to buffer image: %s", path)
	}

	if err = writer.Close(); err != nil {
		return "", errorsx.Wrap(err, "unable to finalize image upload")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, t.instance.JoinPath("pictrs", "image").String(), body)
	if err != nil {
		return "", errorsx.Wrap(err, "unable to build image upload request")
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if t.token != "" {
		req.Header.Set("Authorization", "Bearer "+t.token)
	}

	resp, err := t.http.Do(req)
	if err != nil {
		return "", errorsx.Wrap(err, "unable to upload image")
	}
	defer resp.Body.Close()

	if err := statusError(resp); err != nil {
		return "", err
	}

	var decoded pictrsResponse
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		return "", errorsx.Wrap(err, "unable to decode image upload response")
	}

	if len(decoded.Files) == 0 {
		return "", errorsx.Wrapf(errorsx.String("image upload stored nothing"), "response: %s", decoded.Message)
	}

	return t.instance.JoinPath("pictrs", "image", decoded.Files[0].File).String(), nil
}

// do issues one json api request. body is encoded as the request payload
// when non-nil, and the response is decoded into dst when non-nil.
func (t *Client) do(ctx context.Context, method, path string, body any, dst any) error {
	var payload io.Reader

	if body != nil {
		encoded, err := json.Marshal(body)
		if err != nil {
			return errorsx.Wrap(err, "unable to encode request")
		}
		payload = bytes.NewReader(encoded)
	}

	req, err := http.NewRequestWithContext(ctx, method, t.instance.String()+path, payload)
	if err != nil {
		return errorsx.Wrap(err, "unable to build request")
	}

	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if t.token != "" {
		req.Header.Set("Authorization", "Bearer "+t.token)
	}

	resp, err := t.http.Do(req)
	if err != nil {
		return errorsx.Wrap(err, "request failed")
	}
	defer resp.Body.Close()

	if err := statusError(resp); err != nil {
		return err
	}

	if dst == nil {
		return nil
	}

	return errorsx.Wrap(json.NewDecoder(resp.Body).Decode(dst), "unable to decode response")
}

// statusError turns a non-2xx response into an error, preserving lemmy's
// own error text (its bodies are small and say things like
// "couldnt_find_community", which is the entire diagnostic).
func statusError(resp *http.Response) error {
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	cause, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))

	if resp.StatusCode == http.StatusUnauthorized {
		return errorsx.Wrapf(ErrUnauthorized, "%s", bytes.TrimSpace(cause))
	}

	return fmt.Errorf("lemmy returned %s: %s", resp.Status, bytes.TrimSpace(cause))
}
