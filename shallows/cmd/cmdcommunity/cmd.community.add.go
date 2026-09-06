package cmdcommunity

import (
	"log"
	"net/http"
	"os"

	"github.com/retrovibed/retrovibed/retroapi/jsonx"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/communityapi"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
)

type cmdCommunityAdd struct {
	Community    string `arg:"" name:"community" help:"community ID to add content to" required:"true"`
	MagnetURI    string `flag:"" name:"magnet" help:"magnet URI of the content" required:"true"`
	KnownMediaID string `flag:"" name:"known-media-id" help:"known media ID (e.g., TMDB ID)" default:""`
	ArchivedID   string `flag:"" name:"archived-id" help:"archived content ID in CAS" default:""`
}

func (t cmdCommunityAdd) Run(gctx *cmdopts.Global, dpc cmdopts.DeeppoolClient) (err error) {
	var (
		httpc *http.Client
		resp  *communityapi.PublishContentResponse
	)

	log.Println("community add initiated", t.Community)
	defer log.Println("community add completed", t.Community)

	if httpc, err = dpc.HTTPClient(gctx.Context); err != nil {
		return err
	}

	client := communityapi.NewDeeppoolPublished(httpc)
	pc := &communityapi.PublishedContent{
		KnownMediaId: t.KnownMediaID,
		MagnetUri:    t.MagnetURI,
		ArchivedId:   t.ArchivedID,
	}

	if resp, err = client.Publish(gctx.Context, t.Community, pc); err != nil {
		return errorsx.Wrap(err, "failed to publish content to deeppool")
	}

	log.Println("published content", resp.PublishedContent.Id)

	if err = jsonx.MarshalWrite(os.Stdout, resp.PublishedContent); err != nil {
		return errorsx.Wrap(err, "unable to encode response")
	}

	return nil
}
