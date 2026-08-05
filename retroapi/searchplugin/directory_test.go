package searchplugin

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSanitizeName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "empty", input: "", expected: ""},
		{name: "already safe", input: "my-plugin_1", expected: "my-plugin_1"},
		{name: "path traversal", input: "../../etc/passwd", expected: "etcpasswd"},
		{name: "absolute path", input: "/etc/passwd", expected: "etcpasswd"},
		{name: "path separator", input: "foo/bar", expected: "foobar"},
		{name: "windows path separator", input: "foo\\bar", expected: "foobar"},
		{name: "dot only", input: "..", expected: ""},
		{name: "current dir", input: ".", expected: ""},
		{name: "extension stripped", input: "plugin.wasm", expected: "pluginwasm"},
		{name: "whitespace stripped", input: "my plugin", expected: "myplugin"},
		{name: "null byte stripped", input: "plugin\x00name", expected: "pluginname"},
		{name: "unicode stripped", input: "plugin日本語", expected: "plugin"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected, SanitizeName(tt.input))
		})
	}
}
