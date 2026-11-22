package konggdx

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"syscall"
)

var pidsocket = regexp.MustCompile(`^gdx\.(\d+)\.socket$`)

// resolveSocket picks the gdx socket to dial.
//
// hint is the configured socket path. if it exists it is used as is. otherwise the hint's directory
// is searched for per-process sockets (gdx.<pid>.socket), skipping our own pid (self) and the sockets of
// processes that no longer exist. a single candidate is used directly, multiple candidates are
// prompted for on out/in when interactive, otherwise an error listing them is returned.
func resolveSocket(hint string, self int, interactive bool, in io.Reader, out io.Writer) (string, error) {
	if _, err := os.Stat(hint); err == nil {
		return hint, nil
	}

	dir := filepath.Dir(hint)
	matches, err := filepath.Glob(filepath.Join(dir, "gdx.*.socket"))
	if err != nil {
		return "", fmt.Errorf("unable to search for gdx sockets: %w", err)
	}

	candidates := make([]string, 0, len(matches))
	for _, m := range matches {
		sub := pidsocket.FindStringSubmatch(filepath.Base(m))
		if sub == nil {
			continue
		}
		pid, err := strconv.Atoi(sub[1])
		if err != nil || pid == self {
			continue
		}
		if err := syscall.Kill(pid, 0); errors.Is(err, syscall.ESRCH) {
			continue
		}
		candidates = append(candidates, m)
	}
	sort.Strings(candidates)

	switch len(candidates) {
	case 0:
		return "", fmt.Errorf("no gdx sockets found in %s", dir)
	case 1:
		return candidates[0], nil
	}

	if !interactive {
		return "", fmt.Errorf("multiple gdx sockets found, select one with --socket:\n  %s", strings.Join(candidates, "\n  "))
	}

	fmt.Fprintln(out, "multiple gdx sockets found:")
	for i, c := range candidates {
		fmt.Fprintf(out, "  %d) %s\n", i+1, c)
	}

	r := bufio.NewReader(in)
	for {
		fmt.Fprintf(out, "select a socket [1-%d]: ", len(candidates))
		line, err := r.ReadString('\n')
		if n, perr := strconv.Atoi(strings.TrimSpace(line)); perr == nil && n >= 1 && n <= len(candidates) {
			return candidates[n-1], nil
		}
		if err != nil {
			return "", fmt.Errorf("no gdx socket selected: %w", err)
		}
	}
}

// autoResolveSocket resolves hint using the current process's pid and terminal.
func autoResolveSocket(hint string) (string, error) {
	interactive := false
	if st, err := os.Stdin.Stat(); err == nil {
		interactive = st.Mode()&os.ModeCharDevice != 0
	}
	return resolveSocket(hint, os.Getpid(), interactive, os.Stdin, os.Stderr)
}
