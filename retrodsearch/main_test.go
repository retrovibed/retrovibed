package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/retrovibed/retrovibed/retroapi/ddiscapi"
	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/stretchr/testify/require"
	"golang.org/x/time/rate"
)

func TestPeerTubeRunEncodesResults(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/search/videos", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "ubuntu", r.URL.Query().Get("search"))
		fmt.Fprint(w, `{"total":1,"data":[{"uuid":"abc-123","name":"Ubuntu Documentary","description":"a documentary about ubuntu","thumbnailPath":"/thumb.jpg","views":42}]}`)
	})
	mux.HandleFunc("GET /api/v1/videos/abc-123", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"files":[{"resolution":{"id":1080},"magnetUri":"magnet:?xt=urn:btih:1111111111111111111111111111111111111111","fileDownloadUrl":""}]}`)
	})
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := cli{
		Query:      "ubuntu",
		Categories: []string{"all"},
		Domain:     srv.URL,
		Attempts:   1,
		Workers:    1,
	}

	var buf bytes.Buffer
	require.NoError(t, cmd.run(ctx, srv.Client(), &buf, rate.NewLimiter(rate.Inf, 1)))

	var imp ddiscapi.Import
	require.NoError(t, json.Unmarshal(buf.Bytes(), &imp))
	require.Equal(t, "magnet:?xt=urn:btih:1111111111111111111111111111111111111111", imp.Uri)
	require.Equal(t, mimex.Magnet, imp.Uritype)
	require.Equal(t, "Ubuntu Documentary", imp.Title)
	require.Equal(t, "a documentary about ubuntu", imp.Overview)
	require.Equal(t, srv.URL+"/thumb.jpg", imp.PosterPath)
	require.Equal(t, float64(42), imp.Popularity)
}

func TestPeerTubeRunSkipsFileWithNoMagnet(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/search/videos", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"total":1,"data":[{"uuid":"abc-123","name":"Ubuntu Documentary"}]}`)
	})
	mux.HandleFunc("GET /api/v1/videos/abc-123", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"files":[{"resolution":{"id":1080},"magnetUri":""}]}`)
	})
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := cli{Query: "ubuntu", Categories: []string{"all"}, Domain: srv.URL, Attempts: 1, Workers: 1}

	var buf bytes.Buffer
	require.NoError(t, cmd.run(ctx, srv.Client(), &buf, rate.NewLimiter(rate.Inf, 1)))
	require.Empty(t, buf.String(), "a file without a magnet link must never fall back to a download link")
}

func TestPeerTubeRunPicksHighestResolutionMagnet(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/search/videos", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"total":1,"data":[{"uuid":"abc-123","name":"Ubuntu Documentary"}]}`)
	})
	mux.HandleFunc("GET /api/v1/videos/abc-123", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"files":[
			{"resolution":{"id":480},"magnetUri":"magnet:?xt=urn:btih:4444444444444444444444444444444444444444"},
			{"resolution":{"id":1080},"magnetUri":"magnet:?xt=urn:btih:5555555555555555555555555555555555555555"},
			{"resolution":{"id":720},"magnetUri":""}
		]}`)
	})
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := cli{Query: "ubuntu", Categories: []string{"all"}, Domain: srv.URL, Attempts: 1, Workers: 1}

	var buf bytes.Buffer
	require.NoError(t, cmd.run(ctx, srv.Client(), &buf, rate.NewLimiter(rate.Inf, 1)))

	var imp ddiscapi.Import
	require.NoError(t, json.Unmarshal(buf.Bytes(), &imp))
	require.Equal(t, "magnet:?xt=urn:btih:5555555555555555555555555555555555555555", imp.Uri)
	require.Equal(t, mimex.Magnet, imp.Uritype)
}

func TestPeerTubeRunSkipsVideoWithNoFiles(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/search/videos", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"total":1,"data":[{"uuid":"abc-123","name":"Ubuntu Documentary"}]}`)
	})
	mux.HandleFunc("GET /api/v1/videos/abc-123", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"files":[]}`)
	})
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := cli{Query: "ubuntu", Categories: []string{"all"}, Domain: srv.URL, Attempts: 1, Workers: 1}

	var buf bytes.Buffer
	require.NoError(t, cmd.run(ctx, srv.Client(), &buf, rate.NewLimiter(rate.Inf, 1)))
	require.Empty(t, buf.String())
}

func TestPeerTubeRunPaginatesUntilTotalReached(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/search/videos", func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Query().Get("start") {
		case "0":
			fmt.Fprint(w, `{"total":2,"data":[{"uuid":"v1","name":"one"}]}`)
		case "1":
			fmt.Fprint(w, `{"total":2,"data":[{"uuid":"v2","name":"two"}]}`)
		default:
			t.Fatalf("unexpected start=%s", r.URL.Query().Get("start"))
		}
	})
	mux.HandleFunc("GET /api/v1/videos/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"files":[{"resolution":{"id":720},"magnetUri":"magnet:?xt=urn:btih:2222222222222222222222222222222222222222"}]}`)
	})
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := cli{Query: "ubuntu", Categories: []string{"all"}, Domain: srv.URL, Attempts: 1, Workers: 1}

	var buf bytes.Buffer
	require.NoError(t, cmd.run(ctx, srv.Client(), &buf, rate.NewLimiter(rate.Inf, 1)))

	var count int
	dec := json.NewDecoder(&buf)
	for dec.More() {
		var imp ddiscapi.Import
		require.NoError(t, dec.Decode(&imp))
		count++
	}
	require.Equal(t, 2, count)
}

func TestPeerTubeRunCapsAtMaxResults(t *testing.T) {
	var page2Fetched bool

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/search/videos", func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Query().Get("start") {
		case "0":
			require.Equal(t, "1", r.URL.Query().Get("count"))
			fmt.Fprint(w, `{"total":2,"data":[{"uuid":"v1","name":"one"}]}`)
		default:
			page2Fetched = true
			fmt.Fprint(w, `{"total":2,"data":[{"uuid":"v2","name":"two"}]}`)
		}
	})
	mux.HandleFunc("GET /api/v1/videos/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"files":[{"resolution":{"id":720},"magnetUri":"magnet:?xt=urn:btih:3333333333333333333333333333333333333333"}]}`)
	})
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := cli{Query: "ubuntu", Categories: []string{"all"}, Domain: srv.URL, Attempts: 1, Workers: 1, MaxResults: 1}

	var buf bytes.Buffer
	require.NoError(t, cmd.run(ctx, srv.Client(), &buf, rate.NewLimiter(rate.Inf, 1)))
	require.False(t, page2Fetched, "must not fetch a second page once MaxResults is satisfied")
}

func TestPeerTubeRunRejectsUnsupportedCategory(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := cli{Query: "ubuntu", Categories: []string{"bogus"}, Domain: "http://example.invalid", Attempts: 1, Workers: 1}

	var buf bytes.Buffer
	err := cmd.run(ctx, http.DefaultClient, &buf, rate.NewLimiter(rate.Inf, 1))
	require.ErrorContains(t, err, "unsupported category")
}

func TestPeerTubeRunForwardsAdultToNSFWParam(t *testing.T) {
	var gotNSFW string

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/search/videos", func(w http.ResponseWriter, r *http.Request) {
		gotNSFW = r.URL.Query().Get("nsfw")
		fmt.Fprint(w, `{"total":0,"data":[]}`)
	})
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := cli{Query: "ubuntu", Categories: []string{"all"}, Domain: srv.URL, Attempts: 1, Workers: 1}
	var buf bytes.Buffer
	require.NoError(t, cmd.run(ctx, srv.Client(), &buf, rate.NewLimiter(rate.Inf, 1)))
	require.Equal(t, "false", gotNSFW)

	cmd.Adult = true
	buf.Reset()
	require.NoError(t, cmd.run(ctx, srv.Client(), &buf, rate.NewLimiter(rate.Inf, 1)))
	require.Equal(t, "true", gotNSFW)
}
