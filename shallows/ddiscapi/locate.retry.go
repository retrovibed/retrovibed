package ddiscapi

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/retrovibed/retrovibed/retroapi/httpx"
	"github.com/retrovibed/retrovibed/retroapi/jsonx"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
)

// LocateRetry clears the cooldown on a locate request, identified by id, on the given library endpoint.
func LocateRetry(ctx context.Context, c *http.Client, endpoint string, id string, req *LocateRetryRequest) (resp *LocateRetryResponse, err error) {
	encoded, err := jsonx.Marshal(req)
	if err != nil {
		return nil, errorsx.Wrap(err, "unable to encode request")
	}

	hreq, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/l/%s/retry", endpoint, id), bytes.NewReader(encoded))
	if err != nil {
		return nil, errorsx.Wrap(err, "unable to create http request")
	}

	hresp, err := httpx.AsError(c.Do(hreq))
	if err != nil {
		return nil, errorsx.Wrap(err, "http request failed")
	}

	resp = new(LocateRetryResponse)
	if err = httpx.DecodeJSON(hresp, resp); err != nil {
		return nil, errorsx.Wrap(err, "unable to decode response")
	}

	return resp, nil
}
