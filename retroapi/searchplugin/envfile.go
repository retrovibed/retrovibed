package searchplugin

import (
	"errors"
	"os"

	"github.com/retrovibed/retrovibed/retroapi/envfile"
	"github.com/retrovibed/retrovibed/retroapi/internal/stringsx"
)

// readEnvFile parses a .env-style file into a "KEY=VALUE" pair slice using
// the same convention envfile.Apply writes it with: comment lines and
// trailing "# ..." hints are dropped and a surrounding quote pair is
// stripped, so only the value itself reaches the plugin.
// A missing file is not an error - it yields no pairs. A KEY= with a blank
// value is treated as unset and skipped: the settings form saves every
// declared variable, and an empty env var would otherwise override the
// plugin's own default (including values baked in via -ldflags -X).
func readEnvFile(path string) (pairs []string, err error) {
	content, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	for _, v := range envfile.Parse(string(content)) {
		if stringsx.Blank(v.Value) {
			continue
		}

		pairs = append(pairs, v.Key+"="+v.Value)
	}

	return pairs, nil
}
