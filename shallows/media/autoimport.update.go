package media

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/retrovibed/retrovibed/retroapi/httpx"
	"github.com/retrovibed/retrovibed/retroapi/jsonx"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
)

// AutoimportUpdate replaces the editable fields of a monitored directory, identified by id, on the given library endpoint.
func AutoimportUpdate(ctx context.Context, c *http.Client, endpoint string, id string, req *AutoimportDirectoryUpdateRequest) (resp *AutoimportDirectoryUpdateResponse, err error) {
	encoded, err := jsonx.Marshal(req)
	if err != nil {
		return nil, errorsx.Wrap(err, "unable to encode request")
	}

	hreq, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/autoimport/%s", endpoint, id), bytes.NewReader(encoded))
	if err != nil {
		return nil, errorsx.Wrap(err, "unable to create http request")
	}

	hresp, err := httpx.AsError(c.Do(hreq))
	if err != nil {
		return nil, errorsx.Wrap(err, "http request failed")
	}

	resp = new(AutoimportDirectoryUpdateResponse)
	if err = httpx.DecodeJSON(hresp, resp); err != nil {
		return nil, errorsx.Wrap(err, "unable to decode response")
	}

	return resp, nil
}
