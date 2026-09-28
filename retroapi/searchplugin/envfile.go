package searchplugin

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

// readEnvFile parses a .env-style file (KEY=VALUE per line; blank lines
// and lines beginning with # are skipped) into a "KEY=VALUE" pair slice.
// A missing file is not an error - it yields no pairs. A KEY= with a blank
// value is treated as unset and skipped: the settings form saves every
// declared variable, and an empty env var would otherwise override the
// plugin's own default (including values baked in via -ldflags -X).
func readEnvFile(path string) (pairs []string, err error) {
	f, err := os.Open(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if _, v, ok := strings.Cut(line, "="); ok && strings.TrimSpace(v) != "" {
			pairs = append(pairs, line)
		}
	}

	return pairs, scanner.Err()
}
