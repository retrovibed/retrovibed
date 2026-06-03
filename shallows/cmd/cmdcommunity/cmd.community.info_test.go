package cmdcommunity

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/retrovibed/retrovibed/shallows/communityapi"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/testx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommunityInfo(t *testing.T) {
	t.Run("retrieves community info successfully", func(t *testing.T) {
		ctx, cancel := testx.Context(t)
		defer cancel()

		expected := &communityapi.Community{
			Id:          "test-id",
			Domain:      "test-community",
			Description: "test description",
			Mimetype:    "video/mp4",
		}

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodGet, r.Method)
			require.Equal(t, "/c/test-community", r.URL.Path)
			assert.NoError(t, json.NewEncoder(w).Encode(&communityapi.CommunityFindResponse{Community: expected}))
		}))
		defer srv.Close()

		c := &http.Client{}
		c.Transport = httpx.RewriteHostTransport(testx.Must(url.ParseRequestURI(srv.URL))(t), c.Transport)

		cmd := cmdCommunityInfo{Name: "test-community"}
		err := cmd.run(ctx, c, strings.NewReader(""), &bytes.Buffer{})
		require.NoError(t, err)
	})

	t.Run("returns error on server failure", func(t *testing.T) {
		ctx, cancel := testx.Context(t)
		defer cancel()

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer srv.Close()

		c := &http.Client{}
		c.Transport = httpx.RewriteHostTransport(testx.Must(url.ParseRequestURI(srv.URL))(t), c.Transport)

		cmd := cmdCommunityInfo{Name: "nonexistent-community"}
		err := cmd.run(ctx, c, strings.NewReader(""), &bytes.Buffer{})
		require.Error(t, err)
	})

	t.Run("includes community name in path", func(t *testing.T) {
		ctx, cancel := testx.Context(t)
		defer cancel()

		communityName := "my-special-community"

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/c/"+communityName, r.URL.Path)
			assert.NoError(t, json.NewEncoder(w).Encode(&communityapi.CommunityFindResponse{
				Community: &communityapi.Community{},
			}))
		}))
		defer srv.Close()

		c := &http.Client{}
		c.Transport = httpx.RewriteHostTransport(testx.Must(url.ParseRequestURI(srv.URL))(t), c.Transport)

		cmd := cmdCommunityInfo{Name: communityName}
		err := cmd.run(ctx, c, strings.NewReader(""), &bytes.Buffer{})
		require.NoError(t, err)
	})

	t.Run("copies stdin to stdout", func(t *testing.T) {
		ctx, cancel := testx.Context(t)
		defer cancel()

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.NoError(t, json.NewEncoder(w).Encode(&communityapi.CommunityFindResponse{
				Community: &communityapi.Community{},
			}))
		}))
		defer srv.Close()

		c := &http.Client{}
		c.Transport = httpx.RewriteHostTransport(testx.Must(url.ParseRequestURI(srv.URL))(t), c.Transport)

		const stdinContent = "hello from stdin"
		var out bytes.Buffer

		cmd := cmdCommunityInfo{Name: "test-community"}
		require.NoError(t, cmd.run(ctx, c, strings.NewReader(stdinContent), &out))
		require.Contains(t, out.String(), stdinContent)
	})
}
