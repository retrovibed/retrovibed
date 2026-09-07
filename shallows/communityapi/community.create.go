package communityapi

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/retrovibed/retrovibed/retroapi/env"
	"github.com/retrovibed/retrovibed/retroapi/jsonx"
	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
)

func CommunityCreate(ctx context.Context, c *http.Client, com *CommunityCreateRequest) (resp *CommunityCreateResponse, err error) {
	encoded, err := jsonx.Marshal(com)
	if err != nil {
		return nil, err
	}
	_resp, err := httpx.AsError(c.Post(fmt.Sprintf("https://%s/c/", env.Deeppool()), mimex.JSON, bytes.NewReader(encoded)))
	if err != nil {
		return nil, err
	}

	resp = new(CommunityCreateResponse)

	if err = jsonx.UnmarshalRead(_resp.Body, resp); err != nil {
		return nil, err
	}

	return resp, nil
}
