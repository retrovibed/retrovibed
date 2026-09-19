package tracking

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMetadataOptionDescriptionDefault(t *testing.T) {
	t.Run("should set the description when blank", func(t *testing.T) {
		md := Metadata{}
		MetadataOptionDescriptionDefault("fallback")(&md)
		require.Equal(t, "fallback", md.Description)
	})

	t.Run("should not override an existing description", func(t *testing.T) {
		md := Metadata{Description: "archlinux-2025.07.01-x86_64.iso"}
		MetadataOptionDescriptionDefault("2025.07.01")(&md)
		require.Equal(t, "archlinux-2025.07.01-x86_64.iso", md.Description)
	})

	t.Run("should leave the description blank when the default is blank", func(t *testing.T) {
		md := Metadata{}
		MetadataOptionDescriptionDefault("")(&md)
		require.Equal(t, "", md.Description)
	})
}
