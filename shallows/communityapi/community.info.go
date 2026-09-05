package communityapi

import (
	"context"
	"fmt"
	"net/http"

	"github.com/retrovibed/retrovibed/retroapi/env"
	"github.com/retrovibed/retrovibed/retroapi/jsonx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
)

func CommunityInfo(ctx context.Context, c *http.Client, domainOrId string) (resp *CommunityFindResponse, err error) {
	_resp, err := httpx.AsError(c.Get(fmt.Sprintf("https://%s/c/%s", env.Deeppool(), domainOrId)))
	if err != nil {
		return nil, err
	}

	resp = new(CommunityFindResponse)

	if err = jsonx.UnmarshalRead(_resp.Body, resp); err != nil {
		return nil, err
	}

	return resp, nil
}
