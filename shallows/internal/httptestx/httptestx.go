package httptestx

import (
	"bufio"
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
)

func BuildURL(path string, v url.Values) *url.URL {
	return (&url.URL{
		Scheme:   "http",
		Host:     "example.com",
		Path:     path,
		RawQuery: v.Encode(),
	})
}

// ReadRequest reads a request from a file.
func ReadRequest(path string) (resp *httptest.ResponseRecorder, req *http.Request, err error) {
	var (
		raw []byte
	)

	if raw, err = os.ReadFile(path); err != nil {
		return nil, nil, err
	}

	req, err = http.ReadRequest(bufio.NewReader(bytes.NewReader(raw)))
	return httptest.NewRecorder(), req, err
}

type RequestOption func(*http.Request)

// BuildRequest ...
func BuildRequest(method string, uri string, body []byte, options ...RequestOption) (recorder *httptest.ResponseRecorder, req *http.Request, err error) {

	recorder = httptest.NewRecorder()
	if req, err = http.NewRequest(strings.ToUpper(method), uri, bytes.NewBuffer(body)); err != nil {
		return recorder, req, err
	}

	for _, opt := range options {
		opt(req)
	}

	return recorder, req, nil
}

func RequestOptionURL(uri *url.URL) RequestOption {
	return func(r *http.Request) {
		r.URL = uri
	}
}

func RequestOptionAuthorization(value string) RequestOption {
	return RequestOptionHeader("authorization", value)
}

func RequestOptionContent(value string) RequestOption {
	return RequestOptionHeader("Content-Type", value)
}

func RequestOptionHeader(key string, value string) RequestOption {
	return func(r *http.Request) {
		r.Header.Add(key, value)
	}
}

// BuildGetRequest ...
func BuildGetRequest(body []byte, options ...RequestOption) (*httptest.ResponseRecorder, *http.Request, error) {
	return BuildRequest(http.MethodGet, "http://example.com/", body, options...)
}

// BuildPostRequest ...
func BuildPostRequest(body []byte, options ...RequestOption) (*httptest.ResponseRecorder, *http.Request, error) {
	return BuildRequest(http.MethodPost, "http://example.com/", body, options...)
}

// BuildDeleteRequest ...
func BuildDeleteRequest(body []byte, options ...RequestOption) (*httptest.ResponseRecorder, *http.Request, error) {
	return BuildRequest(http.MethodDelete, "http://example.com/", body, options...)
}

// RoundTripFunc ...
type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip ...
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

// NewTestClient returns *http.Client with Transport replaced to avoid making real calls
func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}
