package httpx

import (
	"bytes"
	"io"
	"net/http"

	"github.com/retrovibed/retrovibed/retroapi/jsonx"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
)

// EncodeJSON encode data into the http.Request body.
func EncodeJSON(req *http.Request, body interface{}) (err error) {
	var (
		encoded []byte
	)

	if encoded, err = jsonx.Marshal(body); err != nil {
		return errorsx.WithStack(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Body = io.NopCloser(bytes.NewReader(encoded))

	return nil
}
