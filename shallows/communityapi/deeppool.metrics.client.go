package communityapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/retrovibed/retrovibed/retroapi/env"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
)

func NewMetrics(c *http.Client) Metrics {
	return Metrics{
		c:        c,
		endpoint: env.Deeppool(),
	}
}

type Metrics struct {
	c        *http.Client
	endpoint string
}

func (t Metrics) Sync(ctx context.Context, communityID string) (*MetricsSyncResponse, error) {
	var (
		err  error
		req  *http.Request
		resp *http.Response
		msg  MetricsSyncResponse
	)

	uri := fmt.Sprintf("https://%s/c/%s/metrics/sync", t.endpoint, communityID)
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

func (t Metrics) Publish(ctx context.Context, content *PublishContentRequest, torrent io.Reader) (*PublishContentResponse, error) {
	var (
		err  error
		req  *http.Request
		resp *http.Response
		msg  PublishContentResponse
	)

	uri := fmt.Sprintf("https://%s/p/", t.endpoint)
	contentType, body, err := httpx.Multipart(func(w *multipart.Writer) error {
		metadata, lerr := w.CreatePart(httpx.NewMultipartHeader("application/json", "metadata", "metadata.json"))
		if lerr != nil {
			return lerr
		}

		if lerr = json.NewEncoder(metadata).Encode(content); lerr != nil {
			return lerr
		}

		if torrent == nil {
			return nil
		}

		tp, lerr := w.CreatePart(httpx.NewMultipartHeader("application/x-bittorrent", "torrent", "content.torrent"))
		if lerr != nil {
			return lerr
		}

		_, lerr = io.Copy(tp, torrent)
		return lerr
	})
	if err != nil {
		return nil, err
	}
	defer body.Close()

	if req, err = http.NewRequestWithContext(ctx, http.MethodPost, uri, body); err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)

	if resp, err = httpx.AsError(t.c.Do(req)); err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&msg); err != nil {
		return nil, err
	}

	return &msg, nil
}
