package ddisc

//go:generate gomodifytags -w -all -quiet -skip-unexported -add-tags json -file genieql.gen.go
//go:generate gomodifytags -w -quiet -struct Discovered -field Bytes -add-options json=string -file genieql.gen.go
