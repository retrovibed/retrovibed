package deeppool

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/retrovibed/retrovibed/retroapi/env"
	"github.com/retrovibed/retrovibed/retroapi/errorsx"
	"github.com/retrovibed/retrovibed/retroapi/httpx"
)

func NewMediaMetadata(c *http.Client) MediaMetadataClient {
	return MediaMetadataClient{
		c:        c,
		endpoint: env.Deeppool(),
	}
}

// golden copy of the generated media metadata for all sources.
type MediaMetadataClient struct {
	c        *http.Client
	endpoint string
}

// date known media imports should resume from for the matching sources.
func (t MediaMetadataClient) Env(ctx context.Context, include []string, exclude []string) (*MediaMetadataEnvResponse, error) {
	var (
		msg MediaMetadataEnvResponse
	)

	params := url.Values{"include_source": include, "exclude_source": exclude}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("https://%s/mm/env?%s", t.endpoint, params.Encode()), nil)
	if err != nil {
		return nil, err
	}

	resp, err := httpx.AsError(t.c.Do(req))
	if err != nil {
		return nil, errorsx.Wrap(err, "env request failed")
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&msg); err != nil {
		return nil, err
	}

	return &msg, nil
}

func (t MediaMetadataClient) Upload(ctx context.Context, items ...*MediaMetadata) (*MediaMetadataUploadResponse, error) {
	var (
		msg MediaMetadataUploadResponse
	)

	encoded, err := json.Marshal(&MediaMetadataUploadRequest{Items: items})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("https://%s/mm/", t.endpoint), bytes.NewReader(encoded))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpx.AsError(t.c.Do(req))
	if err != nil {
		return nil, errorsx.Wrap(err, "upload request failed")
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&msg); err != nil {
		return nil, err
	}

	return &msg, nil
}

// Sync pages through all media metadata changed after the cursor (a uuidv7 sync id).
func (t MediaMetadataClient) Sync(ctx context.Context, cursor string) (*MediaMetadataSyncResponse, error) {
	var (
		msg MediaMetadataSyncResponse
	)

	params := url.Values{}
	params.Set("sync", cursor)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("https://%s/mm/sync?%s", t.endpoint, params.Encode()), nil)
	if err != nil {
		return nil, err
	}

	resp, err := httpx.AsError(t.c.Do(req))
	if err != nil {
		return nil, errorsx.Wrap(err, "sync request failed")
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&msg); err != nil {
		return nil, err
	}

	return &msg, nil
}

func (t MediaMetadataClient) Search(ctx context.Context, sreq *MediaMetadataSearchRequest) (*MediaMetadataSearchResponse, error) {
	var (
		msg MediaMetadataSearchResponse
	)

	params := url.Values{"source": sreq.Source}
	params.Set("query", sreq.Query)
	params.Set("mimetype", sreq.Mimetype)
	params.Set("offset", strconv.FormatUint(sreq.Offset, 10))
	params.Set("limit", strconv.FormatUint(sreq.Limit, 10))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("https://%s/mm/?%s", t.endpoint, params.Encode()), nil)
	if err != nil {
		return nil, err
	}

	resp, err := httpx.AsError(t.c.Do(req))
	if err != nil {
		return nil, errorsx.Wrap(err, "search request failed")
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&msg); err != nil {
		return nil, err
	}

	return &msg, nil
}
