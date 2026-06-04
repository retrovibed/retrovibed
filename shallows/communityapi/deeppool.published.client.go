package communityapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/retrovibed/retrovibed/retroapi/deeppool"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/formx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"

	"github.com/retrovibed/retrovibed/shallows/internal/mimex"
)

func NewDeeppoolPublished(c *http.Client) DeeppoolPublished {
	return DeeppoolPublished{
		c:        c,
		endpoint: deeppool.Deeppool(),
	}
}

type DeeppoolPublished struct {
	c        *http.Client
	endpoint string
}

// Sync pages through all published content using a cursor. cursor is the last-seen ID
// (use empty string to start from the beginning).
func (t DeeppoolPublished) Sync(ctx context.Context, cursor string) (*PublishedContentSearchResponse, error) {
	var (
		err  error
		req  *http.Request
		resp *http.Response
		msg  PublishedContentSearchResponse
	)

	params := url.Values{}
	params.Set("sync", cursor)

	uri := fmt.Sprintf("https://%s/p/sync?%s", t.endpoint, params.Encode())
	if req, err = http.NewRequestWithContext(ctx, http.MethodGet, uri, nil); err != nil {
		return nil, err
	}

	if resp, err = httpx.AsError(t.c.Do(req)); err != nil {
		return nil, errorsx.Wrap(err, "sync request failed")
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&msg); err != nil {
		return nil, err
	}

	return &msg, nil
}

// List returns all published content for a community.
func (t DeeppoolPublished) List(ctx context.Context, communityID string, req *PublishedContentSearchRequest) (*PublishedContentSearchResponse, error) {
	var (
		err  error
		r    *http.Request
		resp *http.Response
		msg  PublishedContentSearchResponse
	)

	req.CommunityId = communityID
	params, err := formx.NewEncoder().Encode(req)
	if err != nil {
		return nil, err
	}

	uri := fmt.Sprintf("https://%s/p?%s", t.endpoint, params.Encode())
	if r, err = http.NewRequestWithContext(ctx, http.MethodGet, uri, nil); err != nil {
		return nil, err
	}

	if resp, err = httpx.AsError(t.c.Do(r)); err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&msg); err != nil {
		return nil, err
	}

	return &msg, nil
}

// Publish publishes content to a community in deeppool.
func (t DeeppoolPublished) Publish(ctx context.Context, communityID string, pc *PublishedContent) (*PublishContentResponse, error) {
	var (
		err  error
		req  *http.Request
		resp *http.Response
		msg  PublishContentResponse
	)

	body, err := json.Marshal(&PublishContentRequest{PublishedContent: pc})
	if err != nil {
		return nil, err
	}

	uri := fmt.Sprintf("https://%s/c/%s/publish", t.endpoint, communityID)
	if req, err = http.NewRequestWithContext(ctx, http.MethodPost, uri, bytes.NewReader(body)); err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", mimex.JSON)

	if resp, err = httpx.AsError(t.c.Do(req)); err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&msg); err != nil {
		return nil, err
	}

	return &msg, nil
}
