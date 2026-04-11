package metaapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/retrovibed/retrovibed/deeppool"
	"github.com/retrovibed/retrovibed/internal/httpx"
	"github.com/retrovibed/retrovibed/meta"
)

func CommunityDelete(ctx context.Context, c *http.Client, domainOrId string) (resp *meta.CommunityDeleteResponse, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, fmt.Sprintf("https://%s/c/%s", deeppool.Deeppool(), domainOrId), nil)
	if err != nil {
		return nil, err
	}

	_resp, err := httpx.AsError(c.Do(req))
	if err != nil {
		return nil, err
	}

	resp = new(meta.CommunityDeleteResponse)

	if err = json.NewDecoder(_resp.Body).Decode(resp); err != nil {
		return nil, err
	}

	return resp, nil
}
