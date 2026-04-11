package library_test

import (
	"math"
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retrovibed/library"
	"github.com/stretchr/testify/require"
)

func TestKnownImportedUintID(t *testing.T) {
	t.Run("zero ID and simple prefix", func(t *testing.T) {
		require.Equal(t, "61626381-0000-0000-0000-000000000000", library.KnownImportedUintID("abc", 0))
	})

	t.Run("excessively long prefix", func(t *testing.T) {
		require.Equal(t, "6c6f6e67-0000-0000-0000-000000000000", library.KnownImportedUintID("long-prefix-for-id", 0))
	})

	t.Run("single digit ID", func(t *testing.T) {
		require.Equal(t, "746d6462-0000-0000-0000-000000000007", library.KnownImportedUintID("tmdb", 7))
	})

	t.Run("maximum uint64", func(t *testing.T) {
		require.Equal(t, "811c9dc5-0000-0000-00ff-ffffffffffff", library.KnownImportedUintID("", math.MaxUint64))
	})

	t.Run("highest byte all 1", func(t *testing.T) {
		require.Equal(t, "811c9dc5-0000-0000-00ff-000000000000", library.KnownImportedUintID("", 0xFFFF000000000000))
	})

	t.Run("empty prefix", func(t *testing.T) {
		require.Equal(t, "811c9dc5-0000-0000-0000-000000000000", library.KnownImportedUintID("", 0))
	})

	t.Run("prefix with special characters", func(t *testing.T) {
		require.Equal(t, "41214023-0000-0000-0000-000000000001", library.KnownImportedUintID("A!@#$", 1))
	})
}

func TestKnownImportedUUID(t *testing.T) {
	mustParse := func(s string) uuid.UUID {
		u, err := uuid.FromString(s)
		require.NoError(t, err)
		return u
	}

	t.Run("retrovibed prefix with v7 uuid", func(t *testing.T) {
		input := mustParse("01234567-89ab-7def-0123-456789abcdef")
		result := library.KnownImportedUUID("retrovibed", input)
		require.Equal(t, "72657472-89ab-7def-0123-456789abcdef", result.String())
	})

	t.Run("preserves uuid bytes after first 4", func(t *testing.T) {
		input := mustParse("ffffffff-1234-5678-9abc-def012345678")
		result := library.KnownImportedUUID("test", input)
		require.Equal(t, "74657374-1234-5678-9abc-def012345678", result.String())
	})

	t.Run("empty prefix", func(t *testing.T) {
		input := mustParse("00000000-0000-0000-0000-000000000000")
		result := library.KnownImportedUUID("", input)
		require.Equal(t, "811c9dc5-0000-0000-0000-000000000000", result.String())
	})

	t.Run("same prefix produces same first 4 bytes", func(t *testing.T) {
		id1 := mustParse("11111111-1111-1111-1111-111111111111")
		id2 := mustParse("22222222-2222-2222-2222-222222222222")
		result1 := library.KnownImportedUUID("retrovibed", id1)
		result2 := library.KnownImportedUUID("retrovibed", id2)
		require.Equal(t, result1.String()[:8], result2.String()[:8])
		require.NotEqual(t, result1.String()[9:], result2.String()[9:])
	})

	t.Run("does not mutate original uuid", func(t *testing.T) {
		original := mustParse("01234567-89ab-cdef-0123-456789abcdef")
		originalCopy := original
		_ = library.KnownImportedUUID("retrovibed", original)
		require.Equal(t, originalCopy, original)
	})
}
