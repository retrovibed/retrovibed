package community

import (
	"context"
	"io"
	"net/http"

	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"

	"github.com/retrovibed/retrovibed/retroapi/blockcache"
	"github.com/retrovibed/retrovibed/shallows/deeppool"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/library"
)

// YouTubeUpload uploads a library media file to YouTube as an unlisted video.
func YouTubeUpload(ctx context.Context, q sqlx.Queryer, httpc *http.Client, mediastorage fsx.Virtual, oauthID string, lmd library.Metadata, title, description string) error {
	var row OAuth2Google
	if err := OAuth2GoogleFindByID(ctx, q, oauthID).Scan(&row); err != nil {
		return errorsx.Wrap(err, "no youtube credentials found")
	}

	ts := deeppool.NewYouTube(httpc).TokenSource(&oauth2.Token{
		AccessToken:  row.AccessToken,
		TokenType:    row.TokenType,
		RefreshToken: row.RefreshToken,
		Expiry:       row.Expiry,
	})

	cache, err := blockcache.NewDirectoryCache(mediastorage.Path(lmd.ID))
	if err != nil {
		return errorsx.Wrap(err, "unable to create media reader")
	}

	svc, err := youtube.NewService(ctx, option.WithTokenSource(ts))
	if err != nil {
		return errorsx.Wrap(err, "unable to create youtube service")
	}

	_, err = svc.Videos.Insert([]string{"snippet", "status"}, &youtube.Video{
		Snippet: &youtube.VideoSnippet{
			Title:       title,
			Description: description,
		},
		Status: &youtube.VideoStatus{
			PrivacyStatus: "unlisted",
		},
	}).Media(io.NewSectionReader(cache, int64(lmd.DiskOffset), int64(lmd.Bytes))).Do()

	return errorsx.Wrap(err, "youtube upload failed")
}
