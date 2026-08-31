package communityapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/retrovibed/retrovibed/retroapi/env"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
)

func CommunityUpdate(ctx context.Context, c *http.Client, domainOrId string, com *CommunityUpdateRequest) (resp *CommunityUpdateResponse, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, fmt.Sprintf("https://%s/c/%s", env.Deeppool(), domainOrId), nil)
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

	resp = new(CommunityUpdateResponse)

	if err = json.NewDecoder(_resp.Body).Decode(resp); err != nil {
		return nil, err
	}

	return resp, nil
}
