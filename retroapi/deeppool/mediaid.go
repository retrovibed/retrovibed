package deeppool

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/retrovibed/retrovibed/retroapi/internal/httpx"
)

func NewMediaID(c *http.Client) MediaID {
	return MediaID{
		c:        c,
		endpoint: Deeppool(),
	}
}

type MediaID struct {
	c        *http.Client
	endpoint string
}

func (t MediaID) Clean(ctx context.Context, text string) (string, error) {
	var (
		buf  bytes.Buffer
		resp *http.Response
		msg  CleanIDResponse
	)

	if err := json.NewEncoder(&buf).Encode(&CleanIDRequest{Text: text}); err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("https://%s/clean/", t.endpoint), &buf)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	if resp, err = httpx.AsError(t.c.Do(req)); err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&msg); err != nil {
		return "", err
	}

	return msg.Text, nil
}
