package ddiscapi

import (
	"context"
	"fmt"
	"net/http"

	"github.com/retrovibed/retrovibed/retroapi/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
)

// LocateFind looks up a single locate request, identified by id, on the given library endpoint.
func LocateFind(ctx context.Context, c *http.Client, endpoint string, id string) (resp *LocateLookupResponse, err error) {
	hreq, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/l/%s", endpoint, id), nil)
	if err != nil {
		return nil, errorsx.Wrap(err, "unable to create http request")
	}

	hresp, err := httpx.AsError(c.Do(hreq))
	if err != nil {
		return nil, errorsx.Wrap(err, "http request failed")
	}

	resp = new(LocateLookupResponse)
	if err = httpx.DecodeJSON(hresp, resp); err != nil {
		return nil, errorsx.Wrap(err, "unable to decode response")
	}

	return resp, nil
}
