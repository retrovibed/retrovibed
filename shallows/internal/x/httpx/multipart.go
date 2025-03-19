package httpx

import (
	"fmt"
	"mime/multipart"
	"net/textproto"
	"os"
	"strings"

	"github.com/retrovibed/retrovibed/internal/x/envx"
	"github.com/retrovibed/retrovibed/internal/x/errorsx"
	"github.com/retrovibed/retrovibed/internal/x/iox"
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

func Multipart(do func(*multipart.Writer) error) (_ string, _ *os.File, err error) {
	buffer, err := os.CreateTemp(envx.String("", "CACHE_DIRECTORY"), "multipart.upload.bin.")
	if err != nil {
		return "", nil, errorsx.Wrap(err, "unable to create tmpfile buffer")
	}

	mw := multipart.NewWriter(buffer)

	if err = do(mw); err != nil {
		return "", nil, err
	}

	// Close the form
	if err = mw.Close(); err != nil {
		return "", nil, errorsx.Wrap(err, "unable to close writer request")
	}

	if err = iox.Rewind(buffer); err != nil {
		return "", nil, errorsx.Wrap(err, "rewind buffer")
	}

	return mw.FormDataContentType(), buffer, nil
}
