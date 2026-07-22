package ddiscapi

import (
	"context"
	"fmt"
	"net/http"

	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/formx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
)

// DiscoverySearch searches discovery entries on the given library endpoint.
func DiscoverySearch(ctx context.Context, c *http.Client, endpoint string, req *DiscoverySearchRequest) (resp *DiscoverySearchResponse, err error) {
	encoded, err := formx.NewEncoder().Encode(req)
	if err != nil {
		return nil, errorsx.Wrap(err, "unable to encode request")
	}

	hreq, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/ddisc/discovery/?%s", endpoint, encoded.Encode()), nil)
	if err != nil {
		return nil, errorsx.Wrap(err, "unable to create http request")
	}

	hresp, err := httpx.AsError(c.Do(hreq))
	if err != nil {
		return nil, errorsx.Wrap(err, "http request failed")
	}

	resp = new(DiscoverySearchResponse)
	if err = httpx.DecodeJSON(hresp, resp); err != nil {
		return nil, errorsx.Wrap(err, "unable to decode response")
	}

	return resp, nil
}
