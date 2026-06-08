package testx

import (
	"math/rand/v2"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/retrovibed/retrovibed/retroapi/internal/langx"
	"github.com/stretchr/testify/require"
)

type faked interface {
	Struct(v any) error
}

func Fake[T any](v *T, options ...func(*T)) error {
	if err := gofakeit.Struct(v); err != nil {
		return err
	}

	*v = langx.Clone(*v, options...)

	return nil
}

func Seeded(src rand.Source) *gofakeit.Faker {
	return gofakeit.NewFaker(src, false)
}

func SeededFake[T any](t testing.TB, n faked, v *T, options ...func(*T)) error {
	if err := n.Struct(v); err != nil {
		require.NoError(t, err)
	}

	*v = langx.Clone(*v, options...)

	return nil
}
