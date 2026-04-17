//go:build localdev

package deeppool

import (
	"github.com/retrovibed/retroapi/internal/env"
	"github.com/retrovibed/retroapi/internal/envx"
)

func Deeppool() string {
	return envx.String("localhost:8081", env.DeeppoolEndpoint)
	// return envx.String("api.retrovibe.space", env.DeeppoolEndpoint)
}
