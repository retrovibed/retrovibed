package media

import (
	"github.com/retrovibed/retrovibed/internal/grpcx"
	"github.com/retrovibed/retrovibed/internal/stringsx"
	"github.com/retrovibed/retrovibed/library"
)

type KnownOption func(*Known)

func KnownOptionFromLibraryKnown(cc library.Known) KnownOption {
	return func(c *Known) {
		c.Id = cc.UID
		c.Image = stringsx.FirstNonBlank(cc.PosterPath, cc.BackdropPath)
		c.Rating = float32(cc.Popularity)
		c.Description = cc.Title
		c.Summary = cc.Overview
		c.Adult = cc.Adult
		c.Released = grpcx.EncodeTime(cc.Released)
	}
}
