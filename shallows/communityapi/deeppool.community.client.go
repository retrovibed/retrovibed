package communityapi

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/retrovibed/retrovibed/retroapi/env"
	"github.com/retrovibed/retrovibed/retroapi/jsonx"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/formx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"

	"github.com/retrovibed/retrovibed/retroapi/mimex"
)

func NewDeeppoolCommunity(c *http.Client) DeeppoolCommunity {
	return DeeppoolCommunity{
		c:        c,
		endpoint: env.Deeppool(),
	}
}

type DeeppoolCommunity struct {
	c        *http.Client
	endpoint string
}

// Search returns communities matching the query from deeppool.
func (t DeeppoolCommunity) Search(ctx context.Context, query string, offset, limit uint64) (*CommunitySearchResponse, error) {
	var (
		err  error
		req  *http.Request
		resp *http.Response
		msg  CommunitySearchResponse
	)

	params, err := formx.NewEncoder().Encode(&CommunitySearchRequest{Query: query, Offset: offset, Limit: limit})
	if err != nil {
		return nil, errorsx.Wrap(err, "invalid params")
	}

	uri := fmt.Sprintf("https://%s/c/?%s", t.endpoint, params.Encode())
	if req, err = http.NewRequestWithContext(ctx, http.MethodGet, uri, nil); err != nil {
		return nil, errorsx.Wrapf(err, "request creation failed: %s", t.endpoint)
	}

	if resp, err = httpx.AsError(t.c.Do(req)); err != nil {
		return nil, errorsx.Wrapf(err, "request failed: %s", uri)
	}
	defer resp.Body.Close()

	if err = jsonx.UnmarshalRead(resp.Body, &msg); err != nil {
		return nil, err
	}

	return &msg, nil
}

// Find returns community info from deeppool.
func (t DeeppoolCommunity) Find(ctx context.Context, communityID string) (*Community, error) {
	var (
		err  error
		req  *http.Request
		resp *http.Response
		msg  CommunityFindResponse
	)

	uri := fmt.Sprintf("https://%s/c/%s", t.endpoint, communityID)
	if req, err = http.NewRequestWithContext(ctx, http.MethodGet, uri, nil); err != nil {
		return nil, err
	}

	if resp, err = httpx.AsError(t.c.Do(req)); err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err = jsonx.UnmarshalRead(resp.Body, &msg); err != nil {
		return nil, err
	}

	return msg.Community, nil
}

// UploadFeed uploads an RSS feed to deeppool for a community.
func (t DeeppoolCommunity) UploadFeed(ctx context.Context, communityID string, feed io.Reader) error {
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
