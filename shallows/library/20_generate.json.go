package library

//go:generate gomodifytags -w -all -quiet -skip-unexported -add-tags json -file genieql.gen.go
//go:generate gomodifytags -w -quiet -struct Known -field Md5Lower -add-options json=string -file genieql.gen.go
//go:generate gomodifytags -w -quiet -struct Metadata -field Bytes -add-options json=string -file genieql.gen.go
//go:generate gomodifytags -w -quiet -struct Metadata -field DiskOffset -add-options json=string -file genieql.gen.go
//go:generate gomodifytags -w -quiet -struct Metadata -field DiskUsage -add-options json=string -file genieql.gen.go
//go:generate gomodifytags -w -quiet -struct Metadata -field QuotaUsage -add-options json=string -file genieql.gen.go
