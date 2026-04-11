//go:build localdev

package deeppool

import (
	"github.com/retrovibed/retrovibed/internal/env"
	"github.com/retrovibed/retrovibed/internal/envx"
)

func Deeppool() string {
	return envx.String("localhost:8081", env.DeeppoolEndpoint)
	// return envx.String("api.retrovibe.space", env.DeeppoolEndpoint)
}
