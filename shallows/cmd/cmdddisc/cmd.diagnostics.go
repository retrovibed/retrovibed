package cmdddisc

import (
	"fmt"

	"github.com/retrovibed/retrovibed/retroapi/authn"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
)

// diagnostics command examples
// go -C shallows run ./cmd/retrovibe/... discovery diagnostics --insecure --library="eg:9998"
type cmdDiagnostics struct {
	Endpoint string `flag:"" name:"library" help:"http address for the library you want to connect to" default:"localhost:9998"`
}

func (t cmdDiagnostics) Run(gctx *cmdopts.Global, tls *cmdopts.TLSConfig, id *cmdopts.SSHID) (err error) {
	signer, err := id.Signer()
	if err != nil {
		return errorsx.Wrap(err, "failed to create signer")
	}

	c := authn.AutoOauth2Client(gctx.Context, tls.Config(), authn.EndpointSSHAuth(fmt.Sprintf("https://%s", t.Endpoint)), authn.SSHTokenSourceOptionSigner(signer))
	cc := authn.AuthzClientLibrary(tls.Config(), c, t.Endpoint)

	result, err := metaapi.DiscoveryMetrics(gctx.Context, cc, t.Endpoint)
	if err != nil {
		return err
	}

	d := result.GetDiscovery()
	_, err = fmt.Printf(
		"enabled=%t ratio=%d partitions=%d workloads=%d local_partition=%s peers=%d peers_ddisc=%d peers_bep51=%d unidentified=%d queued=%d indexing=%d offload=%d indexed=%d\n",
		d.GetEnabled(), d.GetRatio(), d.GetPartitions(), d.GetWorkloads(), d.GetLocalPartition(), d.GetPeers(), d.GetPeersDdisc(), d.GetPeersBep51(),
		d.GetUnidentified(), d.GetQueued(), d.GetIndexing(), d.GetOffload(), d.GetIndexed(),
	)
	return err
}
