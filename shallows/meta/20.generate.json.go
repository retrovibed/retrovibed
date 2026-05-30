package meta

//go:generate gomodifytags -w -all -quiet -skip-unexported -add-tags json -file genieql.gen.go
//go:generate gomodifytags -w -quiet -struct Wireguard -field MaximumConnections -add-options json=string -file genieql.gen.go
