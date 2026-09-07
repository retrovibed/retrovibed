package cmdcommunity

// publisher command examples
// go -C shallows run ./cmd/retrovibe/... community publisher install ../retroapi/publishplugin-activitypub --name lemmy-movies -e LEMMY_INSTANCE=https://lemmy.ml
// go -C shallows run ./cmd/retrovibe/... community publisher link lemmy-movies lemmy-music -e LEMMY_COMMUNITY=music@lemmy.ml
// go -C shallows run ./cmd/retrovibe/... community publisher env lemmy-movies
// go -C shallows run ./cmd/retrovibe/... community publisher config lemmy-movies -e LEMMY_PASSWORD=hunter2
//
// An installed plugin is loaded by a running daemon immediately (the
// publish.d directory is watched), but only becomes selectable for a
// community once the daemon reconciles the directory into the
// plugin_publishers table, which it does at startup.
type publisher struct {
	Install cmdPublisherInstall `cmd:"" help:"compile and install a publishplugin repository into publish.d"`
	Link    cmdPublisherLink    `cmd:"" help:"install a second, independently configured copy of an installed plugin"`
	Config  cmdPublisherConfig  `cmd:"" help:"merge configuration into an installed plugin's .env"`
	Env     cmdPublisherEnv     `cmd:"" help:"print the variables an installed plugin understands"`
}
