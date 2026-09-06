package communityapi

import (
	"log"

	"github.com/davecgh/go-spew/spew"
	"github.com/retrovibed/retrovibed/shallows/community"
	"github.com/retrovibed/retrovibed/shallows/internal/grpcx"
)

func NewPluginPublisher(opts ...func(*PluginPublisher)) *PluginPublisher {
	var p PluginPublisher
	for _, opt := range opts {
		opt(&p)
	}

	return &p
}

// PluginPublisherOptionFromDB converts a local DB plugin publisher to proto options.
// Call site should apply timex.JSONSafeEncodeOption before passing the DB value.
func PluginPublisherOptionFromDB(p community.PluginPublisher) func(*PluginPublisher) {
	return func(dst *PluginPublisher) {
		dst.Id = p.ID
		dst.Path = p.Path
		dst.Description = p.Description
		dst.Mimetype = p.Mimetype
		dst.CreatedAt = grpcx.EncodeTime(p.CreatedAt)
		dst.UpdatedAt = grpcx.EncodeTime(p.UpdatedAt)
		log.Println("DERP DERP 0", spew.Sdump(p))
		log.Println("DERP DERP 1", spew.Sdump(dst))
	}
}

func NewCommunityPublisher(opts ...func(*CommunityPublisher)) *CommunityPublisher {
	var p CommunityPublisher
	for _, opt := range opts {
		opt(&p)
	}

	return &p
}

// CommunityPublisherOptionFromDB converts a local DB community publisher to proto options.
// Call site should apply timex.JSONSafeEncodeOption before passing the DB value.
func CommunityPublisherOptionFromDB(p community.CommunityPublisher) func(*CommunityPublisher) {
	return func(dst *CommunityPublisher) {
		dst.Id = p.ID
		dst.CommunityId = p.CommunityID
		dst.PublisherId = p.PublisherID
		dst.CreatedAt = grpcx.EncodeTime(p.CreatedAt)
		dst.UpdatedAt = grpcx.EncodeTime(p.UpdatedAt)
	}
}

func NewCommunitySocial(opts ...func(*CommunitySocial)) *CommunitySocial {
	var s CommunitySocial
	for _, opt := range opts {
		opt(&s)
	}

	return &s
}
