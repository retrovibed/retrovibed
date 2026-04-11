package metaapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/retrovibed/retrovibed/deeppool"
	"github.com/retrovibed/retrovibed/internal/errorsx"
	"github.com/retrovibed/retrovibed/internal/httpx"
	"github.com/retrovibed/retrovibed/internal/mimex"
	"github.com/retrovibed/retrovibed/meta"
)

func CommunityPublish(ctx context.Context, c *http.Client, id string, in io.Reader) (resp *meta.CommunityUploadResponse, err error) {
	boundary, reader, err := httpx.Multipart(func(w *multipart.Writer) error {
		part, lerr := w.CreatePart(httpx.NewMultipartHeader(mimex.RSS, "feed", mimex.RSS))
		if lerr != nil {
			return errorsx.Wrap(lerr, "unable to create feed part")
		}

		if _, lerr = io.Copy(part, in); lerr != nil {
			return errorsx.Wrap(lerr, "unable to copy feed")
		}

		return nil
	})
	defer reader.Close()

	_resp, err := httpx.AsError(c.Post(fmt.Sprintf("https://%s/c/%s", deeppool.Deeppool(), id), boundary, reader))
	if err != nil {
		return nil, err
	}

	resp = new(meta.CommunityUploadResponse)

	if err = json.NewDecoder(_resp.Body).Decode(resp); err != nil {
		return nil, err
	}

	return resp, nil
}
