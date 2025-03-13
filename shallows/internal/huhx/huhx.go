package huhx

import (
	"log"
	"strconv"
	"strings"

	"github.com/charmbracelet/huh"
)

// parse a boolean value from the input. returning any error that occurs.
func Bool(i *huh.Input) (bool, error) {
	return resolve(i, func(s string) (bool, error) {
		switch strings.ToLower(s) {
		case "y":
			return true, nil
		case "n":
			return false, nil
		default:
			return strconv.ParseBool(s)
		}
	})
}

func Fallback[T any](fallback T) func(T, error) T {
	return func(t T, err error) T {
		if err == nil {
			return t
		}

		log.Println(err)
		return fallback
	}
}

func resolve[T any](input *huh.Input, parse func(string) (T, error)) (_zero T, _ error) {
	var v string
	if err := input.Value(&v).Run(); err != nil {
		return _zero, err
	}

	return parse(v)
}
