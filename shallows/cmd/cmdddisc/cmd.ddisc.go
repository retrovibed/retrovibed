package cmdddisc

// peer command examples
// go -C shallows run ./cmd/retrovibe/... discovery peers create --insecure --library="eg:9998" --name="derp" --peer="34363564353033612d643263352d363338332d30" --partition="033292b1-98c2-5e96-38a4-956548a40b55"
// go -C shallows run ./cmd/retrovibe/... discovery peers delete --insecure --library="eg:9998" --peer="34363564353033612d643263352d363338332d30"
// go -C shallows run ./cmd/retrovibe/... discovery peers list --insecure --library="eg:9998"
type peer struct {
	Create cmdPeerCreate `cmd:"" help:"peer with another library for discovery"`
	Delete cmdPeerDelete `cmd:"" help:"remove a peered library"`
	List   cmdPeerList   `cmd:"" help:"list peered libraries"`
}

// discovery command examples
// go -C shallows run ./cmd/retrovibe/... discovery discovery ls --insecure --library="eg:9998"
// go -C shallows run ./cmd/retrovibe/... discovery discovery create --insecure --library="eg:9998" --infohash="<hex infohash>"
// go -C shallows run ./cmd/retrovibe/... discovery discovery delete --insecure --library="eg:9998" --id="<id from ls>"
// go -C shallows run ./cmd/retrovibe/... discovery discovery identify --insecure --library="eg:9998" --id="<id from ls>"
type discovery struct {
	Ls       cmdDiscoveryList     `cmd:"" help:"list infohashes currently being investigated by discovery"`
	Create   cmdDiscoveryCreate   `cmd:"" help:"start tracking an infohash for discovery"`
	Delete   cmdDiscoveryDelete   `cmd:"" help:"stop tracking an infohash being investigated"`
	Identify cmdDiscoveryIdentify `cmd:"" help:"identify a tracked infohash right now and record it as discovered media"`
}

// media command examples
// go -C shallows run ./cmd/retrovibe/... discovery media ls --insecure --library="eg:9998"
// go -C shallows run ./cmd/retrovibe/... discovery media create --insecure --library="eg:9998" --infohash="<hex infohash>" --title="derp"
// go -C shallows run ./cmd/retrovibe/... discovery media delete --insecure --library="eg:9998" --id="<id from ls>"
// go -C shallows run ./cmd/retrovibe/... discovery media query "<known media id>" "127.0.0.1:3196"
type media struct {
	Ls     cmdMediaLs     `cmd:"" help:"list/search discovered media on a library"`
	Create cmdMediaCreate `cmd:"" help:"create a discovered media record on a library"`
	Delete cmdMediaDelete `cmd:"" help:"remove a discovered media record from a library"`
	Query  cmdMediaQuery  `cmd:"" help:"query the DHT directly for media matching a known-media id"`
}

type Commands struct {
	Peers     peer      `cmd:"" help:"commands for managing library peering"`
	Discovery discovery `cmd:"" help:"commands for managing infohashes currently being investigated"`
	Media     media     `cmd:"" help:"commands for managing discovered media records"`
}
