package communityapi

import (
	"github.com/retrovibed/retrovibed/shallows/community"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/grpcx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/library"
)

// PublishedContentOptionFromDB converts a database model to proto options.
func PublishedContentOptionFromDB(pc community.PublishedContent) func(*PublishedContent) {
	return func(p *PublishedContent) {
		p.Id = pc.ID
		p.Title = pc.Title
		p.Description = pc.Description
		p.CommunityId = pc.CommunityID
		p.KnownMediaId = pc.KnownMediaID
		p.Mimetype = pc.Mimetype
		p.MagnetUri = pc.MagnetURI
		p.LibraryId = pc.LibraryID
		p.OauthGoogleId = pc.OAuthGoogleID
		p.PublishedAt = grpcx.EncodeTime(pc.PublishedAt)
		p.CreatedAt = grpcx.EncodeTime(pc.CreatedAt)
		p.UpdatedAt = grpcx.EncodeTime(pc.UpdatedAt)
		p.Bytes = pc.Bytes
	}
}

// PublishedContentOptionFromProto converts proto fields to a database model option.
func PublishedContentOptionFromProto(pc *PublishedContent) func(*community.PublishedContent) {
	return func(p *community.PublishedContent) {
		p.ID = pc.Id
		p.Mimetype = pc.Mimetype
		p.Title = pc.Title
		p.Description = pc.Description
		p.CommunityID = pc.CommunityId
		p.MagnetURI = pc.MagnetUri
		p.OAuthGoogleID = pc.OauthGoogleId
		p.Bytes = pc.Bytes
		p.PublishedAt = langx.FirstNonZero(errorsx.Zero(grpcx.DecodeTime(pc.PublishedAt)), timex.Inf())
	}
}

// PublishedContentOptionFromLibraryMetadata applies fields the library record is
// authoritative for — the actual byte size and detected mimetype — since these
// describe the file on disk rather than user-supplied input.
func PublishedContentOptionFromLibraryMetadata(lmd library.Metadata) func(*community.PublishedContent) {
	return func(p *community.PublishedContent) {
		p.LibraryID = lmd.ID
		p.Bytes = lmd.Bytes
		p.Mimetype = lmd.Mimetype
	}
}
