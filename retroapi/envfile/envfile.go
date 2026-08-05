// Package envfile implements the client-side convention for .env sidecar
// files: a comment line (or block of comment lines) immediately preceding
// a KEY=VALUE line, with no blank line in between, becomes that variable's
// help text. A trailing "# ..." on the same line as the value takes
// precedence over any preceding comment block. Servers treat .env files as
// an opaque byte blob; this package is what turns that blob into something
// a UI can render, and back.
package envfile

import (
	"regexp"
	"strings"
)

type Variable struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Hint  string `json:"hint"`
}

var assignment = regexp.MustCompile(`^([A-Za-z_][A-Za-z0-9_]*)=(.*)$`)

// splitValueAndHint splits a KEY=VALUE line's value portion into
// (value, hint), stripping a surrounding quote pair from value and
// treating a trailing "# ..." found past the closing quote (or, if
// unquoted, anywhere) as the hint.
func splitValueAndHint(raw string) (value, hint string) {
	if raw == "" {
		return "", ""
	}

	if quote := raw[0]; quote == '"' || quote == '\'' {
		if close := strings.IndexByte(raw[1:], quote); close != -1 {
			close++
			value = raw[1:close]
			rest := strings.TrimLeft(raw[close+1:], " \t")
			if strings.HasPrefix(rest, "#") {
				hint = strings.TrimSpace(rest[1:])
			}
			return value, hint
		}
	}

	if idx := strings.IndexByte(raw, '#'); idx != -1 {
		return strings.TrimSpace(raw[:idx]), strings.TrimSpace(raw[idx+1:])
	}

	return strings.TrimSpace(raw), ""
}

// Parse scans content line-by-line, attaching the nearest preceding
// comment block (cleared by a blank line) as a variable's hint, unless
// that variable's own line has a trailing comment, which wins.
func Parse(content string) []Variable {
	var (
		variables []Variable
		pending   []string
	)

	for _, line := range strings.Split(content, "\n") {
		trimmed := strings.TrimSpace(line)

		switch {
		case strings.HasPrefix(trimmed, "#"):
			pending = append(pending, strings.TrimSpace(strings.TrimPrefix(trimmed, "#")))
			continue
		case trimmed == "":
			pending = nil
			continue
		}

		m := assignment.FindStringSubmatch(trimmed)
		if m == nil {
			pending = nil
			continue
		}

		value, inlineHint := splitValueAndHint(m[2])
		hint := inlineHint
		if hint == "" {
			hint = strings.Join(pending, " ")
		}

		variables = append(variables, Variable{Key: m[1], Value: value, Hint: hint})
		pending = nil
	}

	return variables
}

func quoteIfNeeded(value string) string {
	if value == "" {
		return value
	}
	if strings.ContainsAny(value, " \t#\"") {
		return `"` + strings.ReplaceAll(value, `"`, `\"`) + `"`
	}
	return value
}

// Apply re-serializes content, replacing the value of every KEY=VALUE line
// whose key matches an edit while leaving that line's comment and every
// other line byte-identical. Edits whose key has no matching line are
// appended as new bare KEY=VALUE lines.
func Apply(content string, edits []Variable) string {
	byKey := make(map[string]Variable, len(edits))
	for _, e := range edits {
		byKey[e.Key] = e
	}

	seen := make(map[string]bool, len(edits))
	lines := strings.Split(content, "\n")

	for i, line := range lines {
		m := assignment.FindStringSubmatch(strings.TrimSpace(line))
		if m == nil {
			continue
		}

		edit, ok := byKey[m[1]]
		if !ok {
			continue
		}

		seen[m[1]] = true
		_, hint := splitValueAndHint(m[2])
		comment := ""
		if hint != "" {
			comment = " # " + hint
		}
		lines[i] = m[1] + "=" + quoteIfNeeded(edit.Value) + comment
	}

	for _, e := range edits {
		if seen[e.Key] {
			continue
		}
		line := e.Key + "=" + quoteIfNeeded(e.Value)
		if n := len(lines); n > 0 && lines[n-1] == "" {
			lines = append(lines[:n-1], line, "")
		} else {
			lines = append(lines, line)
		}
	}

	return strings.Join(lines, "\n")
}
