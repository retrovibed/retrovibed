package deeppool

//go:generate protoc --proto_path=../.proto --go_opt=Mmedia.proto=github.com/james-lawrence/deeppool/media --go_opt=paths=source_relative --go_out=media media.proto
