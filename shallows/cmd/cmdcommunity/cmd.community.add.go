package cmdcommunity

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/deeppool"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/meta"
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
		resp  *meta.PublishContentResponse
	)

	log.Println("community add initiated", t.Community)
	defer log.Println("community add completed", t.Community)

	if httpc, err = dpc.HTTPClient(gctx.Context); err != nil {
		return err
	}

	client := deeppool.NewPublished(httpc)
	pc := &meta.PublishedContent{
		KnownMediaId: t.KnownMediaID,
		MagnetUri:    t.MagnetURI,
		ArchivedId:   t.ArchivedID,
	}

	if resp, err = client.Publish(gctx.Context, t.Community, pc); err != nil {
		return errorsx.Wrap(err, "failed to publish content to deeppool")
	}

	log.Println("published content", resp.PublishedContent.Id)

	if err = json.NewEncoder(os.Stdout).Encode(resp.PublishedContent); err != nil {
		return errorsx.Wrap(err, "unable to encode response")
	}

	return nil
}
