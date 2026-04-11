package httpx_test

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/textproto"
	"testing"

	"github.com/retrovibed/retrovibed/internal/cryptox"
	"github.com/retrovibed/retrovibed/internal/httpx"
	"github.com/retrovibed/retrovibed/internal/md5x"
	"github.com/retrovibed/retrovibed/internal/testx"
	"github.com/stretchr/testify/require"
)

func TestMultipart(t *testing.T) {
	t.Run("read the generate parts", func(t *testing.T) {
		var expected = md5.New()

		contentType, content, err := httpx.Multipart(func(mw *multipart.Writer) error {
			if err := mw.WriteField("message", t.Name()); err != nil {
				return fmt.Errorf("failed to write text field: %w", err)
			}

			dst, err := mw.CreateFormFile("upload", "random.bin")
			if err != nil {
				return fmt.Errorf("failed to create form file: %w", err)
			}

			if _, err := io.Copy(dst, io.LimitReader(io.TeeReader(cryptox.NewChaCha8(t.Name()), expected), 8)); err != nil {
				return fmt.Errorf("failed to write file content: %w", err)
			}

			return nil
		})
		require.NoError(t, err)
		defer content.Close()

		// 1. Verify Content-Type
		d, params, err := mime.ParseMediaType(contentType)
		require.NoError(t, err)
		require.Equal(t, "multipart/form-data", d)
		boundary := params["boundary"]

		mr := multipart.NewReader(content, boundary)

		// Verify text part
		p, err := mr.NextPart()
		require.NoError(t, err)
		require.Equal(t, "message", p.FormName())
		require.Equal(t, t.Name(), testx.IOString(p))
		require.NoError(t, p.Close())

		// Verify file part
		p, err = mr.NextPart()
		require.NoError(t, err)
		require.Equal(t, "upload", p.FormName())
		require.Equal(t, "random.bin", p.FileName())
		require.Equal(t, "application/octet-stream", p.Header.Get("Content-Type"))
		v := testx.IOMD5(p)
		require.Equal(t, md5x.FormatUUID(expected), v)
		require.NoError(t, p.Close())

		// Ensure no more parts
		_, err = mr.NextPart()
		require.ErrorIs(t, err, io.EOF)
	})

	t.Run("handle error in do function", func(t *testing.T) {
		expectedErr := errors.New("simulated error from doFunc")
		doFunc := func(mw *multipart.Writer) error {
			return expectedErr
		}

		_, content, err := httpx.Multipart(doFunc)
		require.NoError(t, err)
		defer content.Close()

		_, err = io.ReadAll(content)
		require.ErrorIs(t, err, expectedErr)
		require.NoError(t, content.Close())
	})

	t.Run("handle empty multipart body", func(t *testing.T) {
		doFunc := func(mw *multipart.Writer) error {
			return nil
		}

		contentType, content, err := httpx.Multipart(doFunc)
		require.NoError(t, err)
		defer content.Close()

		d, params, err := mime.ParseMediaType(contentType)
		require.NoError(t, err)
		require.Equal(t, "multipart/form-data", d)
		boundary := params["boundary"]
		require.NotEmpty(t, boundary, "boundary should not be empty")

		mr := multipart.NewReader(content, boundary)

		_, err = mr.NextPart()
		require.ErrorIs(t, err, io.EOF)
	})

	t.Run("handle file part with custom content type", func(t *testing.T) {
		testFileName := "custom_image.jpg"
		testFileContent := "This is some custom image data."
		customContentType := "image/jpeg"

		doFunc := func(mw *multipart.Writer) error {
			header := make(textproto.MIMEHeader)
			header.Set("Content-Disposition", fmt.Sprintf(`form-data; name="custom_file"; filename="%s"`, testFileName))
			header.Set("Content-Type", customContentType)

			partWriter, err := mw.CreatePart(header)
			require.NoError(t, err, "failed to create custom part")

			if _, err := io.WriteString(partWriter, testFileContent); err != nil {
				return fmt.Errorf("failed to write file content: %w", err)
			}
			return nil
		}

		contentType, content, err := httpx.Multipart(doFunc)
		require.NoError(t, err)
		defer content.Close()

		d, params, err := mime.ParseMediaType(contentType)
		require.NoError(t, err)
		require.Equal(t, "multipart/form-data", d)
		boundary := params["boundary"]
		require.NotEmpty(t, boundary, "boundary should not be empty")

		mr := multipart.NewReader(content, boundary)

		p, err := mr.NextPart()
		require.NoError(t, err)
		require.Equal(t, "custom_file", p.FormName())
		require.Equal(t, testFileName, p.FileName())
		require.Equal(t, customContentType, p.Header.Get("Content-Type"))

		require.Equal(t, testFileContent, testx.IOString(p))
		require.NoError(t, p.Close())

		_, err = mr.NextPart()
		require.ErrorIs(t, err, io.EOF)
	})
}
