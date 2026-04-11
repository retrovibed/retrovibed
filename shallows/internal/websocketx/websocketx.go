package websocketx

import (
	"github.com/coder/websocket"
)

// https://datatracker.ietf.org/doc/html/rfc6455#section-7.4.2
func PrivateStatus(httpcode int) websocket.StatusCode {
	return websocket.StatusCode(4000 + httpcode)
}
