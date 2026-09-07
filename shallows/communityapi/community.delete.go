package communityapi

import (
	"context"
	"fmt"
	"net/http"

	"github.com/retrovibed/retrovibed/retroapi/env"
	"github.com/retrovibed/retrovibed/retroapi/jsonx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
)

func CommunityDelete(ctx context.Context, c *http.Client, domainOrId string) (resp *CommunityDeleteResponse, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, fmt.Sprintf("https://%s/c/%s", env.Deeppool(), domainOrId), nil)
	if err != nil {
		return nil, err
	}

	_resp, err := httpx.AsError(c.Do(req))
	if err != nil {
		return nil, err
	}

	resp = new(CommunityDeleteResponse)

	if err = jsonx.UnmarshalRead(_resp.Body, resp); err != nil {
		return nil, err
	}

	return resp, nil
}
