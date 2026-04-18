package cmdopts

import "github.com/retrovibed/retrovibed/shallows/internal/env"

func JWTSecret() []byte {
	return env.JWTSecret()
}
