package communityapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/retrovibed/retrovibed/retroapi/env"
	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
)

func CommunityPublish(ctx context.Context, c *http.Client, id string, in io.Reader) (resp *CommunityUploadResponse, err error) {
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

	_resp, err := httpx.AsError(c.Post(fmt.Sprintf("https://%s/c/%s", env.Deeppool(), id), boundary, reader))
	if err != nil {
		return nil, err
	}

	resp = new(CommunityUploadResponse)

	if err = json.NewDecoder(_resp.Body).Decode(resp); err != nil {
		return nil, err
	}

	return resp, nil
}
