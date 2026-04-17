package metaapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/retrovibed/retrovibed/retroapi/deeppool"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/meta"
)

func CommunityUpdate(ctx context.Context, c *http.Client, domainOrId string, com *meta.CommunityUpdateRequest) (resp *meta.CommunityUpdateResponse, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, fmt.Sprintf("https://%s/c/%s", deeppool.Deeppool(), domainOrId), nil)
	if err != nil {
		return nil, err
	}

	if err = httpx.EncodeJSON(req, com); err != nil {
		return nil, err
	}

	_resp, err := httpx.AsError(c.Do(req))
	if err != nil {
		return nil, err
	}

	resp = new(meta.CommunityUpdateResponse)

	if err = json.NewDecoder(_resp.Body).Decode(resp); err != nil {
		return nil, err
	}

	return resp, nil
}
