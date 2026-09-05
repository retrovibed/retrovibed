package cmdlibrary

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/retrovibed/retrovibed/retroapi/authn"
	"github.com/retrovibed/retrovibed/retroapi/jsonx"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/media"
)

type importJSONL struct{}

func (t importJSONL) Run(gctx *cmdopts.Global, tls *cmdopts.TLSConfig, daemon *cmdopts.Endpoint) error {
	c := authn.AutoOauth2Client(gctx.Context, tls.Config(), authn.EndpointSSHAuth(daemon.Endpoint))
	return t.run(gctx.Context, c, daemon.Endpoint, os.Stdin)
}

func (t importJSONL) run(ctx context.Context, c *http.Client, endpoint string, r io.Reader) error {
	dec := json.NewDecoder(r)
	var imported uint64

	for {
		var hdr exportHeader
		if err := dec.Decode(&hdr); err == io.EOF {
			break
		} else if err != nil {
			return errorsx.Wrap(err, "decode header")
		}

		if err := t.importItem(ctx, c, endpoint, dec, hdr.Chunks); err != nil {
			return err
		}

		if imported++; imported%256 == 0 {
			log.Println("imported", imported, "records")
		}
	}

	log.Println("imported", imported, "records total")
	return nil
}

func (t importJSONL) importItem(ctx context.Context, c *http.Client, endpoint string, dec *json.Decoder, numChunks uint64) error {
	h := md5.New()
	var (
		data    bytes.Buffer
		trailer exportTrailer
	)

	w := io.MultiWriter(h, &data)
	for i := range numChunks {
		var chunk exportChunk
		if err := dec.Decode(&chunk); err != nil {
			return errorsx.Wrapf(err, "decode chunk %d", i)
		}

		if _, err := w.Write(chunk.Data); err != nil {
			return errorsx.Wrap(err, "failed to write chunk")
		}
	}

	if err := dec.Decode(&trailer); err != nil {
		return errorsx.Wrap(err, "decode trailer")
	}
	trailer.Metadata = langx.Clone(trailer.Metadata, timex.JSONSafeDecodeOption)

	got := hex.EncodeToString(h.Sum(nil))
	if got != trailer.MD5 {
		return errorsx.Errorf("MD5 mismatch for %s: got %s want %s", trailer.Metadata.ID, got, trailer.MD5)
	}

	contentType, body, err := httpx.Multipart(func(mw *multipart.Writer) error {
		part, lerr := mw.CreatePart(httpx.NewMultipartHeader(trailer.Metadata.Mimetype, "content", trailer.Metadata.Description))
		if lerr != nil {
			return lerr
		}
		_, lerr = io.Copy(part, &data)
		return lerr
	})
	if err != nil {
		return errorsx.Wrap(err, "create multipart")
	}
	defer body.Close()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint+"/m/", body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", contentType)

	resp, err := httpx.AsError(c.Do(req))
	if err != nil {
		return errorsx.Wrap(err, "upload")
	}
	defer resp.Body.Close()

	var uploadResp media.MediaUploadResponse
	if err := jsonx.UnmarshalRead(resp.Body, &uploadResp); err != nil {
		return errorsx.Wrap(err, "decode upload response")
	}

	patch, err := json.Marshal(&media.MediaUpdateRequest{
		Media: &media.Media{
			Description:  trailer.Metadata.Description,
			KnownMediaId: trailer.Metadata.KnownMediaID,
			ArchiveId:    trailer.Metadata.ArchiveID,
		},
	})
	if err != nil {
		return err
	}

	patchReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		endpoint+"/m/"+uploadResp.Media.Id+"/metadatasync",
		bytes.NewReader(patch))
	if err != nil {
		return err
	}
	patchReq.Header.Set("Content-Type", "application/json")

	patchResp, err := httpx.AsError(c.Do(patchReq))
	if err != nil {
		return errorsx.Wrap(err, "patch metadata")
	}
	patchResp.Body.Close()

	return nil
}
