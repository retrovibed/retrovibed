package media

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gabriel-vasile/mimetype"
	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retrovibed/retroapi/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
)

// NewUploader uploads files into the library served at endpoint.
func NewUploader(c *http.Client, endpoint string) Uploader {
	return Uploader{
		c:        c,
		endpoint: endpoint,
	}
}

type Uploader struct {
	c        *http.Client
	endpoint string
}

// Upload the file at path into the library, filing it into the given library directory,
// returning the id of the resulting media.
func (t Uploader) Upload(ctx context.Context, path string, directoryID string) (_ string, err error) {
	var (
		decoded MediaUploadResponse
	)

	mtype, err := mimetype.DetectFile(path)
	if err != nil {
		return "", errorsx.Wrapf(err, "unable to detect mimetype: %s", path)
	}

	contentType, body, err := httpx.Multipart(func(w *multipart.Writer) error {
		// the library root is the upload default, only name a destination when one was configured.
		if directoryID != uuid.Nil.String() {
			if lerr := w.WriteField("directory_id", directoryID); lerr != nil {
				return errorsx.Wrap(lerr, "unable to write directory")
			}
		}

		src, lerr := os.Open(path)
		if lerr != nil {
			return errorsx.Wrapf(lerr, "unable to read %s", path)
		}
		defer src.Close()

		part, lerr := w.CreatePart(httpx.NewMultipartHeader(mtype.String(), "content", filepath.Base(path)))
		if lerr != nil {
			return errorsx.Wrap(lerr, "unable to create content part")
		}

		if _, lerr = io.Copy(part, src); lerr != nil {
			return errorsx.Wrap(lerr, "unable to copy content")
		}

		return nil
	})
	if err != nil {
		return "", errorsx.Wrap(err, "unable to build multipart request")
	}
	defer body.Close()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/m/", t.endpoint), body)
	if err != nil {
		return "", errorsx.Wrap(err, "unable to create http request")
	}
	req.Header.Set("Content-Type", contentType)

	resp, err := httpx.AsError(t.c.Do(req))
	if err != nil {
		return "", errorsx.Wrap(err, "http request failed")
	}

	if err = httpx.DecodeJSON(resp, &decoded); err != nil {
		return "", errorsx.Wrap(err, "unable to decode response")
	}

	return decoded.Media.Id, nil
}
