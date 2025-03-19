package main

import (
	"github.com/james-lawrence/deeppool/cmd/metaidentity"

	_ "github.com/marcboeker/go-duckdb"
)

type cmdMetaIdentity struct {
	Bootstrap metaidentity.Bootstrap `cmd:"" help:"bootstrap your identity using an existing ssh key"`
	Show      metaidentity.Identity  `cmd:"" help:"display current identity"`
}
