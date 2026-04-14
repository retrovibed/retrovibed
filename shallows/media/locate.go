package media

import (
	"github.com/retrovibed/retrovibed/shallows/internal/grpcx"
	"github.com/retrovibed/retrovibed/shallows/library"
)

type LocateOption func(*Locate)

func LocateOptionFromLibraryLocate(cc library.Locate) LocateOption {
	return func(c *Locate) {
		c.Id = cc.ID
		c.CreatedAt = grpcx.EncodeTime(cc.CreatedAt)
		c.UpdatedAt = grpcx.EncodeTime(cc.UpdatedAt)
		c.KnownMediaId = cc.KnownMediaID
		c.LocatedTorrentId = cc.LocatedTorrentID
	}
}

func NewLibraryLocateFromLocate(cc *Locate) (l library.Locate, err error) {
	err = grpcx.JSONEncode(cc, &l)
	return l, err
}
