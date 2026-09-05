// Command publishplugin-activitypub publishes a community's content to the
// fediverse by posting it into a Lemmy community, as a publish plugin for
// retroapi/publishplugin's Registry.
//
// It speaks Lemmy's own HTTP API (/api/v3) rather than raw ActivityPub:
// server-to-server ActivityPub would require this host to serve a publicly
// reachable actor document and answer signature fetches, which a sandboxed,
// outbound-only wasm guest cannot do. Posting through the instance's client
// api gets the content federated by Lemmy itself, which is the point.
//
// build it for the registry with:
//
//	GOOS=wasip1 GOARCH=wasm go build -o lemmy.wasm ./retroapi/publishplugin-activitypub
//
// and drop lemmy.wasm in the well-known publish.d plugin directory (see
// publishplugin.PublishPluginDir) to have Registry load and run it.
//
// One installed copy posts to exactly one Lemmy community, because its
// configuration is a sidecar .env keyed by the installed filename. To post
// to several, symlink the module once per target and give each link its own
// .env:
//
//	publish.d/lemmy-movies.wasm -> lemmy.wasm   + publish.d/lemmy-movies.env
//	publish.d/lemmy-music.wasm  -> lemmy.wasm   + publish.d/lemmy-music.env
//
// Registry keys a plugin's configuration, cache and config directories off
// that filename, so the links stay independent of each other.
//
// Two kong subcommands are exposed. "publish", invoked by
// publishplugin.Registry.Publish as:
//
//	<binary> publish --title <t> --description <d> --mimetype <m> [--media <path>] [--community-id <id>] [--link <uri>]
//
// and "env", invoked by publishplugin.Registry.Environment as:
//
//	<binary> env
//
// which prints the variables below as a .env document. That output is this
// plugin's configuration schema - it is what the console renders its
// settings form from.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/alecthomas/kong"
	"github.com/retrovibed/retrovibed/retroapi/bytesx"
	"github.com/retrovibed/retrovibed/retroapi/mimex"

	// autohijack points net.DefaultResolver and http.DefaultTransport at
	// wasinet's virtual sockets when this is built for wasip1 (a no-op on
	// every other platform). Registry.Publish mounts the host's real TLS
	// trust store into the sandbox and sets SSL_CERT_DIR before running a
	// plugin, so the plain net/http calls in lemmy.go - TLS verification
	// included - work unmodified once this import is in place. Without it
	// a wasip1 build has no network at all.
	_ "github.com/egdaemon/wasinet/wasinet/autohijack"
)

var cli struct {
	Publish publishCmd `cmd:"" help:"invoked by publishplugin.Registry to publish content to a lemmy community"`
	Env     envCmd     `cmd:"" help:"invoked by publishplugin.Registry to report the variables this plugin understands"`
}

// publishCmd's flags fall into two halves: the content ones the registry
// passes on every invocation, and the configuration ones resolved from this
// installation's .env (kong reads them straight out of the process
// environment the registry populates from the sidecar file).
type publishCmd struct {
	Title       string `flag:"" name:"title" help:"title of the content being published"`
	Description string `flag:"" name:"description" help:"description of the content being published"`
	Mimetype    string `flag:"" name:"mimetype" help:"mimetype of the content being published"`
	Media       string `flag:"" name:"media" help:"guest path to the mounted media file, when the caller provided one"`
	CommunityID string `flag:"" name:"community-id" help:"id of the retrovibed community the content is being published on behalf of"`
	Link        string `flag:"" name:"link" help:"publicly reachable uri for the content, when it has one"`

	Instance     string `flag:"" name:"instance" help:"base url of the lemmy instance to post to" required:"" env:"LEMMY_INSTANCE"`
	Community    string `flag:"" name:"community" help:"lemmy community to post into" required:"" env:"LEMMY_COMMUNITY"`
	Username     string `flag:"" name:"username" help:"lemmy account to post as" env:"LEMMY_USERNAME"`
	Password     string `flag:"" name:"password" help:"password for that account" env:"LEMMY_PASSWORD"`
	Token        string `flag:"" name:"token" help:"pre-issued jwt, used instead of logging in" env:"LEMMY_TOKEN"`
	TOTP         string `flag:"" name:"totp" help:"one time code, when the account has 2fa enabled" env:"LEMMY_TOTP"`
	NSFW         bool   `flag:"" name:"nsfw" help:"mark every post from this installation as nsfw" env:"LEMMY_NSFW"`
	LanguageID   int64  `flag:"" name:"language-id" help:"lemmy language id to tag posts with; 0 leaves it unset" env:"LEMMY_LANGUAGE_ID"`
	ThumbnailMax string `flag:"" name:"thumbnail-max" help:"largest image to upload as a post thumbnail" default:"8 MB" env:"LEMMY_THUMBNAIL_MAX"`
}

// envCmd takes no flags - it is a pure declaration of what this plugin can
// be configured with.
type envCmd struct{}

// environment is this plugin's configuration schema, in the format
// retroapi/envfile parses: the comment block immediately preceding a
// KEY=VALUE line becomes that variable's help text. Every key here matches
// an env tag on publishCmd above; changing one means changing both.
const environment = `# base url of the lemmy instance to post to
LEMMY_INSTANCE=""
# lemmy community to post into, either "movies" for one local to the
# instance above or "movies@lemmy.ml" for a remote one
LEMMY_COMMUNITY=""
# lemmy account to post as; not needed when LEMMY_TOKEN is set
LEMMY_USERNAME=""
# password for that account; not needed when LEMMY_TOKEN is set
LEMMY_PASSWORD=""
# a pre-issued jwt to use instead of logging in with a username/password
LEMMY_TOKEN=""
# one time code, only when the account has 2fa enabled
LEMMY_TOTP=""
# mark every post from this installation as nsfw
LEMMY_NSFW="false"
# lemmy language id to tag posts with; 0 leaves it unset
LEMMY_LANGUAGE_ID="0"
# largest image to upload as a post thumbnail; larger content is posted as
# a plain link and lemmy generates its own thumbnail
LEMMY_THUMBNAIL_MAX="8 MB"
`

// Run prints the declaration verbatim on stdout. Registry.Environment
// returns these bytes untouched, so anything else written here would end up
// in the configuration form.
func (cmd *envCmd) Run(ctx context.Context) error {
	_, err := fmt.Fprint(os.Stdout, environment)
	return err
}

// Run posts the content to lemmy and emits the single JSON object
// Registry.Publish decodes as its *Result. Everything else - progress,
// failures that were recovered from - goes to stderr, which the registry
// logs rather than parses.
func (cmd *publishCmd) Run(ctx context.Context) error {
	client, err := cmd.connect(ctx)
	if err != nil {
		return err
	}

	post, err := cmd.post(ctx, client)
	if errors.Is(err, ErrUnauthorized) && cmd.Username != "" {
		// the cached session outlived its welcome; the credentials are
		// still on hand, so log in again rather than failing a publish
		// that would succeed on the next run anyway.
		fmt.Fprintln(os.Stderr, "publishplugin-activitypub: cached session rejected, logging in again")

		if err = client.Login(ctx, cmd.Username, cmd.Password, cmd.TOTP); err != nil {
			return err
		}

		post, err = cmd.post(ctx, client)
	}

	if err != nil {
		return err
	}

	return json.NewEncoder(os.Stdout).Encode(struct {
		URL        string `json:"url"`
		ExternalID string `json:"external_id"`
		Status     string `json:"status"`
	}{
		URL:        post.ApID,
		ExternalID: strconv.FormatInt(post.ID, 10),
		Status:     "published",
	})
}

// connect builds a client already holding credentials: a configured token
// wins, then a session cached by a previous invocation, and only failing
// both does this spend a round trip logging in.
func (cmd *publishCmd) connect(ctx context.Context) (*Client, error) {
	client, err := NewClient(
		cmd.Instance,
		OptionToken(cmd.Token),
		OptionCacheDir(os.Getenv("CACHE_DIRECTORY")),
	)
	if err != nil {
		return nil, err
	}

	if client.Authenticated() {
		return client, nil
	}

	if client.RestoreSession() {
		return client, nil
	}

	if cmd.Username == "" {
		return nil, errors.New("no credentials configured: set LEMMY_TOKEN, or LEMMY_USERNAME and LEMMY_PASSWORD")
	}

	return client, client.Login(ctx, cmd.Username, cmd.Password, cmd.TOTP)
}

// post resolves the target community and submits the content to it. It is
// the unit retried after a re-login, so it must be safe to call twice - it
// is: nothing here mutates the client or the filesystem.
func (cmd *publishCmd) post(ctx context.Context, client *Client) (*Post, error) {
	community, err := client.ResolveCommunity(ctx, cmd.Community)
	if err != nil {
		return nil, err
	}

	return client.CreatePost(ctx, CreatePost{
		Name:            cmd.Title,
		CommunityID:     community,
		URL:             cmd.Link,
		Body:            cmd.Description,
		NSFW:            cmd.NSFW,
		LanguageID:      cmd.LanguageID,
		CustomThumbnail: cmd.thumbnail(ctx, client),
	})
}

// thumbnail uploads the mounted media as a post thumbnail when it is an
// image small enough to be worth hosting, returning "" in every other case
// - including on failure. A post that lost its thumbnail is still a
// perfectly good post, and lemmy will try to derive its own from the link,
// so nothing here is allowed to fail a publish.
func (cmd *publishCmd) thumbnail(ctx context.Context, client *Client) string {
	if cmd.Media == "" {
		return ""
	}

	// lemmy posts have no attachments - only a thumbnail image referenced
	// by url - so anything that isn't an image has nothing to contribute
	// here. Torrented video, the common case, never reaches the upload.
	if mimex.Category(cmd.Mimetype) != mimex.Image {
		return ""
	}

	limit := bytesx.Parse(cmd.ThumbnailMax)
	if limit == 0 {
		fmt.Fprintln(os.Stderr, "publishplugin-activitypub: unusable thumbnail limit, skipping upload:", cmd.ThumbnailMax)
		return ""
	}

	info, err := os.Stat(cmd.Media)
	if err != nil {
		fmt.Fprintln(os.Stderr, "publishplugin-activitypub: unable to stat media, skipping thumbnail:", err)
		return ""
	}

	if uint64(info.Size()) > limit {
		fmt.Fprintln(os.Stderr, "publishplugin-activitypub: media larger than", cmd.ThumbnailMax, "- posting without a thumbnail")
		return ""
	}

	uploaded, err := client.UploadImage(ctx, cmd.Media)
	if err != nil {
		fmt.Fprintln(os.Stderr, "publishplugin-activitypub: thumbnail upload failed, posting without one:", err)
		return ""
	}

	return uploaded
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	kctx := kong.Parse(&cli, kong.BindTo(ctx, (*context.Context)(nil)))
	kctx.FatalIfErrorf(kctx.Run())
}
