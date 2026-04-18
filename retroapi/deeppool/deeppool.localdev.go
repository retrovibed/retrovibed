//go:build localdev

package deeppool

import (
	"github.com/retrovibed/retrovibed/retroapi/env"
	"github.com/retrovibed/retrovibed/retroapi/internal/envx"
)

func Deeppool() string {
	return envx.String("localhost:8081", env.DeeppoolEndpoint)
}
