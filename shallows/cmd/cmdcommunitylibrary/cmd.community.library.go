package cmdcommunitylibrary

type Library struct {
	Publish cmdPublish `cmd:"" name:"publish" help:"publish a library item to a community"`
}
