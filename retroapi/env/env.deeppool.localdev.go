//go:build localdev

package env

import (
	"github.com/retrovibed/retrovibed/retroapi/internal/envx"
)

func Deeppool() string {
	return envx.String("localhost:8081", DeeppoolEndpoint)
}
