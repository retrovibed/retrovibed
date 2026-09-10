package communityapi

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/retrovibed/retrovibed/retroapi/httpx"
	"github.com/retrovibed/retrovibed/retroapi/jsonx"
	"github.com/retrovibed/retrovibed/retroapi/publishplugin"
	"github.com/retrovibed/retrovibed/shallows/community"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/stringsx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/library"
)

func NewSyncPublishedBackgroundWorker(
	q sqlx.Queryer,
	httpc *http.Client,
	metrics Publisher,
	publisher FeedPublisher,
	publishers publishplugin.T,
	archiver library.Archiver,
	mvfs,
	tvfs fsx.Virtual,
) SyncPublishedBackgroundWorker {
	return SyncPublishedBackgroundWorker{
		q:          q,
		httpc:      httpc,
		metrics:    metrics,
		publisher:  publisher,
		publishers: publishers,
		archiver:   archiver,
		mvfs:       mvfs,
		tvfs:       tvfs,
	}
}

type SyncPublishedBackgroundWorker struct {
	q          sqlx.Queryer
	httpc      *http.Client
	metrics    Publisher
	publisher  FeedPublisher
	publishers publishplugin.T
	archiver   library.Archiver
	mvfs       fsx.Virtual
	tvfs       fsx.Virtual
}

func (t SyncPublishedBackgroundWorker) render(
	tmp string,
	published community.PublishedContent,
	known library.Known,
	md library.Metadata,
) string {
	tmpl, err := template.New("title").Parse(tmp)
	if err != nil {
		log.Println(errorsx.Wrap(err, "unable to parse title template"))
		return ""
	}

	var buf bytes.Buffer
	data := struct {
		Published community.PublishedContent
		Known     library.Known
		Metadata  library.Metadata
	}{
		Published: published,
		Known:     known,
		Metadata:  md,
	}

	if err := tmpl.Execute(&buf, data); err != nil {
		log.Println(errorsx.Wrap(err, "unable to execute title template"))
		return ""
	}

	return buf.String()
}

func (t SyncPublishedBackgroundWorker) Message(ctx context.Context, m []byte) (err error) {
	var (
		decoded community.PublishedContent
	)

	var (
		lmd   library.Metadata
		known library.Known
	)

	if err = jsonx.Unmarshal(m, &decoded); err != nil {
		return err
	}
	decoded = langx.Clone(decoded, timex.JSONSafeDecodeOption)

	if !decoded.TombstonedAt.Equal(timex.Inf()) {
		if _, err := t.metrics.Delete(ctx, decoded.ID); err != nil && !httpx.IgnoreError(err, http.StatusNotFound) {
			return errorsx.Wrap(err, "failed to delete from deeppool")
		}
		log.Printf("deleted published content %s from deeppool", decoded.ID)
		return nil
	}

	if err := library.MetadataFindByID(ctx, t.q, decoded.LibraryID).Scan(&lmd); err != nil {
		return errorsx.Wrap(err, "failed to find library metadata")
	}

	tmd, err := ensureTorrent(ctx, t.q, t.mvfs, t.tvfs, &lmd)
	if err != nil {
		return errorsx.Wrap(err, "failed to ensure torrent")
	}

	decoded.MagnetURI = magnetURI(tmd, lmd.Description)
	if err := community.PublishedContentUpdateMagnetURI(ctx, t.q, decoded.ID, decoded.MagnetURI).Scan(&decoded); err != nil {
		return errorsx.Wrap(err, "failed to update magnet_uri")
	}

	if decoded.KnownMediaID != "" {
		if err := library.KnownFindByID(ctx, t.q, decoded.KnownMediaID).Scan(&known); sqlx.IgnoreNoRows(err) != nil {
			log.Println(errorsx.Wrap(err, "failed to find known media"))
		}
	}

	if decoded.OAuthGoogleID != uuid.Nil.String() {
		if uerr := community.YouTubeUpload(ctx, t.q, t.httpc, t.mvfs, decoded.OAuthGoogleID, lmd, stringsx.FirstNonBlank(known.Title, lmd.Description), known.Overview); uerr != nil {
			log.Println(errorsx.Wrap(uerr, "youtube cross-post failed"))
		}
	}

	enabled := sqlx.Scan(community.CommunityPublisherFindByCommunityID(ctx, t.q, decoded.CommunityID))
	for cp := range enabled.Iter() {
		var pub community.PluginPublisher
		if err := community.PluginPublisherFindByID(ctx, t.q, cp.PublisherID).Scan(&pub); err != nil {
			log.Println(errorsx.Wrap(err, "unable to find plugin publisher"))
			continue
		}

		req := publishplugin.Request{
			MediaPath:   lmd.ID,
			Title:       stringsx.FirstNonBlank(t.render(cp.TemplateTitle, decoded, known, lmd), decoded.Title, known.Title, lmd.Description),
			Description: stringsx.FirstNonBlank(t.render(cp.TemplateDescription, decoded, known, lmd), decoded.Title, known.Title, lmd.Description),
			Mimetype:    stringsx.FirstNonBlank(decoded.Mimetype, known.Mimetype, lmd.Mimetype),
			CommunityID: decoded.CommunityID,
			Magnet:      decoded.MagnetURI,
			Adult:       known.Adult,
		}

		if err := publishToPlugin(ctx, t.mvfs, t.publishers, pub, req); err != nil {
			log.Println(errorsx.Wrap(err, "plugin publish failed"))
		}
	}

	if err := enabled.Err(); err != nil {
		log.Println(errorsx.Wrap(err, "unable to list enabled publishers"))
	}

	if decoded.PublishMode == int32(PublishMode_UNLISTED) {
		if err := community.PublishedContentUpdatePublishedAt(ctx, t.q, decoded.ID, time.Now()).Scan(&decoded); err != nil {
			log.Println(errorsx.Wrap(err, "failed to update published_at"))
		}
		return nil
	}

	if decoded.PublishMode == int32(PublishMode_LISTED) {
		if err := community.PublishedContentUpdatePublishedAt(ctx, t.q, decoded.ID, time.Now()).Scan(&decoded); err != nil {
			return errorsx.Wrap(err, "failed to update published_at")
		}

		if err := RegenerateFeed(ctx, t.q, t.publisher, decoded.CommunityID); err != nil {
			return errorsx.Wrap(err, "feed regeneration failed for community "+decoded.CommunityID)
		}

		log.Printf("synced listed published content %s", decoded.ID)
		return nil
	}

	if lmd.ArchiveID == uuid.Max.String() {
		return nil // archival in progress, waiting
	} else if lmd.ArchiveID == uuid.Nil.String() {
		if err := library.Archive(ctx, t.q, &lmd, t.mvfs, t.archiver); err != nil {
			log.Println(errorsx.Wrap(err, "failed to archive media"))
			return nil
		}
	}

	encoded, err := os.ReadFile(t.tvfs.Path(fmt.Sprintf("%s.torrent", int160.FromBytes(tmd.Infohash).String())))
	if err != nil {
		return errorsx.Wrap(err, "failed to open torrent file")
	}

	req := PublishContentRequest{
		PublishedContent: &PublishedContent{
			Id:             decoded.ID,
			CommunityId:    decoded.CommunityID,
			KnownMediaId:   decoded.KnownMediaID,
			MagnetUri:      decoded.MagnetURI,
			Title:          stringsx.FirstNonBlank(known.Title, lmd.Description),
			Mimetype:       stringsx.FirstNonBlank(known.Mimetype, lmd.Mimetype),
			Description:    known.Overview,
			ArchivedId:     lmd.ArchiveID,
			EncryptionSeed: lmd.EncryptionSeed,
			Bytes:          lmd.Bytes,
		},
	}
	if _, err = t.metrics.Publish(ctx, &req, io.NopCloser(bytes.NewReader(encoded))); err != nil {
		return errorsx.Wrap(err, "failed to sync to deeppool")
	}

	if err := community.PublishedContentUpdatePublishedAt(ctx, t.q, decoded.ID, time.Now()).Scan(&decoded); err != nil {
		return errorsx.Wrap(err, "failed to update published_at")
	}

	if err := RegenerateFeed(ctx, t.q, t.publisher, decoded.CommunityID); err != nil {
		return errorsx.Wrapf(err, "feed regeneration failed for community: %s", decoded.CommunityID)
	}

	log.Printf("synced published content %s to deeppool with archive_id %s", decoded.ID, lmd.ArchiveID)

	return nil
}
