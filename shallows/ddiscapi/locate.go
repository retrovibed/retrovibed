package ddiscapi

import (
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/internal/grpcx"
)

type LocateOption func(*Locate)

func LocateOptionFromDdiscLocate(cc ddisc.Locate) LocateOption {
	return func(c *Locate) {
		c.Id = cc.ID
		c.CreatedAt = grpcx.EncodeTime(cc.CreatedAt)
		c.UpdatedAt = grpcx.EncodeTime(cc.UpdatedAt)
		c.TombstonedAt = grpcx.EncodeTime(cc.TombstonedAt)
		c.KnownMediaId = cc.KnownMediaID
		c.LocatedTorrentId = cc.LocatedTorrentID
		c.Query = cc.Query
		c.Mimetype = cc.Mimetype
	}
}

func NewDdiscLocateFromLocate(cc *Locate) (l ddisc.Locate, err error) {
	err = grpcx.JSONEncode(cc, &l)
	return l, err
}
