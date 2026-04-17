package testx

import (
	faker "github.com/go-faker/faker/v4"
)

func Fake[T any](v *T, options ...func(*T)) error {
	if err := faker.FakeData(v); err != nil {
		return err
	}

	for _, opt := range options {
		opt(v)
	}

	return nil
}
