package httpx

import (
	"net/http"

	"github.com/retrovibed/retrovibed/retroapi/jsonx"
)

// DecodeJSON from a http.Response into the provide destination.
func DecodeJSON(resp *http.Response, dst interface{}) error {
	defer resp.Body.Close()
	return jsonx.UnmarshalRead(resp.Body, dst)
}
