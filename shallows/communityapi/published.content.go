package communityapi

import (
	"github.com/retrovibed/retrovibed/shallows/community"
	"github.com/retrovibed/retrovibed/shallows/internal/grpcx"
)

// PublishedContentOptionFromDB converts a database model to proto options.
func PublishedContentOptionFromDB(pc community.PublishedContent) func(*PublishedContent) {
	return func(p *PublishedContent) {
		p.Id = pc.ID
		p.Title = pc.Title
		p.Description = pc.Description
		p.CommunityId = pc.CommunityID
		p.KnownMediaId = pc.KnownMediaID
		p.MagnetUri = pc.MagnetURI
		p.LibraryId = pc.LibraryID
		p.OauthGoogleId = pc.OAuthGoogleID
		p.PublishedAt = grpcx.EncodeTime(pc.PublishedAt)
		p.CreatedAt = grpcx.EncodeTime(pc.CreatedAt)
		p.UpdatedAt = grpcx.EncodeTime(pc.UpdatedAt)
		p.Bytes = pc.Bytes
	}
}
