package communityapi

import (
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
	}
}
