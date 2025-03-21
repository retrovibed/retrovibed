package retrovibed

//go:generate protoc --proto_path=../.proto --go_opt=Mmedia.proto=github.com/retrovibed/retrovibed/media --go_opt=paths=source_relative --go_out=media media.proto
//go:generate protoc --proto_path=../.proto --go_opt=Mrss.proto=github.com/retrovibed/retrovibed/rss --go_opt=paths=source_relative --go_out=rss rss.proto
