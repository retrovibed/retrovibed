package media

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime/multipart"

	"github.com/james-lawrence/torrent"
	"github.com/james-lawrence/torrent/metainfo"
	"github.com/retrovibed/retrovibed/internal/errorsx"
	"github.com/retrovibed/retrovibed/internal/httpx"
	"github.com/retrovibed/retrovibed/internal/mimex"
)

// PublishRequest a torrent to an http endpoint in its entirety
func PublishRequest(ctx context.Context, md torrent.Metadata, req *PublishedUploadRequest) (boundary string, _ io.ReadCloser, err error) {
	encodedreq, err := json.Marshal(req)
	if err != nil {
		return "", nil, errorsx.Wrap(err, "unable to encode publish metadata")
	}

	encoded, err := metainfo.Encode(md.Metainfo())
	if err != nil {
		return "", nil, errorsx.Wrap(err, "unable to encode torrent file")
	}

	info, err := md.Metainfo().UnmarshalInfo()
	if err != nil {
		return "", nil, errorsx.Wrap(err, "unable to decode torrent info")
	}

	disk, err := md.Storage.OpenTorrent(&info, md.ID)
	if err != nil {
		return "", nil, errorsx.Wrap(err, "unable to open torrent reader")
	}

	return httpx.Multipart(func(w *multipart.Writer) (merr error) {
		part, lerr := w.CreatePart(httpx.NewMultipartHeader(mimex.Bittorrent, "torrentfile", md.DisplayName))
		if lerr != nil {
			return errorsx.Wrap(lerr, "unable to create torrent file part")
		}

		if _, lerr = io.Copy(part, bytes.NewBuffer(encoded)); lerr != nil {
			return errorsx.Wrap(lerr, "unable to copy archive")
		}

		content, lerr := w.CreatePart(httpx.NewMultipartHeader(mimex.Binary, "contents", ""))
		if lerr != nil {
			return errorsx.Wrap(lerr, "unable to create torrent file part")
		}

		if n, lerr := io.Copy(content, io.NewSectionReader(disk, 0, info.TotalLength())); lerr != nil {
			return errorsx.Wrap(lerr, "unable to copy torrent")
		} else if n < info.TotalLength() {
			return errorsx.Errorf("failed to copy entire torrent %d < %d", n, info.TotalLength())
		}

		metadata, lerr := w.CreatePart(httpx.NewMultipartHeader(mimex.JSON, "metadata", "metadata.json"))
		if lerr != nil {
			return errorsx.Wrap(lerr, "unable to create torrent file part")
		}

		if _, lerr = io.Copy(metadata, bytes.NewBuffer(encodedreq)); lerr != nil {
			return errorsx.Wrap(lerr, "unable to copy metadata")
		}
		return nil
	})
}
