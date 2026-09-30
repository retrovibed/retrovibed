package media

import (
	"context"
	"fmt"
	"net/http"

	"github.com/retrovibed/retrovibed/retroapi/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
)

// AutoimportFind retrieves a monitored directory, identified by id, on the given library endpoint.
func AutoimportFind(ctx context.Context, c *http.Client, endpoint string, id string) (resp *AutoimportDirectoryLookupResponse, err error) {
	hreq, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/autoimport/%s", endpoint, id), nil)
	if err != nil {
		return nil, errorsx.Wrap(err, "unable to create http request")
	}

	hresp, err := httpx.AsError(c.Do(hreq))
	if err != nil {
		return nil, errorsx.Wrap(err, "http request failed")
	}

	resp = new(AutoimportDirectoryLookupResponse)
	if err = httpx.DecodeJSON(hresp, resp); err != nil {
		return nil, errorsx.Wrap(err, "unable to decode response")
	}

	return resp, nil
}
