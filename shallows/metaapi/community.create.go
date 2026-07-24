package metaapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/retrovibed/retrovibed/retroapi/env"
	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/shallows/communityapi"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
)

func CommunityCreate(ctx context.Context, c *http.Client, com *communityapi.CommunityCreateRequest) (resp *communityapi.CommunityCreateResponse, err error) {
	encoded, err := json.Marshal(com)
	if err != nil {
		return nil, err
	}
	_resp, err := httpx.AsError(c.Post(fmt.Sprintf("https://%s/c/", env.Deeppool()), mimex.JSON, bytes.NewReader(encoded)))
	if err != nil {
		return nil, err
	}

	resp = new(communityapi.CommunityCreateResponse)

	if err = json.NewDecoder(_resp.Body).Decode(resp); err != nil {
		return nil, err
	}

	return resp, nil
}
