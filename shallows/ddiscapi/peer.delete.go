package ddiscapi

import (
	"context"
	"fmt"
	"net/http"

	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/tracking"
)

// PeerDelete removes a peer, identified by its hex encoded id, on the given library endpoint.
func PeerDelete(ctx context.Context, c *http.Client, endpoint string, hexPeerID string) (resp *PeerDeleteResponse, err error) {
	infohash, err := int160.FromHexEncodedString(hexPeerID)
	if err != nil {
		return nil, errorsx.Wrap(err, "failed to decode peer id")
	}

	hreq, err := http.NewRequestWithContext(ctx, http.MethodDelete, fmt.Sprintf("%s/ddisc/%s", endpoint, tracking.PeerUID(infohash)), nil)
	if err != nil {
		return nil, errorsx.Wrap(err, "unable to create http request")
	}

	hresp, err := httpx.AsError(c.Do(hreq))
	if err != nil {
		return nil, errorsx.Wrap(err, "http request failed")
	}

	resp = new(PeerDeleteResponse)
	if err = httpx.DecodeJSON(hresp, resp); err != nil {
		return nil, errorsx.Wrap(err, "unable to decode response")
	}

	return resp, nil
}
