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

func CommunityInfo(ctx context.Context, c *http.Client, domainOrId string) (resp *meta.CommunityFindResponse, err error) {
	_resp, err := httpx.AsError(c.Get(fmt.Sprintf("https://%s/c/%s", deeppool.Deeppool(), domainOrId)))
	if err != nil {
		return nil, err
	}

	resp = new(meta.CommunityFindResponse)

	if err = json.NewDecoder(_resp.Body).Decode(resp); err != nil {
		return nil, err
	}

	return resp, nil
}
