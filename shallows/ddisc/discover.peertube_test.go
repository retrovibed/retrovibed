package ddisc_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/stretchr/testify/require"
)

// peerTubeTOFUServer serves a single search result whose video detail
// carries two thumbnails of different sizes, for subtests asserting the
// known-media TOFU-record path picks the largest one.
func peerTubeTOFUServer(t *testing.T) *httptest.Server {
	t.Helper()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/search/videos", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"total":1,"data":[{"uuid":"abc-123","name":"Ubuntu Documentary","description":"a documentary about ubuntu"}]}`)
	})
	mux.HandleFunc("GET /api/v1/videos/abc-123", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{
			"files":[{"resolution":{"id":1080},"magnetUri":"magnet:?xt=urn:btih:7777777777777777777777777777777777777777"}],
			"thumbnails":[
				{"height":157,"width":280,"fileUrl":"https://video.example/small.jpg"},
				{"height":480,"width":850,"fileUrl":"https://video.example/large.jpg"}
			]
		}`)
	})
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	return srv
}

func TestPeerTubeStrategy(t *testing.T) {
	t.Run("yields the real infohash parsed from the magnet uri", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc("GET /api/v1/search/videos", func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "ubuntu", r.URL.Query().Get("search"))
			_, _ = fmt.Fprint(w, `{"total":1,"data":[{"uuid":"abc-123","name":"Ubuntu Documentary","description":"a documentary about ubuntu"}]}`)
		})
		mux.HandleFunc("GET /api/v1/videos/abc-123", func(w http.ResponseWriter, r *http.Request) {
			_, _ = fmt.Fprint(w, `{"files":[{"resolution":{"id":1080},"magnetUri":"magnet:?xt=urn:btih:1111111111111111111111111111111111111111"}]}`)
		})
		srv := httptest.NewServer(mux)
		t.Cleanup(srv.Close)

		strategy := ddisc.PeerTubeStrategy(srv.Client(), srv.URL, nil, ddisc.PeerTubeOptionAttempts(1), ddisc.PeerTubeOptionWorkers(1))
		seq := strategy.Discover(t.Context(), ddisc.DiscoverRequest{Query: "ubuntu", Mimetypes: []string{mimex.RetrovibedDiscoveryMovies}})

		var got []ddisc.Discovered
		for d := range seq.Each(t.Context()) {
			got = append(got, d)
		}
		require.NoError(t, seq.Err())
		require.Len(t, got, 1)
		require.Equal(t, "magnet:?xt=urn:btih:1111111111111111111111111111111111111111", got[0].URI)
		require.Equal(t, "Ubuntu Documentary", got[0].Title)
		require.Equal(t, "a documentary about ubuntu", got[0].Description)
		require.Equal(t, mimex.Video, got[0].Category)
		require.NotEqual(t, make([]byte, 20), got[0].Infohash, "infohash must be parsed from the real magnet, not left zeroed")
	})

	t.Run("noops without a query", func(t *testing.T) {
		strategy := ddisc.PeerTubeStrategy(http.DefaultClient, "http://example.invalid", nil, ddisc.PeerTubeOptionAttempts(1), ddisc.PeerTubeOptionWorkers(1))
		seq := strategy.Discover(t.Context(), ddisc.DiscoverRequest{})

		var count int
		for range seq.Each(t.Context()) {
			count++
		}
		require.NoError(t, seq.Err())
		require.Equal(t, 0, count)
	})

	t.Run("skips a file with no magnet link instead of falling back to a download link", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc("GET /api/v1/search/videos", func(w http.ResponseWriter, r *http.Request) {
			_, _ = fmt.Fprint(w, `{"total":1,"data":[{"uuid":"abc-123","name":"Ubuntu Documentary"}]}`)
		})
		mux.HandleFunc("GET /api/v1/videos/abc-123", func(w http.ResponseWriter, r *http.Request) {
			_, _ = fmt.Fprint(w, `{"files":[{"resolution":{"id":1080},"magnetUri":""}]}`)
		})
		srv := httptest.NewServer(mux)
		t.Cleanup(srv.Close)

		strategy := ddisc.PeerTubeStrategy(srv.Client(), srv.URL, nil, ddisc.PeerTubeOptionAttempts(1), ddisc.PeerTubeOptionWorkers(1))
		seq := strategy.Discover(t.Context(), ddisc.DiscoverRequest{Query: "ubuntu"})

		var count int
		for range seq.Each(t.Context()) {
			count++
		}
		require.NoError(t, seq.Err())
		require.Equal(t, 0, count)
	})

	t.Run("picks the highest resolution file that has a magnet link", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc("GET /api/v1/search/videos", func(w http.ResponseWriter, r *http.Request) {
			_, _ = fmt.Fprint(w, `{"total":1,"data":[{"uuid":"abc-123","name":"Ubuntu Documentary"}]}`)
		})
		mux.HandleFunc("GET /api/v1/videos/abc-123", func(w http.ResponseWriter, r *http.Request) {
			_, _ = fmt.Fprint(w, `{"files":[
				{"resolution":{"id":480},"magnetUri":"magnet:?xt=urn:btih:4444444444444444444444444444444444444444"},
				{"resolution":{"id":1080},"magnetUri":"magnet:?xt=urn:btih:5555555555555555555555555555555555555555"},
				{"resolution":{"id":720},"magnetUri":""}
			]}`)
		})
		srv := httptest.NewServer(mux)
		t.Cleanup(srv.Close)

		strategy := ddisc.PeerTubeStrategy(srv.Client(), srv.URL, nil, ddisc.PeerTubeOptionAttempts(1), ddisc.PeerTubeOptionWorkers(1))
		seq := strategy.Discover(t.Context(), ddisc.DiscoverRequest{Query: "ubuntu"})

		var got []ddisc.Discovered
		for d := range seq.Each(t.Context()) {
			got = append(got, d)
		}
		require.NoError(t, seq.Err())
		require.Len(t, got, 1)
		require.Equal(t, "magnet:?xt=urn:btih:5555555555555555555555555555555555555555", got[0].URI)
	})

	t.Run("skips a video with no files", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc("GET /api/v1/search/videos", func(w http.ResponseWriter, r *http.Request) {
			_, _ = fmt.Fprint(w, `{"total":1,"data":[{"uuid":"abc-123","name":"Ubuntu Documentary"}]}`)
		})
		mux.HandleFunc("GET /api/v1/videos/abc-123", func(w http.ResponseWriter, r *http.Request) {
			_, _ = fmt.Fprint(w, `{"files":[]}`)
		})
		srv := httptest.NewServer(mux)
		t.Cleanup(srv.Close)

		strategy := ddisc.PeerTubeStrategy(srv.Client(), srv.URL, nil, ddisc.PeerTubeOptionAttempts(1), ddisc.PeerTubeOptionWorkers(1))
		seq := strategy.Discover(t.Context(), ddisc.DiscoverRequest{Query: "ubuntu"})

		var count int
		for range seq.Each(t.Context()) {
			count++
		}
		require.NoError(t, seq.Err())
		require.Equal(t, 0, count)
	})

	t.Run("paginates the search listing until every result has been seen", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc("GET /api/v1/search/videos", func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Query().Get("start") {
			case "0":
				_, _ = fmt.Fprint(w, `{"total":2,"data":[{"uuid":"v1","name":"one"}]}`)
			case "1":
				_, _ = fmt.Fprint(w, `{"total":2,"data":[{"uuid":"v2","name":"two"}]}`)
			default:
				t.Fatalf("unexpected start=%s", r.URL.Query().Get("start"))
			}
		})
		mux.HandleFunc("GET /api/v1/videos/", func(w http.ResponseWriter, r *http.Request) {
			_, _ = fmt.Fprint(w, `{"files":[{"resolution":{"id":720},"magnetUri":"magnet:?xt=urn:btih:2222222222222222222222222222222222222222"}]}`)
		})
		srv := httptest.NewServer(mux)
		t.Cleanup(srv.Close)

		strategy := ddisc.PeerTubeStrategy(srv.Client(), srv.URL, nil, ddisc.PeerTubeOptionAttempts(1), ddisc.PeerTubeOptionWorkers(1))
		seq := strategy.Discover(t.Context(), ddisc.DiscoverRequest{Query: "ubuntu"})

		var count int
		for range seq.Each(t.Context()) {
			count++
		}
		require.NoError(t, seq.Err())
		require.Equal(t, 2, count)
	})

	t.Run("stops paginating once MaxResults is satisfied", func(t *testing.T) {
		var page2Fetched bool

		mux := http.NewServeMux()
		mux.HandleFunc("GET /api/v1/search/videos", func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Query().Get("start") {
			case "0":
				require.Equal(t, "1", r.URL.Query().Get("count"))
				_, _ = fmt.Fprint(w, `{"total":2,"data":[{"uuid":"v1","name":"one"}]}`)
			default:
				page2Fetched = true
				_, _ = fmt.Fprint(w, `{"total":2,"data":[{"uuid":"v2","name":"two"}]}`)
			}
		})
		mux.HandleFunc("GET /api/v1/videos/", func(w http.ResponseWriter, r *http.Request) {
			_, _ = fmt.Fprint(w, `{"files":[{"resolution":{"id":720},"magnetUri":"magnet:?xt=urn:btih:3333333333333333333333333333333333333333"}]}`)
		})
		srv := httptest.NewServer(mux)
		t.Cleanup(srv.Close)

		strategy := ddisc.PeerTubeStrategy(srv.Client(), srv.URL, nil, ddisc.PeerTubeOptionAttempts(1), ddisc.PeerTubeOptionWorkers(1), ddisc.PeerTubeOptionMaxResults(1))
		seq := strategy.Discover(t.Context(), ddisc.DiscoverRequest{Query: "ubuntu"})

		for range seq.Each(t.Context()) {
		}
		require.NoError(t, seq.Err())
		require.False(t, page2Fetched, "must not fetch a second page once MaxResults is satisfied")
	})

	t.Run("fetches video detail from the row's own origin instance, not the index domain", func(t *testing.T) {
		origin := http.NewServeMux()
		var originHit bool
		origin.HandleFunc("GET /api/v1/videos/abc-123", func(w http.ResponseWriter, r *http.Request) {
			originHit = true
			_, _ = fmt.Fprint(w, `{"files":[{"resolution":{"id":1080},"magnetUri":"magnet:?xt=urn:btih:6666666666666666666666666666666666666666"}]}`)
		})
		originSrv := httptest.NewServer(origin)
		t.Cleanup(originSrv.Close)

		index := http.NewServeMux()
		index.HandleFunc("GET /api/v1/search/videos", func(w http.ResponseWriter, r *http.Request) {
			_, _ = fmt.Fprintf(w, `{"total":1,"data":[{"uuid":"abc-123","name":"Ubuntu Documentary","url":"%s/videos/watch/abc-123"}]}`, originSrv.URL)
		})
		index.HandleFunc("GET /api/v1/videos/abc-123", func(w http.ResponseWriter, r *http.Request) {
			t.Fatal("must not fetch video detail from the index domain when the row carries its own origin url")
		})
		indexSrv := httptest.NewServer(index)
		t.Cleanup(indexSrv.Close)

		strategy := ddisc.PeerTubeStrategy(indexSrv.Client(), indexSrv.URL, nil, ddisc.PeerTubeOptionAttempts(1), ddisc.PeerTubeOptionWorkers(1))
		seq := strategy.Discover(t.Context(), ddisc.DiscoverRequest{Query: "ubuntu"})

		var got []ddisc.Discovered
		for d := range seq.Each(t.Context()) {
			got = append(got, d)
		}
		require.NoError(t, seq.Err())
		require.Len(t, got, 1)
		require.Equal(t, "magnet:?xt=urn:btih:6666666666666666666666666666666666666666", got[0].URI)
		require.True(t, originHit, "must fetch video detail from the row's own origin instance")
	})

	t.Run("forwards adult to the nsfw search param", func(t *testing.T) {
		var gotNSFW string

		mux := http.NewServeMux()
		mux.HandleFunc("GET /api/v1/search/videos", func(w http.ResponseWriter, r *http.Request) {
			gotNSFW = r.URL.Query().Get("nsfw")
			_, _ = fmt.Fprint(w, `{"total":0,"data":[]}`)
		})
		srv := httptest.NewServer(mux)
		t.Cleanup(srv.Close)

		strategy := ddisc.PeerTubeStrategy(srv.Client(), srv.URL, nil, ddisc.PeerTubeOptionAttempts(1), ddisc.PeerTubeOptionWorkers(1))

		for range strategy.Discover(t.Context(), ddisc.DiscoverRequest{Query: "ubuntu"}).Each(t.Context()) {
		}
		require.Equal(t, "false", gotNSFW)

		for range strategy.Discover(t.Context(), ddisc.DiscoverRequest{Query: "ubuntu", Adult: true}).Each(t.Context()) {
		}
		require.Equal(t, "true", gotNSFW)
	})

	// mirrors TestPluginStrategyRecordsKnownMediaTOFU: when the request
	// targets a specific catalog entry, PeerTube's largest thumbnail should
	// get TOFU-recorded onto it via the same ddiscapi.Import path plugins
	// use.
	t.Run("TOFU-records the largest thumbnail when the request targets a known-media id", func(t *testing.T) {
		q := sqltestx.Metadatabase(t)
		srv := peerTubeTOFUServer(t)
		kid := uuid.Must(uuid.NewV4()).String()

		strategy := ddisc.PeerTubeStrategy(srv.Client(), srv.URL, q, ddisc.PeerTubeOptionAttempts(1), ddisc.PeerTubeOptionWorkers(1))
		seq := strategy.Discover(t.Context(), ddisc.DiscoverRequest{KnownMediaID: kid, Query: "ubuntu"})

		var got []ddisc.Discovered
		for d := range seq.Each(t.Context()) {
			got = append(got, d)
		}
		require.NoError(t, seq.Err())
		require.Len(t, got, 1)
		require.Equal(t, kid, got[0].KnownMediaID)

		var known library.Known
		require.NoError(t, library.KnownFindByID(t.Context(), q, kid).Scan(&known))
		require.Equal(t, "Ubuntu Documentary", known.Title)
		require.Equal(t, "a documentary about ubuntu", known.Overview)
		require.Equal(t, "https://video.example/large.jpg", known.PosterPath, "must pick the largest (by pixel area) thumbnail")
		require.Equal(t, "retrovibed.discovery.peertube", known.Source)
	})

	// mirrors TestPluginStrategySkipsSentinelKnownMediaID: a bare free-text
	// search (no target catalog entry) must never create a known-media row.
	t.Run("skips the known-media TOFU record for a sentinel known-media id", func(t *testing.T) {
		for _, kid := range []string{"", uuid.Nil.String(), uuid.Max.String()} {
			t.Run(fmt.Sprintf("known_media_id=%q", kid), func(t *testing.T) {
				q := sqltestx.Metadatabase(t)
				srv := peerTubeTOFUServer(t)

				strategy := ddisc.PeerTubeStrategy(srv.Client(), srv.URL, q, ddisc.PeerTubeOptionAttempts(1), ddisc.PeerTubeOptionWorkers(1))
				seq := strategy.Discover(t.Context(), ddisc.DiscoverRequest{KnownMediaID: kid, Query: "ubuntu"})
				for range seq.Each(t.Context()) {
				}
				require.NoError(t, seq.Err())
				require.Equal(t, 0, sqltestx.Count(t, q, "SELECT COUNT(*) FROM library_known_media"))
			})
		}
	})
}
