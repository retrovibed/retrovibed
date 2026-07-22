package ddiscapi

import (
	"context"
	"fmt"
	"net/http"

	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
)

// MediaDelete removes a media entry, identified by id, on the given library endpoint.
func MediaDelete(ctx context.Context, c *http.Client, endpoint string, id string) (resp *MediaDeleteResponse, err error) {
	hreq, err := http.NewRequestWithContext(ctx, http.MethodDelete, fmt.Sprintf("%s/ddisc/media/%s", endpoint, id), nil)
	if err != nil {
		return nil, errorsx.Wrap(err, "unable to create http request")
	}

	hresp, err := httpx.AsError(c.Do(hreq))
	if err != nil {
		return nil, errorsx.Wrap(err, "http request failed")
	}

	resp = new(MediaDeleteResponse)
	if err = httpx.DecodeJSON(hresp, resp); err != nil {
		return nil, errorsx.Wrap(err, "unable to decode response")
	}

	return resp, nil
}
