package cmdmedia

import (
	"io"
	"os"

	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/internal/jsonl"
)

// examples:
// retrovibed media known genquery "inception" "the dark knight" | retrovibed media known query
// retrovibed media known genquery "inception" "the dark knight" | retrovibed media known detect
type knowngenquery struct {
	Out     cmdopts.IOOut `flag:"" name:"out" default:"-" help:"output destination; '-' for stdout"`
	Queries []string      `arg:"" name:"query" help:"queries to generate"`
}

func (t knowngenquery) Run(gctx *cmdopts.Global) error {
	out, err := t.Out.Open(os.Stdout)
	if err != nil {
		return err
	}
	defer out.Close()

	return t.run(out)
}

func (t knowngenquery) run(out io.Writer) error {
	type output struct {
		Query string `json:"query"`
	}

	enc := jsonl.NewEncoder(out)
	for _, q := range t.Queries {
		if err := enc.Encode(output{Query: q}); err != nil {
			return err
		}
	}

	return nil
}
