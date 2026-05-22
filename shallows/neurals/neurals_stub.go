//go:build !(retrovibed && neural)

package neurals

func predict(_ *Text, input string) (string, error) {
	return input, nil
}
