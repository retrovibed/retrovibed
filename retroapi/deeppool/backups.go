package deeppool

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/retrovibed/retrovibed/retroapi/env"
	"github.com/retrovibed/retrovibed/retroapi/errorsx"
	"github.com/retrovibed/retrovibed/retroapi/internal/httpx"
)

func NewBackups(c *http.Client) Backups {
	return Backups{
		c:        c,
		endpoint: env.Deeppool(),
	}
}

// encrypted database backups. the server stores what it is handed and never holds a key.
type Backups struct {
	c        *http.Client
	endpoint string
}

// the account's half of the encryption key.
func (t Backups) Seed(ctx context.Context) (string, error) {
	var (
		bsr BackupSeedResponse
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("https://%s/backups/seed", t.endpoint), nil)
	if err != nil {
		return "", err
	}

	resp, err := httpx.AsError(t.c.Do(req))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&bsr); err != nil {
		return "", err
	}

	return bsr.Seed, nil
}

func (t Backups) Upload(ctx context.Context, device string, mimetype string, r io.Reader) (m *Media, _ error) {
	var (
		mur MediaUploadResponse
	)

	contenttype, data, err := httpx.Multipart(func(w *multipart.Writer) error {
		part, lerr := w.CreatePart(httpx.NewMultipartHeader(mimetype, "content", "bin"))
		if lerr != nil {
			return errorsx.Wrap(lerr, "unable to create backup part")
		}

		if _, lerr = io.Copy(part, r); lerr != nil {
			return errorsx.Wrap(lerr, "unable to copy backup")
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	defer data.Close()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("https://%s/backups/%s", t.endpoint, device), data)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contenttype)

	resp, err := httpx.AsError(t.c.Do(req))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&mur); err != nil {
		return nil, err
	}

	return mur.Media, nil
}

func (t Backups) Latest(ctx context.Context) (m *Media, _ error) {
	var (
		mfr MediaFindResponse
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("https://%s/backups/latest", t.endpoint), nil)
	if err != nil {
		return nil, err
	}

	resp, err := httpx.AsError(t.c.Do(req))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&mfr); err != nil {
		return nil, err
	}

	return mfr.Media, nil
}

func (t Backups) Download(ctx context.Context, id string, into io.Writer) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("https://%s/backups/%s/download", t.endpoint, id), nil)
	if err != nil {
		return err
	}

	resp, err := httpx.AsError(t.c.Do(req))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(into, resp.Body)
	return err
}
