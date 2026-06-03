package metaapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/retrovibed/retrovibed/retroapi/deeppool"
	"github.com/retrovibed/retrovibed/shallows/communityapi"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
)

func CommunityInfo(ctx context.Context, c *http.Client, domainOrId string) (resp *communityapi.CommunityFindResponse, err error) {
	_resp, err := httpx.AsError(c.Get(fmt.Sprintf("https://%s/c/%s", deeppool.Deeppool(), domainOrId)))
	if err != nil {
		return nil, err
	}

	resp = new(communityapi.CommunityFindResponse)

	if err = json.NewDecoder(_resp.Body).Decode(resp); err != nil {
		return nil, err
	}

	return resp, nil
}
