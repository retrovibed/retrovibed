package tracking

//go:generate gomodifytags -w -all -quiet -skip-unexported -add-tags json -file genieql.gen.go
//go:generate gomodifytags -w -quiet -struct Metadata -field Bytes -add-options json=string -file genieql.gen.go
//go:generate gomodifytags -w -quiet -struct Metadata -field Downloaded -add-options json=string -file genieql.gen.go
//go:generate gomodifytags -w -quiet -struct Metadata -field Uploaded -add-options json=string -file genieql.gen.go
//go:generate gomodifytags -w -quiet -struct UnknownHash -field Attempts -add-options json=string -file genieql.gen.go
//go:generate gomodifytags -w -quiet -struct Peer -field Bep51Available -add-options json=string -file genieql.gen.go
