package media

import (
	"github.com/retrovibed/retrovibed/shallows/internal/grpcx"
	"github.com/retrovibed/retrovibed/shallows/internal/stringsx"
	"github.com/retrovibed/retrovibed/shallows/library"
)

type KnownOption func(*Known)

func KnownOptionFromLibraryKnown(cc library.Known) KnownOption {
	return func(c *Known) {
		c.Id = cc.ID
		c.Uid = cc.UID
		c.Image = stringsx.FirstNonBlank(cc.PosterPath, cc.BackdropPath)
		c.Rating = float32(cc.Popularity)
		c.Description = cc.Title
		c.Summary = cc.Overview
		c.Adult = cc.Adult
		c.Released = grpcx.EncodeTime(cc.Released)
		c.Mimetype = cc.Mimetype
		c.Source = cc.Source
	}
}

func KnownOptionFromRecommendation(cc library.Recommendation) KnownOption {
	return func(c *Known) {
		c.Id = cc.ID
		c.Uid = cc.ContentID
		c.Image = cc.Image
		c.Rating = float32(cc.Popularity)
		c.Description = cc.Title
		c.Summary = cc.Overview
		c.Adult = cc.Adult
		c.Released = grpcx.EncodeTime(cc.Released)
		c.Mimetype = cc.Mimetype
		c.Source = cc.Source
	}
}
