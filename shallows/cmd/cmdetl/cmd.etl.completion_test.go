package cmdetl

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/jsonx"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/internal/jsonl"
	"github.com/stretchr/testify/require"
)

func gctx(t *testing.T) *cmdopts.Global {
	t.Helper()
	ctx, cancel := testx.Context(t)
	t.Cleanup(cancel)
	return &cmdopts.Global{Context: ctx, Shutdown: cancel, Cleanup: &sync.WaitGroup{}}
}

type llamaResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func newLlamaResponse(content string) llamaResponse {
	var r llamaResponse
	r.Choices = append(r.Choices, struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	}{Message: struct {
		Content string `json:"content"`
	}{Content: content}})
	return r
}

func TestCmdCompletion(t *testing.T) {
	t.Run("sends record to server and emits output", func(t *testing.T) {
		var received request
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.NoError(t, jsonx.UnmarshalRead(r.Body, &received))
			require.NoError(t, json.NewEncoder(w).Encode(newLlamaResponse(`{"output":"result"}`)))
		}))
		defer srv.Close()

		var out bytes.Buffer
		cmd := cmdCompletion{URL: srv.URL, Concurrency: 1, Template: "{{.}}"}
		require.NoError(t, cmd.run(gctx(t), strings.NewReader(`{"prompt":"hello"}`+"\n"), &out, io.Discard))

		var result map[string]any
		require.NoError(t, json.Unmarshal(out.Bytes(), &result))
		require.Equal(t, "result", result["output"])
	})

	t.Run("sends one request per record", func(t *testing.T) {
		calls := 0
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			calls++
			require.NoError(t, json.NewEncoder(w).Encode(newLlamaResponse(`{}`)))
		}))
		defer srv.Close()

		input := `{"prompt":"a"}` + "\n" + `{"prompt":"b"}` + "\n" + `{"prompt":"c"}` + "\n"
		cmd := cmdCompletion{URL: srv.URL, Concurrency: 1, Template: "{{.}}"}
		require.NoError(t, cmd.run(gctx(t), strings.NewReader(input), &bytes.Buffer{}, io.Discard))
		require.Equal(t, 3, calls)
	})

	t.Run("template is applied to record", func(t *testing.T) {
		var receivedPrompt string
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var req request
			require.NoError(t, jsonx.UnmarshalRead(r.Body, &req))
			for _, m := range req.Messages {
				if m.Role == "user" {
					receivedPrompt = m.Content
				}
			}
			require.NoError(t, json.NewEncoder(w).Encode(newLlamaResponse(`{}`)))
		}))
		defer srv.Close()

		cmd := cmdCompletion{URL: srv.URL, Concurrency: 1, Template: `classify: {{.}}`}
		require.NoError(t, cmd.run(gctx(t), strings.NewReader(`{"prompt":"hello"}`+"\n"), &bytes.Buffer{}, io.Discard))
		require.Equal(t, `classify: {"prompt":"hello"}`, receivedPrompt)
	})

	t.Run("invalid template returns error", func(t *testing.T) {
		cmd := cmdCompletion{URL: "http://localhost:1", Template: "{{.Unclosed"}
		err := cmd.run(gctx(t), strings.NewReader(`{"prompt":"hello"}`+"\n"), &bytes.Buffer{}, io.Discard)
		require.Error(t, err)
	})

	t.Run("server error emits failure record to errw", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer srv.Close()

		var out, errw bytes.Buffer
		cmd := cmdCompletion{URL: srv.URL, Concurrency: 1, Template: "{{.}}"}
		require.NoError(t, cmd.run(gctx(t), strings.NewReader(`{"prompt":"hello"}`+"\n"), &out, &errw))

		require.Empty(t, out.Bytes())
		var result map[string]any
		require.NoError(t, json.Unmarshal(errw.Bytes(), &result))
		require.Equal(t, true, result["failed"])
		require.NotEmpty(t, result["cause"])
	})

	t.Run("empty input produces no output and no requests", func(t *testing.T) {
		calls := 0
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			calls++
		}))
		defer srv.Close()

		var out bytes.Buffer
		cmd := cmdCompletion{URL: srv.URL, Concurrency: 1, Template: "{{.}}"}
		require.NoError(t, cmd.run(gctx(t), strings.NewReader(""), &out, io.Discard))
		require.Empty(t, out.Bytes())
		require.Equal(t, 0, calls)
	})

	t.Run("output is valid jsonl", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.NoError(t, json.NewEncoder(w).Encode(newLlamaResponse(`{"value":"out"}`)))
		}))
		defer srv.Close()

		input := `{"prompt":"a"}` + "\n" + `{"prompt":"b"}` + "\n"
		var out bytes.Buffer
		cmd := cmdCompletion{URL: srv.URL, Concurrency: 1, Template: "{{.}}"}
		require.NoError(t, cmd.run(gctx(t), strings.NewReader(input), &out, io.Discard))

		dec := jsonl.NewDecoder(bytes.NewReader(out.Bytes()))
		var count int
		for {
			var result map[string]any
			if err := dec.Decode(&result); err != nil {
				break
			}
			require.Equal(t, "out", result["value"])
			count++
		}
		require.Equal(t, 2, count)
	})
}
