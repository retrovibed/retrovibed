package publishplugin

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

// readEnvFile parses a .env-style file (KEY=VALUE per line; blank lines
// and lines beginning with # are skipped) into a "KEY=VALUE" pair slice.
// A missing file is not an error - it yields no pairs.
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

		if _, _, ok := strings.Cut(line, "="); ok {
			pairs = append(pairs, line)
		}
	}

	return pairs, scanner.Err()
}
