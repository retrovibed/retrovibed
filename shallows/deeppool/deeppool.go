package deeppool

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/retrovibed/retrovibed/internal/errorsx"
	"github.com/retrovibed/retrovibed/internal/httpx"
)

func NewArchiver(c *http.Client) Archiver {
	return Archiver{
		c:        c,
		endpoint: Deeppool(),
	}
}

type Archiver struct {
	c        *http.Client
	endpoint string
}

func (t Archiver) Upload(ctx context.Context, mimetype string, r io.Reader) (m *Media, _ error) {
	var (
		mur MediaUploadResponse
	)

	dst := fmt.Sprintf("https://%s/m/", t.endpoint)
	redirectreq, err := http.NewRequestWithContext(ctx, "ENDPOINT", dst, nil)
	if err != nil {
		return nil, err
	}

	resp, err := httpx.AsError(t.c.Do(redirectreq))
	if err != nil {
		return nil, err
	}
	httpx.TryClose(resp)

	if uri, ok := httpx.GetRedirect(resp); ok {
		dst = uri
	}

	contenttype, data, err := httpx.Multipart(func(w *multipart.Writer) error {
		part, lerr := w.CreatePart(httpx.NewMultipartHeader(mimetype, "content", "bin"))
		if lerr != nil {
			return errorsx.Wrap(lerr, "unable to create archive part")
		}

		if _, lerr = io.Copy(part, r); lerr != nil {
			return errorsx.Wrap(lerr, "unable to copy archive")
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	defer data.Close()

	resp, err = httpx.AsError(t.c.Post(dst, contenttype, data))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&mur); err != nil {
		return nil, err
	}

	return mur.Media, nil
}

func NewRetrieval(c *http.Client) Retrieval {
	return Retrieval{
		c:        c,
		endpoint: Deeppool(),
	}
}

type Retrieval struct {
	c        *http.Client
	endpoint string
}

func (t Retrieval) Download(ctx context.Context, id string, into io.Writer) error {
	resp, err := httpx.AsError(t.c.Get(fmt.Sprintf("https://%s/m/%s/download", t.endpoint, id)))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if _, err := io.Copy(into, resp.Body); err != nil {
		return err
	}

	return nil
}

func NewRanger(c *http.Client) Ranger {
	return Ranger{
		c:        c,
		endpoint: Deeppool(),
	}
}

type Ranger struct {
	c        *http.Client
	endpoint string
}

func (t Ranger) Download(ctx context.Context, id string, start, end uint64, into io.Writer) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("https://%s/m/%s/download", t.endpoint, id), nil)
	if err != nil {
		return err
	}
	httpx.RangeHeaders(req, start, end-1)

	resp, err := httpx.AsError(t.c.Do(req))
	if httpx.CheckStatusCode(resp.StatusCode, http.StatusForbidden, http.StatusUnauthorized) {
		err = errorsx.NewUnrecoverable(err)
	}

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if _, err := io.Copy(into, resp.Body); err != nil {
		return err
	}

	return nil
}
