package communityapi

import (
	"github.com/retrovibed/retrovibed/shallows/community"
	"github.com/retrovibed/retrovibed/shallows/internal/grpcx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
)

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
		dst.TemplateTitle = p.TemplateTitle
		dst.TemplateDescription = p.TemplateDescription
	}
}

// CommunityPublisherOptionFromProto converts proto fields to a database model option.
func CommunityPublisherOptionFromProto(p *CommunityPublisher) func(*community.CommunityPublisher) {
	return func(dst *community.CommunityPublisher) {
		dst.ID = p.Id
		dst.CommunityID = p.CommunityId
		dst.PublisherID = p.PublisherId
		dst.TemplateTitle = p.TemplateTitle
		dst.TemplateDescription = p.TemplateDescription
	}
}

func CommunityPublisherOptionTemplates(title, desc string) func(*CommunityPublisher) {
	return func(dst *CommunityPublisher) {
		dst.TemplateTitle = title
		dst.TemplateDescription = desc
	}
}

func NewCommunitySocial(opts ...func(*CommunitySocial)) *CommunitySocial {
	return new(langx.Clone(CommunitySocial{}, opts...))
}
