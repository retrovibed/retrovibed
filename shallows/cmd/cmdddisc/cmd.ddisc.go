package cmdddisc

// peer command examples
// go -C shallows run ./cmd/retrovibe/... discovery peers create --insecure --library="eg:9998" --name="derp" --peer="34363564353033612d643263352d363338332d30" --partition="033292b1-98c2-5e96-38a4-956548a40b55"
// go -C shallows run ./cmd/retrovibe/... discovery peers delete --insecure --library="eg:9998" --peer="34363564353033612d643263352d363338332d30"
type peer struct {
	Create cmdPeerCreate `cmd:"" help:"peer with another library for discovery"`
	Delete cmdPeerDelete `cmd:"" help:"remove a peered library"`
}

type Commands struct {
	Peers peer `cmd:"" help:"commands for managing library peering"`
}
