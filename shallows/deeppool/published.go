package deeppool

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"

	deepool "github.com/retrovibed/retrovibed/retroapi/deeppool"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/formx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"

	"github.com/retrovibed/retrovibed/shallows/internal/mimex"
	"github.com/retrovibed/retrovibed/shallows/meta"
)

func NewPublished(c *http.Client) Published {
	return Published{
		c:        c,
		endpoint: deepool.Deeppool(),
	}
}

type Published struct {
	c        *http.Client
	endpoint string
}

// Search returns communities matching the query from deeppool.
func (t Published) Search(ctx context.Context, query string, offset, limit uint64) (*meta.CommunitySearchResponse, error) {
	var (
		err  error
		req  *http.Request
		resp *http.Response
		msg  meta.CommunitySearchResponse
	)

	params, err := formx.NewEncoder().Encode(&meta.CommunitySearchRequest{Query: query, Offset: offset, Limit: limit})
	if err != nil {
		return nil, err
	}

	uri := fmt.Sprintf("https://%s/c/?%s", t.endpoint, params.Encode())
	if req, err = http.NewRequestWithContext(ctx, http.MethodGet, uri, nil); err != nil {
		return nil, err
	}

	if resp, err = httpx.AsError(t.c.Do(req)); err != nil {
		return nil, errorsx.Wrap(err, "request failed")
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&msg); err != nil {
		return nil, err
	}

	return &msg, nil
}

// Find returns community info from deeppool.
func (t Published) Find(ctx context.Context, communityID string) (*meta.Community, error) {
	var (
		err  error
		req  *http.Request
		resp *http.Response
		msg  meta.CommunityFindResponse
	)

	uri := fmt.Sprintf("https://%s/c/%s", t.endpoint, communityID)
	if req, err = http.NewRequestWithContext(ctx, http.MethodGet, uri, nil); err != nil {
		return nil, err
	}

	if resp, err = httpx.AsError(t.c.Do(req)); err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&msg); err != nil {
		return nil, err
	}

	return msg.Community, nil
}

// PublishedContentSyncRequest is the cursor-based request for the sync endpoint.
type PublishedContentSyncRequest struct {
	Sync   string `json:"sync"`
	Offset uint64 `json:"offset"`
	Limit  uint64 `json:"limit"`
}

// PublishedContentSyncResponse is the response from the sync endpoint.
type PublishedContentSyncResponse struct {
	Next  *PublishedContentSyncRequest `json:"next"`
	Items []*meta.PublishedContent     `json:"items"`
}

// Sync pages through all published content using a cursor. cursor is the last-seen ID
// (use empty string to start from the beginning).
func (t Published) Sync(ctx context.Context, cursor string) (*PublishedContentSyncResponse, error) {
	var (
		err  error
		req  *http.Request
		resp *http.Response
		msg  PublishedContentSyncResponse
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
func (t Published) List(ctx context.Context, communityID string) (*meta.PublishedContentSearchResponse, error) {
	var (
		err  error
		req  *http.Request
		resp *http.Response
		msg  meta.PublishedContentSearchResponse
	)

	uri := fmt.Sprintf("https://%s/c/%s/published", t.endpoint, communityID)
	if req, err = http.NewRequestWithContext(ctx, http.MethodGet, uri, nil); err != nil {
		return nil, err
	}

	if resp, err = httpx.AsError(t.c.Do(req)); err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&msg); err != nil {
		return nil, err
	}

	return &msg, nil
}

// UploadFeed uploads an RSS feed to deeppool for a community.
func (t Published) UploadFeed(ctx context.Context, communityID string, feed io.Reader) error {
	var (
		err  error
		req  *http.Request
		resp *http.Response
		part io.Writer
	)

	contenttype, data, err := httpx.Multipart(func(w *multipart.Writer) error {
		if part, err = w.CreatePart(httpx.NewMultipartHeader(mimex.RSS, "content", "feed.xml")); err != nil {
			return errorsx.Wrap(err, "unable to create feed part")
		}

		if _, err = io.Copy(part, feed); err != nil {
			return errorsx.Wrap(err, "unable to copy feed")
		}

		return nil
	})
	if err != nil {
		return err
	}
	defer data.Close()

	uri := fmt.Sprintf("https://%s/c/%s", t.endpoint, communityID)
	if req, err = http.NewRequestWithContext(ctx, http.MethodPost, uri, data); err != nil {
		return err
	}
	req.Header.Set("Content-Type", contenttype)

	if resp, err = httpx.AsError(t.c.Do(req)); err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// Publish publishes content to a community in deeppool.
func (t Published) Publish(ctx context.Context, communityID string, pc *meta.PublishedContent) (*meta.PublishContentResponse, error) {
	var (
		err  error
		req  *http.Request
		resp *http.Response
		msg  meta.PublishContentResponse
	)

	body, err := json.Marshal(&meta.PublishContentRequest{PublishedContent: pc})
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
