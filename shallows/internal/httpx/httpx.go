package httpx

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/justinas/alice"
	"github.com/retrovibed/retrovibed/internal/backoffx"
	"github.com/retrovibed/retrovibed/internal/bytesx"
	"github.com/retrovibed/retrovibed/internal/debugx"
	"github.com/retrovibed/retrovibed/internal/errorsx"
	"github.com/retrovibed/retrovibed/internal/netx"
	"github.com/retrovibed/retrovibed/internal/stringsx"
	"golang.org/x/time/rate"
)

// CheckStatusCode compares the provided status code with a list of acceptable
// status codes.
func CheckStatusCode(actual int, acceptable ...int) bool {
	for _, code := range acceptable {
		if actual == code {
			return true
		}
	}

	return false
}

// IsSuccess - returns true iff the response code was one of the following:
// http.StatusOK, http.StatusAccepted, http.StatusCreated. Delegates to CheckStatusCode, http.StatusNoContent.
func IsSuccess(actual int) bool {
	return CheckStatusCode(actual, http.StatusNoContent, http.StatusOK, http.StatusAccepted, http.StatusCreated)
}

// Get return a get request for the given endpoint
func Get(ctx context.Context, endpoint string) (*http.Request, error) {
	return http.NewRequestWithContext(ctx, http.MethodGet, endpoint, strings.NewReader(""))
}

// ParseForm automatically triggers a parse of the request form.
func ParseForm(original http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			log.Println("unable to parse form", err)
			http.Error(w, "malformatted form", http.StatusBadRequest)
			return
		}

		if mtype, _, err := mime.ParseMediaType(r.Header.Get("Content-Type")); err == nil && mtype == "multipart/form-data" {
			if err := r.ParseMultipartForm(bytesx.MiB); err != nil {
				log.Println("unable to parse multipart form", err)
				http.Error(w, "malformatted form", http.StatusBadRequest)
				return
			}
		}

		original.ServeHTTP(w, r)
	})
}

// RouteInvoked wraps a http.Handler and emits route invocations.
func RouteInvoked(original http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		p := HTTPRequestScheme(req) + "://" + req.Host + req.URL.Path
		started := time.Now()
		log.Println(p, req.Method, "initiated")
		original.ServeHTTP(resp, req)
		log.Println(p, req.Method, "completed", time.Since(started))
	})
}

// RouteRateLimited applies a rate limit to the http handler.
func RouteRateLimited(l *rate.Limiter) func(http.Handler) http.Handler {
	return func(original http.Handler) http.Handler {
		attempts := int64(0)
		b := backoffx.New(
			backoffx.Exponential(32*time.Millisecond),
			backoffx.Maximum(2*time.Second),
		)

		return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
			if l.Allow() {
				atomic.StoreInt64(&attempts, 0)
				original.ServeHTTP(resp, req)
				return
			}

			nattempt := int(atomic.AddInt64(&attempts, 1))
			resp.Header().Add("Retry-After", fmt.Sprintf("%d", int(b.Backoff(nattempt)/time.Second)))
			resp.WriteHeader(http.StatusTooManyRequests)
		})
	}
}

// DebugRequest dumps the request to STDERR.
func DebugRequest(original http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		raw, err := httputil.DumpRequest(req, true)
		if err != nil {
			log.Println(errorsx.Wrap(err, "failed to dump request"))
		} else {
			log.Println(string(raw))
		}
		for idx, c := range req.Cookies() {
			log.Println("COOKIE:", idx, ":", c.Name)
		}
		original.ServeHTTP(resp, req)
	})
}

func GetOrigin(r *http.Request) string {
	return strings.TrimSuffix(stringsx.FirstNonBlank(r.Header.Get("origin"), r.Header.Get("referer")), "/")
}

// RecordRequestHandler records the request to a temporary file.
// does not clean up the file from disk.
func RecordRequestHandler(original http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		var (
			err error
			raw []byte
			out *os.File
		)

		if raw, err = httputil.DumpRequest(req, true); err != nil {
			log.Println("failed to dump request", err)
			goto next
		}

		if out, err = os.CreateTemp("", "request-recording"); err != nil {
			log.Println("failed to record request", err)
			goto next
		}
		defer out.Close()

		if _, err = out.Write(raw); err != nil {
			log.Println("failed to record contents to file", err)
			goto next
		}
	next:
		original.ServeHTTP(resp, req)
	})
}

func HTTPRequestIPv6(r *http.Request) net.IP {
	return netx.DefaultIfZero(net.IPv6unspecified, netx.IP(r.RemoteAddr).To16())
}

// HTTPRequestScheme return the http scheme for a request.
func HTTPRequestScheme(req *http.Request) string {
	const (
		scheme       = "http"
		secureScheme = "https"
	)

	if req.TLS == nil {
		return scheme
	}

	return secureScheme
}

// HTTPRequestOrigin return the http origin for a request.
func HTTPRequestOrigin(req *http.Request) string {
	o := req.Header.Get("ORIGIN")
	o = strings.TrimPrefix(o, HTTPRequestScheme(req))
	o = strings.TrimPrefix(o, "://")

	return o
}

func HTTPRequestURL(req *http.Request) string {
	return HTTPRequestScheme(req) + "://" + req.Host
}

// WebsocketRequestScheme return the websocket scheme for a request.
func WebsocketRequestScheme(req *http.Request) string {
	const (
		scheme       = "ws"
		secureScheme = "wss"
	)

	if req.TLS == nil {
		return scheme
	}

	return secureScheme
}

type JSONError struct {
	Reason string `json:"reason"`
}

func NewJSONError(reason string) JSONError {
	return JSONError{Reason: reason}
}

// WriteJSON writes a json payload into the provided buffer and sets the context-type header to application/json.
func WriteJSON(resp http.ResponseWriter, buffer *bytes.Buffer, context interface{}) error {
	var (
		err error
	)

	buffer.Reset()
	resp.Header().Set("Content-Type", "application/json")

	if err = json.NewEncoder(buffer).Encode(context); err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		return err
	}

	_, err = io.Copy(resp, buffer)
	return err
}

// WriteEmptyJSON emits empty json with the provided status code.
func WriteJSONCode(resp http.ResponseWriter, code int, buffer *bytes.Buffer, context interface{}) (err error) {
	buffer.Reset()

	resp.WriteHeader(code)
	resp.Header().Set("Content-Type", "application/json")

	if err = json.NewEncoder(buffer).Encode(context); err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		return err
	}

	_, err = io.Copy(resp, buffer)
	return err
}

func WriteIO(resp http.ResponseWriter, mimetype string, code int, r io.Reader) {
	resp.Header().Set("Content-Type", mimetype)
	resp.WriteHeader(code)
	if _, err := io.CopyN(resp, r, 16*bytesx.KiB); err != nil {
		log.Println("unable to copy io to response", err)
	}
}

// WriteEmptyJSONArray emits an empty json array with the provided status code.
func WriteEmptyJSONArray(resp http.ResponseWriter, code int) (err error) {
	const emptyJSON = "[]"
	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(code)
	_, err = resp.Write([]byte(emptyJSON))
	return err
}

// WriteEmptyJSON emits empty json with the provided status code.
func WriteEmptyJSON(resp http.ResponseWriter, code int) (err error) {
	const emptyJSON = "{}"
	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(code)
	_, err = resp.Write([]byte(emptyJSON))
	return err
}

// RedirectHTTPRequest generates a url to redirect to from the provided
// request and destination node
func RedirectHTTPRequest(req *http.Request, dst string, defaultPort string) *url.URL {
	_, port, err := net.SplitHostPort(req.Host)
	if err != nil {
		debugx.Println("using default port error splitting request host", err)
		port = defaultPort
	}

	return &url.URL{
		Scheme:   HTTPRequestScheme(req),
		Host:     net.JoinHostPort(dst, port),
		Path:     req.URL.Path,
		RawQuery: req.URL.Query().Encode(),
	}
}

// AsError converts http response codes to errorsx.
func AsError(r *http.Response, err error) (*http.Response, error) {
	if err != nil {
		return r, err
	}

	if r.StatusCode >= 400 {
		return r, &Error{Code: r.StatusCode, cause: errorsx.New(r.Status)}
	}

	return r, nil
}

// TryClose attempts to close the response body if it exists.
func TryClose(r *http.Response) error {
	if r == nil {
		return nil
	}

	return r.Body.Close()
}

// ErrorCode ...
func ErrorCode(resp *http.Response) error {
	if resp.StatusCode < 400 {
		return nil
	}

	return &Error{Code: resp.StatusCode, cause: errorsx.New(resp.Status)}
}

// Error ...
type Error struct {
	Code  int
	cause error
}

func (t Error) Error() string {
	return t.cause.Error()
}

// IgnoreError ...
func IgnoreError(err error, code ...int) bool {
	var (
		cause Error
		ok    bool
	)

	if cause, ok = errorsx.Cause(err).(Error); !ok {
		return false
	}

	return CheckStatusCode(cause.Code, code...)
}

// MimeType extracts mimetype from request, defaults to application/
func MimeType(h http.Header) string {
	const fallback = "application/octet-stream"
	t, _, err := mime.ParseMediaType(h.Get("Content-Type"))
	if err != nil {
		return fallback
	}

	return stringsx.DefaultIfBlank(t, fallback)
}

type notFoundHandler struct{}

func (t notFoundHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	raw, _ := httputil.DumpRequest(req, false)
	log.Println("requested endpoint not found", string(raw))
	resp.WriteHeader(http.StatusNotFound)
}

// NotFound handles paths that do not exist.
func NotFound(c alice.Chain) http.Handler {
	return c.Then(notFoundHandler{})
}

func NewUpload(fieldname, filename string, in io.Reader) (mime string, b io.ReadCloser, err error) {
	r, w := io.Pipe()

	mp := multipart.NewWriter(w)
	go func() {
		defer w.Close()
		defer mp.Close()
		part, cerr := mp.CreateFormFile(fieldname, filename)
		if cerr != nil {
			w.CloseWithError(cerr)
			return
		}

		if _, cerr = io.Copy(part, in); cerr != nil {
			w.CloseWithError(cerr)
			return
		}
	}()

	return mp.FormDataContentType(), r, nil
}

func Unauthorized(r http.ResponseWriter, cause error) {
	errorsx.Log(log.Output(2, fmt.Sprintln(cause)))
	r.WriteHeader(http.StatusUnauthorized)
}

func Forbidden(r http.ResponseWriter, cause error) {
	errorsx.Log(log.Output(2, fmt.Sprintln(cause)))
	r.WriteHeader(http.StatusForbidden)
}

func ErrorHeader(r http.ResponseWriter, code int, cause error) {
	errorsx.Log(log.Output(2, fmt.Sprintln(cause)))
	r.WriteHeader(code)
}

// determines the ip address of the client, if it can't it'll fallback
// to all zeros.
func MaybeIP(req *http.Request) net.IP {
	host, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		log.Println("unable to split host and port", err)
		return net.IPv6zero
	}

	if i := net.ParseIP(host); i != nil {
		return i
	}

	return net.IPv6zero
}

func AutoClose(r *http.Response) error {
	if r == nil {
		return nil
	}

	return r.Body.Close()
}
