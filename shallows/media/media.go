package media

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/shallows/internal/grpcx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/stringsx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/retrovibed/retrovibed/shallows/tracking"
	"google.golang.org/protobuf/encoding/protojson"
)

func (t *MediaSearchRequest) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(t)
}

func (t *MediaSearchRequest) UnmarshalJSON(b []byte) error {
	return protojson.Unmarshal(b, t)
}

type MediaOption func(*Media)

func MediaOptionFromLibraryMetadata(cc library.Metadata) MediaOption {
	return func(c *Media) {
		c.Id = cc.ID
		c.Description = cc.Description
		c.Mimetype = cc.Mimetype
		c.TorrentId = cc.TorrentID
		c.ArchiveId = cc.ArchiveID
		c.KnownMediaId = cc.KnownMediaID
		c.CreatedAt = grpcx.EncodeTime(cc.CreatedAt)
		c.UpdatedAt = grpcx.EncodeTime(cc.UpdatedAt)
		c.Mimetype = stringsx.FirstNonBlank(cc.Mimetype, mimex.Binary)
		c.EncryptionSeed = cc.EncryptionSeed
	}
}

func MediaOptionImageAuto(r *http.Request) MediaOption {
	return func(c *Media) {
		if strings.HasPrefix(c.Mimetype, mimex.Image) {
			c.Image = fmt.Sprintf("%s%s%s", httpx.HTTPRequestURL(r), r.URL.Path, c.Id)
		}
	}
}

func MediaOptionFromTorrentMetadata(cc tracking.Metadata) MediaOption {
	return func(c *Media) {
		c.Id = cc.ID
		c.TorrentId = cc.ID
		c.Description = cc.Description
		c.KnownMediaId = cc.KnownMediaID
		c.CreatedAt = grpcx.EncodeTime(cc.CreatedAt)
		c.UpdatedAt = grpcx.EncodeTime(cc.UpdatedAt)
		c.Mimetype = stringsx.FirstNonBlank(cc.Mimetype, mimex.Bittorrent)
		c.EncryptionSeed = cc.EncryptionSeed
	}
}
