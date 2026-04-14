package cmdetl

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/internal/asynccompute"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/jsonl"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
)

// examples:
// cat prompts.jsonl | retrovibed etl completion
// cat prompts.jsonl | retrovibed etl completion --url http://192.168.1.10:11434/api/chat
// cat prompts.jsonl | retrovibed etl completion --template prompt.tmpl
// cat prompts.jsonl | retrovibed etl completion --template 'classify this record: {{.}}'
// duckdb example.sqlite -c "COPY (SELECT * FROM table) TO '/dev/stdout' (FORMAT JSON, ARRAY false);" | retrovibed etl completion --template 'clean this title: {{.}}'
type cmdCompletion struct {
	URL           string               `flag:"" name:"url" default:"http://localhost:11434/api/chat" help:"server endpoint"`
	Timeout       time.Duration        `flag:"" name:"timeout" default:"240s" help:"http client timeout per request"`
	Thinking      bool                 `flag:"" name:"thinking" help:"enable/disable thinking" default:"false" negatable:""`
	Concurrency   uint16               `flag:"" name:"concurrency" default:"${vars_cores}" help:"number of concurrent requests"`
	RepeatPenalty float64              `flag:"" name:"repeat-penalty" default:"1.2" help:"repeat penalty for the model"`
	Temperature   float64              `flag:"" name:"temperature" default:"0.0" help:"temperature for the model"`
	Template      cmdopts.FileContents `flag:"" name:"template" default:"{{.}}" help:"text/template used to build the prompt; receives the record as '.'"`
	Out           cmdopts.IOOut        `flag:"" name:"out" default:"-" help:"output destination; '-' for stdout"`
	Err           cmdopts.IOOut        `flag:"" name:"err" default:"-" help:"error output destination; '-' for stderr"`
}

func (t cmdCompletion) Run(gctx *cmdopts.Global) error {
	out, err := t.Out.Open(os.Stdout)
	if err != nil {
		return err
	}
	defer out.Close()

	errw, err := t.Err.Open(os.Stderr)
	if err != nil {
		return err
	}
	defer errw.Close()

	return t.run(gctx, os.Stdin, out, errw)
}

func (t cmdCompletion) run(gctx *cmdopts.Global, in io.Reader, out io.Writer, errw io.Writer) error {
	type LLMResponse struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	tmpl, err := template.New("prompt").Parse(string(t.Template))
	if err != nil {
		return errorsx.Wrap(err, "failed to parse template")
	}

	client := &http.Client{Timeout: t.Timeout}
	client = httpx.BindRetryTransport(client, http.StatusInternalServerError)
	if gctx.Verbosity > 1 {
		client = httpx.DebugClient(client)
	}
	enc := jsonl.NewEncoder(out)
	errenc := jsonl.NewEncoder(errw)

	type insertpayload struct {
		err    error
		record json.RawMessage
		data   json.RawMessage
	}

	insert := asynccompute.New(func(ctx context.Context, v insertpayload) error {
		type failed struct {
			Failed   bool            `json:"failed,omitzero"`
			Error    string          `json:"cause"`
			Original json.RawMessage `json:"record"`
		}
		if v.err != nil {
			return errenc.Encode(failed{Failed: true, Error: v.err.Error(), Original: v.record})
		}
		return enc.Encode(v.data)
	}, asynccompute.Workers[insertpayload](1))

	pool := asynccompute.New(func(ctx context.Context, rec json.RawMessage) error {
		var prompt bytes.Buffer
		if err := tmpl.Execute(&prompt, string(rec)); err != nil {
			return errorsx.Wrap(err, "failed to execute template")
		}

		req := request{
			RepeatPenalty: t.RepeatPenalty,
			Temperature:   t.Temperature,
			Stream:        false,
			Think:         t.Thinking,
			Format:        "json",
			Messages: []message{
				{
					Role:    "system",
					Content: "You are a data processor. Output ONLY raw JSONL. No markdown, no thinking, no preamble. You may only second guess your response twice.",
				},
				{Role: "user", Content: prompt.String()},
			},
			ChatTemplateKwargs: map[string]interface{}{
				"enable_thinking": t.Thinking,
			},
		}
		body, err := json.Marshal(req)
		if err != nil {
			return errorsx.Wrap(err, "failed to marshal request")
		}

		httpreq, err := http.NewRequestWithContext(ctx, http.MethodPost, t.URL, bytes.NewBuffer(body))
		if err != nil {
			return errorsx.Wrap(err, "failed to create request")
		}
		httpreq.Header.Set("Content-Type", "application/json")

		resp, err := httpx.AsError(client.Do(httpreq))
		if err != nil {
			return insert.Run(ctx, insertpayload{record: rec, err: errorsx.Wrap(err, "failed to post request")})
		}

		defer resp.Body.Close()
		var decoded LLMResponse
		if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
			return insert.Run(ctx, insertpayload{record: rec, err: errorsx.Wrap(err, "failed to decode response")})
		}
		content := ""
		if len(decoded.Choices) > 0 {
			content = decoded.Choices[0].Message.Content
		}
		final := json.RawMessage(content)

		return insert.Run(ctx, insertpayload{data: final})
	}, asynccompute.Workers[json.RawMessage](t.Concurrency))

	seq := jsonl.Iter[json.RawMessage](jsonl.NewDecoder(in))
	for rec := range seq.Each(gctx.Context) {
		if err := pool.Run(gctx.Context, rec); err != nil {
			return err
		}
	}

	if err := seq.Err(); err != nil {
		return errorsx.Wrap(err, "failed to decode input")
	}

	return langx.FirstNonZero(
		asynccompute.Shutdown(gctx.Context, pool),
		asynccompute.Shutdown(gctx.Context, insert),
	)
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type request struct {
	RepeatPenalty      float64        `json:"repeat_penalty"`
	Temperature        float64        `json:"temperature"`
	Model              string         `json:"model,omitempty"`
	Messages           []message      `json:"messages,omitempty"`
	Stream             bool           `json:"stream"`
	Think              bool           `json:"think"`
	Format             string         `json:"format"`
	ChatTemplateKwargs map[string]any `json:"chat_template_kwargs"`
	ExtraBody          map[string]any `json:"extra_body"`
}
