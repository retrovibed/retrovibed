package retrovibed

// meta
//go:generate protoc --proto_path=../.proto --go_opt=Mmeta.authz.proto=github.com/retrovibed/retrovibed/metaapi --go_opt=paths=source_relative --go_out=metaapi meta.authz.proto
//go:generate protoc --proto_path=../.proto --go_opt=Mmeta.profile.proto=github.com/retrovibed/retrovibed/metaapi --go_opt=paths=source_relative --go_out=metaapi meta.profile.proto
//go:generate protoc --proto_path=../.proto --go_opt=Mmeta.daemon.proto=github.com/retrovibed/retrovibed/metaapi --go_opt=paths=source_relative --go_out=metaapi meta.daemon.proto

// media
//go:generate protoc --proto_path=../.proto --go_opt=Mmedia.proto=github.com/retrovibed/retrovibed/media --go_opt=paths=source_relative --go_out=media media.proto
//go:generate protoc --proto_path=../.proto --go_opt=Mrss.proto=github.com/retrovibed/retrovibed/rss --go_opt=paths=source_relative --go_out=rss rss.proto
