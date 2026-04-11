package metaapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/retrovibed/retrovibed/deeppool"
	"github.com/retrovibed/retrovibed/internal/httpx"
	"github.com/retrovibed/retrovibed/internal/mimex"
	"github.com/retrovibed/retrovibed/meta"
)

func CommunityCreate(ctx context.Context, c *http.Client, com *meta.CommunityCreateRequest) (resp *meta.CommunityCreateResponse, err error) {
	encoded, err := json.Marshal(com)
	if err != nil {
		return nil, err
	}
	_resp, err := httpx.AsError(c.Post(fmt.Sprintf("https://%s/c/", deeppool.Deeppool()), mimex.JSON, bytes.NewReader(encoded)))
	if err != nil {
		return nil, err
	}

	resp = new(meta.CommunityCreateResponse)

	if err = json.NewDecoder(_resp.Body).Decode(resp); err != nil {
		return nil, err
	}

	return resp, nil
}
