package ddiscapi

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/retrovibed/retrovibed/retroapi/jsonx"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
)

// MediaCreate creates a media entry on the given library endpoint.
func MediaCreate(ctx context.Context, c *http.Client, endpoint string, req *MediaCreateRequest) (resp *MediaCreateResponse, err error) {
	encoded, err := jsonx.Marshal(req)
	if err != nil {
		return nil, errorsx.Wrap(err, "unable to encode request")
	}

	hreq, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/ddisc/media/", endpoint), bytes.NewReader(encoded))
	if err != nil {
		return nil, errorsx.Wrap(err, "unable to create http request")
	}

	hresp, err := httpx.AsError(c.Do(hreq))
	if err != nil {
		return nil, errorsx.Wrap(err, "http request failed")
	}

	resp = new(MediaCreateResponse)
	if err = httpx.DecodeJSON(hresp, resp); err != nil {
		return nil, errorsx.Wrap(err, "unable to decode response")
	}

	return resp, nil
}
