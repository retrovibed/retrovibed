package deeppool

//go:generate protoc --proto_path=../.proto --go_opt=Mmedia.proto=github.com/james-lawrence/deeppool/media --go_opt=paths=source_relative --go_out=media media.proto
//go:generate protoc --proto_path=../.proto --go_opt=Mrss.proto=github.com/james-lawrence/deeppool/rss --go_opt=paths=source_relative --go_out=rss rss.proto
