package cmdddisc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/netip"
	"net/url"
	"strconv"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/james-lawrence/torrent"
	"github.com/james-lawrence/torrent/dht"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/james-lawrence/torrent/storage"
	"github.com/retrovibed/retrovibed/retroapi/authn"
	retronetx "github.com/retrovibed/retrovibed/retroapi/netx"
	"github.com/retrovibed/retrovibed/retroapi/userx"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/ddiscapi"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/formx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/torrentx"
)

type cmdDiscoveryIdentify struct {
	Endpoint    string        `flag:"" name:"library" help:"http address for the library you want to connect to" default:"localhost:9998"`
	ID          string        `flag:"" name:"id" help:"unknown-hash id to identify" required:"true"`
	PeerTimeout time.Duration `flag:"" name:"peer-timeout" help:"how long to wait for a reachable peer" default:"1m"`
	InfoTimeout time.Duration `flag:"" name:"info-timeout" help:"how long to wait for torrent metadata/content" default:"10m"`
	DHTPeers    []string      `flag:"" name:"dht-peers" help:"use these dht peers as the sole bootstrap nodes instead of the public network" hidden:"true"`
	Peer        []string      `flag:"" name:"peer" help:"connect directly to these torrent peer(s) (host:port)" hidden:"true"`
}

func (t cmdDiscoveryIdentify) Run(gctx *cmdopts.Global, tls *cmdopts.TLSConfig, id *cmdopts.SSHID) (err error) {
	signer, err := id.Signer()
	if err != nil {
		return errorsx.Wrap(err, "failed to create signer")
	}

	c := authn.AutoOauth2Client(gctx.Context, tls.Config(), authn.EndpointSSHAuth(fmt.Sprintf("https://%s", t.Endpoint)), authn.SSHTokenSourceOptionSigner(signer))
	cc := authn.AuthzClientLibrary(tls.Config(), c, t.Endpoint)

	disc, err := t.lookup(gctx, cc)
	if err != nil {
		return err
	}

	dhts, tclient, ttstore, err := t.torrentClient()
	if err != nil {
		return errorsx.Wrap(err, "unable to setup torrent client")
	}
	defer tclient.Close()

	tuners, err := t.peerTuners(gctx)
	if err != nil {
		return errorsx.Wrap(err, "unable to resolve peers")
	}

	result, err := ddisc.IdentifyOne(gctx.Context, dhts, tclient, ttstore, t.PeerTimeout, t.InfoTimeout, disc, tuners...)
	if err != nil {
		return errorsx.Wrap(err, "unable to identify media")
	}

	media, err := t.persist(gctx, cc, result)
	if err != nil {
		return err
	}

	if err = t.cleanup(gctx, cc); err != nil {
		return err
	}

	log.Println("media identified", spew.Sdump(media))

	return nil
}

func (t cmdDiscoveryIdentify) lookup(gctx *cmdopts.Global, cc *http.Client) (_ ddisc.Discovered, err error) {
	var (
		encoded url.Values
		req     *http.Request
		resp    *http.Response
		result  ddiscapi.DiscoverySearchResponse
	)

	if encoded, err = formx.NewEncoder().Encode(&ddiscapi.DiscoverySearchRequest{
		Id:    []string{t.ID},
		Limit: 1,
	}); err != nil {
		return ddisc.Discovered{}, errorsx.Wrap(err, "unable to encode request")
	}

	if req, err = http.NewRequestWithContext(gctx.Context, http.MethodGet, fmt.Sprintf("https://%s/ddisc/discovery/?"+encoded.Encode(), t.Endpoint), nil); err != nil {
		return ddisc.Discovered{}, errorsx.Wrap(err, "unable to create http request")
	}

	if resp, err = httpx.AsError(cc.Do(req)); err != nil {
		return ddisc.Discovered{}, errorsx.Wrap(err, "http request failed")
	}

	if err = httpx.DecodeJSON(resp, &result); err != nil {
		return ddisc.Discovered{}, errorsx.Wrap(err, "unable to decode response")
	}

	if len(result.Items) != 1 {
		return ddisc.Discovered{}, errorsx.Errorf("expected exactly one discovery entry for id %s, found %d", t.ID, len(result.Items))
	}

	infohash := int160.FromBytesOrZero(result.Items[0].GetInfohash())

	return ddisc.NewDiscovered(&infohash), nil
}

func (t cmdDiscoveryIdentify) peerTuners(gctx *cmdopts.Global) (_ []torrent.Tuner, err error) {
	if len(t.Peer) == 0 {
		return nil, nil
	}

	peers := make([]torrent.Peer, 0, len(t.Peer))
	for _, p := range t.Peer {
		host, port, err := net.SplitHostPort(p)
		if err != nil {
			return nil, errorsx.Wrapf(err, "invalid peer address %s", p)
		}

		addrs, err := net.DefaultResolver.LookupIP(gctx.Context, "ip4", host)
		if err != nil {
			return nil, errorsx.Wrapf(err, "unable to resolve peer %s", p)
		}

		portn, err := strconv.Atoi(port)
		if err != nil {
			return nil, errorsx.Wrapf(err, "invalid peer port %s", p)
		}

		peers = append(peers, torrent.NewPeerDeprecated(int160.Zero(), addrs[0], uint16(portn), torrent.PeerOptionTrusted(true)))
	}

	return []torrent.Tuner{torrent.TunePeers(peers...)}, nil
}

func (t cmdDiscoveryIdentify) torrentClient() (dhts *dht.Server, tclient *torrent.Client, ttstore storage.ClientImpl, err error) {
	cachedir := userx.DefaultCacheDirectory("torrentddisc")

	dhtOptions := []dht.Option{dht.OptionMuxer(dht.DefaultMuxer())}
	if len(t.DHTPeers) > 0 {
		addrs := make([]dht.Addr, 0, len(t.DHTPeers))
		for _, p := range t.DHTPeers {
			addrport, perr := netip.ParseAddrPort(p)
			if perr != nil {
				return nil, nil, nil, errorsx.Wrapf(perr, "invalid dht peer address %s", p)
			}
			addrs = append(addrs, dht.NewAddr(addrport))
		}
		dhtOptions = append(dhtOptions, dht.OptionBootstrapNodesNone, dht.OptionBootstrapFixedAddrs(addrs...))
	}

	if dhts, err = dht.NewServer(32, dhtOptions...); err != nil {
		return nil, nil, nil, errorsx.Wrap(err, "failed to initialize dht")
	}

	tnetwork, err := torrentx.Autosocket(dhts, 0, retronetx.NewConnUnlimited())
	if err != nil {
		return nil, nil, nil, errorsx.Wrap(err, "unable to setup torrent socket")
	}

	ttstore = storage.NewFile(cachedir)

	torconfig := torrent.NewDefaultClientConfig(
		torrent.NewMetadataCache(cachedir),
		ttstore,
		torrent.ClientConfigCacheDirectory(cachedir),
		torrent.ClientConfigSeed(false),
		torrent.ClientConfigPEX(false),
		torrent.ClientConfigHTTPUserAgent("retrovibed/0.0"),
	)

	if tclient, err = tnetwork.Bind(torrent.NewClient(torconfig)); err != nil {
		return nil, nil, nil, errorsx.Wrap(err, "unable to bind torrent to socket")
	}

	return dhts, tclient, ttstore, nil
}

func (t cmdDiscoveryIdentify) persist(gctx *cmdopts.Global, cc *http.Client, result ddisc.Discovered) (_ *ddiscapi.Media, err error) {
	var (
		encoded []byte
		req     *http.Request
		resp    *http.Response
		mrsp    ddiscapi.MediaCreateResponse
	)

	if encoded, err = json.Marshal(&ddiscapi.MediaCreateRequest{
		Media: ddiscapi.NewMediaFromDiscovered(result),
	}); err != nil {
		return nil, errorsx.Wrap(err, "unable to encode request")
	}

	if req, err = http.NewRequestWithContext(gctx.Context, http.MethodPost, fmt.Sprintf("https://%s/ddisc/media/", t.Endpoint), bytes.NewReader(encoded)); err != nil {
		return nil, errorsx.Wrap(err, "unable to create http request")
	}

	if resp, err = httpx.AsError(cc.Do(req)); err != nil {
		return nil, errorsx.Wrap(err, "http request failed")
	}

	if err = httpx.DecodeJSON(resp, &mrsp); err != nil {
		return nil, errorsx.Wrap(err, "unable to decode response")
	}

	return mrsp.Media, nil
}

func (t cmdDiscoveryIdentify) cleanup(gctx *cmdopts.Global, cc *http.Client) (err error) {
	var (
		req  *http.Request
		resp *http.Response
	)

	if req, err = http.NewRequestWithContext(gctx.Context, http.MethodDelete, fmt.Sprintf("https://%s/ddisc/discovery/%s", t.Endpoint, t.ID), bytes.NewReader(nil)); err != nil {
		return errorsx.Wrap(err, "unable to create http request")
	}

	if resp, err = httpx.AsError(cc.Do(req)); err != nil {
		return errorsx.Wrap(err, "http request failed")
	}

	return errorsx.Wrap(httpx.DecodeJSON(resp, &ddiscapi.DiscoveryDeleteResponse{}), "unable to decode response")
}
