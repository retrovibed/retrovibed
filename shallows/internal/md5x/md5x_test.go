package md5x_test

import (
	"encoding/hex"
	"testing"

	"github.com/retrovibed/retrovibed/shallows/internal/md5x"
	"github.com/stretchr/testify/require"
)

func TestJSON(t *testing.T) {
	type record struct {
		Name string
		Tags []string
	}

	// JSON content addresses records that are already stored, so its digest is
	// a golden value: it has to stay byte for byte what encoding/json produces.
	// json/v2 (jsonx) drops the html escaping and writes a nil slice as [],
	// which yields cb5100deea9d36fd5c18152104de1f90 instead and would silently
	// rekey everything.
	t.Run("digests the encoding/json encoding of the value", func(t *testing.T) {
		in := record{Name: "Tom & Jerry <1940>"}

		require.Equal(t, "3e55f4f5e4b35380d5bdca2c25c94eee", hex.EncodeToString(md5x.JSON(in).Sum(nil)))
	})

	t.Run("digests the same value the same way through a pointer", func(t *testing.T) {
		in := record{Name: "Tom & Jerry <1940>"}

		require.Equal(t, hex.EncodeToString(md5x.JSON(in).Sum(nil)), hex.EncodeToString(md5x.JSON(&in).Sum(nil)))
	})

	t.Run("distinguishes an empty slice from a nil one", func(t *testing.T) {
		nilled := md5x.JSON(record{Name: "derp"})
		empty := md5x.JSON(record{Name: "derp", Tags: []string{}})

		require.NotEqual(t, hex.EncodeToString(nilled.Sum(nil)), hex.EncodeToString(empty.Sum(nil)))
	})
}
