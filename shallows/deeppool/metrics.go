package deeppool

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/retrovibed/retrovibed/internal/httpx"
	"github.com/retrovibed/retrovibed/meta"
)

func NewMetrics(c *http.Client) Metrics {
	return Metrics{
		c:        c,
		endpoint: Deeppool(),
	}
}

type Metrics struct {
	c        *http.Client
	endpoint string
}

func (t Metrics) Sync(ctx context.Context, communityID string) (*meta.MetricsSyncResponse, error) {
	var (
		err  error
		req  *http.Request
		resp *http.Response
		msg  meta.MetricsSyncResponse
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

func (t Metrics) Publish(ctx context.Context, communityID string, content *meta.PublishContentRequest) (*meta.PublishContentResponse, error) {
	var (
		err  error
		req  *http.Request
		resp *http.Response
		msg  meta.PublishContentResponse
	)

	uri := fmt.Sprintf("https://%s/c/%s/publish", t.endpoint, communityID)
	body, err := json.Marshal(content)
	if err != nil {
		return nil, err
	}

	if req, err = http.NewRequestWithContext(ctx, http.MethodPost, uri, bytes.NewReader(body)); err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	if resp, err = httpx.AsError(t.c.Do(req)); err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&msg); err != nil {
		return nil, err
	}

	return &msg, nil
}
