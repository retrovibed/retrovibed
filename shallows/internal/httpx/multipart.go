package httpx

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/textproto"
	"strings"

	"github.com/retrovibed/retrovibed/internal/errorsx"
)

func escapeQuotes(s string) string {
	quoteEscaper := strings.NewReplacer("\\", "\\\\", `"`, "\\\"")
	return quoteEscaper.Replace(s)
}

func NewMultipartHeader(mimetype string, fieldname string, filename string) textproto.MIMEHeader {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
			escapeQuotes(fieldname), escapeQuotes(filename)))
	h.Set("Content-Type", mimetype)
	return h
}

func Multipart(do func(*multipart.Writer) error) (contentType string, _ io.ReadCloser, err error) {
	r, w := io.Pipe()

	mw := multipart.NewWriter(w)

	go func() {
		errorsx.Log(w.CloseWithError(errorsx.Compact(do(mw), mw.Close())))
	}()

	return mw.FormDataContentType(), r, nil
}
