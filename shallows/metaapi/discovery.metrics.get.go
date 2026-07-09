package metaapi

import (
	"context"
	"fmt"
	"net/http"

	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
)

// DiscoveryMetrics reports the state of the infohash identification pipeline on the given library endpoint.
func DiscoveryMetrics(ctx context.Context, c *http.Client, endpoint string) (resp *DiscoveryMetricsResponse, err error) {
	hreq, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("https://%s/diagnostics/discovery/", endpoint), nil)
	if err != nil {
		return nil, errorsx.Wrap(err, "unable to create http request")
	}

	hresp, err := httpx.AsError(c.Do(hreq))
	if err != nil {
		return nil, errorsx.Wrap(err, "http request failed")
	}

	resp = new(DiscoveryMetricsResponse)
	if err = httpx.DecodeJSON(hresp, resp); err != nil {
		return nil, errorsx.Wrap(err, "unable to decode response")
	}

	return resp, nil
}
